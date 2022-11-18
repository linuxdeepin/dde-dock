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

#include "powerplugin.h"
#include "dbus/dbusaccount.h"
#include "../widgets/tipswidget.h"
#include "../frame/util/utils.h"

#include <QIcon>
#include <QGSettings>

#include <DDBusSender>

#define PLUGIN_STATE_KEY    "enable"
#define DELAYTIME           (20 * 1000)

using namespace Dock;
static QGSettings *GSettingsByApp()
{
    static QGSettings settings("com.deepin.dde.dock.module.power");
    return &settings;
}

PowerPlugin::PowerPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_showTimeToFull(true)
    , m_powerStatusWidget(nullptr)
    , m_tipsLabel(new TipsWidget)
    , m_systemPowerInter(nullptr)
    , m_powerInter(nullptr)
    , m_preChargeTimer(new QTimer(this))
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setObjectName("power");
    m_preChargeTimer->setInterval(DELAYTIME);
    m_preChargeTimer->setSingleShot(true);
    connect(m_preChargeTimer,&QTimer::timeout,this,&PowerPlugin::refreshTipsData);
}

const QString PowerPlugin::pluginName() const
{
    return "power";
}

const QString PowerPlugin::pluginDisplayName() const
{
    return tr("Power");
}

QWidget *PowerPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == POWER_KEY)
        return m_powerStatusWidget.data();

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
}

const QString PowerPlugin::itemCommand(const QString &itemKey)
{
    if (itemKey == POWER_KEY)
        return QString("dbus-send --print-reply --dest=org.deepin.dde.ControlCenter1 /org/deepin/dde/ControlCenter1 org.deepin.dde.ControlCenter1.ShowPage \"string:power\"");

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
    // 电池插件不显示在快捷面板上，因此此处返回空图标
    static QIcon batteryIcon;
    const QPixmap pixmap = m_powerStatusWidget->getBatteryIcon(themeType);
    batteryIcon.detach();
    batteryIcon.addPixmap(pixmap);
    return batteryIcon;
}

PluginFlags PowerPlugin::flags() const
{
    // 电池插件只在任务栏上面展示，不在快捷面板展示，并且可以拖动，可以在其前面插入其他插件，可以在控制中心设置是否显示隐藏
    return PluginFlag::Type_Common
            | PluginFlag::Attribute_CanDrag
            | PluginFlag::Attribute_CanInsert
            | PluginFlag::Attribute_CanSetting;
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
        m_proxyInter->itemUpdate(this, POWER_KEY);
    });

    m_powerInter = new DBusPower(this);

    m_systemPowerInter = new SystemPowerInter("org.deepin.dde.Power1", "/org/deepin/dde/Power1", QDBusConnection::systemBus(), this);
    m_systemPowerInter->setSync(true);

    connect(GSettingsByApp(), &QGSettings::changed, this, &PowerPlugin::onGSettingsChanged);
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

    if (GSettingsByApp()->keys().contains("showtimetofull")) {
        const bool isEnable = GSettingsByApp()->keys().contains("showtimetofull") && GSettingsByApp()->get("showtimetofull").toBool();
        m_showTimeToFull = isEnable && GSettingsByApp()->get("showtimetofull").toBool();
    }

    refreshTipsData();
}

void PowerPlugin::refreshTipsData()
{
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();
    const uint percentage = qMin(100.0, qMax(0.0, data.value("Display")));
    const QString value = QString("%1%").arg(std::round(percentage));
    const int batteryState = m_powerInter->batteryState()["Display"];

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
