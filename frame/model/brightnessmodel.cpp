/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
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

#include "brightnessmodel.h"

#include <QDBusArgument>
#include <QDBusInterface>
#include <QDBusPendingCall>
#include <QDebug>
#include <QApplication>
#include <QScreen>

static const QString serviceName("com.deepin.daemon.Display");
static const QString servicePath("/com/deepin/daemon/Display");
static const QString serviceInterface("com.deepin.daemon.Display");
static const QString propertiesInterface("org.freedesktop.DBus.Properties");

BrightnessModel::BrightnessModel(QObject *parent)
    : QObject(parent)
{
    QDBusInterface dbusInter(serviceName, servicePath, serviceInterface, QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        // 读取所有的屏幕的信息
        QString primaryScreenName = dbusInter.property("Primary").value<QString>();
        QList<QDBusObjectPath> paths = dbusInter.property("Monitors").value<QList<QDBusObjectPath>>();
        for (QDBusObjectPath path : paths) {
            BrightMonitor *monitor = new BrightMonitor(path.path(), this);
            monitor->setPrimary(primaryScreenName == monitor->name());
            m_monitor << monitor;
        }
    }

    connect(qApp, &QApplication::primaryScreenChanged, this, &BrightnessModel::primaryScreenChanged);
}

BrightnessModel::~BrightnessModel()
{
}

QList<BrightMonitor *> BrightnessModel::monitors()
{
    return m_monitor;
}

BrightMonitor *BrightnessModel::primaryMonitor() const
{
    for (BrightMonitor *monitor : m_monitor) {
        if (monitor->isPrimary())
            return monitor;
    }

    return nullptr;
}

void BrightnessModel::primaryScreenChanged(QScreen *screen)
{
    BrightMonitor *defaultMonitor = nullptr;
    for (BrightMonitor *monitor : m_monitor) {
        monitor->setPrimary(monitor->name() == screen->name());
        if (monitor->isPrimary())
            defaultMonitor = monitor;
    }

    if (defaultMonitor)
        Q_EMIT primaryChanged(defaultMonitor);
}

/**
 * @brief monitor
 */
BrightMonitor::BrightMonitor(QString path, QObject *parent)
    : QObject(parent)
    , m_path(path)
    , m_brightness(100)
    , m_enabled(false)
    , m_isPrimary(false)
{
    QDBusInterface dbusInter(serviceName, path, serviceInterface + QString(".Monitor"), QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        // 读取所有的屏幕的信息
        m_name = dbusInter.property("Name").toString();
        m_brightness = static_cast<int>(dbusInter.property("Brightness").toDouble() * 100);
        m_enabled = dbusInter.property("Enabled").toBool();
    }

    QDBusConnection::sessionBus().connect(serviceName, path, propertiesInterface,
                     "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(const QDBusMessage &)));
}

BrightMonitor::~BrightMonitor()
{
}

void BrightMonitor::setPrimary(bool primary)
{
    m_isPrimary = primary;
}

int BrightMonitor::brightness()
{
    return m_brightness;
}

bool BrightMonitor::enabled()
{
    return m_enabled;
}

QString BrightMonitor::name()
{
    return m_name;
}

bool BrightMonitor::isPrimary()
{
    return m_isPrimary;
}

void BrightMonitor::setBrightness(int value)
{
    callMethod("SetBrightness", { m_name, static_cast<double>(value * 0.01) });
}

void BrightMonitor::onPropertyChanged(const QDBusMessage &msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QString interfaceName = msg.arguments().at(0).toString();
    if (interfaceName != QString("%1.Monitor").arg(serviceInterface))
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    if (changedProps.contains("Brightness")) {
        int brightness = static_cast<int>(changedProps.value("Brightness").value<double>() * 100);
        if (brightness != m_brightness) {
            m_brightness = brightness;
            Q_EMIT brightnessChanged(brightness);
        }
    }
    if (changedProps.contains("Name")) {
        QString name = changedProps.value("Name").value<QString>();
        if (name != m_name) {
            m_name = name;
            Q_EMIT nameChanged(name);
        }
    }
    if (changedProps.contains("Enabled")) {
        bool enabled = changedProps.value("Enabled").value<bool>();
        if (enabled != m_enabled) {
            m_enabled = enabled;
            Q_EMIT enabledChanged(enabled);
        }
    }
}

QDBusMessage BrightMonitor::callMethod(const QString &methodName, const QList<QVariant> &argument)
{
    QDBusInterface dbusInter(serviceName, servicePath, serviceInterface, QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        QDBusPendingCall reply = dbusInter.asyncCallWithArgumentList(methodName, argument);
        reply.waitForFinished();
        return reply.reply();
    }

    return QDBusMessage();
}
