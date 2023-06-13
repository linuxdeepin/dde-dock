// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dockplugincontroller.h"
#include "docksettings.h"
#include "pluginsiteminterface.h"
#include "pluginsiteminterface_v20.h"
#include "pluginadapter.h"
#include "utils.h"

#include <DNotifySender>
#include <DSysInfo>

#include <QDebug>
#include <QDir>
#include <QMapIterator>
#include <QPluginLoader>

#define PLUGININFO "pluginInfo"

static const QStringList CompatiblePluginApiList {
    "1.1.1",
    "1.2",
    "1.2.1",
    "1.2.2",
    DOCK_PLUGIN_API_VERSION
};

class PluginInfo : public QObject
{
public:
    PluginInfo() : QObject(nullptr), m_loaded(false), m_visible(false) {}
    bool m_loaded;
    bool m_visible;
    QString m_itemKey;
};

DockPluginController::DockPluginController(PluginProxyInterface *proxyInter, QObject *parent)
    : QObject(parent)
    , m_dbusDaemonInterface(QDBusConnection::sessionBus().interface())
    // , m_dockDaemonInter(new DockInter(dockServiceName(), dockServicePath(), QDBusConnection::sessionBus(), this))
    , m_proxyInter(proxyInter)
{
    qApp->installEventFilter(this);

    refreshPluginSettings();

    connect(DockSettings::instance(), &DockSettings::quickPluginsChanged, this, &DockPluginController::onConfigChanged);
    // connect(m_dockDaemonInter, &DockInter::PluginSettingsSynced, this, &DockPluginController::refreshPluginSettings, Qt::QueuedConnection);
}

DockPluginController::~DockPluginController()
{
    for (auto inter : m_pluginsMap.keys()) {
        delete m_pluginsMap.value(inter).value("pluginloader");
        m_pluginsMap[inter]["pluginloader"] = nullptr;
        if (m_pluginsMap[inter].contains(PLUGININFO))
            m_pluginsMap[inter][PLUGININFO]->deleteLater();
        m_pluginsMap.remove(inter);
        delete inter;
        inter = nullptr;
    }
}

QList<PluginsItemInterface *> DockPluginController::plugins() const
{
    return m_pluginsMap.keys();
}

QList<PluginsItemInterface *> DockPluginController::pluginsInSetting() const
{
    // 插件有三种状态
    // 1、所有的插件，不管这个插件是否调用itemAdded方法，只要是通过dock加载的插件（换句话说，也就是在/lib/dde-dock/plugins目录下的插件）
    // 2、在1的基础上，插件自身调用了itemAdded方法的插件，例如机器上没有蓝牙设备，那么蓝牙插件就不会调用itemAdded方法，这时候就不算
    // 3、在2的基础上，由控制中心来决定那些插件是否在任务栏显示的插件
    // 此处返回的是第二种插件
    QList<PluginsItemInterface *> settingPlugins;
    QMap<PluginsItemInterface *, int> pluginSort;
    for (auto it = m_pluginsMap.begin(); it != m_pluginsMap.end(); it++) {
        PluginsItemInterface *plugin = it.key();
        qInfo() << plugin->pluginName();
        if (plugin->pluginDisplayName().isEmpty())
            continue;

        QMap<QString, QObject *> pluginMap = it.value();
        // 如果不包含PLUGININFO这个key值，肯定是未加载
        if (!pluginMap.contains(PLUGININFO))
            continue;

        PluginInfo *pluginInfo = static_cast<PluginInfo *>(pluginMap[PLUGININFO]);
        if (!pluginInfo->m_loaded)
            continue;

        // 这里只需要返回插件为可以在控制中心设置的插件
        if (!(plugin->flags() & PluginFlag::Attribute_CanSetting))
            continue;

        settingPlugins << plugin;
        pluginSort[plugin] = plugin->itemSortKey(pluginInfo->m_itemKey);
    }

    std::sort(settingPlugins.begin(), settingPlugins.end(), [ pluginSort ](PluginsItemInterface *plugin1, PluginsItemInterface *plugin2) {
        return pluginSort[plugin1] < pluginSort[plugin2];
    });

    return settingPlugins;
}

QList<PluginsItemInterface *> DockPluginController::currentPlugins() const
{
    QList<PluginsItemInterface *> loadedPlugins;

    QMap<PluginsItemInterface *, int> pluginSortMap;
    for (auto it = m_pluginsMap.begin(); it != m_pluginsMap.end(); it++) {
        QMap<QString, QObject *> objectMap = it.value();
        if (!objectMap.contains(PLUGININFO))
            continue;

        PluginInfo *pluginInfo = static_cast<PluginInfo *>(objectMap[PLUGININFO]);
        if (!pluginInfo->m_loaded)
            continue;

        PluginsItemInterface *plugin = it.key();
        loadedPlugins << plugin;
        pluginSortMap[plugin] = plugin->itemSortKey(pluginInfo->m_itemKey);
    }

    std::sort(loadedPlugins.begin(), loadedPlugins.end(), [ pluginSortMap ](PluginsItemInterface *pluginItem1, PluginsItemInterface *pluginItem2) {
        return pluginSortMap.value(pluginItem1) < pluginSortMap.value(pluginItem2);
    });
    return loadedPlugins;
}

void DockPluginController::saveValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant &value)
{
    savePluginValue(getPluginInterface(itemInter), key, value);
}

const QVariant DockPluginController::getValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant &fallback)
{
    return getPluginValue(getPluginInterface(itemInter), key, fallback);
}

void DockPluginController::removeValue(PluginsItemInterface *const itemInter, const QStringList &keyList)
{
    removePluginValue(getPluginInterface(itemInter), keyList);
}

void DockPluginController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItemInterface *pluginItem = getPluginInterface(itemInter);

    PluginAdapter *pluginAdapter = dynamic_cast<PluginAdapter *>(pluginItem);
    if (pluginAdapter) {
        // 如果该插件可以正常转换为PluginAdapter插件，表示当前插件是v20插件，为了兼容v20插件
        // 中获取ICON，因此，使用调用插件的itemWidget来截图的方式返回QIcon，所以此处传入itemKey
        pluginAdapter->setItemKey(itemKey);
    }

    // 如果是通过插件来调用m_proxyInter的
    PluginInfo *pluginInfo = nullptr;
    QMap<QString, QObject *> &interfaceData = m_pluginsMap[pluginItem];
    if (interfaceData.contains(PLUGININFO)) {
        pluginInfo = static_cast<PluginInfo *>(interfaceData[PLUGININFO]);
        // 如果插件已经加载，则无需再次加载（此处保证插件出现重复调用itemAdded的情况）
        if (pluginInfo->m_loaded)
            return;
    } else {
        pluginInfo = new PluginInfo;
        interfaceData[PLUGININFO] = pluginInfo;
    }
    pluginInfo->m_itemKey = itemKey;
    pluginInfo->m_loaded = true;

    if (pluginCanDock(pluginItem))
        addPluginItem(pluginItem, itemKey);

    Q_EMIT pluginInserted(pluginItem, itemKey);
}

void DockPluginController::itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    m_proxyInter->itemUpdate(getPluginInterface(itemInter), itemKey);
}

void DockPluginController::itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItemInterface *pluginInter = getPluginInterface(itemInter);
    // 更新字段中的isLoaded字段，表示当前没有加载
    QMap<QString, QObject *> &interfaceData = m_pluginsMap[pluginInter];
    if (interfaceData.contains(PLUGININFO)) {
        PluginInfo *pluginInfo = static_cast<PluginInfo *>(interfaceData[PLUGININFO]);
        // 将是否加载的标记修改为未加载
        pluginInfo->m_loaded = false;
    }

    removePluginItem(pluginInter, itemKey);
    Q_EMIT pluginRemoved(pluginInter);
}

void DockPluginController::requestWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide)
{
    m_proxyInter->requestWindowAutoHide(getPluginInterface(itemInter), itemKey, autoHide);
}

void DockPluginController::requestRefreshWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    m_proxyInter->requestRefreshWindowVisible(getPluginInterface(itemInter), itemKey);
}

// 请求页面显示或者隐藏，由插件内部来调用，例如在移除蓝牙插件后，如果已经弹出了蓝牙插件的面板，则隐藏面板
void DockPluginController::requestSetAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible)
{
    PluginsItemInterface *pluginInter = getPluginInterface(itemInter);
    Q_EMIT requestAppletVisible(pluginInter, itemKey, visible);
    m_proxyInter->requestSetAppletVisible(pluginInter, itemKey, visible);
}

PluginsItemInterface *DockPluginController::getPluginInterface(PluginsItemInterface * const itemInter)
{
    // 先从事先定义好的map中查找，如果没有找到，就是v23插件，直接返回当前插件的指针
    qulonglong pluginAddr = (qulonglong)itemInter;
    if (m_pluginAdapterMap.contains(pluginAddr))
        return m_pluginAdapterMap[pluginAddr];

    return itemInter;
}

void DockPluginController::addPluginItem(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // 如果这个插件都没有加载，那么此处肯定是无需新增
    if (!m_pluginsMap.contains(itemInter))
        return;

    PluginInfo *pluginInfo = nullptr;
    QMap<QString, QObject *> &interfaceData = m_pluginsMap[itemInter];
    // 此处的PLUGININFO的数据已经在前面调用的地方给填充了数据，如果没有获取到这个数据，则无需新增
    if (!interfaceData.contains(PLUGININFO))
        return;

    pluginInfo = static_cast<PluginInfo *>(interfaceData[PLUGININFO]);
    pluginInfo->m_visible = true;

    m_proxyInter->itemAdded(itemInter, itemKey);
}

void DockPluginController::removePluginItem(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    if (!m_pluginsMap.contains(itemInter))
        return;

    // 更新字段中的isLoaded字段，表示当前没有加载
    QMap<QString, QObject *> &interfaceData = m_pluginsMap[itemInter];
    if (!interfaceData.contains(PLUGININFO))
        return;

    PluginInfo *pluginInfo = static_cast<PluginInfo *>(interfaceData[PLUGININFO]);
    // 将是否在任务栏显示的标记改为不显示
    pluginInfo->m_visible = false;

    if (QWidget * popup = itemInter->itemPopupApplet(itemKey))
        popup->hide();

    m_proxyInter->itemRemoved(itemInter, itemKey);
}

QString DockPluginController::itemKey(PluginsItemInterface *itemInter) const
{
    if (!m_pluginsMap.contains(itemInter))
        return QString();

    QMap<QString, QObject *> interfaceData = m_pluginsMap[itemInter];
    if (!interfaceData.contains(PLUGININFO))
        return QString();

    PluginInfo *pluginInfo = static_cast<PluginInfo *>(interfaceData[PLUGININFO]);
    return pluginInfo->m_itemKey;
}

QJsonObject DockPluginController::metaData(PluginsItemInterface *pluginItem)
{
    if (!m_pluginsMap.contains(pluginItem))
        return QJsonObject();

    QPluginLoader *pluginLoader = qobject_cast<QPluginLoader *>(m_pluginsMap[pluginItem].value("pluginloader"));

    if (!pluginLoader)
        return QJsonObject();

    return pluginLoader->metaData().value("MetaData").toObject();
}

void DockPluginController::savePluginValue(PluginsItemInterface * const itemInter, const QString &key, const QVariant &value)
{
    // is it necessary?
    //    refreshPluginSettings();

    // save to local cache
    QJsonObject localObject = m_pluginSettingsObject.value(itemInter->pluginName()).toObject();
    localObject.insert(key, QJsonValue::fromVariant(value)); //Note: QVariant::toJsonValue() not work in Qt 5.7

    // save to daemon
    QJsonObject remoteObject, remoteObjectInter;
    remoteObjectInter.insert(key, QJsonValue::fromVariant(value)); //Note: QVariant::toJsonValue() not work in Qt 5.7
    remoteObject.insert(itemInter->pluginName(), remoteObjectInter);

    if (itemInter->type() == PluginsItemInterface::Fixed && key == "enable" && !value.toBool()) {
        int fixedPluginCount = 0;
        // 遍历FixPlugin插件个数
        for (auto it(m_pluginsMap.begin()); it != m_pluginsMap.end();) {
            if (it.key()->type() == PluginsItemInterface::Fixed) {
                fixedPluginCount++;
            }
            ++it;
        }
        // 修改插件的order值，位置为队尾
        QString name = localObject.keys().last();
        // 此次做一下判断，有可能初始数据不存在pos_*字段，会导致enable字段被修改。或者此处可以循环所有字段是否存在pos_开头的字段？
        if (name != key) {
            localObject.insert(name, QJsonValue::fromVariant(fixedPluginCount)); //Note: QVariant::toJsonValue() not work in Qt 5.7
            // daemon中同样修改
            remoteObjectInter.insert(name, QJsonValue::fromVariant(fixedPluginCount)); //Note: QVariant::toJsonValue() not work in Qt 5.7
            remoteObject.insert(itemInter->pluginName(), remoteObjectInter);
        }
    }

    m_pluginSettingsObject.insert(itemInter->pluginName(), localObject);
    DockSettings::instance()->mergePluginSettings(QJsonDocument(remoteObject).toJson(QJsonDocument::JsonFormat::Compact));
}

const QVariant DockPluginController::getPluginValue(PluginsItemInterface * const itemInter, const QString &key, const QVariant &fallback)
{
    // load from local cache
    QVariant v = m_pluginSettingsObject.value(itemInter->pluginName()).toObject().value(key).toVariant();
    if (v.isNull() || !v.isValid()) {
        v = fallback;
    }

    return v;
}

void DockPluginController::removePluginValue(PluginsItemInterface * const itemInter, const QStringList &keyList)
{
    if (keyList.isEmpty()) {
        m_pluginSettingsObject.remove(itemInter->pluginName());
    } else {
        QJsonObject localObject = m_pluginSettingsObject.value(itemInter->pluginName()).toObject();
        for (auto key : keyList) {
            localObject.remove(key);
        }
        m_pluginSettingsObject.insert(itemInter->pluginName(), localObject);
    }

    DockSettings::instance()->removePluginSettings(itemInter->pluginName(), keyList);
}

void DockPluginController::startLoadPlugin(const QStringList &dirs)
{
    QDir dir;
    for (const QString &path : dirs) {
        if (!dir.exists(path))
            continue;

        startLoader(new PluginLoader(path, this));
    }
}

bool DockPluginController::isPluginLoaded(PluginsItemInterface *itemInter)
{
    if (!m_pluginsMap.contains(itemInter))
        return false;

    QMap<QString, QObject *> pluginObject = m_pluginsMap.value(itemInter);
    if (!pluginObject.contains(PLUGININFO))
        return false;

    PluginInfo *pluginInfo = static_cast<PluginInfo *>(pluginObject.value(PLUGININFO));
    return pluginInfo->m_visible;
}

QObject *DockPluginController::pluginItemAt(PluginsItemInterface *const itemInter, const QString &itemKey) const
{
    if (!m_pluginsMap.contains(itemInter))
        return nullptr;

    return m_pluginsMap[itemInter][itemKey];
}

PluginsItemInterface *DockPluginController::pluginInterAt(const QString &itemKey)
{
    QMapIterator<PluginsItemInterface *, QMap<QString, QObject *>> it(m_pluginsMap);
    while (it.hasNext()) {
        it.next();
        if (it.value().keys().contains(itemKey)) {
            return it.key();
        }
    }

    return nullptr;
}

PluginsItemInterface *DockPluginController::pluginInterAt(QObject *destItem)
{
    QMapIterator<PluginsItemInterface *, QMap<QString, QObject *>> it(m_pluginsMap);
    while (it.hasNext()) {
        it.next();
        if (it.value().values().contains(destItem)) {
            return it.key();
        }
    }

    return nullptr;
}

void DockPluginController::startLoader(PluginLoader *loader)
{
    connect(loader, &PluginLoader::finished, loader, &PluginLoader::deleteLater, Qt::QueuedConnection);
    connect(loader, &PluginLoader::pluginFounded, this, [ = ](const QString &pluginFile) {
        QPair<QString, PluginsItemInterface *> pair;
        pair.first = pluginFile;
        pair.second = nullptr;
        m_pluginLoadMap.insert(pair, false);
    });
    connect(loader, &PluginLoader::pluginFounded, this, &DockPluginController::loadPlugin, Qt::QueuedConnection);

    int delay = Utils::SettingValue("com.deepin.dde.dock", "/com/deepin/dde/dock/", "delay-plugins-time", 0).toInt();
    QTimer::singleShot(delay, loader, [ = ] { loader->start(QThread::LowestPriority); });
}

void DockPluginController::displayModeChanged()
{
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const auto inters = m_pluginsMap.keys();

    for (auto inter : inters)
        inter->displayModeChanged(displayMode);
}

void DockPluginController::positionChanged()
{
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    const auto inters = m_pluginsMap.keys();

    for (auto inter : inters)
        inter->positionChanged(position);
}

void DockPluginController::loadPlugin(const QString &pluginFile)
{
    QPluginLoader *pluginLoader = new QPluginLoader(pluginFile, this);
    const QJsonObject &meta = pluginLoader->metaData().value("MetaData").toObject();
    const QString &pluginApi = meta.value("api").toString();
    bool pluginIsValid = true;
    if (pluginApi.isEmpty() || !CompatiblePluginApiList.contains(pluginApi)) {
        qDebug() << objectName()
                 << "plugin api version not matched! expect versions:" << CompatiblePluginApiList
                 << ", got version:" << pluginApi
                 << ", the plugin file is:" << pluginFile;

        pluginIsValid = false;
    }

    PluginsItemInterface *interface = qobject_cast<PluginsItemInterface *>(pluginLoader->instance());
    if (!interface) {
        // 如果识别当前插件失败，就认为这个插件是v20的插件，将其转换为v20插件接口
        PluginsItemInterface_V20 *interface_v20 = qobject_cast<PluginsItemInterface_V20 *>(pluginLoader->instance());
        if (interface_v20) {
            // 将v20插件接口通过适配器转换成v23的接口，方便在后面识别
            PluginAdapter *pluginAdapter = new PluginAdapter(interface_v20, pluginLoader);
            // 将适配器的地址保存到map列表中，因为适配器自己会调用itemAdded方法，转换成PluginsItemInterface类，但是实际上它
            // 对应的是PluginAdapter类，因此，这个map用于在后面的itemAdded方法中用来查找
            m_pluginAdapterMap[(qulonglong)(interface_v20)] = pluginAdapter;
            interface = pluginAdapter;
        }
    }

    if (!interface) {
        qDebug() << objectName() << "load plugin failed!!!" << pluginLoader->errorString() << pluginFile;

        pluginLoader->unload();
        pluginLoader->deleteLater();

        pluginIsValid = false;
    }

    if (!pluginIsValid) {
        for (auto &pair : m_pluginLoadMap.keys()) {
            if (pair.first == pluginFile) {
                m_pluginLoadMap.remove(pair);
            }
        }
        QString notifyMessage(tr("The plugin %1 is not compatible with the system."));
        Dtk::Core::DUtil::DNotifySender(notifyMessage.arg(QFileInfo(pluginFile).fileName())).appIcon("dialog-warning").call();
        return;
    }

    if (interface->pluginName() == "multitasking" && (Utils::IS_WAYLAND_DISPLAY || Dtk::Core::DSysInfo::deepinType() == Dtk::Core::DSysInfo::DeepinServer)) {
        for (auto &pair : m_pluginLoadMap.keys()) {
            if (pair.first == pluginFile) {
                m_pluginLoadMap.remove(pair);
            }
        }
        return;
    }

    QMapIterator<QPair<QString, PluginsItemInterface *>, bool> it(m_pluginLoadMap);
    while (it.hasNext()) {
        it.next();
        if (it.key().first == pluginFile) {
            m_pluginLoadMap.remove(it.key());
            QPair<QString, PluginsItemInterface *> newPair;
            newPair.first = pluginFile;
            newPair.second = interface;
            m_pluginLoadMap.insert(newPair, false);
            break;
        }
    }

    // 保存 PluginLoader 对象指针
    QMap<QString, QObject *> interfaceData;
    interfaceData["pluginloader"] = pluginLoader;
    m_pluginsMap.insert(interface, interfaceData);
    QString dbusService = meta.value("depends-daemon-dbus-service").toString();
    if (!dbusService.isEmpty() && !m_dbusDaemonInterface->isServiceRegistered(dbusService).value()) {
        qDebug() << objectName() << dbusService << "daemon has not started, waiting for signal";
        connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this,
                [ = ](const QString & name, const QString & oldOwner, const QString & newOwner) {
            Q_UNUSED(oldOwner);
            if (name == dbusService && !newOwner.isEmpty()) {
                qDebug() << objectName() << dbusService << "daemon started, init plugin and disconnect";
                initPlugin(interface);
                disconnect(m_dbusDaemonInterface);
            }
        }
        );
        return;
    }

    // NOTE(justforlxz): 插件的所有初始化工作都在init函数中进行，
    // loadPlugin函数是按队列执行的，initPlugin函数会有可能导致
    // 函数执行被阻塞。
    QTimer::singleShot(1, this, [ = ] {
        initPlugin(interface);
    });
}

void DockPluginController::initPlugin(PluginsItemInterface *interface)
{
    if (!interface)
        return;

    qDebug() << objectName() << "init plugin: " << interface->pluginName();
    interface->init(this);

    for (const auto &pair : m_pluginLoadMap.keys()) {
        if (pair.second == interface)
            m_pluginLoadMap.insert(pair, true);
    }

    bool loaded = true;
    for (int i = 0; i < m_pluginLoadMap.keys().size(); ++i) {
        if (!m_pluginLoadMap.values()[i]) {
            loaded = false;
            break;
        }
    }

    // 插件全部加载完成
    if (loaded) {
        emit pluginLoadFinished();
    }
    qDebug() << objectName() << "init plugin finished: " << interface->pluginName();
}

void DockPluginController::refreshPluginSettings()
{
    const QString &pluginSettings = DockSettings::instance()->getPluginSettings();
    if (pluginSettings.isEmpty()) {
        qDebug() << "Error! get plugin settings from dbus failed!";
        return;
    }

    const QJsonObject &pluginSettingsObject = QJsonDocument::fromJson(pluginSettings.toLocal8Bit()).object();
    if (pluginSettingsObject.isEmpty()) {
        return;
    }

    // nothing changed
    if (pluginSettingsObject == m_pluginSettingsObject) {
        return;
    }

    for (auto pluginsIt = pluginSettingsObject.constBegin(); pluginsIt != pluginSettingsObject.constEnd(); ++pluginsIt) {
        const QString &pluginName = pluginsIt.key();
        const QJsonObject &settingsObject = pluginsIt.value().toObject();
        QJsonObject newSettingsObject = m_pluginSettingsObject.value(pluginName).toObject();
        for (auto settingsIt = settingsObject.constBegin(); settingsIt != settingsObject.constEnd(); ++settingsIt) {
            newSettingsObject.insert(settingsIt.key(), settingsIt.value());
        }
        // TODO: remove not exists key-values
        m_pluginSettingsObject.insert(pluginName, newSettingsObject);
    }

    // not notify plugins to refresh settings if this update is not emit by dock daemon
    // if (sender() != m_dockDaemonInter) {
    //     return;
    // }

    // notify all plugins to reload plugin settings
    for (PluginsItemInterface *pluginInter : m_pluginsMap.keys()) {
        pluginInter->pluginSettingsChanged();
    }

    // reload all plugin items for sort order or container
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> pluginsMapTemp = m_pluginsMap;
    for (auto it = pluginsMapTemp.constBegin(); it != pluginsMapTemp.constEnd(); ++it) {
        const QList<QString> &itemKeyList = it.value().keys();
        for (auto key : itemKeyList) {
            if (key != "pluginloader") {
                itemRemoved(it.key(), key);
            }
        }
        for (auto key : itemKeyList) {
            if (key != "pluginloader") {
                itemAdded(it.key(), key);
            }
        }
    }
}

bool DockPluginController::eventFilter(QObject *object, QEvent *event)
{
    if (object != qApp || event->type() != QEvent::DynamicPropertyChange)
        return false;

    QDynamicPropertyChangeEvent *const dpce = static_cast<QDynamicPropertyChangeEvent *>(event);
    const QString propertyName = dpce->propertyName();

    if (propertyName == PROP_POSITION)
        positionChanged();
    else if (propertyName == PROP_DISPLAY_MODE)
        displayModeChanged();

    return false;
}

bool DockPluginController::pluginCanDock(PluginsItemInterface *plugin) const
{
    const QStringList configPlugins = DockSettings::instance()->getQuickPlugins();
    return pluginCanDock(configPlugins, plugin);
}

bool DockPluginController::pluginCanDock(const QStringList &config, PluginsItemInterface *plugin) const
{
    // 1、如果插件是强制驻留任务栏，则始终显示
    // 2、如果插件是托盘插件，例如U盘插件，则始终显示
    if ((plugin->flags() & PluginFlag::Attribute_ForceDock)
            || (plugin->flags() & PluginFlag::Type_Tray))
        return true;

    // 3、如果该插件并未加载（未调用itemAdde或已经调用itemRemoved)，则该插件不显示
    if (!m_pluginsMap.contains(plugin))
        return false;

    const QMap<QString, QObject *> &pluginMap = m_pluginsMap[plugin];
    // 如果不包含PLUGININFO，说明该插件从未调用itemAdded方法，无需加载
    if (!pluginMap.contains(PLUGININFO))
        return false;

    // 如果该插件信息的m_loaded为true,说明已经调用过itemAdded方法，并且之后又调用了itemRemoved方法,则插件也无需加载
    PluginInfo *pluginInfo = static_cast<PluginInfo *>(pluginMap[PLUGININFO]);
    if (!pluginInfo->m_loaded)
        return false;

    // 4、插件已经驻留在任务栏，则始终显示
    return config.contains(plugin->pluginName());
}

void DockPluginController::updateDockInfo(PluginsItemInterface * const itemInter, const DockPart &part)
{
    m_proxyInter->updateDockInfo(itemInter, part);
    Q_EMIT pluginUpdated(itemInter, part);
}

void DockPluginController::onConfigChanged(const QStringList &pluginNames)
{
    // 这里只处理工具插件(回收站)和系统插件(电源插件)
    for (PluginsItemInterface *plugin : plugins()) {
        QString itemKey = this->itemKey(plugin);
        bool canDock = pluginCanDock(pluginNames, plugin);
        if (!canDock && isPluginLoaded(plugin)) {
            // 如果当前配置中不包含当前插件，但是当前插件已经加载，那么就移除该插件
            removePluginItem(plugin, itemKey);
            QWidget *itemWidget = plugin->itemWidget(itemKey);
            if (itemWidget)
                itemWidget->setVisible(false);
        } else if (canDock && !isPluginLoaded(plugin)) {
            // 如果当前配置中包含当前插件，但是当前插件并未加载，那么就加载该插件
            if (!pluginNames.contains(plugin->pluginName())) {
                // deepin-screen-recorder has Attribute_ForceDock flag
                // FIX https://github.com/linuxdeepin/developer-center/issues/4215
                continue;
            }
            addPluginItem(plugin, itemKey);
            // 工具|固定区域 插件是通过QWidget的方式进行显示的
            if (plugin->flags() & (PluginFlag::Type_Tool | PluginFlag::Type_Fixed)) {
                QWidget *itemWidget = plugin->itemWidget(itemKey);
                if (itemWidget)
                    itemWidget->setVisible(true);
            }
        }
    }
}
