#include "mainwidget.h"

MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),50);
    Panel * mainPanel = new Panel(this);
    mainPanel->setMinimumSize(this->width(),this->height());
    mainPanel->resize(this->width(),this->height());
    mainPanel->move(0,0);

    this->setWindowFlags(Qt::ToolTip);
    this->setAttribute(Qt::WA_TranslucentBackground);
    this->move(0,800);
}

MainWidget::~MainWidget()
{

}
