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

#ifndef SOUNDITEM_H
#define SOUNDITEM_H

#include "soundapplet.h"
#include "org_deepin_dde_audio1_sink.h"

#include <QWidget>
#include <QIcon>

#define SOUND_KEY "sound-item-key"

using DBusSink = org::deepin::dde::audio1::Sink;

namespace Dock {
class TipsWidget;
}
class SoundItem : public QWidget
{
    Q_OBJECT

public:
    explicit SoundItem(QWidget *parent = 0);

    QWidget *tipsWidget();
    QWidget *popupApplet();

    const QString contextMenu();
    void invokeMenuItem(const QString menuId, const bool checked);

    void refreshIcon();
    void refreshTips(const int volume, const bool force = false);
    QPixmap pixmap() const;
    QPixmap pixmap(DGuiApplicationHelper::ColorType colorType, int iconWidth, int iconHeight) const;

signals:
    void requestContextMenu() const;

protected:
    void resizeEvent(QResizeEvent *e);
    void wheelEvent(QWheelEvent *e);
    void paintEvent(QPaintEvent *e);

private slots:
    void sinkChanged(DBusSink *sink);
    void refresh(const int volume);

private:
    Dock::TipsWidget *m_tipsLabel;
    QScopedPointer<SoundApplet> m_applet;
    DBusSink *m_sinkInter;
    QPixmap m_iconPix