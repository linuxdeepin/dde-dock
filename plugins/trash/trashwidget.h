// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef TRASHWIDGET_H
#define TRASHWIDGET_H

#include "popupcontrolwidget.h"

#include "org_freedesktop_filemanager.h"

#include <QWidget>
#include <QPixmap>
#include <QAction>
#include <QIcon>

using  DBusFileManager1 = org::freedesktop::FileManager1;

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
    bool getDragging() const;
    void setDragging(bool state);

signals:
    void requestContextMenu() const;

protected:
    void dragEnterEvent(QDragEnterEvent *e) override;
    void dragMoveEvent(QDragMoveEvent *e) override;
    void dragLeaveEvent(QDragLeaveEvent *e) override;
    void dropEvent(QDropEvent *e) override;
    void paintEvent(QPaintEvent *e) override;

private slots:
    void removeApp(const QString &appKey);
    void moveToTrash(const QUrl &url);

private:
    PopupControlWidget *m_popupApplet;
    DBusFileManager1 *m_fileManagerInter;

    bool m_dragging; // item是否被拖入回收站

    QPixmap m_icon;
    QIcon m_defaulticon;
};

#endif // TRASHWIDGET_H
