#include "previewarrowrectangle.h"

MirrorLabel::MirrorLabel(QWidget *parent) : QLabel(parent)
{

}

void MirrorLabel::storagePixmap(QPixmap pixmap)
{
    m_pixmap = pixmap;
    QLabel::setFixedSize(0, 0);
}

void MirrorLabel::setFixedSize(QSize size)
{
    setPixmap(m_pixmap.scaled(size));
    QLabel::setFixedSize(size);
}

MirrorFrame::MirrorFrame(QWidget *parent) : QLabel(parent)
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);
    mainLayout->setAlignment(Qt::AlignHCenter | Qt::AlignBottom);
    mainLayout->setMargin(0);
    m_mirrorLabel = new MirrorLabel(this);
    mainLayout->addStretch();
    mainLayout->addWidget(m_mirrorLabel);
    setWindowFlags(Qt::SplashScreen);
    setAttribute(Qt::WA_TranslucentBackground);
}

void MirrorFrame::showWithAnimation(QPixmap mirror, QRect rect)
{
    setFixedSize(rect.size());
    move(rect.x(), rect.y());

    m_mirrorLabel->storagePixmap(mirror);
    QPropertyAnimation *animation = new QPropertyAnimation(m_mirrorLabel, "size");
    animation->setDuration(SHOW_ANIMATION_DURATION);
    animation->setEasingCurve(SHOW_EASING_CURVE);
    animation->setEndValue(rect.size());
    connect(animation, &QPropertyAnimation::finished, [=]{
        emit showFinish();
        animation->deleteLater();
        hide();
    });
    connect(this, &MirrorFrame::needStopShow, animation, &QPropertyAnimation::stop);
    animation->start();
    show();
}

void MirrorFrame::hideWithAnimation(QPixmap mirror, QRect rect)
{
    setFixedSize(rect.size());
    move(rect.x(), rect.y());

    m_mirrorLabel->storagePixmap(mirror);
    QPropertyAnimation *animation = new QPropertyAnimation(m_mirrorLabel, "size");
    animation->setDuration(HIDE_ANIMATION_DURATION);
    animation->setEasingCurve(HIDE_EASING_CURVE);
    animation->setStartValue(rect.size());
    animation->setEndValue(QSize(width() / 4, height() / 3));
    connect(animation, &QPropertyAnimation::finished, [=]{
        emit hideFinish();
        animation->deleteLater();
        hide();
        setFixedSize(0,0);
    });
    animation->start();
    show();
}

PreviewArrowRectangle::PreviewArrowRectangle(QWidget *parent)
    :ArrowRectangle(parent)
{
    initShowFrame();
    initHideFrame();
    initDelayHideTimer();
    initDelayShowTimer();
}

void PreviewArrowRectangle::enterEvent(QEvent *)
{
    cancelHide();
    ArrowRectangle::show(m_lastArrowDirection, m_lastX, m_lastY);
}

void PreviewArrowRectangle::leaveEvent(QEvent *)
{
    cancelShow();
    doHide();
}

void PreviewArrowRectangle::showPreview(ArrowDirection direction, int x, int y, int interval)
{
    if (m_delayShowTImer->isActive())
        return;
    m_lastArrowDirection = direction;
    m_lastX = x;
    m_lastY = y;

    m_delayShowTImer->start(interval);
}

void PreviewArrowRectangle::hidePreview(int interval)
{
    cancelShow();
    if (m_delayHideTimer->isActive() || isHidden())
        return;
    m_delayHideTimer->start(interval);
}

void PreviewArrowRectangle::cancelHide()
{
    m_delayHideTimer->stop();
}

void PreviewArrowRectangle::cancelShow()
{
    m_delayShowTImer->stop();

    emit needStopShow();
}

void PreviewArrowRectangle::initDelayHideTimer()
{
    m_delayHideTimer = new QTimer(this);
    m_delayHideTimer->setSingleShot(true);
    connect(m_delayHideTimer, &QTimer::timeout, this, &PreviewArrowRectangle::doHide);
}

void PreviewArrowRectangle::initDelayShowTimer()
{
    m_delayShowTImer = new QTimer(this);
    m_delayShowTImer->setSingleShot(true);
    connect(m_delayShowTImer, &QTimer::timeout, this, &PreviewArrowRectangle::doShow);
}

void PreviewArrowRectangle::initHideFrame()
{
    m_hideFrame = new MirrorFrame();
}

void PreviewArrowRectangle::initShowFrame()
{
    m_showFrame = new MirrorFrame();
    connect(m_showFrame, &MirrorFrame::showFinish, [=]{
        ArrowRectangle::show(m_lastArrowDirection, m_lastX, m_lastY);
    });
    connect(this, &PreviewArrowRectangle::needStopShow,[=]{
        m_showFrame->hide();
        m_showFrame->needStopShow();
    });
}

void PreviewArrowRectangle::doHide()
{
    //update geometry
    ArrowRectangle::show(m_lastArrowDirection, m_lastX, m_lastY);
    ArrowRectangle::move(m_lastX, m_lastY);
    ArrowRectangle::hide();

    QPixmap widgetImg = this->grab();

    m_hideFrame->hideWithAnimation(widgetImg, geometry());
    ArrowRectangle::hide();

    emit hideFinish();
}

void PreviewArrowRectangle::doShow()
{
    //update geometry
    ArrowRectangle::show(m_lastArrowDirection, m_lastX, m_lastY);
    ArrowRectangle::move(m_lastX, m_lastY);
    ArrowRectangle::hide();

    QPixmap widgetImg = this->grab();

    m_showFrame->showWithAnimation(widgetImg, geometry());
}
