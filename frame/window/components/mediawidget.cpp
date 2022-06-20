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
#include "brightnessmodel.h"
#include "mediawidget.h"

#include <DFontSizeManager>

#include <QLabel>
#include <QHBoxLayout>
#include <QVBoxLayout>
#include <QEvent>
#include <QPainter>
#include <QDebug>
#include <QPainterPath>

DWIDGET_USE_NAMESPACE

MediaWidget::MediaWidget(QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_musicIcon(new QLabel(this))
    , m_musicName(new QLabel(this))
    , m_musicSinger(new QLabel(this))
    , m_pausePlayButton(new MusicButton(this))
    , m_nextButton(new MusicButton(this))
{
    initUi();
    initConnection();
}

MediaWidget::~MediaWidget()
{
}

void MediaWidget::showEvent(QShowEvent *event)
{
    DBlurEffectWidget::showEvent(event);
    Q_EMIT visibleChanged(true);
}

void MediaWidget::hideEvent(QHideEvent *event)
{
    DBlurEffectWidget::hideEvent(event);
    Q_EMIT visibleChanged(false);
}

void MediaWidget::statusChanged(const MediaPlayerModel::PlayStatus &newStatus)
{
    switch (newStatus) {
    case MediaPlayerModel::PlayStatus::Play: {
        m_pausePlayButton->setButtonType(MusicButton::ButtonType::Pause);
        break;
    }
    case MediaPlayerModel::PlayStatus::Stop:
    case MediaPlayerModel::PlayStatus::Pause: {
        m_pausePlayButton->setButtonType(MusicButton::ButtonType::Playing);
        break;
    }
    default: break;
    }
}

void MediaWidget::onPlayClicked()
{
    // 设置当前的播放状态
    MediaPlayerModel *player = MediaPlayerModel::instance();
    if (player->status() == MediaPlayerModel::PlayStatus::Play)
        player->setStatus(MediaPlayerModel::PlayStatus::Pause);
    else
        player->setStatus(MediaPlayerModel::PlayStatus::Play);
}

void MediaWidget::onNext()
{
    // 播放下一曲
    MediaPlayerModel *player = MediaPlayerModel::instance();
    player->playNext();
}

void MediaWidget::initUi()
{
    m_pausePlayButton->setFixedWidth(20);
    m_nextButton->setFixedWidth(20);

    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(20, 0, 20, 0);
    mainLayout->addWidget(m_musicIcon);

    QWidget *infoWidget = new QWidget(this);
    QVBoxLayout *infoLayout = new QVBoxLayout(infoWidget);
    infoLayout->addWidget(m_musicName);
    infoLayout->addWidget(m_musicSinger);
    mainLayout->addWidget(infoWidget);
    mainLayout->addWidget(m_pausePlayButton);
    mainLayout->addSpacing(25);
    mainLayout->addWidget(m_nextButton);

    m_musicIcon->setFixedSize(32, 32);
    m_musicName->setFont(DFontSizeManager::instance()->t8());
    m_musicSinger->setFont(DFontSizeManager::instance()->t10());
    setVisible(MediaPlayerModel::instance()->isActived());
}

void MediaWidget::initConnection()
{
    MediaPlayerModel *mediaPlayer = MediaPlayerModel::instance();
    connect(mediaPlayer, &MediaPlayerModel::startStop, this, [ this, mediaPlayer ](bool startOrStop) {
        setVisible(startOrStop);
        m_nextButton->setEnabled(mediaPlayer->canGoNext());
        onUpdateMediaInfo();
        statusChanged(mediaPlayer->status());
    });
    connect(mediaPlayer, &MediaPlayerModel::metadataChanged, this, &MediaWidget::onUpdateMediaInfo);
    connect(mediaPlayer, &MediaPlayerModel::statusChanged, this, &MediaWidget::statusChanged);
    connect(m_pausePlayButton, &MusicButton::clicked, this, &MediaWidget::onPlayClicked);
    connect(m_nextButton, &MusicButton::clicked, this, &MediaWidget::onNext);

    m_pausePlayButton->setButtonType(mediaPlayer->status() == MediaPlayerModel::PlayStatus::Play ?
                                         MusicButton::ButtonType::Pause : MusicButton::ButtonType::Playing);
    m_nextButton->setButtonType(MusicButton::ButtonType::Next);
}

void MediaWidget::onUpdateMediaInfo()
{
    MediaPlayerModel *mediaPlayer = MediaPlayerModel::instance();
    m_musicName->setText(mediaPlayer->name());
    QString file = mediaPlayer->iconUrl();
    if (file.startsWith("file:///"))
        file.replace("file:///", "/");
    m_musicIcon->setPixmap(QPixmap(file).scaled(m_musicIcon->size()));
    m_musicSinger->setText(mediaPlayer->artist());
}

/**
 * @brief 音乐播放的相关按钮
 * @param parent
 */

MusicButton::MusicButton(QWidget *parent)
    : QWidget(parent)
{
    installEventFilter(this);
}

MusicButton::~MusicButton()
{
}

int MusicButton::getIconHeight() const
{
    switch (m_buttonType) {
    case ButtonType::Pause:
        return 21;
    case ButtonType::Next:
    case ButtonType::Playing:
        return 18;
    }

    return 18;
}

void MusicButton::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);
    int ctrlHeight = getIconHeight();

    int width = this->width();
    int height = this->height();
    int startX = 2;
    int startY = (height - ctrlHeight) / 2;
    QColor color(0, 0, 0);
    QPainter painter(this);
    painter.save();
    painter.setRenderHint(QPainter::Antialiasing);
    painter.setPen(color);
    painter.setBrush(color);
    if (m_buttonType == ButtonType::Pause) {
        painter.drawRect(QRect(startX, startY, 6, ctrlHeight));
        painter.drawRect(QRect(width - 6 - 2, startY, 6, ctrlHeight));
    } else {
        QPainterPath trianglePath;
        trianglePath.moveTo(startX, startY);
        trianglePath.lineTo(width - 6, height / 2);
        trianglePath.lineTo(startX, startY + ctrlHeight);
        trianglePath.lineTo(startX, startY);
        painter.drawPath(trianglePath);
        if (m_buttonType == ButtonType::Next)
            painter.drawRect(width - 6, startY, 2, ctrlHeight);
    }
    painter.restore();
}

void MusicButton::mouseReleaseEvent(QMouseEvent *event)
{
    Q_UNUSED(event);
    Q_EMIT clicked();
}
