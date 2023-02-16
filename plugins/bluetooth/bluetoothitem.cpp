// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "bluetoothitem.h"
#include "adaptersmanager.h"
#include "constants.h"
#include "../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"
#include "componments/bluetoothapplet.h"

#include <DApplication>
#include <DDBusSender>
#include <DGuiApplicationHelper>
#include <adapter.h>

#include <QPainter>

// menu actions
#define SHIFT       "shift"
#define SETTINGS    "settings"

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

using namespace Dock;

BluetoothItem::BluetoothItem(AdaptersManager *adapterManager, QWidget *parent)
    : QWidget(parent)
    , m_tipsLabel(new TipsWidget(this))
    , m_applet(new BluetoothApplet(adapterManager, this))
    , m_devState(Device::State::StateUnavailable)
    , m_adapterPowered(m_applet->poweredInitState())
{
    setAccessibleName("BluetoothPluginItem");
    m_applet->setVisible(false);
    m_tipsLabel->setVisible(false);
    refreshIcon();

    connect(m_applet, &BluetoothApplet::powerChanged, [ & ](bool powered) {
        m_adapterPowered = powered;
        refreshIcon();
    });
    connect(m_applet, &BluetoothApplet::deviceStateChanged, [ & ](const Device *device) {
        m_devState = device->state();
        refreshIcon();
        refreshTips();
    });
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &BluetoothItem::refreshIcon);
    connect(m_applet, &BluetoothApplet::noAdapter, this, &BluetoothItem::noAdapter);
    connect(m_applet, &BluetoothApplet::justHasAdapter, this, &BluetoothItem::justHasAdapter);
    connect(m_applet, &BluetoothApplet::requestHide, this, &BluetoothItem::requestHide);
}

QWidget *BluetoothItem::tipsWidget()
{
    refreshTips();
    return m_tipsLabel;
}

QWidget *BluetoothItem::popupApplet()
{
    if (m_applet && m_applet->hasAadapter())
        m_applet->setAdapterRefresh();
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
        shift["isActive"] = !m_applet->airplaneModeEnable();
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
    } else if (menuId == SETTINGS) {
        DDBusSender()
        .service("org.deepin.dde.ControlCenter1")
        .interface("org.deepin.dde.ControlCenter1")
        .path("/org/deepin/dde/ControlCenter1")
        .method(QString("ShowPage"))
        .arg(QString("bluetooth"))
        .call();
    }
}

void BluetoothItem::refreshIcon()
{
    QString stateString;
    QString iconString;

    if (m_adapterPowered) {
        if (m_applet->connectedDevicesName().isEmpty()) {
            stateString = "disable";
        } else {
            stateString = "active";
        }
    } else {
        stateString = "disable";
    }

    iconString = QString("bluetooth-%1-symbolic").arg(stateString);

    const qreal ratio = devicePixelRatioF();
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconString.append(PLUGIN_MIN_ICON_NAME);

    m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);

    update();
}

void BluetoothItem::refreshTips()
{
    QString tipsText;

    if (m_adapterPowered) {
        if (!m_applet->connectedDevicesName().isEmpty() && m_devState != Device::StateAvailable) {
            QStringList textList;
            for (QString devName : m_applet->connectedDevicesName()) {
                textList << tr("%1 connected").arg(devName);
            }
            m_tipsLabel->setTextList(textList);
            return;
        }
        if (m_devState == Device::StateAvailable) {
            tipsText = tr("Connecting...");
        } else {
            tipsText = tr("Bluetooth");
        }
    } else {
        tipsText = tr("Turned off");
    }

    m_tipsLabel->setText(tipsText);
}


bool BluetoothItem::hasAdapter()
{
    return m_applet->hasAadapter();
}

bool BluetoothItem::isPowered()
{
    if (!m_applet->hasAadapter())
        return false;

    QList<const Adapter *> adapters = m_applet->adaptersManager()->adapters();
    for (const Adapter *adapter : adapters) {
        if (adapter->powered())
            return true;
    }
    return false;
}

void BluetoothItem::resizeEvent(QResizeEvent *event)
{
    QWidget::resizeEvent(event);

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

void BluetoothItem::paintEvent(QPaintEvent *event)
{
    QWidget::paintEvent(event);

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_iconPixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / m_iconPixmap.devicePixelRatioF(), m_iconPixmap);
}
