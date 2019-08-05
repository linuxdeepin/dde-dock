/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
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

#ifndef DISKCONTROLWIDGET_H
#define DISKCONTROLWIDGET_H

#include "dbus/dbusdiskmount.h"

#include <QScrollArea>
#include <QVBoxLayout>

class DiskControlWidget : public QScrollArea
{
    Q_OBJECT

public:
    explicit DiskControlWidget(QWidget *parent = 0);

    void unmountAll();

signals:
    void diskCountChanged(const int count) const;

private slots:
    void diskListChanged();
    void unmountDisk(const QString &diskId) const;
    void unmountFinished(const QString &uuid, const QString &info);

private:
    QVBoxLayout *m_centralLayout;
    QWidget *m_centralWidget;
    DBusDiskMount *m_diskInter;

    DiskInfoList m_diskInfoList;
};

#endif // DISKCONTROLWIDGET_H
