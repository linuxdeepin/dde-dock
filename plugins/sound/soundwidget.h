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
#ifndef SOUNDWIDGET_H
#define SOUNDWIDGET_H

#include "org_deepin_dde_audio1.h"
#include "org_deepin_dde_audio1_sink.h"

#include <DBlurEffectWidget>
#include <QWidget>

class QDBusMessage;
class SliderContainer;
class QLabel;
class AudioSink;

DWIDGET_USE_NAMESPACE

using DBusAudio = org::deepin::dde::Audio1;
using DBusSink = org::deepin::dde::audio1::Sink;

class SoundWidget : public QWidget
{
    Q_OBJECT

public:
    explicit SoundWidget(QWidget *parent = nullptr);
    ~SoundWidget() override;

Q_SIGNALS:
    void rightIconClick();

protected:
    void initUi();
    void initConnection();

private:
    const QString leftIcon();
    const QString rightIcon();
    void convertThemePixmap(QPixmap &pixmap);
    bool existActiveOutputDevice() const;

private Q_SLOTS:
    void onThemeTypeChanged();

private:
    DBusAudio *m_dbusAudio;
    SliderContainer *m_sliderContainer;
    DBusSink *m_defaultSink;
};

#endif // VOLUMEWIDGET_H
