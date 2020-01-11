/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#include "onboardplugin.h"

#include <QIcon>
#include <QSettings>

#define PLUGIN_STATE_KEY    "enable"

OnboardPlugin::OnboardPlugin(QObject *parent)
    : QObject(parent),

      m_pluginLoaded(false),
      m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setText(tr("Onboard"));
    m_tipsLabel->setVisible(false);
}

const QString OnboardPlugin::pluginName() const
{
    return "onboard";
}

const QString OnboardPlugin::pluginDisplayName() const
{
    return tr("Onboard");
}

QWidget *OnboardPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_onboardItem;
}

QWidget *OnboardPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_tipsLabel;
}

void OnboardPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void OnboardPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool OnboardPlugin::pluginIsDisable()
{
    return !(m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());
}

const QString OnboardPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return QString("dbus-send --print-reply --dest=org.onboard.Onboard /org/onboard/Onboard/Keyboard org.onboard.Onboard.Keyboard.ToggleVisible");
}

const QString OnboardPlugin::itemContextMenu(const QString &itemKey)
{
    QList<QVariant> items;

    QMap<QString, QVariant> onboardSettings;
    onboardSettings["itemId"] = "onboard-settings";
    onboardSettings["itemText"] = tr("Settings");
    onboardSettings["isActive"] = true;
    items.push_back(onboardSettings);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void OnboardPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "onboard-settings")
        QProcess::startDetached("onboard-settings");
}

void OnboardPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    if (!pluginIsDisable()) {
        m_onboardItem->update();
    }
}

int OnboardPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    return m_proxyInter->getValue(this, key, 3).toInt();
}

void OnboardPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    m_proxyInter->saveValue(this, key, order);
}

void OnboardPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

void OnboardPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        qDebug() << "onboard plugin has been loaded! return";
        return;
    }

    m_pluginLoaded = true;

    m_onboardItem = new OnboardItem;

    m_proxyInter->itemAdded(this, pluginName());
    displayModeChanged(displayMode());
}

void OnboardPlugin::refreshPluginItemsVisible()
{
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
