// Copyright (C) 2015 Digia Plc and/or its subsidiary(-ies).
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DBUSTRAYMANAGER_H_1467094672
#define DBUSTRAYMANAGER_H_1467094672

#include <QtCore/QObject>
#include <QtCore/QByteArray>
#include <QtCore/QList>
#include <QtCore/QMap>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QVariant>
#include <QtDBus/QtDBus>

typedef QList<quint32> TrayList;

/*
 * Proxy class for interface org.deepin.dde.TrayManager1
 */
class DBusTrayManager: public QDBusAbstractInterface
{
    Q_OBJECT

    Q_SLOT void __propertyChanged__(const QDBusMessage& msg)
    {
        QList<QVariant> arguments = msg.arguments();
        if (3 != arguments.count())
            return;
        QString interfaceName = msg.arguments().at(0).toString();
        if (interfaceName !="org.deepin.dde.TrayManager1")
            return;
        QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
        foreach(const QString &prop, changedProps.keys()) {
        const QMetaObject* self = metaObject();
            for (int i=self->propertyOffset(); i < self->propertyCount(); ++i) {
                QMetaProperty p = self->property(i);
                if (p.name() == prop) {
 	            Q_EMIT p.notifySignal().invoke(this);
                }
            }
        }
   }
public:
    static inline const char *staticInterfaceName()
    { return "org.deepin.dde.TrayManager1"; }

public:
    explicit DBusTrayManager(QObject *parent = 0);

    ~DBusTrayManager();

    Q_PROPERTY(TrayList TrayIcons READ trayIcons NOTIFY TrayIconsChanged)
    inline TrayList trayIcons() const
    { return qvariant_cast< TrayList >(property("TrayIcons")); }

public Q_SLOTS: // METHODS
    inline QDBusPendingReply<> EnableNotification(uint in0, bool in1)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0) << QVariant::fromValue(in1);
        return asyncCallWithArgumentList(QStringLiteral("EnableNotification"), argumentList);
    }

    inline QDBusPendingReply<QString> GetName(uint in0)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0);
        return asyncCallWithArgumentList(QStringLiteral("GetName"), argumentList);
    }

    inline QDBusPendingReply<bool> Manage()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("Manage"), argumentList);
    }

    inline QDBusPendingReply<> RetryManager()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("RetryManager"), argumentList);
    }

    inline QDBusPendingReply<bool> Unmanage()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("Unmanage"), argumentList);
    }

Q_SIGNALS: // SIGNALS
    void Added(uint in0);
    void Changed(uint in0);
    void Inited();
    void Removed(uint in0);
// begin property changed signals
void TrayIconsChanged();
};

namespace com {
  namespace deepin {
    namespace dde {
      typedef ::DBusTrayManager TrayManager;
    }
  }
}
#endif
