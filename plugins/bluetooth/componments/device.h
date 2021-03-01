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

#ifndef DEVICE_H
#define DEVICE_H

#include <QObject>
#include <QDebug>

class Device : public QObject
{
    Q_OBJECT

public:
    enum State {
        StateUnavailable = 0,
        StateAvailable   = 1,
        StateConnected   = 2
    };
    Q_ENUM(State)

private:
    static QMap<QString,QString> deviceType2Icon;

public:
    explicit Device(QObject *parent = nullptr);
    ~Device();

    inline QString id() const { return m_id; }
    void setId(const QString &id);

    inline QString name() const { return m_name; }
    void setName(const QString &name);

    inline QString alias() const { return m_alias; }
    void setAlias(const QString &alias);

    inline bool paired() const { return m_paired; }
    void setPaired(bool paired);

    inline State state() const { return m_state; }
    void setState(const State &state);

    inline bool connectState() const { return m_connectState; }
    void setConnectState(const bool connectState);

    inline bool trusted() const { return m_trusted; }
    void setTrusted(bool trusted);

    inline bool connecting() const { return m_connecting; }
    void setConnecting(bool connecting);

    inline int rssi() const { return  m_rssi; }
    void setRssi(int rssi);

    inline void setAdapterId(const QString &id) { m_adapterId = id; }
    inline const QString &getAdapterId() const { return m_adapterId; }

    inline QString deviceType() const { return m_deviceType; }
    void setDeviceType(const QString &deviceType);

    inline int deviceTime() const { return m_time; }
    void updateDeviceTime();

Q_SIGNALS:
    void nameChanged(const QString &name) const;
    void aliasChanged(const QString &alias) const;
    void pairedChanged(const bool paired) const;
    void stateChanged(const State state) const;
    void connectStateChanged(const bool connectState) const;
    void trustedChanged(const bool trusted) const;
    void connectingChanged(const bool &connecting) const;
    void rssiChanged(const int rssi) const;

private:
    QString m_id;
    QString m_name;
    QString m_alias;
    bool m_paired;
    bool m_trusted;
    bool m_connecting;
    int m_rssi;
    State m_state;
    bool m_connectState;
    QString m_adapterId;
    QString m_deviceType;
    int m_time;
};

QDebug &operator<<(QDebug &stream, const Device *device);

#endif // DEVICE_H
