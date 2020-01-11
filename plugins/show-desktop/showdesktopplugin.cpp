/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     wangshaojun <wangshaojun_cm@deepin.com>
 *
 * Maintainer: wangshaojun <wangshaojun_cm@deepin.com>
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

#include "showdesktopplugin.h"

#include <QIcon>

#define PLUGIN_STATE_KEY    "enable"

ShowDesktopPlugin::ShowDesktopPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setObjectName("show-desktop");
}

const QString ShowDesktopPlugin::pluginName() const
{
    return "show-desktop";
}

const QString ShowDesktopPlugin::pluginDisplayName() const
{
    return tr("Show Desktop");
}

QWidget *ShowDesktopPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_showDesktopWidget;
}

QWidget *ShowDesktopPlugin::itemTipsWidget(const QString &itemKey)
{
    m_tipsLabel->setObjectName(itemKey);

    m_tipsLabel->setText(pluginDisplayName());

    return m_tipsLabel;
}

void ShowDesktopPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void ShowDesktopPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool ShowDesktopPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

const QString ShowDesktopPlugin::itemCommand(const QString &itemKey)
{
    if (itemKey == pluginName())
        QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");

    return QString();
}

const QString ShowDesktopPlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey != pluginName()) {
        return QString();
    }

    QList<QVariant> items;
    items.reserve(6);

    QMap<QString, QVariant> desktop;
    desktop["itemId"] = "show-desktop";
    desktop["itemText"] = tr("Show Desktop");
    desktop["isActive"] = true;
    items.push_back(desktop);

    QMap<QString, QVariant> power;
    power["itemId"] = "remove";
    power["itemText"] = tr("Undock");
    power["isActive"] = true;
    items.push_back(power);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void ShowDesktopPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "show-desktop") {
        QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");
    } else if (menuId == "remove") {
        pluginStateSwitched();
    }
}

void ShowDesktopPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == pluginName()) {
        m_showDesktopWidget->refreshIcon();
    }
}

int ShowDesktopPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 1).toInt();
}

void ShowDesktopPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

void ShowDesktopPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

PluginsItemInterface::PluginType ShowDesktopPlugin::type()
{
    return PluginType::Fixed;
}

void ShowDesktopPlugin::updateVisible()
{
    if (pluginIsDisable())
        m_proxyInter->itemRemoved(this, pluginName());
    else
        m_proxyInter->itemAdded(this, pluginName());
}

void ShowDesktopPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        return;
    }

    m_pluginLoaded = true;

    m_showDesktopWidget = new ShowDesktopWidget;

    m_proxyInter->itemAdded(this, pluginName());

    updateVisible();
}

void ShowDesktopPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, pluginName());
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        updateVisible();
    }
}
