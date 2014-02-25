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

#ifndef __GPMCOMMON_H
#define __GPMCOMMON_H

#include <glib.h>
#include <libupower-glib/upower.h>

G_BEGIN_DECLS

/* UPower helpers */
gchar           *gpm_get_timestring                     (guint           time);
const gchar     *gpm_device_to_localised_string         (UpDevice       *device);
const gchar     *gpm_device_kind_to_localised_string    (UpDeviceKind    kind,
                                                         guint           number);
const gchar     *gpm_device_kind_to_icon                (UpDeviceKind    kind);
const gchar     *gpm_device_technology_to_localised_string (UpDeviceTechnology technology_enum);
const gchar     *gpm_device_state_to_localised_string   (UpDeviceState   state);
GIcon           *gpm_upower_get_device_icon             (UpDevice       *device,
                                                         gboolean        use_symbolic);
gchar           *gpm_upower_get_device_summary          (UpDevice       *device);
gchar           *gpm_upower_get_device_description      (UpDevice       *device);

/* Power helpers */
gboolean         gsd_power_is_hardware_a_vm             (void);
guint            gsd_power_enable_screensaver_watchdog  (void);
void             reset_idletime                         (void);

/* Backlight helpers */

/* on ACPI machines we have 4-16 levels, on others it's ~150 */
#define BRIGHTNESS_STEP_AMOUNT(max) ((max) < 20 ? 1 : (max) / 20)

#define ABS_TO_PERCENTAGE(min, max, value) gsd_power_backlight_abs_to_percentage(min, max, value)
#define PERCENTAGE_TO_ABS(min, max, value) (min + (((max - min) * value) / 100))

int              gsd_power_backlight_abs_to_percentage  (int min, int max, int value);
gboolean         backlight_available                    (GnomeRRScreen *rr_screen);
int              backlight_get_abs                      (GnomeRRScreen *rr_screen, GError **error);
int              backlight_get_percentage               (GnomeRRScreen *rr_screen, GError **error);
int              backlight_get_min                      (GnomeRRScreen *rr_screen);
int              backlight_get_max                      (GnomeRRScreen *rr_screen, GError **error);
gboolean         backlight_set_percentage               (GnomeRRScreen *rr_screen,
                                                         guint value,
                                                         GError **error);
int              backlight_step_up                      (GnomeRRScreen *rr_screen, GError **error);
int              backlight_step_down                    (GnomeRRScreen *rr_screen, GError **error);
int              backlight_set_abs                      (GnomeRRScreen *rr_screen,
                                                         guint value,
                                                         GError **error);

/* RandR helpers */
gboolean         external_monitor_is_connected          (GnomeRRScreen *screen);

/* Sound helpers */
void             play_loop_start                        (guint *id);
void             play_loop_stop                         (guint *id);

G_END_DECLS

#endif  /* __GPMCOMMON_H */
