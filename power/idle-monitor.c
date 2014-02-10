#include <X11/Xlib.h>
#include <X11/extensions/sync.h>

#include <gtk/gtk.h>
#include <gdk/gdkx.h>
#define GNOME_DESKTOP_USE_UNSTABLE_API
#include "gnome-idle-monitor.h"

#define IDLE_TIME 1000 /* 1 second */

GHashTable *monitors = NULL; /* key = device id, value = GnomeIdleMonitor */

static void
active_watch_func (GnomeIdleMonitor *monitor,
                   guint             id,
                   gpointer          user_data)
{
    GdkDevice *device;
    int device_id;

    g_object_get (monitor, "device", &device, NULL);
    device_id = gdk_x11_device_get_id (device);
    g_message ("Active watch func called for device %s (id: %d, watch id %d)",
               gdk_device_get_name (device),
               device_id,
               id);
    g_object_unref (device);
}

static void
ensure_active_watch (GnomeIdleMonitor *monitor)
{
    GdkDevice *device;
    guint watch_id;
    int device_id;

    g_object_get (monitor, "device", &device, NULL);
    device_id = gdk_x11_device_get_id (device);
    watch_id = gnome_idle_monitor_add_user_active_watch (monitor,
               active_watch_func,
               NULL,
               NULL);
    g_message ("Added active watch ID %d for device %s (%d)",
               watch_id,
               gdk_device_get_name (device),
               device_id);
}

static void
idle_watch_func (GnomeIdleMonitor      *monitor,
                 guint                  id,
                 gpointer               user_data)
{
    GdkDevice *device;
    int device_id;

    g_object_get (monitor, "device", &device, NULL);
    device_id = gdk_x11_device_get_id (device);
    g_message ("Idle watch func called for device %s (id: %d, watch id %d)",
               gdk_device_get_name (device),
               device_id,
               id);
    g_object_unref (device);

    /*ensure_active_watch (monitor);*/
}

static void
idle_watch_all_func(GnomeIdleMonitor *monitor,
                    guint id,
                    gpointer user_data)
{
    g_message("Idle watch func called for all devices\n\n");
}

static void
device_added_cb (GdkDeviceManager *manager,
                 GdkDevice        *device,
                 gpointer          user_data)
{
    GnomeIdleMonitor *monitor;
    guint watch_id;
    int device_id;

    device_id = gdk_x11_device_get_id (device);
    monitor = gnome_idle_monitor_new_for_device (device);
    if (!monitor)
    {
        g_warning ("Per-device idletime monitor not available");
        return;
    }

    watch_id = gnome_idle_monitor_add_idle_watch (monitor,
               IDLE_TIME,
               idle_watch_func,
               NULL,
               NULL);
    g_message ("Added idle watch ID %d for device %s (%d)",
               watch_id,
               gdk_device_get_name (device),
               device_id);

    /*ensure_active_watch (monitor);*/

    g_hash_table_insert (monitors,
                         GINT_TO_POINTER (device_id),
                         monitor);
}

static void
device_removed_cb (GdkDeviceManager *manager,
                   GdkDevice        *device,
                   gpointer          user_data)
{
    g_hash_table_remove (monitors,
                         GINT_TO_POINTER (gdk_x11_device_get_id (device)));
    g_message ("Removed watch for device %s (%d)",
               gdk_device_get_name (device),
               gdk_x11_device_get_id (device));
}

static void
device_changed_cb (GdkDeviceManager *manager,
                   GdkDevice        *device,
                   gpointer          user_data)
{
    if (gdk_device_get_device_type (device) == GDK_DEVICE_TYPE_FLOATING)
        device_removed_cb (manager, device, NULL);
    else
        device_added_cb (manager, device, NULL);
}

int main (int argc, char **argv)
{
    /*GdkDeviceManager *manager;*/
    /*GList *devices, *l;*/

    gtk_init (&argc, &argv);

    monitors = g_hash_table_new_full (g_direct_hash, g_direct_equal,
                                      NULL, g_object_unref);

    /*manager = gdk_display_get_device_manager (gdk_display_get_default ());*/
    /*g_signal_connect (manager, "device-added",*/
    /*G_CALLBACK (device_added_cb), NULL);*/
    /*g_signal_connect (manager, "device-removed",*/
    /*G_CALLBACK (device_removed_cb), NULL);*/
    /*g_signal_connect (manager, "device-changed",*/
    /*G_CALLBACK (device_changed_cb), NULL);*/
    /*devices = gdk_device_manager_list_devices (manager, GDK_DEVICE_TYPE_SLAVE);*/
    /*for (l = devices; l != NULL; l = l->next)*/
    /*{*/
    /*GdkDevice *device = l->data;*/

    /*device_added_cb (manager, device, NULL);*/
    /*}*/

    /*GnomeIdleMonitor *monitor = gnome_idle_monitor_new();*/
    GnomeIdleMonitor *monitor = gnome_idle_monitor_new();
    guint watch_id;
    if (!monitor)
    {
        g_warning("idle monitor not available");
        return -1;
    }

    watch_id = gnome_idle_monitor_add_idle_watch(monitor,
               3 * IDLE_TIME,
               idle_watch_all_func,
               NULL,
               NULL);
    g_message("Added idle watch ID %d for devices\n", watch_id);

    /*g_hash_table_insert(monitors,*/
    /*GINT_TO_POINTER(),*/
    /*monitor);*/
    gtk_main ();

    return 0;
}
