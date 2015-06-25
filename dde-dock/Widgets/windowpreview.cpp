#include <QImage>
#include <QApplication>
#include <QtX11Extras/QX11Info>
#include <QPaintEvent>
#include <QPainter>
#include <QTimer>
#include <QDebug>

#include <X11/Xlib.h>
#include <X11/Xutil.h>

#include "windowpreview.h"

WindowPreview::WindowPreview(WId sourceWindow, QWidget *parent)
    : QWidget(parent),
      m_sourceWindow(sourceWindow),
      m_cache(NULL),
      m_timer(new QTimer(this))
{
    m_timer->setInterval(500);
    m_timer->setSingleShot(false);
    m_timer->start();
    connect(m_timer, &QTimer::timeout, [=]{ this->updateCache(); this->repaint(); });
}

WindowPreview::~WindowPreview()
{
    clearCache();
}

void WindowPreview::paintEvent(QPaintEvent * event)
{
    if (m_cache) {
        QPainter painter(this);

        QRect rect = m_cache->rect();
        rect.moveCenter(event->rect().center());

        painter.drawImage(rect, *m_cache);
        painter.end();
    }
}

void WindowPreview::clearCache()
{
    if (m_cache) {
        delete m_cache;
        m_cache = NULL;
    }
}

void WindowPreview::updateCache()
{
    clearCache();

    Display *dpy = QX11Info::display();

    XWindowAttributes watts;
    Status status = XGetWindowAttributes(dpy, m_sourceWindow, &watts);

    if (status != 0) {
        XImage *image = XGetImage(dpy, m_sourceWindow,
                                  watts.x, watts.y, watts.width, watts.height,
                                  AllPlanes, ZPixmap);
        if (image) {
            QImage cache(watts.width, watts.height, QImage::Format_RGB32);

            for (int y = 0; y < watts.height; y++) {
                for (int x = 0; x < watts.width; x++) {
                    u_long pixel = XGetPixel(image, x, y);

                    cache.setPixel(x, y, pixel);
                }
            }

            cache = cache.scaledToWidth(width(), Qt::SmoothTransformation);

            m_cache = new QImage(cache);
        }
    }
}
