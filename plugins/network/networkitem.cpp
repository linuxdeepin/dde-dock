#include "networkitem.h"
#include "item/wireditem.h"
#include "item/wirelessitem.h"
#include "../../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"
#include "utils.h"

#include <DHiDPIHelper>
#include <DApplicationHelper>
#include <DDBusSender>
#include <DFontSizeManager>

#include <QVBoxLayout>
#include <QJsonDocument>
#include <QGSettings>
#include <QNetworkInterface>
#include <QHostAddress>
#include <QMap>

extern const int ItemWidth;
extern const int ItemMargin;
extern const int ItemHeight;
const QString MenueEnable = "enable";
const QString MenueWiredEnable = "wireEnable";
const QString MenueWirelessEnable = "wirelessEnable";
const QString MenueSettings = "settings";

#define TITLE_HEIGHT 46
#define ITEM_HEIGHT 36

NetworkItem::NetworkItem(QWidget *parent)
    : QWidget(parent)
    , m_tipsWidget(new Dock::TipsWidget(this))
    , m_applet(new QScrollArea(this))
    , m_switchWire(true)
    , m_timeOut(true)
    , refreshIconTimer(new QTimer(this))
    , m_switchWireTimer(new QTimer(this))
    , m_wirelessScanTimer(new QTimer(this))
    , m_wirelessScanInterval(Utils::SettingValue("com.deepin.dde.dock", QByteArray(), "wireless-scan-interval", 10).toInt())
    , m_firstSeparator(new HorizontalSeperator(this))
    , m_secondSeparator(new HorizontalSeperator(this))
    , m_thirdSeparator(new HorizontalSeperator(this))
    , m_networkInter(new DbusNetwork("com.deepin.daemon.Network", "/com/deepin/daemon/Network", QDBusConnection::sessionBus(), this))
    , m_detectConflictTimer(new QTimer(this))
    , m_ipConflict(false)
    , m_ipConflictChecking(false)
{
    refreshIconTimer->setInterval(100);

    m_tipsWidget->setVisible(false);

    m_wirelessControlPanel = new QWidget(this);

    QLabel *wirelessTitle = new QLabel(m_wirelessControlPanel);
    wirelessTitle->setText(tr("Wireless Network"));
    wirelessTitle->setFixedHeight(TITLE_HEIGHT);
    wirelessTitle->setForegroundRole(QPalette::BrightText);
    DFontSizeManager::instance()->bind(wirelessTitle, DFontSizeManager::T4, QFont::Medium);

    m_switchWirelessBtn = new DSwitchButton(m_wirelessControlPanel);
    m_switchWirelessBtnState = false;

    const QPixmap pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh.svg");

    m_loadingIndicator = new DLoadingIndicator(this);
    m_loadingIndicator->setLoading(false);
    m_loadingIndicator->setSmooth(true);
    m_loadingIndicator->setAniDuration(1000);
    m_loadingIndicator->setAniEasingCurve(QEasingCurve::InOutCirc);
    m_loadingIndicator->installEventFilter(this);
    m_loadingIndicator->setFixedSize(pixmap.size() / devicePixelRatioF());
    m_loadingIndicator->viewport()->setAutoFillBackground(false);
    m_loadingIndicator->setFrameShape(QFrame::NoFrame);
    m_loadingIndicator->installEventFilter(this);

    this->installEventFilter(this);

    m_wirelessLayout = new QVBoxLayout;
    m_wirelessLayout->setMargin(0);
    m_wirelessLayout->setSpacing(0);

    // 无线网络控制器
    QHBoxLayout *switchWirelessLayout = new QHBoxLayout;
    switchWirelessLayout->setMargin(0);
    switchWirelessLayout->setSpacing(0);
    // DSwitchButton 按照设计要求: 在保持现有控件的尺寸下,这里需要预留绘制focusRect的区域,borderWidth为2,间隙宽度为2
    // 所以此处按设计的要求 10-4 = 6 right margin
    switchWirelessLayout->setContentsMargins(20, 0, 6, 0);
    switchWirelessLayout->addWidget(wirelessTitle);
    switchWirelessLayout->addStretch();
    switchWirelessLayout->addWidget(m_loadingIndicator);
    switchWirelessLayout->addSpacing(4);
    switchWirelessLayout->addWidget(m_switchWirelessBtn);
    m_wirelessControlPanel->setLayout(switchWirelessLayout);
    m_wirelessControlPanel->setFixedHeight(TITLE_HEIGHT);

    m_wiredControlPanel = new QWidget(this);
    m_wiredControlPanel->setFixedHeight(TITLE_HEIGHT);

    QLabel *wiredTitle = new QLabel(m_wiredControlPanel);
    wiredTitle->setText(tr("Wired Network"));
    wiredTitle->setForegroundRole(QPalette::BrightText);
    DFontSizeManager::instance()->bind(wiredTitle, DFontSizeManager::T4, QFont::Medium);

    m_switchWiredBtn = new DSwitchButton(m_wiredControlPanel);
    m_switchWiredBtnState = false;
    m_wiredLayout = new QVBoxLayout;
    m_wiredLayout->setMargin(0);
    m_wiredLayout->setSpacing(0);

    // 有线网络控制器
    QHBoxLayout *switchWiredLayout = new QHBoxLayout;
    switchWiredLayout->setContentsMargins(20, 0, 6, 0);
    switchWiredLayout->addWidget(wiredTitle);
    switchWiredLayout->addStretch();
    switchWiredLayout->addWidget(m_switchWiredBtn);
    m_wiredControlPanel->setLayout(switchWiredLayout);
    m_wiredControlPanel->setFixedHeight(TITLE_HEIGHT);

    QWidget *centralWidget = new QWidget;
    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->setContentsMargins(QMargins(ItemMargin, 0, ItemMargin, 0));
    centralLayout->setSpacing(0);
    centralLayout->setMargin(0);

    centralLayout->addWidget(m_wirelessControlPanel);
    centralLayout->addWidget(m_firstSeparator);
    centralLayout->addLayout(m_wirelessLayout);
    centralLayout->addWidget(m_secondSeparator);

    //TODO 先暂时这样写，后面要重构，届时布局要重新修改，直接使用dlistview
    m_wirelessControlPanel->setVisible(m_wirelessItems.count() > 0);
    m_firstSeparator->setVisible(m_wirelessItems.count() > 0);
    m_secondSeparator->setVisible(m_wirelessItems.count() > 0 && m_wiredItems.count() > 0);

    centralLayout->addWidget(m_wiredControlPanel);
    centralLayout->addWidget(m_thirdSeparator);
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

    connect(m_switchWireTimer, &QTimer::timeout, [ = ] {
        m_switchWire = !m_switchWire;
        m_timeOut = true;
    });
    connect(refreshIconTimer, &QTimer::timeout, this, &NetworkItem::refreshIcon);
    connect(m_switchWiredBtn, &DSwitchButton::toggled, this, &NetworkItem::wiredsEnable);
    connect(m_switchWirelessBtn, &DSwitchButton::toggled, this, &NetworkItem::wirelessEnable);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &NetworkItem::onThemeTypeChanged);

    connect(m_networkInter, &DbusNetwork::IPConflict, this, &NetworkItem::ipConflict);
    connect(this, &NetworkItem::sendIpConflictDect, this, &NetworkItem::onSendIpConflictDect);
    connect(m_detectConflictTimer, &QTimer::timeout, this, &NetworkItem::onDetectConflict);
    const QGSettings *gsetting = Utils::SettingsPtr("com.deepin.dde.dock", QByteArray(), this);
    if (gsetting)
        connect(gsetting, &QGSettings::changed, [&](const QString &key) {
            if (key == "wireless-scan-interval") {
                m_wirelessScanInterval = gsetting->get("wireless-scan-interval").toInt() * 1000;
                m_wirelessScanTimer->setInterval(m_wirelessScanInterval);
            }
        });
    connect(m_wirelessScanTimer, &QTimer::timeout, [&] {
        for (auto wirelessItem : m_wirelessItems) {
            if (wirelessItem) {
                wirelessItem->requestWirelessScan();
            }
        }
    });

    m_wirelessScanTimer->setInterval(m_wirelessScanInterval);
}

QWidget *NetworkItem::itemApplet()
{
    m_applet->setVisible(true);
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
            disconnect(wirelessItem, &WirelessItem::apInfoChanged, this, &NetworkItem::refreshIcon);
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

    m_wirelessControlPanel->setVisible(m_wirelessItems.count() > 0);
    m_firstSeparator->setVisible(m_wirelessItems.count() > 0);
    m_secondSeparator->setVisible(m_wirelessItems.count() > 0 && m_wiredItems.count() > 0);

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
    case Disabled:
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
        stateString = getStrengthStateString(strength);
        iconString = QString("wireless-%1-symbolic").arg(stateString);

        //如果无线连接有IP冲突，则显示已连接但是无法访问网络的图标
        if (m_ipConflict && getActiveWirelessList().size() > 0) {
            foreach(auto ip, getActiveWirelessList()) {
                if (m_conflictMap.keys().contains(ip)) {
                    stateString = "offline";
                    iconString = QString("network-wireless-%1-symbolic").arg(stateString);
                    break;
                }
            }
        }

        break;
    case Bconnected:
        stateString = "online";
        iconString = QString("network-%1-symbolic").arg(stateString);

        //如果有线连接有IP冲突，则显示有线连接断开的图标
        if (m_ipConflict && getActiveWiredList().size() > 0) {
            foreach(auto ip, getActiveWiredList()) {
                if (m_conflictMap.keys().contains(ip)) {
                    stateString = "offline";
                    iconString = QString("network-%1-symbolic").arg(stateString);
                    break;
                }
            }
        }

        break;
    case Disconnected:
    case Adisconnected:
        stateString = "0";
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        break;
    case Bdisconnected:
        stateString = "none";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Connecting: {
        refreshIconTimer->start();
        if (m_switchWire) {
            strength = QTime::currentTime().msec() / 10 % 100;
            stateString = getStrengthStateString(strength);
            iconString = QString("wireless-%1-symbolic").arg(stateString);
            if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                    && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                iconString.append(PLUGIN_MIN_ICON_NAME);
            m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
            update();
            return;
        } else {
            refreshIconTimer->start(200);
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
    case Aconnecting: {
        refreshIconTimer->start();
        strength = QTime::currentTime().msec() / 10 % 100;
        stateString = getStrengthStateString(strength);
        iconString = QString("wireless-%1-symbolic").arg(stateString);
        if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
            iconString.append(PLUGIN_MIN_ICON_NAME);
        m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
        update();
        return;
    }
    case Bconnecting: {
        refreshIconTimer->start(200);
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
    case ConnectNoInternet:
    case AconnectNoInternet: //无线已连接但无法访问互联网 offline
        stateString = "offline";
        iconString = QString("network-wireless-%1-symbolic").arg(stateString);
        break;
    case BconnectNoInternet:
        stateString = "warning";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Bfailed://有线连接失败none变为offline
        stateString = "offline";
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Unknow:
    case Nocable:
        stateString = "error";//待图标 暂用错误图标
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    case Afailed:
    case Failed: //无线连接失败改为 disconnect
        stateString = "disconnect";
        iconString = QString("wireless-%1").arg(stateString);
        break;
    }

    refreshIconTimer->stop();

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
            for (auto wirelessItem : m_wirelessItems) {
                if (wirelessItem) {
                    wirelessItem->requestWirelessScan();
                }
            }
            wirelessScan();
        }
    }

    // 用户会鼠标悬浮检查网络状态，此时发起检测，并更新网络状态
    // 当主机有n张激活的网卡（包含无线网卡，有线网卡）时，鼠标移动到网络插件位置，会发起主动ip冲突检测，
    // 主动检测，如果冲突会触发IPConflict信号
    if (obj == this) {
        if (event->type() == QEvent::Enter) {
            onDetectConflict();
        }
    }

    return false;
}

QString NetworkItem::getStrengthStateString(int strength)
{
    if (5 >= strength)
        return "0";
    else if (30 >= strength)
        return "20";
    else if (55 >= strength)
        return "40";
    else if (65 >= strength)
        return "60";
    else
        return "80";

    Q_UNREACHABLE();
}

void NetworkItem::wiredsEnable(bool enable)
{
    // 刷新有线连接,不管是否禁用均显示
    for (auto wiredItem : m_wiredItems) {
        if (wiredItem) {
            wiredItem->setDeviceEnabled(enable);
            m_wiredLayout->addWidget(wiredItem);
        }
    }

    updateView();
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
    //禁用无线网络时对应的分割线设置为不可见防止两分割线叠加增加分割线高度与下面分割线高度不一样
    m_secondSeparator->setVisible(enable && m_wiredItems.count() > 0);
    updateView();
}

void NetworkItem::onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType)
{
    for (auto wiredItem : m_wiredItems) {
        wiredItem->setThemeType(themeType);
    }
    refreshIcon();
}

/**ip冲突以及冲突解除时，更新网络插件显示状态
 * @brief NetworkItem::ipConflict
 * @param ip 本机的ip地址
 * @param mac 与本机冲突的mac地址，不为空，则冲突，为空则ip冲突解除
 */
void NetworkItem::ipConflict(const QString &ip, const QString &mac)
{
    static int conflictCount = 0;
    static int removeCount = 0;
    QStringList ipList = currentIpList();

    //判断缓存冲突列表中的IP是否在本机IP列表中
    foreach (auto tmpIP, m_conflictMap.keys()) {
        if (!ipList.contains(tmpIP)) {
            m_conflictMap.remove(tmpIP);
        }
    }
    // 如果不是本机ip
    if (!ipList.contains(ip))
        return;

    // mac为空时冲突解除或IP地址不冲突
    if (mac.isEmpty()) {
        // 自检为空，则解除ip冲突状态或者冲突列表为空时，更新状态
        m_conflictMap.remove(ip);

        if (m_conflictMap.isEmpty()) {
            conflictCount = 0;
        }

        if (m_conflictMap.isEmpty() && m_ipConflict) {
            // 确认1次解除
            if (removeCount++ < 1) {
                onDetectConflict();
                return;
            }

            // 当mac为空且map也为空时，立即更新状态会导致文字显示由'ip地址冲突'变为'已连接网络但无法访问互联网'
            // 因为加了次数判断
            m_ipConflict = false;
            m_ipConflictChecking = false;
            m_detectConflictTimer->stop();
            updateSelf();
            m_conflictMap.clear();
            removeCount = 0;
        }
        return;
    }

    // 缓存冲突ip和mac地址
    m_conflictMap.insert(ip, mac);
    removeCount = 0;

    if (m_conflictMap.size() && !m_ipConflict) {
        // 确认2次
        if (conflictCount++ < 2) {
            onDetectConflict();
            return;
        }

        conflictCount = 0;
        m_ipConflict = true;
        updateSelf();

        // 有冲突才开启5秒中ip自检，目的是当其他主机主动解除了冲突，我方不知情
        m_detectConflictTimer->start(5000);
    }
}

/**
 * @brief NetworkItem::onSendIpConflictDect 延时发送ip冲突检测
 * @param index 本地ip地址索引号
 */
void NetworkItem::onSendIpConflictDect(int index)
{
    QTimer::singleShot(500, this, [ = ]() mutable {
        const QStringList& ipList = currentIpList();
        if (index >= ipList.size()) {
            m_ipConflictChecking = false;
            return;
        }

        m_networkInter->RequestIPConflictCheck(ipList.at(index), "");

        ++index;
        if (ipList.size() > index) {
            emit sendIpConflictDect(index);
        } else {
            m_ipConflictChecking = false;
        }
    });
}

void NetworkItem::onDetectConflict()
{
    // ip冲突时发起主动检测，如果解除则更新状态
    if (currentIpList().size() <= 0 || m_ipConflictChecking) {
        return;
    }

    m_ipConflictChecking = true;
    onSendIpConflictDect();
}

void NetworkItem::getPluginState()
{
    int wiredState = 0;
    int wirelessState = 0;
    int state = 0;
    int temp = 0;
    // 所有设备状态叠加
    QMapIterator<QString, WirelessItem *> iwireless(m_wirelessItems);
    while (iwireless.hasNext()) {
        iwireless.next();
        auto wirelessItem = iwireless.value();
        if (wirelessItem) {
            temp = wirelessItem->getDeviceState();
            state |= temp;
            if ((temp & WirelessItem::Connected) >> 18) {
                m_connectedWirelessDevice.insert(iwireless.key(), wirelessItem);
                connect(wirelessItem, &WirelessItem::apInfoChanged, this, &NetworkItem::refreshIcon);
            } else {
                disconnect(wirelessItem, &WirelessItem::apInfoChanged, this, &NetworkItem::refreshIcon);
                m_connectedWirelessDevice.remove(iwireless.key());

            }
        }
    }
    // 按如下顺序得到当前无线设备状态
    temp = state;
    if (!temp)
        wirelessState = WirelessItem::Unknown;
    temp = state;
    if ((temp & WirelessItem::Disabled) >> 17)
        wirelessState = WirelessItem::Disabled;
    temp = state;
    if ((temp & WirelessItem::Disconnected) >> 19)
        wirelessState = WirelessItem::Disconnected;
    temp = state;
    if ((temp & WirelessItem::Connecting) >> 20)
        wirelessState = WirelessItem::Connecting;
    temp = state;
    if ((temp & WirelessItem::ConnectNoInternet) >> 24)
        wirelessState = WirelessItem::ConnectNoInternet;
    temp = state;
    if ((temp & WirelessItem::Connected) >> 18) {
        wirelessState = WirelessItem::Connected;
    }
    //将无线获取地址状态中显示为连接中状态
    temp = state;
    if ((temp & WirelessItem::ObtainingIP) >> 22) {
        wirelessState = WirelessItem::ObtainingIP;
    }
    //无线正在认证
    temp = state;
    if ((temp & WirelessItem::Authenticating) >> 17) {
        wirelessState = WirelessItem::Authenticating;
    }

    state = 0;
    temp = 0;
    QMapIterator<QString, WiredItem *> iwired(m_wiredItems);
    while (iwired.hasNext()) {
        iwired.next();
        auto wiredItem = iwired.value();
        if (wiredItem) {
            temp = wiredItem->getDeviceState();
            state |= temp;
            if ((temp & WiredItem::Connected) >> 2) {
                m_connectedWiredDevice.insert(iwired.key(), wiredItem);
            } else {
                m_connectedWiredDevice.remove(iwired.key());
            }
        }
    }
    temp = state;
    if (!temp)
        wiredState = WiredItem::Unknown;
    temp = state;
    if ((temp & WiredItem::Nocable) >> 9)
        wiredState = WiredItem::Nocable;
    temp = state;
    if ((temp & WiredItem::Disabled) >> 1)
        wiredState = WiredItem::Disabled;
    temp = state;
    if ((temp & WiredItem::Disconnected) >> 3)
        wiredState = WiredItem::Disconnected;
    temp = state;
    if ((temp & WiredItem::Connecting) >> 4)
        wiredState = WiredItem::Connecting;
    temp = state;
    if ((temp & WiredItem::ConnectNoInternet) >> 8)
        wiredState = WiredItem::ConnectNoInternet;
    temp = state;
    if ((temp & WiredItem::Connected) >> 2) {
        wiredState = WiredItem::Connected;
    }
    //将有线获取地址状态中显示为连接中状态
    temp = state;
    if ((temp & WiredItem::ObtainingIP) >> 6) {
        wiredState = WiredItem::ObtainingIP;
    }
    //有线正在认证
    temp = state;
    if ((temp & WiredItem::Authenticating) >> 5) {
        wiredState = WiredItem::Authenticating;
    }

    switch (wirelessState | wiredState) {
    case 0:
        m_pluginState = Unknow;
        break;
    case 0x00000001:
        m_pluginState = Bdisconnected;
        break;
    case 0x00000002:
        m_pluginState = Bdisabled;
        break;
    case 0x00000004:
        m_pluginState = Bconnected;
        break;
    case 0x00000008:
        m_pluginState = Bdisconnected;
        break;
    case 0x00000010: //有线正在连接
        m_pluginState = Bconnecting;
        break;
    case 0x00000020: //有线正在认证
        m_pluginState = Bconnecting;
        break;
    case 0x00000040: //有线正在获取ip转换为正在连接
        m_pluginState = Bconnecting;
        break;
    case 0x00000080:
        m_pluginState = Bdisconnected;
        break;
    case 0x00000100:
        m_pluginState = BconnectNoInternet;
        break;
    case 0x00000200:
        m_pluginState = Nocable;
        break;
    case 0x00000400: //只有有线,有线失败
        m_pluginState = Bfailed;
        break;
    case 0x00010000:
        m_pluginState = Adisconnected;
        break;
    case 0x00020000:
        m_pluginState = Adisabled;
        break;
    case 0x00040000:
        m_pluginState = Aconnected;
        break;
    case 0x00080000:
        m_pluginState = Adisconnected;
        break;
    case 0x00100000: //无线正在连接
        m_pluginState = Aconnecting;
        break;
    case 0x00200000: //无线正在认证
        m_pluginState = Aconnecting;
        break;
    case 0x00400000: //无线正在获取ip转换为正在连接
        m_pluginState = Aconnecting;
        break;
    case 0x00800000: //无线获取ip失败
        m_pluginState = Adisconnected;
        break;
    case 0x01000000:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000000: // 只有无线 Adisconnected(无线未连接) 改为 Afailed(无线连接失败)
        m_pluginState = Afailed;
        break;
    case 0x00010001:
        m_pluginState = Disconnected;
        break;
    case 0x00020001:
        m_pluginState = Bdisconnected;//
        break;
    case 0x00040001:
        m_pluginState = Aconnected;
        break;
    case 0x00080001:
        m_pluginState = Disconnected;
        break;
    case 0x00100001: //无线正在连接,有线启用
        m_pluginState = Aconnecting;
        break;
    case 0x00200001: //无线正在认证, 有线启用
        m_pluginState = Aconnecting;
        break;
    case 0x00400001: //无线正在获取ip,有线启用
        m_pluginState = Aconnecting;
        break;
    case 0x00800001:
        m_pluginState = Disconnected;
        break;
    case 0x01000001:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000001:
        m_pluginState = Disconnected;
        break;
    case 0x00010002:
        m_pluginState = Adisconnected;
        break;
    case 0x00020002: //有线无线都禁用
        m_pluginState = Disabled;
        break;
    case 0x00040002:
        m_pluginState = Aconnected;
        break;
    case 0x00080002:
        m_pluginState = Adisconnected;
        break;
    case 0x00100002: //无线正在连接,有线禁用
        m_pluginState = Aconnecting;
        break;
    case 0x00200002: //无线正在认证,有线禁用
        m_pluginState = Aconnecting;
        break;
    case 0x00400002: //无线正在获取ip,有线禁用
        m_pluginState = Aconnecting;
        break;
    case 0x00800002: //有线禁用,无线未连接
        m_pluginState = Adisconnected;
        break;
    case 0x01000002:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000002: //有线禁用,无线连接失败,设为无线连接失败  Adisconnected换成Afailed
        m_pluginState = Afailed;
        break;
    case 0x00010004:
        m_pluginState = Bconnected;
        break;
    case 0x00020004:
        m_pluginState = Bconnected;
        break;
    case 0x00040004: //无线已连接,有线已连接
        m_pluginState = Connected;
        break;
    case 0x00080004: //无线断开连接,有线已连接,状态改为有线已连接
        m_pluginState = Bconnected;
        break;
    case 0x00100004: //无线正在连接,有线已连接
        m_pluginState = Aconnecting;
        break;
    case 0x00200004: // 无线认证中.有线已连接,
        m_pluginState = Aconnecting;
        break;
    case 0x00400004: //无线正在获取ip,有线已连接
        m_pluginState = Aconnecting;
        break;
    case 0x00800004:
        m_pluginState = Bconnected;
        break;
    case 0x01000004:
        m_pluginState = Bconnected;
        break;
    case 0x02000004:
        m_pluginState = Bconnected;
        break;
    case 0x00010008:
        m_pluginState = Disconnected;
        break;
    case 0x00020008:
        m_pluginState = Bdisconnected;
        break;
    case 0x00040008: //无线已连接,有线连接失败
        m_pluginState = Aconnected;
        break;
    case 0x00080008:
        m_pluginState = Disconnected;
        break;
    case 0x00100008: //无线正在连接,有线断开连接
        m_pluginState = Aconnecting;
        break;
    case 0x00200008: //无线正在认证,有线断开连接
        m_pluginState = Aconnecting;
        break;
    case 0x00400008:  //无线正在获取ip,有线断开连接
        m_pluginState = Aconnecting;
        break;
    case 0x00800008:
        m_pluginState = Disconnected;
        break;
    case 0x01000008:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000008:
        m_pluginState = Disconnected;
        break;
    case 0x00010010:
        m_pluginState = Bconnecting;
        break;
    case 0x00020010:
        m_pluginState = Bconnecting;
        break;
    case 0x00040010: //有线正在连接, 无线已连接
        m_pluginState = Bconnecting;
        break;
    case 0x00080010:
        m_pluginState = Bconnecting;
        break;
    case 0x00100010: //无线正在连接,有线正在连接
        m_pluginState = Connecting;
        break;
    case 0x00200010: //无线正在认证, 有线正在连接
        m_pluginState = Connecting;
        break;
    case 0x00400010:  //无线正在获取ip,有线正在连接
        m_pluginState = Connecting;
        break;
    case 0x00800010:
        m_pluginState = Bconnecting;
        break;
    case 0x01000010:
        m_pluginState = Bconnecting;
        break;
    case 0x02000010:
        m_pluginState = Bconnecting;
        break;
    case 0x00010020: //有线正在连接认证 ,无线其余操作
        m_pluginState = Bconnecting;
        break;
    case 0x00020020:
        m_pluginState = Bconnecting;
        break;
    case 0x00040020: //有线正在认证, 无线已连接
        m_pluginState = Bconnecting;
        break;
    case 0x00080020://无线断开连接,有线正在认证
        m_pluginState = Bconnecting;
        break;
    case 0x00100020: //无线正在连接,有线正在认证
        m_pluginState = Connecting;
        break;
    case 0x00200020: //无线正在认证,有线正在认证
        m_pluginState = Connecting;
        break;
    case 0x00400020: // 无线正在获取ip ,有线正在认证
        m_pluginState = Connecting;
        break;
    case 0x00800020: // 无线获取ip失败 ,有线正在认证(仅有线正在连接,显示有线状态)
        m_pluginState = Bconnecting;
        break;
    case 0x01000020: // 无线连上但无法访问 ,有线正在认证
        m_pluginState = Bconnecting;
        break;
    case 0x02000020: // 无线连接失败 ,有线正在认证
        m_pluginState = Bconnecting;
        break;
    case 0x00010040: //有线正在获取ip,无线启用
        m_pluginState = Bconnecting;
        break;
    case 0x00020040:
        m_pluginState = Bconnecting;
        break;
    case 0x00040040: //有线正在获取ip,无线已连接
        m_pluginState = Bconnecting;
        break;
    case 0x00080040:
        m_pluginState = Bconnecting;
        break;
    case 0x00100040: //无线正在连接,有线正在获取ip
        m_pluginState = Connecting;
        break;
    case 0x00200040: //无线正在认证.有线正在获取ip
        m_pluginState = Connecting;
        break;
    case 0x00400040: // 无线正在获取ip ,有线正在获取ip
        m_pluginState = Connecting;
        break;
    case 0x00800040: //有线正在获取ip,无线获取ip失败
        m_pluginState = Bconnecting;
        break;
    case 0x01000040:
        m_pluginState = Bconnecting;
        break;
    case 0x02000040:
        m_pluginState = Bconnecting;
        break;
    case 0x00010080:
        m_pluginState = Disconnected;
        break;
    case 0x00020080:
        m_pluginState = Bdisconnected;
        break;
    case 0x00040080:
        m_pluginState = Aconnected;
        break;
    case 0x00080080:
        m_pluginState = Disconnected;
        break;
    case 0x00100080: //无线正在连接,有线获取ip失败
        m_pluginState = Aconnecting;
        break;
    case 0x00200080: //无线正在认证,有线获取ip失败
        m_pluginState = Aconnecting;
        break;
    case 0x00400080: // 无线正在获取ip ,有线获取ip失败
        m_pluginState = Aconnecting;
        break;
    case 0x00800080:
        m_pluginState = Disconnected;
        break;
    case 0x01000080:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000080:
        m_pluginState = Disconnected;
        break;
    case 0x00010100:
        m_pluginState = BconnectNoInternet;
        break;
    case 0x00020100:
        m_pluginState = BconnectNoInternet;
        break;
    case 0x00040100:
        m_pluginState = Aconnected;
        break;
    case 0x00080100:
        m_pluginState = BconnectNoInternet;
        break;
    case 0x00100100: //无线正在连接,有线已连接但无法访问
        m_pluginState = Aconnecting;
        break;
    case 0x00200100: //无线正在认证,有线连接但无法访问
        m_pluginState = Aconnecting;
        break;
    case 0x00400100: // 无线正在获取ip ,有线已连接但无法访问
        m_pluginState = Aconnecting;
        break;
    case 0x00800100:
        m_pluginState = BconnectNoInternet;
        break;
    case 0x01000100:
        m_pluginState = ConnectNoInternet;
        break;
    case 0x02000100:
        m_pluginState = BconnectNoInternet;
        break;
    case 0x00010200:
        m_pluginState = Adisconnected;
        break;
    case 0x00020200:
        m_pluginState = Nocable;
        break;
    case 0x00040200:
        m_pluginState = Aconnected;
        break;
    case 0x00080200:
        m_pluginState = Adisconnected;
        break;
    case 0x00100200: //无线正在连接,未插入网线
        m_pluginState = Aconnecting;
        break;
    case 0x00200200: //无线正在认证,有线未插入网线
        m_pluginState = Aconnecting;
        break;
    case 0x00400200: // 无线正在获取ip ,有线未插入网线
        m_pluginState = Aconnecting;
        break;
    case 0x00800200:
        m_pluginState = Adisconnected;
        break;
    case 0x01000200:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000200:
        m_pluginState = Adisconnected;
        break;
    case 0x00010400:
        m_pluginState = Adisconnected;
        break;
    case 0x00020400: //有线失败,无线禁用
        m_pluginState = Bfailed;
        break;
    case 0x00040400:
        m_pluginState = Aconnected;
        break;
    case 0x00080400:
        m_pluginState = Adisconnected;
        break;
    case 0x00100400: //无线正在连接,有线连接失败
        m_pluginState = Aconnecting;
        break;
    case 0x00200400: //无线正在认证,有线连接失败
        m_pluginState = Aconnecting;
        break;
    case 0x00400400: // 无线正在获取ip ,有线连接失败
        m_pluginState = Aconnecting;
        break;
    case 0x00800400:
        m_pluginState = Adisconnected;
        break;
    case 0x01000400:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000400: //有线,无线都连接失败,改为无线连接失败
        m_pluginState = Failed;
        break;
    }

    switch (m_pluginState) {
    case Unknow:
    case Disabled:
    case Connected:
    case Disconnected:
    case ConnectNoInternet:
    case Adisabled:
    case Bdisabled:
    case Aconnected:
    case Bconnected:
    case Adisconnected:
    case Bdisconnected:
    case Aconnecting:
    case Bconnecting:
    case AconnectNoInternet:
    case BconnectNoInternet:
    case Bfailed:
    case Nocable:
        m_switchWireTimer->stop();
        m_timeOut = true;
        break;
    case Connecting:
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

void NetworkItem::updateView()
{
    // 固定显示高度即为固定示项目数
    const int constDisplayItemCnt = 10;
    auto wirelessCnt = m_wirelessItems.size();

    if (m_switchWirelessBtnState) {
        for (auto wirelessItem : m_wirelessItems) {
            if (wirelessItem && wirelessItem->device()->enabled())
                // 单个设备开关控制项
                wirelessItem->setControlPanelVisible(wirelessCnt != 1);
        }
    }
    // 设备总控开关只与是否有设备相关
    m_wirelessControlPanel->setVisible(wirelessCnt);
    m_wiredControlPanel->setVisible(m_wiredItems.size());

    m_applet->widget()->adjustSize();
    m_applet->setFixedHeight(qMin(m_applet->widget()->height(), constDisplayItemCnt * ITEM_HEIGHT));

    if (m_wirelessControlPanel->isVisible()) {
        if (!m_wirelessScanTimer->isActive())
            m_wirelessScanTimer->start(m_wirelessScanInterval * 1000);
    } else {
        if (m_wirelessScanTimer->isActive())
            m_wirelessScanTimer->stop();
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

int NetworkItem::getStrongestAp()
{
    int retStrength = -1;
    for (auto wirelessItem : m_connectedWirelessDevice) {
        auto apInfo = wirelessItem->getConnectedApInfo();
        if (apInfo.isEmpty())
            continue;
        auto strength = apInfo.value("Strength").toInt();
        if (retStrength < strength)
            retStrength = strength;
    }
    return retStrength;
}

// check if exist at least one available network
// if exist, open dock pannel instead of control center
bool NetworkItem::isExistAvailableNetwork()
{
    for (auto item : m_wirelessItems) {
        if (item->APcount() > 0)
            return true;
    }
    return false;
}

/**
 * @brief 更新有线（无线）适配器的开关状态，并根据开关状态显示设备列表。
 */
void NetworkItem::updateMasterControlSwitch()
{
    m_switchWiredBtnState = false;
    m_switchWirelessBtnState = false;

    /* 获取有线适配器启用状态 */
    for (WiredItem *wiredItem : m_wiredItems) {
        if (wiredItem && wiredItem->deviceEabled()) {
            m_switchWiredBtnState = wiredItem->deviceEabled();
            break;
        }
    }
    /* 更新有线适配器总开关状态（阻塞信号是为了防止重复设置适配器启用状态）*/
    m_switchWiredBtn->blockSignals(true);
    m_switchWiredBtn->setChecked(m_switchWiredBtnState);
    m_thirdSeparator->setVisible(m_switchWiredBtnState);
    m_switchWiredBtn->blockSignals(false);
    // 刷新有线连接,不管是否禁用均显示
    for (WiredItem *wiredItem : m_wiredItems) {
        if (!wiredItem)
            continue;

        m_wiredLayout->addWidget(wiredItem);
    }

    /* 获取无线适配器启用状态 */
    for (auto wirelessItem : m_wirelessItems) {
        if (wirelessItem && wirelessItem->deviceEanbled()) {
            m_switchWirelessBtnState = wirelessItem->deviceEanbled();
            break;
        }
    }
    /* 更新无线适配器总开关状态（阻塞信号是为了防止重复设置适配器启用状态） */
    m_switchWirelessBtn->blockSignals(true);
    m_switchWirelessBtn->setChecked(m_switchWirelessBtnState);
    m_secondSeparator->setVisible(m_switchWirelessBtnState && m_wiredItems.count() > 0);
    m_switchWirelessBtn->blockSignals(false);
    /* 根据无线适配器启用状态增/删布局中的组件 */
    for (WirelessItem *wirelessItem : m_wirelessItems) {
        if (!wirelessItem)
            continue;

        if (m_switchWirelessBtnState) {
            m_wirelessLayout->addWidget(wirelessItem->itemApplet());
        } else {
            m_wirelessLayout->removeWidget(wirelessItem->itemApplet());
        }
        wirelessItem->itemApplet()->setVisible(m_switchWirelessBtnState);
        wirelessItem->setVisible(m_switchWirelessBtnState);
    }

    m_loadingIndicator->setVisible(m_switchWirelessBtnState || m_switchWiredBtnState);
}

void NetworkItem::refreshTips()
{
    if (m_ipConflict) {
        m_tipsWidget->setText(tr("IP conflict"));
        return;
    }

    switch (m_pluginState) {
    case Disabled:
    case Adisabled:
    case Bdisabled:
        m_tipsWidget->setText(tr("Device disabled"));
        break;
    case Connected: {
        QString strTips;
        QStringList textList;
        int wirelessIndex = 1;
        int wireIndex = 1;
        for (auto wirelessItem : m_connectedWirelessDevice) {
            if (wirelessItem) {
                auto info = wirelessItem->getActiveWirelessConnectionInfo();
                if (!info.contains("Ip4"))
                    continue;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    continue;
                if (m_connectedWirelessDevice.size() == 1) {
                    strTips = tr("Wireless connection: %1").arg(ipv4.value("Address").toString()) + '\n';
                } else {
                    strTips = tr("Wireless Network").append(QString("%1").arg(wirelessIndex++)).append(":"+ipv4.value("Address").toString()) + '\n';
                }
                strTips.chop(1);
                textList << strTips;
            }
        }
        for (auto wiredItem : m_connectedWiredDevice) {
            if (wiredItem) {
                auto info = wiredItem->getActiveWiredConnectionInfo();
                if (!info.contains("Ip4"))
                    continue;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    continue;
                if (m_connectedWiredDevice.size() == 1) {
                    strTips = tr("Wired connection: %1").arg(ipv4.value("Address").toString()) + '\n';
                } else {
                    strTips = tr("Wired Network").append(QString("%1").arg(wireIndex++)).append(":"+ipv4.value("Address").toString()) + '\n';
                }
                strTips.chop(1);
                textList << strTips;
            }
        }
        m_tipsWidget->setTextList(textList);
    }
        break;
    case Aconnected: {
        QString strTips;
        int wirelessIndex=1;
        QStringList textList;
        for (auto wirelessItem : m_connectedWirelessDevice) {
            if (wirelessItem) {
                auto info = wirelessItem->getActiveWirelessConnectionInfo();
                if (!info.contains("Ip4"))
                    continue;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    continue;
                if (m_connectedWirelessDevice.size() == 1) {
                    strTips = tr("Wireless connection: %1").arg(ipv4.value("Address").toString()) + '\n';
                } else {
                    strTips = tr("Wireless Network").append(QString("%1").arg(wirelessIndex++)).append(":"+ipv4.value("Address").toString()) + '\n';
                }
                strTips.chop(1);
                textList << strTips;
            }
        }
        m_tipsWidget->setTextList(textList);
    }
        break;
    case Bconnected: {
        QString strTips;
        QStringList textList;
        int wireIndex = 1;
        for (auto wiredItem : m_connectedWiredDevice) {
            if (wiredItem) {
                auto info = wiredItem->getActiveWiredConnectionInfo();
                if (!info.contains("Ip4"))
                    continue;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    continue;
                if (m_connectedWiredDevice.size() == 1) {
                    strTips = tr("Wired connection: %1").arg(ipv4.value("Address").toString()) + '\n';
                } else {
                    strTips = tr("Wired Network").append(QString("%1").arg(wireIndex++)).append(":"+ipv4.value("Address").toString()) + '\n';
                }
                strTips.chop(1);
                textList << strTips;
            }
        }
        m_tipsWidget->setTextList(textList);
    }
        break;
    case Disconnected:
    case Adisconnected:
    case Bdisconnected:
        m_tipsWidget->setText(tr("Not connected"));
        break;
    case Connecting:
    case Aconnecting:
    case Bconnecting: {
        m_tipsWidget->setText(tr("Connecting"));
        return;
    }
    case ConnectNoInternet:
    case AconnectNoInternet:
    case BconnectNoInternet:
        m_tipsWidget->setText(tr("Connected but no Internet access"));
        break;
    case Failed:
    case Afailed:
    case Bfailed:
        m_tipsWidget->setText(tr("Connection failed"));
        break;
    case Unknow:
    case Nocable:
        m_tipsWidget->setText(tr("Network cable unplugged"));
        break;
    }
}

bool NetworkItem::isShowControlCenter()
{
    bool onlyOneTypeDevice = false;
    if ((m_wiredItems.size() == 0 && m_wirelessItems.size() > 0)
            || (m_wiredItems.size() > 0 && m_wirelessItems.size() == 0))
        onlyOneTypeDevice = true;

    // if exist at least one available network
    if (isExistAvailableNetwork())
        return false;

    if (onlyOneTypeDevice) {
        switch (m_pluginState) {
        case Unknow:
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
        case Unknow:
        case Nocable:
        case Bfailed:
        // if has both type device, wired device is disconned
        // and wireless has no ap, should also show control center
        case Adisconnected:
        case ConnectNoInternet:
        case Disconnected:
        case Disabled:
            return true;
        default:
            break;
        }
    }

    return false;
}

const QStringList NetworkItem::currentIpList()
{
    QStringList nativeIpList = QStringList();

    nativeIpList.append(getActiveWiredList());
    nativeIpList.append(getActiveWirelessList());

    return nativeIpList;
}

const QStringList NetworkItem::getActiveWiredList()
{
    QStringList wiredIpList;
    for (auto wiredItem : m_wiredItems.values()) {
        if (wiredItem) {
            auto info = wiredItem->getActiveWiredConnectionInfo();
            if (!info.contains("Ip4"))
                continue;

            const QJsonObject ipv4 = info.value("Ip4").toObject();
            if (!ipv4.contains("Address"))
                continue;

            if (!wiredIpList.contains(ipv4.value("Address").toString()))
                wiredIpList.append(ipv4.value("Address").toString());
        }
    }
    return wiredIpList;
}

const QStringList NetworkItem::getActiveWirelessList()
{
    QStringList wirelessIpList;
    for (auto wirelessItem : m_wirelessItems.values()) {
        if (wirelessItem) {
            auto info = wirelessItem->getActiveWirelessConnectionInfo();
            if (!info.contains("Ip4"))
                continue;

            const QJsonObject ipv4 = info.value("Ip4").toObject();
            if (!ipv4.contains("Address"))
                continue;

            if (!wirelessIpList.contains(ipv4.value("Address").toString()))
                wirelessIpList.append(ipv4.value("Address").toString());
        }
    }
    return wirelessIpList;
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
