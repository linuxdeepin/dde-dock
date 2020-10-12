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

#ifndef VOLUMESLIDER_H
#define VOLUMESLIDER_H

#include <DSlider>
#include <QTimer>

class VolumeSlider : public DTK_WIDGET_NAMESPACE::DSlider
{
    Q_OBJECT

public:
    explicit VolumeSlider(QWidget *parent = 0);

    void setValue(const int value);

signals:
    void requestPlaySoundEffect() const;

protected:
    void mousePressEvent(QMouseEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;
    void wheelEvent(QWheelEvent *e) override;

private slots:
    void onTimeout();

private:
    bool m_pressed;
    QTimer *m_timer;
};

#endif // VOLUMESLIDER_H
