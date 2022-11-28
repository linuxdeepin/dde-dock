/*
 * Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
 *
 * Author:     zhaoyingzhen <zhaoyingzhen@uniontech.com>
 *
 * Maintainer: zhaoyingzhen <zhaoyingzhen@uniontech.com>
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
#ifndef DEVITEMDELEGATE_H
#define DEVITEMDELEGATE_H

#include <DStyledItemDelegate>

/*!
 * \brief The DevItemDelegate class
 */
class DevItemDelegate : public QStyledItemDelegate
{
    Q_OBJECT
public:
    enum DevItemDataRole {
        StaticDataRole         = Qt::UserRole + 1,  // 静态信息
        MachinePathDataRole    = Qt::UserRole + 2,  // machinePath, 可唯一代表一个设备
        DegreeDataRole         = Qt::UserRole + 3,  // degree 绘制waiting使用的参数
        ResultDataRole         = Qt::UserRole + 4   // 连接结果
    };

    enum ResultState {
        None,
        Connecting,
        Connected
    };

    struct DevItemData {
        QString checkedIconPath;
        QString iconPath;
        QString text;
    };

public:
    explicit DevItemDelegate(QObject *parent = nullptr);

protected:
    void paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const Q_DECL_OVERRIDE;
    QSize sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const Q_DECL_OVERRIDE;

private:
    void drawWaitingState(QPainter *painter, const QRect &rect, int degree) const;
    void drawResultState(QPainter *painter, const QRect &rect) const;
    QList<QColor> createDefaultIndicatorColorList(QColor color) const;
};

Q_DECLARE_METATYPE(DevItemDelegate::DevItemData)

#endif // DEVITEMDELEGATE_H
