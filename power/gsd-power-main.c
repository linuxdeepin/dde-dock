#include "gsd-power-manager.h"
#include "config.h"

#include <stdlib.h>
#include <unistd.h>
#include <libintl.h>
#include <errno.h>
#include <locale.h>
#include <signal.h>
#include <fcntl.h>
#include <sys/wait.h>

#include <glib/gi18n.h>
#include <glib/gstdio.h>
#include <gtk/gtk.h>
#include <libnotify/notify.h>

static gboolean   replace      = FALSE;
static gboolean   debug        = FALSE;
static gboolean   do_timed_exit = FALSE;
static int        term_signal_pipe_fds[2];
static guint      name_id      = 0;
/*static GnomeSettingsManager *manager = NULL;*/

static GOptionEntry entries[] =
{
    {"debug", 0, 0, G_OPTION_ARG_NONE, &debug, N_("Enable debugging code"), NULL },
    { "replace", 'r', 0, G_OPTION_ARG_NONE, &replace, N_("Replace existing daemon"), NULL },
    { "timed-exit", 0, 0, G_OPTION_ARG_NONE, &do_timed_exit, N_("Exit after a time (for debugging)"), NULL },
    {NULL}
};


static void
parse_args (int *argc, char ***argv)
{
    GError *error;
    GOptionContext *context;

    gnome_settings_profile_start (NULL);


    context = g_option_context_new (NULL);

    g_option_context_add_main_entries (context, entries, NULL);
    g_option_context_add_group (context, gtk_get_option_group (FALSE));

    error = NULL;
    if (!g_option_context_parse (context, argc, argv, &error))
    {
        if (error != NULL)
        {
            g_warning ("%s", error->message);
            g_error_free (error);
        }
        else
        {
            g_warning ("Unable to initialize GTK+");
        }
        exit (EXIT_FAILURE);
    }

    g_option_context_free (context);

    gnome_settings_profile_end (NULL);

    if (debug)
        g_setenv ("G_MESSAGES_DEBUG", "all", FALSE);
}

int gsd_power_manager_main(int argc, char **argv)
{
    GsdPowerManager *manager = gsd_power_manager_new();

    GError *error = NULL;
    parse_args(&argc, &argv);
    gtk_init(&argc, &argv);
    notify_init("gsd-power-manager");
    XInitThreads();
    gsd_power_manager_start(manager, &error);
    return 0;
}

#ifdef MAIN
int main(int argc, char **argv)
{
    gsd_power_manager_main(argc, argv);
    gtk_main();
    return 0;
}
#endif
