#include <glib.h>
#include "gsd-power-manager.h"
#include "power-force-idle.h"
#include "libgnome-desktop/gnome-idle-monitor.h"

int locked = 0;

int start_dim()
{
    locked = 1;
    return 0;
}

int stop_dim()
{
    locked = 0;
    return 0;
}
