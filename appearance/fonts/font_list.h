/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __FONT_LIST_H__
#define __FONT_LIST_H__

typedef struct _FcInfo {
	char *family;
	char *familylang;
	char *fullname;
	char *fullnamelang;
	char *style;
	char *lang;
	char *spacing;
	char *filename;
} FcInfo;

int fc_cache_update ();
FcInfo *list_font_info (int *num);
void free_font_info_list(FcInfo *list, int num);

char* font_match(char* family);

#endif
