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

#include <DApplicationHelper>
DGUI_USE_NAMESPACE

extern int ControlHeight;
extern int ItemHeight;
const int Width = 200;

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
    , m_centralWidget(new QWidget(this))
    , m_centrealLayout(new QVBoxLayout)
    , m_adaptersManager(new AdaptersManager(this))
{
    m_line->setVisible(false);

    auto defaultFont = font();
    auto titlefont = QFont(defaultFont.family(), defaultFont.pointSize() + 2);

    m_appletName->setText(tr("Bluetooth"));
    m_appletName->setFont(titlefont);
    initFontColor(m_appletName);
    m_appletName->setVisible(false);

    auto appletNameLayout = new QHBoxLayout;
    appletNameLayout->setMargin(0);
    appletNameLayout->setSpacing(0);
    appletNameLayout->addSpacing(12);
    appletNameLayout->addWidget(m_appletName);
    appletNameLayout->addStretch();

    m_centrealLayout->setMargin(0);
    m_centrealLayout->setSpacing(0);
    m_centrealLayout->addLayout(appletNameLayout);
    m_centrealLayout->addWidget(m_line);
    m_centralWidget->setLayout(m_centrealLayout);
    m_centralWidget->setFixedWidth(Width);
    m_centralWidget->setSizePolicy(QSizePolicy::Fixed, QSizePolicy::Preferred);

    setFixedWidth(Width);
    setWidget(m_centralWidget);
    setFrameShape(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);

    connect(m_adaptersManager, &AdaptersManager::adapterIncreased, this, &BluetoothApplet::addAdapter);
    connect(m_adaptersManager, &AdaptersManager::adapterDecreased, this, &BluetoothApplet::removeAdapter);
}

void BluetoothApplet::setAdapterPowered(bool powered)
{
    for (auto adapterItem : m_adapterItems) {
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

Device::State BluetoothApplet::initDeviceState()
{
    m_initDeviceState = Device::StateUnavailable;
    for (auto adapterItem : m_adapterItems) {
        if (adapterItem)
            if (Device::StateAvailable  == adapterItem->initDeviceState()) {
                m_initDeviceState = Device::StateAvailable;
                continue;
            }
        if (Device::StateConnected  == adapterItem->initDeviceState()) {
            m_initDeviceState = Device::StateConnected;
            break;
        }
    }
    return m_initDeviceState;
}

void BluetoothApplet::onPowerChanged(bool state)
{
    Q_UNUSED(state)
    bool powerState = false;
    for (auto adapterItem : m_adapterItems) {
         if (adapterItem->isPowered()) {
             powerState = true;
             break;
         }
    }
    emit powerChanged(powerState);
}

void BluetoothApplet::onDeviceStateChanged(const Device::State state)
{
    Q_UNUSED(state)

    Device::State deviceState = Device::StateUnavailable;
    for (auto adapterItem : m_adapterItems) {
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
}

void BluetoothApplet::addAdapter(Adapter *adapter)
{
    if (!adapter)
        return;

    if (!m_adapterItems.size()) {
        emit justHasAdapter();
    }

    auto adapterId = adapter->id();
    auto adatpterItem = new AdapterItem(m_adaptersManager, adapter, this);
    m_adapterItems[adapterId] = adatpterItem;
    m_centrealLayout->addWidget(adatpterItem);

    connect(adatpterItem, &AdapterItem::deviceStateChanged, this, &BluetoothApplet::onDeviceStateChanged);
    connect(adatpterItem, &AdapterItem::powerChanged, this, &BluetoothApplet::onPowerChanged);
    connect(adatpterItem, &AdapterItem::sizeChange, this, &BluetoothApplet::updateView);

    updateView();
}

void BluetoothApplet::removeAdapter(Adapter *adapter)
{
    if (adapter) {
        auto adapterId = adapter->id();
        auto adapterItem = m_adapterItems.value(adapterId);
        if (adapterItem) {
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
    for (auto adapterItem : m_adapterItems) {
        if (adapterItem) {
            contentHeight += adapterItem->viewHeight();
            if (adapterItem->isPowered())
                itemCount += adapterItem->deviceCount();
        }
    }

    auto adaptersCnt = m_adapterItems.size();
    if (adaptersCnt > 1) {
        m_line->setVisible(true);
        m_appletName->setVisible(true);
    } else {
        m_line->setVisible(false);
        m_appletName->setVisible(false);
    }

    if (adaptersCnt > 1)
        contentHeight += m_appletName->height();

    if (itemCount <= 10) {
        contentHeight += itemCount * ItemHeight;
        m_centralWidget->setFixedHeight(contentHeight);
        setFixedHeight(contentHeight);
        setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    } else {
        contentHeight += itemCount * ItemHeight;
        m_centralWidget->setFixedHeight(contentHeight);
        setFixedHeight(10 * ItemHeight);
        setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOn);
    }
}
