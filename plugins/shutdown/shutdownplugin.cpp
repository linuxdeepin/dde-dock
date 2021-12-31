/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
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

#include "shutdownplugin.h"
#include "dbus/dbusaccount.h"
#include "../frame/util/utils.h"
#include "../widgets/tipswidget.h"
#include "./dbus/dbuspowermanager.h"

#include <DSysInfo>
#include <DDBusSender>

#include <QIcon>
#include <QSettings>

#define PLUGIN_STATE_KEY "enable"
#define GSETTING_SHOW_SUSPEND "showSuspend"
#define GSETTING_SHOW_HIBERNATE "showHibernate"
#define GSETTING_SHOW_SHUTDOWN "showShutdown"
#define GSETTING_SHOW_LOCK "showLock"

DCORE_USE_NAMESPACE
using namespace Dock;

ShutdownPlugin::ShutdownPlugin(QObject *parent)
    : QObject(parent)
    , m_pluginLoaded(false)
    , m_shutdownWidget(nullptr)
    , m_tipsLabel(new TipsWidget)
    , m_powerManagerInter(new DBusPowerManager("com.deepin.daemon.PowerManager", "/com/deepin/daemon/PowerManager", QDBusConnection::systemBus(), this))
    , m_gsettings(Utils::ModuleSettingsPtr("shutdown", QByteArray(), this))
    , m_sessionShellGsettings(Utils::SettingsPtr("com.deepin.dde.session-shell", "/com/deepin/dde/session-shell/", this))
{
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setAccessibleName("shutdown");
}

const QString ShutdownPlugin::pluginName() const
{
    return "shutdown";
}

const QString ShutdownPlugin::pluginDisplayName() const
{
    return tr("Power");
}

QWidget *ShutdownPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_shutdownWidget.data();
}

QWidget *ShutdownPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    // reset text every time to avoid size of LabelWidget not change after
    // font size be changed in ControlCenter
    m_tipsLabel->setText(tr("Power"));

    return m_tipsLabel.data();
}

void ShutdownPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    // transfer config
    QSettings settings("deepin", "dde-dock-shutdown");
    if (QFile::exists(settings.fileName())) {
        QFile::remove(settings.fileName());
    }

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void ShutdownPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());

    refreshPluginItemsVisible();
}

bool ShutdownPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool();
}

const QString ShutdownPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return QString("dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show");
}

const QString ShutdownPlugin::itemContextMenu(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    QList<QVariant> items;
    items.reserve(6);

    QMap<QString, QVariant> shutdown;
    if (!m_gsettings || (m_gsettings->keys().contains(GSETTING_SHOW_SHUTDOWN) && m_gsettings->get(GSETTING_SHOW_SHUTDOWN).toBool())) {
        shutdown["itemId"] = "Shutdown";
        shutdown["itemText"] = tr("Shut down");
        shutdown["isActive"] = true;
        items.push_back(shutdown);
    }

    QMap<QString, QVariant> reboot;
    reboot["itemId"] = "Restart";
    reboot["itemText"] = tr("Reboot");
    reboot["isActive"] = true;
    items.push_back(reboot);

#ifndef DISABLE_POWER_OPTIONS

    QProcessEnvironment enviromentVar = QProcessEnvironment::systemEnvironment();
    bool can_sleep = enviromentVar.contains("POWER_CAN_SLEEP") ? QVariant(enviromentVar.value("POWER_CAN_SLEEP")).toBool()
                     : valueByQSettings<bool>("Power", "sleep", true) &&
                     m_powerManagerInter->CanSuspend();
    ;
    if (can_sleep) {
        QMap<QString, QVariant> suspend;
        if (!m_gsettings || (m_gsettings->keys().contains(GSETTING_SHOW_SUSPEND) && m_gsettings->get(GSETTING_SHOW_SUSPEND).toBool())) {
            suspend["itemId"] = "Suspend";
            suspend["itemText"] = tr("Suspend");
            suspend["isActive"] = true;
            items.push_back(suspend);
        }
    }

    bool can_hibernate = enviromentVar.contains("POWER_CAN_HIBERNATE") ? QVariant(enviromentVar.value("POWER_CAN_HIBERNATE")).toBool()
                         : checkSwap() && m_powerManagerInter->CanHibernate();

    if (can_hibernate) {
        QMap<QString, QVariant> hibernate;
        if (!m_gsettings || (m_gsettings->keys().contains(GSETTING_SHOW_HIBERNATE) && m_gsettings->get(GSETTING_SHOW_HIBERNATE).toBool())) {
            hibernate["itemId"] = "Hibernate";
            hibernate["itemText"] = tr("Hibernate");
            hibernate["isActive"] = true;
            items.push_back(hibernate);
        }
    }

#endif

    QMap<QString, QVariant> lock;
    if (!m_gsettings || (m_gsettings->keys().contains(GSETTING_SHOW_LOCK) && m_gsettings->get(GSETTING_SHOW_LOCK).toBool())) {
        lock["itemId"] = "Lock";
        lock["itemText"] = tr("Lock");
        lock["isActive"] = true;
        items.push_back(lock);
    }

    QMap<QString, QVariant> logout;
    logout["itemId"] = "Logout";
    logout["itemText"] = tr("Log out");
    logout["isActive"] = true;
    items.push_back(logout);

    if (!QFile::exists(ICBC_CONF_FILE)) {
        // 读取com.deepin.dde.session-shell切换用户配置项
        enum SwitchUserConfig {
            AlwaysShow = 0,
            OnDemand,
            Disabled
        } switchUserConfig = OnDemand;

        if (m_sessionShellGsettings && m_sessionShellGsettings->keys().contains("switchuser")) {
            switchUserConfig = SwitchUserConfig(m_sessionShellGsettings->get("switchuser").toInt());
        }

        // 和登录锁屏界面的逻辑保持一致
        if (AlwaysShow == switchUserConfig ||
                 (OnDemand == switchUserConfig &&
                 (DBusAccount().userList().count() > 1 || DSysInfo::uosType() == DSysInfo::UosType::UosServer))) {
            QMap<QString, QVariant> switchUser;
            switchUser["itemId"] = "SwitchUser";
            switchUser["itemText"] = tr("Switch account");
            switchUser["isActive"] = true;
            items.push_back(switchUser);
        }

#ifndef DISABLE_POWER_OPTIONS
        QMap<QString, QVariant> power;
        power["itemId"] = "power";
        power["itemText"] = tr("Power settings");
        power["isActive"] = true;
        items.push_back(power);
#endif
    }

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void ShutdownPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    // 使得下面逻辑代码延迟200ms执行，保证线程不阻塞
    QTime dieTime = QTime::currentTime().addMSecs(200);
    while (QTime::currentTime() < dieTime)
        QCoreApplication::processEvents(QEventLoop::AllEvents, 200);

    if (menuId == "power") {
        DDBusSender()
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method(QString("ShowModule"))
        .arg(QString("power"))
        .call();
    } else if (menuId == "Lock") {
        if (QFile::exists(ICBC_CONF_FILE)) {
            QDBusMessage send = QDBusMessage::createMethodCall("com.deepin.dde.lockFront", "/com/deepin/dde/lockFront", "com.deepin.dde.lockFront", "SwitchTTYAndShow");
            QDBusConnection conn = QDBusConnection::connectToBus("unix:path=/run/user/1000/bus", "unix:path=/run/user/1000/bus");
            QDBusMessage reply = conn.call(send);
#ifdef QT_DEBUG
            qInfo() << "----------" << reply;
#endif

        } else {
            DDBusSender()
            .service("com.deepin.dde.lockFront")
            .interface("com.deepin.dde.lockFront")
            .path("/com/deepin/dde/lockFront")
            .method(QString("Show"))
            .call();
        }
    } else
        DDBusSender()
        .service("com.deepin.dde.shutdownFront")
        .interface("com.deepin.dde.shutdownFront")
        .path("/com/deepin/dde/shutdownFront")
        .method(QString(menuId))
        .call();
}

void ShutdownPlugin::displayModeChanged(const Dock::DisplayMode displayMode)
{
    Q_UNUSED(displayMode);

    if (!pluginIsDisable()) {
        m_shutdownWidget->update();
    }
}

int ShutdownPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    return m_proxyInter->getValue(this, key, 5).toInt();
}

void ShutdownPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    m_proxyInter->saveValue(this, key, order);
}

void ShutdownPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

void ShutdownPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        qDebug() << "shutdown plugin has been loaded! return";
        return;
    }

    m_pluginLoaded = true;

    m_shutdownWidget.reset(new ShutdownWidget);

    m_proxyInter->itemAdded(this, pluginName());
    displayModeChanged(displayMode());
}

std::pair<bool, qint64> ShutdownPlugin::checkIsPartitionType(const QStringList &list)
{
    std::pair<bool, qint64> result{ false, -1 };

    if (list.length() != 5) {
        return result;
    }

    const QString type{ list[1] };
    const QString size{ list[2] };

    result.first  = type == "partition";
    result.second = size.toLongLong() * 1024.0f;

    return result;
}

qint64 ShutdownPlugin::get_power_image_size()
{
    qint64 size{ 0 };
    QFile  file("/sys/power/image_size");

    if (file.open(QIODevice::Text | QIODevice::ReadOnly)) {
        size = file.readAll().trimmed().toLongLong();
        file.close();
    }

    return size;
}

bool ShutdownPlugin::checkSwap()
{
    if (!valueByQSettings<bool>("Power", "hibernate", false))
        return false;

    bool hasSwap = false;
    QFile file("/proc/swaps");
    if (file.open(QIODevice::Text | QIODevice::ReadOnly)) {
        const QString &body = file.readAll();
        QTextStream    stream(body.toUtf8());
        while (!stream.atEnd()) {
            const std::pair<bool, qint64> result =
                checkIsPartitionType(stream.readLine().simplified().split(
                                         " ", QString::SplitBehavior::SkipEmptyParts));
            qint64 image_size{ get_power_image_size() };

            if (result.first) {
                hasSwap = image_size < result.second;
            }

            if (hasSwap) {
                break;
            }
        }

        file.close();
    } else {
        qDebug() << "open /proc/swaps failed! please check permission!!!";
    }

    return hasSwap;
}

void ShutdownPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, pluginName());
    } else {
        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        m_proxyInter->itemAdded(this, pluginName());
    }
}
