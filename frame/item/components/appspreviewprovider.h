/*
 * Copyright (C) 2011 ~ 2018 uniontech Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
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
#ifndef APPSPREVIEWPROVIDER_H
#define APPSPREVIEWPROVIDER_H

#include "previewcontainer.h"

static PreviewContainer *PreviewWindow(const WindowInfoMap &infos, const WindowList &allowClose, const Dock::Position dockPos)
{
    static PreviewContainer *preview;
    if (!preview) {
        preview = new PreviewContainer;
    }

    preview->disconnect();
    preview->setWindowInfos(infos, allowClose);
    preview->updateSnapshots();
    preview->updateLayoutDirection(dockPos);

    return preview;
}

#endif /* APPSPREVIEWPROVIDER_H */
