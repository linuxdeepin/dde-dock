#include "dockmodedata.h"

DockModeData::DockModeData(QObject *parent) :
    QObject(parent)
{
    initDDS();
}

DockModeData * DockModeData::m_dockModeData = NULL;
DockModeData * DockModeData::instance()
{
    if (m_dockModeData == NULL)
        m_dockModeData = new DockModeData();

    return m_dockModeData;
}

Dock::DockMode DockModeData::getDockMode()
{
  return m_currentMode;
}

void DockModeData::setDockMode(Dock::DockMode value)
{
    m_dockSetting->SetDisplayMode(value);
}

Dock::HideMode DockModeData::getHideMode()
{
    return m_hideMode;
}

void DockModeData::setHideMode(Dock::HideMode value)
{
    m_dockSetting->SetHideMode(value);
    m_hideStateManager->UpdateState();
    QTimer::singleShot(100, m_hideStateManager, SLOT(UpdateState()));
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

void DockModeData::onDockModeChanged(int mode)
{
    Dock::DockMode tmpMode = Dock::DockMode(mode);
    Dock::DockMode oldmode = m_currentMode;
    m_currentMode = tmpMode;

    emit dockModeChanged(tmpMode,oldmode);
}

void DockModeData::onHideModeChanged(int mode)
{
    Dock::HideMode tmpMode = Dock::HideMode(mode);
    Dock::HideMode oldMode = m_hideMode;
    m_hideMode = tmpMode;

    emit hideModeChanged(tmpMode,oldMode);
}

void DockModeData::initDDS()
{
    m_dockSetting = new DBusDockSetting(this);
    connect(m_dockSetting,&DBusDockSetting::DisplayModeChanged,this,&DockModeData::onDockModeChanged);
    connect(m_dockSetting,&DBusDockSetting::HideModeChanged,this,&DockModeData::onHideModeChanged);

    m_currentMode = Dock::DockMode(m_dockSetting->GetDisplayMode().value());
    m_hideMode = Dock::HideMode(m_dockSetting->GetHideMode().value());

    emit dockModeChanged(m_currentMode,m_currentMode);
    emit hideModeChanged(m_hideMode,m_hideMode);
}
