/**
 * Copyright (C) 2012 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <string.h>
#include <gtk/gtk.h>

char *get_icon_file_path(const char*name)
{
    GtkIconTheme* theme = gtk_icon_theme_get_default();
    if (theme == NULL) {
        g_warning("error get default icon theme failed");
        return NULL;
    }
    const int size = 48;
    GtkIconInfo* info = gtk_icon_theme_lookup_icon(theme, name, size, GTK_ICON_LOOKUP_GENERIC_FALLBACK);
    if (info) {
        char* path = g_strdup(gtk_icon_info_get_filename(info));
        g_object_unref(info);
        return path;
    }
    return NULL;
}

static char DE_NAME[100] = "DEEPIN";

void set_desktop_env_name(const char* name)
{
    size_t max_len = strlen(name) + 1;
    memcpy(DE_NAME, name, max_len > 100 ? max_len : 100);
#if GTK_CHECK_VERSION(2, 42, 0)
    g_setenv("XDG_CURRENT_DESKTOP", name, TRUE);
#else
    g_desktop_app_info_set_desktop_env(name);
#endif
}


void init_deepin()
{
    gtk_init(NULL, NULL);
    set_desktop_env_name("Deepin");
}


char* get_data_uri_by_pixbuf(GdkPixbuf* pixbuf)
{
    gchar* buf = NULL;
    gsize size = 0;
    GError *error = NULL;

    gdk_pixbuf_save_to_buffer(pixbuf, &buf, &size, "png", &error, NULL);
    g_assert(buf != NULL);

    if (error != NULL) {
        g_warning("%s\n", error->message);
        g_error_free(error);
        g_free(buf);
        return NULL;
    }

    char* base64 = g_base64_encode((const guchar*)buf, size);
    g_free(buf);
    char* data = g_strconcat("data:image/png;base64,", base64, NULL);
    g_free(base64);

    return data;
}


char* get_data_uri_by_path(const char* path)
{
    GError *error = NULL;
    GdkPixbuf* pixbuf = gdk_pixbuf_new_from_file(path, &error);
    if (error != NULL) {
        g_warning("%s\n", error->message);
        g_error_free(error);
        return NULL;
    }
    char* c = get_data_uri_by_pixbuf(pixbuf);
    g_object_unref(pixbuf);
    return c;

}
