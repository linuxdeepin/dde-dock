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

#define API_VERSION "1.0"

SystemTraysController::SystemTraysController(QObject *parent)
    : QObject(parent)
    , m_dbusDaemonInterface(QDBusConnection::sessionBus().interface())
    , m_pluginsSetting("deepin", "dde-dock")
{
    qApp->installEventFilter(this);
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

void SystemTraysController::requestContextMenu(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    SystemTrayItem *item = pluginItemAt(itemInter, itemKey);
    Q_ASSERT(item);

    //    item->showContextMenu();
}

void SystemTraysController::requestWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide)
{
    SystemTrayItem *item = pluginItemAt(itemInter, itemKey);
    Q_ASSERT(item);

    Q_EMIT item->requestWindowAutoHide(autoHide);
}

void SystemTraysController::requestRefershWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    SystemTrayItem *item = pluginItemAt(itemInter, itemKey);
    Q_ASSERT(item);

    Q_EMIT item->requestRefershWindowVisible();
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
    if (!meta.contains("api") || meta["api"].toString() != API_VERSION)
    {
        qWarning() << "plugin api version not matched!" << pluginFile;
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
    m_pluginsSetting.beginGroup(itemInter->pluginName());
    m_pluginsSetting.setValue(key, value);
    m_pluginsSetting.endGroup();
}

const QVariant SystemTraysController::getValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback) {
    m_pluginsSetting.beginGroup(itemInter->pluginName());
    QVariant value { std::move(m_pluginsSetting.value(key, fallback)) };
    m_pluginsSetting.endGroup();
    return std::move(value);
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
