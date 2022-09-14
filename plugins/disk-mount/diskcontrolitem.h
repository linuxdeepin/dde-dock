// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DISKCONTROLITEM_H
#define DISKCONTROLITEM_H

#include "dbus/dbusdiskmount.h"

#include <dimagebutton.h>

#include <QWidget>
#include <QLabel>
#include <QProgressBar>
#include <QIcon>

class DiskControlItem : public QFrame
{
    Q_OBJECT

public:
    explicit DiskControlItem(const DiskInfo &info, QWidget *parent = 0);

signals:
    void requestUnmount(const QString &diskId) const;

private slots:
    void updateInfo(const DiskInfo &info);
    const QString formatDiskSize(const quint64 size) const;

private:
    DiskInfo m_info;
    QIcon m_unknowIcon;

    QLabel *m_diskIcon;
    QLabel *m_diskName;
    QLabel *m_diskCapacity;
    QProgressBar *m_capacityValueBar;
    Dtk::Widget::DImageButton *m_unmountButton;
};

#endif // DISKCONTROLITEM_H
