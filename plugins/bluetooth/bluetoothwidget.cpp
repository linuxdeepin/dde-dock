/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "bluetoothwidget.h"
#include "adaptersmanager.h"
#include "bloothadapterwidget.h"
#include "adapter.h"
#include "device.h"

#include <DSwitchButton>
#include <DListView>

#include <QVBoxLayout>
#include <QLabel>

BluetoothWidget::BluetoothWidget(AdaptersManager *adapterManager, QWidget *parent)
    : QWidget(parent)
    , m_switchButton(new DSwitchButton(this))
    , m_headerWidget(new QWidget(this))
    , m_adapterWidget(new QWidget(this))
    , m_adaptersManager(adapterManager)
    , m_adapterLayout(new QVBoxLayout(m_adapterWidget))
{
    initUi();
    initConnection();
}

BluetoothWidget::~BluetoothWidget()
{
}

void BluetoothWidget::onAdapterIncreased(Adapter *adapter)
{
    BloothAdapterWidget *adapterWidget = new BloothAdapterWidget(adapter, m_adapterWidget);
    m_adapterLayout->addWidget(adapterWidget);
    connect(adapterWidget, &BloothAdapterWidget::requestConnectDevice, this, [ this, adapter ](Device *device) {
        m_adaptersManager->connectDevice(device, adapter);
    });
    connect(adapterWidget, &BloothAdapterWidget::requestUpdate, this, [ this ] {
        adjustHeight();
    });

    updateCheckStatus();

    QMetaObject::invokeMethod(this, &BluetoothWidget::adjustHeight, Qt::QueuedConnection);
}

void BluetoothWidget::onAdapterDecreased(Adapter *adapter)
{
    for (int i = 0; i < m_adapterLayout->count(); i++) {
        BloothAdapterWidget *adapterWidget = static_cast<BloothAdapterWidget *>(m_adapterLayout->itemAt(i)->widget());
        if (adapterWidget && adapterWidget->adapter() == adapter) {
            m_adapterLayout->removeWidget(adapterWidget);

            updateCheckStatus();
            QMetaObject::invokeMethod(this, &BluetoothWidget::adjustHeight, Qt::QueuedConnection);
            break;
        }
    }
}

void BluetoothWidget::onCheckedChanged(bool checked)
{
    QList<const Adapter *> adapters = m_adaptersManager->adapters();
    for (const Adapter *adapter : adapters)
        m_adaptersManager->setAdapterPowered(adapter, checked);
}

void BluetoothWidget::initUi()
{
    QHBoxLayout *headerLayout = new QHBoxLayout(m_headerWidget);
    headerLayout->addStretch();
    headerLayout->addWidget(m_switchButton);
    headerLayout->addStretch();
    QVBoxLayout *mainLayout = new QVBoxLayout(this);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->setSpacing(0);
    mainLayout->addWidget(m_headerWidget);
    mainLayout->addSpacing(3);
    mainLayout->addWidget(m_adapterWidget);

    m_adapterLayout->setContentsMargins(0, 0, 0, 0);
    m_adapterLayout->setSpacing(0);

    QList<const Adapter *> adapters = m_adaptersManager->adapters();
    for (const Adapter *adapter : adapters) {
        onAdapterIncreased(const_cast<Adapter *>(adapter));
    }
}

void BluetoothWidget::initConnection()
{
    connect(m_adaptersManager, &AdaptersManager::adapterIncreased, this, &BluetoothWidget::onAdapterIncreased);
    connect(m_adaptersManager, &AdaptersManager::adapterDecreased, this, &BluetoothWidget::onAdapterDecreased);
    connect(m_switchButton, &DSwitchButton::checkedChanged, this, &BluetoothWidget::onCheckedChanged);
}

void BluetoothWidget::updateCheckStatus()
{
    bool checked = false;
    QList<const Adapter *> adapters = m_adaptersManager->adapters();
    for (const Adapter *adapter : adapters)
        checked = adapter->powered();

    m_switchButton->setChecked(checked);
}

void BluetoothWidget::adjustHeight()
{
    int height = m_switchButton->height() + m_headerWidget->height();
    for (int i = 0; i < m_adapterLayout->count(); i++) {
        BloothAdapterWidget *adapterWidget = static_cast<BloothAdapterWidget *>(m_adapterLayout->itemAt(i)->widget());
        if (!adapterWidget)
            continue;

        height += adapterWidget->height();
    }

    setFixedHeight(height);
}
