/*
 * Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
 *
 * Author:     zhaoyingzhen <zhaoyingzhen@uniontech.com>
 *
 * Maintainer: zhaoyingzhen <zhaoyingzhen@uniontech.com>
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
#include "brightnessadjwidget.h"
#include "brightnessmodel.h"
#include "slidercontainer.h"
#include "imageutil.h"

#include <QVBoxLayout>

BrightnessAdjWidget::BrightnessAdjWidget(QWidget *parent)
    : QWidget(parent)
    , m_mainLayout(new QVBoxLayout(this))
    , m_brightnessModel(new BrightnessModel(this))
{
    m_mainLayout->setSpacing(5);
    loadBrightnessItem();
}

void BrightnessAdjWidget::loadBrightnessItem()
{
    QList<BrightMonitor *> monitors = m_brightnessModel->monitors();
    for (BrightMonitor *monitor : monitors) {
        SliderContainer *sliderContainer = new SliderContainer(this);
        if (monitors.count() > 1)
            sliderContainer->setTitle(monitor->name());

        QPixmap leftPixmap = ImageUtil::loadSvg(":/icons/resources/brightnesslow", QSize(20, 20));
        QPixmap rightPixmap = ImageUtil::loadSvg(":/icons/resources/brightnesshigh", QSize(20, 20));
        sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon,leftPixmap, QSize(), 12);
        sliderContainer->setIcon(SliderContainer::IconPosition::RightIcon, rightPixmap, QSize(), 12);

        sliderContainer->setFixedWidth(310);
        sliderContainer->setFixedHeight(monitors.count() > 1 ? 56 : 30);
        sliderContainer->updateSliderValue(monitor->brightness());

        SliderProxyStyle *proxy = new SliderProxyStyle(SliderProxyStyle::Normal);
        sliderContainer->setSliderProxyStyle(proxy);
        m_mainLayout->addWidget(sliderContainer);

        connect(monitor, &BrightMonitor::brightnessChanged, sliderContainer, &SliderContainer::updateSliderValue);
        connect(sliderContainer, &SliderContainer::sliderValueChanged, monitor, &BrightMonitor::setBrightness);
    }
}

