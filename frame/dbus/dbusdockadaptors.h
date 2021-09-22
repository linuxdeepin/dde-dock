/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#ifndef DBUSDOCKADAPTORS_H
#define DBUSDOCKADAPTORS_H

#include <QtDBus/QtDBus>

#include "mainwindow.h"

/*
 * Adaptor class for interface com.deepin.dde.Dock
 */
class QGSettings;
class DBusDockAdaptors: public QDBusAbstractAdaptor
{
    Q_OBJECT
    Q_CLASSINFO("D-Bus Interface", "com.deepin.dde.Dock")
    Q_CLASSINFO("D-Bus Introspection", ""
                                       "  <interface name=\"com.deepin.dde.Dock\">\n"
                                       "    <property access=\"read\" type=\"(iiii)\" name=\"geometry\"/>\n"
                                       "    <property access=\"readwrite\" type=\"b\" name=\"showInPrimary\"/>\n"
                                       "    <method name=\"callShow\"/>"
                                       "    <method name=\"ReloadPlugins\"/>"
                                       "    <method name=\"GetLoadedPlugins\">"
                                       "        <arg name=\"list\" type=\"as\" direction=\"out\"/>"
                                       "    </method>"
                                       "    <method name=\"getPluginVisible\">"
                                       "        <arg name=\"pluginName\" type=\"s\" direction=\"in\"/>"
                                       "        <arg name=\"visible\" type=\"b\" direction=\"out\"/>"
                                       "    </method>"
                                       "    <method name=\"setPluginVisible\">"
                                       "        <arg name=\"pluginName\" type=\"s\" direction=\"in\"/>"
                                       "        <arg name=\"visible\" type=\"b\" direction=\"in\"/>"
                                       "    </method>"
                                       "    <signal name=\"pluginVisibleChanged\">"
                                       "        <arg type=\"s\"/>"
                                       "        <arg type=\"b\"/>"
                                       "    </signal>"
                                       "  </interface>\n"
                                       "")
    Q_PROPERTY(QRect geometry READ geometry NOTIFY geometryChanged)
    Q_PROPERTY(bool showInPrimary READ showInPrimary WRITE setShowInPrimary NOTIFY showInPrimaryChanged)

public:
    explicit DBusDockAdaptors(MainWindow *parent);
    virtual ~DBusDockAdaptors();

    MainWindow *parent() const;

public Q_SLOTS: // METHODS
    void callShow();
    void ReloadPlugins();

    QStringList GetLoadedPlugins();

    bool getPluginVisible(const QString &pluginName);
    void setPluginVisible(const QString &pluginName, bool visible);

public: // PROPERTIES
    QRect geometry() const;

    bool showInPrimary() const;
    void setShowInPrimary(bool showInPrimary);

signals:
    void geometryChanged(QRect geometry);
    void showInPrimaryChanged(bool);
    void pluginVisibleChanged(const QString &pluginName, bool visible);

private:
    bool isPluginValid(const QString &name);

private:
    QGSettings *m_gsettings;
};

#endif //DBUSDOCKADAPTORS
