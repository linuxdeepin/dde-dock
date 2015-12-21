#ifndef MOVABLELAYOUT_H
#define MOVABLELAYOUT_H

#include <QFrame>
#include <QList>
#include <QDragEnterEvent>
#include <QHBoxLayout>
#include <QEasingCurve>


class QPropertyAnimation;
class MovableSpacingItem : public QWidget
{
    Q_OBJECT
    Q_PROPERTY(QSize size READ size WRITE setFixedSize)
public:
    explicit MovableSpacingItem(int duration, QEasingCurve::Type easingType, const QSize &targetSize, QWidget *parent = 0);

signals:
    void growFinish();
    void declineFinish();

public slots:
    void StartGrow();
    void StartDecline();

private:
    QSize m_targetSize;
    QPropertyAnimation *m_growAnimation;
    QPropertyAnimation *m_declineAnimation;
};
/////////////////////////////////////////////////////////////////////////

class MovableLayout : public QFrame
{
    Q_OBJECT
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

    QWidget *widget(int index) const;
    QList<QWidget *> widgets() const;
    void addWidget(QWidget *widget);
    void removeWidget(int index);
    void removeWidget(QWidget *widget);

    int count() const;

    QBoxLayout::Direction direction() const;
    void setDirection(QBoxLayout::Direction direction);

    int getLayoutSpacing() const;
    void setLayoutSpacing(int spacing);

    QSize getDefaultSpacingItemSize() const;
    void setDefaultSpacingItemSize(const QSize &defaultSpacingItemSize);

    void setDuration(int v);
    void setEasingCurve(QEasingCurve::Type curve);

    int getAnimationDuration() const;
    void setAnimationDuration(int animationDuration);

    QEasingCurve::Type getAnimationCurve() const;
    void setAnimationCurve(const QEasingCurve::Type &animationCurve);

    bool getAutoResize() const;
    void setAutoResize(bool autoResize);

signals:
    void spacingItemAdded();
    void sizeChanged();

private:
    bool event(QEvent *e);
    void mouseMoveEvent(QMouseEvent *event);
    void dragEnterEvent(QDragEnterEvent *event);
    void dragLeaveEvent(QDragLeaveEvent *event);
    void dragMoveEvent(QDragMoveEvent *event);
    void dropEvent(QDropEvent *event);

private:
    void handleDrag(const QPoint &pos);
    void updateCurrentHoverInfo(int index, const QPoint &pos);
    void addSpacingItem(QWidget *souce, MoveDirection md, const QSize &size);
    void insertSpacingItemToLayout(int index, const QSize &size);
    int getHoverIndextByPos(const QPoint &pos);
    MoveDirection getVMoveDirection(int index, const QPoint &pos);
    MoveDirection getHMoveDirection(int index, const QPoint &pos);

private:
    int m_lastHoverIndex;
    bool m_hoverToSpacing;
    bool m_autoResize;
    QHBoxLayout *m_layout;
    QWidget *m_draginItem;
    QList<QWidget *> m_widgetList;
    QSize m_defaultSpacingItemSize;
    MoveDirection m_vMoveDirection;
    MoveDirection m_hMoveDirection;
    int m_animationDuration;
    QEasingCurve::Type m_animationCurve;
};

#endif // MOVABLELAYOUT_H
