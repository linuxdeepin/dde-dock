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

#ifndef PLUGINWIDGET_H
#define PLUGINWIDGET_H

#include "constants.h"

#include <QWidget>
#include <QIcon>

class ShutdownWidget : public QWidget
{
    Q_OBJECT

public:
    explicit ShutdownWidget(QWidget *parent = 0);

protected:
    void paintEvent(QPaintEvent *e) override;
    void mousePressEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *event) override;
    void leaveEvent(QEvent *event) override;

private:
    const QPixmap loadSvg(const QString &fileName, const QSize &size) const;

private:
    Dock::DisplayMode m_displayMode;
    bool m_hover;
    bool m_pressed;
    QIcon m_icon;
};

#endif // PLUGINWIDGET_H
