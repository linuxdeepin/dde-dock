// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef POWERSTATUSWIDGET_H
#define POWERSTATUSWIDGET_H

#include <QWidget>

#include <DGuiApplicationHelper>

#define POWER_KEY "power"

DGUI_USE_NAMESPACE

class DBusPower;

// from https://upower.freedesktop.org/docs/Device.html#Device:State
enum BatteryState {
    UNKNOWN = 0,        // 未知
    CHARGING = 1,       // 充电中
    DIS_CHARGING = 2,   // 放电
    NOT_CHARGED = 3,    // 未充
    FULLY_CHARGED = 4   // 充满
};

class PowerStatusWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PowerStatusWidget(QWidget *parent = 0);
    QPixmap getBatteryIcon(int themeType);

public Q_SLOTS:
    void refreshIcon();

signals:
    void requestContextMenu(const QString &itemKey) const;
    void iconChanged();

protected:
    void resizeEvent(QResizeEvent *event);
    void paintEvent(QPaintEvent *e);


private:
    DBusPower *m_powerInter;
};

#endif // POWERSTATUSWIDGET_H
