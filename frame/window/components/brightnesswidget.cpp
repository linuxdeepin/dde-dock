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
#include "customslider.h"
#include "brightnessmodel.h"
#include "brightnessmonitorwidget.h"
#include "imageutil.h"

#include <QHBoxLayout>
#include <QDebug>

#define BACKSIZE 36
#define IMAGESIZE 24

BrightnessWidget::BrightnessWidget(QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_slider(new CustomSlider(CustomSlider::SliderType::Normal, this))
    , m_model(new BrightnessModel(this))
{
    initUi();
    initConenction();
}

BrightnessWidget::~BrightnessWidget()
{
}

BrightnessModel *BrightnessWidget::model()
{
    return m_model;
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
    layout->setContentsMargins(20, 0, 20, 0);
    layout->addWidget(m_slider);

    m_slider->setPageStep(1);
    m_slider->setIconSize(QSize(BACKSIZE, BACKSIZE));

    QIcon leftIcon(QPixmap(":/icons/resources/brightness.svg").scaled(IMAGESIZE, IMAGESIZE));
    m_slider->setLeftIcon(leftIcon);
    QPixmap rightPixmap = ImageUtil::getShadowPixmap(QPixmap(QString(":/icons/resources/ICON_Device_Laptop.svg")).scaled(24, 24), Qt::lightGray, QSize(36, 36));
    m_slider->setRightIcon(rightPixmap);

    SliderProxy *style = new SliderProxy;
    style->setParent(m_slider->qtSlider());
    m_slider->qtSlider()->setStyle(style);
}

void BrightnessWidget::initConenction()
{
    connect(m_slider, &CustomSlider::iconClicked, this, [ this ](DSlider::SliderIcons icon, bool) {
        if (icon == DSlider::SliderIcons::RightIcon)
            Q_EMIT rightIconClicked();
    });

    connect(m_slider, &CustomSlider::valueChanged, this, [ this ](int value) {
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
    m_slider->setValue(monitor->brihtness());
}
