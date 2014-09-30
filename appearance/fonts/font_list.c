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

#include <fontconfig/fontconfig.h>
#include <fontconfig/fcfreetype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "font_list.h"

void font_info_free(FontInfo *info);

FontInfo *
get_font_info_list (int *num)
{
	FcInit();
	*num = -1;
	FcPattern *pat = FcPatternCreate();
	if (!pat) {
		fprintf(stderr, "Create FcPattern Failed\n");
		return NULL;
	}

	FcObjectSet *os = FcObjectSetBuild(
	                      FC_FAMILY,
	                      FC_FAMILYLANG,
	                      FC_STYLE,
	                      FC_FILE,
	                      FC_LANG,
	                      FC_SPACING,
	                      NULL);
	if (!os) {
		fprintf(stderr, "Build FcObjectSet Failed\n");
		FcPatternDestroy(pat);
		return NULL;
	}

	FcFontSet *fs = FcFontList(0, pat, os);
	FcObjectSetDestroy(os);
	FcPatternDestroy(pat);
	if (!fs) {
		fprintf(stderr, "List Font Failed\n");
		return NULL;
	}

	int i;
	int cnt = 0;
	FontInfo *list = NULL;
	for (i = 0; i < fs->nfont; i++) {
		FontInfo *info = calloc(1, sizeof(FontInfo));
		if (!info) {
			fprintf(stderr, "Alloc memory failed");
			continue;
		}

		info->family = (char*)FcPatternFormat(fs->fonts[i],
		                                      (FcChar8*)"%{family}");
		info->familylang = (char*)FcPatternFormat(fs->fonts[i],
		                   (FcChar8*)"%{familylang}");
		info->style = (char*)FcPatternFormat(fs->fonts[i],
		                                     (FcChar8*)"%{style}");
		info->filename = (char*)FcPatternFormat(fs->fonts[i],
		                                        (FcChar8*)"%{file}");
		info->lang = (char*)FcPatternFormat(fs->fonts[i],
		                                    (FcChar8*)"%{lang}");
		info->spacing = (char*)FcPatternFormat(fs->fonts[i],
		                                       (FcChar8*)"%{spacing}");
		if (!info->family || !info->familylang ||
		        !info->style || !info->filename ||
		        !info->lang || !info->spacing) {
			font_info_free(info);
			continue;
		}

		FontInfo *tmp = malloc((cnt+1) * sizeof(FontInfo));
		if (!tmp) {
			fprintf(stderr, "Alloc memory failed\n");
			font_info_free(info);
			continue;
		}

		memcpy(tmp+cnt, info, sizeof(FontInfo));
		free(info);
		if (cnt != 0 ) {
			memcpy(tmp, list, cnt * sizeof(FontInfo));
			free(list);
			list = NULL;
		}

		list = tmp;
		tmp = NULL;

		cnt++;
	}
	FcFontSetDestroy(fs);
	FcFini();

	*num = cnt;

	return list;
}

void
font_info_list_free(FontInfo *list, int num)
{
	if (!list) {
		return;
	}

	int i;
	for (i = 0; i < num; i++) {
		font_info_free(list+i);
	}

	free(list);
}

void
font_info_free(FontInfo *info)
{
	if (info == NULL) {
		return;
	}

	free(info->family);
	free(info->familylang);
	free(info->style);
	free(info->lang);
	free(info->spacing);
	free(info->filename);
}
