/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
 *
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

#include "bluetoothitem.h"
#include "constants.h"
#include "../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"
#include "bluetoothapplet.h"

#include <DApplication>
#include <DDBusSender>
#include <DGuiApplicationHelper>

#include <QPainter>

// menu actions
#define SHIFT     "shift"
#define SETTINGS "settings"

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE


BluetoothItem::BluetoothItem(QWidget *parent)
    : QWidget(parent)
    , m_applet(new BluetoothApplet(this))
    , m_timer(new QTimer(this))
{
    m_applet->setVisible(false);
    m_adapterPowered = m_applet->poweredInitState();

    m_devState = m_applet->initDeviceState();

    connect(m_timer, &QTimer::timeout, this, &BluetoothItem::refreshIcon);
    connect(m_applet, &BluetoothApplet::powerChanged, [&](bool powered) {
        m_adapterPowered = powered;
        refreshIcon();
    });
    connect(m_applet, &BluetoothApplet::deviceStateChanged, [&](const Device::State state) {
        m_devState = state;
        refreshIcon();
    });
    connect(m_applet, SIGNAL(noAdapter()), this, SIGNAL(noAdapter()));
    connect(m_applet, SIGNAL(justHasAdapter()), this, SIGNAL(justHasAdapter()));
}

//QWidget *BluetoothItem::tipsWidget()
//{
//}

QWidget *BluetoothItem::popupApplet()
{
    return m_applet->hasAadapter() ? m_applet : nullptr;
}

const QString BluetoothItem::contextMenu() const
{
    QList<QVariant> items;
    if (m_applet->hasAadapter()) {
        items.reserve(2);

        QMap<QString, QVariant> shift;
        shift["itemId"] = SHIFT;
        if (m_adapterPowered)
            shift["itemText"] = tr("Turn off");
        else
            shift["itemText"] = tr("Turn on");
        shift["isActive"] = true;
        items.push_back(shift);

        QMap<QString, QVariant> settings;
        settings["itemId"] = SETTINGS;
        settings["itemText"] = tr("Bluetooth settings");
        settings["isActive"] = true;
        items.push_back(settings);

        QMap<QString, QVariant> menu;
        menu["items"] = items;
        menu["checkableMenu"] = false;
        menu["singleCheck"] = false;
        return QJsonDocument::fromVariant(menu).toJson();
    }
    return QByteArray();
}

void BluetoothItem::invokeMenuItem(const QString menuId, const bool checked)
{
    Q_UNUSED(checked);

    if (menuId == SHIFT) {
        m_applet->setAdapterPowered(!m_adapterPowered);
    }
    else if (menuId == SETTINGS)
        DDBusSender()
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method(QString("ShowModule"))
        .arg(QString("bluetooth"))
        .call();
}

void BluetoothItem::refreshIcon()
{
    if (!m_applet)
        return;

    QString stateString;
    QString iconString;

    if (m_adapterPowered) {
        switch (m_devState) {
        case Device::StateConnected:
            stateString = "active";
            break;
        case Device::StateAvailable: {
            m_timer->start();
            stateString = "waiting";
            iconString = QString("bluetooth-%1-symbolic").arg(stateString);
            const auto ratio = devicePixelRatioF();
            int iconSize = PLUGIN_ICON_MAX_SIZE;
            if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
                iconString.append(PLUGIN_MIN_ICON_NAME);

            m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);

            update();
            return ;
        }
        case Device::StateUnavailable: {
            stateString = "disable";
        }      break;
        }
    } else {
        stateString = "disable";
    }

    m_timer->stop();
    iconString = QString("bluetooth-%1-symbolic").arg(stateString);

    const auto ratio = devicePixelRatioF();
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconString.append(PLUGIN_MIN_ICON_NAME);

    m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);

    update();
}

bool BluetoothItem::hasAdapter()
{
    return m_applet->hasAadapter();
}

void BluetoothItem::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    // 保持横纵比
    if (position == Dock::Bottom || position == Dock::Top) {
        setMaximumWidth(height());
        setMaximumHeight(QWIDGETSIZE_MAX);
    } else {
        setMaximumHeight(width());
        setMaximumWidth(QWIDGETSIZE_MAX);
    }

    refreshIcon();
}

void BluetoothItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_iconPixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / m_iconPixmap.devicePixelRatioF(), m_iconPixmap);
    if (m_devState == Device::StateAvailable) {
        QTime time = QTime::currentTime();
        painter.rotate((time.second() + (time.msec() / 1000.0)) * 6.0);
    }
}

