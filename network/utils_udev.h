
/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <stdlib.h>
#include <libudev.h>

int is_device_has_property(struct udev_device *device, const char *property);
char *get_device_vendor(const char *syspath);
char *get_device_product(const char *syspath);
int is_usb_device(const char *syspath);
