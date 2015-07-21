#include "mainwidget.h"
#include "xcb_misc.h"

MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(rec.width(),m_dmd->getDockHeight());
    this->move(0, rec.height() - this->height());

    mainPanel = new Panel(this);
    connect(mainPanel,&Panel::startShow,this,&MainWidget::showDock);
    connect(mainPanel,&Panel::panelHasHidden,this,&MainWidget::hideDock);

    this->setWindowFlags(Qt::WindowStaysOnTopHint | Qt::FramelessWindowHint);
    this->setAttribute(Qt::WA_TranslucentBackground);

    connect(m_dmd, &DockModeData::dockModeChanged, this, &MainWidget::changeDockMode);

    //For init
    changeDockMode(m_dmd->getDockMode(), m_dmd->getDockMode());
}

void MainWidget::changeDockMode(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (hasHidden)
        return;

    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(rec.width(),m_dmd->getDockHeight());
    this->move(0, rec.height() - this->height());

    XcbMisc::instance()->set_window_type(winId(),
                                         XcbMisc::Dock);

    XcbMisc::instance()->set_strut_partial(winId(),
                                           XcbMisc::OrientationBottom,
                                           height(),
                                           x(),
                                           x() + width());
}

void MainWidget::showDock()
{
    hasHidden = false;
    QRect rec = QApplication::desktop()->screenGeometry();
    this->setFixedSize(rec.width(),m_dmd->getDockHeight());
}

void MainWidget::hideDock()
{
    hasHidden = true;
    QRect rec = QApplication::desktop()->screenGeometry();
    //set height with 0 mean window is hidden,Windows manager will handle it's showing animation
    this->setFixedSize(rec.width(),1);
}

MainWidget::~MainWidget()
{

}
