#include "powertrayloader.h"
#include "powertraywidget.h"

#define PowerItemKey "system-tray-power"
#define PowerService "com.deepin.daemon.Power"

PowerTrayLoader::PowerTrayLoader(QObject *parent)
    : AbstractTrayLoader(PowerService, parent),
      m_powerInter(new DBusPower(this))
{
}

void PowerTrayLoader::load()
{
    if (!m_powerInter->batteryState().isEmpty()) {
        emit systemTrayAdded(PowerItemKey, new PowerTrayWidget);
    }
}
