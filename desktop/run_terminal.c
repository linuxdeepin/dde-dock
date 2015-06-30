#include <glib.h>
#include <gio/gdesktopappinfo.h>
#include <gdk/gdk.h>

GAppInfo *gen_app_info (const char* dir, const char* executable, GAppInfoCreateFlags flags)
{
    GAppInfo *appinfo = NULL;
    GError* error = NULL;
    char* cmd_line = NULL;

    if (executable == NULL)
    {
        char* tmp1 = g_shell_quote(dir);
        char* tmp2 = g_strdup_printf("cd %s && exec $SHELL", tmp1);
        g_free(tmp1);
        tmp1 = g_shell_quote(tmp2);
        g_free(tmp2);

        cmd_line = g_strconcat("sh -c ", tmp1, NULL);
        g_free(tmp1);
    }
    else
    {
        cmd_line = g_strdup(executable);
    }

    appinfo = g_app_info_create_from_commandline(cmd_line, NULL,
            flags, &error);
    g_free(cmd_line);
    if (error!=NULL)
    {
        g_debug("gen_app_info error: %s", error->message);
        g_error_free(error);
        return NULL;
    }
    error = NULL;

    return appinfo;
}


gboolean exec_app_info (const char* dir, const char *executable, GAppInfoCreateFlags flags)
{
    GAppInfo *appinfo = NULL;
    GError *error = NULL;
    gboolean is_ok __attribute__((unused)) = FALSE;

    appinfo = gen_app_info (dir, executable, flags);
    if ( appinfo == NULL ) {
        g_debug ("gen app info failed!");
        return FALSE;
    }

    GdkAppLaunchContext* ctx = gdk_display_get_app_launch_context(gdk_display_get_default());
    gdk_app_launch_context_set_screen(ctx, gdk_screen_get_default()); //must set this otherwise termiator will not work properly
    is_ok = g_app_info_launch (appinfo, NULL, (GAppLaunchContext*)ctx, &error);
    g_object_unref(ctx);
    if (error!=NULL)
    {
        g_debug("exec app info error: %s", error->message);
        g_error_free(error);
        g_object_unref (appinfo);

        return FALSE;
    }

    g_object_unref (appinfo);
    return TRUE;
}


void run_in_terminal(char* dir, char* executable)
{
    gboolean is_ok = FALSE;
    gboolean should_free_dir = FALSE;
    if (dir == NULL) {
        should_free_dir = TRUE;
        dir = g_get_current_dir();
    }

    GSettings* s = g_settings_new("com.deepin.desktop.default-applications.terminal");
    if (s != NULL) {
        char* terminal = g_settings_get_string(s, "exec");
        g_object_unref(s);
        if (terminal != NULL && terminal[0] != '\0') {
            char* quoted_dir = g_shell_quote(dir);
            char* exec = NULL;
            exec = g_strdup_printf("sh -c 'cd %s && %s'", quoted_dir, terminal);
            g_free(quoted_dir);
            g_free(terminal);
            is_ok = exec_app_info(dir, exec, G_APP_INFO_CREATE_NONE);
            g_free(exec);
            if (!is_ok) {
                g_debug("exec app info failed!");
            }

            if (should_free_dir) {
                g_free(dir);
            }
            return;
        }
        g_free(terminal);
    }


    is_ok = exec_app_info(dir, executable, G_APP_INFO_CREATE_NEEDS_TERMINAL);
    if ( !is_ok ) {
        g_debug ("exec app info failed!");
        /*exec_app_info (executable);*/
    }

    if (should_free_dir) {
        g_free(dir);
    }
    return;
}
