// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
