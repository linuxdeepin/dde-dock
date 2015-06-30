#include "panel.h"

Panel::Panel(QWidget *parent)
    : QLabel(parent),parentWidget(parent)
{
    this->setStyleSheet("QWidget{background-color: rgba(0,0,0,0.3);}");

    leftLayout = new DockLayout(this);
    leftLayout->resize(1024,50);
    leftLayout->move(0,0);

    AppItem * b1 = new AppItem("App",":/test/Resources/images/brasero.png");
    AppItem * b2 = new AppItem("App",":/test/Resources/images/crossover.png");
    AppItem * b3 = new AppItem("App",":/test/Resources/images/vim.png");
    AppItem * b4 = new AppItem("App",":/test/Resources/images/google-chrome.png");
    AppItem * b5 = new AppItem("App",":/test/Resources/images/QtProject-qtcreator.png");

    leftLayout->addItem(b1);
    leftLayout->addItem(b2);
    leftLayout->addItem(b3);
    leftLayout->addItem(b4);
    leftLayout->addItem(b5);

    connect(leftLayout,SIGNAL(dragStarted()),this,SLOT(slotDragStarted()));
    connect(leftLayout,SIGNAL(itemDropped()),this,SLOT(slotItemDropped()));

    AppItem * b6 = new AppItem("App","../Resources/images/display-im6.q16.png");
    AppItem * b7 = new AppItem("App","../Resources/images/eog.png");
    rightLayout = new DockLayout(this);
    rightLayout->setSortDirection(DockLayout::RightToLeft);
    rightLayout->resize(300,50);
    rightLayout->move(0,0);
    rightLayout->addItem(b6);
    rightLayout->addItem(b7);
}

void Panel::resize(const QSize &size)
{
    QWidget::resize(size);
    leftLayout->resize(this->width() * 2 / 3,this->height());
    rightLayout->move(this->width() - rightLayout->width(),0);
}

void Panel::resize(int width, int height)
{
    QWidget::resize(width,height);
    leftLayout->resize(this->width() * 2 / 3,this->height());
    rightLayout->move(this->width() - rightLayout->width(),0);
}

void Panel::showScreenMask()
{
//    qWarning() << "[Info:]" << "Show Screen Mask.";
    maskWidget = new ScreenMask();
    connect(maskWidget,SIGNAL(itemDropped(QPoint)),this,SLOT(slotItemDropped()));
    connect(maskWidget,SIGNAL(itemEntered()),this,SLOT(slotEnteredMask()));
    connect(maskWidget,SIGNAL(itemExited()),this,SLOT(slotExitedMask()));

    //TODO change to Other ways to do this,it will hide the drag icon
    parentWidget->hide();
    parentWidget->show();
}

void Panel::hideScreenMask()
{
//    qWarning() << "[Info:]" << "Hide Screen Mask.";
    disconnect(maskWidget,SIGNAL(itemDropped(QPoint)),this,SLOT(slotItemDropped()));
    disconnect(maskWidget,SIGNAL(itemEntered()),this,SLOT(slotEnteredMask()));
    disconnect(maskWidget,SIGNAL(itemExited()),this,SLOT(slotExitedMask()));
    maskWidget->hide();
    maskWidget->deleteLater();
    maskWidget = NULL;
}

void Panel::slotDragStarted()
{
    showScreenMask();
}

void Panel::slotItemDropped()
{
    hideScreenMask();
    leftLayout->relayout();
}

void Panel::slotEnteredMask()
{
    leftLayout->relayout();
}

void Panel::slotExitedMask()
{
    leftLayout->addSpacingItem();
//    leftLayout->relayout();
}

Panel::~Panel()
{

}

