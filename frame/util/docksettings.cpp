#include "docksettings.h"

#include <QDebug>

#define ICON_SIZE_LARGE     48
#define ICON_SIZE_MEDIUM    36
#define ICON_SIZE_SMALL     24

DockSettings::DockSettings(QObject *parent)
    : QObject(parent),

      m_settingsMenu(this),
      m_fashionModeAct(tr("Fashion Mode"), this),
      m_efficientModeAct(tr("Efficient Mode"), this),
      m_topPosAct(tr("Top"), this),
      m_bottomPosAct(tr("Bottom"), this),
      m_leftPosAct(tr("Left"), this),
      m_rightPosAct(tr("Right"), this),
      m_largeSizeAct(tr("Large"), this),
      m_mediumSizeAct(tr("Medium"), this),
      m_smallSizeAct(tr("Small"), this),
      m_keepShownAct(tr("Keep Shown"), this),
      m_keepHiddenAct(tr("Keep Hidden"), this),
      m_smartHideAct(tr("Smart Hide"), this),

      m_displayInter(new DBusDisplay(this)),
      m_dockInter(new DBusDock(this)),
      m_itemController(DockItemController::instance(this))
{
    m_position = Dock::Position(m_dockInter->position());
    m_displayMode = Dock::DisplayMode(m_dockInter->displayMode());

    m_mainWindowSize.setWidth(m_displayInter->primaryRect().width);
    m_mainWindowSize.setHeight(60);

    m_fashionModeAct.setCheckable(true);
    m_efficientModeAct.setCheckable(true);
    m_topPosAct.setCheckable(true);
    m_bottomPosAct.setCheckable(true);
    m_leftPosAct.setCheckable(true);
    m_rightPosAct.setCheckable(true);
    m_largeSizeAct.setCheckable(true);
    m_mediumSizeAct.setCheckable(true);
    m_smallSizeAct.setCheckable(true);
    m_keepShownAct.setCheckable(true);
    m_keepHiddenAct.setCheckable(true);
    m_smartHideAct.setCheckable(true);

    DMenu *modeSubMenu = new DMenu(&m_settingsMenu);
    modeSubMenu->addAction(&m_fashionModeAct);
    modeSubMenu->addAction(&m_efficientModeAct);
    DAction *modeSubMenuAct = new DAction(tr("Mode"), this);
    modeSubMenuAct->setMenu(modeSubMenu);

    DMenu *locationSubMenu = new DMenu(&m_settingsMenu);
    locationSubMenu->addAction(&m_topPosAct);
    locationSubMenu->addAction(&m_bottomPosAct);
    locationSubMenu->addAction(&m_leftPosAct);
    locationSubMenu->addAction(&m_rightPosAct);
    DAction *locationSubMenuAct = new DAction(tr("Location"), this);
    locationSubMenuAct->setMenu(locationSubMenu);

    DMenu *sizeSubMenu = new DMenu(&m_settingsMenu);
    sizeSubMenu->addAction(&m_largeSizeAct);
    sizeSubMenu->addAction(&m_mediumSizeAct);
    sizeSubMenu->addAction(&m_smallSizeAct);
    DAction *sizeSubMenuAct = new DAction(tr("Size"), this);
    sizeSubMenuAct->setMenu(sizeSubMenu);

    DMenu *statusSubMenu = new DMenu(&m_settingsMenu);
    statusSubMenu->addAction(&m_keepShownAct);
    statusSubMenu->addAction(&m_keepHiddenAct);
    statusSubMenu->addAction(&m_smartHideAct);
    DAction *statusSubMenuAct = new DAction(tr("Status"), this);
    statusSubMenuAct->setMenu(statusSubMenu);

    m_settingsMenu.addAction(modeSubMenuAct);
    m_settingsMenu.addAction(locationSubMenuAct);
    m_settingsMenu.addAction(sizeSubMenuAct);
    m_settingsMenu.addAction(statusSubMenuAct);

    connect(&m_settingsMenu, &DMenu::triggered, this, &DockSettings::menuActionClicked);
    connect(m_dockInter, &DBusDock::PositionChanged, this, &DockSettings::positionChanged);
}

Position DockSettings::position() const
{
    return m_position;
}

const QSize DockSettings::windowSize() const
{
    return m_mainWindowSize;
}

void DockSettings::showDockSettingsMenu()
{
    m_fashionModeAct.setChecked(m_displayMode == Fashion);
    m_efficientModeAct.setChecked(m_displayMode == Efficient);
    m_topPosAct.setChecked(m_position == Top);
    m_bottomPosAct.setChecked(m_position == Bottom);
    m_leftPosAct.setChecked(m_position == Left);
    m_rightPosAct.setChecked(m_position == Right);
    m_largeSizeAct.setChecked(m_iconSize == ICON_SIZE_LARGE);
    m_mediumSizeAct.setChecked(m_iconSize == ICON_SIZE_MEDIUM);
    m_smallSizeAct.setChecked(m_iconSize == ICON_SIZE_SMALL);
    m_keepShownAct.setChecked(m_hideMode == KeepShowing);
    m_keepHiddenAct.setChecked(m_hideMode == KeepHidden);
    m_smartHideAct.setChecked(m_hideMode == SmartHide);

    m_settingsMenu.exec();
}

void DockSettings::updateGeometry()
{

}

void DockSettings::menuActionClicked(DAction *action)
{
    Q_ASSERT(action);

    if (action == &m_fashionModeAct)
        return m_dockInter->setDisplayMode(Fashion);
    if (action == &m_efficientModeAct)
        return m_dockInter->setDisplayMode(Efficient);

    if (action == &m_topPosAct)
        return m_dockInter->setPosition(Top);
    if (action == &m_bottomPosAct)
        return m_dockInter->setPosition(Bottom);
    if (action == &m_leftPosAct)
        return m_dockInter->setPosition(Left);
    if (action == &m_rightPosAct)
        return m_dockInter->setPosition(Right);

    if (action == &m_keepShownAct)
        return m_dockInter->setHideMode(KeepShowing);
    if (action == &m_keepHiddenAct)
        return m_dockInter->setHideMode(KeepHidden);
    if (action == &m_smartHideAct)
        return m_dockInter->setHideMode(SmartHide);
}

void DockSettings::positionChanged()
{
    m_position = Dock::Position(m_dockInter->position());

    const QRect primaryRect = m_displayInter->primaryRect();
    const int defaultHeight = 60;
    const int defaultWidth = 60;

    switch (m_position)
    {
    case Top:
    case Bottom:
        m_mainWindowSize.setHeight(defaultHeight);
        m_mainWindowSize.setWidth(primaryRect.width());
        break;

    case Left:
    case Right:
        m_mainWindowSize.setHeight(primaryRect.height());
        m_mainWindowSize.setWidth(defaultWidth);
        break;

    default:
        Q_ASSERT(false);
    }

    emit dataChanged();
}

void DockSettings::calculateWindowConfig()
{

}
