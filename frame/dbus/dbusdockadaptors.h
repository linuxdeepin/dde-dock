/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DBUSDOCKADAPTORS_H
#define DBUSDOCKADAPTORS_H

#include <QtDBus/QtDBus>
#include "window/mainwindow.h"
class MainWindow;
/*
 * Adaptor class for interface com.deepin.dde.Dock
 */

class DBusDockAdaptors: public QDBusAbstractAdaptor
{
    Q_OBJECT
    Q_CLASSINFO("D-Bus Interface", "com.deepin.dde.Dock")
    Q_CLASSINFO("D-Bus Introspection", ""
                                       "  <interface name=\"com.deepin.dde.Dock\">\n"
                                       "    <property access=\"read\" type=\"(iiii)\" name=\"geometry\"/>\n"
                                       "    <signal name=\"geometryChanged\">"
                                                "<arg name=\"geometry\" type=\"(iiii)\"/>"
                                            "</signal>"
                                       "  </interface>\n"
                                       "")

public:
    DBusDockAdaptors(MainWindow *parent);
    virtual ~DBusDockAdaptors();

    MainWindow *parent() const;

public: // PROPERTIES
    Q_PROPERTY(QRect geometry READ geometry NOTIFY geometryChanged)
    QRect geometry() const;

signals:
    void geometryChanged(QRect geometry);
};

#endif //DBUSDOCKADAPTORS
