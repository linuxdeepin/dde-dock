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

#ifndef SLIDERCONTAINER_H
#define SLIDERCONTAINER_H

#include <DSlider>

#include <QProxyStyle>
#include <QTimer>

class QLabel;
class SliderProxyStyle;
class SliderIconWidget;

/**
 * @brief 滚动条的类的封装，封装这个类的原因，是为了方便设置左右图标，dtk中的对应的DSlider类也有这个功能，
 * 但是只能简单的设置左右图标，对于右边图标有阴影的，需要外部提供一个带阴影图标，但是如果由外部来提供，
 * 通过QPixmap绘制的带阴影的图标无法消除锯齿（即使通过反走样也不行），因此在此处封装这个类
 */
class SliderContainer : public QWidget
{
    Q_OBJECT

public:
    enum IconPosition {
        LeftIcon = 0,
        RightIcon
    };

public:
    explicit SliderContainer(QWidget *parent);
    ~SliderContainer() override;

    void setTitle(const QString &text);
    void setSliderProxyStyle(QProxyStyle *proxyStyle);
    void setIcon(const IconPosition &iconPosition, const QIcon &icon);
    void setIcon(const IconPosition &iconPosition, const QPixmap &icon, const QSize &shadowSize, int space);

    void setPageStep(int step);
    void setRange(int min, int max);

Q_SIGNALS:
    void iconClicked(const IconPosition &);
    void sliderValueChanged(int value);

public slots:
    void updateSliderValue(int value);

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;
    QSize getSuitableSize(const QSize &iconSize, const QSize &bgSize);

private:
    SliderIconWidget *m_leftIconWidget;
    QSlider *m_slider;
    QLabel *m_titleLabel;
    SliderIconWidget *m_rightIconWidget;
    QWidget *m_spaceLeftWidget;
    QWidget *m_spaceRightWidget;
};

/**
 * @brief 用来设置滚动条的样式
 * @param drawSpecial: true
 */
class SliderProxyStyle : public QProxyStyle
{
    Q_OBJECT

public:
    enum StyleType {
        RoundHandler = 0,    // 绘制那种黑色圆底滑动条
        Normal               // 绘制那种通用的滑动条
    };

public:
    explicit SliderProxyStyle(StyleType drawSpecial = RoundHandler, QStyle *style = nullptr);
    ~SliderProxyStyle() override;

protected:
    void drawComplexControl(QStyle::ComplexControl control, const QStyleOptionComplex *option, QPainter *painter, const QWidget *widget = nullptr) const override;

private:
    void drawNormalSlider(QPainter *painter, QRect rectGroove, QRect rectHandle, const QWidget *wigdet) const;
    void drawRoundSlider(QPainter *painter, QRect rectGroove, QRect rectHandle, const QWidget *wigdet) const;

private:
    StyleType m_drawSpecial;
};

#endif // VOLUMESLIDER_H
