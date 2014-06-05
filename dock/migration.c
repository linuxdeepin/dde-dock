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

char* guess_app_id(long s_pid, const char* instance_name, const char* wmname, const char* wmclass, const char* icon_name);

#include <glib.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <gio/gdesktopappinfo.h>
#include <glib/gprintf.h>
#include <sys/stat.h>
#include <stdio.h>
#include <gtk/gtk.h>
#include <fcntl.h>


char* get_name_by_pid(int pid);

GKeyFile* load_app_config(const char* name);

char* get_basename_without_extend_name(char const* path);
gboolean is_deepin_icon(char const* icon_path);
char* check_absolute_path_icon(char const* app_id, char const* icon_path);
gboolean is_chrome_app(char const* name);



enum APPID_FINDER_FILTER {
    APPID_FILTER_ARGS=1,
    APPID_FILTER_WMCLASS,
    APPID_FILTER_WMINSTANCE,
    APPID_FILTER_WMNAME,
    APPID_FILTER_ICON_NAME,
    APPID_FILTER_EXEC_NAME,
};

enum APPID_ICON_OPERATOR {
    ICON_OPERATOR_USE_ICONNAME=0,
    ICON_OPERATOR_USE_RUNTIME_WITH_BOARD,
    ICON_OPERATOR_USE_RUNTIME_WITHOUT_BOARD,
    ICON_OPERATOR_USE_PATH,
    ICON_OPERATOR_SET_DOMINANTCOLOR
};

char* find_app_id(const char* exec_name, const char* key, int filter);
void get_pid_info(int pid, char** exec_name, char** exec_args);
gboolean is_app_in_white_list(const char* name);

gboolean is_deepin_app_id(const char* app_id);
int get_deepin_app_id_operator(const char* app_id);
char* get_deepin_app_id_value(const char* app_id);

/*MEMORY_TESTED*/

#define DATA_DIR "/usr/share/dde/data"
#define FILTER_ARGS_PATH DATA_DIR"/filter_arg.ini"
#define FILTER_WMNAME_PATH DATA_DIR"/filter_wmname.ini"
#define FILTER_WMCLASS_PATH DATA_DIR"/filter_wmclass.ini"
#define FILTER_WMINSTANCE_PATH DATA_DIR"/filter_wminstance.ini"
#define FILTER_ICON_NAME_PATH DATA_DIR"/filter_icon_name.ini"
#define FILTER_EXEC_NAME_PATH DATA_DIR"/filter_execname.ini"
#define PROCESS_REGEX_PATH DATA_DIR"/process_regex.ini"
#define DEEPIN_ICONS_PATH DATA_DIR"/deepin_icons.ini"
#define FILTER_FILE "dock/filter.ini"

static GKeyFile* filter_args = NULL;
static GKeyFile* filter_wmname = NULL;
static GKeyFile* filter_wmclass = NULL;
static GKeyFile* filter_wminstance = NULL;
static GKeyFile* filter_icon_name = NULL;
static GKeyFile* filter_exec_name = NULL;
static GKeyFile* deepin_icons = NULL;

static GRegex* prefix_regex = NULL;
static GRegex* suffix_regex = NULL;
static GHashTable* white_apps = NULL;
static gboolean _is_init = FALSE;

static
void _build_filter_info(GKeyFile* filter, const char* path)
{
    if (g_key_file_load_from_file(filter, path, G_KEY_FILE_NONE, NULL)) {
        gsize size;
        char** groups = g_key_file_get_groups(filter, &size);
        gsize i=0;
        for (; i<size; i++) {
            gsize key_len;
            char** keys = g_key_file_get_keys(filter, groups[i], &key_len, NULL);
            gsize j=0;
            for (; j<key_len; j++) {
                g_hash_table_insert(white_apps, g_key_file_get_string(filter, groups[i], keys[j], NULL), NULL);
            }
        }
    }
}


static
void _init()
{
    white_apps = g_hash_table_new_full(g_str_hash, g_str_equal, g_free, NULL);

    // load and build process regex config information
    GKeyFile* process_regex = g_key_file_new();
    if (g_key_file_load_from_file(process_regex, PROCESS_REGEX_PATH, G_KEY_FILE_NONE, NULL)) {
        char* str = g_key_file_get_string(process_regex, "DEEPIN_PREFIX", "skip_prefix", NULL);
        prefix_regex = g_regex_new(str, G_REGEX_OPTIMIZE, 0, NULL);
        g_free(str);

        str = g_key_file_get_string(process_regex, "DEEPIN_PREFIX", "skip_suffix", NULL);
        suffix_regex = g_regex_new(str, G_REGEX_OPTIMIZE, 0, NULL);
        g_free(str);
    }
    if (prefix_regex == NULL) {
        g_warning("Can't build prefix_regex, use fallback config!");
        prefix_regex = g_regex_new(
                "(^gksu(do)?$)|(^sudo$)|(^mono$)|(^ruby$)|(^padsp$)|(^aoss$)|(^python(\\d.\\d)?$)|(^(ba)?sh$)",
                G_REGEX_OPTIMIZE, 0, NULL
                );

    }
    if (suffix_regex == NULL) {
        g_warning("Can't build suffix_regex, use fallback config!");
        suffix_regex = g_regex_new( "((-|.)bin$)|(.py$)", G_REGEX_OPTIMIZE, 0, NULL);
    }
    g_key_file_free(process_regex);

    // load filters and build white_list
    _build_filter_info(filter_args = g_key_file_new(), FILTER_ARGS_PATH);
    _build_filter_info(filter_wmclass = g_key_file_new(), FILTER_WMCLASS_PATH);
    _build_filter_info(filter_wminstance = g_key_file_new(), FILTER_WMINSTANCE_PATH);
    _build_filter_info(filter_wmname = g_key_file_new(), FILTER_WMNAME_PATH);
    _build_filter_info(filter_icon_name = g_key_file_new(), FILTER_ICON_NAME_PATH);
    _build_filter_info(filter_exec_name = g_key_file_new(), FILTER_EXEC_NAME_PATH);

    // set init flag
    _is_init = TRUE;
    g_assert(suffix_regex != NULL);
    g_assert(prefix_regex != NULL);
}


static
void _get_exec_name_args(char** cmdline, gsize length, char** name, char** args)
{
    g_assert(length != 0);
    *args = NULL;

    gsize name_pos = 0;

    if (cmdline[0] != NULL) {
        char* space_pos = NULL;
        if ((space_pos = strchr(cmdline[0], ' ')) != NULL && g_strrstr(cmdline[0], "chrom") != NULL) {
            *space_pos = '\0';
            gsize i= length -1;
            for (; i > 0; --i) {
                cmdline[i + 1] = cmdline[i];
            }
            length += 1;
            cmdline[1] = space_pos + 1;
        }
        char* basename = g_path_get_basename(cmdline[0]);
        if (g_regex_match(prefix_regex, basename, 0, NULL)) {
            g_debug("prefix match");
            while (cmdline[name_pos + 1] && cmdline[++name_pos][0] == '-') {
                g_debug("name pos changed");
            }
        }
        g_free(basename);
    }

    cmdline[length] = NULL;

    int diff = length - name_pos;
    if (diff == 0) {
        *name = g_path_get_basename(cmdline[0]);
        if (length > 1) {
            *args = g_strjoinv(" ", cmdline+1);
        }
    } else if (diff >= 1){
        *name = g_path_get_basename(cmdline[name_pos]);
        if (diff >= 2)
            *args = g_strjoinv(" ", cmdline + name_pos + 1);
    }

    char* tmp = *name;
    g_assert(tmp != NULL);
    g_assert(suffix_regex != NULL);
    *name = g_regex_replace_literal (suffix_regex, tmp, -1, 0, "", 0, NULL);
    g_free(tmp);

    guint i=0;
    for (; i<strlen(*name); i++) {
        if ((*name)[i] == ' ') {
            (*name)[i] = '\0';
            break;
        }
    }
}

static
char* _find_app_id_by_filter(const char* name, const char* keys_str, GKeyFile* filter)
{
    if (filter == NULL) return NULL;
    g_assert(name != NULL && keys_str != NULL);
    if (g_key_file_has_group(filter, name)) {
        gsize size = 0;
        char** keys = g_key_file_get_keys(filter, name, &size, NULL);
        gsize i=0;
        for (; i<size; i++) {
            if (g_strstr_len(keys_str , -1, keys[i])) {
                char* value = g_key_file_get_string(filter, name, keys[i], NULL);
                g_strfreev(keys);
                return value;
            }
        }
        g_strfreev(keys);
        /*g_debug("find \"%s\" in filter.ini but can't find the really desktop file\n", name);*/
    }
    return NULL;
}

char* find_app_id(const char* exec_name, const char* key, int filter)
{
    if (_is_init == FALSE) {
        _init();
    }
    g_assert(exec_name != NULL && key != NULL);
    // g_warning("exec_name: %s, key: %s", exec_name, key);
    switch (filter) {
        case APPID_FILTER_WMCLASS:
            return _find_app_id_by_filter(exec_name, key, filter_wmclass);
        case APPID_FILTER_WMNAME:
            return _find_app_id_by_filter(exec_name, key, filter_wmname);
        case APPID_FILTER_ARGS:
            return _find_app_id_by_filter(exec_name, key, filter_args);
        case APPID_FILTER_WMINSTANCE:
            return _find_app_id_by_filter(exec_name, key, filter_wminstance);
        case APPID_FILTER_ICON_NAME:
            return _find_app_id_by_filter(exec_name, key, filter_icon_name);
        case APPID_FILTER_EXEC_NAME: {
            char* id = _find_app_id_by_filter(exec_name, key, filter_exec_name);
            if (id == NULL)
                id = g_strdup(exec_name);
            return id;
        }
        default:
            g_error("filter %d is not support !", filter);
    }
    return NULL;
}

char* get_exe_name(int pid)
{
#define BUF_LEN 8095
    char buf[BUF_LEN] = {0};
    char* path = g_strdup_printf("/proc/%d/exe", pid);
    // header doesn't work, add this to avoid warning
    extern ssize_t readlink(const char*, char*, size_t);
    gsize len = readlink(path, buf, BUF_LEN);
    g_free(path);
    if (len > BUF_LEN) {
        g_warning("PID:%d's exe is to long!", pid);
        return NULL;
    }
#undef BUF_LEN
    return g_strdup(buf);
}
void get_pid_info(int pid, char** exec_name, char** exec_args)
{
    if (_is_init == FALSE) {
        _init();
    }
    char* cmd_line = NULL;
    char* path = g_strdup_printf("/proc/%d/cmdline", pid);

    gsize size=0;
    if (g_file_get_contents(path, &cmd_line, &size, NULL) && size > 0) {
        char** name_args = g_new(char*, 1024);
        gsize j = 0;
        name_args[j] = cmd_line;
        gsize i=1;
        for (; i<size && j<1024; i++) {
            if (cmd_line[i] == 0) {
                name_args[++j] = cmd_line + i + 1;
            }
        }
        name_args[j ? : j+1] = NULL;

        _get_exec_name_args(name_args, j+1, exec_name, exec_args);

        g_free(name_args);

    } else {
        *exec_name = get_exe_name(pid);
        *exec_args = NULL;
    }
    g_free(path);
    g_free(cmd_line);
}

gboolean is_app_in_white_list(const char* name)
{
    if (!_is_init) {
        _init();
    }
    return is_chrome_app(name) || g_hash_table_contains(white_apps, name);
}


gboolean is_deepin_app_id(const char* app_id)
{
    if (deepin_icons == NULL) {
        deepin_icons = g_key_file_new();
        if (!g_key_file_load_from_file(deepin_icons, DEEPIN_ICONS_PATH, G_KEY_FILE_NONE, NULL)) {
            g_key_file_free(deepin_icons);
            deepin_icons = NULL;
            return FALSE;
        }
    }
    return g_key_file_has_group(deepin_icons, app_id);

}

int get_deepin_app_id_operator(const char* app_id)
{
    g_assert(deepin_icons != NULL);
    return g_key_file_get_integer(deepin_icons, app_id, "operator", NULL);
}

char* get_deepin_app_id_value(const char* app_id)
{
    g_assert(deepin_icons != NULL);
    return g_key_file_get_string(deepin_icons, app_id, "value", NULL);
}





char* guess_app_id(long s_pid, const char* instance_name, const char* wmname, const char* wmclass, const char* icon_name)
{
    // g_setenv("G_MESSAGES_DEBUG", "all", FALSE);
    if (s_pid == 0) return g_strdup(wmclass);
    char* app_id = NULL;

    char* exec_name = NULL;
    char* exec_args = NULL;
    get_pid_info(s_pid, &exec_name, &exec_args);
    if (exec_name != NULL) {
        if (g_str_has_prefix(exec_name, "google-chrome-") || g_strcmp0(exec_name, "chrome") == 0) {
            g_free(exec_name);
            exec_name = g_strdup("google-chrome");
        }
        if (app_id == NULL) {
            GKeyFile* f = load_app_config(FILTER_FILE);
            if (f != NULL && wmname != NULL) {
                app_id = g_key_file_get_string(f, wmname, "appid", NULL);
            }
            g_key_file_unref(f);
            g_debug("[%s] get app id from StartupWMClass filter: %s", __func__, app_id);
        }
        if (app_id == NULL) {
            app_id = find_app_id(exec_name, wmname, APPID_FILTER_WMNAME);
            g_debug("[%s] get from wmname %s", __func__, app_id);
        }
        g_debug("exec_name:%s wmname:%s", exec_name, wmname);
        if (app_id == NULL && wmname != NULL) {
            app_id = find_app_id(exec_name, wmname, APPID_FILTER_WMINSTANCE);
            g_debug("[%s] get from wmname %s", __func__, app_id);
        }
        if (app_id == NULL && wmclass != NULL) {
            app_id = find_app_id(exec_name, wmclass, APPID_FILTER_WMCLASS);
            g_debug("[%s] get from wmclass %s", __func__, app_id);
        }
        if (app_id == NULL && exec_args != NULL && exec_args[0] != '\0') {
            app_id = find_app_id(exec_name, exec_args, APPID_FILTER_ARGS);
            g_debug("[%s] get app id from exec args(%s): %s", __func__, exec_args, app_id);
        }
        if (app_id == NULL && icon_name != NULL) {
            if (icon_name != NULL) {
                app_id = find_app_id(exec_name, icon_name, APPID_FILTER_ICON_NAME);
                g_debug("[%s] get from icon name %s", __func__, app_id);
            }
        }
        if (app_id == NULL && exec_name != NULL) {
            app_id = find_app_id(exec_name, exec_name, APPID_FILTER_EXEC_NAME);
            g_debug("[%s] get app id from exec name(%s): %s", __func__, exec_name, app_id);
        }
    } else {
        g_warning("[%s] exec_name get failed", __func__);
        app_id = g_strdup(wmclass);
    }
    g_free(exec_name);
    g_free(exec_args);
    return app_id;
}

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
    if (them == NULL) {
        g_warning("get theme failed");
        return NULL;
    }

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

void set_default_theme(const char* theme)
{
    GtkSettings* setting = gtk_settings_get_default();
    g_object_set(setting, "gtk-icon-theme-name", theme, NULL);
}

static char DE_NAME[100] = "DEEPIN";

void set_desktop_env_name(const char* name)
{
    size_t max_len = strlen(name) + 1;
    memcpy(DE_NAME, name, max_len > 100 ? max_len : 100);
    g_desktop_app_info_set_desktop_env(name);
}

void init_deepin()
{
    gtk_init(NULL, NULL);
    set_desktop_env_name("Deepin");
    set_default_theme("Deepin");
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

