/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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

#include "brightnessmodel.h"
#include "slidercontainer.h"

#include <DStyle>
#include <DApplicationHelper>
#include <DGuiApplicationHelper>
#include <DPaletteHelper>

#include <QPainterPath>
#include <QMouseEvent>
#include <QDebug>
#include <QTimer>
#include <QGridLayout>
#include <QLabel>

DWIDGET_USE_NAMESPACE

// 用于绘制图标的窗体，此窗体不想让其在外部调用，因此，将其作为一个私有类
class SliderIconWidget : public QWidget
{
public:
    explicit SliderIconWidget(QWidget *parent)
        : QWidget(parent)
        , m_iconSize(QSize(24, 24))
        , m_shadowSize(QSize())
    {}

    void updateData(const QIcon &icon, const QSize &iconSize, const QSize &shadowSize)
    {
        m_icon = icon;
        m_iconSize = iconSize;
        m_shadowSize = shadowSize;
        update();
    }

    void updateIcon(const QIcon &icon)
    {
        m_icon = icon;
        update();
    }

protected:
    void paintEvent(QPaintEvent *e) override;

private:
    QIcon m_icon;
    QSize m_iconSize;
    QSize m_shadowSize;
};

void SliderIconWidget::paintEvent(QPaintEvent *e)
{
    if (m_iconSize.isNull() || m_icon.isNull())
        return QWidget::paintEvent(e);

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);
    if (m_shadowSize.isValid() && !m_shadowSize.isNull() && !m_shadowSize.isEmpty()) {
        // 绘制圆形背景
        painter.setPen(Qt::NoPen);
        // 获取阴影部分背景颜色
        DPalette dpa = DPaletteHelper::instance()->palette(this);
        painter.setBrush(dpa.brush(DPalette::ColorRole::Midlight));
        int x = (rect().width() - m_shadowSize.width() ) / 2;
        int y = (rect().height() - m_shadowSize.height() ) / 2;
        painter.drawEllipse(QRect(x, y, m_shadowSize.width(), m_shadowSize.height()));
    }
    // 绘制图标
    QPixmap pixmap = m_icon.pixmap(m_iconSize);
    int iconWidth = static_cast<int>(m_iconSize.width() / qApp->devicePixelRatio());
    int iconHeight = static_cast<int>(m_iconSize.height() / qApp->devicePixelRatio());
    int x = (rect().width() - iconWidth) / 2;
    int y = (rect().height() - iconHeight) / 2;
    painter.drawPixmap(x, y, iconWidth, iconHeight, pixmap);
}

SliderContainer::SliderContainer(QWidget *parent)
    : QWidget(parent)
    , m_leftIconWidget(new SliderIconWidget(this))
    , m_slider(new QSlider(Qt::Orientation::Horizontal, this))
    , m_titleLabel(new QLabel(this))
    , m_rightIconWidget(new SliderIconWidget(this))
    , m_spaceLeftWidget(new QWidget(this))
    , m_spaceRightWidget(new QWidget(this))
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->setSpacing(0);

    QHBoxLayout *sliderLayout = new QHBoxLayout(this);
    sliderLayout->setContentsMargins(0, 0, 0, 0);
    sliderLayout->setSpacing(0);
    sliderLayout->addWidget(m_leftIconWidget);
    sliderLayout->addWidget(m_spaceLeftWidget);
    sliderLayout->addWidget(m_slider);
    sliderLayout->addWidget(m_spaceRightWidget);
    sliderLayout->addWidget(m_rightIconWidget);

    mainLayout->addWidget(m_titleLabel);
    mainLayout->addLayout(sliderLayout);

    m_titleLabel->setVisible(false);

    m_leftIconWidget->installEventFilter(this);
    m_slider->installEventFilter(this);
    m_rightIconWidget->installEventFilter(this);
}

SliderContainer::~SliderContainer()
{
}

void SliderContainer::setTitle(const QString &text)
{
    m_titleLabel->setText(text);
    m_titleLabel->setVisible(!text.isEmpty());
}

QSize SliderContainer::getSuitableSize(const QSize &iconSize, const QSize &bgSize)
{
    if (bgSize.isValid() && !bgSize.isNull() && !bgSize.isEmpty())
        return bgSize;

    return iconSize;
}

void SliderContainer::updateSlider(const SliderContainer::IconPosition &iconPosition, const SliderContainer::SliderData &sliderData)
{
    switch (iconPosition) {
    case IconPosition::LeftIcon: {
        m_leftIconWidget->setFixedSize(getSuitableSize(sliderData.iconSize, sliderData.shadowSize));
        m_leftIconWidget->updateData(sliderData.icon, sliderData.iconSize, sliderData.shadowSize);
        m_spaceLeftWidget->setFixedWidth(sliderData.space);
        break;
    }
    case IconPosition::RightIcon: {
        m_rightIconWidget->setFixedSize(getSuitableSize(sliderData.iconSize, sliderData.shadowSize));
        m_rightIconWidget->updateData(sliderData.icon, sliderData.iconSize, sliderData.shadowSize);
        m_spaceRightWidget->setFixedWidth(sliderData.space);
        break;
    }
    }
}

void SliderContainer::setIcon(const SliderContainer::IconPosition &iconPosition, const QIcon &icon)
{
    switch (iconPosition) {
    case IconPosition::LeftIcon: {
        m_leftIconWidget->updateIcon(icon);
        break;
    }
    case IconPosition::RightIcon: {
        m_rightIconWidget->updateIcon(icon);
        break;
    }
    }
}

QSlider *SliderContainer::slider()
{
    return m_slider;
}

bool SliderContainer::eventFilter(QObject *watched, QEvent *event)
{
    if (event->type() == QEvent::MouseButtonRelease) {
        if (watched == m_leftIconWidget)
            Q_EMIT iconClicked(IconPosition::LeftIcon);
        else if (watched == m_rightIconWidget)
            Q_EMIT iconClicked(IconPosition::RightIcon);
    }

    return QWidget::eventFilter(watched, event);
}

SliderProxyStyle::SliderProxyStyle(StyleType drawSpecial, QStyle *style)
    : QProxyStyle(style)
    , m_drawSpecial(drawSpecial)
{
}

SliderProxyStyle::~SliderProxyStyle()
{
}

void SliderProxyStyle::drawComplexControl(QStyle::ComplexControl control, const QStyleOptionComplex *option, QPainter *painter, const QWidget *widget) const
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
    if (m_drawSpecial == RoundHandler)
        drawRoundSlider(painter, rectGroove, rectHandle);
    else
        drawNormalSlider(painter, rectGroove, rectHandle, const_cast<QWidget *>(widget));

    painter->restore();
}

// 绘制通用的滑动条
void SliderProxyStyle::drawNormalSlider(QPainter *painter, QRect rectGroove, QRect rectHandle, QWidget *wigdet) const
{
    DPalette dpa = DPaletteHelper::instance()->palette(wigdet);
    QPen penLine = QPen(dpa.color(DPalette::Highlight), 2);
    // 绘制上下的竖线，一根竖线的宽度是2，+4个像素刚好保证中间也是间隔2个像素
    for (int i = rectGroove.x(); i < rectGroove.x() + rectGroove.width(); i = i + 4) {
        if (i < rectHandle.x())
            painter->setPen(penLine);
        else
            painter->setPen(QPen(Qt::gray, 2));

        painter->drawLine(i, rectGroove.y() + 2, i, rectGroove.y() + rectGroove.height() - 2);
    }
    // 绘制滚动区域
    painter->setBrush(dpa.color(DPalette::Highlight));
    painter->setPen(Qt::NoPen);
    QPainterPath path;
    path.addRoundedRect(rectHandle, 6, 6);
    painter->drawPath(path);
}

// 绘制设计师定义的那种圆形滑块，黑色的滑条
void SliderProxyStyle::drawRoundSlider(QPainter *painter, QRect rectGroove, QRect rectHandle) const
{
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
}
