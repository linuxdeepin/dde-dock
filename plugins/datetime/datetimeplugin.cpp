// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "datetimeplugin.h"
#include "../../widgets/tipswidget.h"
#include "../../frame/util/utils.h"

#include <DDBusSender>

#include <QDebug>
#include <QDBusConnectionInterface>

#include <unistd.h>

#define PLUGIN_STATE_KEY "enable"
#define TIME_FORMAT_KEY "Use24HourFormat"
using namespace Dock;
DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent)
    , m_centralWidget(nullptr)
    , m_dateTipsLabel(nullptr)
    , m_refershTimer(nullptr)
    , m_interface(nullptr)
    , m_pluginLoaded(false)
{
    QDBusConnection sessionBus = QDBusConnection::sessionBus();
    sessionBus.connect("org.deepin.dde.Timedate1", "/org/deepin/dde/Timedate1", "org.freedesktop.DBus.Properties",  "PropertiesChanged", this, SLOT(propertiesChanged()));
}

PluginsItemInterface::PluginSizePolicy DatetimePlugin::pluginSizePolicy() const
{
    return PluginsItemInterface::Custom;
}

const QString DatetimePlugin::pluginName() const
{
    return "datetime";
}

const QString DatetimePlugin::pluginDisplayName() const
{
    return tr("Datetime");
}

void DatetimePlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    // transfer config
    QSettings settings("deepin", "dde-dock-datetime");
    if (QFile::exists(settings.fileName())) {
        Dock::DisplayMode mode = displayMode();
        const QString key = QString("pos_%1_%2").arg(pluginName()).arg(mode);
        proxyInter->saveValue(this, key, settings.value(key, mode == Dock::DisplayMode::Fashion ? 6 : -1));
        QFile::remove(settings.fileName());
    }

    if (pluginIsDisable()) {
        return;
    }

    loadPlugin();
}

void DatetimePlugin::loadPlugin()
{
    if (m_pluginLoaded)
        return;

    m_pluginLoaded = true;
    m_dateTipsLabel.reset(new TipsWidget);
    m_refershTimer = new QTimer(this);
    m_dateTipsLabel->setObjectName("datetime");

    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    m_centralWidget.reset(new DatetimeWidget);

    connect(m_centralWidget.data(), &DatetimeWidget::requestUpdateGeometry, [this] { m_proxyInter->itemUpdate(this, pluginName()); });
    connect(m_refershTimer, &QTimer::timeout, this, &DatetimePlugin::updateCurrentTimeString);

    m_proxyInter->itemAdded(this, pluginName());

    pluginSettingsChanged();
}

void DatetimePlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool DatetimePlugin::pluginIsDisable()
{
    return !(m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());
}

int DatetimePlugin::itemSortKey(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    return m_proxyInter->getValue(this, key, 6).toInt();
}

void DatetimePlugin::setSortKey(const QString &itemKey, const int order)
{
    Q_UNUSED(itemKey);

    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);
    m_proxyInter->saveValue(this, key, order);
}

QWidget *DatetimePlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_centralWidget.data();
}

QWidget *DatetimePlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_dateTipsLabel.data();
}

const QString DatetimePlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return "dbus-send --print-reply --dest=org.deepin.dde.Widgets1 /org/deepin/dde/Widgets1 org.deepin.dde.Widgets1.Toggle";
}

const QString DatetimePlugin::itemContextMenu(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    QList<QVariant> items;
    items.reserve(1);

    QMap<QString, QVariant> settings;
    settings["itemId"] = "settings";
    if (m_centralWidget->is24HourFormat())
        settings["itemText"] = tr("12-hour time");
    else
        settings["itemText"] = tr("24-hour time");
    settings["isActive"] = true;
    items.push_back(settings);

    if (!QFile::exists(ICBC_CONF_FILE)) {
        QMap<QString, QVariant> open;
        open["itemId"] = "open";
        open["itemText"] = tr("Time settings");
        open["isActive"] = true;
        items.push_back(open);
    }

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void DatetimePlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "open") {
        DDBusSender()
                .service("org.deepin.dde.ControlCenter1")
                .interface("org.deepin.dde.ControlCenter1")
                .path("/org/deepin/dde/ControlCenter1")
                .method(QString("ShowPage"))
                .arg(QString("datetime"))
                .call();
    } else {
        const bool value = timedateInterface()->property(TIME_FORMAT_KEY).toBool();
        timedateInterface()->setProperty(TIME_FORMAT_KEY, !value);
        m_centralWidget->set24HourFormat(!value);
    }
}

void DatetimePlugin::pluginSettingsChanged()
{
    if (!m_pluginLoaded)
        return;

    const bool value = timedateInterface()->property(TIME_FORMAT_KEY).toBool();

    m_proxyInter->saveValue(this, TIME_FORMAT_KEY, value);
    m_centralWidget->set24HourFormat(value);

    refreshPluginItemsVisible();
}

void DatetimePlugin::updateCurrentTimeString()
{
    const QDateTime currentDateTime = QDateTime::currentDateTime();

    if (m_centralWidget->is24HourFormat())
        m_dateTipsLabel->setText(currentDateTime.date().toString(Qt::SystemLocaleLongDate) + currentDateTime.toString(" HH:mm:ss"));
    else
        m_dateTipsLabel->setText(currentDateTime.date().toString(Qt::SystemLocaleLongDate) + currentDateTime.toString(" hh:mm:ss A"));

    const QString currentString = currentDateTime.toString("yyyy/MM/dd hh:mm");

    if (currentString == m_currentTimeString)
        return;

    m_currentTimeString = currentString;
    //当时间显示格式为12小时制且格式为0:00时，当从9:59变到10:00时，插件宽度需要变化
    m_centralWidget->requestUpdateGeometry();
}

void DatetimePlugin::refreshPluginItemsVisible()
{
    if (!pluginIsDisable()) {

        if (!m_pluginLoaded) {
            loadPlugin();
            return;
        }
        m_proxyInter->itemAdded(this, pluginName());
    } else {
        m_proxyInter->itemRemoved(this, pluginName());
    }
}

void DatetimePlugin::propertiesChanged()
{
    pluginSettingsChanged();
}

QDBusInterface* DatetimePlugin::timedateInterface()
{
    if (!m_interface) {
        if (QDBusConnection::sessionBus().interface()->isServiceRegistered("org.deepin.dde.Timedate1")) {
            m_interface = new QDBusInterface("org.deepin.dde.Timedate1", "/org/deepin/dde/Timedate1", "org.deepin.dde.Timedate1", QDBusConnection::sessionBus(), this);
        } else {
            const QString path = QString("/org/deepin/dde/Accounts1/User%1").arg(QString::number(getuid()));
            QDBusInterface * systemInterface = new QDBusInterface("org.deepin.dde.Accounts1", path, "org.deepin.dde.Accounts1.User",
                                                                  QDBusConnection::systemBus(), this);
            return systemInterface;
        }
    }

    return m_interface;
}
