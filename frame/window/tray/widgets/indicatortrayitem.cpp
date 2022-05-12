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

#include "indicatortrayitem.h"
//#include "util/utils.h"

#include <QLabel>
#include <QBoxLayout>

#include <QDBusConnection>
#include <QDBusInterface>
#include <QFile>
#include <QDebug>
#include <QJsonDocument>
#include <QJsonObject>
#include <QDBusInterface>
#include <QDBusReply>
#include <QApplication>
#include <thread>

IndicatorTrayItem::IndicatorTrayItem(const QString &indicatorName, QWidget *parent, Qt::WindowFlags f)
    : BaseTrayWidget(parent, f)
    , m_indicatorName(indicatorName)
//    , m_gsettings(Utils::ModuleSettingsPtr("keyboard", QByteArray(), this))
    , m_enableClick(true)
{
    setAttribute(Qt::WA_TranslucentBackground);

    auto layout = new QVBoxLayout(this);
    layout->setContentsMargins(0, 0, 0, 0);
    m_label = new QLabel(this);

    QPalette p = m_label->palette();
    p.setColor(QPalette::Foreground, Qt::white);
    p.setColor(QPalette::Background, Qt::transparent);
    m_label->setPalette(p);

    m_label->setAttribute(Qt::WA_TranslucentBackground);

    layout->addWidget(m_label, 0, Qt::AlignCenter);
    setLayout(layout);

    // register dbus
    auto path = QString("/com/deepin/dde/Dock/Indicator/") + m_indicatorName;
    auto interface =  QString("com.deepin.dde.Dock.Indicator.") + m_indicatorName;
    auto sessionBus = QDBusConnection::sessionBus();
    sessionBus.registerObject(path,
                              interface,
                              this,
                              QDBusConnection::ExportScriptableSlots);

    initDBus(m_indicatorName);
//    if (m_gsettings) {
//        // 显示键盘布局时更新label的状态
//        if (m_gsettings->keys().contains("itemEnable"))
//            enableLabel(m_gsettings->get("itemEnable").toBool());

//        connect(m_gsettings, &QGSettings::changed, this, &IndicatorTrayWidget::onGSettingsChanged);
//    }
}

IndicatorTrayItem::~IndicatorTrayItem()
{
}

QString IndicatorTrayItem::itemKeyForConfig()
{
    return toIndicatorKey(m_indicatorName);
}

void IndicatorTrayItem::updateIcon()
{

}

void IndicatorTrayItem::sendClick(uint8_t buttonIndex, int x, int y)
{
    if (m_enableClick)
        Q_EMIT clicked(buttonIndex, x, y);
}

void IndicatorTrayItem::enableLabel(bool enable)
{
    QPalette p = m_label->palette();
    if (!enable) {
        m_enableClick = false;
        p.setColor(QPalette::Disabled, QPalette::Foreground, Qt::lightGray);
        p.setColor(QPalette::Disabled, QPalette::Background, Qt::transparent);
        m_label->setEnabled(enable);
    } else {
        m_enableClick = true;
        p.setColor(QPalette::Normal, QPalette::BrightText, Qt::white);
        p.setColor(QPalette::Normal, QPalette::Background, Qt::transparent);
        m_label->setEnabled(enable);
    }

    m_label->setPalette(p);
    m_label->update();
}

QPixmap IndicatorTrayItem::icon()
{
    return QPixmap();
}

void IndicatorTrayItem::setPixmapData(const QByteArray &data)
{
    auto rawPixmap = QPixmap::fromImage(QImage::fromData(data));
    rawPixmap.setDevicePixelRatio(devicePixelRatioF());
    m_label->setPixmap(rawPixmap);
}

void IndicatorTrayItem::setText(const QString &text)
{
    m_label->setText(text);
}

void IndicatorTrayItem::onGSettingsChanged(const QString &key)
{
    Q_UNUSED(key);

//    if (m_gsettings && m_gsettings->keys().contains("itemEnable")) {
//        const bool itemEnable = m_gsettings->get("itemEnable").toBool();
//        enableLabel(itemEnable);
    //    }
}

template<typename Func>
void IndicatorTrayItem::featData(const QString &key,
              const QJsonObject &data,
              const char *propertyChangedSlot,
              Func const &callback)
{
    auto dataConfig = data.value(key).toObject();
    auto dbusService = dataConfig.value("dbus_service").toString();
    auto dbusPath = dataConfig.value("dbus_path").toString();
    auto dbusInterface = dataConfig.value("dbus_interface").toString();
    auto isSystemBus = dataConfig.value("system_dbus").toBool(false);
    auto bus = isSystemBus ? QDBusConnection::systemBus() : QDBusConnection::sessionBus();

    QDBusInterface interface(dbusService, dbusPath, dbusInterface, bus, this);

    if (dataConfig.contains("dbus_method")) {
        QString methodName = dataConfig.value("dbus_method").toString();
        auto ratio = qApp->devicePixelRatio();
        QDBusReply<QByteArray> reply = interface.call(methodName.toStdString().c_str(), ratio);
        callback(reply.value());
    }

    qInfo() << dataConfig;
    if (dataConfig.contains("dbus_properties")) {
        auto propertyName = dataConfig.value("dbus_properties").toString();
        auto propertyNameCStr = propertyName.toStdString();
        //propertyInterfaceNames.insert(key, dbusInterface);
        //propertyNames.insert(key, QString::fromStdString(propertyNameCStr));
        QDBusConnection::sessionBus().connect(dbusService,
                                              dbusPath,
                                              "org.freedesktop.DBus.Properties",
                                              "PropertiesChanged",
                                              "sa{sv}as",
                                              this,
                                              propertyChangedSlot);

        // FIXME(sbw): hack for qt dbus property changed signal.
        // see: https://bugreports.qt.io/browse/QTBUG-48008
        QDBusConnection::sessionBus().connect(dbusService,
                                              dbusPath,
                                              dbusInterface,
                                              QString("%1Changed").arg(propertyName),
                                              "s",
                                              this,
                                              propertyChangedSlot);

        qInfo() << dbusService << dbusPath << dbusInterface;
        qInfo() << propertyName << propertyNameCStr.c_str();
        callback(interface.property(propertyNameCStr.c_str()));
    }
}

void IndicatorTrayItem::initDBus(const QString &indicatorName)
{
    QString filepath = QString("/etc/dde-dock/indicator/%1.json").arg(indicatorName);
    QFile confFile(filepath);
    if (!confFile.open(QIODevice::ReadOnly)) {
        qInfo() << "read indicator config Error";
    }

    QJsonDocument doc = QJsonDocument::fromJson(confFile.readAll());
    confFile.close();

    QJsonObject config = doc.object();

    auto delay = config.value("delay").toInt(0);

    QTimer::singleShot(delay, [ = ]() {
        QJsonObject data = config.value("data").toObject();
        if (data.contains("text")) {
            featData("text", data, SLOT(textPropertyChanged(QDBusMessage)), [ = ](QVariant v) {
                if (v.toString().isEmpty()) {
                    Q_EMIT removed();
                    return;
                }
                Q_EMIT delayLoaded();
                setText(v.toString());
                //updateContent();
            });
        }

        if (data.contains("icon")) {
            featData("icon", data, SLOT(iconPropertyChanged(QDBusMessage)), [ = ](QVariant v) {
                if (v.toByteArray().isEmpty()) {
                    Q_EMIT removed();
                    return;
                }
                Q_EMIT delayLoaded();
                setPixmapData(v.toByteArray());
                //updateContent();
            });
        }

        const QJsonObject action = config.value("action").toObject();
        if (!action.isEmpty())
            connect(this, &IndicatorTrayItem::clicked, this, [ = ](uint8_t button_index, int x, int y) {
                std::thread t([=]() -> void {
                    auto triggerConfig = action.value("trigger").toObject();
                    auto dbusService = triggerConfig.value("dbus_service").toString();
                    auto dbusPath = triggerConfig.value("dbus_path").toString();
                    auto dbusInterface = triggerConfig.value("dbus_interface").toString();
                    auto methodName = triggerConfig.value("dbus_method").toString();
                    auto isSystemBus = triggerConfig.value("system_dbus").toBool(false);
                    auto bus = isSystemBus ? QDBusConnection::systemBus() : QDBusConnection::sessionBus();

                    QDBusInterface interface(dbusService, dbusPath, dbusInterface, bus);
                    interface.call(methodName, button_index, x, y);
                });
                t.detach();
            });
    });
}

