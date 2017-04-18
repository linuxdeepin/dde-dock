#include "previewwidget.h"

#include <X11/Xlib.h>
#include <X11/X.h>
#include <X11/Xutil.h>

#include <QX11Info>
#include <QPainter>
#include <QTimer>
#include <QVBoxLayout>

#define W 200
#define H 130
#define M 8

PreviewWidget::PreviewWidget(const WId wid, QWidget *parent)
    : QWidget(parent),

      m_wid(wid),
      m_hovered(false)
{
    m_closeButton = new QPushButton;
    m_closeButton->setFixedSize(16, 16);
    m_closeButton->setText("x");
    m_closeButton->setVisible(false);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->setSpacing(0);
    centralLayout->setMargin(0);
    centralLayout->addWidget(m_closeButton);
    centralLayout->setAlignment(m_closeButton, Qt::AlignTop | Qt::AlignRight);

    setFixedSize(W + M * 2, H + M * 2);
    setLayout(centralLayout);

    QTimer::singleShot(1, this, &PreviewWidget::refershImage);
}

void PreviewWidget::setTitle(const QString &title)
{
    m_title = title;

    update();
}

void PreviewWidget::refershImage()
{
    XWindowAttributes attrs;
    XGetWindowAttributes(QX11Info::display(), m_wid, &attrs);
    XImage *ximage = XGetImage(QX11Info::display(), m_wid, 0, 0, attrs.width, attrs.height, AllPlanes, ZPixmap);

    const QImage qimage((const uchar*)(ximage->data), ximage->width, ximage->height, ximage->bytes_per_line, QImage::Format_ARGB32);
    m_image = qimage.scaled(W, H, Qt::KeepAspectRatio, Qt::SmoothTransformation);

    update();
    XDestroyImage(ximage);
}

void PreviewWidget::paintEvent(QPaintEvent *e)
{
    const QRect r = rect().marginsRemoved(QMargins(M, M, M, M));

    QPainter painter(this);
#ifdef QT_DEBUG
    painter.fillRect(r, Qt::red);
#endif

    // draw image
    const QRect ir = m_image.rect();
    const QPoint offset = r.center() - ir.center();

    painter.fillRect(offset.x(), offset.y(), ir.width(), ir.height(), Qt::white);
    painter.drawImage(offset.x(), offset.y(), m_image);

    QWidget::paintEvent(e);
}

void PreviewWidget::enterEvent(QEvent *e)
{
    m_hovered = true;
    m_closeButton->setVisible(true);

    update();

    QWidget::enterEvent(e);
}

void PreviewWidget::leaveEvent(QEvent *e)
{
    m_hovered = false;
    m_closeButton->setVisible(false);

    update();

    QWidget::leaveEvent(e);
}

void PreviewWidget::mouseReleaseEvent(QMouseEvent *e)
{
    QWidget::mouseReleaseEvent(e);

    emit requestActivateWindow(m_wid);
}
