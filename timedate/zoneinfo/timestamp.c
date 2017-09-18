/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include <stdio.h>	/* for stdout, stderr, perror */
#include <string.h>	/* for strcpy */
#include <sys/types.h>	/* for time_t */
#include <stdlib.h>	/* for exit, malloc, atoi */
#include <limits.h>	/* for CHAR_BIT, LLONG_MAX */
#include <ctype.h>	/* for isalpha et al. */
#include "zdump.h"
#include "timestamp.h"

static int is_usec_has_dst(time_t t);
static void reset_tz(char* value);

DSTTime
get_dst_time(const char *zone, int year)
{
	register time_t		cutlotime;
	register time_t		cuthitime;

	cutlotime = yeartot(year);
	cuthitime = yeartot(year+1);

	time_t			t;
	time_t			newt;
	struct tm		tm;
	struct tm		newtm;
	register struct tm *	tmp;
	register struct tm *	newtmp;
	static char	buf[MAX_STRING_LENGTH];
	DSTTime ret = {0, 0};

	char *tz = getenv("TZ");
	if (setenv("TZ", zone, 1) != 0 ) {
		fprintf(stderr, "Set TZ=%s failed\n", zone);
		return ret;
	}

	t = absolute_min_time;
	if (t < cutlotime)
		t = cutlotime;
	tmp = localtime(&t);
	if (tmp != NULL) {
		tm = *tmp;
		strncpy(buf, abbr(&tm), (sizeof buf) - 1);
	}

	int enter_flag = 0;
	for ( ; ; ) {
		newt = (t < absolute_max_time - SECSPERDAY / 2
		        ? t + SECSPERDAY / 2
		        : absolute_max_time);
		if (cuthitime <= newt)
			break;
		newtmp = localtime(&newt);
		if (newtmp != NULL)
			newtm = *newtmp;
		if ((tmp == NULL || newtmp == NULL) ? (tmp != newtmp) :
		        (delta(&newtm, &tm) != (newt - t) ||
		         newtm.tm_isdst != tm.tm_isdst ||
		         strcmp(abbr(&newtm), buf) != 0)) {
			newt = hunt(t, newt);
			newtmp = localtime(&newt);
			if (newtmp != NULL) {
				if (newtmp->tm_isdst) {
					if (!enter_flag) {
						ret.enter = newt;
						enter_flag = 1;
					} else {
						ret.leave = newt;
					}
				} else {
					if (!enter_flag) {
						ret.enter = newt-1;
						enter_flag = 1;
					} else {
						ret.leave = newt-1;
					}
				}

				newtm = *newtmp;
				strncpy(buf, abbr(&newtm),(sizeof buf) - 1);
			}
		}
		t = newt;
		tm = newtm;
		tmp = newtmp;
	}
       reset_tz(tz);

	return ret;
}

long long
get_rawoffset_usec (const char *zone, long long t)
{
	char *tz = getenv("TZ");
	setenv("TZ", zone, 1);

	time_t newt = t-1;
	if (!is_usec_has_dst(newt)) {
               reset_tz(tz);
		return newt;
	}

	newt = t+1;
	if (!is_usec_has_dst(newt)) {
               reset_tz(tz);
		return newt;
	}

       reset_tz(tz);
	return 0;
}

long
get_offset_by_usec (const char *zone, long long t)
{
	char *tz = getenv("TZ");
	setenv("TZ", zone, 1);

	struct tm *tmp = NULL;
	if (t == 0) {
		time_t newt = time(NULL);
		tmp = localtime(&newt);
	} else {
		tmp = localtime((time_t*)&t);
	}
       reset_tz(tz);
	if (!tmp) {
		return -1;
	}

# if defined __USE_MISC
	return tmp->tm_gmtoff;
# else
	return tmp->__tm_gmtoff;
#endif
}

static int
is_usec_has_dst(time_t t)
{
	struct tm *tmp = NULL;
	tmp = localtime(&t);
	if (tmp != NULL && tmp->tm_isdst) {
		return 1;
	}

	return 0;
}

static void
reset_tz(char* value)
{
       if (value) {
               setenv("TZ", value, 1);
       } else {
               unsetenv("TZ");
       }
}
