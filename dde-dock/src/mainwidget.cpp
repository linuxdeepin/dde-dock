#include "mainwidget.h"
#include "xcb_misc.h"

MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    initHideStateManager();
    initDockSetting();

    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(rec.width(),m_dmd->getDockHeight());
    this->move((rec.width() - width()) / 2, rec.height() - this->height());

    m_mainPanel = new Panel(this);
    connect(m_mainPanel,&Panel::startShow,this,&MainWidget::showDock);
    connect(m_mainPanel,&Panel::panelHasHidden,this,&MainWidget::hideDock);
    connect(m_mainPanel, &Panel::sizeChanged, this, &MainWidget::onPanelSizeChanged);

    this->setWindowFlags(Qt::Window);
    this->setAttribute(Qt::WA_TranslucentBackground);

    connect(m_dmd, &DockModeData::dockModeChanged, this, &MainWidget::changeDockMode);

    //For init
    changeDockMode(m_dmd->getDockMode(), m_dmd->getDockMode());

    DockUIDbus *dockUIDbus = new DockUIDbus(this);
    Q_UNUSED(dockUIDbus)

    XcbMisc::instance()->set_window_type(winId(),
                                         XcbMisc::Dock);
}

void MainWidget::changeDockMode(Dock::DockMode, Dock::DockMode)
{
    if (hasHidden)
        return;

    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(rec.width(),m_dmd->getDockHeight());
    this->move((rec.width() - width()) / 2, rec.height() - this->height());

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
        m_mainPanel->setContainMouse(true);
        m_mainPanel->startShow();
    }
}

void MainWidget::leaveEvent(QEvent *)
{
    m_mainPanel->setContainMouse(false);
}

void MainWidget::showDock()
{
    hasHidden = false;
    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(m_mainPanel->width(), m_dmd->getDockHeight());
    this->move((rec.width() - width()) / 2, rec.height() - this->height());
    updateXcbStructPartial();
}

void MainWidget::hideDock()
{
    hasHidden = true;
    QRect rec = QApplication::desktop()->screenGeometry();
    //set height with 0 mean window is hidden,Windows manager will handle it's showing animation
    this->setFixedSize(m_mainPanel->width(), 1);
    this->move((rec.width() - width()) / 2, rec.height() - 1);//1 pixel for grab mouse enter event to show panel
    updateXcbStructPartial();
}

void MainWidget::onPanelSizeChanged()
{
    this->setFixedSize(m_mainPanel->size());
    QRect rec = QApplication::desktop()->screenGeometry();
    this->move((rec.width() - width()) / 2, y());
}

MainWidget::~MainWidget()
{

}

void MainWidget::loadResources()
{
    m_mainPanel->loadResources();
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
