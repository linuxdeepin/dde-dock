#include <gio/gio.h>

#include "disk_listener.h"
#include "_cgo_export.h"

static void handle_disk_changed(char* ev, GVolume* volume);
static void volume_changed_cb(GVolumeMonitor* monitor, GVolume* volume, 
        gpointer data);
static void mount_changed_cb(GVolumeMonitor* monitor, GMount* mount, 
        gpointer data);

static void
handle_volume_changed(char* ev, GVolume* volume)
{
    char* uuid = g_volume_get_identifier(volume, "uuid");
    if (uuid) {
        handleDiskChanged(ev, uuid);
        return;
    }

    char* path = g_volume_get_identifier(volume, "unix-device");
    if (path) {
        handleDiskChanged(ev, path);
        return;
    }
}

static void
volume_changed_cb(GVolumeMonitor* monitor, GVolume* volume, gpointer data)
{
    char* ev = (char*)data;
    handle_volume_changed(ev, volume);
}

static void
mount_changed_cb(GVolumeMonitor* monitor, GMount* mount, gpointer data)
{
    char* ev = (char*)data;
    GVolume* volume = g_mount_get_volume(mount);
    if (!volume) {
        if (g_strcmp0(ev, "mount-removed") == 0) {
            handleDiskChanged(ev, "");
        }
        return;
    }

    handle_volume_changed(ev, volume);
    g_object_unref(G_OBJECT(volume));
}

void
start_disk_listener()
{
    static gboolean running = FALSE;
    if (running) {
        return;
    }

    GVolumeMonitor* monitor = g_volume_monitor_get();

    g_signal_connect(G_OBJECT(monitor), "volume-added", 
            G_CALLBACK(volume_changed_cb), "volume-added");
    g_signal_connect(G_OBJECT(monitor), "volume-removed", 
            G_CALLBACK(volume_changed_cb), "volume-removed");

    g_signal_connect(G_OBJECT(monitor), "mount-added", 
            G_CALLBACK(mount_changed_cb), "mount-added");
    g_signal_connect(G_OBJECT(monitor), "mount-removed", 
            G_CALLBACK(mount_changed_cb), "mount-removed");
    running = TRUE;
}
