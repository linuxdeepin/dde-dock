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

#include "multitaskingplugin.h"
#include "../widgets/tipswidget.h"

#include <QIcon>

#define PLUGIN_STATE_KEY    "enable"

using namespace Dock;
MultitaskingPlugin::MultitaskingPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setObjectName("multitasking");
}

const QString MultitaskingPlugin::pluginName() const
{
    return "multitasking";
}

const QString MultitaskingPlugin::pluginDisplayName() const
{
    return tr("Multitasking View");
}

QWidget *MultitaskingPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_multitaskingWidget;
}

QWidget *MultitaskingPlugin::itemTipsWidget(const QString &itemKey)
{
    m_tipsLabel->setObjectName(itemKey);

    m_tipsLabel->setText(pluginDisplayName());

    return m_tipsLabel;
}

void MultitaskingPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void MultitaskingPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool MultitaskingPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

const QString MultitaskingPlugin::itemCommand(const QString &itemKey)
{
    if (itemKey == PLUGIN_KEY)
        return "dbus-send --session --dest=com.deepin.wm --print-reply /com/deepin/wm com.deepin.wm.PerformAction int32:1";

    return "";
}

const QString MultitaskingPlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey != PLUGIN_KEY) {
        return QString();
    }

    QList<QVariant> items;
    items.reserve(6);

    QMap<QString, QVariant> desktop;
    desktop["itemId"] = "multitasking";
    desktop["itemText"] = tr("Multitasking View");
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

void MultitaskingPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "multitasking") {
        QProcess::startDetached("dbus-send --session --dest=com.deepin.wm --print-reply /com/deepin/wm com.deepin.wm.PerformAction int32:1");
    } else if (menuId == "remove") {
        pluginStateSwitched();
    }
}

void MultitaskingPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == PLUGIN_KEY) {
        m_multitaskingWidget->refreshIcon();
    }
}

int MultitaskingPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 2).toInt();
}

void MultitaskingPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

void MultitaskingPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

PluginsItemInterface::PluginType MultitaskingPlugin::type()
{
    return PluginType::Fixed;
}

void MultitaskingPlugin::updateVisible()
{
    if (pluginIsDisable())
        m_proxyInter->itemRemoved(this, PLUGIN_KEY);
    else
        m_proxyInter->itemAdded(this, PLUGIN_KEY);
}

void MultitaskingPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        return;
    }

    m_pluginLoaded = true;

    m_multitaskingWidget = new MultitaskingWidget;

    m_proxyInter->itemAdded(this, pluginName());

    updateVisible();
}

void MultitaskingPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, PLUGIN_KEY);
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        updateVisible();
    }
}
