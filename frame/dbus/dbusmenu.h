// SPDX-FileCopyrightText: 2015 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DBUSMENU_H_1436158836
#define DBUSMENU_H_1436158836

#include <QtCore/QObject>
#include <QtCore/QByteArray>
#include <QtCore/QList>
#include <QtCore/QMap>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QVariant>
#include <QtDBus/QtDBus>

/*
 * Proxy class for interface com.deepin.menu.Menu
 */
class DBusMenu: public QDBusAbstractInterface
{
    Q_OBJECT
public:
    static inline const char *staticServerPath()
    { return "com.deepin.menu"; }
    static inline const char *staticInterfaceName()
    { return "com.deepin.menu.Menu"; }

public:
    DBusMenu(const QString &path,QObject *parent = 0);

    ~DBusMenu();

public Q_SLOTS: // METHODS
    inline QDBusPendingReply<> SetItemActivity(const QString &itemId, bool isActive)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(itemId) << QVariant::fromValue(isActive);
        return asyncCallWithArgumentList(QStringLiteral("SetItemActivity"), argumentList);
    }

    inline QDBusPendingReply<> SetItemChecked(const QString &itemId, bool checked)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(itemId) << QVariant::fromValue(checked);
        return asyncCallWithArgumentList(QStringLiteral("SetItemChecked"), argumentList);
    }

    inline QDBusPendingReply<> SetItemText(const QString &itemId, const QString &text)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(itemId) << QVariant::fromValue(text);
        return asyncCallWithArgumentList(QStringLiteral("SetItemText"), argumentList);
    }

    inline QDBusPendingReply<> ShowMenu(const QString &menuJsonContent)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(menuJsonContent);
        return asyncCallWithArgumentList(QStringLiteral("ShowMenu"), argumentList);
    }

Q_SIGNALS: // SIGNALS
    void ItemInvoked(const QString &itemId, bool checked);
    void MenuUnregistered();
};

namespace com {
  namespace deepin {
    namespace menu {
      typedef ::DBusMenu Menu;
    }
  }
}
#endif
