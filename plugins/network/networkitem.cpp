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
    , m_switchWire(true)
    , m_timer(new QTimer(this))
    , m_switchWireTimer(new QTimer(this))
{
    setMouseTracking(true);
    m_timer->setInterval(100);

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

    connect(m_switchWireTimer, &QTimer::timeout, [ = ] {
        m_switchWire = !m_switchWire;
    });
    connect(m_timer, &QTimer::timeout, this, &NetworkItem::refreshIcon);
    connect(m_switchWiredBtn, &DSwitchButton::toggled, this, &NetworkItem::wiredsEnable);
    connect(m_switchWirelessBtn, &DSwitchButton::toggled, this, &NetworkItem::wirelessEnable);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &NetworkItem::onThemeTypeChanged);
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
        break;
    case Bconnected:
        stateString = "online";
        iconString = QString("network-%1-symbolic").arg(stateString);
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
        m_timer->start();
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
            const quint64 index = QTime::currentTime().msec() / 10 % 100;
            const int num = (index % 5) + 1;
            iconString = QString("network-wired-symbolic-connecting%1.svg").arg(num);
            if (height() <= PLUGIN_BACKGROUND_MIN_SIZE
                    && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                iconString.append(PLUGIN_MIN_ICON_NAME);
            m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);
            update();
            return;
        }
    }
    case Aconnecting: {
        m_timer->start();
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
        m_timer->start(200);
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
    case Unknow:
    case Nocable:
        stateString = "error";//待图标 暂用错误图标
        iconString = QString("network-%1-symbolic").arg(stateString);
        break;
    }

    m_timer->stop();

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
    return false;
}

QString NetworkItem::getStrengthStateString(int strength)
{
    if (5 >= strength)
        return "0";
    else if (5 < strength && 30 >= strength)
        return "20";
    else if (30 < strength && 55 >= strength)
        return "40";
    else if (55 < strength && 65 >= strength)
        return "60";
    else if (65 < strength)
        return "80";
    else
        return "0";
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
    while (iwireless.hasNext()) {
        iwireless.next();
        auto wirelessItem = iwireless.value();
        if (wirelessItem) {
            temp = wirelessItem->getDeviceState();
            state |= temp;
            if ((temp & WirelessItem::Connected) >> 18) {
                m_connectedWirelessDevice.insert(iwireless.key(), wirelessItem);
            }
        }
    }
    // 按如下顺序得到当前无线设备状态
    temp = state;
    if (!temp)
        wirelessState = WirelessItem::Unknow;
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
            }
        }
    }
    temp = state;
    if (!temp)
        wiredState = WiredItem::Unknow;
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
    case 0x00000010:
        m_pluginState = Bconnecting;
        break;
    case 0x00000020:
        m_pluginState = Bconnecting;
        break;
    case 0x00000040:
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
    case 0x00000400:
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
    case 0x00100000:
        m_pluginState = Aconnecting;
        break;
    case 0x00200000:
        m_pluginState = Aconnecting;
        break;
    case 0x00400000:
        m_pluginState = Aconnecting;
        break;
    case 0x00800000:
        m_pluginState = Adisconnected;
        break;
    case 0x01000000:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000000:
        m_pluginState = Adisconnected;
        break;
    case 0x00010001:
        m_pluginState = Disconnected;
        break;
    case 0x00020001:
        m_pluginState = Bdisconnected;
        break;
    case 0x00040001:
        m_pluginState = Aconnected;
        break;
    case 0x00080001:
        m_pluginState = Disconnected;
        break;
    case 0x00100001:
        m_pluginState = Aconnecting;
        break;
    case 0x00200001:
        m_pluginState = Aconnecting;
        break;
    case 0x00400001:
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
    case 0x00020002:
        m_pluginState = Disabled;
        break;
    case 0x00040002:
        m_pluginState = Aconnected;
        break;
    case 0x00080002:
        m_pluginState = Adisconnected;
        break;
    case 0x00100002:
        m_pluginState = Aconnecting;
        break;
    case 0x00200002:
        m_pluginState = Aconnecting;
        break;
    case 0x00400002:
        m_pluginState = Aconnecting;
        break;
    case 0x00800002:
        m_pluginState = Adisconnected;
        break;
    case 0x01000002:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000002:
        m_pluginState = Adisconnected;
        break;
    case 0x00010004:
        m_pluginState = Bconnected;
        break;
    case 0x00020004:
        m_pluginState = Bconnected;
        break;
    case 0x00040004:
        m_pluginState = Connected;
        break;
    case 0x00080004:
        m_pluginState = Bconnected;
        break;
    case 0x00100004:
        m_pluginState = Bconnected;
        break;
    case 0x00200004:
        m_pluginState = Bconnected;
        break;
    case 0x00400004:
        m_pluginState = Bconnected;
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
    case 0x00040008:
        m_pluginState = Aconnected;
        break;
    case 0x00080008:
        m_pluginState = Disconnected;
        break;
    case 0x00100008:
        m_pluginState = Aconnecting;
        break;
    case 0x00200008:
        m_pluginState = Aconnecting;
        break;
    case 0x00400008:
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
    case 0x00040010:
        m_pluginState = Aconnected;
        break;
    case 0x00080010:
        m_pluginState = Bconnecting;
        break;
    case 0x00100010:
        m_pluginState = Connecting;
        break;
    case 0x00200010:
        m_pluginState = Connecting;
        break;
    case 0x00400010:
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
    case 0x00010020:
        m_pluginState = Bconnecting;
        break;
    case 0x00020020:
        m_pluginState = Bconnecting;
        break;
    case 0x00040020:
        m_pluginState = Aconnected;
        break;
    case 0x00080020:
        m_pluginState = Bconnecting;
        break;
    case 0x00100020:
        m_pluginState = Connecting;
        break;
    case 0x00200020:
        m_pluginState = Connecting;
        break;
    case 0x00400020:
        m_pluginState = Connecting;
        break;
    case 0x00800020:
        m_pluginState = Bconnecting;
        break;
    case 0x01000020:
        m_pluginState = Bconnecting;
        break;
    case 0x02000020:
        m_pluginState = Bconnecting;
        break;
    case 0x00010040:
        m_pluginState = Bconnecting;
        break;
    case 0x00020040:
        m_pluginState = Bconnecting;
        break;
    case 0x00040040:
        m_pluginState = Aconnected;
        break;
    case 0x00080040:
        m_pluginState = Bconnecting;
        break;
    case 0x00100040:
        m_pluginState = Connecting;
        break;
    case 0x00200040:
        m_pluginState = Connecting;
        break;
    case 0x00400040:
        m_pluginState = Connecting;
        break;
    case 0x00800040:
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
    case 0x00100080:
        m_pluginState = Aconnecting;
        break;
    case 0x00200080:
        m_pluginState = Aconnecting;
        break;
    case 0x00400080:
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
    case 0x00100100:
        m_pluginState = Aconnecting;
        break;
    case 0x00200100:
        m_pluginState = Aconnecting;
        break;
    case 0x00400100:
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
    case 0x00100200:
        m_pluginState = Aconnecting;
        break;
    case 0x00200200:
        m_pluginState = Aconnecting;
        break;
    case 0x00400200:
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
    case 0x00020400:
        m_pluginState = Bfailed;
        break;
    case 0x00040400:
        m_pluginState = Aconnected;
        break;
    case 0x00080400:
        m_pluginState = Adisconnected;
        break;
    case 0x00100400:
        m_pluginState = Aconnecting;
        break;
    case 0x00200400:
        m_pluginState = Aconnecting;
        break;
    case 0x00400400:
        m_pluginState = Aconnecting;
        break;
    case 0x00800400:
        m_pluginState = Adisconnected;
        break;
    case 0x01000400:
        m_pluginState = AconnectNoInternet;
        break;
    case 0x02000400:
        m_pluginState = Bfailed;
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
        break;
    case Connecting:
        // 启动2s切换计时
        m_switchWireTimer->start(2000);
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

    // 分割线 都有设备时才有
    auto centralWidget = m_applet->widget();
    centralWidget->setFixedHeight(centralWidget->sizeHint().height());
    m_applet->setFixedHeight(qMin(centralWidget->sizeHint().height(), constDisplayItemCnt * ItemHeight));
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
                    break;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    break;
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
                    break;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    break;
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
                    break;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    break;
                if (m_connectedWiredDevice.size() == 1) {
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
                    break;
                const QJsonObject ipv4 = info.value("Ip4").toObject();
                if (!ipv4.contains("Address"))
                    break;
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
    case Bfailed:
        m_tipsWidget->setText(tr("Network cable unplugged"));
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

void NetworkItem::wirelessScan()
{
    if (m_loadingIndicator->loading())
        return;
    m_loadingIndicator->setLoading(true);
    QTimer::singleShot(1000, this, [ = ] {
        m_loadingIndicator->setLoading(false);
    });
}
