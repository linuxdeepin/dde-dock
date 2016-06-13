#include "docksettings.h"

DockSettings::DockSettings(QObject *parent)
    : QObject(parent)
{

}

DockSettings::DockSide DockSettings::side() const
{
    return Top;
}

const QSize DockSettings::mainWindowSize() const
{
    return m_mainWindowSize;
}

void DockSettings::updateGeometry()
{

}
