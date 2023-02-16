// Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
