// SPDX-FileCopyrightText: 2024 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later
#include "notificationplugin.h"

#include <DGuiApplicationHelper>

#include <QIcon>
#include <QSettings>
#include <QPainter>

Q_LOGGING_CATEGORY(qLcPluginNotification, "dock.plugin.notification")

#define PLUGIN_STATE_KEY    "enable"
#define TOGGLE_DND          "toggle-dnd"

DGUI_USE_NAMESPACE
using namespace Dock;

NotificationPlugin::NotificationPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_notification(nullptr)
    , m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setText(tr("Notification"));
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setAccessibleName("Notification");
    m_tipsLabel->setObjectName("NotificationTipsLabel");
}

QWidget *NotificationPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey)
    return m_notification.data();
}

QWidget *NotificationPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);
    return m_tipsLabel.data();
}

void NotificationPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void NotificationPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());
    refreshPluginItemsVisible();
}

bool NotificationPlugin::pluginIsDisable()
{
    return !(m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());
}

const QString NotificationPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);
    return QString("dbus-send --session --print-reply --dest=org.deepin.dde.Widgets1 /org/deepin/dde/Widgets1 org.deepin.dde.Widgets1.Toggle");
}

const QString NotificationPlugin::itemContextMenu(const QString &itemKey)
{
    QList<QVariant> items;
    QMap<QString, QVariant> toggleDnd;
    toggleDnd["itemId"] = TOGGLE_DND;
    toggleDnd["itemText"] = tr("Do Not Disturb");
    toggleDnd["isCheckable"] = true;
    toggleDnd["isActive"] = true;
    toggleDnd["checked"] = m_notification->dndMode();
    items.push_back(toggleDnd);
    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = true;
    menu["singleCheck"] = false;
    return QJsonDocument::fromVariant(menu).toJson();
}

void NotificationPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    if (menuId == TOGGLE_DND) {
        m_notification->setDndMode(checked);
    }
}

int NotificationPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    return m_proxyInter->getValue(this, key, 3).toInt();
}

void NotificationPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    m_proxyInter->saveValue(this, key, order);
}

void NotificationPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

QIcon NotificationPlugin::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{
    Q_UNUSED(themeType)
    if (dockPart == DockPart::DCCSetting)
        return QIcon::fromTheme("notification");
    return m_notification->icon();
}

void NotificationPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        return;
    }
    m_pluginLoaded = true;
    m_notification.reset(new Notification);
    connect(m_notification.data(), &Notification::iconRefreshed, this, [this]() { m_proxyInter->itemUpdate(this, pluginName()); });
    m_proxyInter->itemAdded(this, pluginName());
}

void NotificationPlugin::refreshPluginItemsVisible()
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

void NotificationPlugin::refreshIcon(const QString &itemKey)
{
    Q_UNUSED(itemKey)
    m_notification->refreshIcon();
}
