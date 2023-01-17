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

#ifndef BLUETOOTHITEM_H
#define BLUETOOTHITEM_H

#include "componments/device.h"

#include <QWidget>

#define BLUETOOTH_KEY "bluetooth-item-key"

class BluetoothApplet;
class AdaptersManager;

namespace Dock {
class TipsWidget;
}
class BluetoothItem : public QWidget
{
    Q_OBJECT

public:
    explicit BluetoothItem(AdaptersManager *adapterManager, QWidget *parent = nullptr);

    QWidget *tipsWidget();
    QWidget *popupApplet();

    const QString contextMenu() const;
    void invokeMenuItem(const QString menuId, const bool checked);

    void refreshIcon();
    void refreshTips();

    bool hasAdapter();
    bool isPowered();

protected:
    void resizeEvent(QResizeEvent *event);
    void paintEvent(QPaintEvent *event);

signals:
    void requestContextMenu() const;
    void noAdapter();
    void justHasAdapter();
    void requestHide();

private:
    Dock::TipsWidget *m_tipsLabel;
    BluetoothApplet *m_applet;

    QPixmap m_iconPixmap;
    Device::State m_devState;
    bool m_adapterPowered;
};

#e