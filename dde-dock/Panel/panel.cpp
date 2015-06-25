#include "panel.h"

Panel::Panel(QWidget *parent) : QWidget(parent)
{
    leftLayout = new DockLayout(this);
    leftLayout->resize(1024,50);
    leftLayout->move(0,0);

    for (int i = 0; i < 5; i ++)
    {
        AppItem * tmpButton = new AppItem("App" + QString::number(i),":/test/Resources/images/google-chrome.png");
        tmpButton->resize(50,50);

        leftLayout->addItem(tmpButton);
    }
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

