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

#ifndef DISKPLUGINITEM_H
#define DISKPLUGINITEM_H

#include "constants.h"

#include <QWidget>
#include <QPixmap>

class DiskPluginItem : public QWidget
{
    Q_OBJECT

public:
    explicit DiskPluginItem(QWidget *parent = 0);

signals:
    void requestContextMenu() const;

public slots:
    void setDockDisplayMode(const Dock::DisplayMode mode);

protected:
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);
    QSize sizeHint() const;

private:
    void updateIcon();

private:
    Dock::DisplayMode m_displayMode;

    QPixmap m_icon;
};

#endif // DISKPLUGINITEM_H
