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

#ifndef ADAPTER_H
#define ADAPTER_H

#include "device.h"

class Adapter : public QObject
{
    Q_OBJECT
public:
    explicit Adapter(QObject *parent = nullptr);

    inline QString name() const { return m_name; }
    void setName(const QString &name);

    inline QString id() const { return m_id; }
    void setId(const QString &id);

    QMap<QString, const Device *> devices() const;
    const Device *deviceById(const QString &id) const;

    inline bool powered() const { return m_powered; }
    void setPowered(bool powered);

    inline bool isCurrent() { return m_current; }
    inline void setCurrent(bool c) { m_current = c; }

    void initDevicesList(const QJsonDocument &doc);
    void addDevice(const QJsonObject &deviceObj);
    void removeDevice(const QString &deviceId);
    void updateDevice(const QJsonObject &json);
    void removeAllDevices();
    const QMap<QString, const Device *> &paredDevices() const;
    int paredDevicesCount() const;

Q_SIGNALS:
    void nameChanged(const QString &name) const;
    void deviceAdded(const Device *device) const;
    void deviceRemoved(const Device *device) const;
    void poweredChanged(const bool powered) const;

private:
    void divideDevice(const Device *device);

private:
    QString m_id;
    QString m_name;
    bool m_powered;
    bool m_current;

    QMap<QString, const Device *> m_devices;
    QMap<QString, const Device *> m_paredDev;
};

#endif // ADAPTER_H
