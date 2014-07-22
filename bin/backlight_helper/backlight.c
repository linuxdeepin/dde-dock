#include <libudev.h>
#include <stdlib.h>
#include <stdio.h>
#include "backlight.h"

static struct udev_device* cached_dev = NULL;

struct udev_device* filter_by_type(struct udev* udev, struct udev_list_entry* entries, const char* type)
{
    struct udev_list_entry* current;
    udev_list_entry_foreach(current, entries) {
	const char* name = udev_list_entry_get_name(current);
	struct udev_device* dev = udev_device_new_from_syspath(udev, name);
	if (strcmp(udev_device_get_sysattr_value(dev, "type"), "firmware") == 0) {
	    return dev;
	}
	udev_device_unref(dev);
    }
    return NULL;
}

void set_cached_dev(struct udev_device* dev)
{
    if (cached_dev != NULL) {
	udev_device_unref(cached_dev);
    }
    cached_dev = dev;
    printf("Found backlight device: %s\n", udev_device_get_syspath(dev));
}

void update_backlight_device()
{
    struct udev* udev = udev_new();
    struct udev_enumerate* enumerate = udev_enumerate_new(udev);

    udev_enumerate_add_match_subsystem(enumerate, "backlight");

    udev_enumerate_scan_devices(enumerate);
    struct udev_list_entry* entries = udev_enumerate_get_list_entry(enumerate);

    struct udev_device* dev = filter_by_type(udev, entries, "firmware");

    if (dev == NULL){
	dev = filter_by_type(udev, entries, "platform");
    } else {
	set_cached_dev(dev);

	udev_enumerate_unref(enumerate);
	udev_unref(udev);
	return;
    }

    if (dev == NULL){
	dev = filter_by_type(udev, entries, "raw");
    } else {
	set_cached_dev(dev);

	udev_enumerate_unref(enumerate);
	udev_unref(udev);
	return;
    }

    set_cached_dev(dev);
    udev_enumerate_unref(enumerate);
    udev_unref(udev);
}

double get_backlight() 
{
    if (cached_dev == NULL) {
	return -1;
    }
    const char* str_v = udev_device_get_sysattr_value(cached_dev, "brightness");
    const char* str_max  = udev_device_get_sysattr_value(cached_dev, "max_brightness");
    if (str_v == NULL || str_max == NULL) {
	return -1;
    }
    int v = atoi(str_v);
    int max = atoi(str_max);
    if (max < v || max == 0) {
	return -1;
    }

    return v * 1.0 / max;
}

void set_backlight(double v)
{
    printf("SetBBB: %lf\n", v);
    if (v > 1 || v < 0) {
	fprintf(stderr, "set_backlight(%lf) failed\n", v);
	return;
    }
    const char* str_max  = udev_device_get_sysattr_value(cached_dev, "max_brightness");
    if (str_max == NULL) {
	fprintf(stderr, "set_backlight(%lf) failed\n", v);
	return;
    }
    char str_v[1000] = {0};
    sprintf(str_v, "%d", (int)(v * atoi(str_max)));
    int r = udev_device_set_sysattr_value(cached_dev, "brightness", str_v);
    if (r == -1) {
	fprintf(stderr, "set_backlight(%lf) failed\n", v);
    }
}
