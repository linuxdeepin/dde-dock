/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *               2016 ~ 2018 dragondjf
 *
 * Author:     sbw <sbw@sbw.so>
 *             dragondjf<dingjiangfeng@deepin.com>
 *             zccrs<zhangjide@deepin.com>
 *             Tangtong<tangtong@deepin.com>
 *
 * Maintainer: dragondjf<dingjiangfeng@deepin.com>
 *             zccrs<zhangjide@deepin.com>
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

#ifndef TRASHWIDGET_H
#define TRASHWIDGET_H

#include "popupcontrolwidget.h"

#include <QWidget>
#include <QPixmap>
#include <QMenu>
#include <QAction>
#include <QIcon>

class TrashWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TrashWidget(QWidget *parent = 0);

    QWidget *popupApplet();

    const QString contextMenu() const;
    int trashItemCount() const;
    void invokeMenuItem(const QString &menuId, const bool checked);
    void updateIcon();
    void updateIconAndRefresh();

signals:
    void requestContextMenu() const;

protected:
    void dragEnterEvent(QDragEnterEvent *e) Q_DECL_OVERRIDE;
    void dragMoveEvent(QDragMoveEvent *e) Q_DECL_OVERRIDE;
    void dropEvent(QDropEvent *e) Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;

private slots:
    void removeApp(const QString &appKey);
    void moveToTrash(const QUrl &url);

private:
    PopupControlWidget *m_popupApplet;

    QPixmap m_icon;
    QIcon m_defaulticon;
};

#endif // TRASHWIDGET_H
