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

#ifndef __FONT_LIST_H__
#define __FONT_LIST_H__

typedef struct _FcInfo {
	char *family;
	char *familylang;
	/* char *fullname; */
	/* char *fullnamelang; */
	char *style;
	char *lang;
	char *spacing;
	/* char *filename; */
} FcInfo;

int fc_cache_update ();
FcInfo *list_font_info (int *num);
void free_font_info_list(FcInfo *list, int num);

char* font_match(char* family);

#endif
