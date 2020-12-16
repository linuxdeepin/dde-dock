#include "networkitem.h"
#include "item/wireditem.h"
#include "item/wirelessitem.h"
#include "../../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"

#include <DHiDPIHelper>
#include <DApplicationHelper>
#include <DDBusSender>

#include <QVBoxLayout>
#include <QJsonDocument>

extern const int ItemWidth;
extern const int ItemMargin;
extern const int ItemHeight;
const int ControlItemHeight = 35;
const QString MenueEnable = "enable";
const QString MenueWiredEnable = "wireEnable";
const QString MenueWirelessEnable = "wirelessEnable";
const QString MenueSettings = "settings";

extern void initFontColor(QWidget *widget)
{
    if (!widget)
        return;

    auto fontChange = [&](QWidget * widget) {
        QPalette defaultPalette = widget->palette();
        defaultPalette.setBrush(QPalette::WindowText, defaultPalette.brightText());
        widget->setPalette(defaultPalette);
    };

    fontChange(widget);

    QObject::connect(DApplicationHelper::instance(), &DApplicationHelper::themeTypeChanged, widget, [ = ] {
        fontChange(widget);
    });
}

NetworkItem::NetworkItem(QWidget *parent)
    : QWidget(parent)
    , m_tipsWidget(new Dock::TipsWidget(this))
    , m_applet(new QScrollArea(this))
    , m_isWireless(true)
    , m_connectingTimer(new QTimer(this))
{
    m_connectingTimer->setInterval(200);
    m_connectingTimer->setSingleShot(false);

    m_tipsWidget->setVisible(false);

    auto defaultFont = font();
    auto titlefont = QFont(defaultFont.family(), defaultFont.pointSize() + 2);

    m_wirelessControlPanel = new QWidget(this);
    m_wirelessTitle = new QLabel(m_wirelessControlPanel);
    m_wirelessTitle->setText(tr("Wireless Network"));
    m_wirelessTitle->setFont(titlefont);
    initFontColor(m_wirelessTitle);
    m_switchWirelessBtn = new DSwitchButton(m_wirelessControlPanel);
    m_switchWirelessBtnState = false;

    const QPixmap pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh.svg");

    m_loadingIndicator = new DLoadingIndicator;
    m_loadingIndicator->setLoading(false);
    m_loadingIndicator->setSmooth(true);
    m_loadingIndicator->setAniDuration(1000);
    m_loadingIndicator->setAniEasingCurve(QEasingCurve::InOutCirc);
    m_loadingIndicator->installEventFilter(this);
    m_loadingIndicator->setFixedSize(pixmap.size() / devicePixelRatioF());
    m_loadingIndicator->viewport()->setAutoFillBackground(false);
    m_loadingIndicator->setFrameShape(QFrame::NoFrame);
    m_loadingIndicator->installEventFilter(this);

    m_wirelessLayout = new QVBoxLayout;
    m_wirelessLayout->setMargin(0);
    m_wirelessLayout->setSpacing(0);
    auto switchWirelessLayout = new QHBoxLayout;
    switchWirelessLayout->setMargin(0);
    switchWirelessLayout->setSpacing(0);
    switchWirelessLayout->addSpacing(2);
    switchWirelessLayout->addWidget(m_wirelessTitle);
    switchWirelessLayout->addStretch();
    switchWirelessLayout->addWidget(m_loadingIndicator);
    switchWirelessLayout->addSpacing(10);
    switchWirelessLayout->addWidget(m_switchWirelessBtn);
    switchWirelessLayout->addSpacing(2);
    m_wirelessControlPanel->setLayout(switchWirelessLayout);
    m_wirelessControlPanel->setFixedHeight(ControlItemHeight);

    m_wiredControlPanel = new QWidget(this);

    m_wiredTitle = new QLabel(m_wiredControlPanel);
    m_wiredTitle->setText(tr("Wired Network"));
    m_wiredTitle->setFont(titlefont);
    initFontColor(m_wiredTitle);
    m_switchWiredBtn = new DSwitchButton(m_wiredControlPanel);
    m_switchWiredBtnState = false;
    m_wiredLayout = new QVBoxLayout;
    m_wiredLayout->setMargin(0);
    m_wiredLayout->setSpacing(0);
    auto switchWiredLayout = new QHBoxLayout;
    switchWiredLayout->setMargin(0);
    switchWiredLayout->setSpacing(0);
    switchWiredLayout->addSpacing(2);
    switchWiredLayout->addWidget(m_wiredTitle);
    switchWiredLayout->addStretch();
    switchWiredLayout->addWidget(m_switchWiredBtn);
    switchWiredLayout->addSpacing(2);
    m_wiredControlPanel->setLayout(switchWiredLayout);
    m_wiredControlPanel->setFixedHeight(ControlItemHeight);

    auto centralWidget = new QWidget(m_applet);
    auto centralLayout = new QVBoxLayout;
    centralLayout->setContentsMargins(QMargins(ItemMargin, 0, ItemMargin, 0));
    centralLayout->setSpacing(0);
    centralLayout->addWidget(m_wirelessControlPanel);
    centralLayout->addLayout(m_wirelessLayout);
    centralLayout->addWidget(m_wiredControlPanel);
    centralLayout->addLayout(m_wiredLayout);
    centralWidget->setLayout(centralLayout);
    centralWidget->setFixedWidth(ItemWidth);
    centralWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Preferred);

    m_applet->setFixedWidth(ItemWidth);
    m_applet->setWidget(centralWidget);
    m_applet->setFrameShape(QFrame::NoFrame);
    m_applet->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_applet->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    centralWidget->setAutoFillBackground(false);
    m_applet->viewport()->setAutoFillBackground(false);
    m_applet->setVisible(false);
    m_applet->verticalScrollBar()->setContextMenuPolicy(Qt::NoContextMenu);

    connect(m_connectingTimer, &QTimer::timeout, this, &NetworkItem::onConnecting);
    connect(m_switchWiredBtn, &DSwitchButton::toggled, this, &NetworkItem::wiredsEnable);
    connect(m_switchWirelessBtn, &DSwitchButton::toggled, this, &NetworkItem::wirelessEnable);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &NetworkItem::onThemeTypeChanged);
}

QWidget *NetworkItem::itemApplet()
{
    m_applet->setVisible(true);
    wirelessItemsRequireScan();

    return m_applet;
}

QWidget *NetworkItem::itemTips()
{
    return m_tipsWidget;
}

void NetworkItem::updateDeviceItems(QMap<QString, WiredItem *> &wiredItems, QMap<QString, WirelessItem *> &wirelessItems)
{
    // 已有设备不重复进行增删操作
    auto tempWiredItems = m_wiredItems;
    auto tempWirelessItems = m_wirelessItems;

    for (auto wirelessItem : wirelessItems) {
        if (wirelessItem) {
            auto path = wirelessItem->path();
            if (m_wirelessItems.contains(path)) {
                m_wirelessItems.value(path)->setDeviceInfo(wirelessItem->deviceInfo());
                tempWirelessItems.remove(path);
                delete wirelessItem;
            } else {
                wirelessItem->setParent(this);
                m_wirelessItems.insert(path, wirelessItem);
            }
        }
    }

    for (auto wiredItem : wiredItems) {
        if (wiredItem) {
            auto path = wiredItem->path();
            if (m_wiredItems.contains(path)) {
                m_wiredItems.value(path)->setTitle(wiredItem->deviceName());
                tempWiredItems.remove(path);
                delete wiredItem;
            } else {
                wiredItem->setParent(this);
                m_wiredItems.insert(path, wiredItem);
                wiredItem->setVisible(true);
                m_wiredLayout->addWidget(wiredItem);
            }
        }
    }

    for (auto wirelessItem : tempWirelessItems) {
        if (wirelessItem) {
            auto path = wirelessItem->device()->path();
            m_wirelessItems.remove(path);
            m_connectedWirelessDevice.remove(path);
            wirelessItem->itemApplet()->setVisible(false);
            m_wirelessLayout->removeWidget(wirelessItem->itemApplet());
            delete wirelessItem;
        }
    }
    for (auto wiredItem : tempWiredItems) {
        if (wiredItem) {
            auto path = wiredItem->device()->path();
            m_wiredItems.remove(path);
            m_connectedWiredDevice.remove(path);
            wiredItem->setVisible(false);
            m_wiredLayout->removeWidget(wiredItem);
            delete wiredItem;
        }
    }
    //正常的刷新流程
    updateSelf();
}

const QString NetworkItem::contextMenu() const
{
    QList<QVariant> items;

    if (m_wirelessItems.size() && m_wiredItems.size()) {
        items.reserve(3);
        QMap<QString, QVariant> wireEnable;
        wireEnable["itemId"] = MenueWiredEnable;
        if (m_switchWiredBtnState)
            wireEnable["itemText"] = tr("Disable wired connection");
        else
            wireEnable["itemText"] = tr("Enable wired connection");
        wireEnable["isActive"] = true;
        items.push_back(wireEnable);

        QMap<QString, QVariant> wirelessEnable;
        wirelessEnable["itemId"] = MenueWirelessEnable;
        if (m_switchWirelessBtnState)
            wirelessEnable["itemText"] = tr("Disable wireless connection");
        else
            wirelessEnable["itemText"] = tr("Enable wireless connection");
        wirelessEnable["isActive"] = true;
        items.push_back(wirelessEnable);
    } else {
        items.reserve(2);
        QMap<QString, QVariant> enable;
        enable["itemId"] = MenueEnable;
        if (m_switchWiredBtnState || m_switchWirelessBtnState)
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

void NetworkItem::invokeMenuItem(const QString &menuId, const bool checked)
{
    Q_UNUSED(checked);

    if (menuId == MenueEnable) {
        wiredsEnable(!m_switchWiredBtnState);
        wirelessEnable(!m_switchWirelessBtnState);
    } else if (menuId == MenueWiredEnable)
        wiredsEnable(!m_switchWiredBtnState);
    else if (menuId == MenueWirelessEnable)
        wirelessEnable(!m_switchWirelessBtnState);
    else if (menuId == MenueSettings)
        DDBusSender()
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method(QString("ShowModule"))
        .arg(QString("network"))
        .call();
}

void NetworkItem::refreshIcon()
{
    // 刷新按钮图标
    QPixmap pixmap;
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh_dark.svg");
    else
        pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh.svg");
    m_loadingIndicator->setImageSource(pixmap);

    QString stateString;
    QString iconString;
    const auto ratio = devicePixelRatioF();
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    int strength = 0;
    switch (m_pluginState) {
    case Adisabled:
        stateString = "disabled";
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        break;
    case Bdisabled:
        stateString = "disabled";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Connected:
    case Aconnected:
        strength = getStrongestAp();
        qDebug() << "strength:" << strength;
        //这个区间是需求文档中规定的
        if (strength > 65) {
            stateString = "80";
        } else if (strength > 55) {
            stateString = "60";
        } else if (strength > 30) {
            stateString = "40";
        } else if (strength > 5) {
            stateString = "20";
        } else {
            stateString = "0";
        }
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        break;
    case Bconnected:
        stateString = "online";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Adisconnected:
        stateString = "0";
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        break;
    case Bdisconnected:
        stateString = "none";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Aconnecting: {
        m_isWireless = true;
        m_connectingTimer->start();
        return;
    }
    case Bconnecting: {
        m_isWireless = false;
        m_connectingTimer->start();
        return;
    }
    case AconnectNoInternet:
        stateString = "offline";
        iconString = QString("network-wireless-%1-symbolic").arg(stateString);
        break;
    case BconnectNoInternet:
        stateString = "offline";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Bfailed:
        stateString = "none";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Unknown:
    case Nocable:
        stateString = "error";//待图标 暂用错误图标
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    }

    m_connectingTimer->stop();

    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconString.append(PLUGIN_MIN_ICON_NAME);

    m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);

    update();
}

void NetworkItem::resizeEvent(QResizeEvent *e)
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

void NetworkItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    const QRectF &rf = rect();
    const QRectF &rfp = QRectF(m_iconPixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / m_iconPixmap.devicePixelRatioF(),
                       m_iconPixmap);
}

bool NetworkItem::eventFilter(QObject *obj, QEvent *event)
{
    if (obj == m_loadingIndicator) {
        if (event->type() == QEvent::MouseButtonPress) {
            qDebug() << Q_FUNC_INFO;
            wirelessItemsRequireScan();
        }
    }
    return false;
}

void NetworkItem::wiredsEnable(bool enable)
{
    for (auto wiredItem : m_wiredItems) {
        if (wiredItem) {
            wiredItem->setDeviceEnabled(enable);
        }
    }
//    updateSelf();
}

void NetworkItem::wirelessEnable(bool enable)
{
    for (auto wirelessItem : m_wirelessItems) {
        if (wirelessItem) {
            wirelessItem->setDeviceEnabled(enable);
            enable ? m_wirelessLayout->addWidget(wirelessItem->itemApplet())
            : m_wirelessLayout->removeWidget(wirelessItem->itemApplet());
            wirelessItem->itemApplet()->setVisible(enable);
        }
    }
    updateSelf();
}

void NetworkItem::onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType)
{
    for (auto wiredItem : m_wiredItems) {
        wiredItem->setThemeType(themeType);
    }
    refreshIcon();
}

void NetworkItem::getPluginState()
{
    int wiredState = 0;
    int wirelessState = 0;
    int state = 0;
    int temp = 0;
    // 所有设备状态叠加
    QMapIterator<QString, WirelessItem *> iwireless(m_wirelessItems);
    //无线网状态获取
    while (iwireless.hasNext()) {
        iwireless.next();
        auto wirelessItem = iwireless.value();
        if (wirelessItem) {
            temp = wirelessItem->getDeviceState();
            state |= temp;
            //如果该网卡处于连接成功状态，则加入到列表中
            if (temp == WirelessItem::Connected) {
                m_connectedWirelessDevice.insert(iwireless.key(), wirelessItem);
            }
        }
    }
    // 按如下顺序得到当前无线设备状态
    temp = state;
    if (!temp)
        wirelessState = WirelessItem::Unknown;
    if (temp & WirelessItem::Disabled)
        wirelessState = WirelessItem::Disabled;
    if (temp & WirelessItem::Disconnected)
        wirelessState = WirelessItem::Disconnected;
    //获取ip和验证的时候则代表还在连接过程中
    if (temp & WirelessItem::Connecting ||
            temp & WirelessItem::Authenticating || temp & WirelessItem::ObtainingIP)
        wirelessState = WirelessItem::Connecting;
    if (temp & WirelessItem::ConnectNoInternet)
        wirelessState = WirelessItem::ConnectNoInternet;
    if (temp & WirelessItem::Connected) {
        wirelessState = WirelessItem::Connected;
    }
    //有线网状态获取
    state = 0;
    temp = 0;
    QMapIterator<QString, WiredItem *> iwired(m_wiredItems);
    while (iwired.hasNext()) {
        iwired.next();
        auto wiredItem = iwired.value();
        if (wiredItem) {
            temp = wiredItem->getDeviceState();
            state |= temp;
            if (temp == WiredItem::Connected) {
                m_connectedWiredDevice.insert(iwired.key(), wiredItem);
            }
        }
    }
    temp = state;
    if (!temp)
        wiredState = WiredItem::Unknow;
    if (temp & WiredItem::Nocable)
        wiredState = WiredItem::Nocable;
    if (temp & WiredItem::Disabled)
        wiredState = WiredItem::Disabled;
    if (temp & WiredItem::Disconnected)
        wiredState = WiredItem::Disconnected;
    if (temp & WiredItem::Connecting ||
            temp & WiredItem::Authenticating || temp & WiredItem::ObtainingIP)
        wiredState = WiredItem::Connecting;
    if (temp & WiredItem::ConnectNoInternet)
        wiredState = WiredItem::ConnectNoInternet;
    if (temp & WiredItem::Connected) {
        wiredState = WiredItem::Connected;
    }
    qDebug() << "wirelessState = " << wirelessState << "wiredState =" << wiredState;
    switch (wirelessState) {
        //连接过程中断开第一个连接会发送该状态，导致前端闪一秒有线断开的图标，所以直接对该状态进行连接中处理
        case WirelessItem::Unknown:
            if(m_connectedWirelessDevice.isEmpty())
                m_pluginState = PluginState(wiredState);
            else
                m_pluginState = Aconnecting;

            break;
        case WirelessItem::Disabled:
            if (wiredState < WiredItem::Disconnected)
                m_pluginState = Adisabled;
            else
                m_pluginState = PluginState(wiredState);
            break;
        case WirelessItem::Connected:
             if (wiredState == WiredItem::Connected)
                 m_pluginState = Connected;
             else
                 m_pluginState = Aconnected;
             break;
        case WirelessItem::Disconnected:
            if (wiredState < WiredItem::Disconnected)
                m_pluginState = Adisconnected;
            else
                m_pluginState = PluginState(wiredState);
            break;
        case WirelessItem::Connecting:
            m_pluginState = Aconnecting;
            break;
        case WirelessItem::ConnectNoInternet:
            if (wiredState == WiredItem::Connected || wiredState == WiredItem::Connecting)
                m_pluginState = PluginState(wiredState);
            else
                m_pluginState = AconnectNoInternet;
            break;
        default:
            m_pluginState = PluginState(wiredState);
    }
}

void NetworkItem::updateView()
{
    // 固定显示高度即为固定示项目数
    const int constDisplayItemCnt = 10;
    int contentHeight = 0;
    int itemCount = 0;

    auto wirelessCnt = m_wirelessItems.size();
    if (m_switchWirelessBtnState) {
        for (auto wirelessItem : m_wirelessItems) {
            if (wirelessItem) {
                if (wirelessItem->device()->enabled())
                    itemCount += wirelessItem->APcount();
                // 单个设备开关控制项
                if (wirelessCnt == 1) {
                    wirelessItem->setControlPanelVisible(false);
                    continue;
                } else {
                    wirelessItem->setControlPanelVisible(true);
                }
                itemCount++;
            }
        }
    }
    // 设备总控开关只与是否有设备相关
    auto wirelessDeviceCnt = m_wirelessItems.size();
    if (wirelessDeviceCnt)
        contentHeight += m_wirelessControlPanel->height();
    m_wirelessControlPanel->setVisible(wirelessDeviceCnt);

    auto wiredDeviceCnt = m_wiredItems.size();
    if (wiredDeviceCnt)
        contentHeight += m_wiredControlPanel->height();
    m_wiredControlPanel->setVisible(wiredDeviceCnt);

    itemCount += wiredDeviceCnt;

    auto centralWidget = m_applet->widget();
    if (itemCount <= constDisplayItemCnt) {
        contentHeight += (itemCount - wiredDeviceCnt) * ItemHeight;
        contentHeight += wiredDeviceCnt * ItemHeight;
        centralWidget->setFixedHeight(contentHeight);
        m_applet->setFixedHeight(contentHeight);
    } else {
        contentHeight += (itemCount - wiredDeviceCnt) * ItemHeight;
        contentHeight += wiredDeviceCnt * ItemHeight;
        centralWidget->setFixedHeight(contentHeight);
        m_applet->setFixedHeight(constDisplayItemCnt * ItemHeight);
        m_applet->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOn);
    }
}

void NetworkItem::updateSelf()
{
    getPluginState();
    updateMasterControlSwitch();
    refreshIcon();
    refreshTips();
    updateView();
}

void NetworkItem::onConnecting()
{
    QString stateString;
    QString iconString;
    int strength = 0;
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    const auto ratio = devicePixelRatioF();
    if (m_isWireless)
    {
        strength = QTime::currentTime().msec() / 10 % 100;
        if (strength == 100) {
            stateString = "80";
        } else if (strength < 20) {
            stateString = "0";
        } else {
            stateString = QString::number(strength / 10 & ~0x1) + "0";
        }
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
            iconString.append(PLUGIN_MIN_ICON_NAME);
        m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
    } else {
        const int index = QTime::currentTime().msec() / 200 % 10;
        const int num = index + 1;
        iconString = QString("network-wired-symbolic-connecting%1").arg(num);
        if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
            iconString.append(PLUGIN_MIN_ICON_NAME);
        m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
    }
    update();
    return;
}


int NetworkItem::getStrongestAp()
{
    int retStrength = -1;
    for (auto wirelessItem : m_connectedWirelessDevice) {
        auto apInfo = wirelessItem->getConnectedApInfo();
        if (apInfo.isEmpty())
            continue;
        auto strength = apInfo.value("Strength").toInt();
        qDebug() << "strength" << strength;
        if (retStrength < strength)
            retStrength = strength;
    }
    return retStrength;
}

void NetworkItem::wirelessItemsRequireScan()
{
    for (auto wirelessItem : m_wirelessItems) {
        if (wirelessItem) {
            Q_EMIT wirelessItem->requestWirelessScan();
        }
    }
    wirelessScan();
}

void NetworkItem::updateMasterControlSwitch()
{
    bool deviceState = false;
    for (auto wirelessItem : m_wirelessItems) {
        if (wirelessItem)
            if (wirelessItem->deviceEanbled()) {
                deviceState = true;
                break;
            }
    }
    m_switchWirelessBtn->blockSignals(true);
    m_switchWirelessBtn->setChecked(deviceState);
    m_loadingIndicator->setVisible(deviceState);
    m_switchWirelessBtn->blockSignals(false);
    if (deviceState) {
        for (auto wirelessItem : m_wirelessItems) {
            if (wirelessItem) {
                m_wirelessLayout->addWidget(wirelessItem->itemApplet());
                wirelessItem->itemApplet()->setVisible(true);
            }
        }
    } else {
        for (auto wirelessItem : m_wirelessItems) {
            if (wirelessItem) {
                m_wirelessLayout->removeWidget(wirelessItem->itemApplet());
                wirelessItem->itemApplet()->setVisible(false);
            }
        }
    }
    m_switchWirelessBtnState = deviceState;

    deviceState = false;
    for (auto wiredItem : m_wiredItems) {
        if (wiredItem)
            if (wiredItem->deviceEabled()) {
                deviceState = true;
                break;
            }
    }
    m_switchWiredBtn->blockSignals(true);
    m_switchWiredBtn->setChecked(deviceState);
    m_switchWiredBtn->blockSignals(false);
    m_switchWiredBtnState = deviceState;
}

void NetworkItem::refreshTips()
{
    switch (m_pluginState) {
    case Adisabled:
    case Bdisabled:
        m_tipsWidget->setText(tr("Device disabled"));
        break;
    case Aconnected: {
        QString strTips;
        for (auto wirelessItem : m_connectedWirelessDevice) {
            if (wirelessItem) {
                auto info = wirelessItem->getActiveWirelessConnectionInfo();
                if (!info.contains("Ip4"))
                    break;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    break;
                strTips = tr("Wireless connection: %1").arg(ipv4.value("Address").toString()) + '\n';
            }
        }
        strTips.chop(1);
        m_tipsWidget->setText(strTips);
    }
    break;
    case Bconnected: {
        QString strTips;
        for (auto wiredItem : m_connectedWiredDevice) {
            if (wiredItem) {
                auto info = wiredItem->getActiveWiredConnectionInfo();
                if (!info.contains("Ip4"))
                    break;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    break;
                strTips = tr("Wired connection: %1").arg(ipv4.value("Address").toString()) + '\n';
            }
        }
        strTips.chop(1);
        m_tipsWidget->setText(strTips);
    }
    break;
    case Adisconnected:
    case Bdisconnected:
        m_tipsWidget->setText(tr("Not connected"));
        break;
    case Aconnecting:
    case Bconnecting: {
        m_tipsWidget->setText(tr("Connecting"));
        return;
    }
    case AconnectNoInternet:
    case BconnectNoInternet:
        m_tipsWidget->setText(tr("Connected but no Internet access"));
        break;
    case Bfailed:
        m_tipsWidget->setText(tr("Network cable unplugged"));
        break;
    case Unknown:
    case Nocable:
        m_tipsWidget->setText(tr("Network cable unplugged"));
        break;
    case Connected: {
        QString strTips;
        QStringList textList;
        for (auto wirelessItem : m_connectedWirelessDevice) {
            if (wirelessItem) {
                auto info = wirelessItem->getActiveWirelessConnectionInfo();
                if (!info.contains("Ip4"))
                    break;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    break;
                strTips = tr("Wireless connection: %1").arg(ipv4.value("Address").toString()) + '\n';
                strTips.chop(1);
                textList << strTips;
            }
        }
        for (auto wiredItem : m_connectedWiredDevice) {
            if (wiredItem) {
                auto info = wiredItem->getActiveWiredConnectionInfo();
                if (!info.contains("Ip4"))
                    break;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    break;
                strTips = tr("Wired connection: %1").arg(ipv4.value("Address").toString()) + '\n';
                strTips.chop(1);
                textList << strTips;
            }
        }
        m_tipsWidget->setTextList(textList);
    }
    break;
    }
}

bool NetworkItem::isShowControlCenter()
{
    bool onlyOneTypeDevice = false;
    if ((m_wiredItems.size() == 0 && m_wirelessItems.size() > 0)
            || (m_wiredItems.size() > 0 && m_wirelessItems.size() == 0))
        onlyOneTypeDevice = true;

    if (onlyOneTypeDevice) {
        switch (m_pluginState) {
        case Unknown:
        case Nocable:
        case Bfailed:
        case AconnectNoInternet:
        case BconnectNoInternet:
        case Adisconnected:
        case Bdisconnected:
        case Adisabled:
        case Bdisabled:
            return true;
        default:
            break;
        }
    } else {
        switch (m_pluginState) {
        case Unknown:
        case Nocable:
        case Bfailed:
            return true;
        default:
            break;
        }
    }

    return false;
}

void NetworkItem::wirelessScan()
{
    if (m_loadingIndicator->loading())
        return;
    m_loadingIndicator->setLoading(true);
    QTimer::singleShot(1000, this, [ = ] {
        m_loadingIndicator->setLoading(false);
    });
}
