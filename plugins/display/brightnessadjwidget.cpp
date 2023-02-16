// Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "brightnessadjwidget.h"
#include "brightnessmodel.h"
#include "slidercontainer.h"
#include "imageutil.h"

#include <QVBoxLayout>

const int ItemSpacing = 5;

BrightnessAdjWidget::BrightnessAdjWidget(QWidget *parent)
    : QWidget(parent)
    , m_mainLayout(new QVBoxLayout(this))
    , m_brightnessModel(new BrightnessModel(this))
{
    m_mainLayout->setMargin(0);
    m_mainLayout->setSpacing(ItemSpacing);

    loadBrightnessItem();
}

void BrightnessAdjWidget::loadBrightnessItem()
{
    QList<BrightMonitor *> monitors = m_brightnessModel->monitors();
    int itemHeight = monitors.count() > 1 ? 56 : 30;

    for (BrightMonitor *monitor : monitors) {
        SliderContainer *sliderContainer = new SliderContainer(this);
        if (monitors.count() > 1)
            sliderContainer->setTitle(monitor->name());

        QPixmap leftPixmap = ImageUtil::loadSvg(":/icons/resources/brightnesslow", QSize(20, 20));
        QPixmap rightPixmap = ImageUtil::loadSvg(":/icons/resources/brightnesshigh", QSize(20, 20));
        sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon,leftPixmap, QSize(), 12);
        sliderContainer->setIcon(SliderContainer::IconPosition::RightIcon, rightPixmap, QSize(), 12);
        // 需求要求调节范围是10%-100%,且调节幅度为1%
        sliderContainer->setRange(10, 100);
        sliderContainer->setPageStep(1);
        sliderContainer->setFixedWidth(310);
        sliderContainer->setFixedHeight(itemHeight);
        sliderContainer->updateSliderValue(monitor->brightness());

        SliderProxyStyle *proxy = new SliderProxyStyle(SliderProxyStyle::Normal);
        sliderContainer->setSliderProxyStyle(proxy);
        m_mainLayout->addWidget(sliderContainer);

        connect(monitor, &BrightMonitor::brightnessChanged, sliderContainer, &SliderContainer::updateSliderValue);
        connect(sliderContainer, &SliderContainer::sliderValueChanged, monitor, &BrightMonitor::setBrightness);
    }

    QMargins margins = this->contentsMargins();
    setFixedHeight(margins.top() + margins.bottom() + monitors.count() * itemHeight + monitors.count() * ItemSpacing);
}

