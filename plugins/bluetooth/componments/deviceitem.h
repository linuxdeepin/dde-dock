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

#ifndef DEVICEITEM_H
#define DEVICEITEM_H

#include "device.h"

#include <DSpinner>

#include <QLabel>

DWIDGET_USE_NAMESPACE

class HorizontalSeparator;
class DeviceItem : public QWidget
{
    Q_OBJECT
public:
    explicit DeviceItem(const QString &title, QWidget *parent = nullptr);
    bool operator <(const DeviceItem &item);

    inline void setTitle(const QString &name) { m_title->setText(name); }

    void setDevice(Device *d);
    inline Device *device() { return m_device; }
    inline const Device *device() const { return  m_device; }

protected:
    void mousePressEvent(QMouseEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;

signals:
    void clicked(Device *);

public slots:
    void changeState(const Device::State state);

private:
    QLabel *m_title;
    QLabel *m_state;
    DSpinner *m_loadingStat;
    Device *m_device = nullptr;
    HorizontalSeparator *m_line;
    QString m_statSuffix;
};

class HorizontalSeparator : public QWidget
{
    Q_OBJECT
public:
    explicit HorizontalSeparator(QWidget *parent = nullptr);

protected:
    void paintEvent(QPaintEvent *e);
};

class MenueItem : public QLabel
{
    Q_OBJECT
public:
    explicit MenueItem(QWidget *parent = nullptr);
signals:
    void clicked();
protected:
    void mousePressEvent(QMouseEvent *event) override;
};

#endif // DEVICEITEM_H
