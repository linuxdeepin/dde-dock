// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef COM_DEEPIN_DDE_DOCK_H
#define COM_DEEPIN_DDE_DOCK_H

#include <QtCore/QObject>
#include <QtCore/QByteArray>
#include <QtCore/QList>
#include <QtCore/QMap>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QVariant>

#include <DBusExtendedAbstractInterface>
#include <QtDBus/QtDBus>


/*
 * Proxy class for interface com.deepin.dde.Dock
 */
class _DockPrivate;
class _Dock : public DBusExtendedAbstractInterface
{
    Q_OBJECT

public:
    static inline const char *staticInterfaceName()
    { return "com.deepin.dde.Dock"; }

public:
    explicit _Dock(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent = 0);

    ~_Dock();

    Q_PROPERTY(bool showInPrimary READ showInPrimary WRITE setShowInPrimary NOTIFY ShowInPrimaryChanged)
    bool showInPrimary();
    void setShowInPrimary(bool value);

public Q_SLOTS: // METHODS
    inline QDBusPendingReply<QStringList> GetLoadedPlugins()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("GetLoadedPlugins"), argumentList);
    }



    inline QDBusPendingReply<> ReloadPlugins()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("ReloadPlugins"), argumentList);
    }

    inline void ReloadPluginsQueued()
    {
        QList<QVariant> argumentList;

        CallQueued(QStringLiteral("ReloadPlugins"), argumentList);
    }


    inline QDBusPendingReply<> callShow()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("callShow"), argumentList);
    }

    inline void callShowQueued()
    {
        QList<QVariant> argumentList;

        CallQueued(QStringLiteral("callShow"), argumentList);
    }


    inline QDBusPendingReply<QString> getPluginKey(const QString &pluginName)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(pluginName);
        return asyncCallWithArgumentList(QStringLiteral("getPluginKey"), argumentList);
    }



    inline QDBusPendingReply<bool> getPluginVisible(const QString &pluginName)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(pluginName);
        return asyncCallWithArgumentList(QStringLiteral("getPluginVisible"), argumentList);
    }



    inline QDBusPendingReply<> resizeDock(int offset, bool dragging)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(offset) << QVariant::fromValue(dragging);
        return asyncCallWithArgumentList(QStringLiteral("resizeDock"), argumentList);
    }

    inline void resizeDockQueued(int offset, bool dragging)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(offset) << QVariant::fromValue(dragging);

        CallQueued(QStringLiteral("resizeDock"), argumentList);
    }


    inline QDBusPendingReply<> setPluginVisible(const QString &pluginName, bool visible)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(pluginName) << QVariant::fromValue(visible);
        return asyncCallWithArgumentList(QStringLiteral("setPluginVisible"), argumentList);
    }

    inline void setPluginVisibleQueued(const QString &pluginName, bool visible)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(pluginName) << QVariant::fromValue(visible);

        CallQueued(QStringLiteral("setPluginVisible"), argumentList);
    }



Q_SIGNALS: // SIGNALS
    void pluginVisibleChanged(const QString &in0, bool in1);
    // begin property changed signals
    void ShowInPrimaryChanged(bool value) const;

public Q_SLOTS:
    void CallQueued(const QString &callName, const QList<QVariant> &args);

private Q_SLOTS:
    void onPendingCallFinished(QDBusPendingCallWatcher *w);
    void onPropertyChanged(const QString &propName, const QVariant &value);

private:
    _DockPrivate *d_ptr;
};

namespace com {
  namespace deepin {
    namespace dde {
      typedef ::_Dock Dock;
    }
  }
}
#endif
