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
#include "fashiontray/fashiontrayitem.h"
#include "snitraywidget.h"

#include <QDir>
#include <QWindow>
#include <QWidget>
#include <QX11Info>

#include "../widgets/tipswidget.h"
#include "xcb/xcb_icccm.h"

#define FASHION_MODE_TRAYS_SORTED   "fashion-mode-trays-sorted"

#define SNI_WATCHER_SERVICE "org.kde.StatusNotifierWatcher"
#define SNI_WATCHER_PATH "/StatusNotifierWatcher"

using org::kde::StatusNotifierWatcher;

TrayPlugin::TrayPlugin(QObject *parent)
    : QObject(parent),
      m_trayInter(new DBusTrayManager(this)),
      m_systemTraysController(new SystemTraysController(this)),
      m_dbusDaemonInterface(QDBusConnection::sessionBus().interface()),
      m_tipsLabel(new TipsWidget)
{
    m_fashionItem = new FashionTrayItem(this);

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

    m_proxyInter = proxyInter;

    if (!proxyInter->getValue(this, "enable", true).toBool()) {
        qDebug() << "hide tray from config disable!!";
        return;
    }

    // registor dock as SNI Host on dbus
    QDBusConnection dbusConn = QDBusConnection::sessionBus();
    m_sniHostService = QString("org.kde.StatusNotifierHost-") + QString::number(qApp->applicationPid());
    dbusConn.registerService(m_sniHostService);
    dbusConn.registerObject("/StatusNotifierHost", this);

    m_sniWatcher = new StatusNotifierWatcher(SNI_WATCHER_SERVICE, SNI_WATCHER_PATH, QDBusConnection::sessionBus(), this);
    if (m_sniWatcher->isValid()) {
        m_sniWatcher->RegisterStatusNotifierHost(m_sniHostService);
    } else {
        qDebug() << "Tray:" << SNI_WATCHER_SERVICE << "SNI watcher daemon is not exist for now!";
    }

    connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this, &TrayPlugin::onDbusNameOwnerChanged);

    connect(m_sniWatcher, &StatusNotifierWatcher::StatusNotifierItemRegistered, this, &TrayPlugin::sniItemsChanged);
    connect(m_sniWatcher, &StatusNotifierWatcher::StatusNotifierItemUnregistered, this, &TrayPlugin::sniItemsChanged);

    connect(m_trayInter, &DBusTrayManager::TrayIconsChanged, this, &TrayPlugin::trayListChanged);
    connect(m_trayInter, &DBusTrayManager::Changed, this, &TrayPlugin::trayChanged);

    connect(m_systemTraysController, &SystemTraysController::pluginItemAdded, this, &TrayPlugin::addTrayWidget);
    connect(m_systemTraysController, &SystemTraysController::pluginItemRemoved, this,
            [=](const QString &itemKey, AbstractTrayWidget *pluginItem) { Q_UNUSED(pluginItem); trayRemoved(itemKey); });

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
    if (!m_proxyInter->getValue(this, "enable", true).toBool()) {
        return;
    }

    m_fashionItem->setDockPosition(position);
}

QWidget *TrayPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == FASHION_MODE_ITEM_KEY) {
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
    const QString &key = "container_" + itemKey;

    return m_proxyInter->getValue(this, key, false).toBool();
}

int TrayPlugin::itemSortKey(const QString &itemKey)
{
    // 如果是系统托盘图标则调用内部插件的相应接口
    if (isSystemTrayItem(itemKey)) {
        return m_systemTraysController->systemTrayItemSortKey(itemKey);
    }

    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());

    return m_proxyInter->getValue(this, key, displayMode() == Dock::DisplayMode::Fashion ? 0 : 0).toInt();
}

void TrayPlugin::setSortKey(const QString &itemKey, const int order)
{
    if (displayMode() == Dock::DisplayMode::Fashion && !traysSortedInFashionMode()) {
        m_proxyInter->saveValue(this, FASHION_MODE_TRAYS_SORTED, true);
    }

    // 如果是系统托盘图标则调用内部插件的相应接口
    if (isSystemTrayItem(itemKey)) {
        return m_systemTraysController->setSystemTrayItemSortKey(itemKey, order);
    }

    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    m_proxyInter->saveValue(this, key, order);
}

void TrayPlugin::setItemIsInContainer(const QString &itemKey, const bool container)
{
    const QString &key = "container_" + itemKey;

    m_proxyInter->saveValue(this, key, container);
}

void TrayPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == FASHION_MODE_ITEM_KEY) {
        for (auto trayWidget : m_trayMap.values()) {
            if (trayWidget) {
                trayWidget->updateIcon();
            }
        }
        return;
    }

    AbstractTrayWidget * const trayWidget = m_trayMap.value(itemKey);
    if (trayWidget) {
        trayWidget->updateIcon();
    }
}

Dock::Position TrayPlugin::dockPosition() const
{
    return position();
}

bool TrayPlugin::traysSortedInFashionMode()
{
    return m_proxyInter->getValue(this, FASHION_MODE_TRAYS_SORTED, false).toBool();
}

void TrayPlugin::saveValue(const QString &itemKey, const QString &key, const QVariant &value)
{
    // 如果是系统托盘图标则调用内部插件的相应接口
    if (isSystemTrayItem(itemKey)) {
        return m_systemTraysController->saveValueSystemTrayItem(itemKey, key, value);
    }

    m_proxyInter->saveValue(this, key, value);
}

const QVariant TrayPlugin::getValue(const QString &itemKey, const QString &key, const QVariant &fallback)
{
    // 如果是系统托盘图标则调用内部插件的相应接口
    if (isSystemTrayItem(itemKey)) {
        return m_systemTraysController->getValueSystemTrayItem(itemKey, key, fallback);
    }

    return m_proxyInter->getValue(this, key, fallback);
}

bool TrayPlugin::isSystemTrayItem(const QString &itemKey)
{
    AbstractTrayWidget * const trayWidget = m_trayMap.value(itemKey, nullptr);

    if (trayWidget && trayWidget->trayTyep() == AbstractTrayWidget::TrayType::SystemTray) {
        return true;
    }

    return false;
}

QString TrayPlugin::itemKeyOfTrayWidget(AbstractTrayWidget *trayWidget)
{
    QString itemKey;

    if (displayMode() == Dock::DisplayMode::Fashion) {
        itemKey = FASHION_MODE_ITEM_KEY;
    } else {
        itemKey = m_trayMap.key(trayWidget);
    }

    return itemKey;
}

void TrayPlugin::sniItemsChanged()
{
    const QStringList &itemServicePaths = m_sniWatcher->registeredStatusNotifierItems();
    QStringList sinTrayKeyList;

    for (auto item : itemServicePaths) {
        sinTrayKeyList << SNITrayWidget::toSNIKey(item);
    }
    for (auto itemKey : m_trayMap.keys()) {
        if (!sinTrayKeyList.contains(itemKey) && SNITrayWidget::isSNIKey(itemKey)) {
            trayRemoved(itemKey);
        }
    }

    for (int i = 0; i < sinTrayKeyList.size(); ++i) {
        traySNIAdded(sinTrayKeyList.at(i), itemServicePaths.at(i));
    }
}

void TrayPlugin::trayListChanged()
{
    QList<quint32> winidList = m_trayInter->trayIcons();
    QStringList trayKeyList;

    for (auto winid : winidList) {
        trayKeyList << XWindowTrayWidget::toTrayWidgetId(winid);
    }

    for (auto tray : m_trayMap.keys()) {
        if (!trayKeyList.contains(tray) && XWindowTrayWidget::isXWindowKey(tray)) {
            trayRemoved(tray);
        }
    }

    for (int i = 0; i < trayKeyList.size(); ++i) {
        trayXWindowAdded(trayKeyList.at(i), winidList.at(i));
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
        m_proxyInter->itemAdded(this, FASHION_MODE_ITEM_KEY);
        m_fashionItem->trayWidgetAdded(itemKey, trayWidget);
    }

    connect(trayWidget, &AbstractTrayWidget::requestWindowAutoHide, this, &TrayPlugin::onRequestWindowAutoHide, Qt::UniqueConnection);
    connect(trayWidget, &AbstractTrayWidget::requestRefershWindowVisible, this, &TrayPlugin::onRequestRefershWindowVisible, Qt::UniqueConnection);
}

void TrayPlugin::trayXWindowAdded(const QString &itemKey, quint32 winId)
{
    if (m_trayMap.contains(itemKey) || !XWindowTrayWidget::isXWindowKey(itemKey)) {
        return;
    }

    AbstractTrayWidget *trayWidget = new XWindowTrayWidget(winId);
    addTrayWidget(itemKey, trayWidget);
}

void TrayPlugin::traySNIAdded(const QString &itemKey, const QString &sniServicePath)
{
    if (m_trayMap.contains(itemKey) || !SNITrayWidget::isSNIKey(itemKey)) {
        return;
    }

    SNITrayWidget *trayWidget = new SNITrayWidget(sniServicePath);
    if (trayWidget->status() == SNITrayWidget::ItemStatus::Passive) {
        m_trayMap.insert(itemKey, trayWidget);
    } else {
        addTrayWidget(itemKey, trayWidget);
    }

    connect(trayWidget, &SNITrayWidget::statusChanged, this, &TrayPlugin::onSNIItemStatusChanged);
}

void TrayPlugin::trayIndicatorAdded(const QString &itemKey)
{
    if (m_trayMap.contains(itemKey) || !IndicatorTrayWidget::isIndicatorKey(itemKey)) {
        return;
    }

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

void TrayPlugin::trayRemoved(const QString &itemKey, const bool deleteObject)
{
    if (!m_trayMap.contains(itemKey)) {
        return;
    }

    AbstractTrayWidget *widget = m_trayMap.value(itemKey);

    if (displayMode() == Dock::Efficient) {
        m_proxyInter->itemRemoved(this, itemKey);
    } else {
        m_fashionItem->trayWidgetRemoved(widget);
    }

    // only delete tray object when it is a tray of applications
    // set the parent of the tray object to avoid be deconstructed by parent(DockItem/PluginsItem/TrayPluginsItem)
    if (widget->trayTyep() == AbstractTrayWidget::TrayType::SystemTray) {
        widget->setParent(nullptr);
    } else if (deleteObject) {
        m_trayMap.remove(itemKey);
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

void TrayPlugin::switchToMode(const Dock::DisplayMode mode)
{
    if (mode == Dock::Fashion) {
        for (auto itemKey : m_trayMap.keys()) {
            m_proxyInter->itemRemoved(this, itemKey);
        }
        if (m_trayMap.isEmpty()) {
            m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM_KEY);
        } else {
            m_fashionItem->setTrayWidgets(m_trayMap);
            m_proxyInter->itemAdded(this, FASHION_MODE_ITEM_KEY);
        }
    } else {
        m_fashionItem->clearTrayWidgets();
        m_proxyInter->itemRemoved(this, FASHION_MODE_ITEM_KEY);
        for (auto itemKey : m_trayMap.keys()) {
            m_proxyInter->itemAdded(this, itemKey);
        }
    }
}

void TrayPlugin::onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner)
{
    Q_UNUSED(oldOwner);

    if (name == SNI_WATCHER_SERVICE && !newOwner.isEmpty()) {
        qDebug() << "Tray:" << SNI_WATCHER_SERVICE << "SNI watcher daemon started, register dock to watcher as SNI Host";
        m_sniWatcher->RegisterStatusNotifierHost(m_sniHostService);
    }
}

void TrayPlugin::onRequestWindowAutoHide(const bool autoHide)
{
    const QString &itemKey = itemKeyOfTrayWidget(static_cast<AbstractTrayWidget *>(sender()));

    if (itemKey.isEmpty()) {
        return;
    }

    m_proxyInter->requestWindowAutoHide(this, itemKey, autoHide);
}

void TrayPlugin::onRequestRefershWindowVisible()
{
    const QString &itemKey = itemKeyOfTrayWidget(static_cast<AbstractTrayWidget *>(sender()));

    if (itemKey.isEmpty()) {
        return;
    }

    m_proxyInter->requestRefreshWindowVisible(this, itemKey);
}

void TrayPlugin::onSNIItemStatusChanged(SNITrayWidget::ItemStatus status)
{
    SNITrayWidget *trayWidget = static_cast<SNITrayWidget *>(sender());
    const QString &itemKey = m_trayMap.key(trayWidget);
    if (!trayWidget || itemKey.isEmpty()) {
        return;
    }

    switch (status) {
    case SNITrayWidget::Passive: {
        trayRemoved(itemKey, false);
        break;
    }
    case SNITrayWidget::Active:
    case SNITrayWidget::NeedsAttention: {
        addTrayWidget(itemKey, trayWidget);
        break;
    }
    default:
        break;
    }
}

void TrayPlugin::loadIndicator()
{
    QDir indicatorConfDir("/etc/dde-dock/indicator");

    for (auto fileInfo : indicatorConfDir.entryInfoList({"*.json"}, QDir::Files | QDir::NoDotAndDotDot)) {
        trayIndicatorAdded(IndicatorTrayWidget::toTrayWidgetId(fileInfo.baseName()));
    }
}
