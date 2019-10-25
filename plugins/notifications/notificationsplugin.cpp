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

#include "notificationsplugin.h"

#include <QDBusConnectionInterface>
#include <QIcon>
#include <QSettings>

#define PLUGIN_STATE_KEY    "enable"

NotificationsPlugin::NotificationsPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setVisible(false);
    getNotifyInterface();
}

QDBusInterface *NotificationsPlugin::getNotifyInterface()
{
    if (!m_interface && QDBusConnection::sessionBus().interface()->isServiceRegistered("com.deepin.dde.Notification"))
        m_interface = new QDBusInterface("com.deepin.dde.Notification", "/com/deepin/dde/Notification", "com.deepin.dde.Notification");

    return m_interface;
}

const QString NotificationsPlugin::pluginName() const
{
    return "notifications";
}

const QString NotificationsPlugin::pluginDisplayName() const
{
    return tr("Notifications");
}

QWidget *NotificationsPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_itemWidget;
}

QWidget *NotificationsPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    if (!getNotifyInterface())
        return nullptr;

    QDBusMessage msg = m_interface->call("recordCount");
    uint recordCount = msg.arguments()[0].toUInt();

    if (recordCount)
        m_tipsLabel->setText(QString(tr("%1 Notifications")).arg(recordCount));
    else
        m_tipsLabel->setText(tr("Notifications"));

    return m_tipsLabel;
}

void NotificationsPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void NotificationsPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());

    refreshPluginItemsVisible();
}

bool NotificationsPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

const QString NotificationsPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    if (getNotifyInterface())
        m_interface->call("Toggle");

    return "";
}

void NotificationsPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    if (!pluginIsDisable()) {
        m_itemWidget->update();
    }
}

int NotificationsPlugin::itemSortKey(const QString &itemKey)
{
    Dock::DisplayMode mode = displayMode();
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(mode);

    if (mode == Dock::DisplayMode::Fashion) {
        return m_proxyInter->getValue(this, key, 2).toInt();
    } else {
        return m_proxyInter->getValue(this, key, 5).toInt();
    }
}

void NotificationsPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    m_proxyInter->saveValue(this, key, order);
}

void NotificationsPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

void NotificationsPlugin::loadPlugin()
{
    if (m_pluginLoaded)
        return;

    m_pluginLoaded = true;

    m_itemWidget = new NotificationsWidget;

    m_proxyInter->itemAdded(this, pluginName());
    displayModeChanged(displayMode());
}

bool NotificationsPlugin::checkSwap()
{
    QFile file("/proc/swaps");
    if (file.open(QIODevice::Text | QIODevice::ReadOnly)) {
        const QString &body = file.readAll();
        file.close();
        QRegularExpression re("\\spartition\\s");
        QRegularExpressionMatch match = re.match(body);
        return match.hasMatch();
    }

    return false;
}

void NotificationsPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, pluginName());
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        m_proxyInter->itemAdded(this, pluginName());
    }
}
