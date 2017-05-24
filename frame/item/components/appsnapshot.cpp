#include "appsnapshot.h"
#include "_previewcontainer.h"

#include <X11/Xlib.h>
#include <X11/X.h>
#include <X11/Xutil.h>
#include <X11/Xatom.h>

#include <QX11Info>
#include <QPainter>

AppSnapshot::AppSnapshot(const WId wid, QWidget *parent)
    : QWidget(parent),

      m_wid(wid),

      m_fetchSnapshotTimer(new QTimer(this))
{
    m_fetchSnapshotTimer->setSingleShot(true);
    m_fetchSnapshotTimer->setInterval(10);

    connect(m_fetchSnapshotTimer, &QTimer::timeout, this, &AppSnapshot::fetchSnapshot, Qt::QueuedConnection);

    QTimer::singleShot(1, this, &AppSnapshot::fetchSnapshot);
}

void AppSnapshot::closeWindow() const
{
    const auto display = QX11Info::display();

    XEvent e;

    memset(&e, 0, sizeof(e));
    e.xclient.type = ClientMessage;
    e.xclient.window = m_wid;
    e.xclient.message_type = XInternAtom(display, "WM_PROTOCOLS", true);
    e.xclient.format = 32;
    e.xclient.data.l[0] = XInternAtom(display, "WM_DELETE_WINDOW", false);
    e.xclient.data.l[1] = CurrentTime;

    XSendEvent(display, m_wid, false, NoEventMask, &e);
    XFlush(display);
}

void AppSnapshot::setWindowTitle(const QString &title)
{
    m_title = title;
}

void AppSnapshot::fetchSnapshot()
{
    if (!isVisible())
        return;

    const auto display = QX11Info::display();

    XWindowAttributes attrs;
    XGetWindowAttributes(display, m_wid, &attrs);
    XImage *ximage = XGetImage(display, m_wid, 0, 0, attrs.width, attrs.height, AllPlanes, ZPixmap);
    const QImage qimage((const uchar*)(ximage->data), ximage->width, ximage->height, ximage->bytes_per_line, QImage::Format_RGB32);
    Q_ASSERT(!qimage.isNull());

    const Atom gtk_frame_extents = XInternAtom(display, "_GTK_FRAME_EXTENTS", true);
    Atom actual_type_return;
    int actual_format_return;
    unsigned long n_items_return;
    unsigned long bytes_after_return;
    unsigned char *prop_to_return;

    const auto r = XGetWindowProperty(display, m_wid, gtk_frame_extents, 0, 4, false, XA_CARDINAL,
                                      &actual_type_return, &actual_format_return, &n_items_return, &bytes_after_return, &prop_to_return);
    if (!r && prop_to_return && n_items_return == 4 && actual_format_return == 32)
    {
        const unsigned long *extents = reinterpret_cast<const unsigned long *>(prop_to_return);
        const int left = extents[0];
        const int right = extents[1];
        const int top = extents[2];
        const int bottom = extents[3];
        const int width = qimage.width();
        const int height = qimage.height();

        m_snapshot = qimage.copy(left, top, width - left - right, height - top - bottom);
    } else {
        m_snapshot = qimage.copy();
    }

//    const int w = width();
//    const int h = height();
//    m_snapshot = m_snapshot.scaled(w, h, Qt::KeepAspectRatioByExpanding, Qt::SmoothTransformation);
//    m_snapshot = m_snapshot.copy((m_snapshot.width() - w) / 2, (m_snapshot.height() - h) / 2, w, h);
    XDestroyImage(ximage);
    XFree(prop_to_return);

    update();
}

void AppSnapshot::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    emit entered(m_wid);
}

void AppSnapshot::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);

    if (m_snapshot.isNull())
        return;

    const QRect r = rect().marginsRemoved(QMargins(8, 8, 8, 8));

    // draw image
//    const QPoint offset = r.center() - ir.center();

//    painter.fillRect(offset.x(), offset.y(), ir.width(), ir.height(), Qt::white);
//    painter.drawImage(offset.x(), offset.y(), m_snapshot);
//    painter.fillRect(r, Qt::white);
    const QImage im = m_snapshot.scaled(r.size(), Qt::KeepAspectRatio, Qt::SmoothTransformation);
    const QRect ir = im.rect();
    const QPoint offset = r.center() - ir.center();
    painter.drawImage(offset.x(), offset.y(), im);
}

void AppSnapshot::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    m_fetchSnapshotTimer->start();
}

void AppSnapshot::mousePressEvent(QMouseEvent *e)
{
    QWidget::mousePressEvent(e);

    emit clicked(m_wid);
}
