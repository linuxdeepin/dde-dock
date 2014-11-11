/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

#include <stdio.h>	/* for stdout, stderr, perror */
#include <string.h>	/* for strcpy */
#include <sys/types.h>	/* for time_t */
#include <time.h>	/* for struct tm */
#include <stdlib.h>	/* for exit, malloc, atoi */
#include <limits.h>	/* for CHAR_BIT, LLONG_MAX */
#include <ctype.h>	/* for isalpha et al. */

#include "timestamp.h"

#ifndef INT_FAST32_MAX
# if INT_MAX >> 31 == 0
typedef long int_fast32_t;
# else
typedef int int_fast32_t;
# endif
#endif

#ifndef INTMAX_MAX
# if defined LLONG_MAX || defined __LONG_LONG_MAX__
typedef long long intmax_t;
#  ifdef LLONG_MAX
#   define INTMAX_MAX LLONG_MAX
#  else
#   define INTMAX_MAX __LONG_LONG_MAX__
#  endif
# else
typedef long intmax_t;
#  define INTMAX_MAX LONG_MAX
# endif
#endif

#ifndef ZDUMP_LO_YEAR
#define ZDUMP_LO_YEAR	(-500)
#endif /* !defined ZDUMP_LO_YEAR */

#ifndef ZDUMP_HI_YEAR
#define ZDUMP_HI_YEAR	2500
#endif /* !defined ZDUMP_HI_YEAR */

#ifndef MAX_STRING_LENGTH
#define MAX_STRING_LENGTH	1024
#endif /* !defined MAX_STRING_LENGTH */

#ifndef SECSPERMIN
#define SECSPERMIN	60
#endif /* !defined SECSPERMIN */

#ifndef MINSPERHOUR
#define MINSPERHOUR	60
#endif /* !defined MINSPERHOUR */

#ifndef SECSPERHOUR
#define SECSPERHOUR	(SECSPERMIN * MINSPERHOUR)
#endif /* !defined SECSPERHOUR */

#ifndef HOURSPERDAY
#define HOURSPERDAY	24
#endif /* !defined HOURSPERDAY */

#ifndef EPOCH_YEAR
#define EPOCH_YEAR	1970
#endif /* !defined EPOCH_YEAR */

#ifndef TM_YEAR_BASE
#define TM_YEAR_BASE	1900
#endif /* !defined TM_YEAR_BASE */

#ifndef DAYSPERNYEAR
#define DAYSPERNYEAR	365
#endif /* !defined DAYSPERNYEAR */

#ifndef isleap
#define isleap(y) (((y) % 4) == 0 && (((y) % 100) != 0 || ((y) % 400) == 0))
#endif /* !defined isleap */

#ifndef isleap_sum
/*
** See tzfile.h for details on isleap_sum.
*/
#define isleap_sum(a, b)	isleap((a) % 400 + (b) % 400)
#endif /* !defined isleap_sum */

#define SECSPERDAY	((int_fast32_t) SECSPERHOUR * HOURSPERDAY)
#define SECSPERNYEAR	(SECSPERDAY * DAYSPERNYEAR)
#define SECSPERLYEAR	(SECSPERNYEAR + SECSPERDAY)
#define SECSPER400YEARS	(SECSPERNYEAR * (intmax_t) (300 + 3)	\
			 + SECSPERLYEAR * (intmax_t) (100 - 3))

/*
** True if SECSPER400YEARS is known to be representable as an
** intmax_t.  It's OK that SECSPER400YEARS_FITS can in theory be false
** even if SECSPER400YEARS is representable, because when that happens
** the code merely runs a bit more slowly, and this slowness doesn't
** occur on any practical platform.
*/
enum { SECSPER400YEARS_FITS = SECSPERLYEAR <= INTMAX_MAX / 400 };

extern char *	tzname[2];

/* The minimum and maximum finite time values.  */
static time_t const absolute_min_time =
    ((time_t) -1 < 0
     ? (time_t) -1 << (CHAR_BIT * sizeof (time_t) - 1)
     : 0);
static time_t const absolute_max_time =
    ((time_t) -1 < 0
     ? - (~ 0 < 0) - ((time_t) -1 << (CHAR_BIT * sizeof (time_t) - 1))
     : -1);

static char *	abbr(struct tm * tmp);
static intmax_t	delta(struct tm * newp, struct tm * oldp);
static time_t	hunt(time_t lot, time_t	hit);
static time_t	yeartot(intmax_t y);

long long*
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

	if (setenv("TZ", zone, 1) != 0 ) {
		printf("Set TZ=%s failed\n", zone);
		return NULL;
	}

	t = absolute_min_time;
	if (t < cutlotime)
		t = cutlotime;
	tmp = localtime(&t);
	if (tmp != NULL) {
		tm = *tmp;
		strncpy(buf, abbr(&tm), (sizeof buf) - 1);
	}

	long long *ret = calloc(3, sizeof(long long));
	if (!ret) {
		return NULL;
	}
	ret[2] = 0;

	int idx = 0;
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
			if (idx < 2) {
				ret[idx++] = newt;
				ret[2]++;
			}
			newtmp = localtime(&newt);
			if (newtmp != NULL) {
				newtm = *newtmp;
				strncpy(buf,
				        abbr(&newtm),
				        (sizeof buf) - 1);
			}
		}
		t = newt;
		tm = newtm;
		tmp = newtmp;
	}

	return ret;
}

long long
get_year_begin_time(const char *zone, int year)
{
	// TODO: check return val
	setenv("TZ", zone, 1);

	struct tm tm;

	tm.tm_sec = 0;
	tm.tm_min = 0;
	tm.tm_hour = 0;
	tm.tm_mday = 1;
	tm.tm_mon = 1;
	tm.tm_year = year - TM_YEAR_BASE;

	return mktime(&tm);
}

long
getoffset (const char *zone, long long t)
{
	struct tm *tp = localtime((time_t*)&t);
	if (!tp) {
		return -1;
	}

# if defined __USE_MISC
	return tp->tm_gmtoff;
# else
	return tp->__tm_gmtoff;
#endif
}

static time_t
yeartot(const intmax_t y)
{
	register intmax_t	myy, seconds, years;
	register time_t		t;

	myy = EPOCH_YEAR;
	t = 0;
	while (myy < y) {
		if (SECSPER400YEARS_FITS && 400 <= y - myy) {
			intmax_t diff400 = (y - myy) / 400;
			if (INTMAX_MAX / SECSPER400YEARS < diff400)
				return absolute_max_time;
			seconds = diff400 * SECSPER400YEARS;
			years = diff400 * 400;
		} else {
			seconds = isleap(myy) ? SECSPERLYEAR : SECSPERNYEAR;
			years = 1;
		}
		myy += years;
		if (t > absolute_max_time - seconds)
			return absolute_max_time;
		t += seconds;
	}
	while (y < myy) {
		if (SECSPER400YEARS_FITS && y + 400 <= myy && myy < 0) {
			intmax_t diff400 = (myy - y) / 400;
			if (INTMAX_MAX / SECSPER400YEARS < diff400)
				return absolute_min_time;
			seconds = diff400 * SECSPER400YEARS;
			years = diff400 * 400;
		} else {
			seconds = isleap(myy - 1) ? SECSPERLYEAR : SECSPERNYEAR;
			years = 1;
		}
		myy -= years;
		if (t < absolute_min_time + seconds)
			return absolute_min_time;
		t -= seconds;
	}
	return t;
}

static time_t
hunt(time_t lot, time_t hit)
{
	time_t			t;
	struct tm		lotm;
	register struct tm *	lotmp;
	struct tm		tm;
	register struct tm *	tmp;
	char			loab[MAX_STRING_LENGTH];

	lotmp = localtime(&lot);
	if (lotmp != NULL) {
		lotm = *lotmp;
		strncpy(loab, abbr(&lotm), (sizeof loab) - 1);
	}
	for ( ; ; ) {
		time_t diff = hit - lot;
		if (diff < 2)
			break;
		t = lot;
		t += diff / 2;
		if (t <= lot)
			++t;
		else if (t >= hit)
			--t;
		tmp = localtime(&t);
		if (tmp != NULL)
			tm = *tmp;
		if ((lotmp == NULL || tmp == NULL) ? (lotmp == tmp) :
		        (delta(&tm, &lotm) == (t - lot) &&
		         tm.tm_isdst == lotm.tm_isdst &&
		         strcmp(abbr(&tm), loab) == 0)) {
			lot = t;
			lotm = tm;
			lotmp = tmp;
		} else	hit = t;
	}
	return hit;
}

/*
** Thanks to Paul Eggert for logic used in delta.
*/

static intmax_t
delta(struct tm * newp, struct tm *oldp)
{
	register intmax_t	result;
	register int		tmy;

	if (newp->tm_year < oldp->tm_year)
		return -delta(oldp, newp);
	result = 0;
	for (tmy = oldp->tm_year; tmy < newp->tm_year; ++tmy)
		result += DAYSPERNYEAR + isleap_sum(tmy, TM_YEAR_BASE);
	result += newp->tm_yday - oldp->tm_yday;
	result *= HOURSPERDAY;
	result += newp->tm_hour - oldp->tm_hour;
	result *= MINSPERHOUR;
	result += newp->tm_min - oldp->tm_min;
	result *= SECSPERMIN;
	result += newp->tm_sec - oldp->tm_sec;
	return result;
}

static char *
abbr(struct tm *tmp)
{
	register char *	result;
	static char	nada;

	if (tmp->tm_isdst != 0 && tmp->tm_isdst != 1)
		return &nada;
	result = tzname[tmp->tm_isdst];
	return (result == NULL) ? &nada : result;
}
