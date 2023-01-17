/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
 *
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

    connect(m_bluetoothItem.data(), &BluetoothItem::justHasAdapter, [&] {
        m_proxyInter->itemAdded(this, BLUETOOTH_KEY);
    });
    connect(m_bluetoothItem.data(), &BluetoothItem::noAdapter, [&] {
        m_proxyInter->itemRemoved(this, BLUETOOTH_KEY);
    });
    connect(m_bluetoothWidget.data(), &BluetoothMainWidget::requestExpand, this, [ = ] {
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
    if (dockPart == DockPart::QuickPanel)
        return QIcon();

    QString iconFile;
    if (themeType == DGuiApplicationHelper::ColorType::DarkType)
        iconFile = ":/bluetooth-active-symbolic.sv";
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
        r