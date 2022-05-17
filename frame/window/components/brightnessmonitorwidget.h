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
#ifndef BRIGHTNESSMONITORWIDGET_H
#define BRIGHTNESSMONITORWIDGET_H

#include <QMap>
#include <QWidget>

class QLabel;
class CustomSlider;
class BrightnessModel;
class QStandardItemModel;
class QVBoxLayout;
class SliderContainer;
class BrightMonitor;
class SettingDelegate;

namespace Dtk { namespace Widget { class DListView; } }

using namespace Dtk::Widget;

class BrightnessMonitorWidget : public QWidget
{
    Q_OBJECT

public:
    explicit BrightnessMonitorWidget(BrightnessModel *model, QWidget *parent = nullptr);
    ~BrightnessMonitorWidget() override;

private:
    void initUi();
    void initConnection();
    void reloadMonitor();

    void resetHeight();

private Q_SLOTS:
    void onBrightChanged(BrightMonitor *monitor);

private:
    QWidget *m_sliderWidget;
    QVBoxLayout *m_sliderLayout;
    QList<QPair<BrightMonitor *, SliderContainer *>> m_sliderContainers;
    QLabel *m_descriptionLabel;
    DListView *m_deviceList;
    BrightnessModel *m_brightModel;
    QStandardItemModel *m_model;
    SettingDelegate *m_delegate;
};

#endif // BRIGHTNESSMONITORWIDGET_H
