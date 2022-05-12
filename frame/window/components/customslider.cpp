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

#include <DStyle>
#include <DApplicationHelper>

#include <QMouseEvent>
#include <QDebug>
#include <QTimer>
#include <QGridLayout>
#include <QLabel>

DWIDGET_USE_NAMESPACE

CustomSlider::CustomSlider(CustomSlider::SliderType type, QWidget *parent)
    : DSlider(Qt::Horizontal, parent)
    , m_separateValue(0)
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

void CustomSlider::setSeparateValue(int value)
{
    m_separateValue = value;
}

void CustomSlider::wheelEvent(QWheelEvent *e)
{
    e->ignore();
}

void CustomSlider::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e)

    if (m_separateValue <= 0)
        return;

    QPainter pa(this);
    auto dpa = DApplicationHelper::instance()->palette(this);
    QPen penLine = QPen(dpa.color(DPalette::ObviousBackground), 2);

    //超过间隔线后需要更换间隔线颜色为活动色
    if (qtSlider()->value() >= m_separateValue ) {
        QPalette pe = this->palette();
        QColor brushColor(pe.color(QPalette::Highlight));
        penLine.setColor(brushColor);
    }

    int margin = DStyle::pixelMetric(style(), DStyle::PM_FocusBorderSpacing) + DStyle::pixelMetric(style(), DStyle::PM_FocusBorderSpacing);
    int offsetSize = style()->pixelMetric(QStyle::PM_SliderLength, nullptr, this) / 2;
    int width = this->qtSlider()->width();

    width -= 2 * offsetSize + margin * 2;
    Qt::Orientation orient = this->orientation();
    QSlider::TickPosition tick = m_tickPosition;
    QSlider* slider = DSlider::slider();
    qreal percentage = (m_separateValue - slider->minimum()) * 1.0 / (slider->maximum() - slider->minimum());

    pa.setPen(penLine);
    int leftIconwidth = 0;
    //获取左边声音图标宽度
    QGridLayout *gridLayout = dynamic_cast<QGridLayout *>(this->layout());
    if (!gridLayout)
        return;
    QLayoutItem* item = gridLayout->itemAtPosition(1, 0);
    if (item) {
        leftIconwidth = item->geometry().size().width();
    }

    qreal startX = offsetSize + margin + leftIconwidth + this->contentsMargins().left();
    qreal startY = slider->y() + 10;
    //分别绘制滑动条上方矩形和下方矩形，避免与滑动条重叠，
    //画笔宽为2个像素,设置绘制时矩形高设为3,可达到高度为5的效果
    if (orient == Qt::Horizontal) {
        qreal sliderX = percentage * width;
        if (slider->value() >= m_separateValue) {
            int num = (sliderX + 2) / 3;
            qAbs(3 * num + 1 - sliderX) - qAbs(sliderX - (3 * num - 2)) >= 0 ? num : num++;
            sliderX = 3 * num - 2;
        } else {
            //将分割线左、右的滚动条进行比较获取距离最相近的位置，绘制分割线
            sliderX = (slider->maximum() - m_separateValue) * 1.0 * width / slider->maximum();
            int num = (sliderX + 1) / 3;
            qAbs(3 * num - sliderX) - qAbs(sliderX - 3 * (num - 1)) >= 0 ? num : num++;
            sliderX = width - (3 * num - 1);
        }
        if (tick == QSlider::TicksAbove || tick == QSlider::TicksBelow || tick == QSlider::NoTicks) {
            pa.drawLine(QPointF(startX + sliderX, startY), QPointF(startX + sliderX, startY + 3));
            pa.drawLine(QPointF(startX + sliderX, startY + 9), QPointF(startX + sliderX, startY + 12));
        }
    }
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
