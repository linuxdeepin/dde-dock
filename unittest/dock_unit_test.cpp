/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng@uniontech.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng@uniontech.com>
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

#include <QtTest>
#include <QDBusInterface>
#include <QDBusMetaType>
#include <QDBusMessage>
#include <QDBusArgument>

#include <com_deepin_daemon_display.h>

#include "dock_unit_test.h"

DockUnitTest::DockUnitTest()
{
    qDBusRegisterMetaType<ScreenRect>();
}

DockUnitTest::~DockUnitTest()
{

}
/**
 * @brief DockUnitTest::dock_geometry_test   比较任务栏自身的位置和通知给后端的位置是否吻合
 */
void DockUnitTest::dock_geometry_check()
{
    ScreenRect daemonDockRect, dockRect;

    {
        QDBusInterface inter("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", "org.freedesktop.DBus.Properties");
        QString interface = "com.deepin.dde.daemon.Dock";
        QString arg = "FrontendWindowRect";
        QDBusMessage msg = inter.call("Get", interface, arg);

        QVariant var = msg.arguments().first();
        QDBusVariant dbvFirst = var.value<QDBusVariant>();
        QVariant vFirst = dbvFirst.variant();
        QDBusArgument dbusArgs = vFirst.value<QDBusArgument>();

        dbusArgs >> daemonDockRect;
        qDebug() << daemonDockRect;
    }

    {
        QDBusInterface inter("com.deepin.dde.Dock", "/com/deepin/dde/Dock", "org.freedesktop.DBus.Properties");
        QString interface = "com.deepin.dde.Dock";
        QString arg = "geometry";
        QDBusMessage msg = inter.call("Get", interface, arg);

        QVariant var = msg.arguments().first();
        QDBusVariant dbvFirst = var.value<QDBusVariant>();
        QVariant vFirst = dbvFirst.variant();
        QDBusArgument dbusArgs = vFirst.value<QDBusArgument>();

        dbusArgs >> dockRect;
        qDebug() << dockRect;
    }

    QCOMPARE(daemonDockRect, dockRect);
}

QTEST_APPLESS_MAIN(DockUnitTest)

#include "dock_unit_test.moc"
