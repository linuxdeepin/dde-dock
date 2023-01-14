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
