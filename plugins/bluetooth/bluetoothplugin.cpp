// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "bluetoothplugin.h"

#define STATE_KEY  "enable"

BluetoothPlugin::BluetoothPlugin(QObject *parent)
    : QObject(parent),
      m_bluetoothItem(nullptr)
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

    m_bluetoothItem.reset(new BluetoothItem);

    connect(m_bluetoothItem.data(), &BluetoothItem::justHasAdapter, [&] {
        m_enableState = true;
        refreshPluginItemsVisible();
    });
    connect(m_bluetoothItem.data(), &BluetoothItem::noAdapter, [&] {
        m_enableState = false;
        refreshPluginItemsVisible();
    });

    m_enableState = m_bluetoothItem->hasAdapter();

    if (!pluginIsDisable())
        m_proxyInter->itemAdded(this, BLUETOOTH_KEY);
}

void BluetoothPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool BluetoothPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, STATE_KEY, m_enableState).toBool();
}

QWidget *BluetoothPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == BLUETOOTH_KEY) {
        return m_bluetoothItem.data();
    }

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

    return nullptr;
}

const QString BluetoothPlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey == BLUETOOTH_KEY) {
        return m_bluetoothItem->contextMenu();
    }

    return QString();
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

void BluetoothPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

void BluetoothPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable())
        m_proxyInter->itemRemoved(this, BLUETOOTH_KEY);
    else
        m_proxyInter->itemAdded(this, BLUETOOTH_KEY);
}
