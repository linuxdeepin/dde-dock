/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "customslider.h"
#include <QPainterPath>

#include <DStyle>
#include <DApplicationHelper>
#include <DGuiApplicationHelper>

#include <QMouseEvent>
#include <QDebug>
#include <QTimer>
#include <QGridLayout>
#include <QLabel>

DWIDGET_USE_NAMESPACE

CustomSlider::CustomSlider(CustomSlider::SliderType type, QWidget *parent)
    : DSlider(Qt::Horizontal, parent)
{
    setType(type);
    DSlider::slider()->setTracking(false);
}

CustomSlider::CustomSlider(Qt::Orientation orientation, QWidget *parent)
    : DSlider(orientation, parent)
{
    DSlider::slider()->setTracking(false);
}

void CustomSlider::setType(CustomSlider::SliderType type)
{
    switch (type) {
    case Vernier: setProperty("handleType", "Vernier"); break;
    case Progress: setProperty("handleType", "None"); break;
    default: setProperty("handleType", "Normal"); break;
    }
}

QSlider *CustomSlider::qtSlider()
{
    return DSlider::slider();
}

void CustomSlider::setRange(int min, int max)
{
    setMinimum(min);
    setMaximum(max);
}

void CustomSlider::setTickPosition(QSlider::TickPosition tick)
{
    m_tickPosition = tick;
}

void CustomSlider::setTickInterval(int ti)
{
    DSlider::slider()->setTickInterval(ti);
}

void CustomSlider::setSliderPosition(int Position)
{
    DSlider::slider()->setSliderPosition(Position);
}

void CustomSlider::setAnnotations(const QStringList &annotations)
{
    switch (m_tickPosition) {
    case QSlider::TicksLeft:
        setLeftTicks(annotations);
        break;
    case QSlider::TicksRight:
        setRightTicks(annotations);
        break;
    default:
        break;
    }
}

void CustomSlider::setOrientation(Qt::Orientation orientation)
{
    Q_UNUSED(orientation)
}

void CustomSlider::wheelEvent(QWheelEvent *e)
{
    e->ignore();
}

SliderContainer::SliderContainer(CustomSlider::SliderType type, QWidget *parent)
    : QWidget (parent)
    , m_slider(new CustomSlider(type, this))
    , m_titleLabel(new QLabel(this))
{
    QVBoxLayout *mainlayout = new QVBoxLayout(this);
    mainlayout->setContentsMargins(0, 0, 0, 0);
    mainlayout->setSpacing(5);
    mainlayout->addWidget(m_titleLabel);
    mainlayout->addWidget(m_slider);
}

SliderContainer::SliderContainer(Qt::Orientation orientation, QWidget *parent)
    : QWidget(parent)
    , m_slider(new CustomSlider(orientation, this))
    , m_titleLabel(new QLabel(this))
{
    QVBoxLayout *mainlayout = new QVBoxLayout(this);
    mainlayout->setContentsMargins(0, 1, 0, 0);
    mainlayout->setSpacing(1);

    m_titleLabel->setFixedHeight(8);
    mainlayout->addWidget(m_titleLabel);
    mainlayout->addWidget(m_slider);
}

SliderContainer::~SliderContainer()
{
}

void SliderContainer::setTitle(const QString &title)
{
    m_titleLabel->setText(title);
}

CustomSlider *SliderContainer::slider()
{
    return m_slider;
}

SliderProxy::SliderProxy(QStyle *style)
    : QProxyStyle(style)
{
}

SliderProxy::~SliderProxy()
{
}

void SliderProxy::drawComplexControl(QStyle::ComplexControl control, const QStyleOptionComplex *option, QPainter *painter, const QWidget *widget) const
{
    if (control != ComplexControl::CC_Slider)
        return;

    // 绘制之前先保存之前的画笔
    painter->save();
    painter->setRenderHint(QPainter::RenderHint::Antialiasing);
    // 获取滑动条和滑块的区域
    const QStyleOptionSlider *sliderOption = static_cast<const QStyleOptionSlider *>(option);
    QRect rectGroove = subControlRect(CC_Slider, sliderOption, SC_SliderGroove, widget);
    QRect rectHandle = subControlRect(CC_Slider, sliderOption, SC_SliderHandle, widget);
    // 深色背景下，滑块和滑动条白色，浅色背景下，滑块和滑动条黑色
    QBrush brush(DGuiApplicationHelper::DarkType == DGuiApplicationHelper::instance()->themeType() ? Qt::white : Qt::black);
    // 此处中绘制圆形滑动条，需要绘制圆角，圆角大小为其高度的一半
    QPainterPath pathGroove;
    int radius = rectGroove.height() / 2;
    pathGroove.addRoundedRect(rectGroove, radius, radius);
    painter->fillPath(pathGroove, brush);

    // 绘制滑块,因为滑块是正圆形，而它本来的区域是一个长方形区域，因此，需要计算当前
    // 区域的正中心区域，将其作为一个正方形区域来绘制圆形滑块
    int handleSize = qMin(rectHandle.width(), rectHandle.height());
    int x = rectHandle.x() + (rectHandle.width() - handleSize) / 2;
    int y = rectHandle.y() + (rectHandle.height() - handleSize) / 2;
    rectHandle.setX(x);
    rectHandle.setY(y);
    rectHandle.setWidth(handleSize);
    rectHandle.setHeight(handleSize);

    QPainterPath pathHandle;
    pathHandle.addEllipse(rectHandle);
    painter->fillPath(pathHandle, brush);
    painter->restore();
}
