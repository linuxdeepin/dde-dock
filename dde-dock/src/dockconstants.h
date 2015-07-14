#ifndef DOCKCONSTANTS_H
#define DOCKCONSTANTS_H

namespace Dock {

enum DockMode {
    FashionMode = 0,
    EfficientMode = 1,
    ClassicMode = 2
};

enum HideMode {
    KeepShowing = 0,
    KeepHidden = 1,
    SmartHide = 3
};

const int APP_PREVIEW_WIDTH = 160;
const int APP_PREVIEW_HEIGHT = 100;
const int APP_PREVIEW_MARGIN = 18 ;

const int APPLET_FASHION_ITEM_HEIGHT = 60;
const int APPLET_FASHION_ITEM_WIDTH = 60;
const int APPLET_FASHION_ITEM_SPACING = 10;
const int APPLET_FASHION_ICON_SIZE = 48;

const int APPLET_EFFICIENT_ITEM_HEIGHT = 50;
const int APPLET_EFFICIENT_ITEM_WIDTH = 50;
const int APPLET_EFFICIENT_ITEM_SPACING = 10;
const int APPLET_EFFICIENT_ICON_SIZE = 24;

const int APPLET_CLASSIC_ITEM_HEIGHT = 40;
const int APPLET_CLASSIC_ITEM_WIDTH = 50;
const int APPLET_CLASSIC_ITEM_SPACING = 10;
const int APPLET_CLASSIC_ICON_SIZE = 24;
}

#endif // DOCKCONSTANTS_H
