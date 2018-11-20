/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             listenerri <listenerri@gmail.com>
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

#include "trayplugin.h"
#include "fashiontrayitem.h"
#include "snitraywidget.h"

#include <QDir>
#include <QWindow>
#include <QWidget>
#include <QX11Info>

#include "../widgets/tipswidget.h"
#include "xcb/xcb_icccm.h"

#define FASHION_MODE_ITEM   "fashion-mode-item"

TrayPlugin::TrayPlugin(QObject *parent)
    : QObject(parent),
      m_trayInter(new DBusTrayManager(this)),
      m_systemTraysController(new SystemTraysController(this)),
      m_tipsLabel(new TipsWidget)
{
    m_fashionItem = new FashionTrayItem(position());

    m_tipsLabel->setObjectName("tray");
    m_tipsLabel->setText(tr("System Tray"));
    m_tipsLabel->setVisible(false);
}

const QString TrayPlugin::pluginName() const
{
    return "tray";
}

void TrayPlugin::init(PluginProxyInterface *proxyInter)
{
    // transfex config
    QSettings settings("deepin", "dde-dock-shutdown");
    if (QFile::exists(settings.fileName())) {
        proxyInter->saveValue(this, "enable", settings.value("enable", true));

        QFile::remove(settings.fileName());
    }

    if (!proxyInter->getValue(this, "enable", true).toBool()) {
        qDebug() << "hide tray from config disable!!";
        return;
    }

    m_proxyInter = proxyInter;

    m_sniWatcher = new StatusNotifierWatcher(this);
    QDBusConnection dbusConn = QDBusConnection::sessionBus();
    const QString &host = QString("org.kde.StatusNotifierHost-") + QString::number(qApp->applicationPid());
    dbusConn.registerService(host);
    dbusConn.registerObject("/StatusNotifierHost", this);
    m_sniWatcher->RegisterStatusNotifierHost(host);

    connect(m_sniWatcher, &StatusNotifierWatcher::StatusNotifierItemRegistered, this, &TrayPlugin::sniItemsChanged);
    connect(m_sniWatcher, &StatusNotifierWatcher::StatusNotifierItemUnregistered, this, &TrayPlugin::sniItemsChanged);

    connect(m_trayInter, &DBusTrayManager::TrayIconsChanged, this, &TrayPlugin::trayListChanged);
    connect(m_trayInter, &DBusTrayManager::Changed, this, &TrayPlugin::trayChanged);

    connect(m_systemTraysController, &SystemTraysController::systemTrayAdded, this, &TrayPlugin::addTrayWidget);
    connect(m_systemTraysController, &SystemTraysController::systemTrayRemoved, this, &TrayPlugin::trayRemoved);

    m_trayInter->Manage();

    switchToMode(displayMode());

    QTimer::singleShot(0, this, &TrayPlugin::trayListChanged);
    QTimer::singleShot(0, this, &TrayPlugin::loadIndicator);
    QTimer::singleShot(0, this, &TrayPlugin::sniItemsChanged);
    QTimer::singleShot(0, m_systemTraysController, &SystemTraysController::startLoader);
}

void TrayPlugin::displayModeChanged(const Dock::DisplayMode mode)
{
    if (!m_proxyInter->getValue(this, "enable", true).toBool()) {
        return;
    }

    switchToMode(mode);
}

void TrayPlugin::positionChanged(const Dock::Position position)
{
    m_fashionItem->setDockPostion(position);
}

QWidget *TrayPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == FASHION_MODE_ITEM) {
        return m_fashionItem;
    }

    return m_trayMap.value(itemKey);
}

QWidget *TrayPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return nullptr;
}

QWidget *TrayPlugin::itemPopupApplet(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return nullptr;
}

bool TrayPlugin::itemAllowContainer(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    AbstractTrayWidget * const trayWidget = m_trayMap.value(itemKey);

    if (trayWidget && trayWidget->trayTyep() == AbstractTrayWidget::TrayType::SystemTray) {
        return false;
    }

    return true;
}

bool TrayPlugin::itemIsInContainer(const QString &itemKey)
{
    const QString widKey = getWindowClass(XWindowTrayWidget::toWinId(itemKey));

    return m_proxyInter
        ->getValue(this, widKey.isEmpty() ? itemKey : widKey, false)
        .toBool();
}

int TrayPlugin::itemSortKey(const QString &itemKey)
{
    Dock::DisplayMode mode = displayMode();
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(mode);

    return m_proxyInter
        ->getValue(this, key, mode == Dock::DisplayMode::Fashion ? 3 : 1)
        .toInt();
}

void TrayPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    m_proxyInter->saveValue(this, key, order);
}

void TrayPlugin::setItemIsInContainer(const QString &itemKey, const bool container)
{
    const QString widKey = getWindowClass(XWindowTrayWidget::toWinId(itemKey));

    m_proxyInter->saveValue(this, widKey.isEmpty() ? itemKey : widKey, container);
}

const QString TrayPlugin::getWindowClass(quint32 winId)
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

void TrayPlugin::sniItemsChanged()
{
    const QStringList &itemServicePaths = m_sniWatcher->RegisteredStatusNotifierItems();
    QStringList sinTrayKeyList;

    for (auto item : itemServicePaths) {
        sinTrayKeyList << SNITrayWidget::toSNIKey(item);
    }
    for (auto itemKey : m_trayMap.keys())
        if (!sinTrayKeyList.contains(itemKey) && SNITrayWidget::isSNIKey(itemKey)) {
            trayRemoved(itemKey);
        }

    for (auto tray : sinTrayKeyList) {
        trayAdded(tray);
    }
}

void TrayPlugin::trayListChanged()
{
    QList<quint32> winidList = m_trayInter->trayIcons();
    QStringList trayList;

    for (auto winid : winidList) {
        trayList << XWindowTrayWidget::toTrayWidgetId(winid);
    }

    for (auto tray : m_trayMap.keys())
        if (!trayList.contains(tray) && XWindowTrayWidget::isWinIdKey(tray)) {
            trayRemoved(tray);
        }

    for (auto tray : trayList) {
        trayAdded(tray);
    }
}

void TrayPlugin::addTrayWidget(const QString &itemKey, AbstractTrayWidget *trayWidget)
{
    if (!trayWidget) {
        return;
    }

    if (!m_trayMap.values().contains(trayWidget)) {
        m_trayMap.insert(itemKey, trayWidget);
    }

    if (displayMode() == Dock::Efficient) {
        m_proxyInter->itemAdded(this, itemKey);
    } else {
        m_proxyInter->itemAdded(this, FASHION_MODE_ITEM);
        m_fashionItem->trayWidgetAdded(trayWidget);
    }
}

void TrayPlugin::trayAdded(const QString &itemKey)
{
    if (m_trayMap.contains(itemKey)) {
        return;
    }

    if (XWindowTrayWidget::isWinIdKey(itemKey)) {
        auto winId = XWindowTrayWidget::toWinId(itemKey);
        getWindowClass(winId);
        AbstractTrayWidget *trayWidget = new XWindowTrayWidget(winId);
        addTrayWidget(itemKey, trayWidget);
    } else if (SNITrayWidget::isSNIKey(itemKey)) {
        const QString &sniServicePath = SNITrayWidget::toSNIServicePath(itemKey);
        AbstractTrayWidget *trayWidget = new SNITrayWidget(sniServicePath);
        connect(trayWidget, &AbstractTrayWidget::iconChanged, this, &TrayPlugin::sniItemIconChanged);
        addTrayWidget(itemKey, trayWidget);
    } else if (IndicatorTrayWidget::isIndicatorKey(itemKey)) {
        IndicatorTray *trayWidget = nullptr;
        QString indicatorKey = IndicatorTrayWidget::toIndicatorId(itemKey);

        if (!m_indicatorMap.keys().contains(itemKey)) {
            trayWidget = new IndicatorTray(indicatorKey);
            m_indicatorMap[indicatorKey] = trayWidget;
        }
        else {
            trayWidget = m_indicatorMap[itemKey];
        }

        connect(trayWidget, &IndicatorTray::delayLoaded,
        trayWidget, [ = ]() {
            addTrayWidget(itemKey, trayWidget->widget());
        });

        connect(trayWidget, &IndicatorTray::removed, this, [=] {
            trayRemoved(itemKey);
            trayWidget->removeWidget();
        });
    }
}

void TrayPlugin::trayRemoved(const QString &itemKey)
{
    if (!m_trayMap.contains(itemKey)) {
        return;
    }

    AbstractTrayWidget *widget = m_trayMap.take(itemKey);

    if (displayMode() == Dock::Efficient) {
        m_proxyInter->itemRemoved(this, itemKey);
    } else {
        m_fashionItem->trayWidgetRemoved(widget);
    }

    // only delete tray object when it is a tray of applications
    // set the parent of the tray object to avoid be deconstructed by parent(DockItem/PluginsItem/TrayPluginsItem)
    if (widget->trayTyep() == AbstractTrayWidget::TrayType::SystemTray) {
        widget->setParent(nullptr);
    } else {
        widget->deleteLater();
    }
}

void TrayPlugin::trayChanged(quint32 winId)
{
    QString itemKey = XWindowTrayWidget::toTrayWidgetId(winId);
    if (!m_trayMap.contains(itemKey)) {
        return;
    }

    m_trayMap.value(itemKey)->updateIcon();
}

void TrayPlugin::sniItemIconChanged()
{
    AbstractTrayWidget *trayWidget = static_cast<AbstractTrayWidget *>(sender());
    if (!m_trayMap.values().contains(trayWidget)) {
        return;
    }
}

void TrayPlugin::switchToMode(const Dock::DisplayMode mode)
{
    if (mode == Dock::Fashion) {
        for (auto itemKey : m_trayMap.keys()) {
            m_proxyInter->itemRemoved(this, itemKey);
        }
        if (m_trayMap.isEmpty()) {
            m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM);
        } else {
            m_fashionItem->setTrayWidgets(m_trayMap.values());
            m_proxyInter->itemAdded(this, FASHION_MODE_ITEM);
        }
    } else {
        m_fashionItem->clearTrayWidgets();
        m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM);
        for (auto itemKey : m_trayMap.keys()) {
            m_proxyInter->itemAdded(this, itemKey);
        }
    }
}

void TrayPlugin::loadIndicator()
{
    QDir indicatorConfDir("/etc/dde-dock/indicator");

    for (auto fileInfo : indicatorConfDir.entryInfoList({"*.json"}, QDir::Files | QDir::NoDotAndDotDot)) {
        trayAdded(IndicatorTrayWidget::toTrayWidgetId(fileInfo.baseName()));
    }
}
