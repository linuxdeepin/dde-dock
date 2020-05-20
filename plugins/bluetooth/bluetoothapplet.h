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

#ifndef BLUETOOTHAPPLET_H
#define BLUETOOTHAPPLET_H

#include "componments/device.h"

#include <QScrollArea>

class QLabel;
class QVBoxLayout;
class QHBoxLayout;
class Adapter;
class AdapterItem;
class HorizontalSeparator;
class MenueItem;
class AdaptersManager;
class BluetoothApplet : public QScrollArea
{
    Q_OBJECT
public:
    explicit BluetoothApplet(QWidget *parent = nullptr);
    void setAdapterPowered(bool powered);
    bool poweredInitState();
    bool hasAadapter();
    QStringList connectedDevsName();

public slots :
    void addAdapter(Adapter *adapter);
    void removeAdapter(Adapter *adapter);

signals:
    void powerChanged(bool powered);
    void deviceStateChanged(const Device::State state);
    void noAdapter();
    void justHasAdapter();

private slots:
    void onPowerChanged();
    void onDeviceStateChanged();

private:
    void updateView();
    void getDevieInitStatus(AdapterItem *item);

private:
    HorizontalSeparator *m_line;
    QLabel *m_appletName;
    QWidget *m_centralWidget;
    QVBoxLayout *m_centrealLayout;
    QVBoxLayout *m_adapterLayout;
    QHBoxLayout *m_menueLayout;
    MenueItem *m_openControlCenter;

    AdaptersManager *m_adaptersManager;

    QMap<QString, AdapterItem *> m_adapterItems;
    Device::State m_initDeviceState;
};

#endif // BLUETOOTHAPPLET_H
