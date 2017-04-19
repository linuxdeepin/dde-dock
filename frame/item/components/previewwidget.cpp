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

    const QImage qimage((const uchar*)(ximage->data), ximage->width, ximage->height, ximage->bytes_per_line, QImage::Format_RGB32);
    m_image = qimage.scaled(W, H, Qt::KeepAspectRatioByExpanding, Qt::SmoothTransformation);
    m_image = m_image.copy((m_image.width() - W) / 2, (m_image.height() - H) / 2, W, H);
    XDestroyImage(ximage);

    update();
}

void PreviewWidget::paintEvent(QPaintEvent *e)
{
    const QRect r = rect().marginsRemoved(QMargins(M, M, M, M));

    QPainter painter(this);

    // draw image
    const QRect ir = m_image.rect();
    const QPoint offset = r.center() - ir.center();

    painter.fillRect(offset.x(), offset.y(), ir.width(), ir.height(), Qt::white);
    painter.drawImage(offset.x(), offset.y(), m_image);

    // bottom black background
    QRect bgr = r;
    bgr.setTop(bgr.bottom() - 30);
    painter.fillRect(bgr, QColor(0, 0, 0, 255 * 0.3));
    // bottom title
    painter.drawText(bgr, Qt::AlignCenter, m_title);

    // draw border
    if (m_hovered)
    {
        const QRect br = r.marginsAdded(QMargins(1, 1, 1, 1));
        QPen p;
        p.setBrush(Qt::white);
        p.setWidth(4);
        painter.setBrush(Qt::transparent);
        painter.setRenderHint(QPainter::Antialiasing);
        painter.drawRoundedRect(br, 3, 3);
    }

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
