#ifndef CONSTANTS_H
#define CONSTANTS_H

namespace Dock {

enum DisplayMode
{
    Fashion     = 0,
    Efficient   = 1,
    // deprecreated
//    Classic     = 2,
};

enum HideMode
{
    KeepShowing     = 0,
    KeepHidden      = 1,
    SmartHide       = 3,
};

enum Position
{
    Top         = 0,
    Right       = 1,
    Bottom      = 2,
    Left        = 3,
};

enum HideState
{
    Unknown     = 0,
    Show        = 1,
    Hide        = 2,
};

}

#endif // CONSTANTS_H
