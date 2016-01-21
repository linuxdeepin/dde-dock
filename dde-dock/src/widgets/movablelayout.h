#ifndef MOVABLELAYOUT_H
#define MOVABLELAYOUT_H

#include <QFrame>
#include <QList>
#include <QDragEnterEvent>
#include <QHBoxLayout>
#include <QEasingCurve>

class QDrag;
class QPropertyAnimation;

class MovableLayout : public QFrame
{
    Q_OBJECT
    Q_PROPERTY(QSize size READ size WRITE setFixedSize)
public:
    enum MoveDirection {
        MoveLeftToRight,
        MoveRightToLeft,
        MoveTopToBottom,
        MoveBottomToTop,
        MoveUnknow
    };

    explicit MovableLayout(QWidget *parent = 0);
    explicit MovableLayout(QBoxLayout::Direction direction, QWidget *parent = 0);

    bool dragable() const;
    bool getAutoResize() const;
    int count() const;
    int hoverIndex() const;
    int getLayoutSpacing() const;
    int getAnimationDuration() const;
    int getHoverIndextByPos(const QPoint &pos);
    int indexOf(QWidget * const widget, int from = 0) const;

    QWidget *dragingWidget() const;
    QWidget *widget(int index) const;
    QList<QWidget *> widgets() const;
    QSize getDefaultSpacingItemSize() const;
    QBoxLayout::Direction direction() const;
    QEasingCurve::Type getAnimationCurve() const;

    void restoreDragingWidget();
    void removeWidget(int index);
    void addWidget(QWidget *widget);
    void removeWidget(QWidget *widget);
    void insertWidget(int index, QWidget *widget);

    void setDuration(int v);
    void setDragable(bool v);
    void setLayoutSpacing(int spacing);
    void setAutoResize(bool autoResize);
    void setAlignment(Qt::Alignment alignment);
    void setEasingCurve(QEasingCurve::Type curve);
    void setAnimationDuration(int animationDuration);
    void setDirection(QBoxLayout::Direction direction);
    void setAnimationCurve(const QEasingCurve::Type &animationCurve);
    void setDefaultSpacingItemSize(const QSize &defaultSpacingItemSize);

signals:
    void startDrag(QDrag*);
    void drop(QDropEvent *event);
    void requestSpacingItemsDestroy(bool immediately);
    void sizeChanged(QResizeEvent *event);
    void dragLeaved(QDragLeaveEvent *event);
    void dragEntered(QDragEnterEvent *event);

private:
    bool event(QEvent *e);
    void mouseMoveEvent(QMouseEvent *event);
    void dragEnterEvent(QDragEnterEvent *event);
    void dragLeaveEvent(QDragLeaveEvent *event);
    void dragMoveEvent(QDragMoveEvent *event);
    void dropEvent(QDropEvent *event);
    void resizeEvent(QResizeEvent *event);

private:
    void initSizeAniamtion();
    void storeDragingWidget();
    void handleDrag(const QPoint &pos);
    void updateCurrentHoverInfo(int index, const QPoint &pos);
    void addSpacingItem(QWidget *souce, MoveDirection md, const QSize &size);
    void insertSpacingItemToLayout(int index, const QSize &size, bool immediately = false);
    MoveDirection getVMoveDirection(int index, const QPoint &pos);
    MoveDirection getHMoveDirection(int index, const QPoint &pos);

private:
    int m_lastHoverIndex;
    int m_animationDuration;
    bool m_hoverToSpacing;
    bool m_autoResize;
    bool m_dragable;
    QWidget *m_draginItem;
    QHBoxLayout *m_layout;
    QList<QWidget *> m_widgetList;
    QSize m_defaultSpacingItemSize;
    MoveDirection m_vMoveDirection;
    MoveDirection m_hMoveDirection;
    QPropertyAnimation *m_sizeAnimation;
    QEasingCurve::Type m_animationCurve;
};

#endif // MOVABLELAYOUT_H
