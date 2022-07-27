/*
 * Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
 *
 * Author:     zhaoyingzhen <zhaoyingzhen@uniontech.com>
 *
 * Maintainer: zhaoyingzhen <zhaoyingzhen@uniontech.com>
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
#include "devcollaborationwidget.h"
#include "collaborationdevmodel.h"
#include "devitemdelegate.h"

#include <DStyle>

#include <QMap>
#include <QTimer>
#include <QLabel>
#include <QVBoxLayout>
#include <QStandardItemModel>

#define TITLE_HEIGHT 16
#define ITME_WIDTH 310
#define ITEM_HEIGHT 36
#define LISTVIEW_ITEM_SPACE 2
#define ITME_SPACE 10
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
    m_deviceModel->scanDevice();

    QWidget::showEvent(event);
}

void DevCollaborationWidget::hideEvent(QHideEvent *event)
{
    m_deviceModel->stopScanDevice();

    QWidget::hideEvent(event);
}

void DevCollaborationWidget::initUI()
{
    m_deviceListView->setModel(m_viewItemModel);

    QLabel *title = new QLabel(tr("Cross-end Collaboration"), this);
    title->setFixedHeight(TITLE_HEIGHT);

    QHBoxLayout *hLayout = new QHBoxLayout();
    hLayout->setContentsMargins(10, 0, 0, 0);
    hLayout->addWidget(title);

    QVBoxLayout *mainLayout = new QVBoxLayout();
    mainLayout->setMargin(0);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->setSpacing(ITME_SPACE);
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

        m_deviceListView->setFixedSize(ITME_WIDTH, m_deviceListView->count() * ITEM_HEIGHT + LISTVIEW_ITEM_SPACE * (m_deviceListView->count() * 2));
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
    item->setData(device->uuid(), DevItemDelegate::UUIDDataRole);
    item->setData(0, DevItemDelegate::DegreeDataRole);
    item->setData(resultState, DevItemDelegate::ResultDataRole);

    m_viewItemModel->appendRow(item);
    m_uuidItemMap[device->uuid()] = item;

    connect(device, &CollaborationDevice::pairedStateChanged, this, &DevCollaborationWidget::itemStatusChanged);
}

void DevCollaborationWidget::updateDeviceListView()
{
    QList<CollaborationDevice *> devices = m_deviceModel->devices();
    if (devices.isEmpty()) {
        m_deviceListView->removeItems(0, m_deviceListView->count());
        return;
    }

    // 删除不存在设备
    for (int row = 0; row < m_deviceListView->count(); row++) {
        QStandardItem *item = m_viewItemModel->item(row);
        if (!item)
            continue;

        QString uuid = item->data(DevItemDelegate::UUIDDataRole).toString();
        if (m_deviceModel->getDevice(uuid))
            continue;

        m_deviceListView->removeItem(row);

        if (m_uuidItemMap.contains(uuid)) {
            m_uuidItemMap.remove(uuid);
        }

        if (m_connectingDevices.contains(uuid)) {
            m_connectingDevices.removeAll(uuid);
        }
    }

    // 处理新增
    for (CollaborationDevice *device : devices) {
        if (!m_uuidItemMap.contains(device->uuid())) {
            addItem(device);
        }
    }
}

void DevCollaborationWidget::resetWidgetSize()
{
    int height = TITLE_HEIGHT + ITME_SPACE + (m_deviceListView->count() ? m_deviceListView->height() : 0);

    setFixedSize(ITME_WIDTH, height);
}

void DevCollaborationWidget::itemClicked(const QModelIndex &index)
{
    QString uuid = index.data(DevItemDelegate::UUIDDataRole).toString();
    const CollaborationDevice *device = m_deviceModel->getDevice(uuid);
    if (!device)
        return;

    if (!device->isPaired()) {
        device->pair();
        m_connectingDevices.append(uuid);
    } else if (!device->isCooperated()) {
        device->requestCooperate();
        m_connectingDevices.append(uuid);
    } else if (device->isCooperated()) {
        device->disconnectDevice();
    }

    if (!m_connectingDevices.isEmpty() && !m_refreshTimer->isActive())
        m_refreshTimer->start(30);
}

void DevCollaborationWidget::itemStatusChanged()
{
    CollaborationDevice *device = qobject_cast<CollaborationDevice *>(sender());
    if (!device)
        return;

    QString uuid = device->uuid();
    if (m_uuidItemMap.contains(uuid) && m_uuidItemMap[uuid]) {
        // 更新item的连接状态
        int resultState = device->isCooperated() ? DevItemDelegate::Connected : DevItemDelegate::None;
        m_uuidItemMap[uuid]->setData(resultState, DevItemDelegate::ResultDataRole);
        if (device->isCooperated())
            m_uuidItemMap[uuid]->setData(0, DevItemDelegate::DegreeDataRole);

        m_deviceListView->update(m_uuidItemMap[uuid]->index());

        if (resultState == DevItemDelegate::Connected && m_connectingDevices.contains(uuid)) {
            m_connectingDevices.removeAll(uuid);
        }
    }
}

void DevCollaborationWidget::refreshViewItem()
{
    if (m_connectingDevices.isEmpty()) {
        m_refreshTimer->stop();
        return;
    }

    for (const QString &uuid : m_connectingDevices) {
        if (m_uuidItemMap.contains(uuid) && m_uuidItemMap[uuid]) {
            int degree = m_uuidItemMap[uuid]->data(DevItemDelegate::DegreeDataRole).toInt();
            degree += PER_DEGREE; // 递进值
            m_uuidItemMap[uuid]->setData(DevItemDelegate::Connecting, DevItemDelegate::ResultDataRole);
            m_uuidItemMap[uuid]->setData(degree, DevItemDelegate::DegreeDataRole);
            m_deviceListView->update(m_uuidItemMap[uuid]->index());
        }
    }
}
