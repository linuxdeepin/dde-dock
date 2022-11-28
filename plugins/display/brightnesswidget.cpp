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
#include "imageutil.h"
#include "slidercontainer.h"

#include <QHBoxLayout>
#include <QDebug>

#define BACKSIZE 36
#define IMAGESIZE 18

BrightnessWidget::BrightnessWidget(BrightnessModel *model, QWidget *parent)
    : QWidget(parent)
    , m_sliderContainer(new SliderContainer(this))
    , m_model(model)
{
    initUi();
    initConnection();
}

BrightnessWidget::~BrightnessWidget()
{
}

void BrightnessWidget::showEvent(QShowEvent *event)
{
    QWidget::showEvent(event);

    // 显示的时候更新一下slider的主屏幕亮度值
    updateSliderValue();
    Q_EMIT visibleChanged(true);
}

void BrightnessWidget::hideEvent(QHideEvent *event)
{
    QWidget::hideEvent(event);

    Q_EMIT visibleChanged(true);
}

void BrightnessWidget::initUi()
{
    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(15, 0, 12, 0);

    QPixmap leftPixmap = ImageUtil::loadSvg(":/brightness.svg", QSize(IMAGESIZE, IMAGESIZE));
    QPixmap rightPixmap = ImageUtil::loadSvg(":/ICON_Device_Laptop.svg", QSize(IMAGESIZE, IMAGESIZE));
    m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon, leftPixmap, QSize(), 10);
    m_sliderContainer->setIcon(SliderContainer::IconPosition::RightIcon, rightPixmap, QSize(BACKSIZE, BACKSIZE), 12);

    // 需求要求调节范围是10%-100%,且调节幅度为1%
    m_sliderContainer->setRange(10, 100);
    m_sliderContainer->setPageStep(1);

    SliderProxyStyle *style = new SliderProxyStyle;
    m_sliderContainer->setSliderProxyStyle(style);

    mainLayout->addWidget(m_sliderContainer);
}

void BrightnessWidget::initConnection()
{
    connect(m_sliderContainer, &SliderContainer::sliderValueChanged, this, [ this ](int value) {
        BrightMonitor *monitor = m_model->primaryMonitor();
        if (monitor)
            monitor->setBrightness(value);
    });

    connect(m_sliderContainer, &SliderContainer::iconClicked, this, [ this ](const SliderContainer::IconPosition &position) {
        if (position == SliderContainer::IconPosition::RightIcon)
            Q_EMIT brightClicked();
    });

    updateSliderValue();
}

void BrightnessWidget::updateSliderValue()
{
    BrightMonitor *monitor = m_model->primaryMonitor();
    if (monitor) {
        m_sliderContainer->updateSliderValue(monitor->brightness());
    }
}
