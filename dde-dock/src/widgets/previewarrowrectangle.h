#ifndef PREVIEWARROWRECTANGLE_H
#define PREVIEWARROWRECTANGLE_H

#include <QTimer>
#include <QVBoxLayout>
#include <QPropertyAnimation>

#include "arrowrectangle.h"

class MirrorLabel : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(QSize size READ size WRITE setFixedSize)
public:
    MirrorLabel(QWidget *parent = 0);
    void storagePixmap(QPixmap pixmap);
    void setFixedSize(QSize size);

private:
    QPixmap m_pixmap;
};

class MirrorFrame : public QLabel
{
    Q_OBJECT
public:
    MirrorFrame(QWidget *parent = 0);

    void showWithAnimation(QPixmap mirror, QRect rect);
    void hideWithAnimation(QPixmap mirror, QRect rect);

signals:
    void showFinish();
    void hideFinish();
    void needStopShow();

private:
    MirrorLabel *m_mirrorLabel = NULL;

    const int SHOW_ANIMATION_DURATION = 200;
    const int HIDE_ANIMATION_DURATION = 200;
    const QEasingCurve SHOW_EASING_CURVE = QEasingCurve::InCirc;
    const QEasingCurve HIDE_EASING_CURVE = QEasingCurve::InCirc;
};

class PreviewArrowRectangle : public ArrowRectangle
{
    Q_OBJECT
public:
    PreviewArrowRectangle(QWidget *parent = 0);

    void showPreview(ArrowDirection direction, int x, int y, int interval = 800);
    void hidePreview(int interval);
    void cancelHide();
    void cancelShow();

signals:
    void needStopShow();
    void hideFinish();

protected:
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);

private:
    void initDelayHideTimer();
    void initDelayShowTimer();
    void initShowFrame();
    void initHideFrame();

    void doHide();
    void doShow();

private:
    QTimer *m_delayHideTimer = NULL;
    QTimer *m_delayShowTImer = NULL;

    MirrorFrame *m_showFrame = NULL;
    MirrorFrame *m_hideFrame = NULL;
    ArrowDirection m_lastArrowDirection = ArrowBottom;
    int m_lastX = 0;
    int m_lastY = 0;
};

#endif // PREVIEWARROWRECTANGLE_H
