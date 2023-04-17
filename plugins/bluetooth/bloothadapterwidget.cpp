// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "bloothadapterwidget.h"
#include "adapter.h"
#include "device.h"

#include <QLabel>
#include <QVBoxLayout>

#define ITEMHEIGHT 45

BloothAdapterWidget::BloothAdapterWidget(Adapter *adapter, QWidget *parent)
    : QWidget(parent)
    , m_adapter(adapter)
    , m_myDeviceLabel(new QLabel(tr("My Devices"), this))
    , m_myDeviceView(new DListView(this))
    , m_otherDeviceLabel(new QLabel(tr("Other Devices"), this))
    , m_otherDeviceView(new DListView(this))
    , m_myDeviceModel(new QStandardItemModel(this))
    , m_otherDeviceModel(new QStandardItemModel(this))
{
    initUi();
    initConnection();
    initDevice();
}

Adapter *BloothAdapterWidget::adapter()
{
    return m_adapter;
}

void BloothAdapterWidget::onDeviceAdded(const Device *device)
{
    if (device->name().isEmpty())
        return;

    DStandardItem *deviceItem = new DStandardItem;
    deviceItem->setData(QVariant::fromValue(const_cast<Device *>(device)), Dtk::UserRole + 1);
    deviceItem->setText(device->name());
    if (device->paired()) {
        // 我的设备
        m_myDeviceModel->insertRow(0, deviceItem);
    } else {
        // 其他设备
        m_otherDeviceModel->insertRow(0, deviceItem);
    }

    updateDeviceVisible();
}

void BloothAdapterWidget::onDeviceRemoved(const Device *device)
{
    auto removeDeviceItem = [ = ](QStandardItemModel *model) {
        for (int i = 0; i < model->rowCount(); i++) {
            Device *tmpDevice = model->item(i)->data(Dtk::UserRole + 1).value<Device *>();
            if (tmpDevice == device) {
                model->removeRow(i);
                return true;
            }
        }

        return false;
    };

    if (!removeDeviceItem(m_myDeviceModel))
        removeDeviceItem(m_otherDeviceModel);

    updateDeviceVisible();
}

void BloothAdapterWidget::onDeviceNameUpdated(const Device *device) const
{
    auto findDeviceItem = [ = ](QStandardItemModel *model)->DStandardItem * {
        for (int i = 0; i < model->rowCount(); i++) {
            DStandardItem *item = static_cast<DStandardItem *>(model->item(i));
            Device *tmpDevice = item->data(Dtk::UserRole + 1).value<Device *>();
            if (tmpDevice == device) {
                return item;
            }
        }

        return nullptr;
    };
    DStandardItem *item = findDeviceItem(m_myDeviceModel);
    if (!item)
        item = findDeviceItem(m_otherDeviceModel);
    if (item)
        item->setText(device->name());
}

void BloothAdapterWidget::onPoweredChanged(const bool powered)
{
    initDevice();
    updateDeviceVisible();
}

void BloothAdapterWidget::onOtherClicked(const QModelIndex &index)
{
    Device *device = index.data(Dtk::UserRole + 1).value<Device *>();
    if (!device || device->state() == Device::State::StateConnected)
        return;

    if ((device->deviceType() == "audio-headset" || device->deviceType() == "audio-headphones")
        && device->state() == Device::State::StateAvailable) {
        return;
    }

    Q_EMIT requestConnectDevice(device);
}

void BloothAdapterWidget::initUi()
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->setSpacing(0);
    mainLayout->addWidget(m_myDeviceLabel);
    mainLayout->addWidget(m_myDeviceView);
    mainLayout->addSpacing(20);
    mainLayout->addWidget(m_otherDeviceLabel);
    mainLayout->addSpacing(6);
    mainLayout->addWidget(m_otherDeviceView);

    m_myDeviceLabel->setVisible(false);
    m_myDeviceView->setVisible(false);
    m_myDeviceView->setModel(m_myDeviceModel);
    m_myDeviceView->setFixedHeight(0);
    m_myDeviceView->setItemSpacing(5);

    m_otherDeviceLabel->setVisible(false);
    m_otherDeviceView->setVisible(false);
    m_otherDeviceView->setModel(m_otherDeviceModel);
    m_otherDeviceView->setFixedHeight(0);
    m_otherDeviceView->setItemSpacing(5);
}

void BloothAdapterWidget::initConnection()
{
    connect(m_adapter, &Adapter::deviceAdded, this, &BloothAdapterWidget::onDeviceAdded);
    connect(m_adapter, &Adapter::deviceRemoved, this, &BloothAdapterWidget::onDeviceRemoved);
    connect(m_adapter, &Adapter::deviceNameUpdated, this, &BloothAdapterWidget::onDeviceNameUpdated);
    connect(m_adapter, &Adapter::poweredChanged, this, &BloothAdapterWidget::onPoweredChanged);

    connect(m_otherDeviceView, &DListView::clicked, this, &BloothAdapterWidget::onOtherClicked);
}

void BloothAdapterWidget::initDevice()
{
    m_myDeviceModel->clear();
    m_otherDeviceModel->clear();
    QMap<QString, const Device *> devices = m_adapter->devices();
    for (auto it = devices.begin(); it != devices.end(); it++)
        onDeviceAdded(it.value());
}

void BloothAdapterWidget::adjustHeight()
{
    int height = m_myDeviceView->height() + 20 + m_otherDeviceView->height() + 5;

    if (m_myDeviceLabel->isVisible())
        height += m_myDeviceLabel->height();
    if (m_otherDeviceLabel->isVisible())
        height += m_otherDeviceLabel->height();

    setFixedHeight(height);
}

void BloothAdapterWidget::updateDeviceVisible()
{
    bool powered = m_adapter->powered();
    if (powered) {
        m_myDeviceLabel->setVisible(m_myDeviceModel->rowCount() > 0);
        m_myDeviceView->setVisible(m_myDeviceModel->rowCount() > 0);
        m_myDeviceView->setFixedHeight(std::min(m_myDeviceModel->rowCount(), 10) * ITEMHEIGHT);

        m_otherDeviceLabel->setVisible(m_adapter->powered() && m_otherDeviceModel->rowCount() > 0);
        m_otherDeviceView->setVisible(m_adapter->powered() && m_otherDeviceModel->rowCount() > 0);
        m_otherDeviceView->setFixedHeight(std::min(m_otherDeviceModel->rowCount(), 10) * ITEMHEIGHT);
    } else {
        m_myDeviceLabel->setVisible(false);
        m_myDeviceView->setVisible(false);
        m_myDeviceView->setFixedHeight(0);

        m_otherDeviceLabel->setVisible(false);
        m_otherDeviceView->setVisible(false);
        m_otherDeviceView->setFixedHeight(0);
    }

    adjustHeight();
    Q_EMIT requestUpdate();

}
