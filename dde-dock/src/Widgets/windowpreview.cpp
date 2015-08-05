#include <QApplication>
#include <QTimer>
#include <QX11Info>
#include <QDebug>
#include <QPainter>
#include <QPaintEvent>
#include <QFile>
#include <QByteArray>


#include <cairo/cairo.h>
#include <cairo/cairo-xlib.h>

#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xatom.h>

#include "windowpreview.h"

static cairo_status_t cairo_write_func (void *widget, const unsigned char *data, unsigned int length)
{
    WindowPreview * wp = (WindowPreview *)widget;

    wp->imageData.append((const char *) data, length);

    return CAIRO_STATUS_SUCCESS;
}

WindowPreview::WindowPreview(WId sourceWindow, QWidget *parent)
    : QWidget(parent),
      m_sourceWindow(sourceWindow)
{
    setAttribute(Qt::WA_TransparentForMouseEvents);

    QTimer *timer = new QTimer(this);
    timer->setInterval(60);
    timer->start();
    connect(timer, &QTimer::timeout, this, &WindowPreview::onTimeout);
}

WindowPreview::~WindowPreview()
{

}

void WindowPreview::paintEvent(QPaintEvent *)
{
    qDebug() << "paintEvent of WindowPreview.";

    QPainter painter;
    painter.begin(this);

    QImage image = QImage::fromData(imageData, "PNG");
    QPixmap pixmap = QPixmap::fromImage(image);
    pixmap = pixmap.scaled(this->size(), Qt::KeepAspectRatio, Qt::SmoothTransformation);
    painter.drawPixmap(0, 0, pixmap);

    painter.end();
}

void WindowPreview::onTimeout()
{
    Display *dsp = QX11Info::display();

    XWindowAttributes s_atts;
    Status ss = XGetWindowAttributes(dsp, m_sourceWindow, &s_atts);

    if (ss != 0) {
        cairo_surface_t *source = cairo_xlib_surface_create(dsp,
                                                            m_sourceWindow,
                                                            s_atts.visual,
                                                            s_atts.width,
                                                            s_atts.height);

        cairo_surface_t * image_surface = cairo_surface_map_to_image(source, NULL);

        imageData.clear();
        cairo_surface_write_to_png_stream(image_surface, cairo_write_func, this);

        cairo_surface_unmap_image(source, image_surface);
        cairo_surface_destroy(source);

        this->repaint();
    }
}
