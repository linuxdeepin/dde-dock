// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "bluetoothitem.h"
#include "constants.h"
#include "../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"
#include "componments/bluetoothapplet.h"

#include <DApplication>
#include <DDBusSender>
#include <DGuiApplicationHelper>

#include <QPainter>

// menu actions
#define SHIFT       "shift"
#define SETTINGS    "settings"

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

using namespace Dock;

BluetoothItem::BluetoothItem(QWidget *parent)
    : QWidget(parent)
    , m_tipsLabel(new TipsWidget(this))
    , m_applet(new BluetoothApplet(this))
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
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method(QString("ShowModule"))
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
