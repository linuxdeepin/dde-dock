/**
 * Copyright (c) 2011 ~ 2012 Deepin, Inc.
 *               2011 ~ 2012 snyh
 *
 * Author:      snyh <snyh@snyh.org>
 * Maintainer:  snyh <snyh@snyh.org>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/
#include "utils.h"
#include <glib.h>
#include <glib/gprintf.h>
#include <sys/stat.h>
#include <stdio.h>
#include <gtk/gtk.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>


char* shell_escape(const char* source)
{
    const unsigned char *p;
    char *dest;
    char *q;

    g_return_val_if_fail (source != NULL, NULL);

    p = (unsigned char *) source;
    q = dest = g_malloc (strlen (source) * 4 + 1);

    while (*p)
    {
        switch (*p)
        {
            case '\'':
                *q++ = '\\';
                *q++ = '\'';
                break;
            case '\\':
                *q++ = '\\';
                *q++ = '\\';
                break;
            case ' ':
                *q++ = '\\';
                *q++ = ' ';
                break;

            default:
                *q++ = *p;
        }
        p++;
    }
    *q = 0;
    return dest;
}

#include <unistd.h>
#include <fcntl.h>
char* get_name_by_pid(int pid)
{
#define LEN 1024
    char content[LEN];

    char* path = g_strdup_printf("/proc/%d/cmdline", pid);
    int fd = open(path, O_RDONLY);
    g_free(path);

    if (fd == -1) {
        return NULL;
    } else {
        int dump __attribute__((unused)) = read(fd, content, LEN);
        close(fd);
    }
    int i= 0;
    for (; i<LEN; i++) {
        if (content[i] == ' ') {
            content[i] = '\0';
            break;
        }
    }


    return g_path_get_basename(content);
}


GKeyFile* load_app_config(const char* name)
{
    char* path = g_build_filename(g_get_user_config_dir(), name, NULL);
    GKeyFile* key = g_key_file_new();
    g_key_file_load_from_file(key, path, G_KEY_FILE_NONE, NULL);
    g_free(path);
    /* no need to test file exitstly */
    return key;
}

void save_key_file(GKeyFile* key, const char* path)
{
    gsize size;
    gchar* content = g_key_file_to_data(key, &size, NULL);
    write_to_file(path, content, size);
    g_free(content);
}

void save_app_config(GKeyFile* key, const char* name)
{
    char* path = g_build_filename(g_get_user_config_dir(), name, NULL);
    save_key_file(key, path);
    g_free(path);
}

gboolean write_to_file(const char* path, const char* content, size_t size/* if 0 will use strlen(content)*/)
{
    char* dir = g_path_get_dirname(path);
    if (g_file_test(dir, G_FILE_TEST_IS_REGULAR)) {
        g_free(dir);
        g_warning("write content to %s, but %s is not directory!!\n",
                path, dir);
        return FALSE;
    } else if (!g_file_test(dir, G_FILE_TEST_EXISTS)) {
        if (g_mkdir_with_parents(dir, 0755) == -1) {
            g_warning("write content to %s, but create %s is failed!!\n",
                    path, dir);
            return FALSE;
        }
    }
    g_free(dir);

    if (size == 0) {
        size = strlen(content);
    }
    FILE* f = fopen(path, "w");
    if (f != NULL) {
        fwrite(content, sizeof(char), size, f);
        fclose(f);
        return TRUE;
    } else {
        return FALSE;
    }
}
int close_std_stream()
{
    //close stdin, stdout, stderr
    //redirect them to /dev/null
    int fd;
    close(STDIN_FILENO);
    fd=open("/dev/null", O_RDWR);
    if(fd!=STDIN_FILENO)
        return 1;
    if(dup2(STDIN_FILENO, STDOUT_FILENO)!=STDOUT_FILENO)
        return 1;
    if(dup2(STDIN_FILENO, STDERR_FILENO)!=STDERR_FILENO)
        return 1;
    return 0;
}
// reparent to init process.
int reparent_to_init ()
{
    switch (fork())
    {
	case -1:
	    return EXIT_FAILURE;
	case 0:
	    return EXIT_SUCCESS;
	default:
	    _exit(EXIT_SUCCESS);
    }
}
static void _consolidate_cmd_line (int subargc, char*** subargv_ptr)
{
    NOUSED(subargc);
    NOUSED(subargv_ptr);
    //recursively consolidate
}
void parse_cmd_line (int* argc_ptr, char*** argv_ptr)
{
    char*** subargv_ptr = argv_ptr;
    int     subargc     = (*argc_ptr);

    gboolean should_reparent = TRUE;
    gboolean enable_debug = FALSE;
    int i=0;
    for (;i<(*argc_ptr);i++)
    {
	if(!g_strcmp0 ((*argv_ptr)[i], "-f"))
	{
	    should_reparent=FALSE;
            //(*argv_ptr)[i]=NULL;
	    //(*argc_ptr)--;
            continue;
	}
	if(!g_strcmp0 ((*argv_ptr)[i], "-d"))
	{
            enable_debug = TRUE;
	    should_reparent=FALSE;
            //(*argv_ptr)[i]=NULL;
	    //(*argc_ptr)--;
            continue;
	}
    }
    //uncomment previous comments
    //consolidate *argv, remove NULL slots.
    _consolidate_cmd_line(subargc, subargv_ptr);

    if (should_reparent)
    {
        //FIXME: shall we exit?
        if (close_std_stream())
            return;
	reparent_to_init();
    }
    if (enable_debug)
    {
	g_setenv("G_MESSAGES_DEBUG", "all", FALSE);
    }
}

char* to_lower_inplace(char* str)
{
    g_assert(str != NULL);
    size_t i = 0;
    for (; i<strlen(str); i++)
        str[i] = g_ascii_tolower(str[i]);
    return str;
}

char* get_desktop_file_basename(GDesktopAppInfo* file)
{
    const char* filename = g_desktop_app_info_get_filename(file);
    return g_path_get_basename(filename);
}

GDesktopAppInfo* guess_desktop_file(char const* app_id)
{
    char* basename = g_strconcat(app_id, ".desktop", NULL);
    GDesktopAppInfo* desktop_file = g_desktop_app_info_new(basename);
    g_free(basename);
    return desktop_file;
}


char* get_basename_without_extend_name(char const* path)
{
    g_assert(path!= NULL);
    char* basename = g_path_get_basename(path);
    char* ext_sep = strrchr(basename, '.');
    if (ext_sep != NULL) {
        char* basename_without_ext = g_strndup(basename, ext_sep - basename);
        g_free(basename);
        return basename_without_ext;
    }

    return basename;
}


gboolean is_deepin_icon(char const* icon_path)
{
    return g_str_has_prefix(icon_path, "/usr/share/icons/Deepin/");
}

char* icon_name_to_path(const char* name, int size)
{
    if (g_path_is_absolute(name))
        return g_strdup(name);
    g_return_val_if_fail(name != NULL, NULL);

    int pic_name_len = strlen(name);
    char* ext = strrchr(name, '.');
    if (ext != NULL) {
        if (g_ascii_strcasecmp(ext+1, "png") == 0 || g_ascii_strcasecmp(ext+1, "svg") == 0 || g_ascii_strcasecmp(ext+1, "jpg") == 0) {
            pic_name_len = ext - name;
            g_debug("desktop's Icon name should an absoulte path or an basename without extension");
        }
    }

    char* pic_name = g_strndup(name, pic_name_len);
    GtkIconTheme* them = gtk_icon_theme_get_default(); //do not ref or unref it

    // This info must not unref, owned by gtk !!!!!!!!!!!!!!!!!!!!!
    GtkIconInfo* info = gtk_icon_theme_lookup_icon(them, pic_name, size, GTK_ICON_LOOKUP_GENERIC_FALLBACK);
    g_free(pic_name);
    if (info) {
        char* path = g_strdup(gtk_icon_info_get_filename(info));
        g_object_unref(info);
        return path;
    } else {
        return NULL;
    }
}

static char* _check(char const* app_id)
{
    char* icon = NULL;
    char* temp_icon_name_holder = icon_name_to_path(app_id, 48);

    if (temp_icon_name_holder != NULL) {
        if (!g_str_has_prefix(temp_icon_name_holder, "data:image"))
            icon = temp_icon_name_holder;
        else
            g_free(temp_icon_name_holder);
    }

    return icon;
}


char* check_absolute_path_icon(char const* app_id, char const* icon_path)
{
    char* icon = NULL;
    if ((icon = _check(app_id)) == NULL) {
        char* basename = get_basename_without_extend_name(icon_path);
        if (basename != NULL) {
            if (g_strcmp0(app_id, basename) == 0
                || (icon = _check(basename)) == NULL)
                icon = g_strdup(icon_path);
            g_free(basename);
        }
    }

    return icon;
}


gboolean is_chrome_app(char const* name)
{
    return g_str_has_prefix(name, "chrome-");
}


char* bg_blur_pict_get_dest_path (const char* src_uri)
{
    g_debug ("[%s] bg_blur_pict_get_dest_path: src_uri=%s", __func__, src_uri);
    g_return_val_if_fail (src_uri != NULL, NULL);

    //1. calculate original picture md5
    GChecksum* checksum;
    checksum = g_checksum_new (G_CHECKSUM_MD5);
    g_checksum_update (checksum, (const guchar *) src_uri, strlen (src_uri));

    guint8 digest[16];
    gsize digest_len = sizeof (digest);
    g_checksum_get_digest (checksum, digest, &digest_len);
    g_assert (digest_len == 16);

    //2. build blurred picture path
    char* file;
    file = g_strconcat (g_checksum_get_string (checksum), ".png", NULL);
    g_checksum_free (checksum);
    char* path;
    path = g_build_filename (g_get_user_cache_dir (),
                    BG_BLUR_PICT_CACHE_DIR,
                    file,
                    NULL);
    g_free (file);

    return path;
}

