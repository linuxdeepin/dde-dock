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
#ifndef BRIGHTNESS_ADJUSTMENT_WIDGET_H
#define BRIGHTNESS_ADJUSTMENT_WIDGET_H

#include <QWidget>

class QVBoxLayout;
class BrightnessModel;

/*!
 * \brief The BrightnessAdjWidget class
 * 显示器亮度调整页面
 */
class BrightnessAdjWidget : public QWidget
{
    Q_OBJECT
public:
    explicit BrightnessAdjWidget(QWidget *parent = nullptr);

private:
    void loadBrightnessItem();

private:
    QVBoxLayout *m_mainLayout;
    BrightnessModel *m_brightnessModel;
};


#endif // BRIGHTNESS_ADJUSTMENT_WIDGET_H
