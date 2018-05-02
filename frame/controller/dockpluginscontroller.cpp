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

#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dockitemcontroller.h"
#include "dockpluginloader.h"

#include <QDebug>
#include <QDir>

#define API_VERSION "1.0"

DockPluginsController::DockPluginsController(DockItemController *itemControllerInter)
    : QObject(itemControllerInter),
      m_itemControllerInter(itemControllerInter)
{
    qApp->installEventFilter(this);

    QTimer::singleShot(2000, this, &DockPluginsController::startLoader);
}

void DockPluginsController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // check if same item added
    if (m_pluginList.contains(itemInter))
        if (m_pluginList[itemInter].contains(itemKey))
            return;

    PluginsItem *item = new PluginsItem(itemInter, itemKey);
    item->setVisible(false);

    m_pluginList[itemInter][itemKey] = item;

    emit pluginItemInserted(item);
}

void DockPluginsController::itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = pluginItemAt(itemInter, itemKey);

    Q_ASSERT(item);

    item->update();

    emit pluginItemUpdated(item);
}

void DockPluginsController::itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = pluginItemAt(itemInter, itemKey);

    if (!item)
        return;

    item->detachPluginWidget();

    emit pluginItemRemoved(item);

    m_pluginList[itemInter].remove(itemKey);

//    QTimer::singleShot(1, this, [=] { delete item; });
    // item->deleteLater();
}

//void DockPluginsController::requestRefershWindowVisible()
//{
//    for (auto list : m_pluginList.values())
//    {
//        for (auto item : list.values())
//        {
//            Q_ASSERT(item);
//            emit item->requestRefershWindowVisible();
//            return;
//        }
//    }
//}

void DockPluginsController::requestContextMenu(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = pluginItemAt(itemInter, itemKey);
    Q_ASSERT(item);

    item->showContextMenu();
}

//void DockPluginsController::requestPopupApplet(PluginsItemInterface * const itemInter, const QString &itemKey)
//{
//    PluginsItem *item = pluginItemAt(itemInter, itemKey);

//    Q_ASSERT(item);
//    item->showPopupApplet();
//}

void DockPluginsController::startLoader()
{
    DockPluginLoader *loader = new DockPluginLoader(this);

    connect(loader, &DockPluginLoader::finished, loader, &DockPluginLoader::deleteLater, Qt::QueuedConnection);
    connect(loader, &DockPluginLoader::pluginFounded, this, &DockPluginsController::loadPlugin, Qt::QueuedConnection);

    QTimer::singleShot(1, loader, [=] { loader->start(QThread::LowestPriority); });
}

void DockPluginsController::displayModeChanged()
{
    const DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    for (auto inter : m_pluginList.keys())
        inter->displayModeChanged(displayMode);
}

void DockPluginsController::positionChanged()
{
    const Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    for (auto inter : m_pluginList.keys())
        inter->positionChanged(position);
}

void DockPluginsController::loadPlugin(const QString &pluginFile)
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

    m_pluginList.insert(interface, QMap<QString, PluginsItem *>());
    qDebug() << "init plugin: " << interface->pluginName();
    interface->init(this);
    qDebug() << "init plugin finished: " << interface->pluginName();
}

bool DockPluginsController::eventFilter(QObject *o, QEvent *e)
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

PluginsItem *DockPluginsController::pluginItemAt(PluginsItemInterface * const itemInter, const QString &itemKey) const
{
    if (!m_pluginList.contains(itemInter))
        return nullptr;

    return m_pluginList[itemInter][itemKey];
}
