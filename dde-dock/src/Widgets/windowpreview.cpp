#include <QApplication>
#include <QTimer>
#include <QtX11Extras/QX11Info>
#include <QDebug>

#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xatom.h>

#include "windowpreview.h"

WindowPreview::WindowPreview(WId sourceWindow, QWidget *parent)
    : QWidget(parent),
      m_sourceWindow(sourceWindow)
{
    Display *dsp = QX11Info::display();

    XWindowAttributes t_atts;
    XGetWindowAttributes(dsp, winId(), &t_atts);

    m_surface = cairo_xlib_surface_create(dsp,
                                          winId(),
                                          t_atts.visual,
                                          t_atts.width,
                                          t_atts.height);
    m_cairo = cairo_create(m_surface);

    QTimer *timer = new QTimer(this);
    timer->setInterval(60);
    timer->start();
    connect(timer, &QTimer::timeout, this, &WindowPreview::onTimeout);
}

WindowPreview::~WindowPreview()
{
    cairo_surface_destroy(m_surface);
    cairo_destroy(m_cairo);
}

void WindowPreview::onTimeout()
{
    Display *dsp = QX11Info::display();

    XWindowAttributes s_atts;
    Status ss = XGetWindowAttributes(dsp, m_sourceWindow, &s_atts);

    XWindowAttributes t_atts;
    Status ts = XGetWindowAttributes(dsp, winId(), &t_atts);

    if (ss != 0 && ts != 0) {
        cairo_surface_t *source = cairo_xlib_surface_create(dsp,
                                                            m_sourceWindow,
                                                            s_atts.visual,
                                                            s_atts.width,
                                                            s_atts.height);
        cairo_xlib_surface_set_size(source, s_atts.width, s_atts.height);
        cairo_xlib_surface_set_size(m_surface, t_atts.width, t_atts.height);

        // clear the target surface.
        cairo_set_source_rgb(m_cairo, 1, 1, 1);
        cairo_set_operator(m_cairo, CAIRO_OPERATOR_CLEAR);
        cairo_paint(m_cairo);
        cairo_set_operator(m_cairo, CAIRO_OPERATOR_OVER);

        // calculate the scale ratio
        float ratio = 0.0f;
        if (s_atts.width > s_atts.height) {
            ratio = t_atts.width * 1.0 / s_atts.width;
        } else {
            ratio = t_atts.height * 1.0 / s_atts.height;
        }
        int x = (t_atts.width - s_atts.width * ratio) / 2.0;
        int y = (t_atts.height - s_atts.height * ratio) / 2.0;

        cairo_save(m_cairo);
        cairo_scale(m_cairo, ratio, ratio);
        cairo_set_source_surface(m_cairo, source, x / ratio, y / ratio);
        cairo_paint(m_cairo);
        cairo_restore(m_cairo);

        cairo_surface_destroy(source);
    }
}
