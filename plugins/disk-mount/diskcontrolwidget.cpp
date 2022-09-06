// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "diskcontrolwidget.h"
#include "diskcontrolitem.h"

#define WIDTH           300

DiskControlWidget::DiskControlWidget(QWidget *parent)
    : QScrollArea(parent),

      m_centralLayout(new QVBoxLayout),
      m_centralWidget(new QWidget),

      m_diskInter(new DBusDiskMount(this))
{
    m_centralWidget->setLayout(m_centralLayout);
    m_centralWidget->setFixedWidth(WIDTH);

    setWidget(m_centralWidget);
    setFixedWidth(WIDTH);
    setFrameStyle(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setStyleSheet("background-color:transparent;");

    connect(m_diskInter, &DBusDiskMount::DiskListChanged, this, &DiskControlWidget::diskListChanged);
    connect(m_diskInter, &DBusDiskMount::Error, this, &DiskControlWidget::unmountFinished);

    QMetaObject::invokeMethod(this, "diskListChanged", Qt::QueuedConnection);
}

void DiskControlWidget::unmountAll()
{
    for (auto disk : m_diskInfoList)
        unmountDisk(disk.m_id);
}

void DiskControlWidget::diskListChanged()
{
    DiskInfoList diskList = m_diskInter->diskList();

    while (QLayoutItem *item = m_centralLayout->takeAt(0))
    {
        delete item->widget();
        delete item;
    }

    int mountedCount = 0;
    for (auto info : diskList)
    {
        if (info.m_mountPoint.isEmpty())
            continue;
        else
            ++mountedCount;

        DiskControlItem *item = new DiskControlItem(info, this);

        connect(item, &DiskControlItem::requestUnmount, this, &DiskControlWidget::unmountDisk);

        m_centralLayout->addWidget(item);
        m_diskInfoList.append(info);
    }

    emit diskCountChanged(mountedCount);

    const int contentHeight = mountedCount * 70;
    const int maxHeight = std::min(contentHeight, 70 * 6);

    m_centralWidget->setFixedHeight(contentHeight);
    setFixedHeight(maxHeight);
}

void DiskControlWidget::unmountDisk(const QString &diskId) const
{
    m_diskInter->Unmount(diskId);
}

void DiskControlWidget::unmountFinished(const QString &uuid, const QString &info)
{
    qDebug() << uuid << info;
}
