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

#ifndef SINKINPUTWIDGET_H
#define SINKINPUTWIDGET_H

#include "componments/volumeslider.h"
#include <com_deepin_daemon_audio_sinkinput.h>

#include <QFrame>
#include <QPainter>

#include <DIconButton>

DWIDGET_USE_NAMESPACE
using DBusSinkInput = com::deepin::daemon::audio::SinkInput;
namespace Dock {
    class TipsWidget;
}
class QLabel;
class SinkInputWidget : public QWidget
{
    Q_OBJECT

public:
    explicit SinkInputWidget(const QString &inputPath, QWidget *parent = nullptr);

private slots:
    void setVolume(const int value);
    void setMute();
    void setMuteIcon();
    void onPlaySoundEffect();
    void onVolumeChanged();

private:
    void refreshIcon();

private:
    DBusSinkInput *m_inputInter;

    DIconButton *m_appBtn;
    DIconButton *m_volumeBtnMin;
    QLabel *m_volumeIconMax;
    VolumeSlider *m_volumeSlider;
    Dock::TipsWidget *m_volumeLabel;
};

#endif // SINKINPUTWIDGET_H
