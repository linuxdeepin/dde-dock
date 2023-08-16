// Copyright (C) 2018 ~ 2025 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef TRAYDELEGATE_H
#define TRAYDELEGATE_H

#include "constants.h"

#include <QStyledItemDelegate>

#define ITEM_SIZE 30
// 托盘图标固定20个像素
// 16 x 16 will make the icon be cut, reserve 2 pixels for each side to prevent being cut
#define ICON_SIZE 20
#define ITEM_SPACING 5

class ExpandIconWidget;
class QListView;
class PluginsItemInterface;

class TrayDelegate : public QStyledItemDelegate
{
    Q_OBJECT

public:
    static TrayDelegate *getDockTrayDelegate(QListView *view, QObject *parent = nullptr);
    static TrayDelegate *getIconTrayDelegate(QListView *view, QObject *parent = nullptr);

    void setPositon(Dock::Position position);

Q_SIGNALS:
    void removeRow(const QModelIndex &) const;
    void requestDrag(bool) const;
    void requestHide();

private Q_SLOTS:
    void onUpdateExpand(bool on);

protected:
    QWidget *createEditor(QWidget *parent, const QStyleOptionViewItem &option, const QModelIndex &index) const Q_DECL_OVERRIDE;
    void setEditorData(QWidget *editor, const QModelIndex &index) const override ;
    QSize sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const Q_DECL_OVERRIDE;
    void updateEditorGeometry(QWidget *editor, const QStyleOptionViewItem &option, const QModelIndex &index) const override;
    void paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const override;

private:
    explicit TrayDelegate(QListView *view, QObject *parent = nullptr);

    ExpandIconWidget *expandWidget();
    bool isPopupTray() const;

private:
    Dock::Position m_position;
    QListView *m_listView;
};

#endif // TRAYDELEGATE_H
