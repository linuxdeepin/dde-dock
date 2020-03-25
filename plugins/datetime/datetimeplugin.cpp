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

#include "datetimeplugin.h"

#include <QApplication>
#include <QDesktopWidget>
#include <DBlurEffectWidget>
#include <QComboBox>
#include <QLineEdit>
#include <QPushButton>
#include <QHBoxLayout>
#include <QVBoxLayout>
#include <DDBusSender>
#include <QDebug>
#include <QStyleFactory>

#include "ddialog.h"
#include "dlineedit.h"
#include "dpushbutton.h"
#include "dlabel.h"
#include "dinputdialog.h"
#include "dboxwidget.h"
#include "dhidpihelper.h"

DWIDGET_USE_NAMESPACE

#define PLUGIN_STATE_KEY "enable"
#define TIME_FORMAT_KEY "24HourFormat"
#define DATETIME_FORMAT_KEY "format"

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent),

      m_dateTipsLabel(new TipsWidget),
      m_refershTimer(new QTimer(this)),
      m_settings("deepin", "dde-dock-datetime")
{
    m_dateTipsLabel->setObjectName("datetime");
    m_dateTipsLabel->setStyleSheet("color:white;"
                                   "padding:0px 3px;");

    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    m_centralWidget = new DatetimeWidget;

    connect(m_centralWidget, &DatetimeWidget::requestUpdateGeometry, [this] { m_proxyInter->itemUpdate(this, pluginName()); });

    connect(m_refershTimer, &QTimer::timeout, this, &DatetimePlugin::updateCurrentTimeString);
}

DatetimePlugin::~DatetimePlugin()
{
    delete m_refershTimer;
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

    if (!pluginIsDisable())
    {
        m_proxyInter->itemAdded(this, pluginName());
        m_centralWidget->set24HourFormat(m_settings.value(TIME_FORMAT_KEY, true).toBool());
        m_centralWidget->setDateTimeFormat(m_settings.value(DATETIME_FORMAT_KEY, "hh:mm:ss\nd MMM yyyy").toString());
    }
}

void DatetimePlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, PLUGIN_STATE_KEY, pluginIsDisable());

    if (!pluginIsDisable())
        m_proxyInter->itemAdded(this, pluginName());
    else
        m_proxyInter->itemRemoved(this, pluginName());
}

bool DatetimePlugin::pluginIsDisable()
{
    return !(m_proxyInter->getValue(this, PLUGIN_STATE_KEY, true).toBool());
}

int DatetimePlugin::itemSortKey(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    Dock::DisplayMode mode = displayMode();
    const QString key = QString("pos_%1").arg(mode);

    if (mode == Dock::DisplayMode::Fashion)
        return m_proxyInter->getValue(this, key, 2).toInt();
    else
        return m_proxyInter->getValue(this, key, 5).toInt();
}

void DatetimePlugin::setSortKey(const QString &itemKey, const int order)
{
    Q_UNUSED(itemKey);

    const QString key = QString("pos_%1").arg(displayMode());
    m_proxyInter->saveValue(this, key, order);
}

QWidget *DatetimePlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_centralWidget;
}

QWidget *DatetimePlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_dateTipsLabel;
}

const QString DatetimePlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return "dbus-send --print-reply --dest=com.deepin.Calendar /com/deepin/Calendar com.deepin.Calendar.RaiseWindow";
}

const QString DatetimePlugin::itemContextMenu(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    QList<QVariant> items;
    items.reserve(1);

    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();

    if (displayMode == Dock::Efficient)
    {
        if (position == Dock::Top || position == Dock::Bottom)
        {
            QMap<QString, QVariant> format;
            format["itemId"] = "format";
            format["itemText"] = tr("Time Format");
            format["isActive"] = true;
            items.push_back(format);
        }
        else
        {
            QMap<QString, QVariant> settings;
            settings["itemId"] = "settings";
            if (m_centralWidget->is24HourFormat())
                settings["itemText"] = tr("12 Hour Time");
            else
                settings["itemText"] = tr("24 Hour Time");
            settings["isActive"] = true;
            items.push_back(settings);
        }
    }
    else
    {
        QMap<QString, QVariant> settings;
        settings["itemId"] = "settings";
        if (m_centralWidget->is24HourFormat())
            settings["itemText"] = tr("12 Hour Time");
        else
            settings["itemText"] = tr("24 Hour Time");
        settings["isActive"] = true;
        items.push_back(settings);
    }

    QMap<QString, QVariant> open;
    open["itemId"] = "open";
    open["itemText"] = tr("Time Settings");
    open["isActive"] = true;
    items.push_back(open);

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

    if (menuId == "open")
    {
        DDBusSender()
            .service("com.deepin.dde.ControlCenter")
            .interface("com.deepin.dde.ControlCenter")
            .path("/com/deepin/dde/ControlCenter")
            .method(QString("ShowModule"))
            .arg(QString("datetime"))
            .call();
    }

    if (menuId == "format")
    {
        dialogFormat();
    }

    if (menuId == "settings")
    {
        const bool value = m_settings.value(TIME_FORMAT_KEY, true).toBool();
        m_settings.setValue(TIME_FORMAT_KEY, !value);
        m_centralWidget->set24HourFormat(!value);
    }
}

void DatetimePlugin::pluginSettingsChanged()
{
    m_centralWidget->set24HourFormat(m_settings.value(TIME_FORMAT_KEY, true).toBool());

    if (!pluginIsDisable())
        m_proxyInter->itemAdded(this, pluginName());
    else
        m_proxyInter->itemRemoved(this, pluginName());
}

void DatetimePlugin::updateCurrentTimeString()
{
    const QDateTime currentDateTime = QDateTime::currentDateTime();

    if (m_centralWidget->is24HourFormat())
        m_dateTipsLabel->setText(currentDateTime.date().toString(Qt::SystemLocaleLongDate) + currentDateTime.toString(" HH:mm:ss"));
    else
        m_dateTipsLabel->setText(currentDateTime.date().toString(Qt::SystemLocaleLongDate) + currentDateTime.toString(" hh:mm:ss A"));

    const QString currentString = currentDateTime.toString("ss");

    if (currentString == m_currentTimeString)
        return;

    m_currentTimeString = currentString;
    m_centralWidget->update();
}

void DatetimePlugin::dialogFormat()
{
    QIcon m_dialogIcon;
    QString format;
    DInputDialog *dialog = new DInputDialog();

    dialog->setInputMode(DInputDialog::InputMode::TextInput);

    dialog->setBackgroundColor(DBlurEffectWidget::DarkColor);

    dialog->setTitle(tr("Time display format"));

    m_dialogIcon.addPixmap(DHiDPIHelper::loadNxPixmap(":/icons/resources/icons/clock.svg"));
    dialog->setIcon(m_dialogIcon, QSize(64, 64));

    format = m_settings.value(DATETIME_FORMAT_KEY, "hh:mm:ss\nd MMM yyyy").toString();
    format = format.replace("\n", "\\n");
    dialog->setTextValue(format);
    dialog->addSpacing(10);

    DHBoxWidget *hbox = new DHBoxWidget;
    hbox->addWidget(new DLabel(tr("Such as:")));
    hbox->addWidget(new DLabel("hh:mm:ss\\nddd, d MMM yyyy\nHH:mm AP ddd\\nyyyy/M/d"));
    dialog->addContent(hbox);

    connect(dialog, SIGNAL(okButtonClicked()), dialog, SLOT(accept()));
    connect(dialog, SIGNAL(cancelButtonClicked()), dialog, SLOT(reject()));

    if (dialog->exec() == DDialog::Accepted)
    {
        format = dialog->textValue();
        format = format.replace("\\n", "\n");
        m_settings.setValue(DATETIME_FORMAT_KEY, format);
        m_centralWidget->setDateTimeFormat(format);
        m_centralWidget->update();
    }

    dialog->close();
}
