/*
 * Copyright (C) 2018 ~ 2025 Deepin Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng_cm@deepin.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng_cm@deepin.com>
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
#ifndef TRAYDELEGATE_H
#define TRAYDELEGATE_H

#include "constants.h"

#include <QStyledItemDelegate>

#define ITEM_SIZE 30
// 托盘图标固定16个像素
#define ICON_SIZE 16
#define ITEM_SPACING 5

struct WinInfo;
class ExpandIconWidget;
class QListView;

class TrayDelegate : public QStyledItemDelegate
{
    Q_OBJECT

public:
    explicit TrayDelegate(QListView *view, QObject *parent = nullptr);
    void setPositon(Dock::Position position);

Q_SIGNALS:
    void removeRow(const QModelIndex &) const;
    void visibleChanged(const QModelIndex &, bool) const;
    void requestDrag(bool) const;

private Q_SLOTS:
    void onRequestDrag(bool on);

protected:
    QWidget *createEditor(QWidget *parent, const QStyleOptionViewItem &option, const QModelIndex &index) const Q_DECL_OVERRIDE;
    void setEditorData(QWidget *editor, const QModelIndex &index) const override ;
    QSize sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const Q_DECL_OVERRIDE;
    void updateEditorGeometry(QWidget *editor, const QStyleOptionViewItem &option, const QModelIndex &index) const override;
    void paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const override;

private:
    ExpandIconWidget *expandWidget();
    bool needDrawBackground() const;

private:
    Dock::Position m_position;
    QListView *m_listView;
};

#endif // TRAYDELEGATE_H
