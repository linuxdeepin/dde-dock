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

#include "shutdownplugin.h"
#include "dbus/dbusaccount.h"

#include <QIcon>
#include <QSettings>

#define PLUGIN_STATE_KEY    "enable"

ShutdownPlugin::ShutdownPlugin(QObject *parent)
    : QObject(parent),

      m_settings("deepin", "dde-dock-power"),
      m_shutdownWidget(new PluginWidget),
      m_powerStatusWidget(new PowerStatusWidget),
      m_tipsLabel(new TipsWidget),

      m_powerInter(new DBusPower(this))
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setObjectName("power");

    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, &ShutdownPlugin::updateBatteryVisible);
    connect(m_shutdownWidget, &PluginWidget::requestContextMenu, this, &ShutdownPlugin::requestContextMenu);
    connect(m_powerStatusWidget, &PowerStatusWidget::requestContextMenu, this, &ShutdownPlugin::requestContextMenu);
}

const QString ShutdownPlugin::pluginName() const
{
    return "shutdown";
}

const QString ShutdownPlugin::pluginDisplayName() const
{
    return tr("Power");
}

QWidget *ShutdownPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == SHUTDOWN_KEY)
        return m_shutdownWidget;
    if (itemKey == POWER_KEY)
        return m_powerStatusWidget;

    return nullptr;
}

QWidget *ShutdownPlugin::itemTipsWidget(const QString &itemKey)
{
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();
    m_tipsLabel->setObjectName(itemKey);

    if (data.isEmpty() || (itemKey == SHUTDOWN_KEY && displayMode() == Dock::Efficient))
    {
        m_tipsLabel->setText(tr("Shut down"));
        return m_tipsLabel;
    }

    const uint percentage = qMin(100.0, qMax(0.0, data.value("Display")));
    const QString value = QString("%1%").arg(std::round(percentage));
    const bool charging = !m_powerInter->onBattery();
    if (!charging)
        m_tipsLabel->setText(tr("Remaining Capacity %1").arg(value));
    else
    {
        const int batteryState = m_powerInter->batteryState()["Display"];

        if (batteryState == BATTERY_FULL || percentage == 100.)
            m_tipsLabel->setText(tr("Charged %1").arg(value));
        else
            m_tipsLabel->setText(tr("Charging %1").arg(value));
    }

    return m_tipsLabel;
}

void ShutdownPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable())
        delayLoader();
}

void ShutdownPlugin::pluginStateSwitched()
{
    m_settings.setValue(PLUGIN_STATE_KEY, !m_settings.value(PLUGIN_STATE_KEY, true).toBool());

    if (pluginIsDisable())
    {
        m_proxyInter->itemRemoved(this, SHUTDOWN_KEY);
        m_proxyInter->itemRemoved(this, POWER_KEY);
    } else {
        m_proxyInter->itemAdded(this, SHUTDOWN_KEY);
        updateBatteryVisible();
    }
}

bool ShutdownPlugin::pluginIsDisable()
{
    return !m_settings.value(PLUGIN_STATE_KEY, true).toBool();
}

const QString ShutdownPlugin::itemCommand(const QString &itemKey)
{
    if (itemKey == SHUTDOWN_KEY)
        return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
    if (itemKey == POWER_KEY)
        return QString("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:power\"");

    return QString();
}

const QString ShutdownPlugin::itemContextMenu(const QString &itemKey)
{
    QList<QVariant> items;
    items.reserve(6);

    const Dock::DisplayMode mode = displayMode();

    if (mode == Dock::Fashion || itemKey == SHUTDOWN_KEY)
    {
        QMap<QString, QVariant> shutdown;
        shutdown["itemId"] = "Shutdown";
        shutdown["itemText"] = tr("Shut down");
        shutdown["isActive"] = true;
        items.push_back(shutdown);

        QMap<QString, QVariant> reboot;
        reboot["itemId"] = "Restart";
        reboot["itemText"] = tr("Restart");
        reboot["isActive"] = true;
        items.push_back(reboot);

        QMap<QString, QVariant> suspend;
        suspend["itemId"] = "Suspend";
        suspend["itemText"] = tr("Suspend");
        suspend["isActive"] = true;
        items.push_back(suspend);

        QMap<QString, QVariant> lock;
        lock["itemId"] = "Lock";
        lock["itemText"] = tr("Lock");
        lock["isActive"] = true;
        items.push_back(lock);

        QMap<QString, QVariant> logout;
        logout["itemId"] = "Logout";
        logout["itemText"] = tr("Log out");
        logout["isActive"] = true;
        items.push_back(logout);

        if (DBusAccount().userList().count() > 1)
        {
            QMap<QString, QVariant> switchUser;
            switchUser["itemId"] = "SwitchUser";
            switchUser["itemText"] = tr("Switch account");
            switchUser["isActive"] = true;
            items.push_back(switchUser);
        }
    }

    if (mode == Dock::Fashion || itemKey == POWER_KEY)
    {
        QMap<QString, QVariant> power;
        power["itemId"] = "power";
        power["itemText"] = tr("Power settings");
        power["isActive"] = true;
        items.push_back(power);
    }

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void ShutdownPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "power")
        QProcess::startDetached("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:power\"");
    else if (menuId == "Lock")
        QProcess::startDetached("dbus-send", QStringList() << "--print-reply"
                                                           << "--dest=com.deepin.dde.lockFront"
                                                           << "/com/deepin/dde/lockFront"
                                                           << QString("com.deepin.dde.lockFront.Show"));
    else
        QProcess::startDetached("dbus-send", QStringList() << "--print-reply"
                                                           << "--dest=com.deepin.dde.shutdownFront"
                                                           << "/com/deepin/dde/shutdownFront"
                                                           << QString("com.deepin.dde.shutdownFront.%1").arg(menuId));
}

void ShutdownPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    m_shutdownWidget->update();

    updateBatteryVisible();
}

int ShutdownPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    return m_settings.value(key, 0).toInt();
}

void ShutdownPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    m_settings.setValue(key, order);
}

void ShutdownPlugin::updateBatteryVisible()
{
    const bool exist = !m_powerInter->batteryPercentage().isEmpty();

    if (!exist || displayMode() == Dock::Fashion)
        m_proxyInter->itemRemoved(this, POWER_KEY);
    else if (exist && !pluginIsDisable())
        m_proxyInter->itemAdded(this, POWER_KEY);
}

void ShutdownPlugin::requestContextMenu(const QString &itemKey)
{
    m_proxyInter->requestContextMenu(this, itemKey);
}

void ShutdownPlugin::delayLoader()
{
    static int retryTimes = 0;

    ++retryTimes;

    if (m_powerInter->isValid() || retryTimes > 10)
    {
        qDebug() << "load power item, dbus valid:" << m_powerInter->isValid();

        m_proxyInter->itemAdded(this, SHUTDOWN_KEY);
        displayModeChanged(displayMode());
    } else {
        qDebug() << "load power failed, wait and retry" << retryTimes;

        // wait and retry
        QTimer::singleShot(1000, this, &ShutdownPlugin::delayLoader);
    }
}
