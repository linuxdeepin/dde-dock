// Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DISPLAY_SETTING_WIDGET_H
#define DISPLAY_SETTING_WIDGET_H

#include <QWidget>

class QPushButton;
class BrightnessAdjWidget;
class DevCollaborationWidget;

/*!
 * \brief The DisplaySettingWidget class
 * 显示设置页面，快捷设置面板-->亮度调节栏右边显示按钮-->此页面
 */
class DisplaySettingWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DisplaySettingWidget(QWidget *parent = nullptr);

Q_SIGNALS:
    void requestHide();

private:
    void initUI();
    void resizeWidgetHeight();

private:
    BrightnessAdjWidget *m_brightnessAdjWidget;     // 亮度调整
    DevCollaborationWidget *m_collaborationWidget;  // 跨端协同
    QPushButton *m_settingBtn;                      // 设置按钮
};


#endif // DISPLAY_SETTING_WIDGET_H
