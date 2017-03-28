#include "previewwidget.h"

#include <X11/Xlib.h>
#include <X11/X.h>
#include <X11/Xutil.h>

#include <QX11Info>

#include <QPainter>
#include <QTimer>

#define W 250
#define H 200

PreviewWidget::PreviewWidget(const WId wid, QWidget *parent)
    : QWidget(parent),

      m_wid(wid)
{
    setFixedSize(W, H);

    QTimer::singleShot(1, this, &PreviewWidget::refershImage);
}

void PreviewWidget::refershImage()
{
    XWindowAttributes attrs;
    XGetWindowAttributes(QX11Info::display(), m_wid, &attrs);
    XImage *ximage = XGetImage(QX11Info::display(), m_wid, 0, 0, attrs.width, attrs.height, AllPlanes, ZPixmap);
    const QImage qimage((uchar*)(ximage->data), attrs.width, attrs.height, QImage::Format_ARGB32);
    XDestroyImage(ximage);

    m_pixmap = QPixmap::fromImage(qimage).scaled(W, H, Qt::KeepAspectRatio, Qt::SmoothTransformation);

    update();
}

void PreviewWidget::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_pixmap.rect().center(), m_pixmap);
}

void PreviewWidget::mouseReleaseEvent(QMouseEvent *e)
{
    QWidget::mouseReleaseEvent(e);

    emit requestActivateWindow(m_wid);
}
