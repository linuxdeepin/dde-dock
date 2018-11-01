/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             listenerri <listenerri@gmail.com>
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

#ifndef WIRELESSTRAYWIDGET_H
#define WIRELESSTRAYWIDGET_H

#include "constants.h"

#include "abstractnetworktraywidget.h"
#include "applet/wirelesslist.h"

#include <QHash>
#include <QLabel>

#include <WirelessDevice>

class TipsWidget;
class WirelessTrayWidget : public AbstractNetworkTrayWidget
{
    Q_OBJECT

public:
    explicit WirelessTrayWidget(dde::network::WirelessDevice *device, QWidget *parent = nullptr);
    ~WirelessTrayWidget();

public:
    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

    QWidget *trayTipsWidget() Q_DECL_OVERRIDE;
    QWidget *trayPopupApplet() Q_DECL_OVERRIDE;

public Q_SLOTS:
    // set the device name displayed
    // in the top-left corner of the applet
    void setDeviceInfo(const int index);

Q_SIGNALS:
    void requestActiveAP(const QString &devPath, const QString &apPath, const QString &uuid) const;
    void requestDeactiveAP(const QString &devPath) const;
    void feedSecret(const QString &connectionPath, const QString &settingName, const QString &password, const bool autoConnect);
    void cancelSecret(const QString &connectionPath, const QString &settingName);
    void queryActiveConnInfo();
    void requestWirelessScan();
    void createApConfig(const QString &devPath, const QString &apPath);
    void queryConnectionSession( const QString &devPath, const QString &uuid );

protected:
    bool eventFilter(QObject *o, QEvent *e) Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;
    void resizeEvent(QResizeEvent *e) Q_DECL_OVERRIDE;

private:
    const QPixmap cachedPix(const QString &key, const int size);

private Q_SLOTS:
    void init();
    void adjustHeight();

private:
    QTimer *m_refershTimer;
    QWidget *m_wirelessApplet;
    TipsWidget *m_wirelessPopup;
    WirelessList *m_APList;

    QHash<QString, QPixmap> m_icons;
    QPixmap m_pixmap;

    bool m_reloadIcon;
};

#endif // WIRELESSTRAYWIDGET_H
