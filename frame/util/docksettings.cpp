#include "docksettings.h"
#include "panel/mainpanel.h"
#include "item/appitem.h"

#include <QDebug>
#include <QX11Info>

#define ICON_SIZE_LARGE     48
#define ICON_SIZE_MEDIUM    36
#define ICON_SIZE_SMALL     30

DockSettings::DockSettings(QWidget *parent)
    : QObject(parent),

      m_autoHide(true),

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
    resetFrontendGeometry();
    m_primaryRect = m_displayInter->primaryRect();
    m_position = Dock::Position(m_dockInter->position());
    m_displayMode = Dock::DisplayMode(m_dockInter->displayMode());
    m_hideMode = Dock::HideMode(m_dockInter->hideMode());
    m_hideState = Dock::HideState(m_dockInter->hideState());
    m_iconSize = m_dockInter->iconSize();
    AppItem::setIconBaseSize(m_iconSize);
    DockItem::setDockPosition(m_position);
    qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_position));
    DockItem::setDockDisplayMode(m_displayMode);
    qApp->setProperty(PROP_DISPLAY_MODE, QVariant::fromValue(m_displayMode));

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
    connect(m_dockInter, &DBusDock::IconSizeChanged, this, &DockSettings::iconSizeChanged);
    connect(m_dockInter, &DBusDock::DisplayModeChanged, this, &DockSettings::displayModeChanged);
    connect(m_dockInter, &DBusDock::HideModeChanged, this, &DockSettings::hideModeChanged);
    connect(m_dockInter, &DBusDock::HideStateChanged, this, &DockSettings::hideStateChanegd);
    connect(m_dockInter, &DBusDock::ServiceRestarted, this, &DockSettings::resetFrontendGeometry);

    connect(m_itemController, &DockItemController::itemInserted, this, &DockSettings::dockItemCountChanged, Qt::QueuedConnection);
    connect(m_itemController, &DockItemController::itemRemoved, this, &DockSettings::dockItemCountChanged, Qt::QueuedConnection);

    connect(m_displayInter, &DBusDisplay::PrimaryRectChanged, this, &DockSettings::primaryScreenChanged);

    calculateWindowConfig();
}

DisplayMode DockSettings::displayMode() const
{
    return m_displayMode;
}

HideMode DockSettings::hideMode() const
{
    return m_hideMode;
}

Position DockSettings::position() const
{
    return m_position;
}

int DockSettings::screenHeight() const
{
    return m_displayInter->screenHeight();
}

int DockSettings::screenWidth() const
{
    return m_displayInter->screenWidth();
}

bool DockSettings::autoHide() const
{
    return m_autoHide;
}

HideState DockSettings::hideState() const
{
    return m_hideState;
}

const QRect DockSettings::primaryRect() const
{
    return m_primaryRect;
}

const QSize DockSettings::windowSize() const
{
    return m_mainWindowSize;
}

void DockSettings::showDockSettingsMenu()
{
    m_autoHide = false;

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

    setAutoHide(true);
}

void DockSettings::updateGeometry()
{

}

void DockSettings::setAutoHide(const bool autoHide)
{
    if (m_autoHide == autoHide)
        return;

    m_autoHide = autoHide;
    emit autoHideChanged(m_autoHide);
}

void DockSettings::menuActionClicked(DAction *action)
{
    Q_ASSERT(action);

    qDebug() << "settings triggered: " << action->text();

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

    if (action == &m_largeSizeAct)
        return m_dockInter->setIconSize(ICON_SIZE_LARGE);
    if (action == &m_mediumSizeAct)
        return m_dockInter->setIconSize(ICON_SIZE_MEDIUM);
    if (action == &m_smallSizeAct)
        return m_dockInter->setIconSize(ICON_SIZE_SMALL);

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
    DockItem::setDockPosition(m_position);
    qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_position));

    calculateWindowConfig();

    emit dataChanged();
}

void DockSettings::iconSizeChanged()
{
    m_iconSize = m_dockInter->iconSize();
    AppItem::setIconBaseSize(m_iconSize);

    calculateWindowConfig();

    emit dataChanged();
}

void DockSettings::displayModeChanged()
{
    m_displayMode = Dock::DisplayMode(m_dockInter->displayMode());
    DockItem::setDockDisplayMode(m_displayMode);
    qApp->setProperty(PROP_DISPLAY_MODE, QVariant::fromValue(m_displayMode));

    calculateWindowConfig();

    emit dataChanged();
}

void DockSettings::hideModeChanged()
{
    m_hideMode = Dock::HideMode(m_dockInter->hideMode());

    emit windowHideModeChanged();
}

void DockSettings::hideStateChanegd()
{
    const Dock::HideState state = Dock::HideState(m_dockInter->hideState());

    if (state == Dock::Unknown)
        return;

    m_hideState = state;

    emit windowVisibleChanegd();
}

void DockSettings::dockItemCountChanged()
{
    if (m_displayMode == Dock::Efficient)
        return;

    calculateWindowConfig();

    emit windowGeometryChanged();
}

void DockSettings::primaryScreenChanged()
{
    m_primaryRect = m_displayInter->primaryRect();

    calculateWindowConfig();

    emit dataChanged();
}

void DockSettings::resetFrontendGeometry()
{
    const QSize size = m_mainWindowSize;
    const QRect primaryRect = m_primaryRect;
    const int offsetX = (primaryRect.width() - size.width()) / 2;
    const int offsetY = (primaryRect.height() - size.height()) / 2;

    QPoint p(0, 0);
    switch (m_position)
    {
    case Top:
        p = QPoint(primaryRect.topLeft().x() + offsetX, 0);               break;
    case Left:
        p = QPoint(primaryRect.topLeft().x(), offsetY);                   break;
    case Right:
        p = QPoint(primaryRect.right() - size.width() + 1, offsetY);      break;
    case Bottom:
        p = QPoint(offsetX, primaryRect.bottom() - size.height() + 1);    break;
    }

    m_dockInter->SetFrontendWindowRect(p.x(), p.y(), size.width(), size.height());
}

void DockSettings::calculateWindowConfig()
{
    const int defaultHeight = AppItem::itemBaseHeight();
    const int defaultWidth = AppItem::itemBaseWidth();

    if (m_displayMode == Dock::Efficient)
    {
        switch (m_position)
        {
        case Top:
        case Bottom:
            m_mainWindowSize.setHeight(defaultHeight + PANEL_BORDER);
            m_mainWindowSize.setWidth(m_primaryRect.width());
            break;

        case Left:
        case Right:
            m_mainWindowSize.setHeight(m_primaryRect.height());
            m_mainWindowSize.setWidth(defaultWidth + PANEL_BORDER);
            break;

        default:
            Q_ASSERT(false);
        }
    }
    else if (m_displayMode == Dock::Fashion)
    {
//        int perfectWidth = 0;
//        int perfectHeight = 0;
//        const QList<DockItem *> itemList = m_itemController->itemList();
//        for (auto item : itemList)
//        {
//            switch (item->itemType())
//            {
//            case DockItem::Launcher:
//            case DockItem::App:         perfectWidth += defaultWidth;
//                                        perfectHeight += defaultHeight;             break;
//            case DockItem::Plugins:     perfectWidth += item->sizeHint().width();
//                                        perfectHeight += item->sizeHint().height(); break;
//            default:;
//            }
//        }

        int visibleItemCount = 0;
        const QList<DockItem *> itemList = m_itemController->itemList();
        for (auto item : itemList)
        {
            switch (item->itemType())
            {
            case DockItem::Launcher:
            case DockItem::App:
            case DockItem::Plugins:
                ++visibleItemCount;
                break;
            default:;
            }
        }

        const int perfectWidth = visibleItemCount * defaultWidth + PANEL_BORDER * 2 + PANEL_PADDING * 2;
        const int perfectHeight = visibleItemCount * defaultHeight + PANEL_BORDER * 2 + PANEL_PADDING * 2;
        const int calcWidth = qMin(m_primaryRect.width(), perfectWidth);
        const int calcHeight = qMin(m_primaryRect.height(), perfectHeight);
        switch (m_position)
        {
        case Top:
        case Bottom:
            m_mainWindowSize.setHeight(defaultHeight + PANEL_BORDER);
            m_mainWindowSize.setWidth(calcWidth);
            break;

        case Left:
        case Right:
            m_mainWindowSize.setHeight(calcHeight);
            m_mainWindowSize.setWidth(defaultWidth + PANEL_BORDER);
            break;

        default:
            Q_ASSERT(false);
        }
    } else {
        Q_ASSERT(false);
    }

    resetFrontendGeometry();
}
