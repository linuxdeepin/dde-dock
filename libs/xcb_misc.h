#ifndef XCB_MISC_H
#define XCB_MISC_H

#include <xcb/xcb_ewmh.h>

class XcbMisc
{

public:
    enum Orientation {
        OrientationLeft,
        OrientationRight,
        OrientationTop,
        OrientationBottom
    };

    virtual ~XcbMisc();

    static XcbMisc * instance();

    void set_strut_partial(int winId, Orientation orientation, uint strut, uint start, uint end);

private:
    XcbMisc();

    xcb_ewmh_connection_t m_ewmh_connection;
};

#endif // XCB_MISC_H
