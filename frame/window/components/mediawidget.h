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
#ifndef MEDIAWIDGET_H
#define MEDIAWIDGET_H

#include "mediaplayermodel.h"

#include <DBlurEffectWidget>

class QLabel;
class MusicButton;

DWIDGET_USE_NAMESPACE

class MediaWidget : public DBlurEffectWidget
{
    Q_OBJECT

public:
    explicit MediaWidget(QWidget *parent = nullptr);
    ~MediaWidget() override;

Q_SIGNALS:
    void visibleChanged(bool);

protected:
    void showEvent(QShowEvent *event) override;
    void hideEvent(QHideEvent *event) override;

private Q_SLOTS:
    void statusChanged(const MediaPlayerModel::PlayStatus &newStatus);
    void onPlayClicked();
    void onNext();
    void onUpdateMediaInfo();

private:
    void initUi();
    void initConnection();

private:
    QLabel *m_musicIcon;
    QLabel *m_musicName;
    QLabel *m_musicSinger;
    MusicButton *m_pausePlayButton;
    MusicButton *m_nextButton;
};

// 音乐播放按钮
class MusicButton : public QWidget
{
    Q_OBJECT

Q_SIGNALS:
    void clicked();

public:
    enum ButtonType { Playing = 0, Pause, Next };

public:
    MusicButton(QWidget *parent = Q_NULLPTR);
    ~MusicButton() override;

    inline void setButtonType(const ButtonType &bt) {
        m_buttonType = bt;
        update();
    }

protected:
    void paintEvent(QPaintEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;

private:
    int getIconHeight() const;

private:
    ButtonType m_buttonType;
};

#endif // MEDIAWIDGER_H
