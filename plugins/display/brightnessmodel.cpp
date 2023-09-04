// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "brightnessmodel.h"

#include <QDBusArgument>
#include <QDBusInterface>
#include <QDBusPendingCall>
#include <QDebug>
#include <QApplication>
#include <QScreen>

static const QString serviceName("org.deepin.dde.Display1");
static const QString servicePath("/org/deepin/dde/Display1");
static const QString serviceInterface("org.deepin.dde.Display1");
static const QString propertiesInterface("org.freedesktop.DBus.Properties");

BrightnessModel::BrightnessModel(QObject *parent)
    : QObject(parent)
{
    QDBusInterface dbusInter(serviceName, servicePath, serviceInterface, QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        // 读取所有的屏幕的信息
        m_primaryScreenName = dbusInter.property("Primary").value<QString>();
        m_monitor = readMonitors(dbusInter.property("Monitors").value<QList<QDBusObjectPath>>());

        QDBusConnection::sessionBus().connect(serviceName, servicePath, propertiesInterface,
                         "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(const QDBusMessage &)));
    }
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

void BrightnessModel::onPropertyChanged(const QDBusMessage &msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QString interfaceName = msg.arguments().at(0).toString();
    if (interfaceName != serviceInterface)
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    if (changedProps.contains("Brightness")) {
        Q_EMIT monitorLightChanged();
    } else if (changedProps.contains("Primary")) {
        m_primaryScreenName = changedProps.value("Primary").toString();
        BrightMonitor *defaultMonitor = nullptr;
        for (BrightMonitor *monitor : m_monitor) {
            monitor->setPrimary(monitor->name() == m_primaryScreenName);
            if (monitor->isPrimary())
                defaultMonitor = monitor;
        }

        if (defaultMonitor)
            Q_EMIT primaryChanged(defaultMonitor);
    } else if (changedProps.contains("Monitors")) {
        int oldSize = m_monitor.size();
        qDeleteAll(m_monitor);
        m_monitor = readMonitors(changedProps.value("Monitors").value<QList<QDBusObjectPath>>());
        if (oldSize == 1 && m_monitor.size() == 0) {
            Q_EMIT screenVisibleChanged(false);
        } else if (oldSize == 0 && m_monitor.size() == 1) {
            Q_EMIT screenVisibleChanged(true);
        }
    }
}

QList<BrightMonitor *> BrightnessModel::readMonitors(const QList<QDBusObjectPath> &paths)
{
    QList<BrightMonitor *> monitors;
    for (QDBusObjectPath path : paths) {
        BrightMonitor *monitor = new BrightMonitor(path.path(), this);
        monitor->setPrimary(m_primaryScreenName == monitor->name());
        monitors << monitor;
    }
    return monitors;
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
