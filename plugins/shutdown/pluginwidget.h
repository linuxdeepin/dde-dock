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

#ifndef PLUGINWIDGET_H
#define PLUGINWIDGET_H

#include "constants.h"
#include "dbus/dbuspower.h"

#include <QWidget>
#include <QTimer>

#define SHUTDOWN_KEY    "shutdown"

class PluginWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PluginWidget(QWidget *parent = 0);

signals:
    void requestContextMenu(const QString &itemKey) const;

protected:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);

private:
    const QPixmap loadSvg(const QString &fileName, const QSize &size) const;

private:
    void refershIconPixmap();

private:
    bool m_hover;
    Dock::DisplayMode m_displayMode;

    DBusPower *m_powerInter;
};

#endif // PLUGINWIDGET_H
