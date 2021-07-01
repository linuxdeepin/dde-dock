/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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

#include "networkpanel.h"
#include "constants.h"
#include "../../widgets/tipswidget.h"
#include "utils.h"
#include <item/NetItem.h>
#include "item/devicestatushandler.h"
#include <imageutil.h>

#include <DHiDPIHelper>
#include <DApplicationHelper>
#include <DDBusSender>

#include <QTimer>
#include <QScroller>
#include <QVBoxLayout>

#include <unetworkcontroller.h>
#include <unetworkcontroller.h>
#include <unetworkdevicebase.h>
#include <uwireddevice.h>
#include <uwirelessdevice.h>

const int ItemWidth = 250;
const QString MenueEnable = "enable";
const QString MenueWiredEnable = "wireEnable";
const QString MenueWirelessEnable = "wirelessEnable";
const QString MenueSettings = "settings";

NetworkPanel::NetworkPanel(QWidget *parent)
    : QWidget(parent)
    , m_refreshIconTimer(new QTimer(this))
    , m_switchWireTimer(new QTimer(this))
    , m_wirelessScanTimer(new QTimer(this))
    , m_wirelessScanInterval(Utils::SettingValue("com.deepin.dde.dock", QByteArray(), "wireless-scan-interval", 10).toInt())
    , m_tipsWidget(new Dock::TipsWidget(this))
    , m_switchWire(true)
    , m_applet(new QScrollArea(this))
    , m_centerWidget(new QWidget(this))
    , m_netListView(new DListView(m_centerWidget))
    , m_timeOut(true)
{
    initUi();
    initConnection();
}

NetworkPanel::~NetworkPanel()
{
}

void NetworkPanel::initUi()
{
    const QGSettings *gsetting = Utils::SettingsPtr("com.deepin.dde.dock", QByteArray(), this);
    if (gsetting)
        connect(gsetting, &QGSettings::changed, [&](const QString &key) {
            if (key == "wireless-scan-interval") {
                m_wirelessScanInterval = gsetting->get("wireless-scan-interval").toInt() * 1000;
                m_wirelessScanTimer->setInterval(m_wirelessScanInterval);
            }
        });

    m_refreshIconTimer->setInterval(100);
    m_tipsWidget->setVisible(false);

    m_netListView->setAccessibleName("list_network");
    m_netListView->setBackgroundType(DStyledItemDelegate::BackgroundType::ClipCornerBackground);
    m_netListView->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
    m_netListView->setFrameShape(QFrame::NoFrame);
    m_netListView->setViewportMargins(0, 0, 0, 0);
    m_netListView->setItemSpacing(1);
    m_netListView->setVerticalScrollMode(QAbstractItemView::ScrollPerPixel);

    NetworkDelegate *delegate = new NetworkDelegate(m_netListView);
    m_netListView->setItemDelegate(delegate);

    m_model = new QStandardItemModel(this);
    m_netListView->setModel(m_model);

    QVBoxLayout *centerLayout = new QVBoxLayout(m_centerWidget);
    centerLayout->setContentsMargins(0, 0, 0, 0);
    centerLayout->addWidget(m_netListView);

    m_applet->setFixedWidth(ItemWidth);
    m_applet->setWidget(m_centerWidget);
    m_applet->setFrameShape(QFrame::NoFrame);
    m_applet->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_applet->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centerWidget->setAutoFillBackground(false);
    m_applet->viewport()->setAutoFillBackground(false);
    m_applet->setVisible(false);

    setControlBackground();
}

void NetworkPanel::initConnection()
{
    // 定期更新网络状态图标
    connect(m_refreshIconTimer, &QTimer::timeout, this, &NetworkPanel::refreshIcon);

    // 主题发生变化触发的信号
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &NetworkPanel::updatePlugView);

    // 连接信号
    UNetworkController *networkController = UNetworkController::instance();
    connect(networkController, &UNetworkController::deviceAdded, this, &NetworkPanel::onDeviceAdded);
    connect(networkController, &UNetworkController::deviceRemoved, this, &NetworkPanel::updatePlugView);
    connect(networkController, &UNetworkController::connectivityChanged, this, &NetworkPanel::updatePlugView);

    // 点击列表的信号
    connect(m_netListView, &DListView::clicked, this, &NetworkPanel::onClickListView);

    // 连接超时的信号
    connect(m_switchWireTimer, &QTimer::timeout, [ = ]() {
        m_switchWire = !m_switchWire;
        m_timeOut = true;
    });
}

void NetworkPanel::getPluginState()
{
    // 所有设备状态叠加
    QList<int> status;
    m_pluginState = DeviceStatusHandler::pluginState();
    switch (m_pluginState) {
    case PluginState::Unknow:
    case PluginState::Disabled:
    case PluginState::Connected:
    case PluginState::Disconnected:
    case PluginState::ConnectNoInternet:
    case PluginState::WirelessDisabled:
    case PluginState::WiredDisabled:
    case PluginState::WirelessConnected:
    case PluginState::WiredConnected:
    case PluginState::WirelessDisconnected:
    case PluginState::WiredDisconnected:
    case PluginState::WirelessConnecting:
    case PluginState::WiredConnecting:
    case PluginState::WirelessConnectNoInternet:
    case PluginState::WiredConnectNoInternet:
    case PluginState::WiredFailed:
    case PluginState::Nocable:
        m_switchWireTimer->stop();
        m_timeOut = true;
        break;
    case PluginState::Connecting:
        // 启动2s切换计时,只有当计时器记满则重新计数
        if (m_timeOut) {
            m_switchWireTimer->start(2000);
            m_timeOut = false;
        }
        break;
    default:
        break;
    }
}

void NetworkPanel::updateItems(QList<NetItem *> &removeItems)
{
    auto findBaseController = [ = ](UDeviceType t)->DeviceControllItem *{
        for (int i = 0; i < m_items.size(); i++) {
            NetItem *item = m_items[i];
            if (item->itemType() != NetItemType::DeviceControllViewItem)
                continue;

            DeviceControllItem *pBaseCtrlItem = static_cast<DeviceControllItem *>(item);
            if (pBaseCtrlItem->deviceType() == t)
                return pBaseCtrlItem;
        }

        return Q_NULLPTR;
    };

    auto findWiredController = [ = ](UWiredDevice *device)->WiredControllItem *{
        for (int i = 0; i < m_items.size(); i++) {
            NetItem *item = m_items[i];
            if (item->itemType() != NetItemType::WiredControllViewItem)
                continue;

            WiredControllItem *wiredCtrlItem = static_cast<WiredControllItem *>(item);
            if (wiredCtrlItem->device() == device)
                return wiredCtrlItem;
        }

        return Q_NULLPTR;
    };

    auto findWiredItem = [ = ](UWiredConnection *conn)->WiredItem *{
        for (int i = 0; i < m_items.size(); i++) {
            NetItem *item = m_items[i];
            if (item->itemType() != NetItemType::WiredViewItem)
                continue;

            WiredItem *wiredItem = static_cast<WiredItem *>(item);
            if (wiredItem->connection() == conn)
                return wiredItem;
        }

        return Q_NULLPTR;
    };

    auto findWirelessController = [ = ](UWirelessDevice *device)->WirelessControllItem *{
        for (int i = 0; i < m_items.size(); i++) {
            NetItem *item = m_items[i];
            if (item->itemType() != NetItemType::WirelessControllViewItem)
                continue;

            WirelessControllItem *wiredCtrlItem = static_cast<WirelessControllItem *>(item);
            if (wiredCtrlItem->device() == device)
                return wiredCtrlItem;
        }

        return Q_NULLPTR;
    };

    auto findWirelessItem = [ = ](const UAccessPoints *ap)->WirelessItem *{
        for (int i = 0; i < m_items.size(); i++) {
            NetItem *item = m_items[i];
            if (item->itemType() != NetItemType::WirelessViewItem)
                continue;

            WirelessItem *wirelessItem = static_cast<WirelessItem *>(item);
            const UAccessPoints *apData = wirelessItem->accessPoint();
            if (apData == ap)
                return wirelessItem;
        }

        return Q_NULLPTR;
    };

    QList<UNetworkDeviceBase *> devices = UNetworkController::instance()->devices();
    QList<UWiredDevice *> wiredDevices;
    QList<UWirelessDevice *> wirelessDevices;

    for (int i = 0; i < devices.size(); i++) {
        UNetworkDeviceBase *device = devices[i];

        if (device->deviceType() == UDeviceType::Wired) {
            UWiredDevice *dev = static_cast<UWiredDevice *>(device);
            wiredDevices << dev;
        } else if (device->deviceType() == UDeviceType::Wireless) {
            UWirelessDevice *dev = static_cast<UWirelessDevice *>(device);
            wirelessDevices << dev;
        }
    }

    // 存在多个无线设备的情况下，需要显示总开关
    QList<NetItem *> items;
    if (wirelessDevices.size() > 1) {
        DeviceControllItem *ctrl = findBaseController(UDeviceType::Wireless);
        if (!ctrl)
            ctrl = new DeviceControllItem(UDeviceType::Wireless, m_netListView->viewport());
        else
            ctrl->updateView();

        ctrl->setDevices(devices);
        items << ctrl;
    }

    // 遍历当前所有的无线网卡
    auto accessPoints = [ & ](UWirelessDevice *device) {
        if (device->isEnabled())
            return device->accessPointItems();

        return QList<UAccessPoints *>();
    };

    for (int i = 0; i < wirelessDevices.size(); i++) {
        UWirelessDevice *device = wirelessDevices[i];
        WirelessControllItem *ctrl = findWirelessController(device);
        if (!ctrl)
            ctrl = new WirelessControllItem(m_netListView->viewport(), static_cast<UWirelessDevice *>(device));
        else
            ctrl->updateView();

        items << ctrl;

        QList<UAccessPoints *> aps = accessPoints(device);
        for (int j = 0; j < aps.size(); j++) {
            UAccessPoints *ap = aps[j];
            WirelessItem *apCtrl = findWirelessItem(ap);
            if (!apCtrl)
                apCtrl = new WirelessItem(m_netListView->viewport(), device, ap);
            else
                apCtrl->updateView();

            items << apCtrl;
        }
    }

    // 存在多个有线设备的情况下，需要显示总开关
    if (wiredDevices.size() > 1) {
        DeviceControllItem *ctrl = findBaseController(UDeviceType::Wired);
        if (!ctrl)
            ctrl = new DeviceControllItem(UDeviceType::Wired, m_netListView->viewport());
        else
            ctrl->updateView();

        ctrl->setDevices(devices);
        items << ctrl;
    }

    auto wiredConnections = [ & ](UWiredDevice *device) {
        if (device->isEnabled())
            return device->items();

        return QList<UWiredConnection *>();
    };

    // 遍历当前所有的有线网卡
    for (int i = 0; i < wiredDevices.size(); i++) {
        UWiredDevice *device = wiredDevices[i];
        WiredControllItem *ctrl = findWiredController(device);
        if (!ctrl)
            ctrl = new WiredControllItem(m_netListView->viewport(), device);
        else
            ctrl->updateView();

        items << ctrl;

        QList<UWiredConnection *> connItems = wiredConnections(device);
        for (int j = 0; j < connItems.size(); j++) {
            UWiredConnection *conn = connItems.at(j);
            WiredItem *connectionCtrl = findWiredItem(conn);
            if (!connectionCtrl)
                connectionCtrl = new WiredItem(m_netListView->viewport(), device, conn);
            else
                connectionCtrl->updateView();

            items << connectionCtrl;
        }
    }

    // 把原来列表中不存在的项放到移除列表中
    removeItems.clear();
    for (int i = m_items.size() - 1; i >=0 ; i--) {
        NetItem *item = m_items[i];
        if (!items.contains(item)) {
            m_items.removeOne(item);
            removeItems << item;
        }
    }

    m_items = items;
}

void NetworkPanel::updateView()
{
    QList<NetItem *> removeItems;

    updateItems(removeItems);

    // 先删除所有不存在的列表
    for (int i = 0; i < removeItems.size(); i++) {
        NetItem *item = removeItems[i];
        m_model->removeRow(item->standardItem()->row());
    }

    qDeleteAll(removeItems);
    removeItems.clear();
    int height = 0;
    int totalHeight = 0;
    for (int i = 0; i < m_items.size(); i++) {
        NetItem *item = m_items[i];
        int nRow = item->standardItem()->row();
        if (nRow < 0) {
            m_model->insertRow(i, item->standardItem());
        } else if (nRow != i) {
            m_model->takeItem(nRow, 0);
            m_model->removeRow(nRow);
            m_model->insertRow(i, item->standardItem());
        }

        QSize size = item->standardItem()->sizeHint();
        if (i < 16)
             height += size.height();

        totalHeight += size.height();
    }

    m_netListView->setFixedSize(PANELWIDTH, totalHeight);
    m_centerWidget->setFixedSize(PANELWIDTH, totalHeight);
    m_applet->setFixedSize(PANELWIDTH, height);
    m_netListView->update();
}

QStringList NetworkPanel::ipTipsMessage(const UDeviceType &devType)
{
    int typeCount = deviceCount(devType);
    UDeviceType type = static_cast<UDeviceType>(devType);
    QStringList tipMessage;
    int deviceIndex = 1;
    QList<UNetworkDeviceBase *> devices = UNetworkController::instance()->devices();
    for (int i = 0; i < devices.size(); i ++) {
        UNetworkDeviceBase *device = devices[i];
        if (device->deviceType() != type)
            continue;

        QString ipv4 = device->ipv4();
        if (ipv4.isEmpty())
            continue;

        switch (type) {
        case UDeviceType::Wired: {
            if (typeCount == 1)
                tipMessage << tr("Wired connection: %1").arg(ipv4);
            else
                tipMessage << tr("Wired Network").append(QString("%1").arg(deviceIndex++)).append(":" + ipv4);
            break;
        }
        case UDeviceType::Wireless: {
            if (typeCount == 1)
                tipMessage << tr("Wireless connection: %1").arg(ipv4);
            else
                tipMessage << tr("Wireless Network").append(QString("%1").arg(deviceIndex++)).append(":" + ipv4);
            break;
        }
        default: break;
        }
    }

    return tipMessage;
}

void NetworkPanel::updateTooltips()
{
    switch (m_pluginState) {
    case PluginState::Connected: {
        QStringList textList;
        textList << ipTipsMessage(UDeviceType::Wireless) << ipTipsMessage(UDeviceType::Wired);
        m_tipsWidget->setTextList(textList);
        break;
    }
    case PluginState::WirelessConnected:
        m_tipsWidget->setTextList(ipTipsMessage(UDeviceType::Wireless));
        break;
    case PluginState::WiredConnected:
        m_tipsWidget->setTextList(ipTipsMessage(UDeviceType::Wired));
        break;
    case PluginState::Disabled:
    case PluginState::WirelessDisabled:
    case PluginState::WiredDisabled:
        m_tipsWidget->setText(tr("Device disabled"));
        break;
    case PluginState::Unknow:
    case PluginState::Nocable:
        m_tipsWidget->setText(tr("Network cable unplugged"));
        break;
    case PluginState::Disconnected:
    case PluginState::WirelessDisconnected:
    case PluginState::WiredDisconnected:
        m_tipsWidget->setText(tr("Not connected"));
        break;
    case PluginState::Connecting:
    case PluginState::WirelessConnecting:
    case PluginState::WiredConnecting:
        m_tipsWidget->setText(tr("Connecting"));
        break;
    case PluginState::ConnectNoInternet:
    case PluginState::WirelessConnectNoInternet:
    case PluginState::WiredConnectNoInternet:
        m_tipsWidget->setText(tr("Connected but no Internet access"));
        break;
    case PluginState::Failed:
    case PluginState::WirelessFailed:
    case PluginState::WiredFailed:
        m_tipsWidget->setText(tr("Connection failed"));
        break;
    }
}

void NetworkPanel::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    const QRectF &rf = rect();
    const QRectF &rfp = QRectF(m_iconPixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / m_iconPixmap.devicePixelRatioF(),
                       m_iconPixmap);
}

void NetworkPanel::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    // 保持横纵比
    if (position == Dock::Bottom || position == Dock::Top) {
        setMaximumWidth(height());
        setMaximumHeight(QWIDGETSIZE_MAX);
    } else {
        setMaximumHeight(width());
        setMaximumWidth(QWIDGETSIZE_MAX);
    }

    refreshIcon();
}

int NetworkPanel::deviceCount(const UDeviceType &devType)
{
    // 获取指定的设备类型的设备数量
    int count = 0;
    QList<UNetworkDeviceBase *> devices = UNetworkController::instance()->devices();
    for (int i = 0; i < devices.size(); i ++) {
        UNetworkDeviceBase *dev = devices[i];
        if (dev->deviceType() == static_cast<UDeviceType>(devType))
            count++;
    }

    return count;
}

void NetworkPanel::onDeviceAdded(QList<UNetworkDeviceBase *> devices)
{
    // 处理新增设备的信号
    for (int i = 0; i < devices.size(); i++) {
        UNetworkDeviceBase *device = devices[i];
        // 当网卡连接状态发生变化的时候重新绘制任务栏的图标
        connect(device, &UNetworkDeviceBase::deviceStatusChanged, this, &NetworkPanel::updatePlugView);
        switch (device->deviceType()) {
        case UDeviceType::Wired: {
            UWiredDevice *wiredDevice = static_cast<UWiredDevice *>(device);

            connect(wiredDevice, &UNetworkDeviceBase::deviceStatusChanged, this, &NetworkPanel::wirelessChanged);
            connect(wiredDevice, &UNetworkDeviceBase::enableChanged, this, &NetworkPanel::wirelessChanged);
            connect(wiredDevice, &UNetworkDeviceBase::connectionChanged, this, &NetworkPanel::wirelessChanged);
            break;
        }
        case UDeviceType::Wireless: {
            UWirelessDevice *wirelessDevice = static_cast<UWirelessDevice *>(device);

            connect(wirelessDevice, &UWirelessDevice::networkAdded, this, &NetworkPanel::wirelessChanged);
            connect(wirelessDevice, &UWirelessDevice::networkRemoved, this, &NetworkPanel::wirelessChanged);
            connect(wirelessDevice, &UWirelessDevice::networkInfoChanged, this, &NetworkPanel::wirelessChanged);
            connect(wirelessDevice, &UWirelessDevice::enableChanged, this, &NetworkPanel::wirelessChanged);
            connect(wirelessDevice, &UWirelessDevice::connectionChanged, this, &NetworkPanel::wirelessChanged);

            wirelessDevice->scanNetwork();
            break;
        }
        default:break;
        }
    }

    updatePlugView();
}

void NetworkPanel::invokeMenuItem(const QString &menuId)
{
    // 有线设备是否可用
    bool wiredEnabled = deviceEnabled(UDeviceType::Wired);
    // 无线设备是否可用
    bool wirelessEnabeld = deviceEnabled(UDeviceType::Wireless);
    if (menuId == MenueEnable) {
        setDeviceEnabled(UDeviceType::Wired, !wiredEnabled);
        setDeviceEnabled(UDeviceType::Wireless, !wirelessEnabeld);
    } else if (menuId == MenueWiredEnable) {
        setDeviceEnabled(UDeviceType::Wired, !wiredEnabled);
    } else if (menuId == MenueWirelessEnable) {
        setDeviceEnabled(UDeviceType::Wireless, !wirelessEnabeld);
    } else if (menuId == MenueSettings) {
        DDBusSender()
                .service("com.deepin.dde.ControlCenter")
                .interface("com.deepin.dde.ControlCenter")
                .path("/com/deepin/dde/ControlCenter")
                .method(QString("ShowModule"))
                .arg(QString("network"))
                .call();
    }
}

bool NetworkPanel::needShowControlCenter()
{
    // 得到有线设备和无线设备的数量
    int wiredCount = deviceCount(UDeviceType::Wired);
    int wirelessCount = deviceCount(UDeviceType::Wireless);
    bool onlyOneTypeDevice = false;
    if ((wiredCount == 0 && wirelessCount > 0)
            || (wiredCount > 0 && wirelessCount == 0))
        onlyOneTypeDevice = true;

    if (onlyOneTypeDevice) {
        switch (m_pluginState) {
        case PluginState::Unknow:
        case PluginState::Nocable:
        case PluginState::WiredFailed:
        case PluginState::WirelessConnectNoInternet:
        case PluginState::WiredConnectNoInternet:
        case PluginState::WirelessDisconnected:
        case PluginState::WiredDisconnected:
        case PluginState::Disabled:
        case PluginState::WiredDisabled:
            return true;
        default:
            return false;
        }
    } else {
        switch (m_pluginState) {
        case PluginState::Unknow:
        case PluginState::Nocable:
        case PluginState::WiredFailed:
        case PluginState::ConnectNoInternet:
        case PluginState::Disconnected:
        case PluginState::Disabled:
            return true;
        default:
            return false;
        }
    }

    Q_UNREACHABLE();
    return true;
}

bool NetworkPanel::deviceEnabled(const UDeviceType &deviceType) const
{
    QList<UNetworkDeviceBase *> devices = UNetworkController::instance()->devices();
    for (int i = 0; i < devices.size(); i ++) {
        UNetworkDeviceBase *device = devices[i];
        if (device->deviceType() == deviceType && device->isEnabled())
            return true;
    }

    return false;
}

void NetworkPanel::setDeviceEnabled(const UDeviceType &deviceType, bool enabeld)
{
    QList<UNetworkDeviceBase *> devices = UNetworkController::instance()->devices();
    for (int i = 0; i < devices.size(); i ++) {
        UNetworkDeviceBase *dev = devices[i];
        if (dev->deviceType() == deviceType)
            dev->setEnabled(enabeld);
    }
}

const QString NetworkPanel::contextMenu() const
{
    bool wiredEnabled = deviceEnabled(UDeviceType::Wired);
    bool wirelessEnabeld = deviceEnabled(UDeviceType::Wireless);
    QList<QVariant> items;
    if (wiredEnabled && wirelessEnabeld) {
        items.reserve(3);
        QMap<QString, QVariant> wireEnable;
        wireEnable["itemId"] = MenueWiredEnable;
        if (wiredEnabled)
            wireEnable["itemText"] = tr("Disable wired connection");
        else
            wireEnable["itemText"] = tr("Enable wired connection");

        wireEnable["isActive"] = true;
        items.push_back(wireEnable);

        QMap<QString, QVariant> wirelessEnable;
        wirelessEnable["itemId"] = MenueWirelessEnable;
        if (wirelessEnabeld)
            wirelessEnable["itemText"] = tr("Disable wireless connection");
        else
            wirelessEnable["itemText"] = tr("Enable wireless connection");

        wirelessEnable["isActive"] = true;
        items.push_back(wirelessEnable);
    } else {
        items.reserve(2);
        QMap<QString, QVariant> enable;
        enable["itemId"] = MenueEnable;
        if (wiredEnabled || wirelessEnabeld)
            enable["itemText"] = tr("Disable network");
        else
            enable["itemText"] = tr("Enable network");

        enable["isActive"] = true;
        items.push_back(enable);
    }

    QMap<QString, QVariant> settings;
    settings["itemId"] = MenueSettings;
    settings["itemText"] = tr("Network settings");
    settings["isActive"] = true;
    items.push_back(settings);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

QWidget *NetworkPanel::itemTips()
{
    return m_tipsWidget;
}

QWidget *NetworkPanel::itemApplet()
{
    m_applet->setVisible(true);
    return m_applet;
}

bool NetworkPanel::hasDevice()
{
    return UNetworkController::instance()->devices().size() > 0;
}

void NetworkPanel::refreshIcon()
{
    setControlBackground();

    QString stateString;
    QString iconString;
    const auto ratio = devicePixelRatioF();
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    int strength = 0;

    switch (m_pluginState) {
    case PluginState::Disabled:
    case PluginState::WirelessDisabled:
        stateString = "disabled";
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        break;
    case PluginState::WiredDisabled:
        stateString = "disabled";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case PluginState::Connected:
    case PluginState::WirelessConnected:
        strength = getStrongestAp();
        stateString = WirelessItem::getStrengthStateString(strength);
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        break;
    case PluginState::WiredConnected:
        stateString = "online";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case PluginState::Disconnected:
    case PluginState::WirelessDisconnected:
        stateString = "0";
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        break;
    case PluginState::WiredDisconnected:
        stateString = "none";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case PluginState::Connecting: {
        m_refreshIconTimer->start();
        if (m_switchWire) {
            strength = QTime::currentTime().msec() / 10 % 100;
            stateString = WirelessItem::getStrengthStateString(strength);
            iconString = QString("wireless-%1-symbolic").arg(stateString);
            if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                    && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                iconString.append(PLUGIN_MIN_ICON_NAME);
            m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
            update();
            return;
        } else {
            m_refreshIconTimer->start(200);
            const int index = QTime::currentTime().msec() / 200 % 10;
            const int num = index + 1;
            iconString = QString("network-wired-symbolic-connecting%1").arg(num);
            if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                    && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                iconString.append(PLUGIN_MIN_ICON_NAME);
            m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
            update();
            return;
        }
    }
    case PluginState::WirelessConnecting: {
        m_refreshIconTimer->start();
        strength = QTime::currentTime().msec() / 10 % 100;
        stateString = WirelessItem::getStrengthStateString(strength);
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
            iconString.append(PLUGIN_MIN_ICON_NAME);
        m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
        update();
        return;
    }
    case PluginState::WiredConnecting: {
        m_refreshIconTimer->start(200);
        const int index = QTime::currentTime().msec() / 200 % 10;
        const int num = index + 1;
        iconString = QString("network-wired-symbolic-connecting%1").arg(num);
        if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
            iconString.append(PLUGIN_MIN_ICON_NAME);
        m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
        update();
        return;
    }
    case PluginState::ConnectNoInternet:
    case PluginState::WirelessConnectNoInternet: {
        // 无线已连接但无法访问互联网 offline
        stateString = "offline";
        iconString = QString("network-wireless-%1-symbolic").arg(stateString);
        break;
    }
    case PluginState::WiredConnectNoInternet: {
        stateString = "warning";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    }
    case PluginState::WiredFailed: {
        // 有线连接失败none变为offline
        stateString = "offline";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    }
    case PluginState::Unknow:
    case PluginState::Nocable: {
        stateString = "error";// 待图标 暂用错误图标
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    }
    case PluginState::WirelessFailed:
    case PluginState::Failed: {
        // 无线连接失败改为 disconnect
        stateString = "disconnect";
        iconString = QString("wireless-%1").arg(stateString);
        break;
    }
    }

    m_refreshIconTimer->stop();

    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconString.append(PLUGIN_MIN_ICON_NAME);

    m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);

    update();
}

void NetworkPanel::setControlBackground()
{
    QPalette backgroud;
    QColor separatorColor;
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        backgroud.setColor(QPalette::Background, QColor(255, 255, 255, 0.03 * 255));
    else
        backgroud.setColor(QPalette::Background, QColor(0, 0, 0, 0.03 * 255));

    m_applet->setAutoFillBackground(true);
    m_applet->setPalette(backgroud);
}

void NetworkPanel::updatePlugView()
{
    getPluginState();
    refreshIcon();
    updateTooltips();
    updateView();
}

void NetworkPanel::wirelessChanged()
{
    updatePlugView();
}

void NetworkPanel::onClickListView(const QModelIndex &index)
{
    NetItemType type = static_cast<NetItemType>(index.data(NetItemRole::TypeRole).toInt());
    switch (type) {
    case WirelessViewItem: {
        UAccessPoints *ap = static_cast<UAccessPoints *>(index.data(NetItemRole::DataRole).value<void *>());
        if (ap && !ap->connected()) {
            UWirelessDevice *device = static_cast<UWirelessDevice *>(index.data(NetItemRole::DeviceDataRole).value<void *>());
            UAccessPoints *activeAp = device->activeAccessPoints();
            if (activeAp != ap)
                device->connectNetwork(ap);
        }
        break;
    }
    case WiredViewItem: {
        UWiredConnection *conn = static_cast<UWiredConnection *>(index.data(NetItemRole::DataRole).value<void *>());
        if (conn && !conn->connected()) {
            UWiredDevice *device = static_cast<UWiredDevice *>(index.data(NetItemRole::DeviceDataRole).value<void *>());
            device->connectNetwork(conn);
        }
        break;
    }
    default: break;
    }
}

int NetworkPanel::getStrongestAp()
{
    int retStrength = -1;
    QList<UNetworkDeviceBase *> devices = UNetworkController::instance()->devices();
    for (int i = 0; i < devices.size(); i ++) {
        UNetworkDeviceBase *device = devices[i];
        if (device->deviceType() != UDeviceType::Wireless)
            continue;

        UWirelessDevice *dev = static_cast<UWirelessDevice *>(device);
        UAccessPoints *ap = dev->activeAccessPoints();
        if (ap && retStrength < ap->strength())
            retStrength = ap->strength();
    }

    return retStrength;
}

// 用于绘制分割线
NetworkDelegate::NetworkDelegate(QAbstractItemView *parent)
    : DStyledItemDelegate(parent)
{
}

NetworkDelegate::~NetworkDelegate()
{
}

void NetworkDelegate::paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    if (needDrawLine(index)) {
        QRect rct = option.rect;
        rct.setY(rct.top() + rct.height() - 2);
        rct.setHeight(2);
        if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
            painter->fillRect(rct, QColor(0, 0, 0, 255 * 0.1));
        else
            painter->fillRect(rct, QColor(255, 255, 255, 255 * 0.05));
    }

    DStyledItemDelegate::paint(painter, option, index);
}

bool NetworkDelegate::needDrawLine(const QModelIndex &index) const
{
    // 如果是最后一行，则无需绘制线条
    QModelIndex siblingIndex = index.siblingAtRow(index.row() + 1);
    if (!siblingIndex.isValid())
        return false;

    // 如果是总控开关，无线开关和有线开关，下面都要分割线
    NetItemType itemType = static_cast<NetItemType>(index.data(TypeRole).toInt());
    if (itemType == NetItemType::DeviceControllViewItem
            || itemType == NetItemType::WirelessControllViewItem
            || itemType == NetItemType::WiredControllViewItem)
        return true;

    NetItemType nextItemType = static_cast<NetItemType>(siblingIndex.data(TypeRole).toInt());
    return itemType != nextItemType;
}
