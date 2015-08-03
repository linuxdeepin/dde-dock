#include "mainwidget.h"
#include "xcb_misc.h"

MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    initHideStateManager();
    initDockSetting();

    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(rec.width(),m_dmd->getDockHeight());
    this->move(0, rec.height() - this->height());

    mainPanel = new Panel(this);
    connect(mainPanel,&Panel::startShow,this,&MainWidget::showDock);
    connect(mainPanel,&Panel::panelHasHidden,this,&MainWidget::hideDock);

    this->setWindowFlags(Qt::Window);
    this->setAttribute(Qt::WA_TranslucentBackground);

    connect(m_dmd, &DockModeData::dockModeChanged, this, &MainWidget::changeDockMode);

    //For init
    changeDockMode(m_dmd->getDockMode(), m_dmd->getDockMode());

    DockUIDbus *dockUIDbus = new DockUIDbus(this);

    XcbMisc::instance()->set_window_type(winId(),
                                         XcbMisc::Dock);
}

void MainWidget::changeDockMode(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (hasHidden)
        return;

    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(rec.width(),m_dmd->getDockHeight());
    this->move(0, rec.height() - this->height());

    updateXcbStructPartial();
}

void MainWidget::updateXcbStructPartial()
{
    int tmpHeight = 0;
    if (m_dds->GetHideMode() == 0)
        tmpHeight = this->height();
    XcbMisc::instance()->set_strut_partial(winId(),
                                           XcbMisc::OrientationBottom,
                                           tmpHeight,
                                           x(),
                                           x() + width());
}

void MainWidget::initHideStateManager()
{
    m_dhsm = new DBusHideStateManager(this);
    m_dhsm->SetState(1);
}

void MainWidget::initDockSetting()
{
    m_dds = new DBusDockSetting(this);
    connect(m_dds, &DBusDockSetting::HideModeChanged, this, &MainWidget::updateXcbStructPartial);
}

void MainWidget::enterEvent(QEvent *event)
{
    if (height() == 1){
        showDock();
        mainPanel->setContainMouse(true);
        mainPanel->startShow();
    }
}

void MainWidget::leaveEvent(QEvent *)
{
    mainPanel->setContainMouse(false);
}

void MainWidget::showDock()
{
    hasHidden = false;
    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(rec.width(),m_dmd->getDockHeight());
    this->move(0, rec.height() - this->height());
    updateXcbStructPartial();
}

void MainWidget::hideDock()
{
    hasHidden = true;
    QRect rec = QApplication::desktop()->screenGeometry();
    //set height with 0 mean window is hidden,Windows manager will handle it's showing animation
    this->setFixedSize(rec.width(),1);
    this->move(0, rec.height() - 1);//1 pixel for grab mouse enter event to show panel
    updateXcbStructPartial();
}

MainWidget::~MainWidget()
{

}

DockUIDbus::DockUIDbus(MainWidget *parent):
    QDBusAbstractAdaptor(parent),
    m_parent(parent)
{
    QDBusConnection::sessionBus().registerObject(DBUS_PATH, parent);
}

DockUIDbus::~DockUIDbus()
{

}

qulonglong DockUIDbus::Xid()
{
    return 0;
}
