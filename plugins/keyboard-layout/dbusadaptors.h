// Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DBUSADAPTORS_H
#define DBUSADAPTORS_H

#include <QMenu>
#include <QtDBus/QtDBus>
#include "org_deepin_dde_inputdevice1_keyboard.h"

using Keyboard = org::deepin::dde::inputdevice1::Keyboard;
class QGSettings;

class DBusAdaptors : public QDBusAbstractAdaptor
{
    Q_OBJECT
    Q_CLASSINFO("D-Bus Interface", "org.deepin.dde.Dock1.KeyboardLayout")
//    Q_CLASSINFO("D-Bus Introspection", ""
//                "  <interface name=\"org.deepin.dde.Dock1.KeyboardLayout\">\n"
//                "    <property access=\"read\" type=\"s\" name=\"layout\"/>\n"
//                "    <signal name=\"layoutChanged\">"
//                "        <arg name=\"layout\" type=\"s\"/>"
//                "    </signal>"
//                "  </interface>\n"
//                "")

public:
     explicit DBusAdaptors(QObject *parent = nullptr);
    ~DBusAdaptors();

public:
    Q_PROPERTY(QString layout READ layout WRITE setLayout NOTIFY layoutChanged)
    QString layout() const;
    void setLayout(const QString &str);

    Keyboard *getCurrentKeyboard();

public slots:
    void onClicked(int button, int x, int y);

signals:
    void layoutChanged(QString text);

private slots:
    void onCurrentLayoutChanged(const QString & value);
    void onUserLayoutListChanged(const QStringList & value);
    void initAllLayoutList();
    void refreshMenu();
    void refreshMenuSelection();
    void handleActionTriggered(QAction *action);

private slots:
    void onGSettingsChanged(const QString &key);

private:
    QString duplicateCheck(const QString &kb);

private:
    Keyboard *m_keyboard;
    QMenu *m_menu;
    QAction *m_addLayoutAction;

    QString m_currentLayoutRaw;
    QString m_currentLayout;
    QStringList m_userLayoutList;
    KeyboardLayoutList m_allLayoutList;
    const QGSettings *m_gsettings;
};

#endif
