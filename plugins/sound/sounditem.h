/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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
#include "dbus/dbussink.h"

#include <QWidget>

class SoundItem : public QWidget
{
    Q_OBJECT

public:
    explicit SoundItem(QWidget *parent = 0);

    QWidget *tipsWidget();
    QWidget *popupApplet();

    const QString contextMenu() const;
    void invokeMenuItem(const QString menuId, const bool checked);

signals:
    void requestContextMenu() const;

protected:
    QSize sizeHint() const;
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void wheelEvent(QWheelEvent *e);
    void paintEvent(QPaintEvent *e);

private slots:
    void refershIcon();
    void refershTips(const bool force = false);
    void sinkChanged(DBusSink *sink);

private:
    QLabel *m_tipsLabel;
    SoundApplet *m_applet;
    DBusSink *m_sinkInter;
    QPixmap m_iconPixmap;
};

#endif // SOUNDITEM_H
