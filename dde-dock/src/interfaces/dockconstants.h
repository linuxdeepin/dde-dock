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

enum HideState {
    HideStateShowing = 0,
    HideStateHiding = 1,
    HideStateShown = 2,
    HideStateHidden = 3
};

////////////////  App  ////////////////////////////////
const int APP_PREVIEW_WIDTH = 200;
const int APP_PREVIEW_HEIGHT = 124;
const int APP_PREVIEW_MARGIN = 20 ;
const int APP_PREVIEW_CLOSEBUTTON_SIZE = 26;

const int APP_ITEM_FASHION_HEIGHT = 70;
const int APP_ITEM_FASHION_NORMAL_WIDTH = 48;
const int APP_ITEM_FASHION_ACTIVE_WIDTH = 48;
const int APP_ITEM_FASHION_SPACING = 3;
const int APP_ITEM_FASHION_ICON_SIZE = 48;

const int APP_ITEM_EFFICIENT_HEIGHT = 48;
const int APP_ITEM_EFFICIENT_NORMAL_WIDTH = 64;
const int APP_ITEM_EFFICIENT_ACTIVE_WIDTH = 64;
const int APP_ITEM_EFFICIENT_SPACING = 4;
const int APP_ITEM_EFFICIENT_ICON_SIZE = 32;

const int APP_ITEM_CLASSIC_HEIGHT = 36;
const int APP_ITEM_CLASSIC_NORMAL_WIDTH = 48;
const int APP_ITEM_CLASSIC_ACTIVE_WIDTH = 160;
const int APP_ITEM_CLASSIC_SPACING = 4;
const int APP_ITEM_CLASSIC_ICON_SIZE = 24;

////////////////  APpplet  ////////////////////////////
const int APPLET_FASHION_ITEM_HEIGHT = 48;
const int APPLET_FASHION_ITEM_WIDTH = 48;
const int APPLET_FASHION_ITEM_SPACING = 3;
const int APPLET_FASHION_ICON_SIZE = 48;

const int APPLET_EFFICIENT_ITEM_HEIGHT = 16;
const int APPLET_EFFICIENT_ITEM_WIDTH = 16;
const int APPLET_EFFICIENT_ITEM_SPACING = 10;
const int APPLET_EFFICIENT_ICON_SIZE = 16;

const int APPLET_CLASSIC_ITEM_HEIGHT = 16;
const int APPLET_CLASSIC_ITEM_WIDTH = 16;
const int APPLET_CLASSIC_ITEM_SPACING = 10;
const int APPLET_CLASSIC_ICON_SIZE = 16;

/////////////  Panel  ////////////////////////////////
const int PANEL_FASHION_HEIGHT = 70;
const int PANEL_EFFICIENT_HEIGHT = 48;
const int PANEL_CLASSIC_HEIGHT = 36;

}

#endif // DOCKCONSTANTS_H
