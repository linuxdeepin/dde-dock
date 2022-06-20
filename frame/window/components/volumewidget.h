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
#ifndef VOLUMEWIDGET_H
#define VOLUMEWIDGET_H

#include <DBlurEffectWidget>
#include <QWidget>

class VolumeModel;
class QDBusMessage;
class SliderContainer;
class QLabel;
class AudioSink;

DWIDGET_USE_NAMESPACE

class VolumeWidget : public DBlurEffectWidget
{
    Q_OBJECT

public:
    explicit VolumeWidget(VolumeModel *model, QWidget *parent = nullptr);
    ~VolumeWidget() override;

Q_SIGNALS:
    void visibleChanged(bool);
    void rightIconClick();

protected:
    void initUi();
    void initConnection();

    void showEvent(QShowEvent *event) override;
    void hideEvent(QHideEvent *event) override;

private:
    const QString leftIcon();
    const QString rightIcon();

private:
    VolumeModel *m_model;
    SliderContainer *m_sliderContainer;
    AudioSink *m_defaultSink;
};

#endif // VOLUMEWIDGET_H
