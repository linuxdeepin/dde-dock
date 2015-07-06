#include "mainwidget.h"

MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockModeData::getInstants()->getDockHeight());
    mainPanel = new Panel(this);
    mainPanel->resize(this->width(),this->height());
    mainPanel->move(0,0);

    this->setWindowFlags(Qt::WindowStaysOnTopHint | Qt::FramelessWindowHint);
    this->setAttribute(Qt::WA_TranslucentBackground);
    this->move(0,rec.height());

    connect(DockModeData::getInstants(), SIGNAL(dockModeChanged(Dock::DockMode,Dock::DockMode)),
            this, SLOT(slotDockModeChanged(Dock::DockMode,Dock::DockMode)));
}

void MainWidget::slotDockModeChanged(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockModeData::getInstants()->getDockHeight());

//    mainPanel->resize(this->width(),this->height());
//    mainPanel->move(0,0);
}

MainWidget::~MainWidget()
{

}
