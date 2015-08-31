#include <QApplication>
#include <QX11Info>
#include <QDebug>
#include <QPainter>
#include <QPaintEvent>
#include <QFile>
#include <QByteArray>
#include <QAbstractNativeEventFilter>

#include <cairo/cairo.h>
#include <cairo/cairo-xlib.h>

#include <xcb/xcb.h>
#include <xcb/damage.h>

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

class Monitor : public QAbstractNativeEventFilter
{
public:
    Monitor(WindowPreview * wp) :
        QAbstractNativeEventFilter(),
        m_wp(wp)
    {
        xcb_connection_t *c = QX11Info::connection();
        xcb_prefetch_extension_data(c, &xcb_damage_id);
        const auto *reply = xcb_get_extension_data(c, &xcb_damage_id);
        m_damageEventBase = reply->first_event;
        if (reply->present) {
            xcb_damage_query_version_unchecked(c, XCB_DAMAGE_MAJOR_VERSION, XCB_DAMAGE_MINOR_VERSION);
        }

        m_damage = xcb_generate_id(c);
        xcb_damage_create(c, m_damage, wp->m_sourceWindow, XCB_DAMAGE_REPORT_LEVEL_RAW_RECTANGLES);
    }

    ~Monitor ()
    {
        xcb_connection_t *c = QX11Info::connection();
        xcb_damage_destroy(c, m_damage);
    }

    bool nativeEventFilter(const QByteArray &eventType, void *message, long *)
    {
        if (eventType=="xcb_generic_event_t") {
            xcb_generic_event_t *event = static_cast<xcb_generic_event_t*>(message);
            const uint8_t responseType = event->response_type & ~0x80;

            if (responseType == m_damageEventBase + XCB_DAMAGE_NOTIFY) {
                auto *ev = reinterpret_cast<xcb_damage_notify_event_t*>(event);

                if (m_wp && m_wp->m_sourceWindow == ev->drawable) {
                    m_wp->prepareRepaint();
                }
            } else if (responseType == XCB_CONFIGURE_NOTIFY) {
                auto *ev = reinterpret_cast<xcb_configure_notify_event_t*>(event);

                if (m_wp && m_wp->m_sourceWindow == ev->window) {
                    m_wp->prepareRepaint();
                }
            }
        }
        return false;
    }

private:
    WindowPreview * m_wp;
    int m_damageEventBase;
    int m_damage;
};

WindowPreview::WindowPreview(WId sourceWindow, QWidget *parent)
    : QFrame(parent),
      m_sourceWindow(sourceWindow),
      m_monitor(NULL)
{
    setObjectName("WindowPreview");

    setAttribute(Qt::WA_TransparentForMouseEvents);

    prepareRepaint();

    installMonitor();
}

WindowPreview::~WindowPreview()
{
    removeMonitor();
}

void WindowPreview::paintEvent(QPaintEvent *)
{
    QPainter painter;
    painter.begin(this);

    QImage image = QImage::fromData(imageData, "PNG");
    QPixmap pixmap = QPixmap::fromImage(image);
    //ignore border
    QRect rec(m_borderWidth, m_borderWidth, width() - m_borderWidth * 2, height() - m_borderWidth * 2);
    pixmap = pixmap.scaled(rec.width(), rec.height(), Qt::IgnoreAspectRatio, Qt::SmoothTransformation);
    painter.drawPixmap(rec.x(), rec.y(), pixmap);

    painter.end();
}
bool WindowPreview::isHover() const
{
    return m_isHover;
}

void WindowPreview::setIsHover(bool isHover)
{
    m_isHover = isHover;

    style()->unpolish(this);
    style()->polish(this);// force a stylesheet recomputation

    repaint();
}


void WindowPreview::installMonitor()
{
    if (!m_monitor) {
        m_monitor = new Monitor(this);

        QCoreApplication * app = QApplication::instance();
        if (app) app->installNativeEventFilter(m_monitor);
    }
}

void WindowPreview::removeMonitor()
{
    if (m_monitor) {
        QCoreApplication * app = QApplication::instance();
        if (app) app->removeNativeEventFilter(m_monitor);

        delete m_monitor;
        m_monitor = NULL;
    }
}

void WindowPreview::prepareRepaint()
{
    Display *dsp = QX11Info::display();

    XWindowAttributes s_atts;
    Status ss = XGetWindowAttributes(dsp, m_sourceWindow, &s_atts);
    QSize contentSize(width() - m_borderWidth * 2, height() - m_borderWidth * 2);

    if (ss != 0) {
        cairo_surface_t *source = cairo_xlib_surface_create(dsp,
                                                            m_sourceWindow,
                                                            s_atts.visual,
                                                            s_atts.width,
                                                            s_atts.height);

        cairo_surface_t * image_surface = cairo_image_surface_create(CAIRO_FORMAT_ARGB32,
                                                                     contentSize.width(),
                                                                     contentSize.height());

        float ratio = 0.0f;
        if (s_atts.width > s_atts.height) {
            ratio = contentSize.width() * 1.0 / s_atts.width;
        } else {
            ratio = contentSize.height() * 1.0 / s_atts.height;
        }
        int x = (contentSize.width() - s_atts.width * ratio) / 2.0;
        int y = (contentSize.height() - s_atts.height * ratio) / 2.0;

        cairo_t * cairo = cairo_create(image_surface);
        cairo_scale(cairo, ratio, ratio);
        cairo_set_source_surface(cairo, source, x / ratio, y / ratio);
        cairo_paint(cairo);

        imageData.clear();
        cairo_surface_write_to_png_stream(image_surface, cairo_write_func, this);

        cairo_surface_destroy(source);
        cairo_surface_destroy(image_surface);
        cairo_destroy(cairo);

        this->repaint();
    }
}


