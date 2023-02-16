// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef APPDRAGWIDGET_H
#define APPDRAGWIDGET_H

#include "constants.h"
#include "screenspliter.h"
#include "utils.h"

#include <QPixmap>
#include <QGraphicsObject>
#include <QGraphicsView>
#include <QPainter>
#include <QMouseEvent>
#include <QTimer>
#include <QPropertyAnimation>
#include <QParallelAnimationGroup>
#include <QWidget>

#include "../widgets/tipswidget.h"
#include "dockpopupwindow.h"
#include "dockitem.h"

class QDrag;
class DockScreen;

class AppGraphicsObject : public QGraphicsObject
{
public:
    explicit AppGraphicsObject(QGraphicsItem *parent = Q_NULLPTR)
        : QGraphicsObject(parent) {}
    ~AppGraphicsObject() override {}

    void setAppPixmap(QPixmap pix)
    {
        m_appPixmap = pix;
        resetProperty();
        update();
    }

    void resetProperty()
    {
        setScale(1.0);
        setRotation(0);
        setOpacity(1.0);
        update();
    }

    QRectF boundingRect() const override
    {
        return m_appPixmap.rect();
    }

    void paint(QPainter *painter, const QStyleOptionGraphicsItem *option, QWidget *widget = Q_NULLPTR) override {
        Q_UNUSED(option);
        Q_UNUSED(widget);

        Q_ASSERT(!m_appPixmap.isNull());

        painter->drawPixmap(QPoint(0, 0), m_appPixmap);
    }

private:
    QPixmap m_appPixmap;
};

class AppDragWidget : public QGraphicsView
{
    Q_OBJECT

public:
    explicit AppDragWidget(QWidget *parent = Q_NULLPTR);

    void execFinished();
    void setAppPixmap(const QPixmap &pix);
    void setDockInfo(Dock::Position dockPosition, const QRect &dockGeometry);
    void setOriginPos(const QPoint position);
    void setPixmapOpacity(qreal opacity);
    void setItem(DockItem *item) { m_item = item; }
    void showRemoveAnimation();
    void showGoBackAnimation();

signals:
    void requestChangedArea(QRect);
    void requestSplitWindow(ScreenSpliter::SplitDirection);

protected:
    void mouseMoveEvent(QMouseEvent *event) override;
    void dragEnterEvent(QDragEnterEvent *event) override;
    void dragMoveEvent(QDragMoveEvent *event) override;
    void dropEvent(QDropEvent *event) override;
    void hideEvent(QHideEvent *event) override;
    void enterEvent(QEvent *event) override;

private:
    void initAnimations();
    void onRemoveAnimationStateChanged(QAbstractAnimation::State newState,
            QAbstractAnimation::State oldState);
    const QPoint popupMarkPoint(Dock::Position pos);
    const QPoint topleftPoint() const;
    bool canSplitWindow(const QPoint &pos) const;
    ScreenSpliter::SplitDirection splitPosition() const;
    QRect splitGeometry(const QPoint &pos) const;
    void initWaylandEnv();

    void dropHandler(const QPoint &pos);
    void moveHandler(const QPoint &pos);
    void moveCurrent(const QPoint &destPos);
    void adjustDesktopGeometry(QRect &rect) const;

private Q_SLOTS:
    void onFollowMouse();
    void onButtonRelease(int, int x, int y, const QString &);
    void onCursorMove(int x, int y, const QString &);

protected:
    QScopedPointer<AppGraphicsObject> m_object;
    QGraphicsScene *m_scene;
    QTimer *m_followMouseTimer;
    QPropertyAnimation *m_animScale;
    QPropertyAnimation *m_animRotation;
    QPropertyAnimation *m_animOpacity;
    QParallelAnimationGroup *m_animGroup;
    QPropertyAnimation *m_goBackAnim;

    Dock::Position m_dockPosition;
    QRect m_dockGeometry;
    QPoint m_originPoint;
    QSize m_iconSize;
    QScopedPointer<DockPopupWindow> m_popupWindow;
    /**
     * @brief m_distanceMultiple: 倍数
     * dock栏上应用区驻留应用被拖拽远离dock的距离除以dock的宽或者高（更小的一个）的比值
     */
    double m_distanceMultiple;

    bool m_bDragDrop = false; // 图标是否被拖拽
    DockItem *m_item;
    QRect m_lastMouseGeometry;
    DockScreen *m_dockScreen;
};

class QuickDragWidget : public AppDragWidget
{
    Q_OBJECT

Q_SIGNALS:
    void requestDropItem(QDropEvent *);
    void requestDragMove(QDragMoveEvent *);

public:
    explicit QuickDragWidget(QWidget *parent = Q_NULLPTR);
    ~QuickDragWidget() override;

protected:
    void dropEvent(QDropEvent *event) override;
    void dragMoveEvent(QDragMoveEvent *event) override;

private:
    bool isRemoveAble(const QPoint &curPos);
};

#endif /* APPDRAGWIDGET_H */
