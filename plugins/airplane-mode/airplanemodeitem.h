// SPDX-FileCopyrightText: 2020 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef AIRPLANEMODEITEM_H
#define AIRPLANEMODEITEM_H

#include <com_deepin_daemon_airplanemode.h>

#include <QWidget>

using DBusAirplaneMode = com::deepin::daemon::AirplaneMode;

namespace Dock {
class TipsWidget;
}

class AirplaneModeApplet;
class AirplaneModeItem : public QWidget
{
    Q_OBJECT

public:
    explicit AirplaneModeItem(QWidget *parent = nullptr);

    QWidget *tipsWidget();
    QWidget *popupApplet();
    const QString contextMenu() const;
    void invokeMenuItem(const QString menuId, const bool checked);
    void refreshIcon();
    void updateTips();

    bool airplaneEnable();

protected:
    void resizeEvent(QResizeEvent *e);
    void paintEvent(QPaintEvent *e);

signals:
    void airplaneEnableChanged(bool enable);

private:
    Dock::TipsWidget *m_tipsLabel;
    AirplaneModeApplet *m_applet;
    DBusAirplaneMode *m_airplaneModeInter;
    QPixmap m_iconPixmap;
};

#endif // AIRPLANEMODEITEM_H
