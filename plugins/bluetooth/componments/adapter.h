// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ADAPTER_H
#define ADAPTER_H

#include <QObject>
#include <QMap>

class Device;
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

    inline bool discover() {return  m_discover;}
    void setDiscover(bool discover);

    void initDevicesList(const QJsonDocument &doc);
    void addDevice(const QJsonObject &deviceObj);
    void removeDevice(const QString &deviceId);
    void updateDevice(const QJsonObject &dviceJson);

Q_SIGNALS:
    void nameChanged(const QString &name) const;
    void deviceAdded(const Device *device) const;
    void deviceRemoved(const Device *device) const;
    void deviceNameUpdated(const Device *device) const;
    void poweredChanged(const bool powered) const;
    void discoveringChanged(const bool discover) const;

private:
    QString m_id;
    QString m_name;
    bool m_powered;
    bool m_current;
    bool m_discover;

    QMap<QString, const Device *> m_devices;
};

#endif // ADAPTER_H
