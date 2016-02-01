/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __TIMESTAMP_H__
#define __TIMESTAMP_H__

typedef struct _DSTTime {
	long long enter;
	long long leave;
} DSTTime;

long long get_rawoffset_usec (const char *zone, long long t);
long get_offset_by_usec (const char *zone, long long t);
DSTTime get_dst_time(const char *zone, int year);

#endif
