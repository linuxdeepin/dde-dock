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
#ifndef BRIGHTNESSWIDGET_H
#define BRIGHTNESSWIDGET_H

#include <DBlurEffectWidget>

DWIDGET_USE_NAMESPACE

class SliderContainer;
class BrightnessModel;
class BrightMonitor;

class BrightnessWidget : public QWidget
{
    Q_OBJECT

public:
    explicit BrightnessWidget(BrightnessModel *model, QWidget *parent = nullptr);
    ~BrightnessWidget() override;

Q_SIGNALS:
    void brightClicked();

protected:
    void showEvent(QShowEvent *event) override;

private:
    void initUi();
    void initConnection();
    void updateSliderValue();

private:
    SliderContainer *m_sliderContainer;
    BrightnessModel *m_model;
};

#endif // LIGHTSETTINGWIDGET_H
