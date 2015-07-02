#include "dockconstants.h"

DockConstants::DockConstants(QObject *parent) :
    QObject(parent)
{
}

DockConstants * DockConstants::dockConstants = NULL;
DockConstants * DockConstants::getInstants()
{
    if (dockConstants == NULL)
        dockConstants = new DockConstants();

    return dockConstants;
}

DockConstants::DockMode DockConstants::getDockMode()
{
  return m_currentMode;
}

void DockConstants::setDockMode(DockMode value)
{
    DockMode tmpValue = m_currentMode;
    m_currentMode = value;

    emit dockModeChanged(value, tmpValue);
}

int DockConstants::getDockHeight()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 60;
    case DockConstants::EfficientMode:
        return 50;
    case DockConstants::ClassicMode:
        return 40;
    default:
        return 40;
    }
}

int DockConstants::getItemHeight()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 60;
    case DockConstants::EfficientMode:
        return 50;
    case DockConstants::ClassicMode:
        return 40;
    default:
        return 40;
    }
}

int DockConstants::getNormalItemWidth()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 60;
    case DockConstants::EfficientMode:
        return 60;
    case DockConstants::ClassicMode:
        return 40;
    default:
        return 40;
    }
}

int DockConstants::getActivedItemWidth()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 60;
    case DockConstants::EfficientMode:
        return 60;
    case DockConstants::ClassicMode:
        return 150;
    default:
        return 60;
    }
}

int DockConstants::getAppItemSpacing()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 10;
    case DockConstants::EfficientMode:
        return 15;
    case DockConstants::ClassicMode:
        return 8;
    default:
        return 8;
    }
}

int DockConstants::getAppIconSize()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 48;
    case DockConstants::EfficientMode:
        return 42;
    case DockConstants::ClassicMode:
        return 24;
    default:
        return 32;
    }
}

int DockConstants::getAppletsItemHeight()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 60;
    case DockConstants::EfficientMode:
        return 50;
    case DockConstants::ClassicMode:
        return 40;
    default:
        return 40;
    }
}

int DockConstants::getAppletsItemWidth()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 60;
    case DockConstants::EfficientMode:
        return 50;
    case DockConstants::ClassicMode:
        return 50;
    default:
        return 50;
    }
}

int DockConstants::getAppletsItemSpacing()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 10;
    case DockConstants::EfficientMode:
        return 10;
    case DockConstants::ClassicMode:
        return 10;
    default:
        return 10;
    }
}

int DockConstants::getAppletsIconSize()
{
    switch (m_currentMode)
    {
    case DockConstants::FashionMode:
        return 48;
    case DockConstants::EfficientMode:
        return 24;
    case DockConstants::ClassicMode:
        return 24;
    default:
        return 24;
    }
}

