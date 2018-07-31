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

#ifndef WIRELESSAPPLET_H
#define WIRELESSAPPLET_H

#include "devicecontrolwidget.h"
#include "accesspoint.h"

#include <QScrollArea>
#include <QVBoxLayout>
#include <QList>
#include <QTimer>
#include <QCheckBox>
#include <QDBusObjectPath>

#include <dpicturesequenceview.h>
#include <dinputdialog.h>
#include <WirelessDevice>

DWIDGET_USE_NAMESPACE

class AccessPointWidget;
class WirelessList : public QScrollArea
{
    Q_OBJECT

public:
    explicit WirelessList(dde::network::WirelessDevice *deviceIter, QWidget *parent = 0);
    ~WirelessList();

    QWidget *controlPanel();

public Q_SLOTS:
    void onNeedSecrets(const QString &info);
    void onNeedSecretsFinished(const QString &info0, const QString &info1);
    void setDeviceInfo(const int index);

signals:
    void requestSetDeviceEnable(const QString &path, const bool enable) const;
    void requestActiveAP(const QString &devPath, const QString &apPath, const QString &uuid) const;
    void requestDeactiveAP(const QString &devPath) const;
    void feedSecret(const QString &connectionPath, const QString &settingName, const QString &password, const bool autoConnect);
    void cancelSecret(const QString &connectionPath, const QString &settingName);
    void requestWirelessScan();

private:
    void loadAPList();

private slots:
    void APAdded(const QJsonObject &apInfo);
    void APRemoved(const QJsonObject &apInfo);
    void APPropertiesChanged(const QJsonObject &apInfo);
    void updateAPList();
    void onEnableButtonToggle(const bool enable);
    void pwdDialogAccepted();
    void pwdDialogCanceled();
    void onPwdDialogTextChanged(const QString &text);
    void onDeviceEnableChanged(const bool enable);
    void activateAP(const QString &apPath, const QString &ssid);
    void deactiveAP();
    void updateIndicatorPos();
    void onActiveConnectionChanged();


private:
    dde::network::WirelessDevice *m_device;

    AccessPoint m_activeAP;
    QList<AccessPoint> m_apList;
    QList<AccessPointWidget*> m_apwList;

    QTimer *m_updateAPTimer;
    Dtk::Widget::DInputDialog *m_pwdDialog;
    QCheckBox *m_autoConnBox;
    Dtk::Widget::DPictureSequenceView *m_indicator;
    AccessPointWidget *m_currentClickAPW;
    AccessPoint m_currentClickAP;

    QString m_lastConnPath;
    QString m_lastConnSecurity;
    QString m_lastConnSecurityType;

    QVBoxLayout *m_centralLayout;
    QWidget *m_centralWidget;
    DeviceControlWidget *m_controlPanel;
};

#endif // WIRELESSAPPLET_H
