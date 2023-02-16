// Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef COLLABORATION_DEV_MODEL_H
#define COLLABORATION_DEV_MODEL_H

#include <QMap>
#include <QObject>

class QTimer;
class QDBusInterface;
class QDBusMessage;
class CollaborationDevice;

/*!
 * \brief The CollaborationDevModel class
 * 协同设备model
 */
class CollaborationDevModel : public QObject
{
    Q_OBJECT
public:
    explicit CollaborationDevModel(QObject *parent = nullptr);

signals:
    void devicesChanged();

public:
    void checkServiceValid();

    QList<CollaborationDevice *> devices() const;
    CollaborationDevice *getDevice(const QString &machinePath);

private slots:
    void onPropertyChanged(const QDBusMessage &msg);

private:
    void updateDevice(const QStringList &devPaths);
    void clear();

private:
    QDBusInterface *m_colDbusInter;
    // machine path : device object
    QMap<QString, CollaborationDevice *> m_devices;

};

/*!
 * \brief The CollaborationDevice class
 * 协同设备类
 */
class CollaborationDevice : public QObject
{
    Q_OBJECT
public:
    explicit CollaborationDevice(const QString &devPath, QObject *parent = nullptr);

signals:
    void pairedStateChanged(bool);

public:
    bool isValid() const;
    void connect() const;
    void requestCooperate() const;
    void disconnectDevice() const;

    QString name() const;
    QString uuid() const;
    QString machinePath() const;
    QString deviceIcon() const;
    bool isConnected() const;
    bool isCooperated() const;
    void setDeviceIsCooperating(bool isCooperating);

private slots:
    void onPropertyChanged(const QDBusMessage &msg);

private:
    QDBusMessage callMethod(const QString &methodName) const;

private:
    enum DeviceType {
        Other = 0,
        UOS,
        Linux,
        Windows,
        MacOS,
        Android
    };

    QString m_path;
    QString m_name;
    QString m_uuid;
    int m_OS;

    bool m_isConnected;
    bool m_isCooperated;
    bool m_isValid;

    // 标记任务栏点击触发协同连接
    bool m_isCooperating;

    QDBusInterface *m_devDbusInter;
};

#endif // COLLABORATION_DEV_MODEL_H
