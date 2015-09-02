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

static void font_info_free(FcInfo *info);

FcInfo *
get_font_info_list (int *num)
{
	/* FcInit(); */
	*num = -1;
	FcPattern *pat = FcPatternCreate();
	if (!pat) {
		fprintf(stderr, "Create FcPattern Failed\n");
		return NULL;
	}

	FcObjectSet *os = FcObjectSetBuild(
	                      FC_FAMILY,
	                      FC_FAMILYLANG,
	                      FC_FULLNAME,
	                      FC_FULLNAMELANG,
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
	FcInfo *list = NULL;
	for (i = 0; i < fs->nfont; i++) {
		FcInfo *info = calloc(1, sizeof(FcInfo));
		if (!info) {
			fprintf(stderr, "Alloc memory failed");
			continue;
		}

		info->family = (char*)FcPatternFormat(fs->fonts[i],
		                                      (FcChar8*)"%{family}");
		info->familylang = (char*)FcPatternFormat(fs->fonts[i],
		                   (FcChar8*)"%{familylang}");
		info->fullname = (char*)FcPatternFormat(fs->fonts[i],
		                                        (FcChar8*)"%{fullname}");
		info->fullnamelang = (char*)FcPatternFormat(fs->fonts[i],
		                     (FcChar8*)"%{fullnamelang}");
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

		FcInfo *tmp = malloc((cnt+1) * sizeof(FcInfo));
		if (!tmp) {
			fprintf(stderr, "Alloc memory failed\n");
			font_info_free(info);
			continue;
		}

		memcpy(tmp+cnt, info, sizeof(FcInfo));
		free(info);
		if (cnt != 0 ) {
			memcpy(tmp, list, cnt * sizeof(FcInfo));
			free(list);
			list = NULL;
		}

		list = tmp;
		tmp = NULL;

		cnt++;
	}
	FcFontSetDestroy(fs);
	//FcFini(); // SIGABRT: FcCacheFini 'assert fcCacheChains[i] == NULL failed'

	*num = cnt;

	return list;
}

void
font_info_list_free(FcInfo *list, int num)
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

char*
font_match(char* family)
{
	// configure the search pattern
	FcPattern* pat = FcNameParse((FcChar8*)family);
	if (!pat) {
		return NULL;
	}

	FcConfigSubstitute(NULL, pat, FcMatchPattern);
	FcDefaultSubstitute(pat);

	FcResult result;
	FcPattern* match = FcFontMatch(NULL, pat, &result);
	FcPatternDestroy(pat);
	if (!match) {
		return NULL;
	}

	FcFontSet* fs = FcFontSetCreate();
	if (!fs) {
		FcPatternDestroy(match);
		return NULL;
	}

	FcFontSetAdd(fs, match);
	FcPattern* font = FcPatternFilter(fs->fonts[0], NULL);
	FcChar8* ret = FcPatternFormat(font, (const FcChar8*)"%{=fcmatch}\n");

	FcPatternDestroy(font);
	FcFontSetDestroy(fs);
	FcPatternDestroy(match);
	//FcFini(); // SIGABRT: FcCacheFini 'assert fcCacheChains[i] == NULL failed'

	if (!ret) {
		return NULL;
	}

	return (char*)ret;
}

static void
font_info_free(FcInfo *info)
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



