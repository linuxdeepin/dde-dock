/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#ifndef BLOOTHADAPTERWIDGET_H
#define BLOOTHADAPTERWIDGET_H

#include <DListView>

#include <QWidget>

class Adapter;
class QLabel;
class Device;
class QStandardItemModel;

using namespace Dtk::Widget;

class BloothAdapterWidget : public QWidget
{
    Q_OBJECT

public:
    explicit BloothAdapterWidget(Adapter *adapter, QWidget *parent = nullptr);

    Adapter *adapter();

Q_SIGNALS:
    void requestConnectDevice(Device *device);
    void requestUpdate() const;

protected Q_SLOTS:
    void onDeviceAdded(const Device *device);
    void onDeviceRemoved(const Device *device);
    void onDeviceNameUpdated(const Device *device) const;
    void onPoweredChanged(const bool powered);

    void onOtherClicked(const QModelIndex &index);

private:
    void initUi();
    void initConnection();
    void initDevice();
    void adjustHeight();
    void updateDeviceVisible();

private:
    Adapter *m_adapter;
    QLabel *m_myDeviceLabel;
    DListView *m_myDeviceView;
    QLabel *m_otherDeviceLabel;
    DListView *m_otherDeviceView;
    QStandardItemModel *m_myDeviceModel;
    QStandardItemModel *m_otherDeviceModel;
};

#endif // BLOOTHADAPTERWIDGET_H
