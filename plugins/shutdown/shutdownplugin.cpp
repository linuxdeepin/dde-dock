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

      m_pluginLoaded(false),
      m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setVisible(false);
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
    Q_UNUSED(itemKey);

    return m_shutdownWidget;
}

QWidget *ShutdownPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    // reset text every time to avoid size of LabelWidget not change after
    // font size be changed in ControlCenter
    m_tipsLabel->setText(tr("Power"));

    return m_tipsLabel;
}

void ShutdownPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    // transfer config
    QSettings settings("deepin", "dde-dock-shutdown");
    if (QFile::exists(settings.fileName())) {
        QFile::remove(settings.fileName());
    }

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void ShutdownPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());

    if (pluginIsDisable())
    {
        m_proxyInter->itemRemoved(this, pluginName());
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        m_proxyInter->itemAdded(this, pluginName());
    }
}

bool ShutdownPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

const QString ShutdownPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
}

const QString ShutdownPlugin::itemContextMenu(const QString &itemKey)
{
    QList<QVariant> items;
    items.reserve(6);

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

#ifndef DISABLE_POWER_OPTIONS
    QMap<QString, QVariant> suspend;
    suspend["itemId"] = "Suspend";
    suspend["itemText"] = tr("Suspend");
    suspend["isActive"] = true;
    items.push_back(suspend);

    if (checkSwap()) {
        QMap<QString, QVariant> hibernate;
        hibernate["itemId"] = "Hibernate";
        hibernate["itemText"] = tr("Hibernate");
        hibernate["isActive"] = true;
        items.push_back(hibernate);
    }
#endif

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

#ifndef DISABLE_POWER_OPTIONS
    QMap<QString, QVariant> power;
    power["itemId"] = "power";
    power["itemText"] = tr("Power settings");
    power["isActive"] = true;
    items.push_back(power);
#endif

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

    if (!pluginIsDisable()) {
        m_shutdownWidget->update();
    }
}

int ShutdownPlugin::itemSortKey(const QString &itemKey)
{
    Dock::DisplayMode mode = displayMode();
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(mode);

    if (mode == Dock::DisplayMode::Fashion) {
        return m_proxyInter->getValue(this, key, 2).toInt();
    } else {
        return m_proxyInter->getValue(this, key, 5).toInt();
    }
}

void ShutdownPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    m_proxyInter->saveValue(this, key, order);
}

void ShutdownPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        qDebug() << "shutdown plugin has been loaded! return";
        return;
    }

    m_pluginLoaded = true;

    m_shutdownWidget = new PluginWidget;

    m_proxyInter->itemAdded(this, pluginName());
    displayModeChanged(displayMode());
}

bool ShutdownPlugin::checkSwap()
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
