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
#include "bluetoothconstants.h"

#include <DDBusSender>

extern void initFontColor(QWidget *widget);

AdapterItem::AdapterItem(AdaptersManager *adapterManager, Adapter *adapter, QWidget *parent)
    : QScrollArea(parent)
    , m_centralWidget(new QWidget(this))
    , m_line(new HorizontalSeparator(this))
    , m_deviceLayout(new QVBoxLayout)
    , m_openControlCenter(new MenueItem(this))
    , m_adaptersManager(adapterManager)
    , m_adapter(adapter)
    , m_switchItem(new SwitchItem(this))
{
    m_centralWidget->setFixedWidth(POPUPWIDTH);
    m_line->setVisible(true);
    m_deviceLayout->setMargin(0);
    m_deviceLayout->setSpacing(0);
    m_openControlCenter->setText(tr("Bluetooth settings"));
    initFontColor(m_openControlCenter);
    m_openControlCenter->setFixedHeight(ITEMHEIGHT);
    m_openControlCenter->setVisible(false);
    m_switchItem->setTitle(adapter->name());
    m_switchItem->setChecked(adapter->powered(),false);

    m_deviceLayout->addWidget(m_switchItem);
    m_deviceLayout->addWidget(m_line);
    m_deviceLayout->addWidget(m_openControlCenter);
    m_centralWidget->setFixedWidth(POPUPWIDTH);
    m_centralWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Preferred);
    m_centralWidget->setLayout(m_deviceLayout);

    setFixedWidth(POPUPWIDTH);
    setWidget(m_centralWidget);
    setFrameShape(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);

    auto myDevices = adapter->devices();
    for (auto constDevice : myDevices) {
        auto device = const_cast<Device *>(constDevice);
        if (device)
            createDeviceItem(device);
    }

    m_initDeviceState = Device::StateUnavailable;
    for (auto constDevice : myDevices) {
        auto device = const_cast<Device *>(constDevice);
        if (device) {
            if (device->state() == Device::StateAvailable) {
                m_initDeviceState = Device::StateConnected;
                continue;
            }
            if (device->state() == Device::StateConnected) {
                m_initDeviceState = Device::StateConnected;
                break;
            }
        }
    }

    connect(m_switchItem, &SwitchItem::checkedChanged, this, &AdapterItem::showAndConnect);
    connect(adapter, &Adapter::nameChanged, m_switchItem, &SwitchItem::setTitle);
    connect(adapter, &Adapter::deviceAdded, this, &AdapterItem::addDeviceItem);
    connect(adapter, &Adapter::deviceRemoved, this, &AdapterItem::removeDeviceItem);
    connect(adapter, &Adapter::poweredChanged, m_switchItem, [=](const bool powered){
        m_switchItem->setChecked(powered,false);
    });
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

int AdapterItem::deviceCount()
{
    return m_deviceItems.size();
}

void AdapterItem::setPowered(bool powered)
{
    m_switchItem->setChecked(powered,true);
}

bool AdapterItem::isPowered()
{
    return m_switchItem->checkState();
}

int AdapterItem::viewHeight()
{
    return m_openControlCenter->isVisible() ? CONTROLHEIGHT + ITEMHEIGHT : CONTROLHEIGHT;
}

void AdapterItem::deviceItemPaired(const bool paired)
{
    auto device = qobject_cast<Device *>(sender());
    if (device) {
        auto deviceItem = m_deviceItems.value(device->id());
        if (paired) {
            m_sortUnConnect.removeOne(deviceItem);
            m_sortConnected << deviceItem;
        } else {
            m_sortConnected.removeOne(deviceItem);
            m_sortUnConnect << deviceItem;
        }
        showDevices(m_adapter->powered());
    }
}

void AdapterItem::deviceRssiChanged()
{
    auto device = qobject_cast<Device *>(sender());
    if (device) {
        auto deviceItem = m_deviceItems.value(device->id());
        auto state = device->state();
        if (deviceItem && Device::StateConnected == state)
            qSort(m_sortConnected);
        else
            qSort(m_sortUnConnect);
        moveDeviceItem(state, deviceItem);
    }
}

void AdapterItem::removeDeviceItem(const Device *device)
{
    if (!device)
        return;

    auto deviceItem = m_deviceItems.value(device->id());
    if (deviceItem) {
        m_deviceItems.remove(device->id());
        m_sortConnected.removeOne(deviceItem);
        m_sortUnConnect.removeOne(deviceItem);
        m_deviceLayout->removeWidget(deviceItem);
        delete deviceItem;
    }
    showDevices(m_adapter->powered());
}

void AdapterItem::showAndConnect(bool change)
{
    m_adaptersManager->setAdapterPowered(m_adapter, change);
    showDevices(change);
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

void AdapterItem::deviceChangeState(const Device::State state)
{
    auto device = qobject_cast<Device *>(sender());
    if (device) {
        auto deviceItem = m_deviceItems.value(device->id());
        if (deviceItem) {
            switch (state) {
            case Device::StateUnavailable: {
                int index = m_sortUnConnect.indexOf(deviceItem);
                if (index < 0) {
                    m_sortConnected.removeOne(deviceItem);
                    m_sortUnConnect << deviceItem;
                    qSort(m_sortUnConnect);
                    moveDeviceItem(state, deviceItem);
                }
            }
                break;
            case Device::StateAvailable:
                break;
            case Device::StateConnected: {
                int index = m_sortConnected.indexOf(deviceItem);
                if (index < 0) {
                    m_sortUnConnect.removeOne(deviceItem);
                    m_sortConnected << deviceItem;
                    qSort(m_sortConnected);
                    moveDeviceItem(state, deviceItem);
                }
            }
                break;
            }
        }
        m_currentDeviceState = state;
        emit deviceStateChanged(state);
    }
}

void AdapterItem::moveDeviceItem(Device::State state, DeviceItem *item)
{
    int size = m_sortConnected.size();
    int index = 0;
    switch (state) {
    case Device::StateUnavailable:
    case Device::StateAvailable: {
        index = m_sortUnConnect.indexOf(item);
        index += size;
    }
        break;
    case Device::StateConnected: {
        index = m_sortUnConnect.indexOf(item);
    }
        break;
    }
    index += 2;
    m_deviceLayout->removeWidget(item);
    m_deviceLayout->insertWidget(index, item);
}

void AdapterItem::createDeviceItem(Device *device)
{
    if (!device)
        return;

    auto deviceId = device->id();
    auto deviceItem = new DeviceItem(device->name(), this);
    deviceItem->setDevice(device);
    m_deviceItems[deviceId] = deviceItem;
    if (device->state() == Device::StateConnected)
        m_sortConnected << deviceItem;
    else
        m_sortUnConnect << deviceItem;

    connect(device, &Device::pairedChanged, this, &AdapterItem::deviceItemPaired);
    connect(device, &Device::nameChanged, deviceItem, &DeviceItem::setTitle);
    connect(device, &Device::stateChanged, deviceItem, &DeviceItem::changeState);
    connect(device, &Device::stateChanged, this, &AdapterItem::deviceChangeState);
    connect(device, &Device::rssiChanged, this, &AdapterItem::deviceRssiChanged);
    connect(deviceItem, &DeviceItem::clicked, m_adaptersManager, &AdaptersManager::connectDevice);
}

void AdapterItem::updateView()
{
    auto contentHeight = m_switchItem->height();
    contentHeight += (m_deviceLayout->count() - 3) * ITEMHEIGHT;
    m_centralWidget->setFixedHeight(contentHeight);
    setFixedHeight(contentHeight);
    emit sizeChange();
}

void AdapterItem::showDevices(bool change)
{
    if (m_sortConnected.size())
        qSort(m_sortConnected);
    if (m_sortUnConnect.size())
        qSort(m_sortUnConnect);

    QList<DeviceItem *> deviceItems;
    deviceItems << m_sortConnected << m_sortUnConnect;

    for (DeviceItem *deviceItem : deviceItems) {
        if (change)
            m_deviceLayout->addWidget(deviceItem);
        else
            m_deviceLayout->removeWidget(deviceItem);
        deviceItem->setVisible(change);
    }

    auto itemCount = m_deviceItems.size();
    m_line->setVisible(change);
    m_openControlCenter->setVisible(!itemCount);
    updateView();
}
