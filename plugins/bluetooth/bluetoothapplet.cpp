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

#include "bluetoothapplet.h"
#include "componments/switchitem.h"
#include "componments/deviceitem.h"
#include "componments/adapter.h"
#include "componments/switchitem.h"
#include "componments/adaptersmanager.h"
#include "componments/adapteritem.h"
#include "componments/bluetoothconstants.h"

#include <DApplicationHelper>
#include <DDBusSender>

#include <QLabel>
#include <QVBoxLayout>

DGUI_USE_NAMESPACE

extern void initFontColor(QWidget *widget)
{
    if (!widget)
        return;

    auto fontChange = [&](QWidget *widget){
        QPalette defaultPalette = widget->palette();
        defaultPalette.setBrush(QPalette::WindowText, defaultPalette.brightText());
        widget->setPalette(defaultPalette);
    };

    fontChange(widget);

    QObject::connect(DApplicationHelper::instance(), &DApplicationHelper::themeTypeChanged, widget, [=]{
        fontChange(widget);
    });
}

BluetoothApplet::BluetoothApplet(QWidget *parent)
    : QScrollArea(parent)
    , m_line(new HorizontalSeparator(this))
    , m_appletName(new QLabel(this))
    , m_centralWidget(new QWidget)
    , m_centrealLayout(new QVBoxLayout)
    , m_adapterLayout(new QVBoxLayout)
    , m_menueLayout(new QHBoxLayout)
    , m_openControlCenter(new MenueItem(this))
    , m_adaptersManager(new AdaptersManager(this))
{
    m_line->setVisible(false);

    QFont defaultFont = font();
    auto titlefont = QFont(defaultFont.family(), defaultFont.pointSize() + 2);

    m_appletName->setText(tr("Bluetooth"));
    m_appletName->setFont(titlefont);
    initFontColor(m_appletName);
    m_appletName->setVisible(false);

    m_openControlCenter->setText(tr("Bluetooth settings"));
    initFontColor(m_openControlCenter);
    m_openControlCenter->setFixedHeight(ITEMHEIGHT);
    m_openControlCenter->setVisible(false);

    QHBoxLayout *appletNameLayout = new QHBoxLayout;
    appletNameLayout->setMargin(0);
    appletNameLayout->setSpacing(0);
    appletNameLayout->addSpacing(MARGIN);
    appletNameLayout->addWidget(m_appletName);
    appletNameLayout->addStretch();

    m_menueLayout->setMargin(0);
    m_menueLayout->setSpacing(0);
    m_menueLayout->addSpacing(MARGIN);

    m_centrealLayout->setMargin(0);
    m_centrealLayout->setSpacing(0);
    m_centrealLayout->addLayout(appletNameLayout);
    m_centrealLayout->addWidget(m_line);
    m_centrealLayout->addLayout(m_adapterLayout);
    m_centrealLayout->addLayout(m_menueLayout);
    m_centralWidget->setLayout(m_centrealLayout);
    m_centralWidget->setFixedWidth(POPUPWIDTH);
    m_centralWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Preferred);

    setFixedWidth(POPUPWIDTH);
    setWidget(m_centralWidget);
    setFrameShape(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);

    connect(m_adaptersManager, &AdaptersManager::adapterIncreased, this, &BluetoothApplet::addAdapter);
    connect(m_adaptersManager, &AdaptersManager::adapterDecreased, this, &BluetoothApplet::removeAdapter);
    connect(m_openControlCenter, &MenueItem::clicked, []{
        DDBusSender()
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method(QString("ShowModule"))
        .arg(QString("bluetooth"))
        .call();
    });
}

void BluetoothApplet::setAdapterPowered(bool powered)
{
    for (AdapterItem *adapterItem : m_adapterItems) {
        if (adapterItem)
            adapterItem->setPowered(powered);
    }
}

bool BluetoothApplet::poweredInitState()
{
    return m_adaptersManager->defaultAdapterInitPowerState();
}

bool BluetoothApplet::hasAadapter()
{
    return m_adaptersManager->adaptersCount();
}

QStringList BluetoothApplet::connectedDevsName()
{
    QStringList devicesName;
    for (AdapterItem *adapterItem : m_adapterItems) {
        if (adapterItem) {
            devicesName << adapterItem->connectedDevsName();
        }
    }
    return devicesName;
}

void BluetoothApplet::onPowerChanged()
{
    bool powerState = false;
    for (AdapterItem *adapterItem : m_adapterItems) {
         if (adapterItem->isPowered()) {
             powerState = true;
             break;
         }
    }
    emit powerChanged(powerState);
}

void BluetoothApplet::onDeviceStateChanged()
{
    Device::State deviceState = Device::StateUnavailable;
    for (AdapterItem *adapterItem : m_adapterItems) {
        if (Device::StateAvailable == adapterItem->currentDeviceState()) {
            deviceState = Device::StateAvailable;
            continue;
        }
        if (Device::StateConnected == adapterItem->currentDeviceState()) {
            deviceState = Device::StateConnected;
            break;
        }
    }

    emit deviceStateChanged(deviceState);
    updateView();
}

void BluetoothApplet::addAdapter(Adapter *adapter)
{
    if (!adapter)
        return;

    if (!m_adapterItems.size()) {
        emit justHasAdapter();
    }

    QString adapterId = adapter->id();
    //dde-session-daemon重启的时候，同一个Id的蓝牙设备会再次添加一次，因此需要先移除以前的
    if (m_adapterItems.contains(adapterId)) {
       AdapterItem *adapterItem = m_adapterItems.value(adapterId);
       if (adapterItem) {
           m_adapterLayout->removeWidget(adapterItem);
           delete  adapterItem;
           m_adapterItems.remove(adapterId);
       }
    }

    auto adatpterItem = new AdapterItem(m_adaptersManager, adapter, this);
    m_adapterItems[adapterId] = adatpterItem;
    m_adapterLayout->addWidget(adatpterItem);
    getDevieInitStatus(adatpterItem);

    connect(adatpterItem, &AdapterItem::deviceStateChanged, this, &BluetoothApplet::onDeviceStateChanged);
    connect(adatpterItem, &AdapterItem::powerChanged, this, &BluetoothApplet::onPowerChanged);
    connect(adatpterItem, &AdapterItem::sizeChange, this, &BluetoothApplet::updateView);

    updateView();
}

void BluetoothApplet::removeAdapter(Adapter *adapter)
{
    if (adapter) {
        QString adapterId = adapter->id();
        AdapterItem *adapterItem = m_adapterItems.value(adapterId);
        if (adapterItem) {
            m_adapterLayout->removeWidget(adapterItem);
            delete  adapterItem;
            m_adapterItems.remove(adapterId);
            updateView();
            if (!m_adapterItems.size())
                emit noAdapter();
        }
    }
}

void BluetoothApplet::updateView()
{
    int contentHeight = 0;
    int itemCount = 0;
    bool isAdapterConnected = true;
    for (AdapterItem *adapterItem : m_adapterItems) {
        if (adapterItem) {
            contentHeight += CONTROLHEIGHT;
            if (adapterItem->isPowered())
                itemCount += adapterItem->deviceCount();
            if (adapterItem->connectedDevsName().size())
                isAdapterConnected = false;
        }
    }

    m_openControlCenter->setVisible(isAdapterConnected);
    if (isAdapterConnected) {
        m_menueLayout->addWidget(m_openControlCenter);
        contentHeight += ITEMHEIGHT;
    } else
        m_menueLayout->removeWidget(m_openControlCenter);

    int adaptersCnt = m_adapterItems.size();
    if (adaptersCnt > 1) {
        m_line->setVisible(true);
        m_appletName->setVisible(true);
    } else {
        m_line->setVisible(false);
        m_appletName->setVisible(false);
    }

    if (adaptersCnt > 1)
        contentHeight += m_appletName->height();

    contentHeight += itemCount * ITEMHEIGHT;
    m_centralWidget->setFixedHeight(contentHeight);
    setFixedHeight(itemCount <= 10 ? contentHeight : 10 * ITEMHEIGHT);
    setVerticalScrollBarPolicy(itemCount <= 10 ? Qt::ScrollBarAlwaysOff : Qt::ScrollBarAlwaysOn);
}

void BluetoothApplet::getDevieInitStatus(AdapterItem *item)
{
    if (!item)
        return;

    bool powered = item->isPowered();
    for (AdapterItem *adapterItem : m_adapterItems) {
        if (adapterItem != item) {
            if (adapterItem->isPowered()) {
                powered = true;
                break;
            }
        }
    }
    emit powerChanged(powered);

    Device::State deviceState = item->initDeviceState();
    Device::State otherDeviceState = Device::StateUnavailable;
    for (AdapterItem *adapterItem : m_adapterItems) {
        if (adapterItem != item) {
            if (Device::StateAvailable == adapterItem->currentDeviceState()) {
                otherDeviceState = Device::StateAvailable;
                continue;
            }
            if (Device::StateConnected == adapterItem->currentDeviceState()) {
                otherDeviceState = Device::StateConnected;
                break;
            }
        }
    }

    switch (deviceState) {
        case Device::StateConnected:
            emit deviceStateChanged(deviceState);
            break;
        case Device::StateUnavailable:
            emit deviceStateChanged(otherDeviceState);
            break;
        case Device::StateAvailable: {
            if (otherDeviceState != Device::StateConnected)
                emit deviceStateChanged(deviceState);
            else
                emit deviceStateChanged(otherDeviceState);
        }
            break;
    }
}

void BluetoothApplet::setAdapterRefresh()
{
    for (AdapterItem *adapterItem : m_adapterItems) {
        if (adapterItem)
            adapterItem->refresh();
    }
}
