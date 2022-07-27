/*
 * Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
 *
 * Author:     zhaoyingzhen <zhaoyingzhen@uniontech.com>
 *
 * Maintainer: zhaoyingzhen <zhaoyingzhen@uniontech.com>
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
#include "collaborationdevmodel.h"

#include <QIcon>
#include <QTimer>
#include <QDebug>
#include <QDBusArgument>
#include <QDBusInterface>
#include <QDBusPendingCall>

#include <DGuiApplicationHelper>

DGUI_USE_NAMESPACE
DCORE_USE_NAMESPACE

static const QString CollaborationService = "com.deepin.Cooperation";
static const QString CollaborationPath = "/com/deepin/Cooperation";
static const QString CollaborationInterface = "com.deepin.Cooperation";
static const QString ColPropertiesInterface = "org.freedesktop.DBus.Properties";

CollaborationDevModel::CollaborationDevModel(QObject *parent)
    : QObject(parent)
    , m_timer(new QTimer(this))
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

    connect(m_timer, &QTimer::timeout, this, &CollaborationDevModel::callScanMethod);
}

void CollaborationDevModel::scanDevice()
{
    callScanMethod();
    m_timer->start(30 * 1000); // 30s
}

void CollaborationDevModel::stopScanDevice()
{
    if (m_timer->isActive())
        m_timer->stop();
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
        QStringList devPaths = changedProps.value("Machines").toStringList();
        updateDevice(devPaths);
    }
}

void CollaborationDevModel::callScanMethod()
{
    // TODO 该功能目前不可用
    // m_dbusInter->asyncCall("Scan");
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

const CollaborationDevice *CollaborationDevModel::getDevice(const QString &uuid)
{
    QList<CollaborationDevice *> devices = m_devices.values();
    for (const CollaborationDevice *device : devices) {
        if (device->uuid() == uuid) {
            return device;
        }
    }

    return nullptr;
}

CollaborationDevice::CollaborationDevice(const QString &devPath, QObject *parent)
    : QObject(parent)
    , m_path(devPath)
    , m_OS(-1)
    , m_isPaired(false)
    , m_isCooperated(false)
    , m_isValid(false)
    , m_devDbusInter(new QDBusInterface(CollaborationService, devPath, CollaborationInterface + QString(".Machine"),
                                        QDBusConnection::sessionBus(), this))
{
    if (m_devDbusInter->isValid()) {
        m_name = m_devDbusInter->property("Name").toString();
        m_OS = m_devDbusInter->property("OS").toInt();
        m_isPaired = m_devDbusInter->property("Paired").toBool();
        m_isCooperated = m_devDbusInter->property("Cooperating").toBool();
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
    return m_isValid;
}

QString CollaborationDevice::name() const
{
    return m_name;
}

QString CollaborationDevice::uuid() const
{
    return m_uuid;
}

QString CollaborationDevice::deviceIcon() const
{
    switch (m_OS) {
        case DeviceType::Android: {
            if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                return QString(":/icons/resources/ICON_Device_Headphone_dark.svg");

            return QString(":/icons/resources/ICON_Device_Headphone.svg");
        }
        default: {
            if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                return QString(":/icons/resources/ICON_Device_Laptop_dark.svg");

            return QString(":/icons/resources/ICON_Device_Laptop.svg");
        }
    }
}

bool CollaborationDevice::isPaired() const
{
    return m_isPaired;
}

bool CollaborationDevice::isCooperated() const
{
    return m_isCooperated;
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
    if (changedProps.contains("Paired")) {
        bool isPaired = changedProps.value("Paired").value<bool>();
        m_isPaired = isPaired;
        if (isPaired) {
            // paired 成功之后再去请求cooperate
            requestCooperate();
        } else {
            Q_EMIT pairedStateChanged(false);
        }
    } else if (changedProps.contains("Cooperating")) {
        m_isCooperated = changedProps.value("Cooperating").value<bool>();

        Q_EMIT pairedStateChanged(m_isCooperated);
    }
}

void CollaborationDevice::requestCooperate() const
{
    callMethod("RequestCooperate");
}

void CollaborationDevice::disconnectDevice() const
{
    callMethod("Disconnect");
}

void CollaborationDevice::pair() const
{
    callMethod("Pair");
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
