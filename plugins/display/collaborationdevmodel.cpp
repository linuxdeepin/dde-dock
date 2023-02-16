// Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "collaborationdevmodel.h"

#include <QIcon>
#include <QTimer>
#include <QDebug>
#include <QDBusArgument>
#include <QDBusInterface>
#include <QDBusPendingCall>
#include <QDBusServiceWatcher>

#include <DGuiApplicationHelper>

DGUI_USE_NAMESPACE
DCORE_USE_NAMESPACE

static const QString CollaborationService = "org.deepin.dde.Cooperation1";
static const QString CollaborationPath = "/org/deepin/dde/Cooperation1";
static const QString CollaborationInterface = "org.deepin.dde.Cooperation1";
static const QString ColPropertiesInterface = "org.freedesktop.DBus.Properties";

CollaborationDevModel::CollaborationDevModel(QObject *parent)
    : QObject(parent)
    , m_colDbusInter(new QDBusInterface(CollaborationService, CollaborationPath, CollaborationInterface, QDBusConnection::sessionBus(), this))
{
    if (m_colDbusInter->isValid()) {
        QList<QDBusObjectPath> paths = m_colDbusInter->property("Machines").value<QList<QDBusObjectPath>>();
        for (const QDBusObjectPath& path : paths) {
            CollaborationDevice *device = new CollaborationDevice(path.path(), this);
            if (device->isValid())
                m_devices[path.path()] = device;
            else
                device->deleteLater();
        }
    } else {
        qWarning() << CollaborationService << " is invalid";
    }

    m_colDbusInter->connection().connect(CollaborationService, CollaborationPath, ColPropertiesInterface,
                                         "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(QDBusMessage)));

    auto *dbusWatcher = new QDBusServiceWatcher(CollaborationService, m_colDbusInter->connection(),
                                                QDBusServiceWatcher::WatchForUnregistration, this);
    connect(dbusWatcher, &QDBusServiceWatcher::serviceUnregistered, this, [this](){
        qWarning() << CollaborationService << "unregistered";
        clear();
    });
}

void CollaborationDevModel::checkServiceValid()
{
    if (!m_colDbusInter->isValid()) {
        clear();
    }
}

QList<CollaborationDevice *> CollaborationDevModel::devices() const
{
    return m_devices.values();
}

void CollaborationDevModel::onPropertyChanged(const QDBusMessage &msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QString interfaceName = msg.arguments().at(0).toString();
    if (interfaceName != CollaborationInterface)
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    if (changedProps.contains("Machines")) {
        QList<QDBusObjectPath> paths = m_colDbusInter->property("Machines").value<QList<QDBusObjectPath>>();
        QStringList devPaths;
        for (const QDBusObjectPath& path : paths) {
            devPaths << path.path();
        }
        updateDevice(devPaths);
    }
}

void CollaborationDevModel::updateDevice(const QStringList &devPaths)
{
    if (devPaths.isEmpty()) {
        qDeleteAll(m_devices);
        m_devices.clear();
    } else {
        // 清除已不存在的设备
        QMapIterator<QString, CollaborationDevice *> it(m_devices);
        while (it.hasNext()) {
            it.next();
            if (!devPaths.contains(it.key())) {
                it.value()->deleteLater();
                m_devices.remove(it.key());
            }
        }

        // 添加新增设备
        for (const QString &path : devPaths) {
            if (!m_devices.contains(path)) {
                CollaborationDevice *device = new CollaborationDevice(path, this);
                if (device->isValid())
                    m_devices[path] = device;
                else
                    device->deleteLater();
            }
        }
    }

    emit devicesChanged();
}

void CollaborationDevModel::clear()
{
    for (CollaborationDevice *device : m_devices) {
        device->deleteLater();
    }
    m_devices.clear();

    Q_EMIT devicesChanged();
}

CollaborationDevice *CollaborationDevModel::getDevice(const QString &machinePath)
{
    return m_devices.value(machinePath, nullptr);
}

CollaborationDevice::CollaborationDevice(const QString &devPath, QObject *parent)
    : QObject(parent)
    , m_path(devPath)
    , m_OS(-1)
    , m_isConnected(false)
    , m_isCooperated(false)
    , m_isValid(false)
    , m_isCooperating(false)
    , m_devDbusInter(new QDBusInterface(CollaborationService, devPath, CollaborationInterface + QString(".Machine"),
                                        QDBusConnection::sessionBus(), this))
{
    if (m_devDbusInter->isValid()) {
        m_name = m_devDbusInter->property("Name").toString();
        m_OS = m_devDbusInter->property("OS").toInt();
        m_isConnected = m_devDbusInter->property("Connected").toBool();
        m_isCooperated = m_devDbusInter->property("DeviceSharing").toBool();
        m_uuid = m_devDbusInter->property("UUID").toString();
        m_isValid = true;
    } else {
        qWarning() << "CollaborationDevice devPath:" << devPath << " is invalid and get properties failed";
    }

    m_devDbusInter->connection().connect(CollaborationService, m_path, ColPropertiesInterface, "PropertiesChanged",
                           this, SLOT(onPropertyChanged(QDBusMessage)));
}

bool CollaborationDevice::isValid() const
{
    // not show android device
    return m_isValid && m_OS != Android;
}

QString CollaborationDevice::name() const
{
    return m_name;
}

QString CollaborationDevice::uuid() const
{
    return m_uuid;
}

QString CollaborationDevice::machinePath() const
{
    return m_path;
}

QString CollaborationDevice::deviceIcon() const
{
    switch (m_OS) {
        case DeviceType::Android: {
            if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                return QString(":/ICON_Device_Headphone_dark.svg");

            return QString(":/ICON_Device_Headphone.svg");
        }
        default: {
            if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                return QString(":/ICON_Device_Laptop_dark.svg");

            return QString(":/ICON_Device_Laptop.svg");
        }
    }
}

bool CollaborationDevice::isConnected() const
{
    return m_isConnected;
}

bool CollaborationDevice::isCooperated() const
{
    return m_isCooperated;
}

void CollaborationDevice::setDeviceIsCooperating(bool isCooperating)
{
    m_isCooperating = isCooperating;
}

void CollaborationDevice::onPropertyChanged(const QDBusMessage &msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QString interfaceName = msg.arguments().at(0).toString();
    if (interfaceName != QString("%1.Machine").arg(CollaborationInterface))
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    if (changedProps.contains("Connected")) {
        bool isConnected = changedProps.value("Connected").value<bool>();
        m_isConnected = isConnected;
        if (isConnected && m_isCooperating) {
            // paired 成功之后再去请求cooperate
            requestCooperate();
        }

        if (!isConnected){
            Q_EMIT pairedStateChanged(false);
        }
    } else if (changedProps.contains("DeviceSharing")) {
        m_isCooperated = changedProps.value("DeviceSharing").value<bool>();

        Q_EMIT pairedStateChanged(m_isCooperated);
    }
}

void CollaborationDevice::requestCooperate() const
{
    callMethod("RequestDeviceSharing");
}

void CollaborationDevice::disconnectDevice() const
{
    callMethod("Disconnect");
}

void CollaborationDevice::connect() const
{
    callMethod("Connect");
}

QDBusMessage CollaborationDevice::callMethod(const QString &methodName) const
{
    if (m_devDbusInter->isValid()) {
        QDBusMessage msg = m_devDbusInter->call(methodName);
        qInfo() << "CollaborationDevice callMethod:" << methodName << " " << msg.errorMessage();
        return msg;
    }

    qWarning() << "CollaborationDevice callMethod: " << methodName << " failed";
    return QDBusMessage();
}
