/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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

#include "systemtrayplugin.h"
#include "fashiontrayitem.h"

#include <QWindow>
#include <QWidget>
#include <QX11Info>

#include "xcb/xcb_icccm.h"

#define FASHION_MODE_ITEM   "fashion-mode-item"

SystemTrayPlugin::SystemTrayPlugin(QObject *parent)
    : QObject(parent),
      m_trayInter(new DBusTrayManager(this)),
      m_trayApplet(new TrayApplet),
      m_tipsLabel(new QLabel),

      m_containerSettings(new QSettings("deepin", "dde-dock-tray"))
{
    m_trayApplet->setObjectName("sys-tray");
    m_fashionItem = new FashionTrayItem;

    m_tipsLabel->setObjectName("sys-tray");
    m_tipsLabel->setText(tr("System Tray"));
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setStyleSheet("color:white;"
                               "padding: 0 3px;");
}

const QString SystemTrayPlugin::pluginName() const
{
    return "system-tray";
}

void SystemTrayPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    connect(m_trayInter, &DBusTrayManager::TrayIconsChanged, this, &SystemTrayPlugin::trayListChanged);
    connect(m_trayInter, &DBusTrayManager::Changed, this, &SystemTrayPlugin::trayChanged);

    m_trayInter->Manage();

    switchToMode(displayMode());

    QTimer::singleShot(1, this, &SystemTrayPlugin::trayListChanged);
}

void SystemTrayPlugin::displayModeChanged(const Dock::DisplayMode mode)
{
    switchToMode(mode);
}

QWidget *SystemTrayPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == FASHION_MODE_ITEM)
    {
        // refresh active tray
        if (!m_fashionItem->activeTray())
            m_fashionItem->setActiveTray(m_trayList.first());

        return m_fashionItem;
    }

    const quint32 trayWinId = itemKey.toUInt();

    return m_trayList[trayWinId];
}

QWidget *SystemTrayPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_tipsLabel;
}

QWidget *SystemTrayPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey != FASHION_MODE_ITEM)
        return nullptr;

    Q_ASSERT(m_trayList.size());

    updateTipsContent();

    if (m_trayList.size() > 1)
        return m_trayApplet;
    else
        return nullptr;
}

bool SystemTrayPlugin::itemAllowContainer(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return true;
}

bool SystemTrayPlugin::itemIsInContainer(const QString &itemKey)
{
    const QString widKey = getWindowClass(itemKey.toInt());
    if (widKey.isEmpty())
        return false;

    return m_containerSettings->value(widKey, false).toBool();
}

int SystemTrayPlugin::itemSortKey(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return 0;
}

void SystemTrayPlugin::setItemIsInContainer(const QString &itemKey, const bool container)
{
//    qDebug() << getWindowClass(itemKey.toInt());
    m_containerSettings->setValue(getWindowClass(itemKey.toInt()), container);
}

void SystemTrayPlugin::updateTipsContent()
{
    auto trayList = m_trayList.values();
//    trayList.removeOne(m_fashionItem->activeTray());

    m_trayApplet->clear();
    m_trayApplet->addWidgets(trayList);
}

const QString SystemTrayPlugin::getWindowClass(quint32 winId)
{
    auto *connection = QX11Info::connection();

    auto *reply = new xcb_icccm_get_wm_class_reply_t;
    auto *error = new xcb_generic_error_t;
    auto cookie = xcb_icccm_get_wm_class(connection, winId);
    auto result = xcb_icccm_get_wm_class_reply(connection, cookie, reply, &error);

    QString ret;
    if (result == 1)
    {
        ret = QString("%1-%2").arg(reply->class_name).arg(reply->instance_name);
        xcb_icccm_get_wm_class_reply_wipe(reply);
    }

    delete reply;
    delete error;

    return ret;
}

void SystemTrayPlugin::trayListChanged()
{
    QList<quint32> trayList = m_trayInter->trayIcons();

    for (auto tray : m_trayList.keys())
        if (!trayList.contains(tray))
            trayRemoved(tray);

    for (auto tray : trayList)
        trayAdded(tray);
}

void SystemTrayPlugin::trayAdded(const quint32 winId)
{
    if (m_trayList.contains(winId))
        return;

    getWindowClass(winId);

    TrayWidget *trayWidget = new TrayWidget(winId);

    m_trayList[winId] = trayWidget;

    m_fashionItem->setMouseEnable(m_trayList.size() == 1);

    if (displayMode() == Dock::Efficient)
        m_proxyInter->itemAdded(this, QString::number(winId));
    else
        m_proxyInter->itemAdded(this, FASHION_MODE_ITEM);
}

void SystemTrayPlugin::trayRemoved(const quint32 winId)
{
    if (!m_trayList.contains(winId))
        return;

    TrayWidget *widget = m_trayList[winId];
    m_proxyInter->itemRemoved(this, QString::number(winId));
    m_trayList.remove(winId);
    widget->deleteLater();

    m_fashionItem->setMouseEnable(m_trayList.size() == 1);

    if (m_trayApplet->isVisible())
        updateTipsContent();

    if (m_fashionItem->activeTray() && m_fashionItem->activeTray() != widget)
        return;

    // reset active tray
    if (m_trayList.values().isEmpty())
    {
        m_fashionItem->setActiveTray(nullptr);
        m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM);
    } else
        m_fashionItem->setActiveTray(m_trayList.values().last());
}

void SystemTrayPlugin::trayChanged(const quint32 winId)
{
    if (!m_trayList.contains(winId))
        return;

    m_trayList[winId]->updateIcon();
    m_fashionItem->setActiveTray(m_trayList[winId]);

    if (m_trayApplet->isVisible())
        updateTipsContent();
}

void SystemTrayPlugin::switchToMode(const Dock::DisplayMode mode)
{
    if (mode == Dock::Fashion)
    {
        for (auto winId : m_trayList.keys())
            m_proxyInter->itemRemoved(this, QString::number(winId));
        if (m_trayList.isEmpty())
            m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM);
        else
            m_proxyInter->itemAdded(this, FASHION_MODE_ITEM);
    }
    else
    {
        m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM);
        for (auto winId : m_trayList.keys())
            m_proxyInter->itemAdded(this, QString::number(winId));
    }
}
