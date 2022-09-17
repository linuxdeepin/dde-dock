// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
                                       "    <method name=\"resizeDock\">"
                                       "        <arg name=\"offset\" type=\"i\" direction=\"in\"/>"
                                       "        <arg name=\"dragging\" type=\"b\" direction=\"in\"/>"
                                       "    </method>"
                                       "    <method name=\"getPluginKey\">"
                                       "        <arg name=\"pluginName\" type=\"s\" direction=\"in\"/>"
                                       "        <arg name=\"key\" type=\"s\" direction=\"out\"/>"
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

    void resizeDock(int offset, bool dragging);

    QString getPluginKey(const QString &pluginName);

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
