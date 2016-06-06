#include "docksettings.h"

DockSettings::DockSettings(QObject *parent)
    : QObject(parent)
{

}

DockSettings::DockSide DockSettings::side() const
{
    return Top;
}
