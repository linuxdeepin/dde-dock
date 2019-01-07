/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "systemtrayscontroller.h"
#include "pluginsiteminterface.h"
#include "systemtrayloader.h"

#include <QDebug>
#include <QDir>

SystemTraysController::SystemTraysController(QObject *parent)
    : QObject(parent)
    , m_dbusDaemonInterface(QDBusConnection::sessionBus().interface())
    , m_dockDaemonInter(new DockDaemonInter("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this))
{
    qApp->installEventFilter(this);

    refreshPluginSettings(QDateTime::currentMSecsSinceEpoch() / 1000 / 1000);

    connect(m_dockDaemonInter, &DockDaemonInter::PluginSettingsUpdated, this, &SystemTraysController::refreshPluginSettings);
}

void SystemTraysController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // check if same item added
    if (m_pluginsMap.contains(itemInter))
        if (m_pluginsMap[itemInter].contains(itemKey))
            return;

    SystemTrayItem *item = new SystemTrayItem(itemInter, itemKey);

    item->setVisible(false);

    m_pluginsMap[itemInter][itemKey] = item;

    emit systemTrayAdded(itemKey, item);
}

void SystemTraysController::itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    SystemTrayItem *item = pluginItemAt(itemInter, itemKey);

    Q_ASSERT(item);

    item->update();

    emit systemTrayUpdated(itemKey);
}

void SystemTraysController::itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    SystemTrayItem *item = pluginItemAt(itemInter, itemKey);

    if (!item)
        return;

    item->detachPluginWidget();

    emit systemTrayRemoved(itemKey);

    m_pluginsMap[itemInter].remove(itemKey);

    // do not delete the itemWidget object(specified in the plugin interface)
    item->centralWidget()->setParent(nullptr);

    // just delete our wrapper object(PluginsItem)
    item->deleteLater();
}

void SystemTraysController::requestWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide)
{
    SystemTrayItem *item = pluginItemAt(itemInter, itemKey);
    Q_ASSERT(item);

    Q_EMIT item->requestWindowAutoHide(autoHide);
}

void SystemTraysController::requestRefreshWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    SystemTrayItem *item = pluginItemAt(itemInter, itemKey);
    Q_ASSERT(item);

    Q_EMIT item->requestRefershWindowVisible();
}

void SystemTraysController::requestSetAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible)
{
    SystemTrayItem *item = pluginItemAt(itemInter, itemKey);
    Q_ASSERT(item);

    if (visible) {
        item->showPopupApplet(itemInter->itemPopupApplet(itemKey));
    } else {
        item->hidePopup();
    }
}

void SystemTraysController::startLoader()
{
    SystemTrayLoader *loader = new SystemTrayLoader(this);

    connect(loader, &SystemTrayLoader::finished, loader, &SystemTrayLoader::deleteLater, Qt::QueuedConnection);
    connect(loader, &SystemTrayLoader::pluginFounded, this, &SystemTraysController::loadPlugin, Qt::QueuedConnection);

    QTimer::singleShot(1, loader, [=] { loader->start(QThread::LowestPriority); });
}

void SystemTraysController::displayModeChanged()
{
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const auto inters = m_pluginsMap.keys();

    for (auto inter : inters)
        inter->displayModeChanged(displayMode);
}

void SystemTraysController::positionChanged()
{
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    const auto inters = m_pluginsMap.keys();

    for (auto inter : inters)
        inter->positionChanged(position);
}

void SystemTraysController::loadPlugin(const QString &pluginFile)
{
    QPluginLoader *pluginLoader = new QPluginLoader(pluginFile);
    const auto meta = pluginLoader->metaData().value("MetaData").toObject();
    if (!meta.contains("api") || meta["api"].toString() != DOCK_PLUGIN_API_VERSION)
    {
        qWarning() << "plugin api version not matched! expect version:" << DOCK_PLUGIN_API_VERSION << pluginFile;
        return;
    }

    PluginsItemInterface *interface = qobject_cast<PluginsItemInterface *>(pluginLoader->instance());
    if (!interface)
    {
        qWarning() << "load plugin failed!!!" << pluginLoader->errorString() << pluginFile;
        pluginLoader->unload();
        pluginLoader->deleteLater();
        return;
    }

    m_pluginsMap.insert(interface, QMap<QString, SystemTrayItem *>());

    QString dbusService = meta.value("depends-daemon-dbus-service").toString();
    if (!dbusService.isEmpty() && !m_dbusDaemonInterface->isServiceRegistered(dbusService).value()) {
        qDebug() << "SystemTray:" << dbusService << "daemon has not started, waiting for signal";
        connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this,
            [=](const QString &name, const QString &oldOwner, const QString &newOwner) {
                if (name == dbusService && !newOwner.isEmpty()) {
                    qDebug() << "SystemTray:" << dbusService << "daemon started, init plugin and disconnect";
                    initPlugin(interface);
                    disconnect(m_dbusDaemonInterface);
                }
            }
        );
        return;
    }

    initPlugin(interface);
}

void SystemTraysController::initPlugin(PluginsItemInterface *interface) {
    qDebug() << "SystemTray:" << "init plugin: " << interface->pluginName();
    interface->init(this);
    qDebug() << "SystemTray:" << "init plugin finished: " << interface->pluginName();
}

void SystemTraysController::refreshPluginSettings(qlonglong ts)
{
    // TODO: handle nano seconds

    const QString &pluginSettings = m_dockDaemonInter->GetPluginSettings().value();
    if (pluginSettings.isEmpty()) {
        qDebug() << "Error! get plugin settings from dbus failed!";
        return;
    }

    const QJsonObject &settingsObject = QJsonDocument::fromJson(pluginSettings.toLocal8Bit()).object();
    if (settingsObject.isEmpty()) {
        qDebug() << "Error! parse plugin settings from json failed!";
        return;
    }

    m_pluginSettingsObject = settingsObject;

    // not notify plugins to refresh settings if this update is not emit by dock daemon
    if (sender() != m_dockDaemonInter) {
        return;
    }

//    static_cast<DockDaemonInter *>(sender())
    // TODO: notify all plugins to reload plugin settings
}

bool SystemTraysController::eventFilter(QObject *o, QEvent *e)
{
    if (o != qApp)
        return false;
    if (e->type() != QEvent::DynamicPropertyChange)
        return false;

    QDynamicPropertyChangeEvent * const dpce = static_cast<QDynamicPropertyChangeEvent *>(e);
    const QString propertyName = dpce->propertyName();

    if (propertyName == PROP_POSITION)
        positionChanged();
    else if (propertyName == PROP_DISPLAY_MODE)
        displayModeChanged();

    return false;
}

SystemTrayItem *SystemTraysController::pluginItemAt(PluginsItemInterface * const itemInter, const QString &itemKey) const
{
    if (!m_pluginsMap.contains(itemInter))
        return nullptr;

    return m_pluginsMap[itemInter][itemKey];
}

PluginsItemInterface *SystemTraysController::pluginInterAt(const QString &itemKey) const
{
    for (auto it = m_pluginsMap.constBegin(); it != m_pluginsMap.constEnd(); ++it) {
        for (auto key : it.value().keys()) {
            if (key == itemKey) {
                return it.key();
            }
        }
    }

    return nullptr;
}

PluginsItemInterface *SystemTraysController::pluginInterAt(SystemTrayItem *systemTrayItem) const
{
    for (auto it = m_pluginsMap.constBegin(); it != m_pluginsMap.constEnd(); ++it) {
        for (auto item : it.value().values()) {
            if (item == systemTrayItem) {
                return it.key();
            }
        }
    }

    return nullptr;
}

void SystemTraysController::saveValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant &value) {
    QJsonObject valueObject = m_pluginSettingsObject.value(itemInter->pluginName()).toObject();
    valueObject.insert(key, value.toJsonValue());
    m_pluginSettingsObject.insert(itemInter->pluginName(), valueObject);

    m_dockDaemonInter->SetPluginSettings(
                QDateTime::currentMSecsSinceEpoch() / 1000 / 1000,
                QJsonDocument(m_pluginSettingsObject).toJson());
}

const QVariant SystemTraysController::getValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback) {
    QVariant v = m_pluginSettingsObject.value(itemInter->pluginName()).toObject().value(key).toVariant();
    if (v.isNull() || !v.isValid()) {
        v = fallback;
    }
    return v;
}

int SystemTraysController::systemTrayItemSortKey(const QString &itemKey)
{
    auto inter = pluginInterAt(itemKey);

    if (!inter) {
        return -1;
    }

    return inter->itemSortKey(itemKey);
}

void SystemTraysController::setSystemTrayItemSortKey(const QString &itemKey, const int order)
{
    auto inter = pluginInterAt(itemKey);

    if (!inter) {
        return;
    }

    inter->setSortKey(itemKey, order);
}

const QVariant SystemTraysController::getValueSystemTrayItem(const QString &itemKey, const QString &key, const QVariant &fallback)
{
    auto inter = pluginInterAt(itemKey);

    if (!inter) {
        return QVariant();
    }

    return getValue(inter, key, fallback);
}

void SystemTraysController::saveValueSystemTrayItem(const QString &itemKey, const QString &key, const QVariant &value)
{
    auto inter = pluginInterAt(itemKey);

    if (!inter) {
        return;
    }

    saveValue(inter, key, value);
}
