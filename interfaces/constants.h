#ifndef CONSTANTS_H
#define CONSTANTS_H

#include <QtCore>

namespace Dock {

#define PROP_DISPLAY_MODE   "DisplayMode"
enum DisplayMode
{
    Fashion     = 0,
    Efficient   = 1,
    // deprecreated
//    Classic     = 2,
};

#define PROP_HIDE_MODE      "HideMode"
enum HideMode
{
    KeepShowing     = 0,
    KeepHidden      = 1,
    SmartHide       = 3,
};

#define PROP_POSITION       "Position"
enum Position
{
    Top         = 0,
    Right       = 1,
    Bottom      = 2,
    Left        = 3,
};

#define PROP_HIDE_STATE     "HideState"
enum HideState
{
    Unknown     = 0,
    Show        = 1,
    Hide        = 2,
};

}

Q_DECLARE_METATYPE(Dock::DisplayMode)
Q_DECLARE_METATYPE(Dock::Position)

#endif // CONSTANTS_H
