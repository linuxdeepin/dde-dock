#ifndef PREVIEWFRAME_H
#define PREVIEWFRAME_H

#include <QTimer>
#include <QPropertyAnimation>
#include "arrowrectangle.h"

class PreviewFrame : public ArrowRectangle
{
    Q_OBJECT
    Q_PROPERTY(QPoint arrowPos READ QPoint(0, 0) WRITE setArrowPos)
public:
    explicit PreviewFrame(QWidget *parent = 0);
    ~PreviewFrame();

    void showPreview(ArrowDirection direction, int x, int y, int interval);
    void hidePreview(int interval = 0);
    void setContent(QWidget *content);
    void setArrowPos(const QPoint &pos);

signals:
    void hideFinish();

protected:
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);

private:
    void onShowTimerTriggered();
    void onHideTimerTriggered();

private:
    QTimer *m_showTimer = NULL;
    QTimer *m_hideTimer = NULL;
    QWidget *m_tmpContent = NULL;
    QPropertyAnimation *m_animation = NULL;
    ArrowDirection m_direction = ArrowBottom;
    QPoint m_lastPos = QPoint(0, 0);
    int m_x = 0;
    int m_y = 0;
    const int MOVE_ANIMATION_DURATION = 300;
    const QEasingCurve MOVE_ANIMATION_CURVE = QEasingCurve::OutCirc;
};

#endif // PREVIEWFRAME_H
