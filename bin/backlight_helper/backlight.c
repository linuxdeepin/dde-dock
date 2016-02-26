/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <libudev.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "backlight.h"

struct udev_enumerate* cached_enumerate = NULL;

struct udev_device* filter_by_type(struct udev* udev, struct udev_list_entry* entries, const char* type)
{
    struct udev_list_entry* current;
    udev_list_entry_foreach(current, entries) {
	const char* name = udev_list_entry_get_name(current);
	struct udev_device* dev = udev_device_new_from_syspath(udev, name);
	if (strcmp(udev_device_get_sysattr_value(dev, "type"), type) == 0) {
	    return dev;
	}
	udev_device_unref(dev);
    }
    return NULL;
}

struct udev_device*
get_device_by_type(char* driver_type)
{
    struct udev* udev = udev_new();

    if (cached_enumerate != NULL) {
        udev_enumerate_unref(cached_enumerate);
        cached_enumerate = NULL;
    }
    cached_enumerate = udev_enumerate_new(udev);

    udev_enumerate_add_match_subsystem(cached_enumerate, "backlight");
    udev_enumerate_scan_devices(cached_enumerate);
    struct udev_list_entry* entries = udev_enumerate_get_list_entry(cached_enumerate);
    struct udev_device* dev = NULL;

    if (strcmp(driver_type, "backlight-raw") == 0 ) {
        dev = filter_by_type(udev, entries, "raw");
    } else if (strcmp(driver_type, "backlight-platform") == 0) {
        dev = filter_by_type(udev, entries, "platform");
    } else if (strcmp(driver_type, "backlight-firmware") == 0 ) {
        dev = filter_by_type(udev, entries, "firmware");
    }
    if (dev != NULL) {
        goto out;
    }

    // Auto detect
    dev = filter_by_type(udev, entries, "firmware");
    if (dev != NULL) {
        goto out;
    }

    dev = filter_by_type(udev, entries, "raw");
    if (dev != NULL) {
        goto out;
    }

    dev = filter_by_type(udev, entries, "platform");
    if (dev != NULL) {
        goto out;
    }

out:
    udev_unref(udev);
    return dev;
}

double get_backlight(char* driver_type)
{
    struct udev_device* driver = get_device_by_type(driver_type);
    if (driver == NULL) {
	return -1;
    }

    struct udev* udev = udev_device_get_udev(driver);
    struct udev_device* dev = udev_device_new_from_syspath(udev, udev_device_get_syspath(driver));
    const char* str_v = udev_device_get_sysattr_value(dev, "brightness");
    const char* str_max  = udev_device_get_sysattr_value(dev, "max_brightness");

    double ret = -1.0;
    if (str_v == NULL || str_max == NULL) {
        goto out;
    }
    int v = atoi(str_v);
    int max = atoi(str_max);
    if (max < v || max == 0) {
	goto out;
    }

    ret = v * 1.0 / max;

out:
    udev_device_unref(dev);
    udev_device_unref(driver);
    return ret;
}

void set_backlight(double v, char* driver_type)
{
    if (v > 1 || v < 0) {
	fprintf(stderr, "set_backlight(%lf) type(%s) failed\n", v, driver_type);
	return;
    }

    struct udev_device* driver = get_device_by_type(driver_type);
    if (driver == NULL) {
        fprintf(stderr, "Get udev driver for '%s' failed\n", driver_type);
        return;
    }

    struct udev* udev = udev_device_get_udev(driver);
    struct udev_list_entry* entries = udev_enumerate_get_list_entry(cached_enumerate);

    struct udev_list_entry* current;
    udev_list_entry_foreach(current, entries) {
	const char* name = udev_list_entry_get_name(current);
	struct udev_device* dev = udev_device_new_from_syspath(udev, name);

        const char* str_max = udev_device_get_sysattr_value(dev, "max_brightness");
        if (str_max == NULL) {
            fprintf(stderr, "get max_brightness failed(driver:%s)\n", name);
            udev_device_unref(dev);
            continue;
        }
        char str_v[1000] = {0};
        sprintf(str_v, "%d", (int)(v * atoi(str_max)));
        int r = udev_device_set_sysattr_value(dev, "brightness", str_v);
        if (r != 0) {
            fprintf(stderr, "set_backlight to %lf(%s/%s) %d failed(driver:%s)\n", v, str_v, str_max, r, name);
            udev_device_unref(dev);
            continue;
        }

        fprintf(stdout, "set_backlight to %lf(%s) (driver:%s)\n", v, str_v, name);
	udev_device_unref(dev);
    }

    udev_device_unref(driver);
}
