#include "mainwidget.h"

MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockModeData::instance()->getDockHeight());
    mainPanel = new Panel(this);
    mainPanel->resize(this->width(),this->height());
    mainPanel->move(0,0);
    connect(mainPanel,&Panel::startShow,this,&MainWidget::showDock);
    connect(mainPanel,&Panel::panelHasHidden,this,&MainWidget::hideDock);

    this->setWindowFlags(Qt::WindowStaysOnTopHint | Qt::FramelessWindowHint);
    this->setAttribute(Qt::WA_TranslucentBackground);
    this->move(0, rec.height() - this->height());

    connect(DockModeData::instance(), SIGNAL(dockModeChanged(Dock::DockMode,Dock::DockMode)),
            this, SLOT(slotDockModeChanged(Dock::DockMode,Dock::DockMode)));
}

void MainWidget::slotDockModeChanged(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (hasHidden)
        return;

    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockModeData::instance()->getDockHeight());

//    mainPanel->resize(this->width(),this->height());
//    mainPanel->move(0,0);
}

void MainWidget::showDock()
{
    hasHidden = false;
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockModeData::instance()->getDockHeight());
}

void MainWidget::hideDock()
{
    hasHidden = true;
    QRect rec = QApplication::desktop()->screenGeometry();
    //set height with 0 mean window is hidden,Windows manager will handle it's showing animation
    this->resize(rec.width(),1);
}

MainWidget::~MainWidget()
{

}
