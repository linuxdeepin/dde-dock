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
#include <QGSettings/QGSettings>

#include <com_deepin_daemon_display.h>
#include <com_deepin_dde_daemon_dock.h>

#include "dock_unit_test.h"

using DBusDock = com::deepin::dde::daemon::Dock;
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
/**
 * @brief DockUnitTest::dock_position_check   比较Dbus和QGSettings获取的坐标信息是否一致
 */
void DockUnitTest::dock_position_check()
{
    DBusDock *dockInter = new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this);
    int nPos = dockInter->position();
    QString postion = "";
    qDebug() << nPos;

    switch (nPos) {
    case 0 :
        postion = "top";
        break;
    case 1:
        postion = "right";
        break;
    case 2:
        postion = "bottom";
        break;
    case 3:
        postion = "left";
        break;
    default:
        break;
    }

    QGSettings *setting = new QGSettings("com.deepin.dde.dock");
    if (setting->keys().contains("position")) {
        qDebug() << setting->get("position");
        QCOMPARE(postion,setting->get("position").toString());
    }
}
/**
 * @brief DockUnitTest::dock_displayMode_check   比较Dbus和QGSettings获取的显示模式是否一致
 */
void DockUnitTest::dock_displayMode_check()
{
    DBusDock *dockInter = new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this);
    int nMode = dockInter->displayMode();
    QString displayMode = "";
    qDebug() << nMode;

    switch (nMode) {
    case 0 :
        displayMode = "fashion";
        break;
    case 1:
        displayMode = "efficient";
        break;
    case 2:
        displayMode = "classic";
        break;
    default:
        break;
    }

    QGSettings *setting = new QGSettings("com.deepin.dde.dock");
    if (setting->keys().contains("displayMode")) {
        qDebug() << setting->get("displayMode");
        QCOMPARE(displayMode,setting->get("displayMode").toString());
    }
}

void DockUnitTest::dock_appItemCount_check()
{
    DBusDock *dockInter = new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this);
    qDebug() << dockInter->entries().size();
    for (auto inter : dockInter->entries()) {
        qDebug() << inter.path();
    }

}

QTEST_APPLESS_MAIN(DockUnitTest)

#include "dock_unit_test.moc"
