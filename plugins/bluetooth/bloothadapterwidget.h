// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
