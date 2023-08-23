// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SOUNDDEVICESWIDGET_H
#define SOUNDDEVICESWIDGET_H

#include "org_deepin_dde_audio1.h"
#include "org_deepin_dde_audio1_sink.h"

#include <DStyledItemDelegate>

#include <QWidget>
#include <DGuiApplicationHelper>
#include <QTimer>

namespace Dtk { namespace Widget { class DListView; } }

using namespace Dtk::Widget;
using namespace Dtk::Gui;
namespace Dock {
class TipsWidget;
}

class SliderContainer;
class QStandardItemModel;
class QLabel;
class VolumeModel;
class AudioSink;
class SoundDevicePort;

using DBusAudio = org::deepin::dde::Audio1;
using DBusSink = org::deepin::dde::audio1::Sink;

class SoundDevicesWidget : public QWidget
{
    Q_OBJECT

public:
    explicit SoundDevicesWidget(QWidget *parent = nullptr);
    ~SoundDevicesWidget() override;

    QWidget* tipsWidget();
    QPixmap pixmap() const;
    QPixmap pixmap(DGuiApplicationHelper::ColorType colorType, int iconWidth, int iconHeight) const;

Q_SIGNALS:
    void enableChanged(bool);
    void requestHide();
    void iconChanged();

protected:
    bool eventFilter(QObject *watcher, QEvent *event) override;

private:
    void initUi();
    void initConnection();

    void resizeHeight();

    void resetVolumeInfo();
    uint audioPortCardId(const AudioPort &audioport) const;

    SoundDevicePort *findPort(const QString &portId, const uint &cardId) const;
    void startAddPort(SoundDevicePort *port);
    void startRemovePort(const QString &portId, const uint &cardId);

    void addPort(const SoundDevicePort *port);
    void removePort(const QString &portId, const uint &cardId);

    void activePort(const QString &portId, const uint &cardId);

    void removeDisabledDevice(QString portId, unsigned int cardId);

    void deviceEnabled(bool enable);


    void refreshTips(const int volume, const bool force);

private Q_SLOTS:
    void onSelectIndexChanged(const QModelIndex &index);
    void onDefaultSinkChanged(const QDBusObjectPath & value);
    void onAudioDevicesChanged();

private:
    Dock::TipsWidget *m_tipsLabel;
    SliderContainer *m_sliderContainer;
    QLabel *m_descriptionLabel;
    DListView *m_deviceList;
    DBusAudio *m_soundInter;
    DBusSink *m_sinkInter;
    QStandardItemModel *m_model;
    QList<SoundDevicePort *> m_ports;
    QTimer *m_setVolumeTimer;
};

#endif // VOLUMEDEVICESWIDGET_H
