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

#ifndef CONTAINERWIDGET_H
#define CONTAINERWIDGET_H

#include <QWidget>
#include <QHBoxLayout>

class ContainerWidget : public QWidget
{
    Q_OBJECT

public:
    explicit ContainerWidget(QWidget *parent = 0);

    void addWidget(QWidget * const w);
    void removeWidget(QWidget * const w);
    int itemCount() const;
    const QList<QWidget *> itemList() const;

    bool allowDragEnter(QDragEnterEvent *e);

protected:
    void dragEnterEvent(QDragEnterEvent *e);
    void mousePressEvent(QMouseEvent *e);

private:
    QHBoxLayout *m_centralLayout;

    QList<QWidget *> m_itemList;
};

#endif // CONTAINERWIDGET_H
