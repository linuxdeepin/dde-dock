// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "soundwidget.h"
#include "imageutil.h"
#include "imageutil.h"
#include "slidercontainer.h"

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
#include <QPainter>
#include <QIcon>
#include <QPixmap>

DGUI_USE_NAMESPACE

#define ICON_SIZE 18
#define BACKSIZE 36

SoundWidget::SoundWidget(QWidget *parent)
    : QWidget(parent)
    , m_dbusAudio(new DBusAudio("org.deepin.dde.Audio1", "/org/deepin/dde/Audio1", QDBusConnection::sessionBus(), this))
    , m_sliderContainer(new SliderContainer(this))
    , m_defaultSink(new DBusSink("org.deepin.dde.Audio1", m_dbusAudio->defaultSink().path(), QDBusConnection::sessionBus(), this))
{
    initUi();
    initConnection();
}

SoundWidget::~SoundWidget()
{
}

void SoundWidget::initUi()
{
    if (m_defaultSink)
        m_sliderContainer->updateSliderValue(m_defaultSink->volume() * 100);

    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(17, 0, 12, 0);
    mainLayout->addWidget(m_sliderContainer);

    onThemeTypeChanged();
    m_sliderContainer->setRange(0, std::round(m_dbusAudio->maxUIVolume() * 100.00));
    m_sliderContainer->setPageStep(2);

    SliderProxyStyle *proxy = new SliderProxyStyle;
    m_sliderContainer->setSliderProxyStyle(proxy);

    setEnabled(existActiveOutputDevice());
}

void SoundWidget::initConnection()
{
    connect(m_defaultSink, &DBusSink::VolumeChanged, this, [ this ](double value) { m_sliderContainer->updateSliderValue(std::round(value * 100.00)); });
    connect(m_defaultSink, &DBusSink::MuteChanged, this, [ = ] { m_sliderContainer->updateSliderValue(m_defaultSink->volume() * 100); });

    connect(m_dbusAudio, &DBusAudio::DefaultSinkChanged, this, [ this ](const QDBusObjectPath &value) {
        if (m_defaultSink)
            delete m_defaultSink;

        m_defaultSink = new DBusSink("org.deepin.dde.Audio1", value.path(), QDBusConnection::sessionBus(), this);
        connect(m_defaultSink, &DBusSink::VolumeChanged, this, [ this ](double value) { m_sliderContainer->updateSliderValue(std::round(value * 100.00)); });
        connect(m_defaultSink, &DBusSink::MuteChanged, this, [ = ] {
            m_sliderContainer->updateSliderValue(m_defaultSink->volume() * 100);
            m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon,
            QIcon::fromTheme(leftIcon()).pixmap(ICON_SIZE, ICON_SIZE), QSize(), 10);
        });

        m_sliderContainer->updateSliderValue(std::round(m_defaultSink->volume() * 100.00));
    });

    connect(m_dbusAudio, &DBusAudio::MaxUIVolumeChanged, this, [ = ] (double maxValue) {
        m_sliderContainer->setRange(0, std::round(maxValue * 100.00));
    });

    connect(m_sliderContainer, &SliderContainer::sliderValueChanged, this, [ this ](int value) {
        m_defaultSink->SetVolume(value * 0.01, true);
        if (m_defaultSink->mute()) {
            m_defaultSink->SetMuteQueued(false);
        }
    });

    connect(m_defaultSink, &DBusSink::MuteChanged, this, [ this ] {
        m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon,
            QIcon::fromTheme(leftIcon()).pixmap(ICON_SIZE, ICON_SIZE), QSize(), 10);
    });

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &SoundWidget::onThemeTypeChanged);

    connect(m_sliderContainer, &SliderContainer::iconClicked, this, [ this ](const SliderContainer::IconPosition icon) {
        switch (icon) {
        case SliderContainer::IconPosition::LeftIcon: {
            if (existActiveOutputDevice())
                m_defaultSink->SetMute(!m_defaultSink->mute());
            break;
        }
        case SliderContainer::IconPosition::RightIcon: {
            // 弹出音量选择对话框
            Q_EMIT rightIconClick();
            break;
        }
        }
    });
}

const QString SoundWidget::leftIcon()
{
    const bool mute = existActiveOutputDevice() ? m_defaultSink->mute() : true;
    return QString("audio-volume-%1-symbolic").arg(mute? "muted": "medium");
}

const QString SoundWidget::rightIcon()
{
    // TODO: broadcast ???
    // svg from display plugins
    return QString(":/icons/resources/broadcast.svg");
}

void SoundWidget::convertThemePixmap(QPixmap &pixmap)
{
    // 图片是黑色的，如果当前主题为白色主题，则无需转换
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::ColorType::LightType)
        return;

    // 如果是黑色主题，则转换成白色图像
    QPainter painter(&pixmap);
    painter.setCompositionMode(QPainter::CompositionMode_SourceIn);
    painter.fillRect(pixmap.rect(), Qt::white);
    painter.end();
}

/** 判断是否存在未禁用的声音输出设备
 * @brief SoundApplet::existActiveOutputDevice
 * @return 存在返回true,否则返回false
 */
bool SoundWidget::existActiveOutputDevice() const
{
    QString info = m_dbusAudio->property("CardsWithoutUnavailable").toString();

    QJsonDocument doc = QJsonDocument::fromJson(info.toUtf8());
    QJsonArray jCards = doc.array();
    for (QJsonValue cV : jCards) {
        QJsonObject jCard = cV.toObject();
        QJsonArray jPorts = jCard["Ports"].toArray();

        for (QJsonValue pV : jPorts) {
            QJsonObject jPort = pV.toObject();
            if (jPort["Direction"].toInt() == 1 && jPort["Enabled"].toBool())
                return true;
        }
    }

    return false;
}

void SoundWidget::onThemeTypeChanged()
{
    QPixmap leftPixmap = QIcon::fromTheme(leftIcon()).pixmap(ICON_SIZE, ICON_SIZE);
    QPixmap rightPixmap = QIcon::fromTheme(rightIcon()).pixmap(ICON_SIZE, ICON_SIZE);
    convertThemePixmap(rightPixmap);
    m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon, leftPixmap, QSize(), 10);
    m_sliderContainer->setIcon(SliderContainer::IconPosition::RightIcon, rightPixmap, QSize(BACKSIZE, BACKSIZE), 12);
}
