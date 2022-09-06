// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef APPDRAGWIDGET_H
#define APPDRAGWIDGET_H

#include "constants.h"

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

    void setAppPixmap(const QPixmap &pix);
    void setDockInfo(Dock::Position dockPosition, const QRect &dockGeometry);
    void setOriginPos(const QPoint position);
    void setPixmapOpacity(qreal opacity);
    bool isRemoveAble(const QPoint &curPos);
    void setItem(DockItem *item) { m_item = item; }
    static bool isRemoveable(const Dock::Position &dockPos, const QRect &doctRect);
    void showRemoveAnimation();
    void showGoBackAnimation();

signals:
    void requestRemoveItem();
    void requestRemoveSelf(bool);

protected:
    void mouseMoveEvent(QMouseEvent *event) override;
    void dragEnterEvent(QDragEnterEvent *event) override;
    void dragMoveEvent(QDragMoveEvent *event) override;
    void dropEvent(QDropEvent *event) override;
    void hideEvent(QHideEvent *event) override;
    void moveEvent(QMoveEvent *event) override;
    void enterEvent(QEvent *event) override;
    bool event(QEvent *event) override;

private:
    void initAnimations();
    void onRemoveAnimationStateChanged(QAbstractAnimation::State newState,
            QAbstractAnimation::State oldState);
    const QPoint popupMarkPoint(Dock::Position pos);
    const QPoint topleftPoint() const;
    void showRemoveTips();
    bool isRemoveItem();

private Q_SLOTS:
    void onFollowMouse();

private:
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
    Dock::TipsWidget *m_removeTips;
    QScopedPointer<DockPopupWindow> m_popupWindow;
    /**
     * @brief m_distanceMultiple: 倍数
     * dock栏上应用区驻留应用被拖拽远离dock的距离除以dock的宽或者高（更小的一个）的比值
     */
    double m_distanceMultiple;

    bool m_bDragDrop = false; // 图标是否被拖拽
    DockItem *m_item;
    QPoint m_cursorPosition;
};

#endif /* APPDRAGWIDGET_H */
