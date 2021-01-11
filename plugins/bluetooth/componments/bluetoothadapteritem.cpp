/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     chenwei <chenwei@uniontech.com>
 *
 * Maintainer: chenwei <chenwei@uniontech.com>
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

#include "bluetoothadapteritem.h"
#include "componments/adapter.h"
#include "bluetoothapplet.h"
#include "bluetoothconstants.h"

#include <QBoxLayout>
#include <QStandardItemModel>

#include <DFontSizeManager>
#include <DLabel>
#include <DSwitchButton>
#include <DListView>
#include <DSpinner>
#include <DApplicationHelper>

BluetoothDeviceItem::BluetoothDeviceItem(QStyle *style, const Device *device, DListView *parent)
    : m_style(style)
    , m_device(device)
    , m_standarditem(new DStandardItem())
    , m_loading(new DSpinner(parent))
{
    initActionList();
    initConnect();
}

BluetoothDeviceItem::~BluetoothDeviceItem()
{
    if (m_loading != nullptr) {
        delete m_loading;
        m_loading = nullptr;
    }
}

void BluetoothDeviceItem::initActionList()
{
    m_labelAction = new DViewItemAction(Qt::AlignLeft | Qt::AlignVCenter, QSize(), QSize(), false);
    m_stateAction = new DViewItemAction(Qt::AlignLeft | Qt::AlignVCenter, QSize(), QSize(), true);

    m_loading->setFixedSize(QSize(24, 24));
    m_stateAction->setWidget(m_loading);

    m_standarditem->setAccessibleText(m_device->name());
    m_standarditem->setActionList(Qt::RightEdge, {m_stateAction});
    m_standarditem->setActionList(Qt::LeftEdge, {m_labelAction});

    m_labelAction->setText(m_device->name());
    updateDeviceState(m_device->state());
    updateIconTheme(DGuiApplicationHelper::instance()->themeType());
}

void BluetoothDeviceItem::initConnect()
{
    connect(DApplicationHelper::instance(), &DApplicationHelper::themeTypeChanged, this, &BluetoothDeviceItem::updateIconTheme);
    connect(m_device, &Device::stateChanged, this, &BluetoothDeviceItem::updateDeviceState);
}

void BluetoothDeviceItem::updateIconTheme(DGuiApplicationHelper::ColorType type)
{
    if (type == DGuiApplicationHelper::LightType) {
        if (!m_device->deviceType().isEmpty()) {
            m_deviceIcon = LightString.arg(m_device->deviceType());
        } else {
            m_deviceIcon = LightString.arg("other");
        }
    } else {
        if (!m_device->deviceType().isEmpty()) {
            m_deviceIcon = DarkString.arg(m_device->deviceType());
        } else {
            m_deviceIcon = DarkString.arg("other");
        }
    }
    m_labelAction->setIcon(QIcon(m_deviceIcon));
}

void BluetoothDeviceItem::updateDeviceState(Device::State state)
{
    m_labelAction->setText(m_device->name());
    if (state == Device::StateAvailable) {
        m_loading->start();
        m_stateAction->setVisible(true);
        m_standarditem->setCheckState(Qt::Unchecked);
    } else if (state == Device::StateConnected){
        m_loading->stop();
        m_stateAction->setVisible(false);
        m_standarditem->setCheckState(Qt::Checked);
        emit requestTopDeviceItem(m_standarditem);
    } else {
        m_loading->stop();
        m_stateAction->setVisible(false);
        m_standarditem->setCheckState(Qt::Unchecked);
    }
    emit deviceStateChanged(m_device);
}

BluetoothAdapterItem::BluetoothAdapterItem(Adapter *adapter, QWidget *parent)
    : QWidget(parent)
    , m_adapter(adapter)
    , m_adapterLabel(new SettingLabel(adapter->name(), this))
    , m_adapterStateBtn(new DSwitchButton(this))
    , m_deviceListview(new DListView(this))
    , m_deviceModel(new QStandardItemModel(m_deviceListview))
{
    initData();
    initUi();
    initConnect();
}

BluetoothAdapterItem::~BluetoothAdapterItem()
{
    qDeleteAll(m_deviceItems);
}

void BluetoothAdapterItem::onConnectDevice(const QModelIndex &index)
{
    const QStandardItemModel *deviceModel = dynamic_cast<const QStandardItemModel *>(index.model());
    if (!deviceModel)
        return;
    DStandardItem *deviceitem = dynamic_cast<DStandardItem *>(deviceModel->item(index.row()));

    foreach(const auto item, m_deviceItems) {
        if (item->standardItem() == deviceitem) {
            emit connectDevice(item->device(), m_adapter);
        }
    }
}

void BluetoothAdapterItem::onTopDeviceItem(DStandardItem *item)
{
    if (!item || item->row() == -1 || item->row() == 0)
        return;

    int index1 = item->row();
    QStandardItem *index = m_deviceModel->takeItem(index1, 0);
    m_deviceModel->removeRow(index1);
    m_deviceModel->insertRow(0, index);
}

void BluetoothAdapterItem::onAdapterNameChanged(const QString name)
{
    m_adapterLabel->label()->setText(name);
}

int BluetoothAdapterItem::currentDeviceCount()
{
    return m_deviceItems.size();
}

QStringList BluetoothAdapterItem::connectedDevicesName()
{
    QStringList devsName;
    for (BluetoothDeviceItem *devItem : m_deviceItems) {
        if (devItem && devItem->device()->state() == Device::StateConnected) {
            devsName << devItem->device()->name();
        }
    }

    return devsName;
}

void BluetoothAdapterItem::initData()
{
    if (!m_adapter->powered())
        return;

    foreach(const auto device, m_adapter->devices()) {
        if (!m_deviceItems.contains(device->id()))
            onDeviceAdded(device);
    }
    emit deviceCountChanged();
}

void BluetoothAdapterItem::onDeviceAdded(const Device *device)
{
    int insertRow = 0;
    foreach(const auto item, m_deviceItems) {
        if (item->device()->connectState()) {
            insertRow++;
        }
    }

    BluetoothDeviceItem *item = new BluetoothDeviceItem(style(), device, m_deviceListview);
    connect(item, &BluetoothDeviceItem::requestTopDeviceItem, this, &BluetoothAdapterItem::onTopDeviceItem);
    connect(item, &BluetoothDeviceItem::deviceStateChanged, this, &BluetoothAdapterItem::deviceStateChanged);

    m_deviceItems.insert(device->id(), item);
    m_deviceModel->insertRow(insertRow, item->standardItem());
    emit deviceCountChanged();
}

void BluetoothAdapterItem::onDeviceRemoved(const Device *device)
{
    if(m_deviceItems.isEmpty())
        return;

    m_deviceModel->removeRow(m_deviceItems.value(device->id())->standardItem()->row());
    m_deviceItems.value(device->id())->deleteLater();
    m_deviceItems.remove(device->id());
    emit deviceCountChanged();
}

void BluetoothAdapterItem::initUi()
{
    setAccessibleName(m_adapter->name());
    setContentsMargins(0, 0, 0, 0);
    m_adapterLabel->setFixedSize(ItemWidth, TitleHeight);
    m_adapterLabel->addSwichButton(m_adapterStateBtn);
    DFontSizeManager::instance()->bind(m_adapterLabel->label(), DFontSizeManager::T4);

    m_adapterStateBtn->setChecked(m_adapter->powered());

    QVBoxLayout *mainLayout = new QVBoxLayout(this);
    mainLayout->setMargin(0);
    mainLayout->setSpacing(0);

    m_deviceListview->setAccessibleName("DeviceItemList");
    m_deviceListview->setModel(m_deviceModel);
    m_deviceListview->setItemSpacing(1);
    m_deviceListview->setItemSize(QSize(ItemWidth, DeviceItemHeight));
    m_deviceListview->setBackgroundType(DStyledItemDelegate::ClipCornerBackground);
    m_deviceListview->setItemRadius(0);
    m_deviceListview->setEditTriggers(QAbstractItemView::NoEditTriggers);
    m_deviceListview->setSelectionMode(QAbstractItemView::NoSelection);
    m_deviceListview->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_deviceListview->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_deviceListview->setSizeAdjustPolicy(QAbstractScrollArea::AdjustToContents);
    m_deviceListview->setSizePolicy(QSizePolicy::Preferred, QSizePolicy::Preferred);

    mainLayout->addWidget(m_adapterLabel);
    mainLayout->addSpacing(2);
    mainLayout->addWidget(m_deviceListview);
}

void BluetoothAdapterItem::initConnect()
{
    connect(m_adapter, &Adapter::deviceAdded, this, &BluetoothAdapterItem::onDeviceAdded);
    connect(m_adapter, &Adapter::deviceRemoved, this, &BluetoothAdapterItem::onDeviceRemoved);
    connect(m_adapter, &Adapter::nameChanged, this, &BluetoothAdapterItem::onAdapterNameChanged);
    connect(m_deviceListview, &DListView::clicked, this, &BluetoothAdapterItem::onConnectDevice);
    connect(m_adapter, &Adapter::poweredChanged, this, [ = ] (bool state) {
        initData();
        m_deviceListview->setVisible(state);
        m_adapterStateBtn->setChecked(state);
        m_adapterStateBtn->setEnabled(true);
        emit adapterPowerChanged();
    });
    connect(m_adapterStateBtn, &DSwitchButton::clicked, this, [ = ] (bool state){
        qDeleteAll(m_deviceItems);
        m_deviceItems.clear();
        m_deviceModel->clear();
        m_deviceListview->setVisible(false);
        m_adapterStateBtn->setEnabled(false);
        emit requestSetAdapterPower(m_adapter, state);
    });
}
