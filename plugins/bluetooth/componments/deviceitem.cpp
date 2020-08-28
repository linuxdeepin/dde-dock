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

#include "deviceitem.h"
#include "constants.h"
#include "bluetoothconstants.h"
#include "util/imageutil.h"
#include "util/statebutton.h"

#include <DApplicationHelper>
#include <DStyle>

#include <QHBoxLayout>
#include <QPainter>

DGUI_USE_NAMESPACE

extern void initFontColor(QWidget *widget);

DeviceItem::DeviceItem(Device *d, QWidget *parent)
    : QWidget(parent)
    , m_title(new QLabel(this))
    , m_state(new StateButton(this))
    , m_loadingStat(new DSpinner)
    , m_line(new HorizontalSeparator(this))
    , m_typeIcon(new QLabel(this))
{
    m_device = d;

    setFixedHeight(ITEMHEIGHT);
    auto themeChanged = [&](DApplicationHelper::ColorType themeType){
        QString iconPrefix;
        QString iconSuffix;
        switch (themeType) {
            case DApplicationHelper::UnknownType:
            case DApplicationHelper::LightType: {
                iconPrefix = LIGHTICON;
                iconSuffix = LIGHTICONSUFFIX;
            }
                break;
            case DApplicationHelper::DarkType: {
                iconPrefix = DARKICON;
                iconSuffix = DARKICONSUFFIX;
            }
        }

        QString iconString;
        if (!m_device->deviceType().isEmpty())
            iconString = iconPrefix + m_device->deviceType() + iconSuffix;
        else
            iconString = iconPrefix + QString("other") + iconSuffix;
        m_typeIcon->setPixmap(ImageUtil::loadSvg(iconString, ":/", PLUGIN_ICON_MIN_SIZE, devicePixelRatioF()));
    };

    themeChanged(DApplicationHelper::instance()->themeType());

    m_state->setType(StateButton::Check);
    m_state->setFixedSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE);
    m_state->setVisible(false);

    m_title->setText(nameDecorated(m_device->name()));
    initFontColor(m_title);

    m_line->setVisible(true);
    m_loadingStat->setFixedSize(20, 20);
    m_loadingStat->setVisible(false);

    auto deviceLayout = new QVBoxLayout;
    deviceLayout->setMargin(0);
    deviceLayout->setSpacing(0);
    deviceLayout->addWidget(m_line);
    auto itemLayout = new QHBoxLayout;
    itemLayout->setMargin(0);
    itemLayout->setSpacing(0);
    itemLayout->addSpacing(MARGIN);
    itemLayout->addWidget(m_typeIcon);
    itemLayout->addSpacing(5);
    itemLayout->addWidget(m_title);
    itemLayout->addStretch();
    itemLayout->addWidget(m_state);
    itemLayout->addWidget(m_loadingStat);
    itemLayout->addSpacing(MARGIN);
    deviceLayout->addLayout(itemLayout);
    setLayout(deviceLayout);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, themeChanged);

    changeState(m_device->state());
}

bool DeviceItem::operator <(const DeviceItem &item)
{
    return  this->device()->deviceTime() < item.device()->deviceTime();
}

void DeviceItem::setTitle(const QString &name)
{
    m_title->setText(nameDecorated(name));
}

void DeviceItem::mousePressEvent(QMouseEvent *event)
{
    m_device->updateDeviceTime();
    emit clicked(m_device);
    QWidget::mousePressEvent(event);
}

void DeviceItem::changeState(const Device::State state)
{
    switch (state) {
        case Device::StateUnavailable: {
            m_state->setVisible(false);
            m_loadingStat->stop();
            m_loadingStat->hide();
            m_loadingStat->setVisible(false);
        }
            break;
        case Device::StateAvailable: {
            m_state->setVisible(false);
            m_loadingStat->start();
            m_loadingStat->show();
            m_loadingStat->setVisible(true);
        }
            break;
        case Device::StateConnected: {
            m_loadingStat->stop();
            m_loadingStat->hide();
            m_loadingStat->setVisible(false);
            m_state->setVisible(true);
        }
            break;
    }
}

QString DeviceItem::nameDecorated(const QString &name)
{
    return QFontMetrics(m_title->font()).elidedText(name, Qt::ElideRight, POPUPWIDTH - MARGIN * 2 - PLUGIN_ICON_MIN_SIZE - 30);
}

HorizontalSeparator::HorizontalSeparator(QWidget *parent)
    : QWidget(parent)
{
    setFixedHeight(1);
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
}

void HorizontalSeparator::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.fillRect(rect(), QColor(0, 0, 0, 0));
}

MenueItem::MenueItem(QWidget *parent)
    : QLabel(parent)
{
}

void MenueItem::mousePressEvent(QMouseEvent *event)
{
    QLabel::mousePressEvent(event);
    emit clicked();
}
