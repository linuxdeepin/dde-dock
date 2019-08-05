/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     rekols <rekols@foxmail.com>
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

#ifndef DBUSADAPTORS_H
#define DBUSADAPTORS_H

#include <QMenu>
#include <QtDBus/QtDBus>
#include <com_deepin_daemon_inputdevice_keyboard.h>

using Keyboard = com::deepin::daemon::inputdevice::Keyboard;

class DBusAdaptors : public QDBusAbstractAdaptor
{
    Q_OBJECT
    Q_CLASSINFO("D-Bus Interface", "com.deepin.dde.Keyboard")
//    Q_CLASSINFO("D-Bus Introspection", ""
//                "  <interface name=\"com.deepin.dde.Keyboard\">\n"
//                "    <property access=\"read\" type=\"s\" name=\"layout\"/>\n"
//                "    <signal name=\"layoutChanged\">"
//                "        <arg name=\"layout\" type=\"s\"/>"
//                "    </signal>"
//                "  </interface>\n"
//                "")

public:
     DBusAdaptors(QObject *parent = nullptr);
    ~DBusAdaptors();

public:
    Q_PROPERTY(QString layout READ layout NOTIFY layoutChanged)
    QString layout() const;

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
};

#endif
