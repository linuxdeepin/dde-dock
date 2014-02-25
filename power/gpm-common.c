/* -*- Mode: C; tab-width: 8; indent-tabs-mode: nil; c-basic-offset: 8 -*-
 *
 * Copyright (C) 2005-2011 Richard Hughes <richard@hughsie.com>
 *
 * Licensed under the GNU General Public License Version 2
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
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

#include "config.h"

#include <stdlib.h>
#include <stdio.h>
#include <sys/wait.h>
#include <math.h>
#include <glib.h>
#include <glib/gi18n.h>
#include <gdk/gdkx.h>
#include <X11/extensions/XTest.h>
#include <X11/extensions/dpms.h>
#include <canberra-gtk.h>

#define GNOME_DESKTOP_USE_UNSTABLE_API
#include <libgnome-desktop/gnome-rr.h>

#include "gpm-common.h"
#include "gsd-power-constants.h"
#include "gsd-power-manager.h"
#include "gsd-backlight-linux.h"

#define XSCREENSAVER_WATCHDOG_TIMEOUT           120 /* seconds */
#define UPS_SOUND_LOOP_ID                        99
#define GSD_POWER_MANAGER_CRITICAL_ALERT_TIMEOUT  5 /* seconds */

/* take a discrete value with offset and convert to percentage */
int
gsd_power_backlight_abs_to_percentage (int min, int max, int value)
{
        g_return_val_if_fail (max > min, -1);
        g_return_val_if_fail (value >= min, -1);
        g_return_val_if_fail (value <= max, -1);
        return (((value - min) * 100) / (max - min));
}

#define GPM_UP_TIME_PRECISION                   5*60
#define GPM_UP_TEXT_MIN_TIME                    120

/**
 * Return value: The time string, e.g. "2 hours 3 minutes"
 **/
gchar *
gpm_get_timestring (guint time_secs)
{
        char* timestring = NULL;
        gint  hours;
        gint  minutes;

        /* Add 0.5 to do rounding */
        minutes = (int) ( ( time_secs / 60.0 ) + 0.5 );

        if (minutes == 0) {
                timestring = g_strdup (_("Unknown time"));
                return timestring;
        }

        if (minutes < 60) {
                timestring = g_strdup_printf (ngettext ("%i minute",
                                                        "%i minutes",
                                                        minutes), minutes);
                return timestring;
        }

        hours = minutes / 60;
        minutes = minutes % 60;
        if (minutes == 0)
                timestring = g_strdup_printf (ngettext (
                                "%i hour",
                                "%i hours",
                                hours), hours);
        else
                /* TRANSLATOR: "%i %s %i %s" are "%i hours %i minutes"
                 * Swap order with "%2$s %2$i %1$s %1$i if needed */
                timestring = g_strdup_printf (_("%i %s %i %s"),
                                hours, ngettext ("hour", "hours", hours),
                                minutes, ngettext ("minute", "minutes", minutes));
        return timestring;
}

static const gchar *
gpm_upower_get_device_icon_index (UpDevice *device)
{
        gdouble percentage;
        /* get device properties */
        g_object_get (device, "percentage", &percentage, NULL);
        if (percentage < 10)
                return "000";
        else if (percentage < 30)
                return "020";
        else if (percentage < 50)
                return "040";
        else if (percentage < 70)
                return "060";
        else if (percentage < 90)
                return "080";
        return "100";
}

static const gchar *
gpm_upower_get_device_icon_suffix (UpDevice *device)
{
        gdouble percentage;
        /* get device properties */
        g_object_get (device, "percentage", &percentage, NULL);
        if (percentage < 10)
                return "caution";
        else if (percentage < 30)
                return "low";
        else if (percentage < 60)
                return "good";
        return "full";
}

GIcon *
gpm_upower_get_device_icon (UpDevice *device, gboolean use_symbolic)
{
        GString *filename;
        gchar **iconnames;
        const gchar *kind_str;
        const gchar *suffix_str;
        const gchar *index_str;
        UpDeviceKind kind;
        UpDeviceState state;
        gboolean is_present;
        gdouble percentage;
        GIcon *icon = NULL;

        g_return_val_if_fail (device != NULL, NULL);

        /* get device properties */
        g_object_get (device,
                      "kind", &kind,
                      "state", &state,
                      "percentage", &percentage,
                      "is-present", &is_present,
                      NULL);

        /* get correct icon prefix */
        filename = g_string_new (NULL);

        /* get the icon from some simple rules */
        if (kind == UP_DEVICE_KIND_LINE_POWER) {
                if (use_symbolic)
                        g_string_append (filename, "ac-adapter-symbolic;");
                g_string_append (filename, "ac-adapter;");

        } else if (kind == UP_DEVICE_KIND_MONITOR) {
                if (use_symbolic)
                        g_string_append (filename, "gpm-monitor-symbolic;");
                g_string_append (filename, "gpm-monitor;");

        } else {

                kind_str = up_device_kind_to_string (kind);
                if (!is_present) {
                        if (use_symbolic)
                                g_string_append (filename, "battery-missing-symbolic;");
                        g_string_append_printf (filename, "gpm-%s-missing;", kind_str);
                        g_string_append_printf (filename, "gpm-%s-000;", kind_str);
                        g_string_append (filename, "battery-missing;");

                } else {
                        switch (state) {
                        case UP_DEVICE_STATE_EMPTY:
                                if (use_symbolic)
                                        g_string_append (filename, "battery-empty-symbolic;");
                                g_string_append_printf (filename, "gpm-%s-empty;", kind_str);
                                g_string_append_printf (filename, "gpm-%s-000;", kind_str);
                                g_string_append (filename, "battery-empty;");
                                break;
                        case UP_DEVICE_STATE_FULLY_CHARGED:
                                if (use_symbolic) {
                                        g_string_append (filename, "battery-full-charged-symbolic;");
                                        g_string_append (filename, "battery-full-charging-symbolic;");
                                }
                                g_string_append_printf (filename, "gpm-%s-full;", kind_str);
                                g_string_append_printf (filename, "gpm-%s-100;", kind_str);
                                g_string_append (filename, "battery-full-charged;");
                                g_string_append (filename, "battery-full-charging;");
                                break;
                        case UP_DEVICE_STATE_CHARGING:
                        case UP_DEVICE_STATE_PENDING_CHARGE:
                                suffix_str = gpm_upower_get_device_icon_suffix (device);
                                index_str = gpm_upower_get_device_icon_index (device);
                                if (use_symbolic)
                                        g_string_append_printf (filename, "battery-%s-charging-symbolic;", suffix_str);
                                g_string_append_printf (filename, "gpm-%s-%s-charging;", kind_str, index_str);
                                g_string_append_printf (filename, "battery-%s-charging;", suffix_str);
                                break;
                        case UP_DEVICE_STATE_DISCHARGING:
                        case UP_DEVICE_STATE_PENDING_DISCHARGE:
                                suffix_str = gpm_upower_get_device_icon_suffix (device);
                                index_str = gpm_upower_get_device_icon_index (device);
                                if (use_symbolic)
                                        g_string_append_printf (filename, "battery-%s-symbolic;", suffix_str);
                                g_string_append_printf (filename, "gpm-%s-%s;", kind_str, index_str);
                                g_string_append_printf (filename, "battery-%s;", suffix_str);
                                break;
                        default:
                                if (use_symbolic)
                                        g_string_append (filename, "battery-missing-symbolic;");
                                g_string_append (filename, "gpm-battery-missing;");
                                g_string_append (filename, "battery-missing;");
                        }
                }
        }

        /* nothing matched */
        if (filename->len == 0) {
                g_warning ("nothing matched, falling back to default icon");
                g_string_append (filename, "dialog-warning;");
        }

        g_debug ("got filename: %s", filename->str);

        iconnames = g_strsplit (filename->str, ";", -1);
        icon = g_themed_icon_new_from_names (iconnames, -1);

        g_strfreev (iconnames);
        g_string_free (filename, TRUE);
        return icon;
}

/**
 * gpm_precision_round_down:
 * @value: The input value
 * @smallest: The smallest increment allowed
 *
 * 101, 10      100
 * 95,  10      90
 * 0,   10      0
 * 112, 10      110
 * 100, 10      100
 **/
static gint
gpm_precision_round_down (gfloat value, gint smallest)
{
        gfloat division;
        if (fabs (value) < 0.01)
                return 0;
        if (smallest == 0) {
                g_warning ("divisor zero");
                return 0;
        }
        division = (gfloat) value / (gfloat) smallest;
        division = floorf (division);
        division *= smallest;
        return (gint) division;
}

gchar *
gpm_upower_get_device_summary (UpDevice *device)
{
        const gchar *kind_desc = NULL;
        const gchar *device_desc = NULL;
        GString *description;
        guint time_to_full_round;
        guint time_to_empty_round;
        gchar *time_to_full_str = NULL;
        gchar *time_to_empty_str = NULL;
        UpDeviceKind kind;
        UpDeviceState state;
        gdouble percentage;
        gboolean is_present;
        gint64 time_to_full;
        gint64 time_to_empty;

        /* get device properties */
        g_object_get (device,
                      "kind", &kind,
                      "state", &state,
                      "percentage", &percentage,
                      "is-present", &is_present,
                      "time-to-full", &time_to_full,
                      "time-to-empty", &time_to_empty,
                      NULL);

        description = g_string_new (NULL);
        kind_desc = gpm_device_kind_to_localised_string (kind, 1);
        device_desc = gpm_device_to_localised_string (device);

        /* not installed */
        if (!is_present) {
                g_string_append (description, device_desc);
                goto out;
        }

        /* don't display all the extra stuff for keyboards and mice */
        if (kind == UP_DEVICE_KIND_MOUSE ||
            kind == UP_DEVICE_KIND_KEYBOARD ||
            kind == UP_DEVICE_KIND_PDA) {
                g_string_append (description, kind_desc);
                g_string_append_printf (description, " (%.0f%%)", percentage);
                goto out;
        }

        /* we care if we are on AC */
        if (kind == UP_DEVICE_KIND_PHONE) {
                if (state == UP_DEVICE_STATE_CHARGING || !(state == UP_DEVICE_STATE_DISCHARGING)) {
                        g_string_append (description, device_desc);
                        g_string_append_printf (description, " (%.0f%%)", percentage);
                        goto out;
                }
                g_string_append (description, kind_desc);
                g_string_append_printf (description, " (%.0f%%)", percentage);
                goto out;
        }

        /* precalculate so we don't get Unknown time remaining */
        time_to_full_round = gpm_precision_round_down (time_to_full, GPM_UP_TIME_PRECISION);
        time_to_empty_round = gpm_precision_round_down (time_to_empty, GPM_UP_TIME_PRECISION);

        /* we always display "Laptop battery 16 minutes remaining" as we need to clarify what device we are refering to */
        if (state == UP_DEVICE_STATE_FULLY_CHARGED) {

                g_string_append (description, device_desc);

                if (kind == UP_DEVICE_KIND_BATTERY && time_to_empty_round > GPM_UP_TEXT_MIN_TIME) {
                        time_to_empty_str = gpm_get_timestring (time_to_empty_round);
                        g_string_append (description, " - ");
                        /* TRANSLATORS: The laptop battery is charged, and we know a time.
                         * The parameter is the time, e.g. 7 hours 6 minutes */
                        g_string_append_printf (description, _("provides %s laptop runtime"), time_to_empty_str);
                }
                goto out;
        }
        if (state == UP_DEVICE_STATE_DISCHARGING) {

                if (time_to_empty_round > GPM_UP_TEXT_MIN_TIME) {
                        time_to_empty_str = gpm_get_timestring (time_to_empty_round);
                        /* TRANSLATORS: the device is discharging, and we have a time remaining
                         * The first parameter is the device type, e.g. "Laptop battery" and
                         * the second is the time, e.g. 7 hours 6 minutes */
                        g_string_append_printf (description, _("%s %s remaining"),
                                                kind_desc, time_to_empty_str);
                        g_string_append_printf (description, " (%.0f%%)", percentage);
                } else {
                        g_string_append (description, device_desc);
                        g_string_append_printf (description, " (%.0f%%)", percentage);
                }
                goto out;
        }
        if (state == UP_DEVICE_STATE_CHARGING) {

                if (time_to_full_round > GPM_UP_TEXT_MIN_TIME &&
                    time_to_empty_round > GPM_UP_TEXT_MIN_TIME) {

                        /* display both discharge and charge time */
                        time_to_full_str = gpm_get_timestring (time_to_full_round);
                        time_to_empty_str = gpm_get_timestring (time_to_empty_round);

                        /* TRANSLATORS: device is charging, and we have a time to full and a percentage
                         * The first parameter is the device type, e.g. "Laptop battery" and
                         * the second is the time, e.g. "7 hours 6 minutes" */
                        g_string_append_printf (description, _("%s %s until charged"),
                                                kind_desc, time_to_full_str);
                        g_string_append_printf (description, " (%.0f%%)", percentage);

                        g_string_append (description, " - ");
                        /* TRANSLATORS: the device is charging, and we have a time to full and empty.
                         * The parameter is a time string, e.g. "7 hours 6 minutes" */
                        g_string_append_printf (description, _("provides %s battery runtime"),
                                                time_to_empty_str);
                } else if (time_to_full_round > GPM_UP_TEXT_MIN_TIME) {

                        /* display only charge time */
                        time_to_full_str = gpm_get_timestring (time_to_full_round);

                        /* TRANSLATORS: device is charging, and we have a time to full and a percentage.
                         * The first parameter is the device type, e.g. "Laptop battery" and
                         * the second is the time, e.g. "7 hours 6 minutes" */
                        g_string_append_printf (description, _("%s %s until charged"),
                                                kind_desc, time_to_full_str);
                        g_string_append_printf (description, " (%.0f%%)", percentage);
                } else {
                        g_string_append (description, device_desc);
                        g_string_append_printf (description, " (%.0f%%)", percentage);
                }
                goto out;
        }
        if (state == UP_DEVICE_STATE_PENDING_DISCHARGE) {
                g_string_append (description, device_desc);
                g_string_append_printf (description, " (%.0f%%)", percentage);
                goto out;
        }
        if (state == UP_DEVICE_STATE_PENDING_CHARGE) {
                g_string_append (description, device_desc);
                g_string_append_printf (description, " (%.0f%%)", percentage);
                goto out;
        }
        if (state == UP_DEVICE_STATE_EMPTY) {
                g_string_append (description, device_desc);
                goto out;
        }

        /* fallback */
        g_warning ("in an undefined state we are not charging or "
                     "discharging and the batteries are also not charged");
        g_string_append (description, device_desc);
        g_string_append_printf (description, " (%.0f%%)", percentage);
out:
        g_free (time_to_full_str);
        g_free (time_to_empty_str);
        return g_string_free (description, FALSE);
}

gchar *
gpm_upower_get_device_description (UpDevice *device)
{
        GString *details;
        const gchar *text;
        gchar *time_str;
        UpDeviceKind kind;
        UpDeviceState state;
        UpDeviceTechnology technology;
        gdouble percentage;
        gdouble capacity;
        gdouble energy;
        gdouble energy_full;
        gdouble energy_full_design;
        gdouble energy_rate;
        gboolean is_present;
        gint64 time_to_full;
        gint64 time_to_empty;
        gchar *vendor = NULL;
        gchar *serial = NULL;
        gchar *model = NULL;

        g_return_val_if_fail (device != NULL, NULL);

        /* get device properties */
        g_object_get (device,
                      "kind", &kind,
                      "state", &state,
                      "percentage", &percentage,
                      "is-present", &is_present,
                      "time-to-full", &time_to_full,
                      "time-to-empty", &time_to_empty,
                      "technology", &technology,
                      "capacity", &capacity,
                      "energy", &energy,
                      "energy-full", &energy_full,
                      "energy-full-design", &energy_full_design,
                      "energy-rate", &energy_rate,
                      "vendor", &vendor,
                      "serial", &serial,
                      "model", &model,
                      NULL);

        details = g_string_new ("");
        text = gpm_device_kind_to_localised_string (kind, 1);
        /* TRANSLATORS: the type of data, e.g. Laptop battery */
        g_string_append_printf (details, "<b>%s</b> %s\n", _("Product:"), text);

        if (!is_present) {
                /* TRANSLATORS: device is missing */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Status:"), _("Missing"));
        } else if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                /* TRANSLATORS: device is charged */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Status:"), _("Charged"));
        } else if (state == UP_DEVICE_STATE_CHARGING) {
                /* TRANSLATORS: device is charging */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Status:"), _("Charging"));
        } else if (state == UP_DEVICE_STATE_DISCHARGING) {
                /* TRANSLATORS: device is discharging */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Status:"), _("Discharging"));
        }

        if (percentage >= 0) {
                /* TRANSLATORS: percentage */
                g_string_append_printf (details, "<b>%s</b> %.1f%%\n", _("Percentage charge:"), percentage);
        }
        if (vendor) {
                /* TRANSLATORS: manufacturer */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Vendor:"), vendor);
        }
        if (technology != UP_DEVICE_TECHNOLOGY_UNKNOWN) {
                text = gpm_device_technology_to_localised_string (technology);
                /* TRANSLATORS: how the battery is made, e.g. Lithium Ion */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Technology:"), text);
        }
        if (serial) {
                /* TRANSLATORS: serial number of the battery */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Serial number:"), serial);
        }
        if (model) {
                /* TRANSLATORS: model number of the battery */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Model:"), model);
        }
        if (time_to_full > 0) {
                time_str = gpm_get_timestring (time_to_full);
                /* TRANSLATORS: time to fully charged */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Charge time:"), time_str);
                g_free (time_str);
        }
        if (time_to_empty > 0) {
                time_str = gpm_get_timestring (time_to_empty);
                /* TRANSLATORS: time to empty */
                g_string_append_printf (details, "<b>%s</b> %s\n", _("Discharge time:"), time_str);
                g_free (time_str);
        }
        if (capacity > 0) {
                const gchar *condition;
                if (capacity > 99) {
                        /* TRANSLATORS: Excellent, Good, Fair and Poor are all related to battery Capacity */
                        condition = _("Excellent");
                } else if (capacity > 90) {
                        condition = _("Good");
                } else if (capacity > 70) {
                        condition = _("Fair");
                } else {
                        condition = _("Poor");
                }
                /* TRANSLATORS: %.1f is a percentage and %s the condition (Excellent, Good, ...) */
                g_string_append_printf (details, "<b>%s</b> %.1f%% (%s)\n",
                                        _("Capacity:"), capacity, condition);
        }
        if (kind == UP_DEVICE_KIND_BATTERY) {
                if (energy > 0) {
                        /* TRANSLATORS: current charge */
                        g_string_append_printf (details, "<b>%s</b> %.1f Wh\n",
                                                _("Current charge:"), energy);
                }
                if (energy_full > 0 &&
                    energy_full_design != energy_full) {
                        /* TRANSLATORS: last full is the charge the battery was seen to charge to */
                        g_string_append_printf (details, "<b>%s</b> %.1f Wh\n",
                                                _("Last full charge:"), energy_full);
                }
                if (energy_full_design > 0) {
                        /* Translators:  */
                        /* TRANSLATORS: Design charge is the amount of charge the battery is designed to have when brand new */
                        g_string_append_printf (details, "<b>%s</b> %.1f Wh\n",
                                                _("Design charge:"), energy_full_design);
                }
                if (energy_rate > 0) {
                        /* TRANSLATORS: the charge or discharge rate */
                        g_string_append_printf (details, "<b>%s</b> %.1f W\n",
                                                _("Charge rate:"), energy_rate);
                }
        }
        if (kind == UP_DEVICE_KIND_MOUSE ||
            kind == UP_DEVICE_KIND_KEYBOARD) {
                if (energy > 0) {
                        /* TRANSLATORS: the current charge for CSR devices */
                        g_string_append_printf (details, "<b>%s</b> %.0f/7\n",
                                                _("Current charge:"), energy);
                }
                if (energy_full_design > 0) {
                        /* TRANSLATORS: the design charge for CSR devices */
                        g_string_append_printf (details, "<b>%s</b> %.0f/7\n",
                                                _("Design charge:"), energy_full_design);
                }
        }
        /* remove the last \n */
        g_string_truncate (details, details->len-1);

        g_free (vendor);
        g_free (serial);
        g_free (model);
        return g_string_free (details, FALSE);
}

const gchar *
gpm_device_kind_to_localised_string (UpDeviceKind kind, guint number)
{
        const gchar *text = NULL;
        switch (kind) {
        case UP_DEVICE_KIND_LINE_POWER:
                /* TRANSLATORS: system power cord */
                text = ngettext ("AC adapter", "AC adapters", number);
                break;
        case UP_DEVICE_KIND_BATTERY:
                /* TRANSLATORS: laptop primary battery */
                text = ngettext ("Laptop battery", "Laptop batteries", number);
                break;
        case UP_DEVICE_KIND_UPS:
                /* TRANSLATORS: battery-backed AC power source */
                text = ngettext ("UPS", "UPSs", number);
                break;
        case UP_DEVICE_KIND_MONITOR:
                /* TRANSLATORS: a monitor is a device to measure voltage and current */
                text = ngettext ("Monitor", "Monitors", number);
                break;
        case UP_DEVICE_KIND_MOUSE:
                /* TRANSLATORS: wireless mice with internal batteries */
                text = ngettext ("Mouse", "Mice", number);
                break;
        case UP_DEVICE_KIND_KEYBOARD:
                /* TRANSLATORS: wireless keyboard with internal battery */
                text = ngettext ("Keyboard", "Keyboards", number);
                break;
        case UP_DEVICE_KIND_PDA:
                /* TRANSLATORS: portable device */
                text = ngettext ("PDA", "PDAs", number);
                break;
        case UP_DEVICE_KIND_PHONE:
                /* TRANSLATORS: cell phone (mobile...) */
                text = ngettext ("Cell phone", "Cell phones", number);
                break;
#if UP_CHECK_VERSION(0,9,5)
        case UP_DEVICE_KIND_MEDIA_PLAYER:
                /* TRANSLATORS: media player, mp3 etc */
                text = ngettext ("Media player", "Media players", number);
                break;
        case UP_DEVICE_KIND_TABLET:
                /* TRANSLATORS: tablet device */
                text = ngettext ("Tablet", "Tablets", number);
                break;
        case UP_DEVICE_KIND_COMPUTER:
                /* TRANSLATORS: tablet device */
                text = ngettext ("Computer", "Computers", number);
                break;
#endif
        default:
                g_warning ("enum unrecognised: %i", kind);
                text = up_device_kind_to_string (kind);
        }
        return text;
}

const gchar *
gpm_device_kind_to_icon (UpDeviceKind kind)
{
        const gchar *icon = NULL;
        switch (kind) {
        case UP_DEVICE_KIND_LINE_POWER:
                icon = "ac-adapter";
                break;
        case UP_DEVICE_KIND_BATTERY:
                icon = "battery";
                break;
        case UP_DEVICE_KIND_UPS:
                icon = "network-wired";
                break;
        case UP_DEVICE_KIND_MONITOR:
                icon = "application-certificate";
                break;
        case UP_DEVICE_KIND_MOUSE:
                icon = "input-mouse";
                break;
        case UP_DEVICE_KIND_KEYBOARD:
                icon = "input-keyboard";
                break;
        case UP_DEVICE_KIND_PDA:
                icon = "pda";
                break;
        case UP_DEVICE_KIND_PHONE:
                icon = "phone";
                break;
#if UP_CHECK_VERSION(0,9,5)
        case UP_DEVICE_KIND_MEDIA_PLAYER:
                icon = "multimedia-player";
                break;
        case UP_DEVICE_KIND_TABLET:
                icon = "input-tablet";
                break;
        case UP_DEVICE_KIND_COMPUTER:
                icon = "computer-apple-ipad";
                break;
#endif
        default:
                g_warning ("enum unrecognised: %i", kind);
                icon = "gtk-help";
        }
        return icon;
}

const gchar *
gpm_device_technology_to_localised_string (UpDeviceTechnology technology_enum)
{
        const gchar *technology = NULL;
        switch (technology_enum) {
        case UP_DEVICE_TECHNOLOGY_LITHIUM_ION:
                /* TRANSLATORS: battery technology */
                technology = _("Lithium Ion");
                break;
        case UP_DEVICE_TECHNOLOGY_LITHIUM_POLYMER:
                /* TRANSLATORS: battery technology */
                technology = _("Lithium Polymer");
                break;
        case UP_DEVICE_TECHNOLOGY_LITHIUM_IRON_PHOSPHATE:
                /* TRANSLATORS: battery technology */
                technology = _("Lithium Iron Phosphate");
                break;
        case UP_DEVICE_TECHNOLOGY_LEAD_ACID:
                /* TRANSLATORS: battery technology */
                technology = _("Lead acid");
                break;
        case UP_DEVICE_TECHNOLOGY_NICKEL_CADMIUM:
                /* TRANSLATORS: battery technology */
                technology = _("Nickel Cadmium");
                break;
        case UP_DEVICE_TECHNOLOGY_NICKEL_METAL_HYDRIDE:
                /* TRANSLATORS: battery technology */
                technology = _("Nickel metal hydride");
                break;
        case UP_DEVICE_TECHNOLOGY_UNKNOWN:
                /* TRANSLATORS: battery technology */
                technology = _("Unknown technology");
                break;
        default:
                g_assert_not_reached ();
                break;
        }
        return technology;
}

const gchar *
gpm_device_state_to_localised_string (UpDeviceState state)
{
        const gchar *state_string = NULL;

        switch (state) {
        case UP_DEVICE_STATE_CHARGING:
                /* TRANSLATORS: battery state */
                state_string = _("Charging");
                break;
        case UP_DEVICE_STATE_DISCHARGING:
                /* TRANSLATORS: battery state */
                state_string = _("Discharging");
                break;
        case UP_DEVICE_STATE_EMPTY:
                /* TRANSLATORS: battery state */
                state_string = _("Empty");
                break;
        case UP_DEVICE_STATE_FULLY_CHARGED:
                /* TRANSLATORS: battery state */
                state_string = _("Charged");
                break;
        case UP_DEVICE_STATE_PENDING_CHARGE:
                /* TRANSLATORS: battery state */
                state_string = _("Waiting to charge");
                break;
        case UP_DEVICE_STATE_PENDING_DISCHARGE:
                /* TRANSLATORS: battery state */
                state_string = _("Waiting to discharge");
                break;
        default:
                g_assert_not_reached ();
                break;
        }
        return state_string;
}

const gchar *
gpm_device_to_localised_string (UpDevice *device)
{
        UpDeviceState state;
        UpDeviceKind kind;
        gboolean present;

        /* get device parameters */
        g_object_get (device,
                      "is-present", &present,
                      "kind", &kind,
                      "state", &state,
                      NULL);

        /* laptop battery */
        if (kind == UP_DEVICE_KIND_BATTERY) {

                if (!present) {
                        /* TRANSLATORS: device not present */
                        return _("Laptop battery not present");
                }
                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Laptop battery is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Laptop battery is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("Laptop battery is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("Laptop battery is charged");
                }
                if (state == UP_DEVICE_STATE_PENDING_CHARGE) {
                        /* TRANSLATORS: battery state */
                        return _("Laptop battery is waiting to charge");
                }
                if (state == UP_DEVICE_STATE_PENDING_DISCHARGE) {
                        /* TRANSLATORS: battery state */
                        return _("Laptop battery is waiting to discharge");
                }
        }

        /* UPS */
        if (kind == UP_DEVICE_KIND_UPS) {

                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("UPS is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("UPS is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("UPS is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("UPS is charged");
                }
        }

        /* mouse */
        if (kind == UP_DEVICE_KIND_MOUSE) {

                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Mouse is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Mouse is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("Mouse is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("Mouse is charged");
                }
        }

        /* keyboard */
        if (kind == UP_DEVICE_KIND_KEYBOARD) {

                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Keyboard is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Keyboard is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("Keyboard is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("Keyboard is charged");
                }
        }

        /* PDA */
        if (kind == UP_DEVICE_KIND_PDA) {

                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("PDA is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("PDA is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("PDA is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("PDA is charged");
                }
        }

        /* phone */
        if (kind == UP_DEVICE_KIND_PHONE) {

                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Cell phone is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Cell phone is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("Cell phone is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("Cell phone is charged");
                }
        }
#if UP_CHECK_VERSION(0,9,5)

        /* media player */
        if (kind == UP_DEVICE_KIND_MEDIA_PLAYER) {

                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Media player is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Media player is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("Media player is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("Media player is charged");
                }
        }

        /* tablet */
        if (kind == UP_DEVICE_KIND_TABLET) {

                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Tablet is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Tablet is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("Tablet is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("Tablet is charged");
                }
        }

        /* computer */
        if (kind == UP_DEVICE_KIND_COMPUTER) {

                if (state == UP_DEVICE_STATE_CHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Computer is charging");
                }
                if (state == UP_DEVICE_STATE_DISCHARGING) {
                        /* TRANSLATORS: battery state */
                        return _("Computer is discharging");
                }
                if (state == UP_DEVICE_STATE_EMPTY) {
                        /* TRANSLATORS: battery state */
                        return _("Computer is empty");
                }
                if (state == UP_DEVICE_STATE_FULLY_CHARGED) {
                        /* TRANSLATORS: battery state */
                        return _("Computer is charged");
                }
        }
#endif

        return gpm_device_kind_to_localised_string (kind, 1);
}

static gboolean
parse_vm_kernel_cmdline (gboolean *is_virtual_machine)
{
        gboolean ret = FALSE;
        GRegex *regex;
        GMatchInfo *match;
        char *contents;
        char *word;
        const char *arg;

        if (!g_file_get_contents ("/proc/cmdline", &contents, NULL, NULL))
                return ret;

        regex = g_regex_new ("gnome.is_vm=(\\S+)", 0, G_REGEX_MATCH_NOTEMPTY, NULL);
        if (!g_regex_match (regex, contents, G_REGEX_MATCH_NOTEMPTY, &match))
                goto out;

        word = g_match_info_fetch (match, 0);
        g_debug ("Found command-line match '%s'", word);
        arg = word + strlen ("gnome.is_vm=");
        if (*arg != '0' && *arg != '1') {
                g_warning ("Invalid value '%s' for gnome.is_vm passed in kernel command line.\n", arg);
        } else {
                *is_virtual_machine = atoi (arg);
                ret = TRUE;
        }
        g_free (word);

out:
        g_match_info_free (match);
        g_regex_unref (regex);
        g_free (contents);

        if (ret)
                g_debug ("Kernel command-line parsed to %d", *is_virtual_machine);

        return ret;
}

gboolean
gsd_power_is_hardware_a_vm (void)
{
        const gchar *str;
        gboolean ret = FALSE;
        GError *error = NULL;
        GVariant *inner;
        GVariant *variant = NULL;
        GDBusConnection *connection;

        if (parse_vm_kernel_cmdline (&ret))
                return ret;

        connection = g_bus_get_sync (G_BUS_TYPE_SYSTEM,
                                     NULL,
                                     &error);
        if (connection == NULL) {
                g_warning ("system bus not available: %s", error->message);
                g_error_free (error);
                goto out;
        }
        variant = g_dbus_connection_call_sync (connection,
                                               "org.freedesktop.systemd1",
                                               "/org/freedesktop/systemd1",
                                               "org.freedesktop.DBus.Properties",
                                               "Get",
                                               g_variant_new ("(ss)",
                                                              "org.freedesktop.systemd1.Manager",
                                                              "Virtualization"),
                                               NULL,
                                               G_DBUS_CALL_FLAGS_NONE,
                                               -1,
                                               NULL,
                                               &error);
        if (variant == NULL) {
                g_debug ("Failed to get property '%s': %s", "Virtualization", error->message);
                g_error_free (error);
                goto out;
        }

        /* on bare-metal hardware this is the empty string,
         * otherwise an identifier such as "kvm", "vmware", etc. */
        g_variant_get (variant, "(v)", &inner);
        str = g_variant_get_string (inner, NULL);
        if (str != NULL && str[0] != '\0')
                ret = TRUE;
out:
        if (connection != NULL)
                g_object_unref (connection);
        if (variant != NULL)
                g_variant_unref (variant);
        return ret;
}

/* This timer goes off every few minutes, whether the user is idle or not,
   to try and clean up anything that has gone wrong.

   It calls disable_builtin_screensaver() so that if xset has been used,
   or some other program (like xlock) has messed with the XSetScreenSaver()
   settings, they will be set back to sensible values (if a server extension
   is in use, messing with xlock can cause the screensaver to never get a wakeup
   event, and could cause monitor power-saving to occur, and all manner of
   heinousness.)

   This code was originally part of gnome-screensaver, see
   http://git.gnome.org/browse/gnome-screensaver/tree/src/gs-watcher-x11.c?id=fec00b12ec46c86334cfd36b37771cc4632f0d4d#n530
 */
static gboolean
disable_builtin_screensaver (gpointer unused)
{
        int current_server_timeout, current_server_interval;
        int current_prefer_blank,   current_allow_exp;
        int desired_server_timeout, desired_server_interval;
        int desired_prefer_blank,   desired_allow_exp;

        XGetScreenSaver (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()),
                         &current_server_timeout,
                         &current_server_interval,
                         &current_prefer_blank,
                         &current_allow_exp);

        desired_server_timeout  = current_server_timeout;
        desired_server_interval = current_server_interval;
        desired_prefer_blank    = current_prefer_blank;
        desired_allow_exp       = current_allow_exp;

        desired_server_interval = 0;

        /* I suspect (but am not sure) that DontAllowExposures might have
           something to do with powering off the monitor as well, at least
           on some systems that don't support XDPMS?  Who know... */
        desired_allow_exp = AllowExposures;

        /* When we're not using an extension, set the server-side timeout to 0,
           so that the server never gets involved with screen blanking, and we
           do it all ourselves.  (However, when we *are* using an extension,
           we tell the server when to notify us, and rather than blanking the
           screen, the server will send us an X event telling us to blank.)
        */
        desired_server_timeout = 0;

        if (desired_server_timeout     != current_server_timeout
            || desired_server_interval != current_server_interval
            || desired_prefer_blank    != current_prefer_blank
            || desired_allow_exp       != current_allow_exp) {

                g_debug ("disabling server builtin screensaver:"
                         " (xset s %d %d; xset s %s; xset s %s)",
                         desired_server_timeout,
                         desired_server_interval,
                         (desired_prefer_blank ? "blank" : "noblank"),
                         (desired_allow_exp ? "expose" : "noexpose"));

                XSetScreenSaver (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()),
                                 desired_server_timeout,
                                 desired_server_interval,
                                 desired_prefer_blank,
                                 desired_allow_exp);

                XSync (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), FALSE);
        }

        return TRUE;
}

guint
gsd_power_enable_screensaver_watchdog (void)
{
        int dummy;

        /* Make sure that Xorg's DPMS extension never gets in our
         * way. The defaults are now applied in Fedora 20 from
         * being "0" by default to being "600" by default */
        gdk_error_trap_push ();
        if (DPMSQueryExtension(GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), &dummy, &dummy))
                DPMSSetTimeouts (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), 0, 0, 0);
        gdk_error_trap_pop_ignored ();
        return g_timeout_add_seconds (XSCREENSAVER_WATCHDOG_TIMEOUT,
                                      disable_builtin_screensaver,
                                      NULL);
}

static GnomeRROutput *
get_primary_output (GnomeRRScreen *rr_screen)
{
        GnomeRROutput *output = NULL;
        GnomeRROutput **outputs;
        guint i;

        /* search all X11 outputs for the device id */
        outputs = gnome_rr_screen_list_outputs (rr_screen);
        if (outputs == NULL)
                goto out;

        for (i = 0; outputs[i] != NULL; i++) {
                if (gnome_rr_output_is_connected (outputs[i]) &&
                    gnome_rr_output_is_laptop (outputs[i]) &&
                    gnome_rr_output_get_backlight_min (outputs[i]) >= 0 &&
                    gnome_rr_output_get_backlight_max (outputs[i]) > 0) {
                        output = outputs[i];
                        break;
                }
        }
out:
        return output;
}

#ifdef GSD_MOCK
static void
backlight_set_mock_value (gint value)
{
	const char *filename;
	char *contents;

	g_debug ("Settings mock brightness: %d", value);

	filename = "GSD_MOCK_brightness";
	contents = g_strdup_printf ("%d", value);
	g_file_set_contents (filename, contents, -1, NULL);
	g_free (contents);
}

static gint64
backlight_get_mock_value (const char *argument)
{
	const char *filename;
	char *contents;
	gint64 ret;

	if (g_str_equal (argument, "get-max-brightness")) {
		g_debug ("Returning max mock brightness: %d", GSD_MOCK_MAX_BRIGHTNESS);
		return GSD_MOCK_MAX_BRIGHTNESS;
	}

	if (g_str_equal (argument, "get-brightness")) {
		filename = "GSD_MOCK_brightness";
		ret = GSD_MOCK_DEFAULT_BRIGHTNESS;
	} else {
		g_assert_not_reached ();
	}

	if (g_file_get_contents (filename, &contents, NULL, NULL)) {
		ret = g_ascii_strtoll (contents, NULL, 0);
		g_free (contents);
		g_debug ("Returning mock brightness: %"G_GINT64_FORMAT, ret);
	} else {
		ret = GSD_MOCK_DEFAULT_BRIGHTNESS;
		backlight_set_mock_value (GSD_MOCK_DEFAULT_BRIGHTNESS);
		g_debug ("Returning default mock brightness: %"G_GINT64_FORMAT, ret);
	}

	return ret;
}
#endif /* GSD_MOCK */

gboolean
backlight_available (GnomeRRScreen *rr_screen)
{
        char *path;

#ifdef GSD_MOCK
	return TRUE;
#endif
        if (get_primary_output (rr_screen) != NULL)
                return TRUE;
        path = gsd_backlight_helper_get_best_backlight ();
        if (path == NULL)
                return FALSE;

        g_free (path);
        return TRUE;
}

/**
 * backlight_helper_get_value:
 *
 * Gets a brightness value from the PolicyKit helper.
 *
 * Return value: the signed integer value from the helper, or -1
 * for failure. If -1 then @error is set.
 **/
static gint64
backlight_helper_get_value (const gchar *argument, GError **error)
{
        gboolean ret;
        gchar *stdout_data = NULL;
        gint exit_status = 0;
        gint64 value = -1;
        gchar *command = NULL;
        gchar *endptr = NULL;

#ifdef GSD_MOCK
        return backlight_get_mock_value (argument);
#endif

#ifndef __linux__
        /* non-Linux platforms won't have /sys/class/backlight */
        g_set_error_literal (error,
                             GSD_POWER_MANAGER_ERROR,
                             GSD_POWER_MANAGER_ERROR_FAILED,
                             "The sysfs backlight helper is only for Linux");
        goto out;
#endif

        /* get the data */
        command = g_strdup_printf (LIBEXECDIR "/gsd-backlight-helper --%s",
                                   argument);
        ret = g_spawn_command_line_sync (command,
                                         &stdout_data,
                                         NULL,
                                         &exit_status,
                                         error);
        g_debug ("executed %s retval: %i", command, exit_status);

        if (!ret)
                goto out;

        if (WEXITSTATUS (exit_status) != 0) {
                 g_set_error (error,
                             GSD_POWER_MANAGER_ERROR,
                             GSD_POWER_MANAGER_ERROR_FAILED,
                             "gsd-backlight-helper failed: %s",
                             stdout_data ? stdout_data : "No reason");
                goto out;
        }

        /* parse */
        value = g_ascii_strtoll (stdout_data, &endptr, 10);

        /* parsing error */
        if (endptr == stdout_data) {
                value = -1;
                g_set_error (error,
                             GSD_POWER_MANAGER_ERROR,
                             GSD_POWER_MANAGER_ERROR_FAILED,
                             "failed to parse value: %s",
                             stdout_data);
                goto out;
        }

        /* out of range */
        if (value > G_MAXINT) {
                value = -1;
                g_set_error (error,
                             GSD_POWER_MANAGER_ERROR,
                             GSD_POWER_MANAGER_ERROR_FAILED,
                             "value out of range: %s",
                             stdout_data);
                goto out;
        }

        /* Fetching the value failed, for some other reason */
        if (value < 0) {
                g_set_error (error,
                             GSD_POWER_MANAGER_ERROR,
                             GSD_POWER_MANAGER_ERROR_FAILED,
                             "value negative, but helper did not fail: %s",
                             stdout_data);
                goto out;
        }

out:
        g_free (command);
        g_free (stdout_data);
        return value;
}

/**
 * backlight_helper_set_value:
 *
 * Sets a brightness value using the PolicyKit helper.
 *
 * Return value: Success. If FALSE then @error is set.
 **/
static gboolean
backlight_helper_set_value (const gchar *argument,
                            gint value,
                            GError **error)
{
        gboolean ret = FALSE;
        gint exit_status = 0;
        gchar *command = NULL;

#ifdef GSD_MOCK
	backlight_set_mock_value (value);
	return TRUE;
#endif

#ifndef __linux__
        /* non-Linux platforms won't have /sys/class/backlight */
        g_set_error_literal (error,
                             GSD_POWER_MANAGER_ERROR,
                             GSD_POWER_MANAGER_ERROR_FAILED,
                             "The sysfs backlight helper is only for Linux");
        goto out;
#endif

        /* get the data */
        command = g_strdup_printf ("pkexec " LIBEXECDIR "/gsd-backlight-helper --%s %i",
                                   argument, value);
        ret = g_spawn_command_line_sync (command,
                                         NULL,
                                         NULL,
                                         &exit_status,
                                         error);

        g_debug ("executed %s retval: %i", command, exit_status);

        if (!ret || WEXITSTATUS (exit_status) != 0)
                goto out;

out:
        g_free (command);
        return ret;
}

int
backlight_get_abs (GnomeRRScreen *rr_screen, GError **error)
{
        GnomeRROutput *output;

        /* prefer xbacklight */
        output = get_primary_output (rr_screen);
        if (output != NULL) {
                return gnome_rr_output_get_backlight (output,
                                                      error);
        }

        /* fall back to the polkit helper */
        return backlight_helper_get_value ("get-brightness", error);
}

int
backlight_get_percentage (GnomeRRScreen *rr_screen, GError **error)
{
        GnomeRROutput *output;
        gint now;
        gint value = -1;
        gint min = 0;
        gint max;

        /* prefer xbacklight */
        output = get_primary_output (rr_screen);
        if (output != NULL) {

                min = gnome_rr_output_get_backlight_min (output);
                max = gnome_rr_output_get_backlight_max (output);
                now = gnome_rr_output_get_backlight (output, error);
                if (now < 0)
                        goto out;
                value = ABS_TO_PERCENTAGE (min, max, now);
                goto out;
        }

        /* fall back to the polkit helper */
        max = backlight_helper_get_value ("get-max-brightness", error);
        if (max < 0)
                goto out;
        now = backlight_helper_get_value ("get-brightness", error);
        if (now < 0)
                goto out;
        value = ABS_TO_PERCENTAGE (min, max, now);
out:
        return value;
}

int
backlight_get_min (GnomeRRScreen *rr_screen)
{
        GnomeRROutput *output;

        /* if we have no xbacklight device, then hardcode zero as sysfs
         * offsets everything to 0 as min */
        output = get_primary_output (rr_screen);
        if (output == NULL)
                return 0;

        /* get xbacklight value, which maybe non-zero */
        return gnome_rr_output_get_backlight_min (output);
}

int
backlight_get_max (GnomeRRScreen *rr_screen, GError **error)
{
        gint value;
        GnomeRROutput *output;

        /* prefer xbacklight */
        output = get_primary_output (rr_screen);
        if (output != NULL) {
                value = gnome_rr_output_get_backlight_max (output);
                if (value < 0) {
                        g_set_error (error,
                                     GSD_POWER_MANAGER_ERROR,
                                     GSD_POWER_MANAGER_ERROR_FAILED,
                                     "failed to get backlight max");
                }
                return value;
        }

        /* fall back to the polkit helper */
        return  backlight_helper_get_value ("get-max-brightness", error);
}

gboolean
backlight_set_percentage (GnomeRRScreen *rr_screen,
                          guint value,
                          GError **error)
{
        GnomeRROutput *output;
        gboolean ret = FALSE;
        gint min = 0;
        gint max;
        guint discrete;

        /* prefer xbacklight */
        output = get_primary_output (rr_screen);
        if (output != NULL) {
                min = gnome_rr_output_get_backlight_min (output);
                max = gnome_rr_output_get_backlight_max (output);
                if (min < 0 || max < 0) {
                        g_warning ("no xrandr backlight capability");
                        return ret;
                }
                discrete = PERCENTAGE_TO_ABS (min, max, value);
                ret = gnome_rr_output_set_backlight (output,
                                                     discrete,
                                                     error);
                return ret;
        }

        /* fall back to the polkit helper */
        max = backlight_helper_get_value ("get-max-brightness", error);
        if (max < 0)
                return ret;
        discrete = PERCENTAGE_TO_ABS (min, max, value);
        ret = backlight_helper_set_value ("set-brightness",
                                          discrete,
                                          error);

        return ret;
}

int
backlight_step_up (GnomeRRScreen *rr_screen, GError **error)
{
        GnomeRROutput *output;
        gboolean ret = FALSE;
        gint percentage_value = -1;
        gint min = 0;
        gint max;
        gint now;
        gint step;
        guint discrete;
        GnomeRRCrtc *crtc;

        /* prefer xbacklight */
        output = get_primary_output (rr_screen);
        if (output != NULL) {

                crtc = gnome_rr_output_get_crtc (output);
                if (crtc == NULL) {
                        g_set_error (error,
                                     GSD_POWER_MANAGER_ERROR,
                                     GSD_POWER_MANAGER_ERROR_FAILED,
                                     "no crtc for %s",
                                     gnome_rr_output_get_name (output));
                        return percentage_value;
                }
                min = gnome_rr_output_get_backlight_min (output);
                max = gnome_rr_output_get_backlight_max (output);
                now = gnome_rr_output_get_backlight (output, error);
                if (now < 0)
                       return percentage_value;
                step = BRIGHTNESS_STEP_AMOUNT (max - min + 1);
                discrete = MIN (now + step, max);
                ret = gnome_rr_output_set_backlight (output,
                                                     discrete,
                                                     error);
                if (ret)
                        percentage_value = ABS_TO_PERCENTAGE (min, max, discrete);
                return percentage_value;
        }

        /* fall back to the polkit helper */
        now = backlight_helper_get_value ("get-brightness", error);
        if (now < 0)
                return percentage_value;
        max = backlight_helper_get_value ("get-max-brightness", error);
        if (max < 0)
                return percentage_value;
        step = BRIGHTNESS_STEP_AMOUNT (max - min + 1);
        discrete = MIN (now + step, max);
        ret = backlight_helper_set_value ("set-brightness",
                                          discrete,
                                          error);
        if (ret)
                percentage_value = ABS_TO_PERCENTAGE (min, max, discrete);

        return percentage_value;
}

int
backlight_step_down (GnomeRRScreen *rr_screen, GError **error)
{
        GnomeRROutput *output;
        gboolean ret = FALSE;
        gint percentage_value = -1;
        gint min = 0;
        gint max;
        gint now;
        gint step;
        guint discrete;
        GnomeRRCrtc *crtc;

        /* prefer xbacklight */
        output = get_primary_output (rr_screen);
        if (output != NULL) {

                crtc = gnome_rr_output_get_crtc (output);
                if (crtc == NULL) {
                        g_set_error (error,
                                     GSD_POWER_MANAGER_ERROR,
                                     GSD_POWER_MANAGER_ERROR_FAILED,
                                     "no crtc for %s",
                                     gnome_rr_output_get_name (output));
                        return percentage_value;
                }
                min = gnome_rr_output_get_backlight_min (output);
                max = gnome_rr_output_get_backlight_max (output);
                now = gnome_rr_output_get_backlight (output, error);
                if (now < 0)
                       return percentage_value;
                step = BRIGHTNESS_STEP_AMOUNT (max - min + 1);
                discrete = MAX (now - step, 0);
                ret = gnome_rr_output_set_backlight (output,
                                                     discrete,
                                                     error);
                if (ret)
                        percentage_value = ABS_TO_PERCENTAGE (min, max, discrete);
                return percentage_value;
        }

        /* fall back to the polkit helper */
        now = backlight_helper_get_value ("get-brightness", error);
        if (now < 0)
                return percentage_value;
        max = backlight_helper_get_value ("get-max-brightness", error);
        if (max < 0)
                return percentage_value;
        step = BRIGHTNESS_STEP_AMOUNT (max - min + 1);
        discrete = MAX (now - step, 0);
        ret = backlight_helper_set_value ("set-brightness",
                                          discrete,
                                          error);
        if (ret)
                percentage_value = ABS_TO_PERCENTAGE (min, max, discrete);

        return percentage_value;
}

int
backlight_set_abs (GnomeRRScreen *rr_screen,
                   guint value,
                   GError **error)
{
        GnomeRROutput *output;
        gboolean ret = FALSE;

        /* prefer xbacklight */
        output = get_primary_output (rr_screen);
        if (output != NULL) {
                ret = gnome_rr_output_set_backlight (output,
                                                     value,
                                                     error);
                return ret;
        }

        /* fall back to the polkit helper */
        ret = backlight_helper_set_value ("set-brightness",
                                          value,
                                          error);

        return ret;
}

void
reset_idletime (void)
{
        static gboolean inited = FALSE;
        static KeyCode keycode1, keycode2;
        static gboolean first_keycode = FALSE;

        if (inited == FALSE) {
                keycode1 = XKeysymToKeycode (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), GDK_KEY_Alt_L);
                keycode2 = XKeysymToKeycode (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), GDK_KEY_Alt_R);
        }

        gdk_error_trap_push ();
        /* send a left or right alt key; first press, then release */
        XTestFakeKeyEvent (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), first_keycode ? keycode1 : keycode2, True, CurrentTime);
        XTestFakeKeyEvent (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), first_keycode ? keycode1 : keycode2, False, CurrentTime);
        first_keycode = !first_keycode;
        gdk_error_trap_pop_ignored ();
}

static gboolean
randr_output_is_on (GnomeRROutput *output)
{
        GnomeRRCrtc *crtc;

        crtc = gnome_rr_output_get_crtc (output);
        if (!crtc)
                return FALSE;
        return gnome_rr_crtc_get_current_mode (crtc) != NULL;
}

gboolean
external_monitor_is_connected (GnomeRRScreen *screen)
{
        GnomeRROutput **outputs;
        guint i;

#ifdef GSD_MOCK
	char *mock_external_monitor_contents;

	if (g_file_get_contents ("GSD_MOCK_EXTERNAL_MONITOR", &mock_external_monitor_contents, NULL, NULL)) {
		if (mock_external_monitor_contents[0] == '1') {
			g_free (mock_external_monitor_contents);
			return TRUE;
		} else if (mock_external_monitor_contents[0] == '0') {
			g_free (mock_external_monitor_contents);
			return FALSE;
		}

		g_warning ("Unhandled value for GSD_MOCK_EXTERNAL_MONITOR contents: %s", mock_external_monitor_contents);
		g_free (mock_external_monitor_contents);
	}
#endif /* GSD_MOCK */

        /* see if we have more than one screen plugged in */
        outputs = gnome_rr_screen_list_outputs (screen);
        for (i = 0; outputs[i] != NULL; i++) {
                if (randr_output_is_on (outputs[i]) &&
                    !gnome_rr_output_is_laptop (outputs[i]))
                        return TRUE;
        }

        return FALSE;
}

static void
play_sound (void)
{
        ca_context_play (ca_gtk_context_get (), UPS_SOUND_LOOP_ID,
                         CA_PROP_EVENT_ID, "battery-caution",
                         CA_PROP_EVENT_DESCRIPTION, _("Battery is critically low"), NULL);
}

static gboolean
play_loop_timeout_cb (gpointer user_data)
{
        play_sound ();
        return TRUE;
}

void
play_loop_start (guint *id)
{
        if (*id != 0)
                return;

        *id = g_timeout_add_seconds (GSD_POWER_MANAGER_CRITICAL_ALERT_TIMEOUT,
                                     (GSourceFunc) play_loop_timeout_cb,
                                     NULL);
        play_sound ();
}

void
play_loop_stop (guint *id)
{
        if (*id == 0)
                return;

        ca_context_cancel (ca_gtk_context_get (), UPS_SOUND_LOOP_ID);
        g_source_remove (*id);
        *id = 0;
}
