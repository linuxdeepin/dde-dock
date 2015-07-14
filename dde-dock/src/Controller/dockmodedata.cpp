#include "dockmodedata.h"

DockModeData::DockModeData(QObject *parent) :
    QObject(parent)
{
    initDDS();
}

DockModeData * DockModeData::dockModeData = NULL;
DockModeData * DockModeData::instance()
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
    m_dds->SetDisplayMode(value);
}

Dock::HideMode DockModeData::getHideMode()
{
    return m_hideMode;
}

void DockModeData::setHideMode(Dock::HideMode value)
{
    m_dds->SetHideMode(value);
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

void DockModeData::slotDockModeChanged(int mode)
{
    Dock::DockMode tmpMode = Dock::DockMode(mode);
    Dock::DockMode oldmode = m_currentMode;
    m_currentMode = tmpMode;

    emit dockModeChanged(tmpMode,oldmode);
}

void DockModeData::slotHideModeChanged(int mode)
{
    Dock::HideMode tmpMode = Dock::HideMode(mode);
    Dock::HideMode oldMode = m_hideMode;
    m_hideMode = tmpMode;

    emit hideModeChanged(tmpMode,oldMode);
}

void DockModeData::initDDS()
{
    m_dds = new DBusDockSetting(this);
    connect(m_dds,&DBusDockSetting::DisplayModeChanged,this,&DockModeData::slotDockModeChanged);
    connect(m_dds,&DBusDockSetting::HideModeChanged,this,&DockModeData::slotHideModeChanged);

    m_currentMode = Dock::DockMode(m_dds->GetDisplayMode().value());
    m_hideMode = Dock::HideMode(m_dds->GetHideMode().value());

    emit dockModeChanged(m_currentMode,m_currentMode);
    emit hideModeChanged(m_hideMode,m_hideMode);
}
