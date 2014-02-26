/* -*- Mode: C; tab-width: 8; indent-tabs-mode: nil; c-basic-offset: 8 -*-
 *
 * Copyright (C) 2007 William Jon McCann <mccann@jhu.edu>
 * Copyright (C) 2011-2012 Richard Hughes <richard@hughsie.com>
 * Copyright (C) 2011 Ritesh Khadgaray <khadgaray@gmail.com>
 * Copyright (C) 2012-2013 Red Hat Inc.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA 02111-1307, USA.
 *
 */

/*#include "config.h"*/

#include <stdlib.h>
#include <string.h>
#include <glib/gi18n.h>
#include <gtk/gtk.h>
#define UPOWER_ENABLE_DEPRECATED 1
#include <libupower-glib/upower.h>
#include <libnotify/notify.h>
#include <canberra-gtk.h>
#include <glib-unix.h>
#include <gio/gunixfdlist.h>

#define GNOME_DESKTOP_USE_UNSTABLE_API
#include <libgnome-desktop/gnome-rr.h>
#include <libgnome-desktop/gnome-idle-monitor.h>

/*#include <gsd-input-helper.h>*/

#include "gsd-power-constants.h"
#include "gsm-inhibitor-flag.h"
#include "gsm-presence-flag.h"
#include "gsm-manager-logout-mode.h"
#include "gpm-common.h"
#include "gnome-settings-plugin.h"
/*#include "gnome-settings-profile.h"*/
/*#include "gnome-settings-session.h"*/
#include "gsd-enums.h"
#include "gsd-power-manager.h"

#define GNOME_SESSION_DBUS_NAME                 "org.gnome.SessionManager"
#define GNOME_SESSION_DBUS_PATH_PRESENCE        "/org/gnome/SessionManager/Presence"
#define GNOME_SESSION_DBUS_INTERFACE_PRESENCE   "org.gnome.SessionManager.Presence"

#define UPOWER_DBUS_NAME                        "org.freedesktop.UPower"
#define UPOWER_DBUS_PATH                        "/org/freedesktop/UPower"
#define UPOWER_DBUS_PATH_KBDBACKLIGHT           "/org/freedesktop/UPower/KbdBacklight"
#define UPOWER_DBUS_INTERFACE                   "org.freedesktop.UPower"
#define UPOWER_DBUS_INTERFACE_KBDBACKLIGHT      "org.freedesktop.UPower.KbdBacklight"

/*#define GSD_POWER_SETTINGS_SCHEMA             "org.gnome.settings-daemon.plugins.power"*/
#define DEEPIN_POWER_PROFILE_SCHEMA             "com.deepin.daemon.power"
#define DEEPIN_POWER_SETTINGS_SCHEMA            "com.deepin.daemon.power.settings"
#define DEEPIN_POWER_SETTINGS_PATH_PRE              "/com/deepin/daemon/power/profiles/"
#define GSD_XRANDR_SETTINGS_SCHEMA              "org.gnome.settings-daemon.plugins.xrandr"

#define GSD_POWER_DBUS_NAME                     GSD_DBUS_NAME ".Power"
#define GSD_POWER_DBUS_PATH                     GSD_DBUS_PATH "/Power"
#define GSD_POWER_DBUS_INTERFACE                GSD_DBUS_BASE_INTERFACE ".Power"
#define GSD_POWER_DBUS_INTERFACE_SCREEN         GSD_POWER_DBUS_INTERFACE ".Screen"
#define GSD_POWER_DBUS_INTERFACE_KEYBOARD       GSD_POWER_DBUS_INTERFACE ".Keyboard"

#define GS_DBUS_NAME                            "org.gnome.ScreenSaver"
#define GS_DBUS_PATH                            "/org/gnome/ScreenSaver"
#define GS_DBUS_INTERFACE                       "org.gnome.ScreenSaver"

#define GSD_POWER_MANAGER_NOTIFY_TIMEOUT_SHORT          10 * 1000 /* ms */
#define GSD_POWER_MANAGER_NOTIFY_TIMEOUT_LONG           30 * 1000 /* ms */

#define GSD_POWER_MANAGER_RECALL_DELAY                  30 /* seconds */
#define GSD_POWER_MANAGER_LID_CLOSE_SAFETY_TIMEOUT      30 /* seconds */

#define SYSTEMD_DBUS_NAME                       "org.freedesktop.login1"
#define SYSTEMD_DBUS_PATH                       "/org/freedesktop/login1"
#define SYSTEMD_DBUS_INTERFACE                  "org.freedesktop.login1.Manager"

/* Keep this in sync with gnome-shell */
#define SCREENSAVER_FADE_TIME                           10 /* seconds */

/* Time between notifying the user about a critical action and executing it.
 * This can be changed with the GSD_ACTION_DELAY constant. */
#ifndef GSD_ACTION_DELAY
#define GSD_ACTION_DELAY 20
#endif /* !GSD_ACTION_DELAY */

static const gchar introspection_xml[] =
    "<node>"
    "  <interface name='org.gnome.SettingsDaemon.Power'>"
    "    <property name='Icon' type='s' access='read'/>"
    "    <property name='Tooltip' type='s' access='read'/>"
    "    <property name='Percentage' type='d' access='read'/>"
    "    <method name='GetPrimaryDevice'>"
    "      <arg name='device' type='(susdut)' direction='out' />"
    "    </method>"
    "    <method name='GetDevices'>"
    "      <arg name='devices' type='a(susdut)' direction='out' />"
    "    </method>"
    "  </interface>"
    "  <interface name='org.gnome.SettingsDaemon.Power.Screen'>"
    "    <method name='StepUp'>"
    "      <arg type='u' name='new_percentage' direction='out'/>"
    "    </method>"
    "    <method name='StepDown'>"
    "      <arg type='u' name='new_percentage' direction='out'/>"
    "    </method>"
    "    <method name='GetPercentage'>"
    "      <arg type='u' name='percentage' direction='out'/>"
    "    </method>"
    "    <method name='SetPercentage'>"
    "      <arg type='u' name='percentage' direction='in'/>"
    "      <arg type='u' name='new_percentage' direction='out'/>"
    "    </method>"
    "    <signal name='Changed'/>"
    "  </interface>"
    "  <interface name='org.gnome.SettingsDaemon.Power.Keyboard'>"
    "    <method name='StepUp'>"
    "      <arg type='u' name='new_percentage' direction='out'/>"
    "    </method>"
    "    <method name='StepDown'>"
    "      <arg type='u' name='new_percentage' direction='out'/>"
    "    </method>"
    "    <method name='Toggle'>"
    "      <arg type='u' name='new_percentage' direction='out'/>"
    "    </method>"
    "  </interface>"
    "</node>";

#define GSD_POWER_MANAGER_GET_PRIVATE(o) (G_TYPE_INSTANCE_GET_PRIVATE ((o), GSD_TYPE_POWER_MANAGER, GsdPowerManagerPrivate))

typedef enum
{
    GSD_POWER_IDLE_MODE_NORMAL,
    GSD_POWER_IDLE_MODE_DIM,
    GSD_POWER_IDLE_MODE_BLANK,
    GSD_POWER_IDLE_MODE_SLEEP
} GsdPowerIdleMode;

struct GsdPowerManagerPrivate
{
    /* D-Bus */
    GDBusProxy              *session;
    guint                    name_id;
    GDBusNodeInfo           *introspection_data;
    GDBusConnection         *connection;
    GCancellable            *bus_cancellable;
    GDBusProxy              *session_presence_proxy;

    /* Settings */
    GSettings               *settings_profile;
    GSettings               *settings;
    GSettings               *settings_session;
    GSettings               *settings_screensaver;
    GSettings               *settings_xrandr;

    gboolean                 use_time_primary;
    guint                    action_percentage;
    guint                    action_time;
    guint                    critical_percentage;
    guint                    critical_time;
    guint                    low_percentage;
    guint                    low_time;

    /* Screensaver */
    guint                    screensaver_watch_id;
    GCancellable            *screensaver_cancellable;
    GDBusProxy              *screensaver_proxy;
    gboolean                 screensaver_active;

    /* State */
    gchar                   *settings_path;
    gboolean                 lid_is_closed;
    UpClient                *up_client;
    gchar                   *previous_summary;
    GIcon                   *previous_icon;
    GPtrArray               *devices_array;
    UpDevice                *device_composite;
    GnomeRRScreen           *rr_screen;
    NotifyNotification      *notification_ups_discharging;
    NotifyNotification      *notification_low;
    NotifyNotification      *notification_sleep_warning;
    NotifyNotification      *notification_logout_warning;
    GsdPowerActionType       sleep_action_type;
    gboolean                 battery_is_low; /* laptop battery low, or UPS discharging */

    /* Brightness */
    gboolean                 backlight_available;
    gint                     pre_dim_brightness; /* level, not percentage */

    /* Keyboard */
    GDBusProxy              *upower_kdb_proxy;
    gint                     kbd_brightness_now;
    gint                     kbd_brightness_max;
    gint                     kbd_brightness_old;
    gint                     kbd_brightness_pre_dim;

    /* Sound */
    guint32                  critical_alert_timeout_id;

    /* systemd stuff */
    GDBusProxy              *logind_proxy;
    gint                     inhibit_lid_switch_fd;
    gboolean                 inhibit_lid_switch_taken;
    gint                     inhibit_suspend_fd;
    gboolean                 inhibit_suspend_taken;
    guint                    inhibit_lid_switch_timer_id;
    gboolean                 is_virtual_machine;

    /* Idles */
    GnomeIdleMonitor        *idle_monitor;
    guint                    idle_dim_id;
    guint                    idle_blank_id;
    guint                    idle_sleep_warning_id;
    guint                    idle_sleep_id;
    GsdPowerIdleMode         current_idle_mode;

    guint                    temporary_unidle_on_ac_id;
    GsdPowerIdleMode         previous_idle_mode;

    guint                    xscreensaver_watchdog_timer_id;
};

enum
{
    PROP_0,
};

static void     gsd_power_manager_class_init  (GsdPowerManagerClass *klass);
static void     gsd_power_manager_init        (GsdPowerManager      *power_manager);

static UpDevice *engine_get_composite_device (GsdPowerManager *manager, UpDevice *original_device);
static UpDevice *engine_update_composite_device (GsdPowerManager *manager, UpDevice *original_device);
static GIcon    *engine_get_icon (GsdPowerManager *manager);
static gchar    *engine_get_summary (GsdPowerManager *manager);
static gdouble   engine_get_percentage (GsdPowerManager *manager);
static void      do_power_action_type (GsdPowerManager *manager, GsdPowerActionType action_type);
static void      do_lid_closed_action (GsdPowerManager *manager);
static void      uninhibit_lid_switch (GsdPowerManager *manager);
static void      main_battery_or_ups_low_changed (GsdPowerManager *manager, gboolean is_low);
static gboolean  idle_is_session_inhibited (GsdPowerManager *manager, guint mask, gboolean *is_inhibited);
static void      idle_set_mode (GsdPowerManager *manager, GsdPowerIdleMode mode);
static void      idle_triggered_idle_cb (GnomeIdleMonitor *monitor, guint watch_id, gpointer user_data);
static void      idle_became_active_cb (GnomeIdleMonitor *monitor, guint watch_id, gpointer user_data);

static void
engine_charge_low (GsdPowerManager *manager, UpDevice *device);
static void
engine_charge_critical (GsdPowerManager *manager, UpDevice *device);
static void
engine_charge_action (GsdPowerManager *manager, UpDevice *device);

G_DEFINE_TYPE (GsdPowerManager, gsd_power_manager, G_TYPE_OBJECT)

static gpointer manager_object = NULL;

GQuark
gsd_power_manager_error_quark (void)
{
    static GQuark quark = 0;
    if (!quark)
        quark = g_quark_from_static_string ("gsd_power_manager_error");
    return quark;
}

static void
notify_close_if_showing (NotifyNotification **notification)
{
    if (*notification == NULL)
        return;
    notify_notification_close (*notification, NULL);
    g_clear_object (notification);
}

typedef enum
{
    WARNING_NONE            = 0,
    WARNING_DISCHARGING     = 1,
    WARNING_LOW             = 2,
    WARNING_CRITICAL        = 3,
    WARNING_ACTION          = 4
} GsdPowerManagerWarning;

static GVariant *
engine_get_icon_property_variant (GsdPowerManager  *manager)
{
    GIcon *icon;
    GVariant *retval;

    icon = engine_get_icon (manager);
    if (icon != NULL)
    {
        char *str;
        str = g_icon_to_string (icon);
        g_object_unref (icon);
        retval = g_variant_new_string (str);
        g_free (str);
    }
    else
    {
        retval = g_variant_new_string ("");
    }
    return retval;
}

static GVariant *
engine_get_tooltip_property_variant (GsdPowerManager  *manager)
{
    char *tooltip;
    GVariant *retval;

    tooltip = engine_get_summary (manager);
    retval = g_variant_new_string (tooltip != NULL ? tooltip : "");
    g_free (tooltip);

    return retval;
}

static void
engine_emit_changed (GsdPowerManager *manager,
                     gboolean         icon_changed,
                     gboolean         state_changed)
{
    GVariantBuilder props_builder;
    GVariant *props_changed = NULL;
    GError *error = NULL;

    /* not yet connected to the bus */
    if (manager->priv->connection == NULL)
        return;

    g_variant_builder_init (&props_builder, G_VARIANT_TYPE ("a{sv}"));

    if (icon_changed)
        g_variant_builder_add (&props_builder, "{sv}", "Icon",
                               engine_get_icon_property_variant (manager));
    if (state_changed)
        g_variant_builder_add (&props_builder, "{sv}", "Tooltip",
                               engine_get_tooltip_property_variant (manager));
    g_variant_builder_add (&props_builder, "{sv}", "Percentage",
                           g_variant_new_double (engine_get_percentage (manager)));

    props_changed = g_variant_new ("(s@a{sv}@as)", GSD_POWER_DBUS_INTERFACE,
                                   g_variant_builder_end (&props_builder),
                                   g_variant_new_strv (NULL, 0));
    g_variant_ref_sink (props_changed);

    if (!g_dbus_connection_emit_signal (manager->priv->connection,
                                        NULL,
                                        GSD_POWER_DBUS_PATH,
                                        "org.freedesktop.DBus.Properties",
                                        "PropertiesChanged",
                                        props_changed,
                                        &error))
        goto out;

out:
    if (error)
    {
        g_warning ("%s", error->message);
        g_clear_error (&error);
    }
    if (props_changed)
        g_variant_unref (props_changed);
}

static GsdPowerManagerWarning
engine_get_warning_csr (GsdPowerManager *manager, UpDevice *device)
{
    gdouble percentage;

    /* get device properties */
    g_object_get (device, "percentage", &percentage, NULL);

    if (percentage < 26.0f)
        return WARNING_LOW;
    else if (percentage < 13.0f)
        return WARNING_CRITICAL;
    return WARNING_NONE;
}

/**
 * get the warning based on the percentage
 */
static GsdPowerManagerWarning
engine_get_warning_percentage (GsdPowerManager *manager, UpDevice *device)
{
    gdouble percentage;

    /* get device properties */
    g_object_get (device, "percentage", &percentage, NULL);

    if (percentage <= manager->priv->action_percentage)
        return WARNING_ACTION;
    if (percentage <= manager->priv->critical_percentage)
        return WARNING_CRITICAL;
    if (percentage <= manager->priv->low_percentage)
        return WARNING_LOW;
    return WARNING_NONE;
}

static GsdPowerManagerWarning
engine_get_warning_time (GsdPowerManager *manager, UpDevice *device)
{
    UpDeviceKind kind;
    gint64 time_to_empty;

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  "time-to-empty", &time_to_empty,
                  NULL);

    /* this is probably an error condition */
    if (time_to_empty == 0)
    {
        g_debug ("time zero, falling back to percentage for %s",
                 up_device_kind_to_string (kind));
        return engine_get_warning_percentage (manager, device);
    }

    if (time_to_empty <= manager->priv->action_time)
        return WARNING_ACTION;
    if (time_to_empty <= manager->priv->critical_time)
        return WARNING_CRITICAL;
    if (time_to_empty <= manager->priv->low_time)
        return WARNING_LOW;
    return WARNING_NONE;
}

/**
 * This gets the possible engine state for the device according to the
 * policy, which could be per-percent, or per-time.
 **/
static GsdPowerManagerWarning
engine_get_warning (GsdPowerManager *manager, UpDevice *device)
{
    UpDeviceKind kind;
    UpDeviceState state;
    GsdPowerManagerWarning warning_type;

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  "state", &state,
                  NULL);

    /* default to no engine */
    warning_type = WARNING_NONE;

    /* if the device in question is on ac, don't give a warning */
    if (state == UP_DEVICE_STATE_CHARGING)
        goto out;

    if (kind == UP_DEVICE_KIND_MOUSE ||
            kind == UP_DEVICE_KIND_KEYBOARD)
    {

        warning_type = engine_get_warning_csr (manager, device);

    }
    else if (kind == UP_DEVICE_KIND_UPS ||
             kind == UP_DEVICE_KIND_MEDIA_PLAYER ||
             kind == UP_DEVICE_KIND_TABLET ||
             kind == UP_DEVICE_KIND_COMPUTER ||
             kind == UP_DEVICE_KIND_PDA)
    {

        warning_type = engine_get_warning_percentage (manager, device);

    }
    else if (kind == UP_DEVICE_KIND_PHONE)
    {

        warning_type = engine_get_warning_percentage (manager, device);

    }
    else if (kind == UP_DEVICE_KIND_BATTERY)
    {
        /* only use the time when it is accurate, and settings is not disabled */
        if (manager->priv->use_time_primary)
            warning_type = engine_get_warning_time (manager, device);
        else
            warning_type = engine_get_warning_percentage (manager, device);
    }

    /* If we have no important engines, we should test for discharging */
    if (warning_type == WARNING_NONE)
    {
        if (state == UP_DEVICE_STATE_DISCHARGING)
            warning_type = WARNING_DISCHARGING;
    }

out:
    return warning_type;
}

static gchar *
engine_get_summary (GsdPowerManager *manager)
{
    guint i;
    GPtrArray *array;
    UpDevice *device;
    UpDeviceState state;
    GString *tooltip = NULL;
    gchar *part;
    gboolean is_present;


    /* need to get AC state */
    tooltip = g_string_new ("");

    /* do we have specific device types? */
    array = manager->priv->devices_array;
    for (i = 0; i < array->len; i++)
    {
        device = g_ptr_array_index (array, i);
        g_object_get (device,
                      "is-present", &is_present,
                      "state", &state,
                      NULL);
        if (!is_present)
            continue;
        if (state == UP_DEVICE_STATE_EMPTY)
            continue;
        part = gpm_upower_get_device_summary (device);
        if (part != NULL)
            g_string_append_printf (tooltip, "%s\n", part);
        g_free (part);
    }

    /* remove the last \n */
    g_string_truncate (tooltip, tooltip->len - 1);

    g_debug ("tooltip: %s", tooltip->str);

    return g_string_free (tooltip, FALSE);
}

static gdouble
engine_get_percentage (GsdPowerManager *manager)
{
    guint i;
    GPtrArray *array;
    UpDevice *device;
    UpDeviceKind kind;
    gboolean is_present;
    gdouble percentage;

    array = manager->priv->devices_array;
    for (i = 0; i < array->len ; i++)
    {
        device = g_ptr_array_index (array, i);

        /* get device properties */
        g_object_get (device,
                      "kind", &kind,
                      "is-present", &is_present,
                      NULL);

        /* if battery then use composite device to cope with multiple batteries */
        if (kind == UP_DEVICE_KIND_BATTERY)
            device = engine_get_composite_device (manager, device);

        if (is_present)
        {
            /* Doing it here as it could be a composite device */
            g_object_get (device, "percentage", &percentage, NULL);
            return percentage;
        }
    }
    return -1;

}

static GIcon *
engine_get_icon_priv (GsdPowerManager *manager,
                      UpDeviceKind device_kind,
                      GsdPowerManagerWarning warning,
                      gboolean use_state)
{
    guint i;
    GPtrArray *array;
    UpDevice *device;
    GsdPowerManagerWarning warning_temp;
    UpDeviceKind kind;
    UpDeviceState state;
    gboolean is_present;

    /* do we have specific device types? */
    array = manager->priv->devices_array;
    for (i = 0; i < array->len; i++)
    {
        device = g_ptr_array_index (array, i);

        /* get device properties */
        g_object_get (device,
                      "kind", &kind,
                      "state", &state,
                      "is-present", &is_present,
                      NULL);

        /* if battery then use composite device to cope with multiple batteries */
        if (kind == UP_DEVICE_KIND_BATTERY)
            device = engine_get_composite_device (manager, device);

        warning_temp = GPOINTER_TO_INT(g_object_get_data (G_OBJECT(device),
                                       "engine-warning-old"));
        if (kind == device_kind && is_present)
        {
            if (warning != WARNING_NONE)
            {
                if (warning_temp == warning)
                    return gpm_upower_get_device_icon (device, TRUE);
                continue;
            }
            if (use_state)
            {
                if (state == UP_DEVICE_STATE_CHARGING ||
                        state == UP_DEVICE_STATE_DISCHARGING)
                    return gpm_upower_get_device_icon (device, TRUE);
                continue;
            }
            return gpm_upower_get_device_icon (device, TRUE);
        }
    }
    return NULL;
}

static GIcon *
engine_get_icon (GsdPowerManager *manager)
{
    GIcon *icon = NULL;


    /* we try CRITICAL: BATTERY, UPS, MOUSE, KEYBOARD */
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_BATTERY, WARNING_CRITICAL, FALSE);
    if (icon != NULL)
        return icon;
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_UPS, WARNING_CRITICAL, FALSE);
    if (icon != NULL)
        return icon;
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_MOUSE, WARNING_CRITICAL, FALSE);
    if (icon != NULL)
        return icon;
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_KEYBOARD, WARNING_CRITICAL, FALSE);
    if (icon != NULL)
        return icon;

    /* we try CRITICAL: BATTERY, UPS, MOUSE, KEYBOARD */
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_BATTERY, WARNING_LOW, FALSE);
    if (icon != NULL)
        return icon;
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_UPS, WARNING_LOW, FALSE);
    if (icon != NULL)
        return icon;
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_MOUSE, WARNING_LOW, FALSE);
    if (icon != NULL)
        return icon;
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_KEYBOARD, WARNING_LOW, FALSE);
    if (icon != NULL)
        return icon;

    /* we try (DIS)CHARGING: BATTERY, UPS */
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_BATTERY, WARNING_NONE, TRUE);
    if (icon != NULL)
        return icon;
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_UPS, WARNING_NONE, TRUE);
    if (icon != NULL)
        return icon;

    /* we try PRESENT: BATTERY, UPS */
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_BATTERY, WARNING_NONE, FALSE);
    if (icon != NULL)
        return icon;
    icon = engine_get_icon_priv (manager, UP_DEVICE_KIND_UPS, WARNING_NONE, FALSE);
    if (icon != NULL)
        return icon;

    /* do not show an icon */
    return NULL;
}

static gboolean
engine_recalculate_state_icon (GsdPowerManager *manager)
{
    GIcon *icon;

    /* show a different icon if we are disconnected */
    icon = engine_get_icon (manager);

    if (g_icon_equal (icon, manager->priv->previous_icon))
    {
        g_clear_object (&icon);
        return FALSE;
    }

    g_clear_object (&manager->priv->previous_icon);
    manager->priv->previous_icon = icon;

    g_debug ("Icon changed");

    return TRUE;
}

static gboolean
engine_recalculate_state_summary (GsdPowerManager *manager)
{
    char *summary;

    summary = engine_get_summary (manager);

    if (g_strcmp0 (manager->priv->previous_summary, summary) == 0)
    {
        g_free (summary);
        return FALSE;
    }

    g_free (manager->priv->previous_summary);
    manager->priv->previous_summary = summary;

    g_debug ("Summary changed");

    return TRUE;
}

static void
engine_recalculate_state (GsdPowerManager *manager)
{
    gboolean icon_changed = FALSE;
    gboolean state_changed = FALSE;

    icon_changed = engine_recalculate_state_icon (manager);
    state_changed = engine_recalculate_state_summary (manager);

    /* only emit if the icon or summary has changed */
    if (icon_changed || state_changed)
        engine_emit_changed (manager, icon_changed, state_changed);
}

static UpDevice *
engine_get_composite_device (GsdPowerManager *manager,
                             UpDevice *original_device)
{
    guint battery_devices = 0;
    GPtrArray *array;
    UpDevice *device;
    UpDeviceKind kind;
    UpDeviceKind original_kind;
    guint i;

    /* get the type of the original device */
    g_object_get (original_device,
                  "kind", &original_kind,
                  NULL);

    /* find out how many batteries in the system */
    array = manager->priv->devices_array;
    for (i = 0; i < array->len; i++)
    {
        device = g_ptr_array_index (array, i);
        g_object_get (device,
                      "kind", &kind,
                      NULL);
        if (kind == original_kind)
            battery_devices++;
    }

    /* just use the original device if only one primary battery */
    if (battery_devices <= 1)
        return original_device;

    /* use the composite device */
    device = manager->priv->device_composite;

    /* return composite device or original device */
    return device;
}

static UpDevice *
engine_update_composite_device (GsdPowerManager *manager,
                                UpDevice *original_device)
{
    guint i;
    gdouble percentage = 0.0;
    gdouble energy = 0.0;
    gdouble energy_full = 0.0;
    gdouble energy_rate = 0.0;
    gdouble energy_total = 0.0;
    gdouble energy_full_total = 0.0;
    gdouble energy_rate_total = 0.0;
    gint64 time_to_empty = 0;
    gint64 time_to_full = 0;
    guint battery_devices = 0;
    gboolean is_charging = FALSE;
    gboolean is_discharging = FALSE;
    gboolean is_fully_charged = TRUE;
    GPtrArray *array;
    UpDevice *device;
    UpDeviceState state;
    UpDeviceKind kind;
    UpDeviceKind original_kind;

    /* get the type of the original device */
    g_object_get (original_device,
                  "kind", &original_kind,
                  NULL);

    /* update the composite device */
    array = manager->priv->devices_array;
    for (i = 0; i < array->len; i++)
    {
        device = g_ptr_array_index (array, i);
        g_object_get (device,
                      "kind", &kind,
                      "state", &state,
                      "energy", &energy,
                      "energy-full", &energy_full,
                      "energy-rate", &energy_rate,
                      NULL);
        if (kind != original_kind)
            continue;

        /* one of these will be charging or discharging */
        if (state == UP_DEVICE_STATE_CHARGING)
            is_charging = TRUE;
        if (state == UP_DEVICE_STATE_DISCHARGING)
            is_discharging = TRUE;
        if (state != UP_DEVICE_STATE_FULLY_CHARGED)
            is_fully_charged = FALSE;

        /* sum up composite */
        energy_total += energy;
        energy_full_total += energy_full;
        energy_rate_total += energy_rate;
        battery_devices++;
    }

    /* just use the original device if only one primary battery */
    if (battery_devices == 1)
    {
        g_debug ("using original device as only one primary battery");
        device = original_device;
        goto out;
    }

    /* use percentage weighted for each battery capacity */
    if (energy_full_total > 0.0)
        percentage = 100.0 * energy_total / energy_full_total;

    /* set composite state */
    if (is_charging)
        state = UP_DEVICE_STATE_CHARGING;
    else if (is_discharging)
        state = UP_DEVICE_STATE_DISCHARGING;
    else if (is_fully_charged)
        state = UP_DEVICE_STATE_FULLY_CHARGED;
    else
        state = UP_DEVICE_STATE_UNKNOWN;

    /* calculate a quick and dirty time remaining value */
    if (energy_rate_total > 0)
    {
        if (state == UP_DEVICE_STATE_DISCHARGING)
            time_to_empty = 3600 * (energy_total / energy_rate_total);
        else if (state == UP_DEVICE_STATE_CHARGING)
            time_to_full = 3600 * ((energy_full_total - energy_total) / energy_rate_total);
    }

    /* okay, we can use the composite device */
    device = manager->priv->device_composite;

    g_debug ("printing composite device");
    g_object_set (device,
                  "energy", energy,
                  "energy-full", energy_full,
                  "energy-rate", energy_rate,
                  "time-to-empty", time_to_empty,
                  "time-to-full", time_to_full,
                  "percentage", percentage,
                  "state", state,
                  NULL);

    /* force update of icon */
    if (engine_recalculate_state_icon (manager))
        engine_emit_changed (manager, TRUE, FALSE);
out:
    /* return composite device or original device */
    return device;
}

typedef struct
{
    GsdPowerManager *manager;
    UpDevice        *device;
} GsdPowerManagerRecallData;

static void
device_perhaps_recall_response_cb (GtkDialog *dialog,
                                   gint response_id,
                                   GsdPowerManagerRecallData *recall_data)
{
    GdkScreen *screen;
    GtkWidget *dialog_error;
    GError *error = NULL;
    gboolean ret;
    gchar *website = NULL;

    /* don't show this again */
    if (response_id == GTK_RESPONSE_CANCEL)
    {
        g_settings_set_boolean (recall_data->manager->priv->settings,
                                "notify-perhaps-recall",
                                FALSE);
        goto out;
    }

    /* visit recall website */
    if (response_id == GTK_RESPONSE_OK)
    {

        g_object_get (recall_data->device,
                      "recall-url", &website,
                      NULL);

        screen = gdk_screen_get_default();
        ret = gtk_show_uri (screen,
                            website,
                            gtk_get_current_event_time (),
                            &error);
        if (!ret)
        {
            dialog_error = gtk_message_dialog_new (NULL,
                                                   GTK_DIALOG_MODAL,
                                                   GTK_MESSAGE_INFO,
                                                   GTK_BUTTONS_OK,
                                                   "Failed to show url %s",
                                                   error->message);
            gtk_dialog_run (GTK_DIALOG (dialog_error));
            g_error_free (error);
        }
    }
out:
    gtk_widget_destroy (GTK_WIDGET (dialog));
    g_object_unref (recall_data->device);
    g_object_unref (recall_data->manager);
    g_free (recall_data);
    g_free (website);
    return;
}

static gboolean
device_perhaps_recall_delay_cb (gpointer user_data)
{
    gchar *vendor;
    const gchar *title = NULL;
    GString *message = NULL;
    GtkWidget *dialog;
    GsdPowerManagerRecallData *recall_data = (GsdPowerManagerRecallData *) user_data;

    g_object_get (recall_data->device,
                  "recall-vendor", &vendor,
                  NULL);

    /* TRANSLATORS: the battery may be recalled by its vendor */
    title = _("Battery may be recalled");
    message = g_string_new ("");
    g_string_append_printf (message,
                            _("A battery in your computer may have been "
                              "recalled by %s and you may be at risk."), vendor);
    g_string_append (message, "\n\n");
    g_string_append (message, _("For more information visit the battery recall website."));
    dialog = gtk_message_dialog_new_with_markup (NULL,
             GTK_DIALOG_DESTROY_WITH_PARENT,
             GTK_MESSAGE_INFO,
             GTK_BUTTONS_CLOSE,
             "<span size='larger'><b>%s</b></span>",
             title);
    gtk_message_dialog_format_secondary_markup (GTK_MESSAGE_DIALOG (dialog),
            "%s", message->str);

    /* TRANSLATORS: button text, visit the manufacturers recall website */
    gtk_dialog_add_button (GTK_DIALOG (dialog), _("Visit recall website"),
                           GTK_RESPONSE_OK);

    /* TRANSLATORS: button text, do not show this bubble again */
    gtk_dialog_add_button (GTK_DIALOG (dialog), _("Do not show me this again"),
                           GTK_RESPONSE_CANCEL);

    gtk_widget_show (dialog);
    g_signal_connect (dialog, "response",
                      G_CALLBACK (device_perhaps_recall_response_cb),
                      recall_data);

    g_string_free (message, TRUE);
    g_free (vendor);
    return FALSE;
}

static void
device_perhaps_recall (GsdPowerManager *manager, UpDevice *device)
{
    gboolean ret;
    guint timer_id;
    GsdPowerManagerRecallData *recall_data;

    /* don't show when running under GDM */
    if (g_getenv ("RUNNING_UNDER_GDM") != NULL)
    {
        g_debug ("running under gdm, so no notification");
        return;
    }

    /* already shown, and dismissed */
    ret = g_settings_get_boolean (manager->priv->settings,
                                  "notify-perhaps-recall");
    if (!ret)
    {
        g_debug ("settings prevents recall notification");
        return;
    }

    recall_data = g_new0 (GsdPowerManagerRecallData, 1);
    recall_data->manager = g_object_ref (manager);
    recall_data->device = g_object_ref (device);

    /* delay by a few seconds so the session can load */
    timer_id = g_timeout_add_seconds (GSD_POWER_MANAGER_RECALL_DELAY,
                                      device_perhaps_recall_delay_cb,
                                      recall_data);
    g_source_set_name_by_id (timer_id, "[GsdPowerManager] perhaps-recall");
}

static void
engine_device_add (GsdPowerManager *manager, UpDevice *device)
{
    gboolean recall_notice;
    GsdPowerManagerWarning warning;
    UpDeviceState state;
    UpDeviceKind kind;
    UpDevice *composite;

    /* assign warning */
    warning = engine_get_warning (manager, device);

    //in case the battery percentage is low/critical at startup time
    if (warning == WARNING_LOW)
    {
        g_debug ("** EMIT: charge-low");
        engine_charge_low (manager, device);
    }
    else if (warning == WARNING_CRITICAL)
    {
        g_debug ("** EMIT: charge-critical");
        engine_charge_critical (manager, device);
    }
    else if (warning == WARNING_ACTION)
    {
        g_debug ("charge-action");
        engine_charge_action (manager, device);
    }
    /* save new state */
    g_object_set_data (G_OBJECT(device),
                       "engine-warning-old",
                       GUINT_TO_POINTER(warning));

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  "state", &state,
                  "recall-notice", &recall_notice,
                  NULL);

    /* add old state for transitions */
    g_debug ("adding %s with state %s",
             up_device_get_object_path (device), up_device_state_to_string (state));
    g_object_set_data (G_OBJECT(device),
                       "engine-state-old",
                       GUINT_TO_POINTER(state));

    if (kind == UP_DEVICE_KIND_BATTERY)
    {
        g_debug ("updating because we added a device");
        composite = engine_update_composite_device (manager, device);

        /* get the same values for the composite device */
        warning = engine_get_warning (manager, composite);
        g_object_set_data (G_OBJECT(composite),
                           "engine-warning-old",
                           GUINT_TO_POINTER(warning));
        g_object_get (composite, "state", &state, NULL);
        g_object_set_data (G_OBJECT(composite),
                           "engine-state-old",
                           GUINT_TO_POINTER(state));
    }

    /* the device is recalled */
    if (recall_notice)
        device_perhaps_recall (manager, device);
}

static gboolean
engine_check_recall (GsdPowerManager *manager, UpDevice *device)
{
    UpDeviceKind kind;
    gboolean recall_notice = FALSE;
    gchar *recall_vendor = NULL;
    gchar *recall_url = NULL;

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  "recall-notice", &recall_notice,
                  "recall-vendor", &recall_vendor,
                  "recall-url", &recall_url,
                  NULL);

    /* not battery */
    if (kind != UP_DEVICE_KIND_BATTERY)
        goto out;

    /* no recall data */
    if (!recall_notice)
        goto out;

    /* emit signal for manager */
    g_debug ("** EMIT: perhaps-recall");
    g_debug ("%s-%s", recall_vendor, recall_url);
out:
    g_free (recall_vendor);
    g_free (recall_url);
    return recall_notice;
}

static gboolean
engine_coldplug (GsdPowerManager *manager)
{
    guint i;
    GPtrArray *array = NULL;
    UpDevice *device;
    gboolean ret;
    GError *error = NULL;

    /* get devices from UPower */
    ret = up_client_enumerate_devices_sync (manager->priv->up_client, NULL, &error);
    if (!ret)
    {
        g_warning ("failed to get device list: %s", error->message);
        g_error_free (error);
        goto out;
    }

    engine_recalculate_state (manager);

    /* add to database */
    array = up_client_get_devices (manager->priv->up_client);
    if (array == NULL)
        goto out;

    for (i = 0; i < array->len; i++)
    {
        device = g_ptr_array_index (array, i);
        engine_device_add (manager, device);
        engine_check_recall (manager, device);
    }
out:
    if (array != NULL)
        g_ptr_array_unref (array);
    /* never repeat */
    return FALSE;
}

static void
engine_device_added_cb (UpClient *client, UpDevice *device, GsdPowerManager *manager)
{
    /* add to list */
    g_ptr_array_add (manager->priv->devices_array, g_object_ref (device));
    engine_check_recall (manager, device);

    engine_recalculate_state (manager);
}

static void
engine_device_removed_cb (UpClient *client, UpDevice *device, GsdPowerManager *manager)
{
    gboolean ret;
    ret = g_ptr_array_remove (manager->priv->devices_array, device);
    if (!ret)
        return;
    engine_recalculate_state (manager);
}

static void
engine_profile_changed_cb (GSettings *settings,
                           const gchar *key,
                           GsdPowerManager *manager)
{
    /*if (g_strcmp0 (key, "use-time-for-policy") == 0)*/
    /*{*/
    /*manager->priv->use_time_primary = g_settings_get_boolean (settings, key);*/
    /*return;*/
    /*}*/
    /*if (g_str_has_prefix (key, "sleep-inactive") ||*/
    /*g_str_equal (key, "idle-delay") ||*/
    /*g_str_equal (key, "idle-dim"))*/
    /*{*/
    /*g_debug("engine_settings_key_changed_cb()\n");*/
    /*idle_configure (manager);*/
    /*return;*/
    /*}*/
    GError *error = NULL;
    gsd_power_manager_stop(manager);
    gsd_power_manager_start(manager, &error);

    g_debug("Exit\n");
    gtk_main_quit();

    return;
}

static void
on_notification_closed (NotifyNotification *notification, gpointer data)
{
    g_object_unref (notification);
}

static const gchar *
get_first_themed_icon_name (GIcon *icon)
{
    const gchar* const *icon_names;
    const gchar *icon_name = NULL;

    /* no icon */
    if (icon == NULL)
        goto out;

    /* just use the first icon */
    icon_names = g_themed_icon_get_names (G_THEMED_ICON (icon));
    if (icon_names != NULL)
        icon_name = icon_names[0];
out:
    return icon_name;
}

static void
create_notification (const char *summary,
                     const char *body,
                     GIcon      *icon,
                     NotifyNotification **weak_pointer_location)
{
    NotifyNotification *notification;

    notification = notify_notification_new (summary,
                                            body,
                                            icon ? get_first_themed_icon_name (icon) : NULL);
    *weak_pointer_location = notification;
    g_object_add_weak_pointer (G_OBJECT (notification),
                               (gpointer *) weak_pointer_location);
    g_signal_connect (notification, "closed",
                      G_CALLBACK (on_notification_closed), NULL);
}

static void
engine_ups_discharging (GsdPowerManager *manager, UpDevice *device)
{
    const gchar *title;
    gchar *remaining_text = NULL;
    gdouble percentage;
    GIcon *icon = NULL;
    gint64 time_to_empty;
    GString *message;
    UpDeviceKind kind;

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  "percentage", &percentage,
                  "time-to-empty", &time_to_empty,
                  NULL);

    if (kind != UP_DEVICE_KIND_UPS)
        return;

    main_battery_or_ups_low_changed (manager, TRUE);

    /* only show text if there is a valid time */
    if (time_to_empty > 0)
        remaining_text = gpm_get_timestring (time_to_empty);

    /* TRANSLATORS: UPS is now discharging */
    title = _("UPS Discharging");

    message = g_string_new ("");
    if (remaining_text != NULL)
    {
        /* TRANSLATORS: tell the user how much time they have got */
        g_string_append_printf (message, _("%s of UPS backup power remaining"),
                                remaining_text);
    }
    else
    {
        g_string_append (message, gpm_device_to_localised_string (device));
    }
    g_string_append_printf (message, " (%.0f%%)", percentage);

    icon = gpm_upower_get_device_icon (device, TRUE);

    /* close any existing notification of this class */
    notify_close_if_showing (&manager->priv->notification_ups_discharging);

    /* create a new notification */
    create_notification (title, message->str,
                         icon,
                         &manager->priv->notification_ups_discharging);
    notify_notification_set_timeout (manager->priv->notification_ups_discharging,
                                     GSD_POWER_MANAGER_NOTIFY_TIMEOUT_LONG);
    notify_notification_set_urgency (manager->priv->notification_ups_discharging,
                                     NOTIFY_URGENCY_NORMAL);
    /* TRANSLATORS: this is the notification application name */
    notify_notification_set_app_name (manager->priv->notification_ups_discharging, _("Power"));
    notify_notification_set_hint (manager->priv->notification_ups_discharging,
                                  "transient", g_variant_new_boolean (TRUE));

    notify_notification_show (manager->priv->notification_ups_discharging, NULL);

    g_string_free (message, TRUE);
    if (icon != NULL)
        g_object_unref (icon);
    g_free (remaining_text);
}

static GsdPowerActionType
manager_critical_action_get (GsdPowerManager *manager,
                             gboolean         is_ups)
{
    GsdPowerActionType policy;

    policy = g_settings_get_enum (manager->priv->settings, "critical-battery-action");
    if (policy == GSD_POWER_ACTION_SUSPEND)
    {
        if (is_ups == FALSE &&
                up_client_get_can_suspend (manager->priv->up_client))
            return policy;
        return GSD_POWER_ACTION_SHUTDOWN;
    }
    else if (policy == GSD_POWER_ACTION_HIBERNATE)
    {
        if (up_client_get_can_hibernate (manager->priv->up_client))
            return policy;
        return GSD_POWER_ACTION_SHUTDOWN;
    }

    return policy;
}

static gboolean
manager_critical_action_do (GsdPowerManager *manager,
                            gboolean         is_ups)
{
    GsdPowerActionType action_type;

    /* stop playing the alert as it's too late to do anything now */
    play_loop_stop (&manager->priv->critical_alert_timeout_id);

    action_type = manager_critical_action_get (manager, is_ups);
    do_power_action_type (manager, action_type);

    return FALSE;
}

static gboolean
manager_critical_action_do_cb (GsdPowerManager *manager)
{
    manager_critical_action_do (manager, FALSE);
    return FALSE;
}

static gboolean
manager_critical_ups_action_do_cb (GsdPowerManager *manager)
{
    manager_critical_action_do (manager, TRUE);
    return FALSE;
}

static gboolean
engine_just_laptop_battery (GsdPowerManager *manager)
{
    UpDevice *device;
    UpDeviceKind kind;
    GPtrArray *array;
    gboolean ret = TRUE;
    guint i;

    /* find if there are any other device types that mean we have to
     * be more specific in our wording */
    array = manager->priv->devices_array;
    for (i = 0; i < array->len; i++)
    {
        device = g_ptr_array_index (array, i);
        g_object_get (device, "kind", &kind, NULL);
        if (kind != UP_DEVICE_KIND_BATTERY)
        {
            ret = FALSE;
            break;
        }
    }
    return ret;
}

static void
engine_charge_low (GsdPowerManager *manager, UpDevice *device)
{
    const gchar *title = NULL;
    gboolean ret;
    gchar *message = NULL;
    gchar *tmp;
    gchar *remaining_text;
    gdouble percentage;
    GIcon *icon = NULL;
    gint64 time_to_empty;
    UpDeviceKind kind;

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  "percentage", &percentage,
                  "time-to-empty", &time_to_empty,
                  NULL);

    /* check to see if the batteries have not noticed we are on AC */
    if (kind == UP_DEVICE_KIND_BATTERY)
    {
        if (!up_client_get_on_battery (manager->priv->up_client))
        {
            g_warning ("ignoring low message as we are not on battery power");
            goto out;
        }
    }

    if (kind == UP_DEVICE_KIND_BATTERY)
    {

        /* if the user has no other batteries, drop the "Laptop" wording */
        ret = engine_just_laptop_battery (manager);
        if (ret)
        {
            /* TRANSLATORS: laptop battery low, and we only have one battery */
            title = _("Battery low");
        }
        else
        {
            /* TRANSLATORS: laptop battery low, and we have more than one kind of battery */
            title = _("Laptop battery low");
        }
        tmp = gpm_get_timestring (time_to_empty);
        remaining_text = g_strconcat ("<b>", tmp, "</b>", NULL);
        g_free (tmp);

        /* TRANSLATORS: tell the user how much time they have got */
        message = g_strdup_printf (_("Approximately %s remaining (%.0f%%)"), remaining_text, percentage);
        g_free (remaining_text);

        main_battery_or_ups_low_changed (manager, TRUE);

    }
    else if (kind == UP_DEVICE_KIND_UPS)
    {
        /* TRANSLATORS: UPS is starting to get a little low */
        title = _("UPS low");
        tmp = gpm_get_timestring (time_to_empty);
        remaining_text = g_strconcat ("<b>", tmp, "</b>", NULL);
        g_free (tmp);

        /* TRANSLATORS: tell the user how much time they have got */
        message = g_strdup_printf (_("Approximately %s of remaining UPS backup power (%.0f%%)"),
                                   remaining_text, percentage);
        g_free (remaining_text);
    }
    else if (kind == UP_DEVICE_KIND_MOUSE)
    {
        /* TRANSLATORS: mouse is getting a little low */
        title = _("Mouse battery low");

        /* TRANSLATORS: tell user more details */
        message = g_strdup_printf (_("Wireless mouse is low in power (%.0f%%)"), percentage);

    }
    else if (kind == UP_DEVICE_KIND_KEYBOARD)
    {
        /* TRANSLATORS: keyboard is getting a little low */
        title = _("Keyboard battery low");

        /* TRANSLATORS: tell user more details */
        message = g_strdup_printf (_("Wireless keyboard is low in power (%.0f%%)"), percentage);

    }
    else if (kind == UP_DEVICE_KIND_PDA)
    {
        /* TRANSLATORS: PDA is getting a little low */
        title = _("PDA battery low");

        /* TRANSLATORS: tell user more details */
        message = g_strdup_printf (_("PDA is low in power (%.0f%%)"), percentage);

    }
    else if (kind == UP_DEVICE_KIND_PHONE)
    {
        /* TRANSLATORS: cell phone (mobile) is getting a little low */
        title = _("Cell phone battery low");

        /* TRANSLATORS: tell user more details */
        message = g_strdup_printf (_("Cell phone is low in power (%.0f%%)"), percentage);

    }
    else if (kind == UP_DEVICE_KIND_MEDIA_PLAYER)
    {
        /* TRANSLATORS: media player, e.g. mp3 is getting a little low */
        title = _("Media player battery low");

        /* TRANSLATORS: tell user more details */
        message = g_strdup_printf (_("Media player is low in power (%.0f%%)"), percentage);

    }
    else if (kind == UP_DEVICE_KIND_TABLET)
    {
        /* TRANSLATORS: graphics tablet, e.g. wacom is getting a little low */
        title = _("Tablet battery low");

        /* TRANSLATORS: tell user more details */
        message = g_strdup_printf (_("Tablet is low in power (%.0f%%)"), percentage);

    }
    else if (kind == UP_DEVICE_KIND_COMPUTER)
    {
        /* TRANSLATORS: computer, e.g. ipad is getting a little low */
        title = _("Attached computer battery low");

        /* TRANSLATORS: tell user more details */
        message = g_strdup_printf (_("Attached computer is low in power (%.0f%%)"), percentage);
    }

    /* get correct icon */
    icon = gpm_upower_get_device_icon (device, TRUE);

    /* close any existing notification of this class */
    notify_close_if_showing (&manager->priv->notification_low);

    /* create a new notification */
    create_notification (title, message,
                         icon,
                         &manager->priv->notification_low);
    notify_notification_set_timeout (manager->priv->notification_low,
                                     GSD_POWER_MANAGER_NOTIFY_TIMEOUT_LONG);
    notify_notification_set_urgency (manager->priv->notification_low,
                                     NOTIFY_URGENCY_NORMAL);
    notify_notification_set_app_name (manager->priv->notification_low, _("Power"));
    notify_notification_set_hint (manager->priv->notification_low,
                                  "transient", g_variant_new_boolean (TRUE));

    notify_notification_show (manager->priv->notification_low, NULL);

    /* play the sound, using sounds from the naming spec */
    ca_context_play (ca_gtk_context_get (), 0,
                     CA_PROP_EVENT_ID, "battery-low",
                     /* TRANSLATORS: this is the sound description */
                     CA_PROP_EVENT_DESCRIPTION, _("Battery is low"), NULL);

out:
    if (icon != NULL)
        g_object_unref (icon);
    g_free (message);
}

static void
engine_charge_critical (GsdPowerManager *manager, UpDevice *device)
{
    const gchar *title = NULL;
    gboolean ret;
    gchar *message = NULL;
    gdouble percentage;
    GIcon *icon = NULL;
    gint64 time_to_empty;
    GsdPowerActionType policy;
    UpDeviceKind kind;

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  "percentage", &percentage,
                  "time-to-empty", &time_to_empty,
                  NULL);

    /* check to see if the batteries have not noticed we are on AC */
    if (kind == UP_DEVICE_KIND_BATTERY)
    {
        if (!up_client_get_on_battery (manager->priv->up_client))
        {
            g_warning ("ignoring critically low message as we are not on battery power");
            goto out;
        }
    }

    if (kind == UP_DEVICE_KIND_BATTERY)
    {

        /* if the user has no other batteries, drop the "Laptop" wording */
        ret = engine_just_laptop_battery (manager);
        if (ret)
        {
            /* TRANSLATORS: laptop battery critically low, and only have one kind of battery */
            title = _("Battery critically low");
        }
        else
        {
            /* TRANSLATORS: laptop battery critically low, and we have more than one type of battery */
            title = _("Laptop battery critically low");
        }

        /* we have to do different warnings depending on the policy */
        policy = manager_critical_action_get (manager, FALSE);

        /* use different text for different actions */
        if (policy == GSD_POWER_ACTION_NOTHING)
        {
            /* TRANSLATORS: tell the use to insert the plug, as we're not going to do anything */
            message = g_strdup (_("Plug in your AC adapter to avoid losing data."));

        }
        else if (policy == GSD_POWER_ACTION_SUSPEND)
        {
            /* TRANSLATORS: give the user a ultimatum */
            message = g_strdup_printf (_("Computer will suspend very soon unless it is plugged in."));

        }
        else if (policy == GSD_POWER_ACTION_HIBERNATE)
        {
            /* TRANSLATORS: give the user a ultimatum */
            message = g_strdup_printf (_("Computer will hibernate very soon unless it is plugged in."));

        }
        else if (policy == GSD_POWER_ACTION_SHUTDOWN)
        {
            /* TRANSLATORS: give the user a ultimatum */
            message = g_strdup_printf (_("Computer will shutdown very soon unless it is plugged in."));
        }

        main_battery_or_ups_low_changed (manager, TRUE);

    }
    else if (kind == UP_DEVICE_KIND_UPS)
    {
        gchar *remaining_text;
        gchar *tmp;

        /* TRANSLATORS: the UPS is very low */
        title = _("UPS critically low");
        tmp = gpm_get_timestring (time_to_empty);
        remaining_text = g_strconcat ("<b>", tmp, "</b>", NULL);
        g_free (tmp);

        /* TRANSLATORS: give the user a ultimatum */
        message = g_strdup_printf (_("Approximately %s of remaining UPS power (%.0f%%). "
                                     "Restore AC power to your computer to avoid losing data."),
                                   remaining_text, percentage);
        g_free (remaining_text);
    }
    else if (kind == UP_DEVICE_KIND_MOUSE)
    {
        /* TRANSLATORS: the mouse battery is very low */
        title = _("Mouse battery low");

        /* TRANSLATORS: the device is just going to stop working */
        message = g_strdup_printf (_("Wireless mouse is very low in power (%.0f%%). "
                                     "This device will soon stop functioning if not charged."),
                                   percentage);
    }
    else if (kind == UP_DEVICE_KIND_KEYBOARD)
    {
        /* TRANSLATORS: the keyboard battery is very low */
        title = _("Keyboard battery low");

        /* TRANSLATORS: the device is just going to stop working */
        message = g_strdup_printf (_("Wireless keyboard is very low in power (%.0f%%). "
                                     "This device will soon stop functioning if not charged."),
                                   percentage);
    }
    else if (kind == UP_DEVICE_KIND_PDA)
    {

        /* TRANSLATORS: the PDA battery is very low */
        title = _("PDA battery low");

        /* TRANSLATORS: the device is just going to stop working */
        message = g_strdup_printf (_("PDA is very low in power (%.0f%%). "
                                     "This device will soon stop functioning if not charged."),
                                   percentage);

    }
    else if (kind == UP_DEVICE_KIND_PHONE)
    {

        /* TRANSLATORS: the cell battery is very low */
        title = _("Cell phone battery low");

        /* TRANSLATORS: the device is just going to stop working */
        message = g_strdup_printf (_("Cell phone is very low in power (%.0f%%). "
                                     "This device will soon stop functioning if not charged."),
                                   percentage);
    }
    else if (kind == UP_DEVICE_KIND_MEDIA_PLAYER)
    {

        /* TRANSLATORS: the cell battery is very low */
        title = _("Cell phone battery low");

        /* TRANSLATORS: the device is just going to stop working */
        message = g_strdup_printf (_("Media player is very low in power (%.0f%%). "
                                     "This device will soon stop functioning if not charged."),
                                   percentage);
    }
    else if (kind == UP_DEVICE_KIND_TABLET)
    {

        /* TRANSLATORS: the cell battery is very low */
        title = _("Tablet battery low");

        /* TRANSLATORS: the device is just going to stop working */
        message = g_strdup_printf (_("Tablet is very low in power (%.0f%%). "
                                     "This device will soon stop functioning if not charged."),
                                   percentage);
    }
    else if (kind == UP_DEVICE_KIND_COMPUTER)
    {

        /* TRANSLATORS: the cell battery is very low */
        title = _("Attached computer battery low");

        /* TRANSLATORS: the device is just going to stop working */
        message = g_strdup_printf (_("Attached computer is very low in power (%.0f%%). "
                                     "The device will soon shutdown if not charged."),
                                   percentage);
    }

    /* get correct icon */
    icon = gpm_upower_get_device_icon (device, TRUE);

    /* close any existing notification of this class */
    notify_close_if_showing (&manager->priv->notification_low);

    /* create a new notification */
    create_notification (title, message,
                         icon,
                         &manager->priv->notification_low);
    notify_notification_set_timeout (manager->priv->notification_low,
                                     NOTIFY_EXPIRES_NEVER);
    notify_notification_set_urgency (manager->priv->notification_low,
                                     NOTIFY_URGENCY_CRITICAL);
    notify_notification_set_app_name (manager->priv->notification_low, _("Power"));

    notify_notification_show (manager->priv->notification_low, NULL);

    switch (kind)
    {

    case UP_DEVICE_KIND_BATTERY:
    case UP_DEVICE_KIND_UPS:
        g_debug ("critical charge level reached, starting sound loop");
        play_loop_start (&manager->priv->critical_alert_timeout_id);
        break;

    default:
        /* play the sound, using sounds from the naming spec */
        ca_context_play (ca_gtk_context_get (), 0,
                         CA_PROP_EVENT_ID, "battery-caution",
                         /* TRANSLATORS: this is the sound description */
                         CA_PROP_EVENT_DESCRIPTION, _("Battery is critically low"), NULL);
        break;
    }
out:
    if (icon != NULL)
        g_object_unref (icon);
    g_free (message);
}

static void
engine_charge_action (GsdPowerManager *manager, UpDevice *device)
{
    const gchar *title = NULL;
    gchar *message = NULL;
    GIcon *icon = NULL;
    GsdPowerActionType policy;
    guint timer_id;
    UpDeviceKind kind;

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  NULL);

    /* check to see if the batteries have not noticed we are on AC */
    if (kind == UP_DEVICE_KIND_BATTERY)
    {
        if (!up_client_get_on_battery (manager->priv->up_client))
        {
            g_warning ("ignoring critically low message as we are not on battery power");
            goto out;
        }
    }

    if (kind == UP_DEVICE_KIND_BATTERY)
    {

        /* TRANSLATORS: laptop battery is really, really, low */
        title = _("Laptop battery critically low");

        /* we have to do different warnings depending on the policy */
        policy = manager_critical_action_get (manager, FALSE);

        /* use different text for different actions */
        if (policy == GSD_POWER_ACTION_NOTHING)
        {
            /* TRANSLATORS: computer will shutdown without saving data */
            message = g_strdup (_("The battery is below the critical level and "
                                  "this computer will <b>power-off</b> when the "
                                  "battery becomes completely empty."));

        }
        else if (policy == GSD_POWER_ACTION_SUSPEND)
        {
            /* TRANSLATORS: computer will suspend */
            message = g_strdup (_("The battery is below the critical level and "
                                  "this computer is about to suspend.\n"
                                  "<b>NOTE:</b> A small amount of power is required "
                                  "to keep your computer in a suspended state."));

        }
        else if (policy == GSD_POWER_ACTION_HIBERNATE)
        {
            /* TRANSLATORS: computer will hibernate */
            message = g_strdup (_("The battery is below the critical level and "
                                  "this computer is about to hibernate."));

        }
        else if (policy == GSD_POWER_ACTION_SHUTDOWN)
        {
            /* TRANSLATORS: computer will just shutdown */
            message = g_strdup (_("The battery is below the critical level and "
                                  "this computer is about to shutdown."));
        }

        /* wait 20 seconds for user-panic */
        timer_id = g_timeout_add_seconds (GSD_ACTION_DELAY,
                                          (GSourceFunc) manager_critical_action_do_cb,
                                          manager);
        g_source_set_name_by_id (timer_id, "[GsdPowerManager] battery critical-action");

    }
    else if (kind == UP_DEVICE_KIND_UPS)
    {
        //Uninterruptible power supply
        /* TRANSLATORS: UPS is really, really, low */
        title = _("UPS critically low");

        /* we have to do different warnings depending on the policy */
        policy = manager_critical_action_get (manager, TRUE);

        /* use different text for different actions */
        if (policy == GSD_POWER_ACTION_NOTHING)
        {
            /* TRANSLATORS: computer will shutdown without saving data */
            message = g_strdup (_("UPS is below the critical level and "
                                  "this computer will <b>power-off</b> when the "
                                  "UPS becomes completely empty."));

        }
        else if (policy == GSD_POWER_ACTION_HIBERNATE)
        {
            /* TRANSLATORS: computer will hibernate */
            message = g_strdup (_("UPS is below the critical level and "
                                  "this computer is about to hibernate."));

        }
        else if (policy == GSD_POWER_ACTION_SHUTDOWN)
        {
            /* TRANSLATORS: computer will just shutdown */
            message = g_strdup (_("UPS is below the critical level and "
                                  "this computer is about to shutdown."));
        }

        /* wait 20 seconds for user-panic */
        timer_id = g_timeout_add_seconds (GSD_ACTION_DELAY,
                                          (GSourceFunc) manager_critical_ups_action_do_cb,
                                          manager);
        g_source_set_name_by_id (timer_id, "[GsdPowerManager] ups critical-action");
    }

    /* not all types have actions */
    if (title == NULL)
        return;

    /* get correct icon */
    icon = gpm_upower_get_device_icon (device, TRUE);

    /* close any existing notification of this class */
    notify_close_if_showing (&manager->priv->notification_low);

    /* create a new notification */
    create_notification (title, message,
                         icon,
                         &manager->priv->notification_low);
    notify_notification_set_timeout (manager->priv->notification_low,
                                     NOTIFY_EXPIRES_NEVER);
    notify_notification_set_urgency (manager->priv->notification_low,
                                     NOTIFY_URGENCY_CRITICAL);
    notify_notification_set_app_name (manager->priv->notification_low, _("Power"));

    /* try to show */
    notify_notification_show (manager->priv->notification_low, NULL);

    /* play the sound, using sounds from the naming spec */
    ca_context_play (ca_gtk_context_get (), 0,
                     CA_PROP_EVENT_ID, "battery-caution",
                     /* TRANSLATORS: this is the sound description */
                     CA_PROP_EVENT_DESCRIPTION, _("Battery is critically low"), NULL);
out:
    if (icon != NULL)
        g_object_unref (icon);
    g_free (message);
}

/**
 * callback for upower device changes,such as power percentage drops
 */
static void
engine_device_changed_cb (UpClient *client, UpDevice *device, GsdPowerManager *manager)
{
    UpDeviceKind kind;
    UpDeviceState state;
    UpDeviceState state_old;
    GsdPowerManagerWarning warning_old;
    GsdPowerManagerWarning warning;

    /* get device properties */
    g_object_get (device,
                  "kind", &kind,
                  NULL);

    /* if battery then use composite device to cope with multiple batteries */
    if (kind == UP_DEVICE_KIND_BATTERY)
    {
        g_debug ("updating because %s changed", up_device_get_object_path (device));
        device = engine_update_composite_device (manager, device);
    }

    /* get device properties (may be composite) */
    g_object_get (device,
                  "state", &state,
                  NULL);

    g_debug ("%s state is now %s", up_device_get_object_path (device), up_device_state_to_string (state));

    /* see if any interesting state changes have happened */
    state_old = GPOINTER_TO_INT(g_object_get_data (G_OBJECT(device), "engine-state-old"));
    if (state_old != state)
    {
        if (state == UP_DEVICE_STATE_DISCHARGING)
        {
            g_debug ("discharging");
            engine_ups_discharging (manager, device);
        }
        else if (state == UP_DEVICE_STATE_FULLY_CHARGED ||
                 state == UP_DEVICE_STATE_CHARGING)
        {
            g_debug ("fully charged or charging, hiding notifications if any");
            notify_close_if_showing (&manager->priv->notification_low);
            notify_close_if_showing (&manager->priv->notification_ups_discharging);
            main_battery_or_ups_low_changed (manager, FALSE);
        }

        /* save new state */
        g_object_set_data (G_OBJECT(device), "engine-state-old", GUINT_TO_POINTER(state));
    }

    /* check the warning state has not changed */
    warning_old = GPOINTER_TO_INT(g_object_get_data (G_OBJECT(device), "engine-warning-old"));
    warning = engine_get_warning (manager, device);
    if (warning != warning_old)
        /*if (1)*/
    {
        if (warning == WARNING_LOW)
        {
            g_debug ("** EMIT: charge-low");
            engine_charge_low (manager, device);
        }
        else if (warning == WARNING_CRITICAL)
        {
            g_debug ("** EMIT: charge-critical");
            engine_charge_critical (manager, device);
        }
        else if (warning == WARNING_ACTION)
        {
            g_debug ("charge-action");
            engine_charge_action (manager, device);
        }
        /* save new state */
        g_object_set_data (G_OBJECT(device), "engine-warning-old", GUINT_TO_POINTER(warning));
    }

    engine_recalculate_state (manager);
}

static UpDevice *
engine_get_primary_device (GsdPowerManager *manager)
{
    guint i;
    UpDevice *device = NULL;
    UpDevice *device_tmp;
    UpDeviceKind kind;
    UpDeviceState state;
    gboolean is_present;

    for (i = 0; i < manager->priv->devices_array->len; i++)
    {
        device_tmp = g_ptr_array_index (manager->priv->devices_array, i);

        /* get device properties */
        g_object_get (device_tmp,
                      "kind", &kind,
                      "state", &state,
                      "is-present", &is_present,
                      NULL);

        /* not present */
        if (!is_present)
            continue;

        /* not discharging */
        if (state != UP_DEVICE_STATE_DISCHARGING)
            continue;

        /* not battery */
        if (kind != UP_DEVICE_KIND_BATTERY)
            continue;

        /* use composite device to cope with multiple batteries */
        device = g_object_ref (engine_get_composite_device (manager, device_tmp));
        break;
    }
    return device;
}

static void
gnome_session_shutdown_cb (GObject *source_object,
                           GAsyncResult *res,
                           gpointer user_data)
{
    GVariant *result;
    GError *error = NULL;

    result = g_dbus_proxy_call_finish (G_DBUS_PROXY (source_object),
                                       res,
                                       &error);
    if (result == NULL)
    {
        g_warning ("couldn't shutdown using gnome-session: %s",
                   error->message);
        g_error_free (error);
    }
    else
    {
        g_variant_unref (result);
    }
}

static void
gnome_session_shutdown (GsdPowerManager *manager)
{
    g_dbus_proxy_call (manager->priv->session,
                       "Shutdown",
                       NULL,
                       G_DBUS_CALL_FLAGS_NONE,
                       -1, NULL,
                       gnome_session_shutdown_cb, NULL);
}

static void
gnome_session_logout_cb (GObject *source_object,
                         GAsyncResult *res,
                         gpointer user_data)
{
    GVariant *result;
    GError *error = NULL;

    result = g_dbus_proxy_call_finish (G_DBUS_PROXY (source_object),
                                       res,
                                       &error);
    if (result == NULL)
    {
        g_warning ("couldn't log out using gnome-session: %s",
                   error->message);
        g_error_free (error);
    }
    else
    {
        g_variant_unref (result);
    }
}

static void
gnome_session_logout (GsdPowerManager *manager,
                      guint            logout_mode)
{
    g_dbus_proxy_call (manager->priv->session,
                       "Logout",
                       g_variant_new ("(u)", logout_mode),
                       G_DBUS_CALL_FLAGS_NONE,
                       -1, NULL,
                       gnome_session_logout_cb, NULL);
}

static void
action_poweroff (GsdPowerManager *manager)
{
    if (manager->priv->logind_proxy == NULL)
    {
        g_warning ("no systemd support");
        return;
    }
    g_dbus_proxy_call (manager->priv->logind_proxy,
                       "PowerOff",
                       g_variant_new ("(b)", FALSE),
                       G_DBUS_CALL_FLAGS_NONE,
                       G_MAXINT,
                       NULL,
                       NULL,
                       NULL);
}

static void
action_suspend (GsdPowerManager *manager)
{
    if (manager->priv->logind_proxy == NULL)
    {
        g_warning ("no systemd support");
        return;
    }
    g_dbus_proxy_call (manager->priv->logind_proxy,
                       "Suspend",
                       g_variant_new ("(b)", FALSE),
                       G_DBUS_CALL_FLAGS_NONE,
                       G_MAXINT,
                       NULL,
                       NULL,
                       NULL);
}

static void
action_hibernate (GsdPowerManager *manager)
{
    if (manager->priv->logind_proxy == NULL)
    {
        g_warning ("no systemd support");
        return;
    }
    g_dbus_proxy_call (manager->priv->logind_proxy,
                       "Hibernate",
                       g_variant_new ("(b)", FALSE),
                       G_DBUS_CALL_FLAGS_NONE,
                       G_MAXINT,
                       NULL,
                       NULL,
                       NULL);
}

static void
backlight_enable (GsdPowerManager *manager)
{
    gboolean ret;
    GError *error = NULL;

    ret = gnome_rr_screen_set_dpms_mode (manager->priv->rr_screen,
                                         GNOME_RR_DPMS_ON,
                                         &error);
    if (!ret)
    {
        g_warning ("failed to turn the panel on: %s",
                   error->message);
        g_error_free (error);
    }

    g_debug ("TESTSUITE: Unblanked screen");
}

static void
backlight_disable (GsdPowerManager *manager)
{
    gboolean ret;
    GError *error = NULL;

    ret = gnome_rr_screen_set_dpms_mode (manager->priv->rr_screen,
                                         GNOME_RR_DPMS_OFF,
                                         &error);
    if (!ret)
    {
        g_warning ("failed to turn the panel off: %s",
                   error->message);
        g_error_free (error);
    }
    g_debug ("TESTSUITE: Blanked screen");
}

static void
do_power_action_type (GsdPowerManager *manager,
                      GsdPowerActionType action_type)
{
    switch (action_type)
    {
    case GSD_POWER_ACTION_SUSPEND:
        action_suspend (manager);
        break;
    case GSD_POWER_ACTION_INTERACTIVE:
        gnome_session_shutdown (manager);
        break;
    case GSD_POWER_ACTION_HIBERNATE:
        action_hibernate (manager);
        break;
    case GSD_POWER_ACTION_SHUTDOWN:
        /* this is only used on critically low battery where
         * hibernate is not available and is marginally better
         * than just powering down the computer mid-write */
        action_poweroff (manager);
        break;
    case GSD_POWER_ACTION_BLANK:
        backlight_disable (manager);
        break;
    case GSD_POWER_ACTION_NOTHING:
        break;
    case GSD_POWER_ACTION_LOGOUT:
        gnome_session_logout (manager, GSM_MANAGER_LOGOUT_MODE_FORCE);
        break;
    }
}

static GsmInhibitorFlag
get_idle_inhibitors_for_action (GsdPowerActionType action_type)
{
    switch (action_type)
    {
    case GSD_POWER_ACTION_BLANK:
    case GSD_POWER_ACTION_SHUTDOWN:
    case GSD_POWER_ACTION_INTERACTIVE:
        return GSM_INHIBITOR_FLAG_IDLE;
    case GSD_POWER_ACTION_HIBERNATE:
    case GSD_POWER_ACTION_SUSPEND:
        return GSM_INHIBITOR_FLAG_SUSPEND; /* in addition to idle */
    case GSD_POWER_ACTION_NOTHING:
        return 0;
    case GSD_POWER_ACTION_LOGOUT:
        return GSM_INHIBITOR_FLAG_LOGOUT; /* in addition to idle */
    }
    return 0;
}

static gboolean
is_action_inhibited (GsdPowerManager *manager, GsdPowerActionType action_type)
{
    GsmInhibitorFlag flag;
    gboolean is_inhibited;

    flag = get_idle_inhibitors_for_action (action_type);
    if (!flag)
        return FALSE;
    idle_is_session_inhibited (manager,
                               flag,
                               &is_inhibited);
    return is_inhibited;
}

static gboolean
upower_kbd_set_brightness (GsdPowerManager *manager, guint value, GError **error)
{
    GVariant *retval;

    /* same as before */
    if (manager->priv->kbd_brightness_now == value)
        return TRUE;

    /* update h/w value */
    retval = g_dbus_proxy_call_sync (manager->priv->upower_kdb_proxy,
                                     "SetBrightness",
                                     g_variant_new ("(i)", (gint) value),
                                     G_DBUS_CALL_FLAGS_NONE,
                                     -1,
                                     NULL,
                                     error);
    if (retval == NULL)
        return FALSE;

    /* save new value */
    manager->priv->kbd_brightness_now = value;
    g_variant_unref (retval);
    return TRUE;
}

static gboolean
upower_kbd_toggle (GsdPowerManager *manager,
                   GError **error)
{
    gboolean ret;

    if (manager->priv->kbd_brightness_old >= 0)
    {
        g_debug ("keyboard toggle off");
        ret = upower_kbd_set_brightness (manager,
                                         manager->priv->kbd_brightness_old,
                                         error);
        if (ret)
        {
            /* succeeded, set to -1 since now no old value */
            manager->priv->kbd_brightness_old = -1;
        }
    }
    else
    {
        g_debug ("keyboard toggle on");
        /* save the current value to restore later when untoggling */
        manager->priv->kbd_brightness_old = manager->priv->kbd_brightness_now;
        ret = upower_kbd_set_brightness (manager, 0, error);
        if (!ret)
        {
            /* failed, reset back to -1 */
            manager->priv->kbd_brightness_old = -1;
        }
    }

    return ret;
}

static gboolean
suspend_on_lid_close (GsdPowerManager *manager)
{
    GsdXrandrBootBehaviour val;

    if (!external_monitor_is_connected (manager->priv->rr_screen))
        return TRUE;

    val = g_settings_get_enum(manager->priv->settings, "lid-close-ac-action");
    return val == GSD_POWER_ACTION_SUSPEND;
    /*val = g_settings_get_enum (manager->priv->settings_xrandr, "default-monitors-setup");*/
    /*return val == GSD_XRANDR_BOOT_BEHAVIOUR_DO_NOTHING;*/
}

static gboolean
inhibit_lid_switch_timer_cb (GsdPowerManager *manager)
{
    if (suspend_on_lid_close (manager))
    {
        g_debug ("no external monitors for a while; uninhibiting lid close");
        uninhibit_lid_switch (manager);
        manager->priv->inhibit_lid_switch_timer_id = 0;
        return G_SOURCE_REMOVE;
    }

    g_debug ("external monitor still there; trying again later");
    return G_SOURCE_CONTINUE;
}

/* Sets up a timer to be triggered some seconds after closing the laptop lid
 * when the laptop is *not* suspended for some reason.  We'll check conditions
 * again in the timeout handler to see if we can suspend then.
 */
static void
setup_inhibit_lid_switch_timer (GsdPowerManager *manager)
{
    if (manager->priv->inhibit_lid_switch_timer_id != 0)
    {
        g_debug ("lid close safety timer already set up");
        return;
    }

    g_debug ("setting up lid close safety timer");

    manager->priv->inhibit_lid_switch_timer_id = g_timeout_add_seconds (GSD_POWER_MANAGER_LID_CLOSE_SAFETY_TIMEOUT,
            (GSourceFunc) inhibit_lid_switch_timer_cb,
            manager);
    g_source_set_name_by_id (manager->priv->inhibit_lid_switch_timer_id, "[GsdPowerManager] lid close safety timer");
}

static void
restart_inhibit_lid_switch_timer (GsdPowerManager *manager)
{
    if (manager->priv->inhibit_lid_switch_timer_id != 0)
    {
        g_debug ("restarting lid close safety timer");
        g_source_remove (manager->priv->inhibit_lid_switch_timer_id);
        manager->priv->inhibit_lid_switch_timer_id = 0;
        setup_inhibit_lid_switch_timer (manager);
    }
}

static void
do_lid_open_action (GsdPowerManager *manager)
{
    /* play a sound, using sounds from the naming spec */
    ca_context_play (ca_gtk_context_get (), 0,
                     CA_PROP_EVENT_ID, "lid-open",
                     /* TRANSLATORS: this is the sound description */
                     CA_PROP_EVENT_DESCRIPTION, _("Lid has been opened"),
                     NULL);

    /* This might already have happened when resuming, but
     * if we didn't sleep, we'll need to wake it up */
    reset_idletime ();
}

static void
lock_screensaver (GsdPowerManager *manager)
{
    gboolean do_lock;

    do_lock = g_settings_get_boolean (manager->priv->settings_screensaver,
                                      "lock-enabled");
    if (!do_lock)
    {
        g_dbus_proxy_call_sync (manager->priv->screensaver_proxy,
                                "SetActive",
                                g_variant_new ("(b)", TRUE),
                                G_DBUS_CALL_FLAGS_NONE,
                                -1, NULL, NULL);
        return;
    }

    g_dbus_proxy_call_sync (manager->priv->screensaver_proxy,
                            "Lock",
                            NULL,
                            G_DBUS_CALL_FLAGS_NONE,
                            -1, NULL, NULL);
}

static void
do_lid_closed_action (GsdPowerManager *manager)
{
    /* play a sound, using sounds from the naming spec */
    ca_context_play (ca_gtk_context_get (), 0,
                     CA_PROP_EVENT_ID, "lid-close",
                     /* TRANSLATORS: this is the sound description */
                     CA_PROP_EVENT_DESCRIPTION, _("Lid has been closed"),
                     NULL);

    /* refresh RANDR so we get an accurate view of what monitors are plugged in when the lid is closed */
    gnome_rr_screen_refresh (manager->priv->rr_screen, NULL); /* NULL-GError */

    restart_inhibit_lid_switch_timer (manager);

    if (suspend_on_lid_close (manager))
    {
        gboolean is_inhibited;

        idle_is_session_inhibited (manager,
                                   GSM_INHIBITOR_FLAG_SUSPEND,
                                   &is_inhibited);
        if (is_inhibited)
        {
            g_debug ("Suspend is inhibited but lid is closed, locking the screen");
            /* We put the screensaver on * as we're not suspending,
             * but the lid is closed */
            lock_screensaver (manager);
        }
    }
}

static void
up_client_changed_cb (UpClient *client, GsdPowerManager *manager)
{
    gboolean tmp;

    if (!up_client_get_on_battery (client))
    {
        /* if we are playing a critical charge sound loop on AC, stop it */
        play_loop_stop (&manager->priv->critical_alert_timeout_id);
        notify_close_if_showing (&manager->priv->notification_low);
        main_battery_or_ups_low_changed (manager, FALSE);
    }

    /* same state */
    tmp = up_client_get_lid_is_closed (manager->priv->up_client);
    if (manager->priv->lid_is_closed == tmp)
        return;
    manager->priv->lid_is_closed = tmp;
    g_debug ("up changed: lid is now %s", tmp ? "closed" : "open");

    if (manager->priv->lid_is_closed)
        do_lid_closed_action (manager);
    else
        do_lid_open_action (manager);
}

static const gchar *
idle_mode_to_string (GsdPowerIdleMode mode)
{
    if (mode == GSD_POWER_IDLE_MODE_NORMAL)
        return "normal";
    if (mode == GSD_POWER_IDLE_MODE_DIM)
        return "dim";
    if (mode == GSD_POWER_IDLE_MODE_BLANK)
        return "blank";
    if (mode == GSD_POWER_IDLE_MODE_SLEEP)
        return "sleep";
    return "unknown";
}

static const char *
idle_watch_id_to_string (GsdPowerManager *manager, guint id)
{
    if (id == manager->priv->idle_dim_id)
        return "dim";
    if (id == manager->priv->idle_blank_id)
        return "blank";
    if (id == manager->priv->idle_sleep_id)
        return "sleep";
    if (id == manager->priv->idle_sleep_warning_id)
        return "sleep-warning";
    return NULL;
}

static void
backlight_emit_changed (GsdPowerManager *manager)
{
    gboolean ret;
    GError *error = NULL;

    /* not yet connected to the bus */
    if (manager->priv->connection == NULL)
        return;
    ret = g_dbus_connection_emit_signal (manager->priv->connection,
                                         NULL,
                                         GSD_POWER_DBUS_PATH,
                                         GSD_POWER_DBUS_INTERFACE_SCREEN,
                                         "Changed",
                                         NULL,
                                         &error);
    if (!ret)
    {
        g_warning ("failed to emit Changed: %s", error->message);
        g_error_free (error);
    }
}

static gboolean
display_backlight_dim (GsdPowerManager *manager,
                       gint idle_percentage,
                       GError **error)
{
    gint min;
    gint max;
    gint now;
    gint idle;
    gboolean ret = FALSE;

    if (!manager->priv->backlight_available)
        return TRUE;

    now = backlight_get_abs (manager->priv->rr_screen, error);
    if (now < 0)
    {
        goto out;
    }

    /* is the dim brightness actually *dimmer* than the
     * brightness we have now? */
    min = backlight_get_min (manager->priv->rr_screen);
    max = backlight_get_max (manager->priv->rr_screen, error);
    if (max < 0)
    {
        goto out;
    }
    idle = PERCENTAGE_TO_ABS (min, max, idle_percentage);
    if (idle > now)
    {
        g_debug ("brightness already now %i/%i, so "
                 "ignoring dim to %i/%i",
                 now, max, idle, max);
        ret = TRUE;
        goto out;
    }
    ret = backlight_set_abs (manager->priv->rr_screen,
                             idle,
                             error);
    if (!ret)
    {
        goto out;
    }

    /* save for undim */
    manager->priv->pre_dim_brightness = now;

out:
    return ret;
}

static gboolean
kbd_backlight_dim (GsdPowerManager *manager,
                   gint idle_percentage,
                   GError **error)
{
    gboolean ret;
    gint idle;
    gint max;
    gint now;

    if (manager->priv->upower_kdb_proxy == NULL)
        return TRUE;

    now = manager->priv->kbd_brightness_now;
    max = manager->priv->kbd_brightness_max;
    idle = PERCENTAGE_TO_ABS (0, max, idle_percentage);
    if (idle > now)
    {
        g_debug ("kbd brightness already now %i/%i, so "
                 "ignoring dim to %i/%i",
                 now, max, idle, max);
        return TRUE;
    }
    ret = upower_kbd_set_brightness (manager, idle, error);
    if (!ret)
        return FALSE;

    /* save for undim */
    manager->priv->kbd_brightness_pre_dim = now;
    return TRUE;
}

static gboolean
is_session_active (GsdPowerManager *manager)
{
    GVariant *variant;
    gboolean is_session_active = FALSE;

    variant = g_dbus_proxy_get_cached_property (manager->priv->session,
              "SessionIsActive");
    if (variant)
    {
        is_session_active = g_variant_get_boolean (variant);
        g_variant_unref (variant);
    }

    return is_session_active;
}

static void
idle_set_mode (GsdPowerManager *manager, GsdPowerIdleMode mode)
{
    gboolean ret = FALSE;
    GError *error = NULL;
    gint idle_percentage;
    GsdPowerActionType action_type;
    gboolean is_active = FALSE;

    /* Ignore attempts to set "less idle" modes */
    if (mode <= manager->priv->current_idle_mode &&
            mode != GSD_POWER_IDLE_MODE_NORMAL)
    {
        g_debug ("Not going to 'less idle' mode %s (current: %s)",
                 idle_mode_to_string (mode),
                 idle_mode_to_string (manager->priv->current_idle_mode));
        return;
    }

    /* ensure we're still on an active console */
    is_active = is_session_active (manager);
    if (!is_active)
    {
        g_debug ("ignoring state transition to %s as inactive",
                 idle_mode_to_string (mode));
        return;
    }

    /* don't do any power saving if we're a VM */
    if (manager->priv->is_virtual_machine)
    {
        g_debug ("ignoring state transition to %s as virtual machine",
                 idle_mode_to_string (mode));
        return;
    }

    manager->priv->current_idle_mode = mode;
    g_debug ("Doing a state transition: %s", idle_mode_to_string (mode));

    /* if we're moving to an idle mode, make sure
     * we add a watch to take us back to normal */
    if (mode != GSD_POWER_IDLE_MODE_NORMAL)
    {
        gnome_idle_monitor_add_user_active_watch (manager->priv->idle_monitor,
                idle_became_active_cb,
                manager,
                NULL);
    }

    /* save current brightness, and set dim level */
    if (mode == GSD_POWER_IDLE_MODE_DIM)
    {
        /* display backlight */
        idle_percentage = g_settings_get_int (manager->priv->settings,
                                              "idle-brightness");
        ret = display_backlight_dim (manager, idle_percentage, &error);
        if (!ret)
        {
            g_warning ("failed to set dim backlight to %i%%: %s",
                       idle_percentage,
                       error->message);
            g_clear_error (&error);
        }

        /* keyboard backlight */
        ret = kbd_backlight_dim (manager, idle_percentage, &error);
        if (!ret)
        {
            g_warning ("failed to set dim kbd backlight to %i%%: %s",
                       idle_percentage,
                       error->message);
            g_clear_error (&error);
        }

        /* turn off screen and kbd */
    }
    else if (mode == GSD_POWER_IDLE_MODE_BLANK)
    {

        backlight_disable (manager);

        /* only toggle keyboard if present and not already toggled */
        if (manager->priv->upower_kdb_proxy &&
                manager->priv->kbd_brightness_old == -1)
        {
            ret = upower_kbd_toggle (manager, &error);
            if (!ret)
            {
                g_warning ("failed to turn the kbd backlight off: %s",
                           error->message);
                g_error_free (error);
            }
        }

        /* sleep */
    }
    else if (mode == GSD_POWER_IDLE_MODE_SLEEP)
    {

        if (up_client_get_on_battery (manager->priv->up_client))
        {
            action_type = g_settings_get_enum (manager->priv->settings,
                                               "sleep-inactive-battery-type");
        }
        else
        {
            action_type = g_settings_get_enum (manager->priv->settings,
                                               "sleep-inactive-ac-type");
        }
        do_power_action_type (manager, action_type);

        /* turn on screen and restore user-selected brightness level */
    }
    else if (mode == GSD_POWER_IDLE_MODE_NORMAL)
    {

        backlight_enable (manager);

        /* reset brightness if we dimmed */
        if (manager->priv->pre_dim_brightness >= 0)
        {
            ret = backlight_set_abs (manager->priv->rr_screen,
                                     manager->priv->pre_dim_brightness,
                                     &error);
            if (!ret)
            {
                g_warning ("failed to restore backlight to %i: %s",
                           manager->priv->pre_dim_brightness,
                           error->message);
                g_clear_error (&error);
            }
            else
            {
                manager->priv->pre_dim_brightness = -1;
            }
        }

        /* only toggle keyboard if present and already toggled off */
        if (manager->priv->upower_kdb_proxy &&
                manager->priv->kbd_brightness_old != -1)
        {
            ret = upower_kbd_toggle (manager, &error);
            if (!ret)
            {
                g_warning ("failed to turn the kbd backlight on: %s",
                           error->message);
                g_clear_error (&error);
            }
        }

        /* reset kbd brightness if we dimmed */
        if (manager->priv->kbd_brightness_pre_dim >= 0)
        {
            ret = upower_kbd_set_brightness (manager,
                                             manager->priv->kbd_brightness_pre_dim,
                                             &error);
            if (!ret)
            {
                g_warning ("failed to restore kbd backlight to %i: %s",
                           manager->priv->kbd_brightness_pre_dim,
                           error->message);
                g_error_free (error);
            }
            manager->priv->kbd_brightness_pre_dim = -1;
        }

    }
}

static gboolean
idle_is_session_inhibited (GsdPowerManager  *manager,
                           GsmInhibitorFlag  mask,
                           gboolean         *is_inhibited)
{
    GVariant *variant;
    GsmInhibitorFlag inhibited_actions;

    /* not yet connected to gnome-session */
    if (manager->priv->session == NULL)
    {
        g_debug("not yet connected to gnome-session");
        return FALSE;
    }

    variant = g_dbus_proxy_get_cached_property (manager->priv->session,
              "InhibitedActions");
    if (!variant)
        return FALSE;

    inhibited_actions = g_variant_get_uint32 (variant);
    g_variant_unref (variant);

    *is_inhibited = (inhibited_actions & mask);

    return TRUE;
}

static void
clear_idle_watch (GnomeIdleMonitor *monitor,
                  guint            *id)
{
    if (*id == 0)
        return;
    gnome_idle_monitor_remove_watch (monitor, *id);
    *id = 0;
}

static void
idle_configure (GsdPowerManager *manager)
{
    gboolean is_idle_inhibited;
    GsdPowerActionType action_type;
    guint timeout_blank;
    guint timeout_sleep;
    guint timeout_dim;
    gboolean on_battery;

    if (!idle_is_session_inhibited (manager,
                                    GSM_INHIBITOR_FLAG_IDLE,
                                    &is_idle_inhibited))
    {
        /* Session isn't available yet, postpone */
        g_debug("Session isn't available yet,postpone\n");
        return;
    }

    /* are we inhibited from going idle */
    if (!is_session_active (manager) || is_idle_inhibited)
    {
        g_debug ("inhibited or inactive, so using normal state");
        idle_set_mode (manager, GSD_POWER_IDLE_MODE_NORMAL);

        clear_idle_watch (manager->priv->idle_monitor,
                          &manager->priv->idle_blank_id);
        clear_idle_watch (manager->priv->idle_monitor,
                          &manager->priv->idle_sleep_id);
        clear_idle_watch (manager->priv->idle_monitor,
                          &manager->priv->idle_dim_id);
        clear_idle_watch (manager->priv->idle_monitor,
                          &manager->priv->idle_sleep_warning_id);
        notify_close_if_showing (&manager->priv->notification_sleep_warning);
        return;
    }

    /* set up blank callback only when the screensaver is on,
     * as it's what will drive the blank */
    on_battery = up_client_get_on_battery (manager->priv->up_client);
    timeout_blank = 0;
    if (manager->priv->screensaver_active)
    {
        /* The tail is wagging the dog.
         * The screensaver coming on will blank the screen.
         * If an event occurs while the screensaver is on,
         * the aggressive idle watch will handle it */
        timeout_blank = SCREENSAVER_TIMEOUT_BLANK;
    }

    clear_idle_watch (manager->priv->idle_monitor,
                      &manager->priv->idle_blank_id);

    if (timeout_blank != 0)
    {
        g_debug ("setting up blank callback for %is", timeout_blank);

        manager->priv->idle_blank_id = gnome_idle_monitor_add_idle_watch (manager->priv->idle_monitor,
                                       timeout_blank * 1000,
                                       idle_triggered_idle_cb, manager, NULL);
    }

    /* only do the sleep timeout when the session is idle
     * and we aren't inhibited from sleeping (or logging out, etc.) */
    action_type = g_settings_get_enum (manager->priv->settings, on_battery ?
                                       "sleep-inactive-battery-type" : "sleep-inactive-ac-type");
    timeout_sleep = 0;
    if (!is_action_inhibited (manager, action_type))
    {
        timeout_sleep = g_settings_get_int (manager->priv->settings, on_battery ?
                                            "sleep-inactive-battery-timeout" : "sleep-inactive-ac-timeout");
    }

    clear_idle_watch (manager->priv->idle_monitor,
                      &manager->priv->idle_sleep_id);
    clear_idle_watch (manager->priv->idle_monitor,
                      &manager->priv->idle_sleep_warning_id);

    g_debug ("setting up sleep callback %is", timeout_sleep);
    if (timeout_sleep != 0)
    {
        g_debug ("setting up sleep callback %is", timeout_sleep);

        manager->priv->idle_sleep_id = gnome_idle_monitor_add_idle_watch (manager->priv->idle_monitor,
                                       timeout_sleep * 1000,
                                       idle_triggered_idle_cb, manager, NULL);
        if (action_type == GSD_POWER_ACTION_LOGOUT ||
                action_type == GSD_POWER_ACTION_SUSPEND ||
                action_type == GSD_POWER_ACTION_HIBERNATE)
        {
            guint timeout_sleep_warning;

            manager->priv->sleep_action_type = action_type;
            timeout_sleep_warning = timeout_sleep * IDLE_DELAY_TO_IDLE_DIM_MULTIPLIER;
            if (timeout_sleep_warning < MINIMUM_IDLE_DIM_DELAY)
                timeout_sleep_warning = 0;

            g_debug ("setting up sleep warning callback %is", timeout_sleep_warning);

            manager->priv->idle_sleep_warning_id = gnome_idle_monitor_add_idle_watch (manager->priv->idle_monitor,
                                                   timeout_sleep_warning * 1000,
                                                   idle_triggered_idle_cb, manager, NULL);
        }
    }

    if (manager->priv->idle_sleep_warning_id == 0)
        notify_close_if_showing (&manager->priv->notification_sleep_warning);

    /* set up dim callback for when the screen lock is not active,
     * but only if we actually want to dim. */
    timeout_dim = 0;
    if (manager->priv->screensaver_active)
    {
        /* Don't dim when the screen lock is active */
    }
    else if (!on_battery)
    {
        /* Don't dim when charging */
    }
    else if (manager->priv->battery_is_low)
    {
        /* Aggressively blank when battery is low */
        timeout_dim = SCREENSAVER_TIMEOUT_BLANK;
    }
    else
    {
        if (g_settings_get_boolean (manager->priv->settings, "idle-dim"))
        {
            timeout_dim = g_settings_get_uint (manager->priv->settings_session,
                                               "idle-delay");
            if (timeout_dim == 0)
            {
                timeout_dim = IDLE_DIM_BLANK_DISABLED_MIN;
            }
            else
            {
                timeout_dim *= IDLE_DELAY_TO_IDLE_DIM_MULTIPLIER;
                /* Don't bother dimming if the idle-delay is
                 * too low, we'll do that when we bring down the
                 * screen lock */
                if (timeout_dim < MINIMUM_IDLE_DIM_DELAY)
                    timeout_dim = 0;
            }
        }
    }

    clear_idle_watch (manager->priv->idle_monitor,
                      &manager->priv->idle_dim_id);

    if (timeout_dim != 0)
    {
        g_debug ("setting up dim callback for %is", timeout_dim);

        manager->priv->idle_dim_id = gnome_idle_monitor_add_idle_watch (manager->priv->idle_monitor,
                                     timeout_dim * 1000,
                                     idle_triggered_idle_cb, manager, NULL);
    }
}

static void
main_battery_or_ups_low_changed (GsdPowerManager *manager,
                                 gboolean         is_low)
{
    if (is_low == manager->priv->battery_is_low)
        return;
    manager->priv->battery_is_low = is_low;
    g_debug("main_battery_or_ups_low_changed():\n");
    idle_configure (manager);
}

static gboolean
temporary_unidle_done_cb (GsdPowerManager *manager)
{
    idle_set_mode (manager, manager->priv->previous_idle_mode);
    manager->priv->temporary_unidle_on_ac_id = 0;
    return FALSE;
}

static void
set_temporary_unidle_on_ac (GsdPowerManager *manager,
                            gboolean         enable)
{
    if (!enable)
    {
        if (manager->priv->temporary_unidle_on_ac_id != 0)
        {
            g_source_remove (manager->priv->temporary_unidle_on_ac_id);
            manager->priv->temporary_unidle_on_ac_id = 0;

            idle_set_mode (manager, manager->priv->previous_idle_mode);
        }
    }
    else
    {
        /* Don't overwrite the previous idle mode when an unidle is
         * already on-going */
        if (manager->priv->temporary_unidle_on_ac_id != 0)
        {
            g_source_remove (manager->priv->temporary_unidle_on_ac_id);
        }
        else
        {
            manager->priv->previous_idle_mode = manager->priv->current_idle_mode;
            idle_set_mode (manager, GSD_POWER_IDLE_MODE_NORMAL);
        }
        manager->priv->temporary_unidle_on_ac_id = g_timeout_add_seconds (POWER_UP_TIME_ON_AC,
                (GSourceFunc) temporary_unidle_done_cb,
                manager);
    }
}

static void
up_client_on_battery_cb (UpClient *client,
                         GParamSpec *pspec,
                         GsdPowerManager *manager)
{
    g_debug("up_client_on_battery_cb()\n");
    idle_configure (manager);

    if (manager->priv->lid_is_closed)
        return;

    if (manager->priv->current_idle_mode == GSD_POWER_IDLE_MODE_BLANK ||
            manager->priv->current_idle_mode == GSD_POWER_IDLE_MODE_DIM ||
            manager->priv->temporary_unidle_on_ac_id != 0)
        set_temporary_unidle_on_ac (manager, TRUE);
}

static void
gsd_power_manager_finalize (GObject *object)
{
    GsdPowerManager *manager;

    g_return_if_fail (object != NULL);
    g_return_if_fail (GSD_IS_POWER_MANAGER (object));

    manager = GSD_POWER_MANAGER (object);

    g_return_if_fail (manager->priv != NULL);

    g_clear_object (&manager->priv->connection);

    if (manager->priv->name_id != 0)
        g_bus_unown_name (manager->priv->name_id);

    G_OBJECT_CLASS (gsd_power_manager_parent_class)->finalize (object);
}

static void
gsd_power_manager_class_init (GsdPowerManagerClass *klass)
{
    GObjectClass *object_class = G_OBJECT_CLASS (klass);

    object_class->finalize = gsd_power_manager_finalize;

    g_type_class_add_private (klass, sizeof (GsdPowerManagerPrivate));
}

static void
session_presence_proxy_ready_cb (GObject *source_object,
                                 GAsyncResult *res,
                                 gpointer user_data)
{
    GError *error = NULL;
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);

    manager->priv->session_presence_proxy = g_dbus_proxy_new_for_bus_finish (res, &error);
    if (manager->priv->session_presence_proxy == NULL)
    {
        g_warning ("Could not connect to gnome-sesson: %s",
                   error->message);
        g_error_free (error);
        return;
    }
}

static void
handle_screensaver_active (GsdPowerManager *manager,
                           GVariant        *parameters)
{
    gboolean active;

    g_variant_get (parameters, "(b)", &active);
    g_debug ("Received screensaver ActiveChanged signal: %d (old: %d)", active, manager->priv->screensaver_active);
    if (manager->priv->screensaver_active != active)
    {
        manager->priv->screensaver_active = active;
        g_debug("handle_screensaver_active():\n");
        idle_configure (manager);

        /* Setup blank as soon as the screensaver comes on,
         * and its fade has finished.
         *
         * See also idle_configure() */
        if (active)
            idle_set_mode (manager, GSD_POWER_IDLE_MODE_BLANK);
    }
}

static void
screensaver_signal_cb (GDBusProxy *proxy,
                       const gchar *sender_name,
                       const gchar *signal_name,
                       GVariant *parameters,
                       gpointer user_data)
{
    if (g_strcmp0 (signal_name, "ActiveChanged") == 0)
        handle_screensaver_active (GSD_POWER_MANAGER (user_data), parameters);
}

static void
get_active_cb (GDBusProxy *proxy,
               GAsyncResult *result,
               GsdPowerManager *manager)
{
    GVariant *res;
    GError *error = NULL;

    res = g_dbus_proxy_call_finish (proxy, result, &error);
    if (!res)
    {
        g_warning ("Failed to run GetActive() function on screensaver: %s", error->message);
        g_error_free (error);
        return;
    }

    handle_screensaver_active (manager, res);
    g_variant_unref (res);
}

static void
screensaver_proxy_ready_cb (GObject         *source_object,
                            GAsyncResult    *res,
                            GsdPowerManager *manager)
{
    GError *error = NULL;
    GDBusProxy *proxy;

    proxy = g_dbus_proxy_new_finish (res, &error);

    if (proxy == NULL)
    {
        if (!g_error_matches (error, G_IO_ERROR, G_IO_ERROR_CANCELLED))
            g_warning ("Could not connect to screensaver: %s", error->message);
        g_error_free (error);
        return;
    }

    manager->priv->screensaver_proxy = proxy;

    g_signal_connect (manager->priv->screensaver_proxy, "g-signal",
                      G_CALLBACK (screensaver_signal_cb), manager);
    g_dbus_proxy_call (manager->priv->screensaver_proxy,
                       "GetActive",
                       NULL,
                       0,
                       G_MAXINT,
                       NULL,
                       (GAsyncReadyCallback)get_active_cb,
                       manager);

}

static void
screensaver_appeared_cb (GDBusConnection *connection,
                         const char      *name,
                         const char      *name_owner,
                         GsdPowerManager *manager)
{
    g_dbus_proxy_new (connection,
                      0,
                      NULL,
                      GS_DBUS_NAME,
                      GS_DBUS_PATH,
                      GS_DBUS_INTERFACE,
                      manager->priv->screensaver_cancellable,
                      (GAsyncReadyCallback) screensaver_proxy_ready_cb,
                      manager);

}

static void
screensaver_vanished_cb (GDBusConnection *connection,
                         const char      *name,
                         GsdPowerManager *manager)
{
    g_clear_object (&manager->priv->screensaver_proxy);
}

static void
power_keyboard_proxy_ready_cb (GObject             *source_object,
                               GAsyncResult        *res,
                               gpointer             user_data)
{
    GVariant *k_now = NULL;
    GVariant *k_max = NULL;
    GError *error = NULL;
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);

    manager->priv->upower_kdb_proxy = g_dbus_proxy_new_for_bus_finish (res, &error);
    if (manager->priv->upower_kdb_proxy == NULL)
    {
        g_warning ("Could not connect to UPower: %s",
                   error->message);
        g_error_free (error);
        goto out;
    }

    k_now = g_dbus_proxy_call_sync (manager->priv->upower_kdb_proxy,
                                    "GetBrightness",
                                    NULL,
                                    G_DBUS_CALL_FLAGS_NONE,
                                    -1,
                                    NULL,
                                    &error);
    if (k_now == NULL)
    {
        if (error->domain != G_DBUS_ERROR ||
                error->code != G_DBUS_ERROR_UNKNOWN_METHOD)
        {
            g_warning ("Failed to get brightness: %s",
                       error->message);
        }
        g_error_free (error);
        goto out;
    }

    k_max = g_dbus_proxy_call_sync (manager->priv->upower_kdb_proxy,
                                    "GetMaxBrightness",
                                    NULL,
                                    G_DBUS_CALL_FLAGS_NONE,
                                    -1,
                                    NULL,
                                    &error);
    if (k_max == NULL)
    {
        g_warning ("Failed to get max brightness: %s", error->message);
        g_error_free (error);
        goto out;
    }

    g_variant_get (k_now, "(i)", &manager->priv->kbd_brightness_now);
    g_variant_get (k_max, "(i)", &manager->priv->kbd_brightness_max);

    /* set brightness to max if not currently set so is something
     * sensible */
    if (manager->priv->kbd_brightness_now <= 0)
    {
        gboolean ret;
        ret = upower_kbd_set_brightness (manager,
                                         manager->priv->kbd_brightness_max,
                                         &error);
        if (!ret)
        {
            g_warning ("failed to initialize kbd backlight to %i: %s",
                       manager->priv->kbd_brightness_max,
                       error->message);
            g_error_free (error);
        }
    }
out:
    if (k_now != NULL)
        g_variant_unref (k_now);
    if (k_max != NULL)
        g_variant_unref (k_max);
}

static void
show_sleep_warning (GsdPowerManager *manager)
{
    /* close any existing notification of this class */
    notify_close_if_showing (&manager->priv->notification_sleep_warning);

    /* create a new notification */
    switch (manager->priv->sleep_action_type)
    {
    case GSD_POWER_ACTION_LOGOUT:
        create_notification (_("Automatic logout"), _("You will soon log out because of inactivity."),
                             NULL,
                             &manager->priv->notification_sleep_warning);
        break;
    case GSD_POWER_ACTION_SUSPEND:
        create_notification (_("Automatic suspend"), _("Computer will suspend very soon because of inactivity."),
                             NULL,
                             &manager->priv->notification_sleep_warning);
        break;
    case GSD_POWER_ACTION_HIBERNATE:
        create_notification (_("Automatic hibernation"), _("Computer will suspend very soon because of inactivity."),
                             NULL,
                             &manager->priv->notification_sleep_warning);
        break;
    default:
        g_assert_not_reached ();
        break;
    }
    notify_notification_set_timeout (manager->priv->notification_sleep_warning,
                                     NOTIFY_EXPIRES_NEVER);
    notify_notification_set_urgency (manager->priv->notification_sleep_warning,
                                     NOTIFY_URGENCY_CRITICAL);
    notify_notification_set_app_name (manager->priv->notification_sleep_warning, _("Power"));

    notify_notification_show (manager->priv->notification_sleep_warning, NULL);

    if (manager->priv->sleep_action_type == GSD_POWER_ACTION_LOGOUT)
        set_temporary_unidle_on_ac (manager, TRUE);
}

static void
idle_triggered_idle_cb (GnomeIdleMonitor *monitor,
                        guint             watch_id,
                        gpointer          user_data)
{
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);
    const char *id_name;

    id_name = idle_watch_id_to_string (manager, watch_id);
    if (id_name == NULL)
        g_debug ("idletime watch: %i", watch_id);
    else
        g_debug ("idletime watch: %s (%i)", id_name, watch_id);

    if (watch_id == manager->priv->idle_dim_id)
    {
        idle_set_mode (manager, GSD_POWER_IDLE_MODE_DIM);
    }
    else if (watch_id == manager->priv->idle_blank_id)
    {
        idle_set_mode (manager, GSD_POWER_IDLE_MODE_BLANK);
    }
    else if (watch_id == manager->priv->idle_sleep_id)
    {
        idle_set_mode (manager, GSD_POWER_IDLE_MODE_SLEEP);
    }
    else if (watch_id == manager->priv->idle_sleep_warning_id)
    {
        show_sleep_warning (manager);
    }
}

static void
idle_became_active_cb (GnomeIdleMonitor *monitor,
                       guint             watch_id,
                       gpointer          user_data)
{
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);

    g_debug ("idletime reset");

    set_temporary_unidle_on_ac (manager, FALSE);

    /* close any existing notification about idleness */
    notify_close_if_showing (&manager->priv->notification_sleep_warning);

    idle_set_mode (manager, GSD_POWER_IDLE_MODE_NORMAL);
}

static void
engine_settings_key_changed_cb (GSettings *settings,
                                const gchar *key,
                                GsdPowerManager *manager)
{
    if (g_strcmp0 (key, "use-time-for-policy") == 0)
    {
        manager->priv->use_time_primary = g_settings_get_boolean (settings, key);
        return;
    }
    if (g_str_has_prefix (key, "sleep-inactive") ||
            g_str_equal (key, "idle-delay") ||
            g_str_equal (key, "idle-dim"))
    {
        g_debug("engine_settings_key_changed_cb()\n");
        idle_configure (manager);
        return;
    }
}

static void
engine_session_properties_changed_cb (GDBusProxy      *session,
                                      GVariant        *changed,
                                      char           **invalidated,
                                      GsdPowerManager *manager)
{
    GVariant *v;

    v = g_variant_lookup_value (changed, "SessionIsActive", G_VARIANT_TYPE_BOOLEAN);
    if (v)
    {
        gboolean active;

        active = g_variant_get_boolean (v);
        g_debug ("Received session is active change: now %s", active ? "active" : "inactive");
        /* when doing the fast-user-switch into a new account,
         * ensure the new account is undimmed and with the backlight on */
        if (active)
            idle_set_mode (manager, GSD_POWER_IDLE_MODE_NORMAL);
        g_variant_unref (v);

    }

    v = g_variant_lookup_value (changed, "InhibitedActions", G_VARIANT_TYPE_UINT32);
    if (v)
    {
        g_variant_unref (v);
        g_debug ("Received gnome session inhibitor change");
        g_debug ("engine_session_properties_changed_cb()\n");
        idle_configure (manager);
    }
}

static void
inhibit_lid_switch_done (GObject      *source,
                         GAsyncResult *result,
                         gpointer      user_data)
{
    GDBusProxy *proxy = G_DBUS_PROXY (source);
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);
    GError *error = NULL;
    GVariant *res;
    GUnixFDList *fd_list = NULL;
    gint idx;

    res = g_dbus_proxy_call_with_unix_fd_list_finish (proxy, &fd_list, result, &error);
    if (res == NULL)
    {
        g_warning ("Unable to inhibit lid switch: %s", error->message);
        g_error_free (error);
    }
    else
    {
        g_variant_get (res, "(h)", &idx);
        manager->priv->inhibit_lid_switch_fd = g_unix_fd_list_get (fd_list, idx, &error);
        if (manager->priv->inhibit_lid_switch_fd == -1)
        {
            g_warning ("Failed to receive system inhibitor fd: %s", error->message);
            g_error_free (error);
        }
        g_debug ("System inhibitor fd is %d", manager->priv->inhibit_lid_switch_fd);
        g_object_unref (fd_list);
        g_variant_unref (res);
    }
}

static void
inhibit_lid_switch (GsdPowerManager *manager)
{
    GVariant *params;

    if (manager->priv->inhibit_lid_switch_taken)
    {
        g_debug ("already inhibited lid-switch");
        return;
    }
    g_debug ("Adding lid switch system inhibitor");
    manager->priv->inhibit_lid_switch_taken = TRUE;

    params = g_variant_new ("(ssss)",
                            "handle-lid-switch",
                            g_get_user_name (),
                            "Multiple displays attached",
                            "block");
    g_dbus_proxy_call_with_unix_fd_list (manager->priv->logind_proxy,
                                         "Inhibit",
                                         params,
                                         0,
                                         G_MAXINT,
                                         NULL,
                                         NULL,
                                         inhibit_lid_switch_done,
                                         manager);
}

static void
uninhibit_lid_switch (GsdPowerManager *manager)
{
    if (manager->priv->inhibit_lid_switch_fd == -1)
    {
        g_debug ("no lid-switch inhibitor");
        return;
    }
    g_debug ("Removing lid switch system inhibitor");
    close (manager->priv->inhibit_lid_switch_fd);
    manager->priv->inhibit_lid_switch_fd = -1;
    manager->priv->inhibit_lid_switch_taken = FALSE;
}

static void
inhibit_suspend_done (GObject      *source,
                      GAsyncResult *result,
                      gpointer      user_data)
{
    GDBusProxy *proxy = G_DBUS_PROXY (source);
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);
    GError *error = NULL;
    GVariant *res;
    GUnixFDList *fd_list = NULL;
    gint idx;

    res = g_dbus_proxy_call_with_unix_fd_list_finish (proxy, &fd_list, result, &error);
    if (res == NULL)
    {
        g_warning ("Unable to inhibit suspend: %s", error->message);
        g_error_free (error);
    }
    else
    {
        g_variant_get (res, "(h)", &idx);
        manager->priv->inhibit_suspend_fd = g_unix_fd_list_get (fd_list, idx, &error);
        if (manager->priv->inhibit_suspend_fd == -1)
        {
            g_warning ("Failed to receive system inhibitor fd: %s", error->message);
            g_error_free (error);
        }
        g_debug ("System inhibitor fd is %d", manager->priv->inhibit_suspend_fd);
        g_object_unref (fd_list);
        g_variant_unref (res);
    }
}

/* We take a delay inhibitor here, which causes logind to send a
 * PrepareForSleep signal, which gives us a chance to lock the screen
 * and do some other preparations.
 */
static void
inhibit_suspend (GsdPowerManager *manager)
{
    if (manager->priv->inhibit_suspend_taken)
    {
        g_debug ("already inhibited lid-switch");
        return;
    }
    g_debug ("Adding suspend delay inhibitor");
    manager->priv->inhibit_suspend_taken = TRUE;
    g_dbus_proxy_call_with_unix_fd_list (manager->priv->logind_proxy,
                                         "Inhibit",
                                         g_variant_new ("(ssss)",
                                                 "sleep",
                                                 g_get_user_name (),
                                                 "GNOME needs to lock the screen",
                                                 "delay"),
                                         0,
                                         G_MAXINT,
                                         NULL,
                                         NULL,
                                         inhibit_suspend_done,
                                         manager);
}

static void
uninhibit_suspend (GsdPowerManager *manager)
{
    if (manager->priv->inhibit_suspend_fd == -1)
    {
        g_debug ("no suspend delay inhibitor");
        return;
    }
    g_debug ("Removing suspend delay inhibitor");
    close (manager->priv->inhibit_suspend_fd);
    manager->priv->inhibit_suspend_fd = -1;
    manager->priv->inhibit_suspend_taken = FALSE;
}

static void
on_randr_event (GnomeRRScreen *screen, gpointer user_data)
{
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);

    if (suspend_on_lid_close (manager))
    {
        restart_inhibit_lid_switch_timer (manager);
        return;
    }

    /* when a second monitor is plugged in, we take the
     * handle-lid-switch inhibitor lock of logind to prevent
     * it from suspending.
     *
     * Uninhibiting is done in the inhibit_lid_switch_timer,
     * since we want to give users a few seconds when unplugging
     * and replugging an external monitor, not suspend right away.
     */
    inhibit_lid_switch (manager);
    setup_inhibit_lid_switch_timer (manager);
}

#ifdef GSD_MOCK
static gboolean
received_sigusr2 (GsdPowerManager *manager)
{
    on_randr_event (NULL, manager);
    return TRUE;
}
#endif /* GSD_MOCK */

static void
handle_suspend_actions (GsdPowerManager *manager)
{
    backlight_disable (manager);
    uninhibit_suspend (manager);
}

static void
handle_resume_actions (GsdPowerManager *manager)
{
    /* close existing notifications on resume, the system power
     * state is probably different now */
    notify_close_if_showing (&manager->priv->notification_low);
    notify_close_if_showing (&manager->priv->notification_ups_discharging);
    main_battery_or_ups_low_changed (manager, FALSE);

    /* ensure we turn the panel back on after resume */
    backlight_enable (manager);

    /* And work-around Xorg bug:
     * https://bugs.freedesktop.org/show_bug.cgi?id=59576 */
    reset_idletime ();

    /* set up the delay again */
    inhibit_suspend (manager);
}

static void
logind_proxy_signal_cb (GDBusProxy  *proxy,
                        const gchar *sender_name,
                        const gchar *signal_name,
                        GVariant    *parameters,
                        gpointer     user_data)
{
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);
    gboolean is_about_to_suspend;

    if (g_strcmp0 (signal_name, "PrepareForSleep") != 0)
        return;
    g_variant_get (parameters, "(b)", &is_about_to_suspend);
    if (is_about_to_suspend)
    {
        handle_suspend_actions (manager);
    }
    else
    {
        handle_resume_actions (manager);
    }
}

gboolean
gsd_power_manager_start (GsdPowerManager *manager,
                         GError **error)
{
    g_debug ("Starting power manager");
    gnome_settings_profile_start (NULL);

    /* coldplug the list of screens */
    manager->priv->rr_screen = gnome_rr_screen_new (gdk_screen_get_default (), error);
    if (manager->priv->rr_screen == NULL)
    {
        g_debug ("Couldn't detect any screens, disabling plugin");
        return FALSE;
    }

    /* Check for XTEST support */
    if (supports_xtest () == FALSE)
    {
        g_debug ("XTEST extension required, disabling plugin");
        return FALSE;
    }

    /* Set up the logind proxy */
    manager->priv->logind_proxy =
        g_dbus_proxy_new_for_bus_sync (G_BUS_TYPE_SYSTEM,
                                       0,
                                       NULL,
                                       SYSTEMD_DBUS_NAME,
                                       SYSTEMD_DBUS_PATH,
                                       SYSTEMD_DBUS_INTERFACE,
                                       NULL,
                                       error);
    if (manager->priv->logind_proxy == NULL)
    {
        g_debug ("No systemd (logind) support, disabling plugin");
        return FALSE;
    }
    g_signal_connect (manager->priv->logind_proxy, "g-signal",
                      G_CALLBACK (logind_proxy_signal_cb),
                      manager);
    /* Set up a delay inhibitor to be informed about suspend attempts */
    inhibit_suspend (manager);

    /* track the active session */
    manager->priv->session = gnome_settings_session_get_session_proxy ();
    g_debug("manager->priv->session: %p", manager->priv->session);
    g_signal_connect (manager->priv->session, "g-properties-changed",
                      G_CALLBACK (engine_session_properties_changed_cb),
                      manager);

    manager->priv->kbd_brightness_old = -1;
    manager->priv->kbd_brightness_pre_dim = -1;
    manager->priv->pre_dim_brightness = -1;
    /*manager->priv->settings = g_settings_new(GSD_POWER_SETTINGS_SCHEMA);*/

    manager->priv->settings_profile = g_settings_new(DEEPIN_POWER_PROFILE_SCHEMA);
    g_signal_connect(manager->priv->settings_profile, "changed",
                     G_CALLBACK(engine_profile_changed_cb),
                     manager);
    char *s = g_settings_get_string(manager->priv->settings_profile,
                                    "current-profile");
    manager->priv->settings_path = (char*)malloc(
                                       sizeof(DEEPIN_POWER_SETTINGS_PATH_PRE) + strlen(s) + 1);
    strcpy(manager->priv->settings_path, DEEPIN_POWER_SETTINGS_PATH_PRE);
    strcat(manager->priv->settings_path, s);
    strcat(manager->priv->settings_path, "/");

    g_debug("created new setttings with path :%s\n", manager->priv->settings_path);
    manager->priv->settings = g_settings_new_with_path(
                                  DEEPIN_POWER_SETTINGS_SCHEMA,
                                  manager->priv->settings_path);
    g_signal_connect (manager->priv->settings, "changed",
                      G_CALLBACK (engine_settings_key_changed_cb), manager);
    manager->priv->settings_screensaver = g_settings_new ("org.gnome.desktop.screensaver");
    manager->priv->settings_session = g_settings_new ("org.gnome.desktop.session");
    /*g_signal_connect (manager->priv->settings_session, "changed",*/
    /*G_CALLBACK (engine_settings_key_changed_cb), manager);*/
    /*manager->priv->settings_xrandr = g_settings_new (GSD_XRANDR_SETTINGS_SCHEMA);*/
    manager->priv->up_client = up_client_new ();
    manager->priv->lid_is_closed = up_client_get_lid_is_closed (manager->priv->up_client);
    g_signal_connect (manager->priv->up_client, "device-added",
                      G_CALLBACK (engine_device_added_cb), manager);
    g_signal_connect (manager->priv->up_client, "device-removed",
                      G_CALLBACK (engine_device_removed_cb), manager);
    g_signal_connect (manager->priv->up_client, "device-changed",
                      G_CALLBACK (engine_device_changed_cb), manager);
    g_signal_connect_after (manager->priv->up_client, "changed",
                            G_CALLBACK (up_client_changed_cb), manager);
    g_signal_connect (manager->priv->up_client, "notify::on-battery",
                      G_CALLBACK (up_client_on_battery_cb), manager);

    /* connect to UPower for keyboard backlight control */
    g_dbus_proxy_new_for_bus (G_BUS_TYPE_SYSTEM,
                              G_DBUS_PROXY_FLAGS_DO_NOT_LOAD_PROPERTIES,
                              NULL,
                              UPOWER_DBUS_NAME,
                              UPOWER_DBUS_PATH_KBDBACKLIGHT,
                              UPOWER_DBUS_INTERFACE_KBDBACKLIGHT,
                              NULL,
                              power_keyboard_proxy_ready_cb,
                              manager);

    /* connect to the session */
    g_dbus_proxy_new_for_bus (G_BUS_TYPE_SESSION,
                              0,
                              NULL,
                              GNOME_SESSION_DBUS_NAME,
                              GNOME_SESSION_DBUS_PATH_PRESENCE,
                              GNOME_SESSION_DBUS_INTERFACE_PRESENCE,
                              NULL,
                              session_presence_proxy_ready_cb,
                              manager);

    manager->priv->screensaver_watch_id =
        g_bus_watch_name (G_BUS_TYPE_SESSION,
                          GS_DBUS_NAME,
                          G_BUS_NAME_WATCHER_FLAGS_NONE,
                          (GBusNameAppearedCallback) screensaver_appeared_cb,
                          (GBusNameVanishedCallback) screensaver_vanished_cb,
                          manager,
                          NULL);

    manager->priv->devices_array = g_ptr_array_new_with_free_func (g_object_unref);

    /* create a fake virtual composite battery */
    manager->priv->device_composite = up_device_new ();
    g_object_set (manager->priv->device_composite,
                  "kind", UP_DEVICE_KIND_BATTERY,
                  "is-rechargeable", TRUE,
                  "native-path", "dummy:composite_battery",
                  "power-supply", TRUE,
                  "is-present", TRUE,
                  NULL);

    /* get percentage policy */
    manager->priv->low_percentage = g_settings_get_int (manager->priv->settings,
                                    "percentage-low");
    manager->priv->critical_percentage = g_settings_get_int (manager->priv->settings,
                                         "percentage-critical");
    manager->priv->action_percentage = g_settings_get_int (manager->priv->settings,
                                       "percentage-action");

    /* get time policy */
    manager->priv->low_time = g_settings_get_int (manager->priv->settings,
                              "time-low");
    manager->priv->critical_time = g_settings_get_int (manager->priv->settings,
                                   "time-critical");
    manager->priv->action_time = g_settings_get_int (manager->priv->settings,
                                 "time-action");

    /* we can disable this if the time remaining is inaccurate or just plain wrong */
    manager->priv->use_time_primary = g_settings_get_boolean (manager->priv->settings,
                                      "use-time-for-policy");

    /* create IDLETIME watcher */
    manager->priv->idle_monitor = gnome_idle_monitor_new ();

    /* set up the screens */
    g_signal_connect (manager->priv->rr_screen, "changed", G_CALLBACK (on_randr_event), manager);
    on_randr_event (manager->priv->rr_screen, manager);

#ifdef GSD_MOCK
    g_unix_signal_add (SIGUSR2, (GSourceFunc) received_sigusr2, manager);
#endif /* GSD_MOCK */

    /* check whether a backlight is available */
    manager->priv->backlight_available = backlight_available (manager->priv->rr_screen);

    /* ensure the default dpms timeouts are cleared */
    backlight_enable (manager);

    /* coldplug the engine */
    engine_coldplug (manager);
    idle_configure (manager);

    manager->priv->xscreensaver_watchdog_timer_id = gsd_power_enable_screensaver_watchdog ();

    /* don't blank inside a VM */
    manager->priv->is_virtual_machine = gsd_power_is_hardware_a_vm ();

    gnome_settings_profile_end (NULL);
    return TRUE;
}

void
gsd_power_manager_stop (GsdPowerManager *manager)
{
    g_debug ("Stopping power manager");

    if (manager->priv->inhibit_lid_switch_timer_id != 0)
    {
        g_source_remove (manager->priv->inhibit_lid_switch_timer_id);
        manager->priv->inhibit_lid_switch_timer_id = 0;
    }

    if (manager->priv->screensaver_cancellable != NULL)
    {
        g_cancellable_cancel (manager->priv->screensaver_cancellable);
        g_clear_object (&manager->priv->screensaver_cancellable);
    }

    if (manager->priv->screensaver_watch_id != 0)
    {
        g_bus_unwatch_name (manager->priv->screensaver_watch_id);
        manager->priv->screensaver_watch_id = 0;
    }

    if (manager->priv->bus_cancellable != NULL)
    {
        g_cancellable_cancel (manager->priv->bus_cancellable);
        g_object_unref (manager->priv->bus_cancellable);
        manager->priv->bus_cancellable = NULL;
    }

    if (manager->priv->introspection_data)
    {
        g_dbus_node_info_unref (manager->priv->introspection_data);
        manager->priv->introspection_data = NULL;
    }

    g_signal_handlers_disconnect_by_data (manager->priv->up_client, manager);

    g_clear_object (&manager->priv->session);
    g_clear_object (&manager->priv->settings);
    g_clear_object (&manager->priv->settings_screensaver);
    g_clear_object (&manager->priv->settings_session);
    g_clear_object (&manager->priv->up_client);

    if (manager->priv->inhibit_lid_switch_fd != -1)
    {
        close (manager->priv->inhibit_lid_switch_fd);
        manager->priv->inhibit_lid_switch_fd = -1;
        manager->priv->inhibit_lid_switch_taken = FALSE;
    }
    if (manager->priv->inhibit_suspend_fd != -1)
    {
        close (manager->priv->inhibit_suspend_fd);
        manager->priv->inhibit_suspend_fd = -1;
        manager->priv->inhibit_suspend_taken = FALSE;
    }

    g_clear_object (&manager->priv->logind_proxy);
    g_clear_object (&manager->priv->rr_screen);

    g_ptr_array_unref (manager->priv->devices_array);
    manager->priv->devices_array = NULL;
    g_clear_object (&manager->priv->device_composite);
    g_clear_object (&manager->priv->previous_icon);

    g_clear_pointer (&manager->priv->previous_summary, g_free);

    g_clear_object (&manager->priv->session_presence_proxy);
    g_clear_object (&manager->priv->screensaver_proxy);

    play_loop_stop (&manager->priv->critical_alert_timeout_id);

    g_clear_object (&manager->priv->idle_monitor);

    if (manager->priv->xscreensaver_watchdog_timer_id > 0)
    {
        g_source_remove (manager->priv->xscreensaver_watchdog_timer_id);
        manager->priv->xscreensaver_watchdog_timer_id = 0;
    }
}

static void
gsd_power_manager_init (GsdPowerManager *manager)
{
    manager->priv = GSD_POWER_MANAGER_GET_PRIVATE (manager);
    manager->priv->inhibit_lid_switch_fd = -1;
    manager->priv->inhibit_suspend_fd = -1;
    manager->priv->screensaver_cancellable = g_cancellable_new ();
    manager->priv->bus_cancellable = g_cancellable_new ();
}

/* returns new level */
static void
handle_method_call_keyboard (GsdPowerManager *manager,
                             const gchar *method_name,
                             GVariant *parameters,
                             GDBusMethodInvocation *invocation)
{
    gint step;
    gint value = -1;
    gboolean ret;
    guint percentage;
    GError *error = NULL;

    if (g_strcmp0 (method_name, "StepUp") == 0)
    {
        g_debug ("keyboard step up");
        step = BRIGHTNESS_STEP_AMOUNT (manager->priv->kbd_brightness_max);
        value = MIN (manager->priv->kbd_brightness_now + step,
                     manager->priv->kbd_brightness_max);
        ret = upower_kbd_set_brightness (manager, value, &error);

    }
    else if (g_strcmp0 (method_name, "StepDown") == 0)
    {
        g_debug ("keyboard step down");
        step = BRIGHTNESS_STEP_AMOUNT (manager->priv->kbd_brightness_max);
        value = MAX (manager->priv->kbd_brightness_now - step, 0);
        ret = upower_kbd_set_brightness (manager, value, &error);

    }
    else if (g_strcmp0 (method_name, "Toggle") == 0)
    {
        ret = upower_kbd_toggle (manager, &error);
    }
    else
    {
        g_assert_not_reached ();
    }

    /* return value */
    if (!ret)
    {
        g_dbus_method_invocation_take_error (invocation,
                                             error);
    }
    else
    {
        percentage = ABS_TO_PERCENTAGE (0,
                                        manager->priv->kbd_brightness_max,
                                        value);
        g_dbus_method_invocation_return_value (invocation,
                                               g_variant_new ("(u)",
                                                       percentage));
    }
}

static void
handle_method_call_screen (GsdPowerManager *manager,
                           const gchar *method_name,
                           GVariant *parameters,
                           GDBusMethodInvocation *invocation)
{
    gboolean ret = FALSE;
    gint value = -1;
    guint value_tmp;
    GError *error = NULL;

    if (!manager->priv->backlight_available)
    {
        g_set_error_literal (&error,
                             GSD_POWER_MANAGER_ERROR,
                             GSD_POWER_MANAGER_ERROR_FAILED,
                             "Screen backlight not available");
        goto out;
    }

    if (g_strcmp0 (method_name, "GetPercentage") == 0)
    {
        g_debug ("screen get percentage");
        value = backlight_get_percentage (manager->priv->rr_screen, &error);

    }
    else if (g_strcmp0 (method_name, "SetPercentage") == 0)
    {
        g_debug ("screen set percentage");
        g_variant_get (parameters, "(u)", &value_tmp);
        ret = backlight_set_percentage (manager->priv->rr_screen, value_tmp, &error);
        if (ret)
        {
            value = value_tmp;
            backlight_emit_changed (manager);
        }

    }
    else if (g_strcmp0 (method_name, "StepUp") == 0)
    {
        g_debug ("screen step up");
        value = backlight_step_up (manager->priv->rr_screen, &error);
        if (value != -1)
            backlight_emit_changed (manager);
    }
    else if (g_strcmp0 (method_name, "StepDown") == 0)
    {
        g_debug ("screen step down");
        value = backlight_step_down (manager->priv->rr_screen, &error);
        if (value != -1)
            backlight_emit_changed (manager);
    }
    else
    {
        g_assert_not_reached ();
    }

out:
    /* return value */
    if (value < 0)
    {
        g_dbus_method_invocation_take_error (invocation,
                                             error);
    }
    else
    {
        g_dbus_method_invocation_return_value (invocation,
                                               g_variant_new ("(u)",
                                                       value));
    }
}

static GVariant *
device_to_variant_blob (UpDevice *device)
{
    const gchar *object_path;
    gchar *device_icon;
    gdouble percentage;
    GIcon *icon;
    guint64 time_empty, time_full;
    guint64 time_state = 0;
    GVariant *value;
    UpDeviceKind kind;
    UpDeviceState state;

    icon = gpm_upower_get_device_icon (device, TRUE);
    device_icon = g_icon_to_string (icon);
    g_object_get (device,
                  "kind", &kind,
                  "percentage", &percentage,
                  "state", &state,
                  "time-to-empty", &time_empty,
                  "time-to-full", &time_full,
                  NULL);

    /* only return time for these simple states */
    if (state == UP_DEVICE_STATE_DISCHARGING)
        time_state = time_empty;
    else if (state == UP_DEVICE_STATE_CHARGING)
        time_state = time_full;

    /* get an object path, even for the composite device */
    object_path = up_device_get_object_path (device);
    if (object_path == NULL)
        object_path = GSD_DBUS_PATH;

    /* format complex object */
    value = g_variant_new ("(susdut)",
                           object_path,
                           kind,
                           device_icon,
                           percentage,
                           state,
                           time_state);
    g_free (device_icon);
    g_object_unref (icon);
    return value;
}

static void
handle_method_call_main (GsdPowerManager *manager,
                         const gchar *method_name,
                         GVariant *parameters,
                         GDBusMethodInvocation *invocation)
{
    GPtrArray *array;
    guint i;
    GVariantBuilder *builder;
    GVariant *tuple = NULL;
    GVariant *value = NULL;
    UpDevice *device;

    /* return object */
    if (g_strcmp0 (method_name, "GetPrimaryDevice") == 0)
    {

        /* get the virtual device */
        device = engine_get_primary_device (manager);
        if (device == NULL)
        {
            g_dbus_method_invocation_return_dbus_error (invocation,
                    "org.gnome.SettingsDaemon.Power.Failed",
                    "There is no primary device.");
            return;
        }

        /* return the value */
        value = device_to_variant_blob (device);
        tuple = g_variant_new_tuple (&value, 1);
        g_dbus_method_invocation_return_value (invocation, tuple);
        g_object_unref (device);
        return;
    }

    /* return array */
    if (g_strcmp0 (method_name, "GetDevices") == 0)
    {

        /* create builder */
        builder = g_variant_builder_new (G_VARIANT_TYPE("a(susdut)"));

        /* add each tuple to the array */
        array = manager->priv->devices_array;
        for (i = 0; i < array->len; i++)
        {
            device = g_ptr_array_index (array, i);
            value = device_to_variant_blob (device);
            g_variant_builder_add_value (builder, value);
        }

        /* return the value */
        value = g_variant_builder_end (builder);
        tuple = g_variant_new_tuple (&value, 1);
        g_dbus_method_invocation_return_value (invocation, tuple);
        g_variant_builder_unref (builder);
        return;
    }

    g_assert_not_reached ();
}

static void
handle_method_call (GDBusConnection       *connection,
                    const gchar           *sender,
                    const gchar           *object_path,
                    const gchar           *interface_name,
                    const gchar           *method_name,
                    GVariant              *parameters,
                    GDBusMethodInvocation *invocation,
                    gpointer               user_data)
{
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);

    /* Check session pointer as a proxy for whether the manager is in the
       start or stop state */
    if (manager->priv->session == NULL)
    {
        return;
    }

    g_debug ("Calling method '%s.%s' for Power",
             interface_name, method_name);

    if (g_strcmp0 (interface_name, GSD_POWER_DBUS_INTERFACE) == 0)
    {
        handle_method_call_main (manager,
                                 method_name,
                                 parameters,
                                 invocation);
    }
    else if (g_strcmp0 (interface_name, GSD_POWER_DBUS_INTERFACE_SCREEN) == 0)
    {
        handle_method_call_screen (manager,
                                   method_name,
                                   parameters,
                                   invocation);
    }
    else if (g_strcmp0 (interface_name, GSD_POWER_DBUS_INTERFACE_KEYBOARD) == 0)
    {
        handle_method_call_keyboard (manager,
                                     method_name,
                                     parameters,
                                     invocation);
    }
    else
    {
        g_warning ("not recognised interface: %s", interface_name);
    }
}

static GVariant *
handle_get_property (GDBusConnection *connection,
                     const gchar *sender,
                     const gchar *object_path,
                     const gchar *interface_name,
                     const gchar *property_name,
                     GError **error, gpointer user_data)
{
    GsdPowerManager *manager = GSD_POWER_MANAGER (user_data);
    GVariant *retval = NULL;

    /* Check session pointer as a proxy for whether the manager is in the
       start or stop state */
    if (manager->priv->session == NULL)
    {
        return NULL;
    }

    if (g_strcmp0 (property_name, "Icon") == 0)
    {
        retval = engine_get_icon_property_variant (manager);
    }
    else if (g_strcmp0 (property_name, "Tooltip") == 0)
    {
        retval = engine_get_tooltip_property_variant (manager);
    }
    else if (g_strcmp0 (property_name, "Percentage") == 0)
    {
        gdouble percentage;
        percentage = engine_get_percentage (manager);
        if (percentage >= 0)
            retval = g_variant_new_double (percentage);
    }

    return retval;
}

static const GDBusInterfaceVTable interface_vtable =
{
    handle_method_call,
    handle_get_property,
    NULL, /* SetProperty */
};

static void
on_bus_gotten (GObject             *source_object,
               GAsyncResult        *res,
               GsdPowerManager     *manager)
{
    GDBusConnection *connection;
    GDBusInterfaceInfo **infos;
    GError *error = NULL;
    guint i;

    connection = g_bus_get_finish (res, &error);
    if (connection == NULL)
    {
        if (!g_error_matches (error, G_IO_ERROR, G_IO_ERROR_CANCELLED))
            g_warning ("Could not get session bus: %s", error->message);
        g_error_free (error);
        return;
    }

    manager->priv->connection = connection;

    infos = manager->priv->introspection_data->interfaces;
    for (i = 0; infos[i] != NULL; i++)
    {
        g_dbus_connection_register_object (connection,
                                           GSD_POWER_DBUS_PATH,
                                           infos[i],
                                           &interface_vtable,
                                           manager,
                                           NULL,
                                           NULL);
    }

    manager->priv->name_id = g_bus_own_name_on_connection (connection,
                             GSD_POWER_DBUS_NAME,
                             G_BUS_NAME_OWNER_FLAGS_NONE,
                             NULL,
                             NULL,
                             NULL,
                             NULL);
}

static void
register_manager_dbus (GsdPowerManager *manager)
{
    manager->priv->introspection_data = g_dbus_node_info_new_for_xml (introspection_xml, NULL);
    g_assert (manager->priv->introspection_data != NULL);

    g_bus_get (G_BUS_TYPE_SESSION,
               manager->priv->bus_cancellable,
               (GAsyncReadyCallback) on_bus_gotten,
               manager);
}

GsdPowerManager *
gsd_power_manager_new (void)
{
    if (manager_object != NULL)
    {
        g_object_ref (manager_object);
    }
    else
    {
        manager_object = g_object_new (GSD_TYPE_POWER_MANAGER, NULL);
        g_object_add_weak_pointer (manager_object,
                                   (gpointer *) &manager_object);
        register_manager_dbus (manager_object);
    }
    return GSD_POWER_MANAGER (manager_object);
}



