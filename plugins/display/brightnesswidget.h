// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
    void convertThemePixmap(QPixmap &pixmap);

private Q_SLOTS:
    void updateSliderValue();
    void onThemeTypeChanged();

private:
    SliderContainer *m_sliderContainer;
    BrightnessModel *m_model;
};

#endif // LIGHTSETTINGWIDGET_H
