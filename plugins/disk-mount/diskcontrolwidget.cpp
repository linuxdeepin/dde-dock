/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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
