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
#ifndef SOUNDDEVICESWIDGET_H
#define SOUNDDEVICESWIDGET_H

#include "org_deepin_daemon_audio.h"
#include "org_deepin_daemon_audio_sink.h"

#include <DStyledItemDelegate>

#include <QWidget>

namespace Dtk { namespace Widget { class DListView; } }

using namespace Dtk::Widget;

class SliderContainer;
class QStandardItemModel;
class QLabel;
class VolumeModel;
class AudioSink;
class SettingDelegate;

using DBusAudio = org::deepin::daemon::Audio1;
using DBusSink = org::deepin::daemon::audio1::Sink;

class SoundDevicesWidget : public QWidget
{
    Q_OBJECT

public:
    explicit SoundDevicesWidget(QWidget *parent = nullptr);
    ~SoundDevicesWidget() override;

protected:
    bool eventFilter(QObject *watcher, QEvent *event) override;

private:
    void initUi();
    void initConnection();
    QString leftIcon();
    QString rightIcon();
    const QString soundIconFile(const AudioPort port) const;

    void resizeHeight();

    void resetVolumeInfo();
    uint audioPortCardId(const AudioPort &audioport) const;

private Q_SLOTS:
    void onSelectIndexChanged(const QModelIndex &index);
    void onDefaultSinkChanged(const QDBusObjectPath & value);
    void onAudioDevicesChanged();

private:
    QWidget *m_sliderParent;
    SliderContainer *m_sliderContainer;
    QLabel *m_descriptionLabel;
    DListView *m_deviceList;
    DBusAudio *m_volumeModel;
    DBusSink *m_audioSink;
    QStandardItemModel *m_model;
    SettingDelegate *m_delegate;
};

#endif // VOLUMEDEVICESWIDGET_H
