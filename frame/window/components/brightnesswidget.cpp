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
#include "brightnesswidget.h"
#include "brightnessmodel.h"
#include "brightnessmonitorwidget.h"
#include "imageutil.h"
#include "slidercontainer.h"

#include <QHBoxLayout>
#include <QDebug>

#define BACKSIZE 36
#define IMAGESIZE 18

BrightnessWidget::BrightnessWidget(BrightnessModel *model, QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_sliderContainer(new SliderContainer(this))
    , m_model(model)
{
    initUi();
    initConnection();
}

BrightnessWidget::~BrightnessWidget()
{
}

SliderContainer *BrightnessWidget::sliderContainer()
{
    return m_sliderContainer;
}

void BrightnessWidget::showEvent(QShowEvent *event)
{
    DBlurEffectWidget::showEvent(event);
    Q_EMIT visibleChanged(true);
}

void BrightnessWidget::hideEvent(QHideEvent *event)
{
    DBlurEffectWidget::hideEvent(event);
    Q_EMIT visibleChanged(true);
}

void BrightnessWidget::initUi()
{
    QHBoxLayout *layout = new QHBoxLayout(this);
    layout->setContentsMargins(15, 0, 12, 0);
    layout->addWidget(m_sliderContainer);

    QPixmap leftPixmap = ImageUtil::loadSvg(":/icons/resources/brightness.svg", QSize(IMAGESIZE, IMAGESIZE));
    QPixmap rightPixmap = ImageUtil::loadSvg(":/icons/resources/ICON_Device_Laptop.svg", QSize(IMAGESIZE, IMAGESIZE));
    m_sliderContainer->updateSlider(SliderContainer::IconPosition::LeftIcon, { leftPixmap.size(), QSize(), leftPixmap, 10 });
    m_sliderContainer->updateSlider(SliderContainer::IconPosition::RightIcon, { rightPixmap.size(), QSize(BACKSIZE, BACKSIZE), rightPixmap, 12});

    SliderProxyStyle *style = new SliderProxyStyle;
    style->setParent(m_sliderContainer->slider());
    m_sliderContainer->slider()->setStyle(style);
}

void BrightnessWidget::initConnection()
{
    connect(m_sliderContainer->slider(), &QSlider::valueChanged, this, [ this ](int value) {
        BrightMonitor *monitor = m_model->primaryMonitor();
        if (monitor)
            m_model->setBrightness(monitor, value);
    });

    connect(m_model, &BrightnessModel::brightnessChanged, this, &BrightnessWidget::onUpdateBright);

    BrightMonitor *monitor = m_model->primaryMonitor();
    if (monitor)
        onUpdateBright(monitor);
}

void BrightnessWidget::onUpdateBright(BrightMonitor *monitor)
{
    if (!monitor->isPrimary())
        return;
    // 此处只显示主屏的亮度
    m_sliderContainer->slider()->setValue(monitor->brihtness());
}
