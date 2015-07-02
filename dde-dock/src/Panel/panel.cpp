#include "panel.h"
#include "systraymanager.h"

Panel::Panel(QWidget *parent)
    : QLabel(parent),parentWidget(parent)
{
    this->setObjectName("Panel");

    rightLayout = new DockLayout(this);
    rightLayout->setSortDirection(DockLayout::RightToLeft);
    rightLayout->setSpacing(dockCons->getAppletsItemSpacing());
    rightLayout->resize(80,dockCons->getDockHeight());

    leftLayout = new DockLayout(this);
    leftLayout->setSpacing(dockCons->getAppItemSpacing());
    leftLayout->resize(this->width() - rightLayout->width(),dockCons->getDockHeight());
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

    connect(leftLayout, SIGNAL(contentsWidthChange()),this, SLOT(slotLayoutContentsWidthChanged()));
    connect(rightLayout, SIGNAL(contentsWidthChange()), this, SLOT(slotLayoutContentsWidthChanged()));

    connect(dockCons, SIGNAL(dockModeChanged(DockConstants::DockMode,DockConstants::DockMode)),
            this, SLOT(slotDockModeChanged(DockConstants::DockMode,DockConstants::DockMode)));

    SystrayManager *manager = new SystrayManager();
    foreach (AbstractDockItem *item, manager->trayIcons()) {
        rightLayout->addItem(item);
    }

    panelMenu = new PanelMenu();

    slotDockModeChanged(dockCons->getDockMode(),dockCons->getDockMode());






    ///////////////////////////
    AppManager *appManager = new AppManager(this);
}

void Panel::resize(const QSize &size)
{
    QWidget::resize(size);

    reanchorsLayout(dockCons->getDockMode());
}

void Panel::resize(int width, int height)
{
    QWidget::resize(width,height);

    reanchorsLayout(dockCons->getDockMode());
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
//    leftLayout->relayout();
}

void Panel::slotDockModeChanged(DockConstants::DockMode newMode, DockConstants::DockMode oldMode)
{
    leftLayout->relayout();
    rightLayout->relayout();

    reanchorsLayout(newMode);

    this->resize(leftLayout->width() + rightLayout->width(),dockCons->getDockHeight());
    this->move((parentWidget->width() - leftLayout->width() - rightLayout->width()) / 2,0);
}

void Panel::slotLayoutContentsWidthChanged()
{
    reanchorsLayout(dockCons->getDockMode());

    if (dockCons->getDockMode() == DockConstants::FashionMode)
    {
        this->resize(leftLayout->getContentsWidth() + rightLayout->getContentsWidth(),dockCons->getDockHeight());
        this->move((parentWidget->width() - leftLayout->getContentsWidth() - rightLayout->getContentsWidth()) / 2,0);
    }
    else
    {
        this->resize(leftLayout->width() + rightLayout->width(),dockCons->getDockHeight());
        this->move((parentWidget->width() - leftLayout->width() - rightLayout->width()) / 2,0);
    }

}

void Panel::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::RightButton)
        showMenu();
}

void Panel::mouseReleaseEvent(QMouseEvent *event)
{

}

void Panel::reanchorsLayout(DockConstants::DockMode mode)
{
    if (mode == DockConstants::FashionMode)
    {
        leftLayout->resize(leftLayout->getContentsWidth() + dockCons->getAppItemSpacing(),dockCons->getDockHeight());

        rightLayout->setSortDirection(DockLayout::LeftToRight);
        rightLayout->resize(rightLayout->getContentsWidth(),dockCons->getDockHeight());
        rightLayout->move(leftLayout->width() - dockCons->getAppItemSpacing(),0);
    }
    else
    {
        rightLayout->setSortDirection(DockLayout::RightToLeft);
        rightLayout->resize(rightLayout->getContentsWidth(),dockCons->getDockHeight());
        rightLayout->move(parentWidget->width() - rightLayout->width(),0);

        leftLayout->resize(parentWidget->width() - rightLayout->width() ,dockCons->getDockHeight());
    }
}

void Panel::showMenu()
{
    QPoint tmpPos = QCursor::pos();

    panelMenu->move(tmpPos.x(),tmpPos.y() - panelMenu->height());
    panelMenu->show();
}

void Panel::hideMenu()
{

}

Panel::~Panel()
{

}

