#include "panel.h"
#include "dockpluginproxy.h"
#include "dockpluginmanager.h"
#include "Controller/dockmodedata.h"
#include <QHBoxLayout>

Panel::Panel(QWidget *parent)
    : QLabel(parent),parentWidget(parent)
{
    this->setObjectName("Panel");

    rightLayout = new DockLayout(this);
    rightLayout->setSortDirection(DockLayout::RightToLeft);
    rightLayout->setSpacing(dockCons->getAppletsItemSpacing());
    rightLayout->resize(0,dockCons->getItemHeight());

    leftLayout = new DockLayout(this);
    leftLayout->setSpacing(dockCons->getAppItemSpacing());
    leftLayout->resize(this->width() - rightLayout->width(),dockCons->getItemHeight());
    leftLayout->move(0,1);
    connect(leftLayout,SIGNAL(dragStarted()),this,SLOT(slotDragStarted()));
    connect(leftLayout,SIGNAL(itemDropped()),this,SLOT(slotItemDropped()));

    connect(leftLayout, SIGNAL(contentsWidthChange()),this, SLOT(slotLayoutContentsWidthChanged()));
    connect(rightLayout, SIGNAL(contentsWidthChange()), this, SLOT(slotLayoutContentsWidthChanged()));

    connect(dockCons, SIGNAL(dockModeChanged(Dock::DockMode,Dock::DockMode)),
            this, SLOT(slotDockModeChanged(Dock::DockMode,Dock::DockMode)));

    DockPluginManager *pluginManager = new DockPluginManager(this);
    connect(DockModeData::instance(), &DockModeData::dockModeChanged,
            pluginManager, &DockPluginManager::onDockModeChanged);

    QList<DockPluginProxy*> proxies = pluginManager->getAll();
    foreach (DockPluginProxy* proxy, proxies) {
        connect(proxy, &DockPluginProxy::itemAdded, [=](AbstractDockItem* item) {
            rightLayout->addItem(item);
        });
        connect(proxy, &DockPluginProxy::itemRemoved, [=](AbstractDockItem* item) {
            int index = rightLayout->indexOf(item);
            if (index != -1) {
                rightLayout->removeItem(index);
            }
        });

        proxy->plugin()->init(proxy);
    }

    initAppManager();
}

void Panel::showScreenMask()
{
//    qWarning() << "[Info:]" << "Show Screen Mask.";
    maskWidget = new ScreenMask();
    connect(maskWidget,SIGNAL(itemDropped(QPoint)),this,SLOT(slotItemDropped()));
    connect(maskWidget,SIGNAL(itemEntered()),this,SLOT(slotEnteredMask()));
    connect(maskWidget,SIGNAL(itemExited()),this,SLOT(slotExitedMask()));
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

void Panel::slotDockModeChanged(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    leftLayout->relayout();
    rightLayout->relayout();

    reanchorsLayout(newMode);

    qWarning() << "AppCount:********" << leftLayout->getItemCount();
}

void Panel::slotLayoutContentsWidthChanged()
{
    reanchorsLayout(dockCons->getDockMode());
}

void Panel::slotAddAppItem(AppItem *item)
{
    leftLayout->addItem(item);
}

void Panel::slotRemoveAppItem(const QString &id)
{
    QList<AbstractDockItem *> tmpList = leftLayout->getItemList();
    for (int i = 0; i < tmpList.count(); i ++)
    {
        AppItem *tmpItem = qobject_cast<AppItem *>(tmpList.at(i));
        if (tmpItem->itemId() == id)
        {
            leftLayout->removeItem(i);
            return;
        }
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

void Panel::reanchorsLayout(Dock::DockMode mode)
{
    if (mode == Dock::FashionMode)
    {
        leftLayout->resize(leftLayout->getContentsWidth() + dockCons->getAppItemSpacing(),dockCons->getItemHeight());
        rightLayout->setSortDirection(DockLayout::LeftToRight);
        rightLayout->resize(rightLayout->getContentsWidth(),dockCons->getItemHeight());
        rightLayout->move(leftLayout->width() - dockCons->getAppItemSpacing(),1);

        this->resize(leftLayout->getContentsWidth() + rightLayout->getContentsWidth(),dockCons->getDockHeight());
        this->move((parentWidget->width() - leftLayout->getContentsWidth() - rightLayout->getContentsWidth()) / 2,0);
    }
    else
    {
        rightLayout->setSortDirection(DockLayout::RightToLeft);
        rightLayout->resize(rightLayout->getContentsWidth(),dockCons->getItemHeight());
        rightLayout->move(parentWidget->width() - rightLayout->width(),1);

        leftLayout->resize(parentWidget->width() - rightLayout->width() ,dockCons->getItemHeight());

        this->resize(leftLayout->width() + rightLayout->width(),dockCons->getDockHeight());
        this->move((parentWidget->width() - leftLayout->width() - rightLayout->width()) / 2,0);
    }
}

void Panel::showMenu()
{
    QPoint tmpPos = QCursor::pos();

    PanelMenu::instance()->showMenu(tmpPos.x(),tmpPos.y());
}

void Panel::initAppManager()
{
    m_appManager = new AppManager(this);
    connect(m_appManager,SIGNAL(entryAdded(AppItem*)),this, SLOT(slotAddAppItem(AppItem*)));
    connect(m_appManager, SIGNAL(entryRemoved(QString)),this, SLOT(slotRemoveAppItem(QString)));
    m_appManager->updateEntries();
}

Panel::~Panel()
{

}

