#include "docksettings.h"
#include "panel/mainpanel.h"
#include "item/appitem.h"

#include <QDebug>
#include <QX11Info>

#include <DApplication>

#define ICON_SIZE_LARGE         48
#define ICON_SIZE_MEDIUM        36
#define ICON_SIZE_SMALL         30
#define FASHION_MODE_PADDING    30

DWIDGET_USE_NAMESPACE

DockSettings::DockSettings(QWidget *parent)
    : QObject(parent),

      m_autoHide(true),

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
    m_primaryRect = m_displayInter->primaryRect();
    m_position = Dock::Position(m_dockInter->position());
    m_displayMode = Dock::DisplayMode(m_dockInter->displayMode());
    m_hideMode = Dock::HideMode(m_dockInter->hideMode());
    m_hideState = Dock::HideState(m_dockInter->hideState());
    m_iconSize = m_dockInter->iconSize();
    AppItem::setIconBaseSize(m_iconSize * qApp->devicePixelRatio());
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

    WhiteMenu *modeSubMenu = new WhiteMenu(&m_settingsMenu);
    modeSubMenu->addAction(&m_fashionModeAct);
    modeSubMenu->addAction(&m_efficientModeAct);
    QAction *modeSubMenuAct = new QAction(tr("Mode"), this);
    modeSubMenuAct->setMenu(modeSubMenu);

    WhiteMenu *locationSubMenu = new WhiteMenu(&m_settingsMenu);
    locationSubMenu->addAction(&m_topPosAct);
    locationSubMenu->addAction(&m_bottomPosAct);
    locationSubMenu->addAction(&m_leftPosAct);
    locationSubMenu->addAction(&m_rightPosAct);
    QAction *locationSubMenuAct = new QAction(tr("Location"), this);
    locationSubMenuAct->setMenu(locationSubMenu);

    WhiteMenu *sizeSubMenu = new WhiteMenu(&m_settingsMenu);
    sizeSubMenu->addAction(&m_largeSizeAct);
    sizeSubMenu->addAction(&m_mediumSizeAct);
    sizeSubMenu->addAction(&m_smallSizeAct);
    QAction *sizeSubMenuAct = new QAction(tr("Size"), this);
    sizeSubMenuAct->setMenu(sizeSubMenu);

    WhiteMenu *statusSubMenu = new WhiteMenu(&m_settingsMenu);
    statusSubMenu->addAction(&m_keepShownAct);
    statusSubMenu->addAction(&m_keepHiddenAct);
    statusSubMenu->addAction(&m_smartHideAct);
    QAction *statusSubMenuAct = new QAction(tr("Status"), this);
    statusSubMenuAct->setMenu(statusSubMenu);

    m_settingsMenu.addAction(modeSubMenuAct);
    m_settingsMenu.addAction(locationSubMenuAct);
    m_settingsMenu.addAction(sizeSubMenuAct);
    m_settingsMenu.addAction(statusSubMenuAct);
    m_settingsMenu.setTitle("Settings Menu");

    connect(&m_settingsMenu, &WhiteMenu::triggered, this, &DockSettings::menuActionClicked);
    connect(m_dockInter, &DBusDock::PositionChanged, this, &DockSettings::onPositionChanged);
    connect(m_dockInter, &DBusDock::IconSizeChanged, this, &DockSettings::iconSizeChanged);
    connect(m_dockInter, &DBusDock::DisplayModeChanged, this, &DockSettings::displayModeChanged);
    connect(m_dockInter, &DBusDock::HideModeChanged, this, &DockSettings::hideModeChanged, Qt::QueuedConnection);
    connect(m_dockInter, &DBusDock::HideStateChanged, this, &DockSettings::hideStateChanged);
    connect(m_dockInter, &DBusDock::ServiceRestarted, this, &DockSettings::resetFrontendGeometry);

    connect(m_itemController, &DockItemController::itemInserted, this, &DockSettings::dockItemCountChanged, Qt::QueuedConnection);
    connect(m_itemController, &DockItemController::itemRemoved, this, &DockSettings::dockItemCountChanged, Qt::QueuedConnection);

    connect(m_displayInter, &DBusDisplay::PrimaryRectChanged, this, &DockSettings::primaryScreenChanged, Qt::QueuedConnection);
    connect(m_displayInter, &DBusDisplay::ScreenHeightChanged, this, &DockSettings::primaryScreenChanged, Qt::QueuedConnection);
    connect(m_displayInter, &DBusDisplay::ScreenWidthChanged, this, &DockSettings::primaryScreenChanged, Qt::QueuedConnection);

    DApplication *app = qobject_cast<DApplication*>(qApp);
    if (app) {
        connect(app, &DApplication::iconThemeChanged, this, &DockSettings::gtkIconThemeChanged);
    }

    calculateWindowConfig();
    resetFrontendGeometry();
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

int DockSettings::expandTimeout() const
{
    return m_dockInter->showTimeout();
}

int DockSettings::narrowTimeout() const
{
    return 0;
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

const QRect DockSettings::windowRect(const Position position, const bool hide) const
{
    QSize size = m_mainWindowSize;
    if (hide)
    {
        switch (position)
        {
        case Top:
        case Bottom:    size.setHeight(1);      break;
        case Left:
        case Right:     size.setWidth(1);       break;
        }
    }

    const QRect primaryRect = m_primaryRect;
    const int offsetX = (primaryRect.width() - size.width()) / 2;
    const int offsetY = (primaryRect.height() - size.height()) / 2;

    QPoint p(0, 0);
    switch (position)
    {
    case Top:
        p = QPoint(offsetX, 0);                                     break;
    case Left:
        p = QPoint(0, offsetY);                                     break;
    case Right:
        p = QPoint(primaryRect.width() - size.width(), offsetY);    break;
    case Bottom:
        p = QPoint(offsetX, primaryRect.height() - size.height());  break;
    default:Q_UNREACHABLE();
    }

//    qDebug() << Q_FUNC_INFO << m_mainWindowSize << primaryRect << p << size;

    return QRect(primaryRect.topLeft() + p, size);
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

    m_settingsMenu.exec(QCursor::pos());

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

void DockSettings::menuActionClicked(QAction *action)
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

void DockSettings::onPositionChanged()
{
//    qDebug() << Q_FUNC_INFO;
    const Position prevPos = m_position;
    const Position nextPos = Dock::Position(m_dockInter->position());

    if (prevPos == nextPos)
        return;

    m_position = nextPos;
    DockItem::setDockPosition(m_position);
    qApp->setProperty(PROP_POSITION, QVariant::fromValue(m_position));

    calculateWindowConfig();

    emit positionChanged(prevPos);

    m_itemController->refershItemsIcon();
}

void DockSettings::iconSizeChanged()
{
//    qDebug() << Q_FUNC_INFO;
    m_iconSize = m_dockInter->iconSize();
    AppItem::setIconBaseSize(m_iconSize * qApp->devicePixelRatio());

    calculateWindowConfig();

    emit dataChanged();
}

void DockSettings::displayModeChanged()
{
//    qDebug() << Q_FUNC_INFO;
    m_displayMode = Dock::DisplayMode(m_dockInter->displayMode());
    DockItem::setDockDisplayMode(m_displayMode);
    qApp->setProperty(PROP_DISPLAY_MODE, QVariant::fromValue(m_displayMode));

    calculateWindowConfig();

    emit dataChanged();
}

void DockSettings::hideModeChanged()
{
//    qDebug() << Q_FUNC_INFO;
    m_hideMode = Dock::HideMode(m_dockInter->hideMode());

    emit windowHideModeChanged();
}

void DockSettings::hideStateChanged()
{
//    qDebug() << Q_FUNC_INFO;
    const Dock::HideState state = Dock::HideState(m_dockInter->hideState());

    if (state == Dock::Unknown)
        return;

    m_hideState = state;

    emit windowVisibleChanged();
}

void DockSettings::dockItemCountChanged()
{
//    qDebug() << Q_FUNC_INFO;
    if (m_displayMode == Dock::Efficient)
        return;

    calculateWindowConfig();

    emit windowGeometryChanged();
}

void DockSettings::primaryScreenChanged()
{
//    qDebug() << Q_FUNC_INFO;
    m_primaryRect = m_displayInter->primaryRect();

    calculateWindowConfig();

    emit dataChanged();
}

void DockSettings::resetFrontendGeometry()
{
    const QRect r = windowRect(m_position);
//    qDebug() << Q_FUNC_INFO << r;

    m_dockInter->SetFrontendWindowRect(r.x(), r.y(), r.width(), r.height());
}

void DockSettings::calculateWindowConfig()
{
//    qDebug() << Q_FUNC_INFO;
    const auto ratio = qApp->devicePixelRatio();
    const int defaultHeight = AppItem::itemBaseHeight() / ratio;
    const int defaultWidth = AppItem::itemBaseWidth() / ratio;

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
        const int calcWidth = qMin(m_primaryRect.width() - FASHION_MODE_PADDING * 2, perfectWidth);
        const int calcHeight = qMin(m_primaryRect.height() - FASHION_MODE_PADDING * 2, perfectHeight);
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

void DockSettings::gtkIconThemeChanged()
{
    qDebug() << Q_FUNC_INFO;
    m_itemController->refershItemsIcon();
}
