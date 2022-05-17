/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#ifndef TRAYGRIDVIEW_H
#define TRAYGRIDVIEW_H

#include "constants.h"

#include <DListView>

#include <QPropertyAnimation>

DWIDGET_USE_NAMESPACE

class TrayGridView : public DListView
{
    Q_OBJECT

public:
    explicit TrayGridView(QWidget *parent = Q_NULLPTR);

    void setPosition(Dock::Position position);
    QSize suitableSize() const;
    void setDragDistance(int pixel);
    void setAnimationProperty(const QEasingCurve::Type easing, const int duringTime = 250);
    void moveAnimation();
    const QModelIndex modelIndex(const int index) const;
    const QRect indexRect(const QModelIndex &index) const;
    void dropSwap();

Q_SIGNALS:
    void requestRemove(const QString &);
    void dragLeaved();
    void dragEntered();

public Q_SLOTS:
    void clearDragModelIndex();

protected:
    void mousePressEvent(QMouseEvent *e) Q_DECL_OVERRIDE;
    void mouseMoveEvent(QMouseEvent *e) Q_DECL_OVERRIDE;
    void mouseReleaseEvent(QMouseEvent *e) Q_DECL_OVERRIDE;

    void dragEnterEvent(QDragEnterEvent *e) Q_DECL_OVERRIDE;
    void dragLeaveEvent(QDragLeaveEvent *e) Q_DECL_OVERRIDE;
    void dragMoveEvent(QDragMoveEvent *e) Q_DECL_OVERRIDE;
    void dropEvent(QDropEvent *e) Q_DECL_OVERRIDE;
    bool beginDrag(Qt::DropActions supportedActions);

private:
    void initUi();
    void createAnimation(const int pos, const bool moveNext, const bool isLastAni);
    const QModelIndex getIndexFromPos(QPoint currentPoint) const;

private:
    QEasingCurve::Type m_aniCurveType;
    int m_aniDuringTime;

    QPoint m_dragPos;
    QPoint m_dropPos;

    int m_dragDistance;

    QTimer *m_aniStartTime;
    bool m_pressed;
    bool m_aniRunning;
    Dock::Position m_positon;
};

#endif // GRIDVIEW_H
