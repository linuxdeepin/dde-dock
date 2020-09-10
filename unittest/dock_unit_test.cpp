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
#include <QThread>
#include <QProcess>

#include <com_deepin_daemon_display.h>

#include "dock_unit_test.h"

#define SLEEP1 QThread::sleep(1);

DockUnitTest::DockUnitTest()
    : m_dockInter(new QDBusInterface("com.deepin.dde.Dock", "/com/deepin/dde/Dock", "org.freedesktop.DBus.Properties"))
    , m_daemonDockInter(new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this))
{
    qDBusRegisterMetaType<ScreenRect>();
}

DockUnitTest::~DockUnitTest()
{
    delete m_dockInter;
    delete m_daemonDockInter;
}

const DockRect DockUnitTest::dockGeometry()
{
    DockRect dockRect;
    QDBusInterface inter("com.deepin.dde.Dock", "/com/deepin/dde/Dock", "org.freedesktop.DBus.Properties");
    QString interface = "com.deepin.dde.Dock";
    QString arg = "geometry";
    QDBusMessage msg = inter.call("Get", interface, arg);

    QVariant var = msg.arguments().first();
    QDBusVariant dbvFirst = var.value<QDBusVariant>();
    QVariant vFirst = dbvFirst.variant();
    QDBusArgument dbusArgs = vFirst.value<QDBusArgument>();
    dbusArgs >> dockRect;

    return dockRect;
}

const DockRect DockUnitTest::frontendWindowRect()
{
    DBusDock dockInter("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this);
    return dockInter.frontendWindowRect();
}

void DockUnitTest::setPosition(Dock::Position pos)
{
    m_daemonDockInter->setPosition(pos);
}

/**
 * @brief DockUnitTest::dock_defaultGsettings_check
 * 验证任务栏的默认配置值是否正确，用于检查新安装系统或新创建用户的默认配置文件的正确性。
 * 后续如果还有类似验证，应一并放到这里：
 * 1.任务栏默认显示模式。
 * 2.任务栏默认显示状态。
 * 3.任务栏默认位置。
 */
void DockUnitTest::dock_defaultGsettings_check()
{
    QGSettings setting("com.deepin.dde.dock", "/com/deepin/dde/dock/");

    if (setting.keys().contains("displayMode")) {
        QString currentDisplayMode = setting.get("display-mode").toString();
        QString defaultDisplayMode = "efficient";
        QCOMPARE(currentDisplayMode, defaultDisplayMode);
    }
    if (setting.keys().contains("hideMode")) {
        QString currentHideMode = setting.get("hide-mode").toString();
        QString defaultHideMode = "keep-showing";
        QCOMPARE(currentHideMode, defaultHideMode);
    }
    if (setting.keys().contains("position")) {
        QString currentPosition = setting.get("position").toString();
        QString defaultPosition = "bottom";
        QCOMPARE(currentPosition, defaultPosition);
    }
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

/**
 * @brief DockUnitTest::dock_frontWindowRect_check
 * 原理：将任务栏设置为一直显示模式，然后切换四个位置，查看其frontendWindowRect与实际位置是否一致．
 * 前提条件：　开启系统缩放(大于１即可)
 * Tips:　任务栏在副屏时，frontendWindowRect对应的是上一次在主屏设置过的值．
 * 可测出问题：任务栏在5.3.0.2版本时，开启缩放后，启动器打开时位置和任务栏有重叠，5.3.0.5版本修复了这个问题
 * 对应Bug: https://pms.uniontech.com/zentao/bug-view-42095.html
 */
void DockUnitTest::dock_frontWindowRect_check()
{
    setPosition(Dock::Position::Top);
    SLEEP1;
    QVERIFY(dockGeometry() == frontendWindowRect());

    setPosition(Dock::Position::Right);
    SLEEP1;
    QVERIFY(dockGeometry() == frontendWindowRect());

    setPosition(Dock::Position::Bottom);
    SLEEP1;
    QVERIFY(dockGeometry() == frontendWindowRect());

    setPosition(Dock::Position::Left);
    SLEEP1;
    QVERIFY(dockGeometry() == frontendWindowRect());
}

/**
 * @brief DockUnitTest::dock_multi_process
 * 检查dde-dock是否在没进程存在时能否正常启动，在已有dde-dock进程存在时能否正常退出
 */
void DockUnitTest::dock_multi_process()
{
    QProcess *dockProc = new QProcess();
    dockProc->start("dde-dock");
    connect(dockProc, static_cast<void (QProcess::*)(int, QProcess::ExitStatus)>(&QProcess::finished), this, [=](int exitCode, QProcess::ExitStatus exitStatus) {
        QCOMPARE(exitCode, 255);
        QCOMPARE(exitStatus, QProcess::ExitStatus::NormalExit);
    });
    connect(dockProc, &QProcess::errorOccurred, this, [=](QProcess::ProcessError error) {
        qDebug() << "dde-dock error occurred: " << error;
        QFAIL("control center error occurred");
    });
    dockProc->waitForFinished();

    delete dockProc;
}

/**
 * @brief DockUnitTest::dock_defaultVolume_Check　判断音量实际值是否与默认值是否相等
 * @param defaultVolume 默认音量
 */
void DockUnitTest::dock_defaultVolume_Check(float defaultVolume)
{
    float volume = 0;
    QDBusInterface audioInterface("com.deepin.daemon.Audio", "/com/deepin/daemon/Audio", "com.deepin.daemon.Audio", QDBusConnection::sessionBus(), this);
    QDBusObjectPath defaultPath = audioInterface.property("DefaultSink").value<QDBusObjectPath>();
    if (defaultPath.path() == "/") { //路径为空
        qDebug() << "defaultPath" << defaultPath.path();
        return;
    } else {
         QDBusInterface sinkInterface("com.deepin.daemon.Audio", defaultPath.path(), "com.deepin.daemon.Audio.Sink", QDBusConnection::sessionBus(), this);
         volume = sinkInterface.property("Volume").toFloat() * 100.0f;
    }
    QCOMPARE(volume, defaultVolume);
}
/**
 * @brief DockUnitTest::dock_coreDump_check
 *  间隔一段时间判断dock是不是同一个pid,判断是否一直在崩溃
 *
 */
void DockUnitTest::dock_coreDump_check()
{
    auto process = new QProcess();
    process->start("pidof -s  dde-dock");
    process->waitForFinished();
    QByteArray pid = process->readAllStandardOutput();
    process->close();

    QThread::sleep(1);

    process->start("pidof -s  dde-dock");
    process->waitForFinished();
    QByteArray pid2 = process->readAllStandardOutput();
    process->close();

    QCOMPARE(pid,pid2);

    delete process;
}

/**
 * @brief DockUnitTest::dock_appIconSize_check
 * 判断dbus和gsettings获取的任务栏图标大小是否一致
 */
void DockUnitTest::dock_appIconSize_check()
{
    DBusDock *dockInter = new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this);
    QGSettings *setting = new QGSettings("com.deepin.dde.dock");
    unsigned int iconSize = dockInter->iconSize();
    qDebug() << "Please check the size of icons:" << iconSize;
    QCOMPARE(iconSize, setting->get("icon-size").toUInt());
}

QTEST_APPLESS_MAIN(DockUnitTest)

#include "dock_unit_test.moc"
