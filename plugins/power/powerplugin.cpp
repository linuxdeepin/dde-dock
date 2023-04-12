// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "powerplugin.h"
#include "dbus/dbusaccount.h"
#include "../widgets/tipswidget.h"
#include "../frame/util/utils.h"

#include <QIcon>
#include <QGSettings>
#include <QHBoxLayout>

#include <DFontSizeManager>
#include <DDBusSender>
#include <DConfig>

#define PLUGIN_STATE_KEY    "enable"
#define DELAYTIME           (20 * 1000)

DCORE_USE_NAMESPACE
using namespace Dock;

PowerPlugin::PowerPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_showTimeToFull(true)
    , m_powerStatusWidget(nullptr)
    , m_tipsLabel(new TipsWidget)
    , m_systemPowerInter(nullptr)
    , m_powerInter(nullptr)
    , m_dconfig(new DConfig(QString("org.deepin.dde.dock.power"), QString()))
    , m_preChargeTimer(new QTimer(this))
    , m_quickPanel(nullptr)
{
    initUi();
    initConnection();
}

const QString PowerPlugin::pluginName() const
{
    return "power";
}

const QString PowerPlugin::pluginDisplayName() const
{
    return tr("Battery");
}

QWidget *PowerPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == POWER_KEY)
        return m_powerStatusWidget.data();
    if (itemKey == QUICK_ITEM_KEY)
        return m_quickPanel;

    return nullptr;
}

QWidget *PowerPlugin::itemTipsWidget(const QString &itemKey)
{
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();

    if (data.isEmpty()) {
        return nullptr;
    }

    m_tipsLabel->setObjectName(itemKey);

    refreshTipsData();

    return m_tipsLabel.data();
}

void PowerPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    loadPlugin();

    onThemeTypeChanged(DGuiApplicationHelper::instance()->themeType());
}

const QString PowerPlugin::itemCommand(const QString &itemKey)
{
    if (itemKey == POWER_KEY)
        return QString("dbus-send --print-reply --dest=org.deepin.dde.ControlCenter1 /org/deepin/dde/ControlCenter1 org.deepin.dde.ControlCenter1.ShowPage string:power");

    return QString();
}

void PowerPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "power") {
        DDBusSender()
        .service("org.deepin.dde.ControlCenter1")
        .interface("org.deepin.dde.ControlCenter1")
        .path("/org/deepin/dde/ControlCenter1")
        .method(QString("ShowPage"))
        .arg(QString("power"))
        .call();
     }
}

void PowerPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == POWER_KEY) {
        m_powerStatusWidget->refreshIcon();
    }
}

int PowerPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 5).toInt();
}

void PowerPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

QIcon PowerPlugin::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{
    const QPixmap pixmap = m_powerStatusWidget->getBatteryIcon(themeType);
    static QIcon batteryIcon;
    batteryIcon.detach();
    batteryIcon.addPixmap(pixmap);
    return batteryIcon;
}

PluginFlags PowerPlugin::flags() const
{
    // 电池插件在任务栏上面展示，在快捷面板展示，并且可以拖动，可以在其前面插入其他插件，能在控制中心设置是否显示隐藏
    return PluginFlag::Type_Common
            | PluginFlag::Attribute_CanDrag
            | PluginFlag::Attribute_CanInsert
            | PluginFlag::Attribute_CanSetting
            | PluginFlag::Quick_Single;
}

void PowerPlugin::updateBatteryVisible()
{
    const bool exist = !m_powerInter->batteryPercentage().isEmpty();

    if (exist)
        m_proxyInter->itemAdded(this, POWER_KEY);
    else
        m_proxyInter->itemRemoved(this, POWER_KEY);
}

void PowerPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        qDebug() << "power plugin has been loaded! return";
        return;
    }

    m_pluginLoaded = true;

    m_powerStatusWidget.reset(new PowerStatusWidget);

    connect(m_powerStatusWidget.get(), &PowerStatusWidget::iconChanged, this, [ this ] {
        m_proxyInter->updateDockInfo(this, DockPart::QuickPanel);
        m_proxyInter->updateDockInfo(this, DockPart::QuickShow);
        m_proxyInter->itemUpdate(this, POWER_KEY);
    });

    m_powerInter = new DBusPower(this);

    m_systemPowerInter = new SystemPowerInter("org.deepin.dde.Power1", "/org/deepin/dde/Power1", QDBusConnection::systemBus(), this);
    m_systemPowerInter->setSync(true);

    connect(m_dconfig, &DConfig::valueChanged, this, &PowerPlugin::onGSettingsChanged);
    connect(m_systemPowerInter, &SystemPowerInter::BatteryStatusChanged, [&](uint  value) {
        if (value == BatteryState::CHARGING)
            m_preChargeTimer->start();
        refreshTipsData();
    });
    connect(m_systemPowerInter, &SystemPowerInter::BatteryTimeToEmptyChanged, this, &PowerPlugin::refreshTipsData);
    connect(m_systemPowerInter, &SystemPowerInter::BatteryTimeToFullChanged, this, &PowerPlugin::refreshTipsData);

    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, &PowerPlugin::updateBatteryVisible);

    updateBatteryVisible();

    onGSettingsChanged("showtimetofull");
}

void PowerPlugin::onGSettingsChanged(const QString &key)
{
    if (key != "showtimetofull") {
        return;
    }

    if (m_dconfig->isValid())
        m_showTimeToFull = m_dconfig->value("showtimetofull").toBool();

    refreshTipsData();
}

void PowerPlugin::refreshTipsData()
{
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();
    const uint percentage = qMin(100.0, qMax(0.0, data.value("Display")));
    QString value = QString("%1%").arg(std::round(percentage));
    const int batteryState = m_powerInter->batteryState()["Display"];
    QFontMetrics ftm(m_labelText->font());
    value = ftm.elidedText(value, Qt::TextElideMode::ElideMiddle, m_labelText->width());
    m_labelText->setText(value);
    if (m_preChargeTimer->isActive() && m_showTimeToFull) {
        // 插入电源后，20秒内算作预充电时间，此时计算剩余充电时间是不准确的
        QString tips = tr("Capacity %1 ...").arg(value);
        m_tipsLabel->setText(tips);
        return;
    }

    if (batteryState == BatteryState::DIS_CHARGING || batteryState == BatteryState::NOT_CHARGED || batteryState == BatteryState::UNKNOWN) {
        QString tips;
        qulonglong timeToEmpty = m_systemPowerInter->batteryTimeToEmpty();
        QDateTime time = QDateTime::fromTime_t(timeToEmpty).toUTC();
        uint hour = time.toString("hh").toUInt();
        uint min = time.toString("mm").toUInt();
        uint sencond = time.toString("ss").toInt();
        if (sencond > 0)
            min += 1;
        if (!m_showTimeToFull) {
            tips = tr("Capacity %1").arg(value);
        } else {
            if (hour == 0) {
                tips = tr("Capacity %1, %2 min remaining").arg(value).arg(min);
            } else {
                tips = tr("Capacity %1, %2 hr %3 min remaining").arg(value).arg(hour).arg(min);
            }
        }
        m_tipsLabel->setText(tips);
    } else if (batteryState == BatteryState::FULLY_CHARGED || percentage == 100.) {
        m_tipsLabel->setText(tr("Capacity %1, fully charged").arg(value));
    } else {
        qulonglong timeToFull = m_systemPowerInter->batteryTimeToFull();
        QDateTime time = QDateTime::fromTime_t(timeToFull).toUTC();
        uint hour = time.toString("hh").toUInt();
        uint min = time.toString("mm").toUInt();
        QString tips;
        if (!m_showTimeToFull) {
            tips = tr("Charging %1").arg(value);
        } else {
            if (timeToFull == 0) {  // 电量已充満或电量计算中,剩余充满时间会返回0
                tips = tr("Capacity %1 ...").arg(value);
            } else {
                hour == 0 ? tips = tr("Charging %1, %2 min until full").arg(value).arg(min)
                          : tips = tr("Charging %1, %2 hr %3 min until full").arg(value).arg(hour).arg(min);
            }
        }
        m_tipsLabel->setText(tips);
    }
}

void PowerPlugin::onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType)
{
    const QPixmap pixmap = m_powerStatusWidget->getBatteryIcon(themeType);
    m_imageLabel->setPixmap(pixmap);
}

void PowerPlugin::initUi()
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setObjectName("power");
    m_preChargeTimer->setInterval(DELAYTIME);
    m_preChargeTimer->setSingleShot(true);

    m_quickPanel = new QWidget();
    QVBoxLayout *layout = new QVBoxLayout(m_quickPanel);
    layout->setAlignment(Qt::AlignVCenter);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->setSpacing(0);
    m_imageLabel = new QLabel(m_quickPanel);
    m_imageLabel->setObjectName("imageLabel");
    m_imageLabel->setFixedHeight(24);
    m_imageLabel->setAlignment(Qt::AlignCenter);

    m_labelText = new QLabel(m_quickPanel);
    m_labelText->setObjectName("textLabel");
    m_labelText->setFixedHeight(15);
    m_labelText->setAlignment(Qt::AlignCenter);
    m_labelText->setFont(Dtk::Widget::DFontSizeManager::instance()->t10());
    layout->addWidget(m_imageLabel);
    layout->addSpacing(7);
    layout->addWidget(m_labelText);
}

void PowerPlugin::initConnection()
{
    connect(m_preChargeTimer,&QTimer::timeout,this,&PowerPlugin::refreshTipsData);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &PowerPlugin::onThemeTypeChanged);
}
