// Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "showdesktopplugin.h"
#include "../widgets/tipswidget.h"

#include <QIcon>
#include <QDebug>
#include <DDBusSender>

using namespace Dock;
ShowDesktopPlugin::ShowDesktopPlugin(QObject *parent)
    : QObject(parent)
    , m_showDesktopWidget(nullptr)
    , m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setAccessibleName("show-desktop");
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

    return m_showDesktopWidget.data();
}

QWidget *ShowDesktopPlugin::itemTipsWidget(const QString &itemKey)
{
    m_tipsLabel->setObjectName(itemKey);

    m_tipsLabel->setText(pluginDisplayName());

    return m_tipsLabel.data();
}

void ShowDesktopPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    m_showDesktopWidget.reset(new ShowDesktopWidget);
    m_proxyInter->itemAdded(this, pluginName());
}

const QString ShowDesktopPlugin::itemCommand(const QString &itemKey)
{
    if (itemKey == pluginName())
        QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle", QStringList());

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
        QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle", QStringList());
    } else if (menuId == "remove") {
        // m_proxyInter->itemRemoved(this, pluginName());
        DDBusSender()
            .service("org.deepin.dde.Dock1")
            .interface("org.deepin.dde.Dock1")
            .path("/org/deepin/dde/Dock1")
            .method(QString("setItemOnDock"))
            .arg(QString("Dock_Quick_Plugins"))
            .arg(QString("show-desktop"))
            .arg(false)
            .call();
    }
}

void ShowDesktopPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == pluginName()) {
        m_showDesktopWidget->refreshIcon();
    }
}

QIcon ShowDesktopPlugin::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{

    if (dockPart == DockPart::DCCSetting) {
        return QIcon::fromTheme("dcc-show-desktop", QIcon(":/icons/icons/dcc-show-desktop.svg"));
    }

    return QIcon();
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

PluginsItemInterface::PluginType ShowDesktopPlugin::type()
{
    return PluginType::Fixed;
}

PluginFlags ShowDesktopPlugin::flags() const
{
    return PluginFlag::Type_Fixed | Attribute_CanSetting;
}
