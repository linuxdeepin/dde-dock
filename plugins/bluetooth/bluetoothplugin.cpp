// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "bluetoothplugin.h"
#include "adaptersmanager.h"
#include "bluetoothmainwidget.h"
#include "imageutil.h"

#include <DGuiApplicationHelper>

#define STATE_KEY  "enable"

DGUI_USE_NAMESPACE

BluetoothPlugin::BluetoothPlugin(QObject *parent)
    : QObject(parent)
    , m_adapterManager(new AdaptersManager(this))
    , m_bluetoothItem(nullptr)
    , m_bluetoothWidget(nullptr)
{
}

const QString BluetoothPlugin::pluginName() const
{
    return "bluetooth";
}

const QString BluetoothPlugin::pluginDisplayName() const
{
    return tr("Bluetooth");
}

void BluetoothPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (m_bluetoothItem)
        return;

    m_bluetoothItem.reset(new BluetoothItem(m_adapterManager));

    m_bluetoothWidget.reset(new BluetoothMainWidget(m_adapterManager));

    connect(m_bluetoothItem.data(), &BluetoothItem::justHasAdapter, [ this ] {
        m_proxyInter->itemAdded(this, BLUETOOTH_KEY);
    });
    connect(m_bluetoothItem.data(), &BluetoothItem::requestHide, [ this ] {
        m_proxyInter->requestSetAppletVisible(this, QUICK_ITEM_KEY, false);
    });
    connect(m_bluetoothItem.data(), &BluetoothItem::noAdapter, [ this ] {
        m_proxyInter->requestSetAppletVisible(this, QUICK_ITEM_KEY, false);
        m_proxyInter->requestSetAppletVisible(this, BLUETOOTH_KEY, false);
        m_proxyInter->itemRemoved(this, BLUETOOTH_KEY);
    });
    connect(m_bluetoothWidget.data(), &BluetoothMainWidget::requestExpand, this, [ this ] {
        m_proxyInter->requestSetAppletVisible(this, QUICK_ITEM_KEY, true);
    });

    if (m_bluetoothItem->hasAdapter())
        m_proxyInter->itemAdded(this, BLUETOOTH_KEY);
}

QWidget *BluetoothPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == BLUETOOTH_KEY) {
        return m_bluetoothItem.data();
    }

    if (itemKey == QUICK_ITEM_KEY)
        return m_bluetoothWidget.data();

    return nullptr;
}

QWidget *BluetoothPlugin::itemTipsWidget(const QString &itemKey)
{
    if (itemKey == BLUETOOTH_KEY) {
        return m_bluetoothItem->tipsWidget();
    }

    return nullptr;
}

QWidget *BluetoothPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == BLUETOOTH_KEY) {
        return m_bluetoothItem->popupApplet();
    }

    if (itemKey == QUICK_ITEM_KEY) {
        return m_bluetoothItem->popupApplet();
    }

    return nullptr;
}

void BluetoothPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    if (itemKey == BLUETOOTH_KEY) {
        m_bluetoothItem->invokeMenuItem(menuId, checked);
    }
}

int BluetoothPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 4).toInt();
}

void BluetoothPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

void BluetoothPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == BLUETOOTH_KEY) {
        m_bluetoothItem->refreshIcon();
    }
}

QIcon BluetoothPlugin::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{
    QString iconFile;
    if (themeType == DGuiApplicationHelper::ColorType::DarkType)
        iconFile = ":/bluetooth-active-symbolic.svg";
    else
        iconFile = ":/bluetooth-active-symbolic-dark.svg";

    switch (dockPart) {
    case DockPart::DCCSetting:
        return ImageUtil::loadSvg(iconFile, QSize(18, 18));
    case DockPart::QuickShow:
        return ImageUtil::loadSvg(iconFile, QSize(18, 16));
    default:
        break;
    }

    return QIcon();
}

PluginsItemInterface::PluginMode BluetoothPlugin::status() const
{
    if (m_bluetoothItem.data()->isPowered())
        return PluginMode::Active;

    return PluginMode::Deactive;
}

QString BluetoothPlugin::description() const
{
    if (m_bluetoothItem.data()->isPowered())
        return tr("Turn on");

    return tr("Turn off");
}

PluginFlags BluetoothPlugin::flags() const
{
    return PluginFlag::Type_Common
            | PluginFlag::Quick_Multi
            | PluginFlag::Attribute_CanDrag
            | PluginFlag::Attribute_CanInsert
            | PluginFlag::Attribute_CanSetting;
}

