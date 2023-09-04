// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "brightnesswidget.h"
#include "brightnessmodel.h"
#include "imageutil.h"
#include "slidercontainer.h"

#include <DGuiApplicationHelper>

#include <QHBoxLayout>
#include <QDebug>

#define BACKSIZE 36
#define IMAGESIZE 18

DGUI_USE_NAMESPACE

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
}

void BrightnessWidget::initUi()
{
    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(15, 0, 12, 0);

    onThemeTypeChanged();
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

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &BrightnessWidget::onThemeTypeChanged);
    connect(m_model, &BrightnessModel::monitorLightChanged, this, &BrightnessWidget::updateSliderValue);
    updateSliderValue();
}

void BrightnessWidget::updateSliderValue()
{
    BrightMonitor *monitor = m_model->primaryMonitor();
    if (monitor) {
        m_sliderContainer->updateSliderValue(monitor->brightness());
    }
}

void BrightnessWidget::convertThemePixmap(QPixmap &pixmap)
{
    // 图片是黑色的，如果当前主题为白色主题，则无需转换
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::ColorType::LightType)
        return;

    // 如果是黑色主题，则转换成白色图像
    QPainter painter(&pixmap);
    painter.setCompositionMode(QPainter::CompositionMode_SourceIn);
    painter.fillRect(pixmap.rect(), Qt::white);
    painter.end();
}

void BrightnessWidget::onThemeTypeChanged()
{
    QPixmap leftPixmap = ImageUtil::loadSvg(":/brightness.svg", QSize(IMAGESIZE, IMAGESIZE));
    QPixmap rightPixmap = ImageUtil::loadSvg(":/ICON_Device_Laptop.svg", QSize(IMAGESIZE, IMAGESIZE));
    convertThemePixmap(leftPixmap);
    convertThemePixmap(rightPixmap);
    m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon, leftPixmap, QSize(), 10);
    m_sliderContainer->setIcon(SliderContainer::IconPosition::RightIcon, rightPixmap, QSize(BACKSIZE, BACKSIZE), 12);
}
