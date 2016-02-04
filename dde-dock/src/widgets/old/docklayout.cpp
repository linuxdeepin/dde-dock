/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "docklayout.h"
#include "abstractdockitem.h"

DockLayout::DockLayout(QWidget *parent) :
    QWidget(parent)
{
    this->setAttribute(Qt::WA_TranslucentBackground);

    this->installEventFilter(this);
    this->setMouseTracking(true);
}

void DockLayout::addItem(AbstractDockItem *item, bool delayShow)
{
    if (!item)
        return;

    if (m_lastHoverIndex == -1)
        insertItem(item, m_appList.count(), delayShow);
    else
        insertItem(item, m_lastHoverIndex, delayShow);
}

void DockLayout::appendItem(AbstractDockItem *item, bool delayShow)
{
    if (!item)
        return;

    insertItem(item, m_appList.count(), delayShow);
}

void DockLayout::insertItem(AbstractDockItem *item, int index, bool delayShow)
{
    QPointer<AbstractDockItem> pItem = item;
    if (pItem.isNull())
        return;

    pItem->setParent(this);
    pItem->show();
    int appCount = m_appList.count();
    index = index > appCount ? appCount : (index < 0 ? 0 : index);

    m_appList.insert(index,pItem);

    connect(pItem, &AbstractDockItem::frameUpdate, this, &DockLayout::frameUpdate);
    connect(pItem, &AbstractDockItem::posChanged, this, &DockLayout::frameUpdate);
    connect(pItem, &AbstractDockItem::mouseRelease, this, &DockLayout::slotItemRelease);
    connect(pItem, &AbstractDockItem::dragStart, this, &DockLayout::slotItemDrag);
    connect(pItem, &AbstractDockItem::dragEntered, this, &DockLayout::slotItemEntered);
    connect(pItem, &AbstractDockItem::dragExited, this, &DockLayout::slotItemExited);
    connect(pItem, &AbstractDockItem::widthChanged, this, &DockLayout::relayout);
    connect(pItem, &AbstractDockItem::moveAnimationFinished,this, &DockLayout::slotAnimationFinish);
    connect(this, &DockLayout::itemHoverableChange, pItem, &AbstractDockItem::setHoverable);

    m_ddam->Sort(itemsIdList());

    if (delayShow) {
        //hide for delay show
        pItem->setVisible(false);
        //Qt5.3.* not support singleshot with lamda expressions
        QTimer *delayTimer = new QTimer(this);
        connect(delayTimer, &QTimer::timeout, [=] {
            delayTimer->stop();
            delayTimer->deleteLater();

            if (!pItem.isNull())
                item->setVisible(true);
        });
        delayTimer->start(m_addItemDelayInterval);
    }

    relayout();

    //reset state
    m_movingLeftward = true;
}

void DockLayout::moveItem(int from, int to)
{
    int toIndex = to < 0 ? 0 : to;
    m_appList.move(from, toIndex);
    relayout();
}

void DockLayout::removeItem(int index)
{
    m_appList.removeAt(index);
    relayout();
}

void DockLayout::removeItem(AbstractDockItem *item)
{
    int i = indexOf(item);
    if (i != -1){
        m_appList.removeAt(i);
        relayout();
    }
}

void DockLayout::setSpacing(qreal spacing)
{
    m_itemSpacing = spacing;
}

void DockLayout::setVerticalAlignment(Qt::Alignment value)
{
    m_verticalAlignment = value;
}

void DockLayout::setSortDirection(DockLayout::Direction value)
{
    m_sortDirection = value;
}

int DockLayout::indexOf(AbstractDockItem *item) const
{
    if (item)
        return m_appList.indexOf(item);
    else
        return -1;
}

//relative coordinates, not global
int DockLayout::indexOf(int x, int y) const
{
    for (int i = 0; i < m_appList.count(); i ++) {
        if (m_appList.at(i)->geometry().contains(x, y))
            return i;
    }

    return -1;
}

int DockLayout::getContentsWidth()
{
    int tmpWidth = m_appList.count() * m_itemSpacing;
    for (int i = 0; i < m_appList.count(); i ++)
        tmpWidth += m_appList.at(i)->width();

    if (spacingItemIndex() != -1){
        if (!m_dragItemMap.isEmpty() && m_dragItemMap.firstKey())
            tmpWidth += m_dragItemMap.firstKey()->width() + m_itemSpacing;
        else    //spacing add by launcher or desktop item drag enter
            tmpWidth += DockModeData::instance()->getNormalItemWidth() + m_itemSpacing;
    }

    return tmpWidth;
}

int DockLayout::getItemCount() const
{
    return m_appList.count();
}

QList<AbstractDockItem *> DockLayout::getItemList() const
{
    return m_appList;
}

AbstractDockItem *DockLayout::getDraggingItem() const
{
    if (m_dragItemMap.isEmpty())
        return NULL;
    else
        return m_dragItemMap.firstKey();
}

//to recover some damage cause by error
//e.g: item has been drag to some place which can't receive drop event cause item miss
void DockLayout::restoreTmpItem()
{
    if (m_dragItemMap.isEmpty())
        return;

    AbstractDockItem * tmpItem = m_dragItemMap.firstKey();
    m_dragItemMap.remove(tmpItem);
    tmpItem->setVisible(true);

    if (indexOf(tmpItem) == -1)
    {
        if (m_movingLeftward)
            insertItem(tmpItem,m_lastHoverIndex, false);
        else
            insertItem(tmpItem,m_lastHoverIndex + 1, false);
    }

    emit itemDropped();
}

void DockLayout::clearTmpItem()
{
    m_dragItemMap.clear();
}

void DockLayout::relayout()
{
    switch (m_sortDirection)
    {
    case LeftToRight:
        sortLeftToRight();
        break;
    case TopToBottom:
        sortTopToBottom();
        break;
    default:
        break;
    }

    emit contentsWidthChange();
}

//for handle item drop in some area which not accept drop event and respond nothing to the drop event
bool DockLayout::eventFilter(QObject *obj, QEvent *event)
{
    switch (event->type()) {
    case QEvent::MouseMove:
//    case QEvent::Enter:
//    case QEvent::Leave:
        emit itemHoverableChange(true);
        if (m_dragItemMap.isEmpty())
            break;
        restoreTmpItem();
        emit itemDropped();
        break;
    default:
        break;
    }
    return QWidget::eventFilter(obj, event);
}

void DockLayout::dragEnterEvent(QDragEnterEvent *event)
{
    event->setDropAction(Qt::MoveAction);
    event->accept();
}

void DockLayout::dropEvent(QDropEvent *event)
{
    AbstractDockItem *sourceItem = qobject_cast<AbstractDockItem *>(event->source());

    //from launcher
    if (!sourceItem && event->mimeData()->formats().indexOf("RequestDock") != -1){
        QJsonObject dataObj = QJsonDocument::fromJson(event->mimeData()->data("RequestDock")).object();
        if (dataObj.isEmpty() || m_ddam->IsDocked(dataObj.value("appKey").toString()))
            relayout();
        else {
            m_ddam->ReqeustDock(dataObj.value("appKey").toString(), "", "", "");
            emit itemDocking(dataObj.value("appKey").toString());

            qDebug() << "App drop to dock: " << dataObj.value("appKey").toString();
        }
    }
    else {
        //from desktop file
        QList<QUrl> urls = event->mimeData()->urls();
        if (!urls.isEmpty()) {
            for (QUrl url : urls) {
                QString us = url.toString();
                if (us.endsWith(".desktop")) {
                    QString appKey = us.split(QDir::separator()).last();
                    appKey = appKey.mid(0, appKey.length() - 8);
                    if (!m_ddam->IsDocked(appKey)) {
                        m_ddam->ReqeustDock(appKey, "", "", "");
                        emit itemDocking(appKey);

                        qDebug() << "Desktop file drop to dock: " << appKey;
                    }
                }
            }
        }
    }

    restoreTmpItem();
}

void DockLayout::slotItemDrag()
{
    AbstractDockItem *item = qobject_cast<AbstractDockItem*>(sender());

    int tmpIndex = indexOf(item);
    if (tmpIndex != -1)
    {
        m_lastHoverIndex = tmpIndex;
        m_lastPost = QCursor::pos();
        dragoutFromLayout(tmpIndex);

        emit startDrag();
    }
}

void DockLayout::slotItemRelease()
{
    //outside frame,destroy it
    //inside frame,insert it
    AbstractDockItem *item = qobject_cast<AbstractDockItem*>(sender());

    item->setVisible(true);
    if (indexOf(item) == -1)
        insertItem(item,m_lastHoverIndex);
}

void DockLayout::slotItemEntered(QDragEnterEvent *)
{
    //for launcher or desktop item drag enter
    emit startDrag();

    AbstractDockItem *item = qobject_cast<AbstractDockItem*>(sender());

    int tmpIndex = indexOf(item);
    m_lastHoverIndex = tmpIndex;
    if (spacingItemIndex() == -1 && animatingItemCount() <= 0){  //if some animation still running ,there must has spacing item
        addSpacingItem();
        return;
    }

    QPoint tmpPos = QCursor::pos();

    if (tmpPos.x() - m_lastPost.x() == 0)
        return;

    bool lastState = m_movingLeftward;
    switch (m_sortDirection)
    {
    case LeftToRight:
        m_movingLeftward = tmpPos.x() - m_lastPost.x() < 0;
        if (m_movingLeftward != lastState && animatingItemCount() > 0)
        {
            m_movingLeftward = lastState;
            return;
        }
        leftToRightMove(tmpIndex);
        break;
    case TopToBottom:
        //TODO
        topToBottomMove(tmpIndex);
        break;
    }

    m_lastPost = tmpPos;

}

void DockLayout::slotItemExited(QDragLeaveEvent *)
{

}

void DockLayout::slotAnimationFinish()
{
    if (animatingItemCount() > 0){
        //now the animation count should be 0
        //for overlap
        //e.g: spacingIndex is 4 and now if drag item hover item(1) and out of dock suddenly
        //item(1~3) will move to index 4 witch is no longer a spacingItem
        if (animatingItemCount() == 1 && spacingItemIndex() == -1)
            relayout();
    }
}

void DockLayout::sortLeftToRight()
{
    if (m_appList.count() <= 0)
        return;

    for (int i = 0; i < m_appList.count(); i ++)
    {
        AbstractDockItem * toItem = m_appList.at(i);
        toItem->requestAnimationFinish();//make sure the move-animation stop for resort

        int nextX = 0;
        int nextY = 0;
        if (i > 0){
            AbstractDockItem * frontItem = m_appList.at(i - 1);
            nextX = frontItem->x() + frontItem->width() + m_itemSpacing;
        }
        else
            nextX = m_itemSpacing;

        switch (m_verticalAlignment)
        {
        case Qt::AlignTop:
            nextY = 0;
            break;
        case Qt::AlignVCenter:
            nextY = (height() - toItem->height()) / 2;
            break;
        case Qt::AlignBottom:
            nextY = height() - toItem->height();
            break;
        }

        toItem->move(QPoint(nextX, nextY));
        toItem->setNextPos(toItem->pos());
    }
}

void DockLayout::sortTopToBottom()
{

}

void DockLayout::leftToRightMove(int hoverIndex)
{
    int itemWidth = spacingItemWidth();
    int spacintIndex = spacingItemIndex();
    if (spacintIndex == -1)
        return;

    if (spacintIndex > hoverIndex)
    {
        for (int i = hoverIndex; i < spacintIndex; i ++)
        {
            AbstractDockItem *targetItem = m_appList.at(i);
            QPoint nextPos = QPoint(targetItem->x() + itemWidth + m_itemSpacing,0);
            if (targetItem->x() != targetItem->getNextPos().x())    //animation not finish
                break;
            targetItem->moveWithAnimation(nextPos, MOVE_ANIMATION_DURATION_BASE + i * 25);
        }
    }
    else
    {
        for (int i = spacintIndex; i <= hoverIndex; i ++)
        {
            AbstractDockItem *targetItem = m_appList.at(i);
            QPoint nextPos = QPoint(targetItem->x() - itemWidth - m_itemSpacing,0);
            if (targetItem->x() != targetItem->getNextPos().x())    //animation not finish
                break;
            targetItem->moveWithAnimation(nextPos, MOVE_ANIMATION_DURATION_BASE + i * 25);
        }
    }
}

void DockLayout::topToBottomMove(int hoverIndex)
{
    Q_UNUSED(hoverIndex)
}

void DockLayout::addSpacingItem()
{
    if (spacingItemIndex() != -1 || animatingItemCount() > 0)
        return;

    int spacingValue = 0;
    if (m_dragItemMap.isEmpty())
        spacingValue = DockModeData::instance()->getNormalItemWidth();
    else
        spacingValue = m_dragItemMap.firstKey()->width();

    for (int i = m_appList.count() -1;i >= m_lastHoverIndex; i-- )
    {
        AbstractDockItem *targetItem = m_appList.at(i);
        targetItem->moveWithAnimation(QPoint(targetItem->x() + spacingValue + m_itemSpacing,0),
                                      MOVE_ANIMATION_DURATION_BASE + i * 25);
    }

    emit contentsWidthChange();
}

void DockLayout::removeSpacingItem()
{
    if (spacingItemIndex() == -1)
        return;
    if (animatingItemCount() > 0){//try to remove spacing again
        QTimer::singleShot(100, this, SLOT(removeSpacingItem()));
        return;
    }

    int spacingValue = 0;
    if (m_dragItemMap.isEmpty())
        spacingValue = DockModeData::instance()->getNormalItemWidth();
    else
        spacingValue = m_dragItemMap.firstKey()->width();

    for (int i = spacingItemIndex(); i < m_appList.count(); i ++)
    {
        AbstractDockItem *targetItem = m_appList.at(i);
        targetItem->moveWithAnimation(QPoint(targetItem->x() - spacingValue - m_itemSpacing, 0),
                                      MOVE_ANIMATION_DURATION_BASE + i * 25);
    }

    //emit the width change signal after the last animation is finish
    QTimer::singleShot(MOVE_ANIMATION_DURATION_BASE + m_appList.count() * 25, this, SIGNAL(contentsWidthChange()));
}

void DockLayout::dragoutFromLayout(int index)
{
    AbstractDockItem * tmpItem = m_appList.takeAt(index);
    tmpItem->setVisible(false);
    m_dragItemMap.insert(tmpItem,index);
}

int DockLayout::spacingItemWidth() const
{
    if (m_dragItemMap.isEmpty())
        return DockModeData::instance()->getNormalItemWidth();
    else
        return m_dragItemMap.firstKey()->width();
}

int DockLayout::spacingItemIndex() const
{
    if (m_sortDirection == TopToBottom)
        return -1;
    if (m_appList.count() <= 1)
        return -1;
    if (m_appList.at(0)->getNextPos().x() > m_itemSpacing)
        return 0;

    for (int i = 1; i < m_appList.count(); i ++)
    {
        if (m_appList.at(i)->getNextPos().x() - m_itemSpacing != m_appList.at(i - 1)->getNextPos().x() + m_appList.at(i - 1)->width())
            return i;
    }

    return -1;
}

int DockLayout::animatingItemCount()
{
    int tmpCount = 0;
    foreach (AbstractDockItem *item, m_appList) {
        if (item->pos() != item->getNextPos())
            tmpCount ++;
    }
    return tmpCount;
}


//return the docked app id list,just for app
QStringList DockLayout::itemsIdList() const
{
    QStringList dockedAppList = m_ddam->DockedAppList().value();

    QStringList idList;
    foreach (AbstractDockItem *item, m_appList) {
        QString itemId = item->getItemId();
        if (!itemId.isEmpty() && dockedAppList.indexOf(itemId) != -1)
            idList << itemId;
    }
    return idList;
}
int DockLayout::removeItemDelayInterval() const
{
    return m_removeItemDelayInterval;
}

void DockLayout::setRemoveItemDelayInterval(int removeItemDelayInterval)
{
    m_removeItemDelayInterval = removeItemDelayInterval;
}

int DockLayout::addItemDelayInterval() const
{
    return m_addItemDelayInterval;
}

void DockLayout::setaddItemDelayInterval(int addItemDelayInterval)
{
    m_addItemDelayInterval = addItemDelayInterval;
}

