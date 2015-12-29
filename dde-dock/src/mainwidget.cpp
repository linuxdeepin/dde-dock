#include "mainwidget.h"
#include "xcb_misc.h"

const int ENTER_DELAY_INTERVAL = 600;
MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    this->setWindowFlags(Qt::Window | Qt::WindowDoesNotAcceptFocus);
    this->setAttribute(Qt::WA_TranslucentBackground);
    //the attribute "Qt::WA_X11DoNotAcceptFocus" will not tack effect, not know the reason
    //this->setAttribute(Qt::WA_X11DoNotAcceptFocus);

    initHideStateManager();

#ifdef NEW_DOCK_LAYOUT
    m_mainPanel = new DockPanel(this);
    connect(m_mainPanel,&DockPanel::panelHasHidden,this,&MainWidget::hideDock);
    connect(m_mainPanel, &DockPanel::sizeChanged, this, &MainWidget::onPanelSizeChanged);
#else
    m_mainPanel = new Panel(this);
    connect(m_mainPanel,&Panel::panelHasHidden,this,&MainWidget::hideDock);
    connect(m_mainPanel, &Panel::sizeChanged, this, &MainWidget::onPanelSizeChanged);
 #endif

    connect(m_dmd, &DockModeData::dockModeChanged, this, &MainWidget::onDockModeChanged);

    //For init
    m_display = new DBusDisplay(this);
    updatePosition();

    DockUIDbus *dockUIDbus = new DockUIDbus(this);
    Q_UNUSED(dockUIDbus)

    XcbMisc::instance()->set_window_type(winId(), XcbMisc::Dock);

#ifdef NEW_DOCK_LAYOUT
    connect(m_display, &DBusDisplay::PrimaryChanged, [=] {
//        m_mainPanel->resizeWithContent();
        updatePosition();
    });
    connect(m_display, &DBusDisplay::PrimaryRectChanged, [=] {
//        m_mainPanel->resizeWithContent();
        updatePosition();
    });
#else
    connect(m_display, &DBusDisplay::PrimaryChanged, [=] {
        m_mainPanel->resizeWithContent();
        updatePosition();
    });
    connect(m_display, &DBusDisplay::PrimaryRectChanged, [=] {
        m_mainPanel->resizeWithContent();
        updatePosition();
    });
#endif
}

void MainWidget::onDockModeChanged()
{
    updatePosition();
}

void MainWidget::updatePosition()
{
    DisplayRect rec = m_display->primaryRect();

    if (m_hasHidden) {
        //set height with 0 mean window is hidden,Windows manager will handle it's showing animation
        this->setFixedSize(m_mainPanel->width(), 1);
        this->move(rec.x + (rec.width - width()) / 2,
                   rec.y + rec.height - 1);//1 pixel for grab mouse enter event to show panel
    }
    else {
        this->setFixedSize(m_mainPanel->width(), m_dmd->getDockHeight());
        this->move(rec.x + (rec.width - width()) / 2,
                   rec.y + rec.height - this->height());
    }

    updateXcbStructPartial();
}

void MainWidget::updateXcbStructPartial()
{
    int tmpHeight = 0;
    DBusDockSetting dds;
    if (dds.GetHideMode() == 0)
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
    m_dhsm->SetState(Dock::HideStateHiding);
    connect(m_dhsm, &DBusHideStateManager::ChangeState, [=](int state) {
        if (state == Dock::HideStateShowing) {
            showDock();
        }
    });
}

void MainWidget::enterEvent(QEvent *)
{
    if (height() == 1){
        QTimer *st = new QTimer(this);
        connect(st, &QTimer::timeout, this, [=] {
            //make sure the panel will show by mouse-enter
            if (geometry().contains(QCursor::pos())) {
                qDebug() << "MouseEntered, show dock...";
                showDock();
            }
            sender()->deleteLater();
        });
        st->start(ENTER_DELAY_INTERVAL);
    }
}

void MainWidget::leaveEvent(QEvent *)
{
    if (!this->geometry().contains(QCursor::pos()))
        m_dhsm->UpdateState();
}

void MainWidget::showDock()
{
    m_hasHidden = false;
    updatePosition();
    m_mainPanel->startShow();
}

void MainWidget::hideDock()
{
    m_hasHidden = true;
    updatePosition();
}

void MainWidget::onPanelSizeChanged()
{
    updatePosition();
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
    return m_parent->winId();
}
