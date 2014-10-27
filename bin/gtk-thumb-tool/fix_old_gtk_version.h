#include <gtk/gtk.h>
#ifndef __CREATE_FROM_PIXBUF__
#define __CREATE_FROM_PIXBUF__


#if !GTK_CHECK_VERSION(3, 10, 0)
cairo_surface_t* gdk_cairo_surface_create_from_pixbuf(GdkPixbuf* pixbuf, int scale, GdkWindow* w);
#endif

#if !GTK_CHECK_VERSION(3, 8, 0)
void gtk_widget_set_opacity(GtkWidget* widget, double opacity);
#endif

#endif
