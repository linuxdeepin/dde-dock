#include "dockapplayout.h"

DockAppLayout::DockAppLayout(QWidget *parent) : MovableLayout(parent)
{
    initAppManager();
}

QSize DockAppLayout::sizeHint() const
{
    QSize size;
    int w = 0;
    int h = 0;
    switch (direction()) {
    case QBoxLayout::LeftToRight:
    case QBoxLayout::RightToLeft:
        size.setHeight(height());
        for (QWidget * widget : widgets()) {
            w += widget->width();
        }
        size.setWidth(w + getLayoutSpacing() * widgets().count());
        break;
    case QBoxLayout::TopToBottom:
    case QBoxLayout::BottomToTop:
        size.setWidth(width());
        for (QWidget * widget : widgets()) {
            h += widget->height();
        }
        size.setHeight(h + getLayoutSpacing() * widgets().count());
        break;
    }

    return size;
}

void DockAppLayout::initEntries()
{
    m_appManager->initEntries();
}

void DockAppLayout::initAppManager()
{
    m_appManager = new DockAppManager(this);
    connect(m_appManager, &DockAppManager::entryAdded, this, &DockAppLayout::onAppItemAdd);
    connect(m_appManager, &DockAppManager::entryRemoved, this, &DockAppLayout::onAppItemRemove);

    //Make sure the item which was dragged to the dock can be show at once
    //    connect(m_appLayout, &DockLayout::itemDocking, m_appManager, &AppManager::setDockingItemId);
}

void DockAppLayout::onAppItemRemove(const QString &id)
{
    QList<QWidget *> tmpList = this->widgets();
    for (QWidget * item : tmpList) {
        DockAppItem *tmpItem = qobject_cast<DockAppItem *>(item);
        if (tmpItem && tmpItem->getItemId() == id) {
            this->removeWidget(item);
            tmpItem->setVisible(false);
            tmpItem->deleteLater();
            return;
        }
    }
}

void DockAppLayout::onAppItemAdd(DockAppItem *item, bool delayShow)
{
    this->addWidget(item);
    connect(item, &DockAppItem::needPreviewShow, this, &DockAppLayout::needPreviewShow);
    connect(item, &DockAppItem::needPreviewHide, this, &DockAppLayout::needPreviewHide);
    connect(item, &DockAppItem::needPreviewUpdate, this, &DockAppLayout::needPreviewUpdate);
}

