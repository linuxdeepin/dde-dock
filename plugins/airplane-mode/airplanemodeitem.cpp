// SPDX-FileCopyrightText: 2020 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "airplanemodeitem.h"
#include "constants.h"
#include "tipswidget.h"
#include "imageutil.h"
#include "utils.h"
#include "airplanemodeapplet.h"

#include <DGuiApplicationHelper>
#include <DDBusSender>

#include <QPainter>
#include <QJsonDocument>
#include <QDBusConnection>
#include <QIcon>

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
    m_tipsLabel->setText(tr("Airplane mode enabled"));
    m_tipsLabel->setVisible(false);
    m_applet->setVisible(false);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &AirplaneModeItem::refreshIcon);
    connect(m_applet, &AirplaneModeApplet::enableChanged, m_airplaneModeInter, &DBusAirplaneMode::Enable);
    connect(m_airplaneModeInter, &DBusAirplaneMode::EnabledChanged, this, [this](bool enable) {
        m_applet->setEnabled(enable);
        refreshIcon();
        Q_EMIT airplaneEnableChanged(enable);
        updateTips();
    });

    m_applet->setEnabled(m_airplaneModeInter->enabled());
    refreshIcon();
    updateTips();
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
    Q_UNUSED(menuId);
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
    if (m_airplaneModeInter->enabled())
        iconString = "airplane-on";
    else
        iconString = "airplane-off";

    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconString.append(PLUGIN_MIN_ICON_NAME);

    const auto ratio = devicePixelRatioF();
    QSize pixmapSize = QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? QSize(20, 20) : (QSize(20, 20) * ratio);
    m_iconPixmap = QIcon::fromTheme(iconString, QIcon::fromTheme("://" + iconString + ".svg")).pixmap(pixmapSize);
    m_iconPixmap.setDevicePixelRatio(ratio);

    update();
}

void AirplaneModeItem::updateTips()
{
    if (m_airplaneModeInter->enabled())
        m_tipsLabel->setText(tr("Airplane mode enabled"));
    else
        m_tipsLabel->setText(tr("Airplane mode disabled"));
}

bool AirplaneModeItem::airplaneEnable()
{
    return m_airplaneModeInter->enabled();
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

    const auto ratio = devicePixelRatioF();

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_iconPixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / ratio, m_iconPixmap);
}
