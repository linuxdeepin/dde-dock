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

#ifndef SOUNDAPPLET_H
#define SOUNDAPPLET_H

#include "componments/volumeslider.h"
#include "dbus/dbusaudio.h"
#include "dbus/dbussink.h"

#include <QScrollArea>
#include <QVBoxLayout>
#include <QLabel>
#include <QSlider>
#include <dimagebutton.h>

namespace Dock {
class TipsWidget;
}
class SoundApplet : public QScrollArea
{
    Q_OBJECT

public:
    explicit SoundApplet(QWidget *parent = 0);

    int volumeValue() const;
    int maxVolumeValue() const;
    VolumeSlider *mainSlider();
signals:
    void volumeChanged(const int value) const;
    void defaultSinkChanged(DBusSink *sink) const;

private slots:
    void defaultSinkChanged();
    void onVolumeChanged();
    void volumeSliderValueChanged();
    void sinkInputsChanged();
    void toggleMute();
    void onPlaySoundEffect();
    void increaseVolumeChanged();

private:
    void refreshIcon();

private:
    QWidget *m_centralWidget;
    QWidget *m_applicationTitle;
    Dtk::Widget::DImageButton *m_volumeBtn;
    QLabel *m_volumeIconMax;
    VolumeSlider *m_volumeSlider;
    Dock::TipsWidget *m_soundShow;
    QVBoxLayout *m_centralLayout;

    DBusAudio *m_audioInter;
    DBusSink *m_defSinkInter;
};

#endif // SOUNDAPPLET_H
