/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "dbusdockadaptors.h"
#include <QScreen>

DBusDockAdaptors::DBusDockAdaptors(MainWindow* parent): QDBusAbstractAdaptor(parent)
{
    connect(parent, &MainWindow::panelGeometryChanged, this, [this] {
        emit DBusDockAdaptors::geometryChanged(geometry());
    });
}

DBusDockAdaptors::~DBusDockAdaptors()
{

}

MainWindow *DBusDockAdaptors::parent() const
{
    return static_cast<MainWindow *>(QObject::parent());
}

QRect DBusDockAdaptors::geometry() const
{
    return parent()->panelGeometry();
}

