#include "fix_old_gtk_version.h"

#if !GTK_CHECK_VERSION(3, 10, 0)

static void gdk_cairo_surface_paint_pixbuf (cairo_surface_t *surface, const GdkPixbuf *pixbuf);

cairo_surface_t* gdk_cairo_surface_create_from_pixbuf(GdkPixbuf* pixbuf, int scale, GdkWindow* w)
{

    cairo_surface_t* surface = cairo_image_surface_create(
            gdk_pixbuf_get_n_channels(pixbuf)==4 ?  CAIRO_FORMAT_ARGB32 : CAIRO_FORMAT_RGB24,
            gdk_pixbuf_get_width(pixbuf),
            gdk_pixbuf_get_height(pixbuf));
    gdk_cairo_surface_paint_pixbuf(surface, pixbuf);
    return surface;
}

//this static method is copied from gdk/gdkcairo.h
static void
gdk_cairo_surface_paint_pixbuf (cairo_surface_t *surface,
                                const GdkPixbuf *pixbuf)
{
  gint width, height;
  guchar *gdk_pixels, *cairo_pixels;
  int gdk_rowstride, cairo_stride;
  int n_channels;
  int j;

  /* This function can't just copy any pixbuf to any surface, be
   * sure to read the invariants here before calling it */

  g_assert (cairo_surface_get_type (surface) == CAIRO_SURFACE_TYPE_IMAGE);
  g_assert (cairo_image_surface_get_format (surface) == CAIRO_FORMAT_RGB24 ||
            cairo_image_surface_get_format (surface) == CAIRO_FORMAT_ARGB32);
  g_assert (cairo_image_surface_get_width (surface) == gdk_pixbuf_get_width (pixbuf));
  g_assert (cairo_image_surface_get_height (surface) == gdk_pixbuf_get_height (pixbuf));

  cairo_surface_flush (surface);

  width = gdk_pixbuf_get_width (pixbuf);
  height = gdk_pixbuf_get_height (pixbuf);
  gdk_pixels = gdk_pixbuf_get_pixels (pixbuf);
  gdk_rowstride = gdk_pixbuf_get_rowstride (pixbuf);
  n_channels = gdk_pixbuf_get_n_channels (pixbuf);
  cairo_stride = cairo_image_surface_get_stride (surface);
  cairo_pixels = cairo_image_surface_get_data (surface);

  for (j = height; j; j--)
    {
      guchar *p = gdk_pixels;
      guchar *q = cairo_pixels;

      if (n_channels == 3)
        {
          guchar *end = p + 3 * width;

          while (p < end)
            {
#if G_BYTE_ORDER == G_LITTLE_ENDIAN
              q[0] = p[2];
              q[1] = p[1];
              q[2] = p[0];
#else
              q[1] = p[0];
              q[2] = p[1];
              q[3] = p[2];
#endif
              p += 3;
              q += 4;
            }
        }
      else
        {
          guchar *end = p + 4 * width;
          guint t1,t2,t3;

#define MULT(d,c,a,t) G_STMT_START { t = c * a + 0x80; d = ((t >> 8) + t) >> 8; } G_STMT_END

          while (p < end)
            {
#if G_BYTE_ORDER == G_LITTLE_ENDIAN
              MULT(q[0], p[2], p[3], t1);
              MULT(q[1], p[1], p[3], t2);
              MULT(q[2], p[0], p[3], t3);
              q[3] = p[3];
#else
              q[0] = p[3];
              MULT(q[1], p[0], p[3], t1);
              MULT(q[2], p[1], p[3], t2);
              MULT(q[3], p[2], p[3], t3);
#endif

              p += 4;
              q += 4;
            }

#undef MULT
        }

      gdk_pixels += gdk_rowstride;
      cairo_pixels += cairo_stride;
    }

  cairo_surface_mark_dirty (surface);
}
#endif

#if !GTK_CHECK_VERSION(3, 8, 0)
void gtk_widget_set_opacity(GtkWidget* widget, double opacity)
{
    g_return_if_fail(GTK_IS_WINDOW(widget));
    gtk_window_set_opacity(GTK_WINDOW(widget), opacity);
}
#endif
