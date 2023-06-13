// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef APPMULTIITEM_H
#define APPMULTIITEM_H

#include "dockitem.h"
#include "dbusutil.h"
#include "taskmanager/windowinfomap.h"

class AppItem;

class AppMultiItem : public DockItem
{
    Q_OBJECT

    friend class AppItem;

public:
    AppMultiItem(AppItem *appItem, WId winId, const WindowInfo &windowInfo, QWidget *parent = Q_NULLPTR);
    ~AppMultiItem() override;

    QSize suitableSize(int size) const;
    AppItem *appItem() const;
    quint32 winId() const;
    const WindowInfo &windowInfo() const;

    ItemType itemType() const override;

protected:
    void paintEvent(QPaintEvent *) override;
    void mouseReleaseEvent(QMouseEvent *event) override;

private:
    void initMenu();
    void initConnection();

private Q_SLOTS:
    void onOpen();
    void onCurrentWindowChanged(uint32_t value);

private:
    AppItem *m_appItem;
    WindowInfo m_windowInfo;
    QPixmap m_pixmap;
    WId m_winId;
    QMenu *m_menu;
};

#endif // APPMULTIITEM_H
