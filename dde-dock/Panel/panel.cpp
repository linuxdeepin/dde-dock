#include "panel.h"

Panel::Panel(QWidget *parent) : QWidget(parent)
{
    leftLayout = new DockLayout(this);
    leftLayout->resize(1024,50);
    leftLayout->move(0,0);

    AppItem * b1 = new AppItem("App",":/test/Resources/images/brasero.png");b1->resize(50,50);
    AppItem * b2 = new AppItem("App",":/test/Resources/images/crossover.png");b2->resize(50,50);
    AppItem * b3 = new AppItem("App",":/test/Resources/images/gcr-gnupg.png");b3->resize(50,50);
    AppItem * b4 = new AppItem("App",":/test/Resources/images/display-im6.q16.png");b4->resize(50,50);
    AppItem * b5 = new AppItem("App",":/test/Resources/images/eog.png");b5->resize(50,50);

    leftLayout->addItem(b1);
    leftLayout->addItem(b2);
    leftLayout->addItem(b3);
    leftLayout->addItem(b4);
    leftLayout->addItem(b5);
}

void Panel::resize(const QSize &size)
{
    QWidget::resize(size);
    leftLayout->resize(this->width() * 2 / 3,this->height());
}

void Panel::resize(int width, int height)
{
    QWidget::resize(width,height);
    leftLayout->resize(this->width() * 2 / 3,this->height());
}

Panel::~Panel()
{

}

