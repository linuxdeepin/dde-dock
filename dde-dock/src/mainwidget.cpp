#include "mainwidget.h"

MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockConstants::getInstants()->getDockHeight());
    mainPanel = new Panel(this);
    mainPanel->resize(this->width(),this->height());
    mainPanel->move(0,0);

    this->setWindowFlags(Qt::ToolTip);
    this->setAttribute(Qt::WA_TranslucentBackground);
    this->move(0,rec.height()  - 100);

    connect(DockConstants::getInstants(), SIGNAL(dockModeChanged(DockConstants::DockMode,DockConstants::DockMode)),
            this, SLOT(slotDockModeChanged(DockConstants::DockMode,DockConstants::DockMode)));
}

void MainWidget::slotDockModeChanged(DockConstants::DockMode newMode, DockConstants::DockMode oldMode)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockConstants::getInstants()->getDockHeight());

//    mainPanel->resize(this->width(),this->height());
//    mainPanel->move(0,0);
}

MainWidget::~MainWidget()
{

}
