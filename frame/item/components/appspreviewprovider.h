
#ifndef APPSPREVIEWPROVIDER_H
#define APPSPREVIEWPROVIDER_H

#include "previewcontainer.h"

static PreviewContainer *PreviewWindow(const WindowInfoMap &infos, const Dock::Position dockPos)
{
    static PreviewContainer *preview;
    if (!preview) {
        preview = new PreviewContainer;
    }

    preview->disconnect();
    preview->setWindowInfos(infos);
    preview->updateSnapshots();
    preview->updateLayoutDirection(dockPos);

    return preview;
}

#endif /* APPSPREVIEWPROVIDER_H */
