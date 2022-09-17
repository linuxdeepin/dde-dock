// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "com_deepin_dde_dock.h"

/*
 * Implementation of interface class _Dock
 */

class _DockPrivate
{
public:
   _DockPrivate() = default;

    // begin member variables
    bool showInPrimary;

public:
    QMap<QString, QDBusPendingCallWatcher *> m_processingCalls;
    QMap<QString, QList<QVariant>> m_waittingCalls;
};

_Dock::_Dock(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent)
    : DBusExtendedAbstractInterface(service, path, staticInterfaceName(), connection, parent)
    , d_ptr(new _DockPrivate)
{
    connect(this, &_Dock::propertyChanged, this, &_Dock::onPropertyChanged);

}

_Dock::~_Dock()
{
    qDeleteAll(d_ptr->m_processingCalls.values());
    delete d_ptr;
}

void _Dock::onPropertyChanged(const QString &propName, const QVariant &value)
{
    if (propName == QStringLiteral("showInPrimary"))
    {
        const bool &showInPrimary = qvariant_cast<bool>(value);
        if (d_ptr->showInPrimary != showInPrimary)
        {
            d_ptr->showInPrimary = showInPrimary;
            Q_EMIT ShowInPrimaryChanged(d_ptr->showInPrimary);
        }
        return;
    }

    qWarning() << "property not handle: " << propName;
    return;
}

bool _Dock::showInPrimary()
{
    return qvariant_cast<bool>(internalPropGet("showInPrimary", &d_ptr->showInPrimary));
}

void _Dock::setShowInPrimary(bool value)
{

   internalPropSet("showInPrimary", QVariant::fromValue(value), &d_ptr->showInPrimary);
}

void _Dock::CallQueued(const QString &callName, const QList<QVariant> &args)
{
    if (d_ptr->m_waittingCalls.contains(callName))
    {
        d_ptr->m_waittingCalls[callName] = args;
        return;
    }
    if (d_ptr->m_processingCalls.contains(callName))
    {
        d_ptr->m_waittingCalls.insert(callName, args);
    } else {
        QDBusPendingCallWatcher *watcher = new QDBusPendingCallWatcher(asyncCallWithArgumentList(callName, args));
        connect(watcher, &QDBusPendingCallWatcher::finished, this, &_Dock::onPendingCallFinished);
        d_ptr->m_processingCalls.insert(callName, watcher);
    }
}

void _Dock::onPendingCallFinished(QDBusPendingCallWatcher *w)
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
