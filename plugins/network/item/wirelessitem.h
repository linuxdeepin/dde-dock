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

#ifndef WIRELESSITEM_H
#define WIRELESSITEM_H

#include "constants.h"

#include "deviceitem.h"
#include "applet/wirelessapplet.h"

#include <QHash>

class WirelessItem : public DeviceItem
{
    Q_OBJECT

public:
    explicit WirelessItem(const QString &path);
    ~WirelessItem();

    NetworkDevice::NetworkType type() const;
    NetworkDevice::NetworkState state() const;
    QWidget *itemApplet();
    QWidget *itemPopup();

protected:
    bool eventFilter(QObject *o, QEvent *e);
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);

private:
    const QPixmap iconPix(const Dock::DisplayMode displayMode, const int size);
    const QPixmap backgroundPix(const int size);
    const QPixmap cachedPix(const QString &key, const int size);

private slots:
    void init();
    void adjustHeight();
    void refreshIcon();

private:
    QHash<QString, QPixmap> m_icons;

    QTimer *m_refershTimer;
    QWidget *m_wirelessApplet;
    QLabel *m_wirelessPopup;
    WirelessList *m_APList;
};

#endif // WIRELESSITEM_H
