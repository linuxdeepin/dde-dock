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

#include <QIcon>

#define PLUGIN_STATE_KEY    "enable"

PowerPlugin::PowerPlugin(QObject *parent)
    : QObject(parent),

      m_pluginLoaded(false),
      m_tipsLabel(new TipsWidget),
      m_uPowerInter(new QDBusInterface("org.freedesktop.UPower",
                                       "/org/freedesktop/UPower",
                                       "org.freedesktop.UPower",
                                       QDBusConnection::systemBus())),
      m_uBatteryDeviceInter(nullptr)
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setObjectName("power");

    if (!m_uPowerInter->isValid()) {
        qDebug() << "DBusConnection to org.freedesktop.UPower is invalid";
        return;
    }

    QDBusReply<QList<QDBusObjectPath>> reply = m_uPowerInter->call("EnumerateDevices");
    QList<QDBusObjectPath> paths = reply.value();
    QDBusObjectPath batteryPath;

    foreach(auto objectPath, paths) {
        qDebug() << "EnumerateDevices: " << objectPath.path();

        if (objectPath.path().contains("battery")) {
            batteryPath = objectPath;
            break;
        }
    }

    if (batteryPath.path().isEmpty())
        return;
    
    m_uBatteryDeviceInter = new QDBusInterface(
        "org.freedesktop.UPower",
        batteryPath.path(),
        "org.freedesktop.UPower.Device",
        QDBusConnection::systemBus()
    );
    if(!m_uBatteryDeviceInter->isValid()) {
        qDebug() << QString("DBusConnection to %1 is invalid").arg(batteryPath.path());
    }
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
        return m_powerStatusWidget;

    return nullptr;
}

QWidget *PowerPlugin::itemTipsWidget(const QString &itemKey)
{
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();

    if (data.isEmpty()) {
        return nullptr;
    }

    m_tipsLabel->setObjectName(itemKey);

    const uint percentage = qMin(100.0, qMax(0.0, data.value("Display")));
    const QString value = QString("%1%").arg(std::round(percentage));
    const bool charging = !m_powerInter->onBattery();

    if (!charging) {
        qint64 timeToEmpty = -1;
        if(m_uBatteryDeviceInter && m_uBatteryDeviceInter->property("TimeToEmpty").isValid())
            timeToEmpty = m_uBatteryDeviceInter->property("TimeToEmpty").toInt();
        
        m_tipsLabel->setText(
            tr("Remaining Capacity: %1, %2 Until Empty")
                .arg(value)
                .arg(QDateTime::fromTime_t(timeToEmpty).toUTC().toString("hh:mm:ss"))
        );
    } else {
        const int batteryState = m_powerInter->batteryState()["Display"];

        if (batteryState == BatteryState::FULLY_CHARGED || percentage == 100.)
            m_tipsLabel->setText(tr("Charged %1 Battery Is Charged").arg(value));
        else {
            qint64 timeToFull = -1;
            if(m_uBatteryDeviceInter && m_uBatteryDeviceInter->property("TimeToFull").isValid())
                timeToFull = m_uBatteryDeviceInter->property("TimeToFull").toInt();
            m_tipsLabel->setText(
                tr("Charging %1, %2 Until Full")
                    .arg(value)
                    .arg(QDateTime::fromTime_t(timeToFull).toUTC().toString("hh:mm:ss"))
            );
        }
    }

    return m_tipsLabel;
}

void PowerPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void PowerPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool PowerPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

const QString PowerPlugin::itemCommand(const QString &itemKey)
{
    if (itemKey == POWER_KEY)
        return QString("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:power\"");

    return QString();
}

const QString PowerPlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey != POWER_KEY) {
        return QString();
    }

    QList<QVariant> items;
    items.reserve(6);

    QMap<QString, QVariant> power;
    power["itemId"] = "power";
    power["itemText"] = tr("Power settings");
    power["isActive"] = true;
    items.push_back(power);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void PowerPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "power")
        QProcess::startDetached("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:power\"");
}

void PowerPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == POWER_KEY) {
        m_powerStatusWidget->refreshIcon();
    }
}

int PowerPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());

    return m_proxyInter->getValue(this, key, displayMode() == Dock::DisplayMode::Fashion ? 3 : 3).toInt();
}

void PowerPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());

    m_proxyInter->saveValue(this, key, order);
}

void PowerPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

void PowerPlugin::updateBatteryVisible()
{
    const bool exist = !m_powerInter->batteryPercentage().isEmpty();

    if (!exist)
        m_proxyInter->itemRemoved(this, POWER_KEY);
    else if (exist && !pluginIsDisable())
        m_proxyInter->itemAdded(this, POWER_KEY);
}

void PowerPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        qDebug() << "power plugin has been loaded! return";
        return;
    }

    m_pluginLoaded = true;

    m_powerStatusWidget = new PowerStatusWidget;
    m_powerInter = new DBusPower(this);

    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, &PowerPlugin::updateBatteryVisible);

    updateBatteryVisible();
}

void PowerPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, POWER_KEY);
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        updateBatteryVisible();
    }
}
