/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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

#include "entryinterface.h"

/*
 * Implementation of interface class __Entry
 */

#ifdef USE_AM

void registerWindowListMetaType()
{
    qRegisterMetaType<WindowList>();
    qDBusRegisterMetaType<WindowList>();
}

void registerWindowInfoMapMetaType()
{
    registerWindowInfoMetaType();

    qRegisterMetaType<WindowInfoMap>("WindowInfoMap");
    qDBusRegisterMetaType<WindowInfoMap>();
}

void registerWindowInfoMetaType()
{
    qRegisterMetaType<WindowInfo>("WindowInfo");
    qDBusRegisterMetaType<WindowInfo>();
}

QDebug operator<<(QDebug argument, const WindowInfo &info)
{
    argument << '(' << info.title << ',' << info.attention << info.uuid << ')';

    return argument;
}

QDBusArgument &operator<<(QDBusArgument &argument, const WindowInfo &info)
{
    argument.beginStructure();
    argument << info.title << info.attention << info.uuid;
    argument.endStructure();

    return argument;
}

const QDBusArgument &operator>>(const QDBusArgument &argument, WindowInfo &info)
{
    argument.beginStructure();
    argument >> info.title >> info.attention >> info.uuid;
    argument.endStructure();

    return argument;
}

bool WindowInfo::operator==(const WindowInfo &rhs) const
{
    return (attention == rhs.attention &&
           title == rhs.title &&
           uuid == rhs.uuid);
}

class EntryPrivate
{
public:
    EntryPrivate()
        : CurrentWindow(0)
        , IsActive(false)
        , IsDocked(false)
        , mode(0)
    {}

    // begin member variables
    uint CurrentWindow;
    QString DesktopFile;
    QString Icon;
    QString Id;
    bool IsActive;
    bool IsDocked;
    QString Menu;
    QString Name;

    WindowInfoMap WindowInfos;
    int mode;

public:
    QMap<QString, QDBusPendingCallWatcher *> m_processingCalls;
    QMap<QString, QList<QVariant>> m_waittingCalls;
};

Dock_Entry::Dock_Entry(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent)
    : QDBusAbstractInterface(service, path, staticInterfaceName(), connection, parent)
    , d_ptr(new EntryPrivate)
{
    QDBusConnection::sessionBus().connect(this->service(), this->path(),
                                          "org.freedesktop.DBus.Properties",
                                          "PropertiesChanged","sa{sv}as",
                                          this,
                                          SLOT(onPropertyChanged(const QDBusMessage &)));

    if (QMetaType::type("WindowList") == QMetaType::UnknownType)
        registerWindowListMetaType();
    if (QMetaType::type("WindowInfoMap") == QMetaType::UnknownType)
        registerWindowInfoMapMetaType();
}

Dock_Entry::~Dock_Entry()
{
    qDeleteAll(d_ptr->m_processingCalls.values());
    delete d_ptr;
}

void Dock_Entry::onPropertyChanged(const QString &propName, const QVariant &value)
{
    if (propName == QStringLiteral("CurrentWindow")) {
        const uint &CurrentWindow = qvariant_cast<uint>(value);
        if (d_ptr->CurrentWindow != CurrentWindow) {
            d_ptr->CurrentWindow = CurrentWindow;
            Q_EMIT CurrentWindowChanged(d_ptr->CurrentWindow);
        }
        return;
    }

    if (propName == QStringLiteral("DesktopFile")) {
        const QString &DesktopFile = qvariant_cast<QString>(value);
        if (d_ptr->DesktopFile != DesktopFile) {
            d_ptr->DesktopFile = DesktopFile;
            Q_EMIT DesktopFileChanged(d_ptr->DesktopFile);
        }
        return;
    }

    if (propName == QStringLiteral("Icon")) {
        const QString &Icon = qvariant_cast<QString>(value);
        if (d_ptr->Icon != Icon)
        {
            d_ptr->Icon = Icon;
            Q_EMIT IconChanged(d_ptr->Icon);
        }
        return;
    }

    if (propName == QStringLiteral("IsActive")) {
        const bool &IsActive = qvariant_cast<bool>(value);
        if (d_ptr->IsActive != IsActive) {
            d_ptr->IsActive = IsActive;
            Q_EMIT IsActiveChanged(d_ptr->IsActive);
        }
        return;
    }

    if (propName == QStringLiteral("IsDocked")) {
        const bool &IsDocked = qvariant_cast<bool>(value);
        if (d_ptr->IsDocked != IsDocked) {
            d_ptr->IsDocked = IsDocked;
            Q_EMIT IsDockedChanged(d_ptr->IsDocked);
        }
        return;
    }

    if (propName == QStringLiteral("Menu")) {
        const QString &Menu = qvariant_cast<QString>(value);
        if (d_ptr->Menu != Menu) {
            d_ptr->Menu = Menu;
            Q_EMIT MenuChanged(d_ptr->Menu);
        }
        return;
    }

    if (propName == QStringLiteral("Name")) {
        const QString &Name = qvariant_cast<QString>(value);
        if (d_ptr->Name != Name) {
            d_ptr->Name = Name;
            Q_EMIT NameChanged(d_ptr->Name);
        }
        return;
    }

    if (propName == QStringLiteral("WindowInfos")) {
        const WindowInfoMap &WindowInfos = qvariant_cast<WindowInfoMap>(value);
        if (d_ptr->WindowInfos != WindowInfos) {
            d_ptr->WindowInfos = WindowInfos;
            Q_EMIT WindowInfosChanged(d_ptr->WindowInfos);
        }
        return;
    }

    if (propName == QStringLiteral("Mode")) {
        const int mode = qvariant_cast<int>(value);
        if (d_ptr->mode != mode) {
            d_ptr->mode = mode;
            Q_EMIT ModeChanged(d_ptr->mode);
        }
    }

    qWarning() << "property not handle: " << propName;
    return;
}

uint Dock_Entry::currentWindow()
{
    return qvariant_cast<uint>(property("CurrentWindow"));
}

QString Dock_Entry::desktopFile()
{
    return qvariant_cast<QString>(property("DesktopFile"));
}

QString Dock_Entry::icon()
{
    return qvariant_cast<QString>(property("Icon"));
}

QString Dock_Entry::id()
{
    return qvariant_cast<QString>(property("Id"));
}

bool Dock_Entry::isActive()
{
    return qvariant_cast<bool>(property("IsActive"));
}

bool Dock_Entry::isDocked()
{
    return qvariant_cast<bool>(property("IsDocked"));
}

int Dock_Entry::mode() const
{
    return qvariant_cast<int>(property("Mode"));
}

QString Dock_Entry::menu()
{
    return qvariant_cast<QString>(property("Menu"));
}

QString Dock_Entry::name()
{
    return qvariant_cast<QString>(property("Name"));
}

WindowInfoMap Dock_Entry::windowInfos()
{
    return qvariant_cast<WindowInfoMap>(property("WindowInfos"));
}

void Dock_Entry::CallQueued(const QString &callName, const QList<QVariant> &args)
{
    if (d_ptr->m_waittingCalls.contains(callName)) {
        d_ptr->m_waittingCalls[callName] = args;
        return;
    }
    if (d_ptr->m_processingCalls.contains(callName)) {
        d_ptr->m_waittingCalls.insert(callName, args);
    } else {
        QDBusPendingCallWatcher *watcher = new QDBusPendingCallWatcher(asyncCallWithArgumentList(callName, args));
        connect(watcher, &QDBusPendingCallWatcher::finished, this, &Dock_Entry::onPendingCallFinished);
        d_ptr->m_processingCalls.insert(callName, watcher);
    }
}

void Dock_Entry::onPendingCallFinished(QDBusPendingCallWatcher *w)
{
    w->deleteLater();
    const auto callName = d_ptr->m_processingCalls.key(w);
    Q_ASSERT(!callName.isEmpty());
    if (callName.isEmpty())
        return;

    d_ptr->m_processingCalls.remove(callName);
    if (!d_ptr->m_waittingCalls.contains(callName))
        return;

    const auto args = d_ptr->m_waittingCalls.take(callName);
    CallQueued(callName, args);
}

#endif
