/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong<zhaolong@uniontech.com>
 *
 * Maintainer:  xiehui<xiehui@uniontech.com>
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

#include "airplanemodeitem.h"
#include "constants.h"
#include "../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"
#include "airplanemodeapplet.h"

#include <DGuiApplicationHelper>
#include <DDBusSender>

#include <QPainter>
#include <QJsonDocument>
#include <QDBusConnection>

DGUI_USE_NAMESPACE

#define SHIFT       "shift"
#define SETTINGS    "settings"

AirplaneModeItem::AirplaneModeItem(QWidget *parent)
    : QWidget(parent)
    , m_tipsLabel(new Dock::TipsWidget(this))
    , m_applet(new AirplaneModeApplet(this))
    , m_airplaneModeInter(new DBusAirplaneMode("com.deepin.daemon.AirplaneMode",
                                               "/com/deepin/daemon/AirplaneMode",
                                               QDBusConnection::systemBus(),
                                               this))
{
    m_tipsLabel->setText(tr("Airplane mode"));
    m_tipsLabel->setVisible(false);
    m_applet->setVisible(false);

    connect(m_applet, &AirplaneModeApplet::enableChanged, m_airplaneModeInter, &DBusAirplaneMode::Enable);
    connect(m_airplaneModeInter, &DBusAirplaneMode::EnabledChanged, this, [this](bool enable) {
        m_applet->setEnabled(enable);
        refreshIcon();
    });

    m_applet->setEnabled(m_airplaneModeInter->enabled());
    refreshIcon();
}

QWidget *AirplaneModeItem::tipsWidget()
{
    return m_tipsLabel;
}

QWidget *AirplaneModeItem::popupApplet()
{
    return m_applet;
}

const QString AirplaneModeItem::contextMenu() const
{
    QList<QVariant> items;
    items.reserve(2);

    QMap<QString, QVariant> shift;
    shift["itemId"] = SHIFT;
    if (m_airplaneModeInter->enabled())
        shift["itemText"] = tr("Turn off");
    else
        shift["itemText"] = tr("Turn on");
    shift["isActive"] = true;
    items.push_back(shift);

    QMap<QString, QVariant> settings;
    settings["itemId"] = SETTINGS;
    settings["itemText"] = tr("Airplane Mode settings");
    settings["isActive"] = true;
    items.push_back(settings);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void AirplaneModeItem::invokeMenuItem(const QString menuId, const bool checked)
{
    Q_UNUSED(checked);

    if (menuId == SHIFT)
        m_airplaneModeInter->Enable(!m_airplaneModeInter->enabled());
    else if (menuId == SETTINGS)
        DDBusSender()
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method(QString("ShowPage"))
        .arg(QString("network"))
        .arg(QString("Airplane Mode"))
        .call();
}

void AirplaneModeItem::refreshIcon()
{
    QString iconString;

    qDebug() << "ariplane enable:" << m_airplaneModeInter->enabled();
    if (m_airplaneModeInter->enabled())
        iconString = "3";
    else
        iconString = "4";


    const auto ratio = devicePixelRatioF();
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconString.append(PLUGIN_MIN_ICON_NAME);

    m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);

    update();
}

void AirplaneModeItem::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    if (position == Dock::Bottom || position == Dock::Top) {
        setMaximumWidth(height());
        setMaximumHeight(QWIDGETSIZE_MAX);
    } else {
        setMaximumHeight(width());
        setMaximumWidth(QWIDGETSIZE_MAX);
    }

    refreshIcon();
}

void AirplaneModeItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_iconPixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / m_iconPixmap.devicePixelRatioF(), m_iconPixmap);
}
