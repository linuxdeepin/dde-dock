// Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "devcollaborationwidget.h"
#include "collaborationdevmodel.h"
#include "devitemdelegate.h"

#include <DStyle>

#include <QMap>
#include <QTimer>
#include <QLabel>
#include <QVBoxLayout>
#include <QStandardItemModel>

#define TITLE_HEIGHT 30
#define ITEM_WIDTH 310
#define ITEM_HEIGHT 36
#define LISTVIEW_ITEM_SPACE 5
#define PER_DEGREE 14

DevCollaborationWidget::DevCollaborationWidget(QWidget *parent)
    : QWidget(parent)
    , m_deviceModel(new CollaborationDevModel(this))
    , m_deviceListView(new DListView(this))
    , m_viewItemModel(new QStandardItemModel(m_deviceListView))
    , m_refreshTimer(new QTimer(this))
{
    initUI();
    loadDevice();

    connect(m_deviceModel, &CollaborationDevModel::devicesChanged, this, &DevCollaborationWidget::loadDevice);
    connect(m_deviceListView, &DListView::clicked, this, &DevCollaborationWidget::itemClicked);
    connect(m_refreshTimer, &QTimer::timeout, this, &DevCollaborationWidget::refreshViewItem);
}

void DevCollaborationWidget::showEvent(QShowEvent *event)
{
    m_deviceModel->checkServiceValid();

    QWidget::showEvent(event);
}

void DevCollaborationWidget::resizeEvent(QResizeEvent *event)
{
    Q_EMIT sizeChanged();

    QWidget::resizeEvent(event);
}

void DevCollaborationWidget::initUI()
{
    m_deviceListView->setModel(m_viewItemModel);

    QLabel *title = new QLabel(tr("PC collaboration"), this);
    title->setFixedHeight(TITLE_HEIGHT);

    QHBoxLayout *hLayout = new QHBoxLayout();
    hLayout->setContentsMargins(10, 0, 0, 0);
    hLayout->addWidget(title);

    QVBoxLayout *mainLayout = new QVBoxLayout();
    mainLayout->setMargin(0);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->setSpacing(0);
    mainLayout->addLayout(hLayout);
    mainLayout->addWidget(m_deviceListView);

    setLayout(mainLayout);

    m_deviceListView->setContentsMargins(0, 0, 0, 0);
    m_deviceListView->setFrameShape(QFrame::NoFrame);
    m_deviceListView->setVerticalScrollBarPolicy(Qt::ScrollBarAsNeeded);
    m_deviceListView->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_deviceListView->setVerticalScrollMode(QAbstractItemView::ScrollPerPixel);
    m_deviceListView->setResizeMode(QListView::Adjust);
    m_deviceListView->setViewportMargins(0, 0, 0, 0);
    m_deviceListView->setSpacing(LISTVIEW_ITEM_SPACE);
    m_deviceListView->setEditTriggers(QAbstractItemView::NoEditTriggers);
    m_deviceListView->setItemDelegate(new DevItemDelegate(this));
}

void DevCollaborationWidget::loadDevice()
{
    if (!m_deviceListView->count()) {
        for (CollaborationDevice *device : m_deviceModel->devices()) {
            addItem(device);
        }
    } else {
        updateDeviceListView();
    }

    if(!m_deviceListView->count()) {
        m_deviceListView->hide();
    } else {
        if (!m_deviceListView->isVisible())
            m_deviceListView->setVisible(true);

        m_deviceListView->setFixedSize(ITEM_WIDTH, m_deviceListView->count() * ITEM_HEIGHT + LISTVIEW_ITEM_SPACE * (m_deviceListView->count() * 2));
    }

    resetWidgetSize();
}

void DevCollaborationWidget::addItem(const CollaborationDevice *device)
{
    if (!device)
        return;

    QStandardItem *item = new QStandardItem();
    DevItemDelegate::DevItemData data;
    data.checkedIconPath = device->deviceIcon(); // TODO
    data.text = device->name();
    data.iconPath = device->deviceIcon();
    int resultState = device->isCooperated() ? DevItemDelegate::Connected : DevItemDelegate::None;

    item->setData(QVariant::fromValue(data), DevItemDelegate::StaticDataRole);
    item->setData(device->machinePath(), DevItemDelegate::MachinePathDataRole);
    item->setData(0, DevItemDelegate::DegreeDataRole);
    item->setData(resultState, DevItemDelegate::ResultDataRole);

    m_viewItemModel->appendRow(item);
    m_deviceItemMap[device->machinePath()] = item;

    connect(device, &CollaborationDevice::pairedStateChanged, this, &DevCollaborationWidget::itemStatusChanged);
}

void DevCollaborationWidget::updateDeviceListView()
{
    QList<CollaborationDevice *> devices = m_deviceModel->devices();
    if (devices.isEmpty()) {
        m_deviceListView->removeItems(0, m_deviceListView->count());
        m_deviceItemMap.clear();
        m_connectingDevices.clear();
        return;
    }

    // 删除不存在设备
    for (int row = 0; row < m_deviceListView->count(); row++) {
        QStandardItem *item = m_viewItemModel->item(row);
        if (!item)
            continue;

        QString machinePath = item->data(DevItemDelegate::MachinePathDataRole).toString();
        if (m_deviceModel->getDevice(machinePath))
            continue;

        m_deviceListView->removeItem(row);

        if (m_deviceItemMap.contains(machinePath)) {
            m_deviceItemMap.remove(machinePath);
        }

        if (m_connectingDevices.contains(machinePath)) {
            m_connectingDevices.removeAll(machinePath);
        }
    }

    // 处理新增
    for (CollaborationDevice *device : devices) {
        if (!m_deviceItemMap.contains(device->machinePath())) {
            addItem(device);
        }
    }
}

void DevCollaborationWidget::resetWidgetSize()
{
    int height = TITLE_HEIGHT + (m_deviceListView->count() ? m_deviceListView->height() : 0);

    setFixedSize(ITEM_WIDTH, height);
}

void DevCollaborationWidget::itemClicked(const QModelIndex &index)
{
    QString machinePath = index.data(DevItemDelegate::MachinePathDataRole).toString();
    CollaborationDevice *device = m_deviceModel->getDevice(machinePath);
    if (!device)
        return;

    if (!device->isConnected()) {
        device->setDeviceIsCooperating(true);
        device->connect();
        if (!m_connectingDevices.contains(machinePath))
            m_connectingDevices.append(machinePath);
    } else if (!device->isCooperated()) {
        device->requestCooperate();
        if (!m_connectingDevices.contains(machinePath))
            m_connectingDevices.append(machinePath);
    } else if (device->isCooperated()) {
        device->disconnectDevice();
        if (m_connectingDevices.contains(machinePath))
            m_connectingDevices.removeOne(machinePath);
    }

    if (!m_connectingDevices.isEmpty() && !m_refreshTimer->isActive())
        m_refreshTimer->start(80);
}

void DevCollaborationWidget::itemStatusChanged()
{
    CollaborationDevice *device = qobject_cast<CollaborationDevice *>(sender());
    if (!device)
        return;

    device->setDeviceIsCooperating(false);
    QString machinePath = device->machinePath();
    if (m_deviceItemMap.contains(machinePath) && m_deviceItemMap[machinePath]) {
        // 更新item的连接状态
        int resultState = device->isCooperated() ? DevItemDelegate::Connected : DevItemDelegate::None;
        m_deviceItemMap[machinePath]->setData(resultState, DevItemDelegate::ResultDataRole);
        if (device->isCooperated() || !device->isConnected())
            m_deviceItemMap[machinePath]->setData(0, DevItemDelegate::DegreeDataRole);

        m_deviceListView->update(m_deviceItemMap[machinePath]->index());

        if ((resultState == DevItemDelegate::Connected || !device->isConnected()) && m_connectingDevices.contains(machinePath)) {
            m_connectingDevices.removeAll(machinePath);
        }
    }
}

void DevCollaborationWidget::refreshViewItem()
{
    if (m_connectingDevices.isEmpty()) {
        m_refreshTimer->stop();
        return;
    }

    for (const QString &machinePath : m_connectingDevices) {
        if (m_deviceItemMap.contains(machinePath) && m_deviceItemMap[machinePath]) {
            int degree = m_deviceItemMap[machinePath]->data(DevItemDelegate::DegreeDataRole).toInt();
            degree += PER_DEGREE; // 递进值
            m_deviceItemMap[machinePath]->setData(DevItemDelegate::Connecting, DevItemDelegate::ResultDataRole);
            m_deviceItemMap[machinePath]->setData(degree, DevItemDelegate::DegreeDataRole);
            m_deviceListView->update(m_deviceItemMap[machinePath]->index());
        }
    }
}
