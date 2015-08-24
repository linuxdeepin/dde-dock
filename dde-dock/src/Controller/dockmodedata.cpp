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
    m_dhsm->UpdateState();
    QTimer::singleShot(100, m_dhsm, SLOT(UpdateState()));
}

int DockModeData::getDockHeight()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::PANEL_FASHION_HEIGHT;
    case Dock::EfficientMode:
        return Dock::PANEL_EFFICIENT_HEIGHT;
    case Dock::ClassicMode:
        return Dock::PANEL_CLASSIC_HEIGHT;
    default:
        return Dock::PANEL_FASHION_HEIGHT;
    }
}

int DockModeData::getItemHeight()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APP_ITEM_FASHION_HEIGHT;
    case Dock::EfficientMode:
        return Dock::APP_ITEM_EFFICIENT_HEIGHT;
    case Dock::ClassicMode:
        return Dock::APP_ITEM_CLASSIC_HEIGHT;
    default:
        return Dock::APP_ITEM_FASHION_HEIGHT;
    }
}

int DockModeData::getNormalItemWidth()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APP_ITEM_FASHION_NORMAL_WIDTH;
    case Dock::EfficientMode:
        return Dock::APP_ITEM_EFFICIENT_NORMAL_WIDTH;
    case Dock::ClassicMode:
        return Dock::APP_ITEM_CLASSIC_NORMAL_WIDTH;
    default:
        return Dock::APP_ITEM_FASHION_NORMAL_WIDTH;
    }
}

int DockModeData::getActivedItemWidth()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APP_ITEM_FASHION_ACTIVE_WIDTH;
    case Dock::EfficientMode:
        return Dock::APP_ITEM_EFFICIENT_ACTIVE_WIDTH;
    case Dock::ClassicMode:
        return Dock::APP_ITEM_CLASSIC_ACTIVE_WIDTH;
    default:
        return Dock::APP_ITEM_FASHION_ACTIVE_WIDTH;
    }
}

int DockModeData::getAppItemSpacing()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APP_ITEM_FASHION_SPACING;
    case Dock::EfficientMode:
        return Dock::APP_ITEM_EFFICIENT_SPACING;
    case Dock::ClassicMode:
        return Dock::APP_ITEM_CLASSIC_SPACING;
    default:
        return Dock::APP_ITEM_FASHION_SPACING;
    }
}

int DockModeData::getAppIconSize()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APP_ITEM_FASHION_ICON_SIZE;
    case Dock::EfficientMode:
        return Dock::APP_ITEM_EFFICIENT_ICON_SIZE;
    case Dock::ClassicMode:
        return Dock::APP_ITEM_CLASSIC_ICON_SIZE;
    default:
        return Dock::APP_ITEM_FASHION_ICON_SIZE;
    }
}

int DockModeData::getAppletsItemHeight()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APPLET_FASHION_ITEM_HEIGHT;
    case Dock::EfficientMode:
        return Dock::APPLET_EFFICIENT_ITEM_HEIGHT;
    case Dock::ClassicMode:
        return Dock::APPLET_CLASSIC_ITEM_HEIGHT;
    default:
        return Dock::APPLET_FASHION_ITEM_HEIGHT;
    }
}

int DockModeData::getAppletsItemWidth()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APPLET_FASHION_ITEM_WIDTH;
    case Dock::EfficientMode:
        return Dock::APPLET_EFFICIENT_ITEM_WIDTH;
    case Dock::ClassicMode:
        return Dock::APPLET_CLASSIC_ITEM_WIDTH;
    default:
        return Dock::APPLET_FASHION_ITEM_WIDTH;
    }
}

int DockModeData::getAppletsItemSpacing()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APPLET_FASHION_ITEM_SPACING;
    case Dock::EfficientMode:
        return Dock::APPLET_EFFICIENT_ITEM_SPACING;
    case Dock::ClassicMode:
        return Dock::APPLET_CLASSIC_ITEM_SPACING;
    default:
        return Dock::APPLET_FASHION_ITEM_SPACING;
    }
}

int DockModeData::getAppletsIconSize()
{
    switch (m_currentMode)
    {
    case Dock::FashionMode:
        return Dock::APPLET_FASHION_ICON_SIZE;
    case Dock::EfficientMode:
        return Dock::APPLET_EFFICIENT_ICON_SIZE;
    case Dock::ClassicMode:
        return Dock::APPLET_CLASSIC_ICON_SIZE;
    default:
        return Dock::APPLET_FASHION_ICON_SIZE;
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
