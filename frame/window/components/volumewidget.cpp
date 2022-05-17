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
#include "volumewidget.h"
#include "customslider.h"
#include "imageutil.h"
#include "volumemodel.h"

#include <DGuiApplicationHelper>

#include <QDBusConnectionInterface>
#include <QDBusInterface>
#include <QDBusPendingCall>
#include <QDBusPendingReply>
#include <QDebug>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QLabel>
#include <QEvent>
#include <QHBoxLayout>
#include <QMetaMethod>

DGUI_USE_NAMESPACE

#define ICON_SIZE 24

VolumeWidget::VolumeWidget(QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_volumeController(new VolumeModel(this))
    , m_volumnCtrl(new CustomSlider(Qt::Horizontal, this))
    , m_defaultSink(m_volumeController->defaultSink())
{
    initUi();
    initConnection();
}

VolumeWidget::~VolumeWidget()
{
}

void VolumeWidget::initUi()
{
    if (m_defaultSink)
        m_volumnCtrl->setValue(m_defaultSink->volume());

    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(20, 0, 20, 0);
    mainLayout->addWidget(m_volumnCtrl);

    m_volumnCtrl->setIconSize(QSize(36, 36));
    m_volumnCtrl->setLeftIcon(QIcon(leftIcon()));
    m_volumnCtrl->setRightIcon(QIcon(rightIcon()));

    bool existActiveOutputDevice = m_volumeController->existActiveOutputDevice();
    setEnabled(existActiveOutputDevice);
}

void VolumeWidget::initConnection()
{
    auto setCtrlVolumeValue = [this](int volume) {
        m_volumnCtrl->blockSignals(true);
        m_volumnCtrl->setValue(volume);
        m_volumnCtrl->blockSignals(false);
    };
    if (m_defaultSink)
        connect(m_defaultSink, &AudioSink::volumeChanged, this, setCtrlVolumeValue);

    connect(m_volumeController, &VolumeModel::defaultSinkChanged, this, [ this, setCtrlVolumeValue ](AudioSink *sink) {
        if (m_defaultSink)
            disconnect(m_defaultSink);

        m_defaultSink = sink;
        if (sink) {
            setCtrlVolumeValue(sink->volume());
            connect(m_defaultSink, &AudioSink::volumeChanged, this, setCtrlVolumeValue);
        }
    });

    connect(m_volumnCtrl, &DTK_WIDGET_NAMESPACE::DSlider::valueChanged, this, [ this ](int value) {
        AudioSink *sink = m_volumeController->defaultSink();
        if (sink)
            sink->setVolume(value, true);
    });

    connect(m_volumeController, &VolumeModel::muteChanged, this, [ this ] {
        m_volumnCtrl->setLeftIcon(QIcon(leftIcon()));
    });

    connect(m_volumnCtrl, &CustomSlider::iconClicked, this, [ this ](DSlider::SliderIcons icon, bool) {
        switch (icon) {
        case DSlider::SliderIcons::LeftIcon: {
            if (m_volumeController->existActiveOutputDevice())
                m_volumeController->setMute(!m_volumeController->isMute());
            break;
        }
        case DSlider::SliderIcons::RightIcon: {
            // 弹出音量选择对话框
            Q_EMIT rightIconClick();
            break;
        }
        }
    });
}


VolumeModel *VolumeWidget::model()
{
    return m_volumeController;
}

void VolumeWidget::showEvent(QShowEvent *event)
{
    DBlurEffectWidget::showEvent(event);
    Q_EMIT visibleChanged(true);
}

void VolumeWidget::hideEvent(QHideEvent *event)
{
    DBlurEffectWidget::hideEvent(event);
    Q_EMIT visibleChanged(false);
}

const QString VolumeWidget::leftIcon()
{
    bool existActiveOutputDevice = m_volumeController->existActiveOutputDevice();
    const bool mute = existActiveOutputDevice ? m_volumeController->isMute() : true;
    if (mute)
        return QString(":/icons/resources/audio-volume-muted-dark");

    return QString(":/icons/resources/volume");
}

const QString VolumeWidget::rightIcon()
{
    return QString(":/icons/resources/broadcast");
}
