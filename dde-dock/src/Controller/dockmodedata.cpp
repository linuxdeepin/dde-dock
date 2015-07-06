#include "dockmodedata.h"

DockModeData::DockModeData(QObject *parent) :
    QObject(parent)
{
}

DockModeData * DockModeData::dockModeData = NULL;
DockModeData * DockModeData::getInstants()
{
    if (dockModeData == NULL)
        dockModeData = new DockModeData();

    return dockModeData;
}

Dock::DockMode DockModeData::getDockMode()
{
  return m_currentMode;
}

void DockModeData::setDockMode(Dock::DockMode value)
{
    Dock::DockMode tmpValue = m_currentMode;
    m_currentMode = value;

    emit dockModeChanged(value, tmpValue);
}

int DockModeData::getDockHeight()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 60;
    case Dock::EfficientMode:
        return 50;
    case Dock::ClassicMode:
        return 40;
    default:
        return 40;
    }
}

int DockModeData::getItemHeight()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 60;
    case Dock::EfficientMode:
        return 50;
    case Dock::ClassicMode:
        return 40;
    default:
        return 40;
    }
}

int DockModeData::getNormalItemWidth()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 60;
    case Dock::EfficientMode:
        return 60;
    case Dock::ClassicMode:
        return 40;
    default:
        return 40;
    }
}

int DockModeData::getActivedItemWidth()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 60;
    case Dock::EfficientMode:
        return 60;
    case Dock::ClassicMode:
        return 150;
    default:
        return 60;
    }
}

int DockModeData::getAppItemSpacing()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 10;
    case Dock::EfficientMode:
        return 15;
    case Dock::ClassicMode:
        return 8;
    default:
        return 8;
    }
}

int DockModeData::getAppIconSize()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 48;
    case Dock::EfficientMode:
        return 48;
    case Dock::ClassicMode:
        return 32;
    default:
        return 32;
    }
}

int DockModeData::getAppletsItemHeight()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 60;
    case Dock::EfficientMode:
        return 50;
    case Dock::ClassicMode:
        return 40;
    default:
        return 40;
    }
}

int DockModeData::getAppletsItemWidth()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 60;
    case Dock::EfficientMode:
        return 50;
    case Dock::ClassicMode:
        return 50;
    default:
        return 50;
    }
}

int DockModeData::getAppletsItemSpacing()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 10;
    case Dock::EfficientMode:
        return 10;
    case Dock::ClassicMode:
        return 10;
    default:
        return 10;
    }
}

int DockModeData::getAppletsIconSize()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return 48;
    case Dock::EfficientMode:
        return 24;
    case Dock::ClassicMode:
        return 24;
    default:
        return 24;
    }
}

