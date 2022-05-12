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
    TrayGridView(QWidget *parent = Q_NULLPTR);

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
    void beginDrag(Qt::DropActions supportedActions);

private:
    void initUi();
    void createAnimation(const int pos, const bool moveNext, const bool isLastAni);

private:
    QEasingCurve::Type m_aniCurveType;
    int m_aniDuringTime;

    QPoint m_dragPos;
    QPoint m_dropPos;

    int m_dragDistance;

    QTimer *m_aniStartTime;
    bool m_pressed;
    bool m_aniRunning;
};

#endif // GRIDVIEW_H
