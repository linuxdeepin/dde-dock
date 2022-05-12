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

#ifndef CUSTOMCTRL_H
#define CUSTOMCTRL_H

#include <DSlider>
#include <QTimer>

class QLabel;

class CustomSlider : public DTK_WIDGET_NAMESPACE::DSlider
{
    Q_OBJECT

public:
    enum SliderType {
        Normal,
        Vernier,
        Progress
    };

public:
    explicit CustomSlider(SliderType type = Normal, QWidget *parent = nullptr);
    explicit CustomSlider(Qt::Orientation orientation, QWidget *parent = nullptr);

    inline CustomSlider *slider() const { return const_cast<CustomSlider *>(this); }
    QSlider *qtSlider();

    void setType(SliderType type);
    void setRange(int min, int max);
    void setTickPosition(QSlider::TickPosition tick);
    void setTickInterval(int ti);
    void setSliderPosition(int Position);
    void setAnnotations(const QStringList &annotations);
    void setOrientation(Qt::Orientation orientation);
    //当value大于0时，在slider中插入一条分隔线
    void setSeparateValue(int value = 0);

protected:
    void wheelEvent(QWheelEvent *e);
    void paintEvent(QPaintEvent *e);
private:
    QSlider::TickPosition m_tickPosition = QSlider::TicksBelow;
    int m_separateValue;
};

class SliderContainer : public QWidget
{
    Q_OBJECT

public:
    explicit SliderContainer(CustomSlider::SliderType type = CustomSlider::Normal, QWidget *parent = nullptr);
    explicit SliderContainer(Qt::Orientation orientation, QWidget *parent);
    ~SliderContainer();
    void setTitle(const QString &title);
    CustomSlider *slider();

private:
    CustomSlider *m_slider;
    QLabel *m_titleLabel;
};

#endif // VOLUMESLIDER_H
