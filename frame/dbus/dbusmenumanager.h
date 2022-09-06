// SPDX-FileCopyrightText: 2015 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DBUSMENUMANAGER_H_1436158928
#define DBUSMENUMANAGER_H_1436158928

#include <QtCore/QObject>
#include <QtCore/QByteArray>
#include <QtCore/QList>
#include <QtCore/QMap>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QVariant>
#include <QtDBus/QtDBus>

/*
 * Proxy class for interface com.deepin.menu.Manager
 */
class DBusMenuManager: public QDBusAbstractInterface
{
    Q_OBJECT
public:
    static inline const char *staticServerPath()
    { return "com.deepin.menu"; }
    static inline const char *staticInterfacePath()
    { return "/com/deepin/menu"; }
    static inline const char *staticInterfaceName()
    { return "com.deepin.menu.Manager"; }

public:
    explicit DBusMenuManager(QObject *parent = 0);

    ~DBusMenuManager();

public Q_SLOTS: // METHODS
    inline QDBusPendingReply<QDBusObjectPath> RegisterMenu()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("RegisterMenu"), argumentList);
    }

    inline QDBusPendingReply<> UnregisterMenu(const QString &menuObjectPath)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(menuObjectPath);
        return asyncCallWithArgumentList(QStringLiteral("UnregisterMenu"), argumentList);
    }

Q_SIGNALS: // SIGNALS
};

namespace com {
  namespace deepin {
    namespace menu {
      typedef ::DBusMenuManager Manager;
    }
  }
}
#endif
