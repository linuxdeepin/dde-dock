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

#include "systemtrayplugin.h"
#include "fashiontrayitem.h"

#include <QDir>
#include <QWindow>
#include <QWidget>
#include <QX11Info>

#include "../widgets/tipswidget.h"
#include "xcb/xcb_icccm.h"

#define FASHION_MODE_ITEM   "fashion-mode-item"

SystemTrayPlugin::SystemTrayPlugin(QObject *parent)
    : QObject(parent),
      m_trayInter(new DBusTrayManager(this)),
      m_trayApplet(new TrayApplet),
      m_tipsLabel(new TipsWidget),

      m_containerSettings(new QSettings("deepin", "dde-dock-tray"))
{
    m_trayApplet->setObjectName("sys-tray");
    m_fashionItem = new FashionTrayItem;

    m_tipsLabel->setObjectName("sys-tray");
    m_tipsLabel->setText(tr("System Tray"));
    m_tipsLabel->setVisible(false);
}

const QString SystemTrayPlugin::pluginName() const
{
    return "system-tray";
}

void SystemTrayPlugin::init(PluginProxyInterface *proxyInter)
{
    if (!m_containerSettings->value("enable", true).toBool()) {
        qDebug() << "hide tray from config disable!!";
        return;
    }

    m_proxyInter = proxyInter;

    connect(m_trayInter, &DBusTrayManager::TrayIconsChanged, this, &SystemTrayPlugin::trayListChanged);
    connect(m_trayInter, &DBusTrayManager::Changed, this, &SystemTrayPlugin::trayChanged);

    m_trayInter->Manage();

    switchToMode(displayMode());

    QTimer::singleShot(1, this, &SystemTrayPlugin::trayListChanged);
    QTimer::singleShot(2, this, &SystemTrayPlugin::loadIndicator);
}

void SystemTrayPlugin::displayModeChanged(const Dock::DisplayMode mode)
{
    if (!m_containerSettings->value("enable", true).toBool()) {
        return;
    }

    switchToMode(mode);
}

QWidget *SystemTrayPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == FASHION_MODE_ITEM) {
        return m_fashionItem;
    }

    return m_trayList.value(itemKey);
}

QWidget *SystemTrayPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

#ifdef DOCK_TRAY_USE_NATIVE_POPUP
    return nullptr;
#else
    return m_tipsLabel;
#endif
}

QWidget *SystemTrayPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey != FASHION_MODE_ITEM) {
        return nullptr;
    }

    Q_ASSERT(m_trayList.size());

    updateTipsContent();

    if (m_trayList.size() > 1) {
        return m_trayApplet;
    } else {
        return nullptr;
    }
}

bool SystemTrayPlugin::itemAllowContainer(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return true;
}

bool SystemTrayPlugin::itemIsInContainer(const QString &itemKey)
{
    const QString widKey = getWindowClass(XWindowTrayWidget::toWinId(itemKey));
    if (!widKey.isEmpty())
        return m_containerSettings->value(widKey, false).toBool();
    else
        return m_containerSettings->value(itemKey, false).toBool();
}

int SystemTrayPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    return m_containerSettings->value(key, 0).toInt();
}

void SystemTrayPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    m_containerSettings->setValue(key, order);
}

void SystemTrayPlugin::setItemIsInContainer(const QString &itemKey, const bool container)
{
    const QString widKey = getWindowClass(XWindowTrayWidget::toWinId(itemKey));
    if (widKey.isEmpty())
        m_containerSettings->setValue(itemKey, container);
    else
        m_containerSettings->setValue(widKey, container);
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
    if (result == 1) {
        ret = QString("%1-%2").arg(reply->class_name).arg(reply->instance_name);
        xcb_icccm_get_wm_class_reply_wipe(reply);
    }

    delete reply;
    delete error;

    return ret;
}

void SystemTrayPlugin::trayListChanged()
{
    QList<quint32> winidList = m_trayInter->trayIcons();
    QStringList trayList;

    for (auto winid : winidList) {
        trayList << XWindowTrayWidget::toTrayWidgetId(winid);
    }

    for (auto tray : m_trayList.keys())
        if (!trayList.contains(tray) && XWindowTrayWidget::isWinIdKey(tray)) {
            trayRemoved(tray);
        }

    for (auto tray : trayList) {
        trayAdded(tray);
    }
}

void SystemTrayPlugin::trayAdded(const QString itemKey)
{
    if (m_trayList.contains(itemKey)) {
        return;
    }

    auto addTrayWidget = [ = ](AbstractTrayWidget * trayWidget) {
        if (trayWidget) {
            if (!m_trayList.values().contains(trayWidget)) {
                m_trayList.insert(itemKey, trayWidget);
            }

            m_fashionItem->setMouseEnable(m_trayList.size() == 1);
            if (!m_fashionItem->activeTray()) {
                m_fashionItem->setActiveTray(trayWidget);
            }

            if (displayMode() == Dock::Efficient) {
                m_proxyInter->itemAdded(this, itemKey);
            } else {
                m_proxyInter->itemAdded(this, FASHION_MODE_ITEM);
            }
        }
    };

    if (XWindowTrayWidget::isWinIdKey(itemKey)) {
        auto winId = XWindowTrayWidget::toWinId(itemKey);
        getWindowClass(winId);
        AbstractTrayWidget *trayWidget = new XWindowTrayWidget(winId);
        addTrayWidget(trayWidget);
    }

    if (IndicatorTrayWidget::isIndicatorKey(itemKey)) {
        IndicatorTray *trayWidget = nullptr;
        QString indicatorKey = IndicatorTrayWidget::toIndicatorId(itemKey);

        if (!m_indicatorList.keys().contains(itemKey)) {
            trayWidget = new IndicatorTray(indicatorKey);
            m_indicatorList[indicatorKey] = trayWidget;
        }
        else {
            trayWidget = m_indicatorList[itemKey];
        }

        connect(trayWidget, &IndicatorTray::delayLoaded,
        trayWidget, [ = ]() {
            addTrayWidget(trayWidget->widget());
        });

        connect(trayWidget, &IndicatorTray::removed, this, [=] {
            trayRemoved(itemKey);
            trayWidget->removeWidget();
        });
    }
}

void SystemTrayPlugin::trayRemoved(const QString itemKey)
{
    if (!m_trayList.contains(itemKey)) {
        return;
    }

    QWidget *widget = m_trayList.take(itemKey);
    m_proxyInter->itemRemoved(this, itemKey);
    widget->deleteLater();

    m_fashionItem->setMouseEnable(m_trayList.size() == 1);

    if (m_trayApplet->isVisible()) {
        updateTipsContent();
    }

    if (m_fashionItem->activeTray() && m_fashionItem->activeTray() != widget) {
        return;
    }

    // reset active tray
    if (m_trayList.values().isEmpty()) {
        m_fashionItem->setActiveTray(nullptr);
        m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM);
    } else {
        m_fashionItem->setActiveTray(m_trayList.values().last());
    }
}

void SystemTrayPlugin::trayChanged(quint32 winId)
{
    QString itemKey = XWindowTrayWidget::toTrayWidgetId(winId);
    if (!m_trayList.contains(itemKey)) {
        return;
    }

    m_trayList.value(itemKey)->updateIcon();
    m_fashionItem->setActiveTray(m_trayList.value(itemKey));

    if (m_trayApplet->isVisible()) {
        updateTipsContent();
    }
}

void SystemTrayPlugin::switchToMode(const Dock::DisplayMode mode)
{
    if (mode == Dock::Fashion) {
        for (auto itemKey : m_trayList.keys()) {
            m_proxyInter->itemRemoved(this, itemKey);
        }
        if (m_trayList.isEmpty()) {
            m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM);
        } else {
            m_proxyInter->itemAdded(this, FASHION_MODE_ITEM);
        }
    } else {
        m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM);
        for (auto itemKey : m_trayList.keys()) {
            m_proxyInter->itemAdded(this, itemKey);
        }
    }
}

void SystemTrayPlugin::loadIndicator()
{
    QDir indicatorConfDir("/etc/dde-dock/indicator");

    for (auto fileInfo : indicatorConfDir.entryInfoList({"*.json"}, QDir::Files | QDir::NoDotAndDotDot)) {
        trayAdded(IndicatorTrayWidget::toTrayWidgetId(fileInfo.baseName()));
    }
}
