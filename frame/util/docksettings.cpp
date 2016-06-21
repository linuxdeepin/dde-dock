#include "docksettings.h"

#include <QDebug>

DockSettings::DockSettings(QObject *parent)
    : QObject(parent),
      m_dockInter(new DBusDock(this)),
      m_itemController(DockItemController::instance(this))
{
    m_position = Dock::Position(m_dockInter->position());
}

Position DockSettings::position() const
{
    return m_position;
}

const QSize DockSettings::mainWindowSize() const
{
    return m_mainWindowSize;
}

void DockSettings::updateGeometry()
{

}
