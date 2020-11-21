/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
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

#ifndef ACCESSPOINTWIDGET_H
#define ACCESSPOINTWIDGET_H

#include "accesspoint.h"
#include "../wireditem.h"

#include <DImageButton>

#include <NetworkDevice>

#include <QWidget>
#include <QLabel>
#include <QDBusObjectPath>

class StateButton;
class StateLabel;
class SsidButton : public QLabel
{
    Q_OBJECT
public:
    SsidButton(QWidget *parent = nullptr) : QLabel(parent) {}
    virtual ~SsidButton() {}

signals:
    void clicked();

protected:
    void mouseReleaseEvent(QMouseEvent *event) override
    {
        QLabel::mouseReleaseEvent(event);

        Q_EMIT clicked();
    }
};

class AccessPointWidget : public QFrame
{
    Q_OBJECT
public:

    enum ApState{
        Unknown = 0,
        Activating = 1,
        Activated  = 2
    };

    explicit AccessPointWidget(const QJsonObject &apInfo);
    ~AccessPointWidget();

   //const AccessPoint ap() const { return m_ap; }
//    void updateAP(const AccessPoint &ap);

    int strength() const {return m_ap.strength();}
    const QString ssid() const {return m_ap.ssid();}
    const QString path() const {return m_ap.path();}
    const QString uuid() const {return m_ap.uuid();}
    bool secured() const {return m_ap.secured();}
    void updateApInfo(const QJsonObject &apInfo);
    void setActiveState(ApState state);

signals:
    /**
     * @def requestConnectAP
     * @brief 连接wifi
     * @param apPath
     * @param ssid
     */
    void requestConnectAP(const QString &apPath, const QString &uuid) const;
    /**
     * @def requestDisconnectAP
     * @brief 断开连接
     * @param apPath
     */
    void requestDisconnectAP(const QString &apUuid) const;
    /**
     * @def apChange
     * @brief 数据改变
     */
    void apChange();


private:
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    bool operator==(const AccessPointWidget *ap) const;
    void setStrengthIcon(const int strength);
    void initUI();
    void initConnect();

private slots:
    void ssidClicked();
    void disconnectBtnClicked();

private:
    ApState m_activeState;

    AccessPoint m_ap;
    //Ssid显示
    SsidButton *m_ssidBtn;
    QLabel *m_securityLabel;
    QLabel *m_strengthLabel;
    /**
     * @def 连接成功的图标
     * @brief m_stateButton
     */
    StateButton *m_stateButton;
    /**
     * @def m_loadingStat
     * @brief 连接中图标
     */
    DSpinner *m_loadingStat;

    QPixmap m_securityPixmap;
    QSize m_securityIconSize;
};

#endif // ACCESSPOINTWIDGET_H
