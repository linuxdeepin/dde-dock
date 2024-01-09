// SPDX-FileCopyrightText: 2024 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later
#include "notification.h"
#include "constants.h"

#include <QPainter>
#include <QPainterPath>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>
#include <QDBusInterface>
#include <QDBusReply>
#include <QtConcurrent/QtConcurrent>

#include <DStyle>
#include <DGuiApplicationHelper>

Q_DECLARE_LOGGING_CATEGORY(qLcPluginNotification)

DWIDGET_USE_NAMESPACE;
DCORE_USE_NAMESPACE;
Notification::Notification(QWidget *parent)
    : QWidget(parent)
    , m_icon(QIcon::fromTheme("notification"))
    , m_notificationCount(0)
    , m_dbus(nullptr)
    , m_dndMode(false)
{
    setMinimumSize(PLUGIN_BACKGROUND_MIN_SIZE, PLUGIN_BACKGROUND_MIN_SIZE);
    connect(this, &Notification::dndModeChanged, this, &Notification::refreshIcon);
    QtConcurrent::run([this](){
        m_dbus.reset(new QDBusInterface("org.deepin.dde.Notification1", "/org/deepin/dde/Notification1", "org.deepin.dde.Notification1"));
        // Refresh icon for the first time, cause org.deepin.dde.Notification1 might depend on dock's DBus,
        // we should not call org.deepin.dde.Notification1 in the main thread before dock's dbus is initialized.
        // Just refresh icon in the other thread.
        QDBusReply<QDBusVariant> dnd = m_dbus->call(QLatin1String("GetSystemInfo"), QVariant::fromValue(0u));
        if (!dnd.isValid()) {
            qCWarning(qLcPluginNotification) << dnd.error();
        } else {
            m_dndMode = dnd.value().variant().toBool();
            refreshIcon();
        }
        QDBusConnection::sessionBus().connect("org.deepin.dde.Notification1",
                                              "/org/deepin/dde/Notification1",
                                              "org.deepin.dde.Notification1",
                                              "SystemInfoChanged",
                                              this,
                                              SLOT(onSystemInfoChanged(quint32,QDBusVariant))
                                              );
        auto recordCountVariant = m_dbus->property("recordCount");
        if (!recordCountVariant.isValid()) {
            qCWarning(qLcPluginNotification) << dnd.error();
        } else {
            setNotificationCount(recordCountVariant.toUInt());
        }
        QDBusConnection::sessionBus().connect("org.deepin.dde.Notification1",
                                              "/org/deepin/dde/Notification1",
                                              "org.deepin.dde.Notification1",
                                              "recordCountChanged",
                                              this,
                                              SLOT(setNotificationCount(uint))
                                              );
    });
}

QIcon Notification::icon() const
{
    return m_icon;
}

void Notification::refreshIcon()
{
    m_icon = QIcon::fromTheme(dndMode() ? "notification-off" : "notification");
    Q_EMIT iconRefreshed();
}

bool Notification::dndMode() const
{
    return m_dndMode;
}

void Notification::setDndMode(bool dnd)
{
    if (m_dbus) {
        m_dbus->call(QLatin1String("SetSystemInfo"), QVariant::fromValue(0u), QVariant::fromValue(QDBusVariant(dnd)));
    }
}

uint Notification::notificationCount() const
{
    return m_notificationCount;
}

void Notification::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e)
    QPainter p(this);
    m_icon.paint(&p, rect());
}

void Notification::onSystemInfoChanged(quint32 info, QDBusVariant value)
{
    if (info == 0) {
        // DND mode
        m_dndMode = value.variant().toBool();
        Q_EMIT dndModeChanged(m_dndMode);
    }
}

void Notification::setNotificationCount(uint count)
{
    if (m_notificationCount == count) {
        return;
    }
    m_notificationCount = count;
    Q_EMIT this->notificationCountChanged(count);
}
