#ifndef WINDOWPREVIEW_H
#define WINDOWPREVIEW_H

#include <QWidget>

#include <cairo/cairo.h>
#include <cairo/cairo-xlib.h>

class WindowPreview : public QWidget
{
    Q_OBJECT

public:
    WindowPreview(WId sourceWindow, QWidget *parent = 0);
    ~WindowPreview();

    void onTimeout();

private:
    WId m_sourceWindow;
    cairo_t *m_cairo;
    cairo_surface_t *m_surface;
};

#endif // WINDOWPREVIEW_H
