// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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

#endif // BLUETOOTHITEM_H
