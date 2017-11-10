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

#include <fontconfig/fontconfig.h>
#include <fontconfig/fcfreetype.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "font_list.h"

static int append_font_info(FcInfo** list, FcPattern* pat, int idx);
static void free_font_info(FcInfo *info);

int
fc_cache_update ()
{
    static int first = 0;
    if (first == 0) {
        FcInit();
        first = 1;
        return 1;
    }

    return !FcConfigUptoDate (NULL) && FcInitReinitialize ();
}

FcInfo *
list_font_info (int *num)
{
     *num = -1;
     FcPattern *pat = FcPatternCreate();
     if (!pat) {
          fprintf(stderr, "Create FcPattern Failed\n");
          return NULL;
     }

     FcObjectSet *os = FcObjectSetBuild(
          FC_FAMILY,
          FC_FAMILYLANG,
          /* FC_FULLNAME, */
          /* FC_FULLNAMELANG, */
          FC_STYLE,
          /* FC_FILE, */
          FC_LANG,
          FC_SPACING,
          FC_CHARSET,
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
          char *charset = (char*)FcPatternFormat(fs->fonts[i], (FcChar8*)"%{charset}");
          if (charset == NULL || strlen(charset) == 0) {
              free(charset);
              continue;
          }
          free(charset);

          if (append_font_info(&list, fs->fonts[i], cnt) == -1) {
               continue;
          }
          cnt++;
     }
     FcFontSetDestroy(fs);
     //FcFini(); // SIGABRT: FcCacheFini 'assert fcCacheChains[i] == NULL failed'

     *num = cnt;

     return list;
}

void
free_font_info_list(FcInfo *list, int num)
{
     if (!list) {
          return;
     }

     int i;
     for (i = 0; i < num; i++) {
          free_font_info(list+i);
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

static int
append_font_info(FcInfo** list, FcPattern* pat, int idx)
{
     // realloc equals to malloc if the first argument is NULL.
     // realloc will auto free old memory chunk.
     FcInfo* tmp = (FcInfo*)realloc(*list, (idx+1)*sizeof(FcInfo));
     if (!tmp) {
          fprintf(stderr, "Alloc memory at append %d font info failed\n", idx+1);
          return -1;
     }

     *list = tmp;

     tmp[idx].family = (char*)FcPatternFormat(pat, (FcChar8*)"%{family}");
     tmp[idx].familylang = (char*)FcPatternFormat(pat, (FcChar8*)"%{familylang}");
     /* tmp[idx].fullname = (char*)FcPatternFormat(pat, (FcChar8*)"%{fullname}"); */
     /* tmp[idx].fullnamelang = (char*)FcPatternFormat(pat, (FcChar8*)"%{fullnamelang}"); */
     tmp[idx].style = (char*)FcPatternFormat(pat, (FcChar8*)"%{style}");
     /* tmp[idx].filename = (char*)FcPatternFormat(pat, (FcChar8*)"%{file}"); */
     tmp[idx].lang = (char*)FcPatternFormat(pat, (FcChar8*)"%{lang}");
     tmp[idx].spacing = (char*)FcPatternFormat(pat, (FcChar8*)"%{spacing}");

     return 0;
}

static void
free_font_info(FcInfo *info)
{
     if (info == NULL) {
          return;
     }

     free(info->family);
     free(info->familylang);
     /* free(info->fullname); */
     /* free(info->fullnamelang); */
     free(info->style);
     free(info->lang);
     free(info->spacing);
     /* free(info->filename); */
}
