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

#include "adapteritem.h"
#include "switchitem.h"
#include "deviceitem.h"
#include "adapter.h"
#include "adaptersmanager.h"

#include <DDBusSender>

const int Width = 200;

AdapterItem::AdapterItem(AdaptersManager *adapterManager, Adapter *adapter, QWidget *parent)
    : QScrollArea(parent)
    , m_centralWidget(new QWidget(this))
    , m_line(new HorizontalSeparator(this))
    , m_devGoupName(new QLabel(this))
    , m_deviceLayout(new QVBoxLayout)
    , m_openControlCenter(new MenueItem(this))
    , m_adaptersManager(adapterManager)
    , m_adapter(adapter)
    , m_switchItem(new SwitchItem(this))
{
    m_centralWidget->setFixedWidth(Width);
    m_line->setVisible(true);
    m_devGoupName->setText(tr("My Device"));
    m_devGoupName->setVisible(false);
    m_deviceLayout->setMargin(0);
    m_deviceLayout->setSpacing(0);
    m_openControlCenter->setText("Bluetooth settings");
    m_openControlCenter->setVisible(false);
    m_switchItem->setTitle(adapter->name());
    m_switchItem->setChecked(adapter->powered());

    m_deviceLayout->addWidget(m_switchItem);
    m_deviceLayout->addWidget(m_line);
    m_deviceLayout->addWidget(m_devGoupName);
    m_deviceLayout->addWidget(m_openControlCenter);
    m_centralWidget->setFixedWidth(Width);
    m_centralWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Preferred);
    m_centralWidget->setLayout(m_deviceLayout);

    setFixedWidth(Width);
    setWidget(m_centralWidget);
    setFrameShape(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);

    auto myDevices = adapter->devices();
    for (auto constDevice : myDevices) {
        auto device = const_cast<Device *>(constDevice);
        if (device) {
            createDeviceItem(device);
        }
    }

    connect(m_switchItem, &SwitchItem::checkedChanged, this, &AdapterItem::showAndConnect);
    connect(adapter, &Adapter::nameChanged, m_switchItem, &SwitchItem::setTitle);
    connect(adapter, &Adapter::deviceAdded, this, &AdapterItem::addDeviceItem);
    connect(adapter, &Adapter::deviceRemoved, this, &AdapterItem::removeDeviceItem);
    connect(adapter, &Adapter::poweredChanged, m_switchItem, &SwitchItem::setChecked);
    connect(m_openControlCenter, &MenueItem::clicked, []{
        DDBusSender()
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method(QString("ShowModule"))
        .arg(QString("bluetooth"))
        .call();
    });

    showDevices(adapter->powered());
}

int AdapterItem::pairedDeviceCount()
{
    return  m_pairedDeviceItems.size();
}

int AdapterItem::deviceCount()
{
    return m_deviceItems.size();
}

void AdapterItem::setPowered(bool powered)
{
    m_switchItem->setChecked(powered);
}

void AdapterItem::deviceItemPaired(const bool paired)
{
    auto device = qobject_cast<Device *>(sender());
    if (device) {
        auto deviceId = device->id();
        auto deviceItem = m_deviceItems.value(deviceId);
        if (deviceItem) {
            if (paired)
                m_pairedDeviceItems[deviceId] = deviceItem;
            else
                m_pairedDeviceItems.remove(deviceId);
        }
        showDevices(m_adapter->powered());
    }
}

void AdapterItem::removeDeviceItem(const Device *device)
{
    if (!device)
        return;

    auto deviceItem = m_deviceItems.value(device->id());
    if (deviceItem) {
        m_deviceItems.remove(device->id());
        if (device->paired()) {
            m_pairedDeviceItems.remove(device->id());
            m_deviceLayout->removeWidget(deviceItem);
        }
        delete deviceItem;
        showDevices(m_adapter->powered());
    }
}

void AdapterItem::showAndConnect(bool change)
{
    showDevices(change);

    m_adaptersManager->setAdapterPowered(m_adapter, change);
    if (change) {
        m_adaptersManager->connectAllPairedDevice(m_adapter);
    }
    emit powerChanged(change);
}

void AdapterItem::addDeviceItem(const Device *constDevice)
{
    auto device = const_cast<Device *>(constDevice);
    if (!device)
        return;

    createDeviceItem(device);
    showDevices(m_adapter->powered());
}

void AdapterItem::createDeviceItem(Device *device)
{
    if (!device)
        return;

    auto paired = device->paired();
    auto deviceId = device->id();
    auto deviceItem = new DeviceItem(device->name(), this);
    deviceItem->setDevice(device);
    m_deviceItems[deviceId] = deviceItem;
    if (paired)
        m_pairedDeviceItems[deviceId] = deviceItem;
    deviceItem->setVisible(paired);

    connect(device, &Device::pairedChanged, this, &AdapterItem::deviceItemPaired);
    connect(device, &Device::nameChanged, deviceItem, &DeviceItem::setTitle);
    connect(device, &Device::stateChanged, deviceItem, &DeviceItem::chaneState);
    connect(device, &Device::stateChanged, this, &AdapterItem::deviceStateChanged);
    connect(deviceItem, &DeviceItem::clicked, m_adaptersManager, &AdaptersManager::connectDevice);
}

void AdapterItem::updateView()
{   
    auto contentHeight = m_centralWidget->sizeHint().height();
    m_centralWidget->setFixedHeight(contentHeight);
    setFixedHeight(contentHeight);
    emit sizeChange();
}

void AdapterItem::showDevices(bool change)
{
    for (auto deviceItem : m_pairedDeviceItems) {
        if (change)
            m_deviceLayout->addWidget(deviceItem);
        else {
            m_deviceLayout->removeWidget(deviceItem);
        }
        deviceItem->setVisible(change);
    }
    auto itemCount = m_pairedDeviceItems.size();
    m_line->setVisible(change);
    m_devGoupName->setVisible(itemCount && change);
    m_openControlCenter->setVisible(!itemCount);
    updateView();
}
