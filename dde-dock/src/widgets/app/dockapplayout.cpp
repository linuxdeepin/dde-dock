/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QDrag>
#include <QtConcurrent>
#include "dockapplayout.h"
#include "../../controller/dockmodedata.h"

/*
class DropMask : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(int rValue READ getRValue WRITE setRValue)
    Q_PROPERTY(double sValue  READ getSValue WRITE setSValue)
    int m_rValue;

    double m_sValue;

public:
    DropMask(QWidget *parent = 0);

    int getRValue() const {return m_rValue;}
    double getSValue() const {return m_sValue;}

public slots:
    void setRValue(int rValue);
    void setSValue(double sValue);

signals:
    void droped();
    void invalidDroped();

protected:
    void dropEvent(QDropEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
};

DropMask::DropMask(QWidget *parent) :
    QLabel(parent)
{
    setAcceptDrops(true);
    setWindowFlags(Qt::ToolTip);
    setAttribute(Qt::WA_TranslucentBackground);
    setFixedWidth(DockModeData::instance()->getAppIconSize());
    setFixedHeight(DockModeData::instance()->getAppIconSize());
}

void DropMask::setRValue(int rValue)
{
    if (!pixmap())
        return;
    QTransform rt;
    rt.translate(width() / 2, height() / 2);
    rt.rotate(rValue);
    rt.translate(-width() / 2, -height() / 2);
    setPixmap(pixmap()->transformed(rt));
    m_rValue = rValue;
}

void DropMask::setSValue(double sValue)
{
    if (!pixmap())
        return;
    QTransform st(1, 0, 0, 1, width()/2, height()/2);
    st.scale(sValue, sValue);
    st.rotate(90);//TODO work around here
    setPixmap(pixmap()->transformed(st));
    m_sValue = sValue;
}

void DropMask::dropEvent(QDropEvent *e)
{
    DockAppLayout *layout = dynamic_cast<DockAppLayout *>(e->source());
    if (!layout)
        return;
    DockAppItem *item = qobject_cast<DockAppItem *>(layout->dragingWidget());
    if (item)
    {
        //restore item to dock if item is actived
        if (item->itemData().isActived) {
            emit invalidDroped();
            return;
        }

        DBusDockedAppManager dda;
        if (dda.IsDocked(item->itemData().id).value()) {
            dda.RequestUndock(item->itemData().id);
        }

        qDebug() << "Item drop to mask:" << e->mimeData()->hasImage();
        QImage image = qvariant_cast<QImage>(e->mimeData()->imageData());
        if (!image.isNull()) {
            setPixmap(QPixmap::fromImage(image).scaled(size()));

            QPropertyAnimation *scaleAnimation = new QPropertyAnimation(this, "sValue");
            scaleAnimation->setDuration(1000);
            scaleAnimation->setStartValue(1);
            scaleAnimation->setEndValue(0.3);

            QPropertyAnimation *rotationAnimation = new QPropertyAnimation(this, "rValue");
            rotationAnimation->setDuration(1000);
            rotationAnimation->setStartValue(0);
            rotationAnimation->setEndValue(360);

            QParallelAnimationGroup * group = new QParallelAnimationGroup();
            group->addAnimation(scaleAnimation);
//            group->addAnimation(rotationAnimation);

            group->start();
            emit droped();

            connect(group, &QPropertyAnimation::finished, [=]{
                hide();

                scaleAnimation->deleteLater();
                rotationAnimation->deleteLater();
                group->deleteLater();
            });
        }
        else {
            qWarning() << "Item drop to mask, Image is NULL!";
        }
    }


}

void DropMask::dragEnterEvent(QDragEnterEvent *e)
{
    e->accept();
}

#include "dockapplayout.moc"

/////////////////////////////////////////////////////////////////////////////////////////////////
*/

DockAppLayout::DockAppLayout(QWidget *parent) :
    MovableLayout(parent), m_isDraging(false)
{
//    initDropMask();
    initAppManager();

//    qApp->installEventFilter(this);
    m_ddam = new DBusDockedAppManager(this);
    connect(this, &DockAppLayout::drop, this, &DockAppLayout::onDrop);
    connect(this, &DockAppLayout::dragLeaved, this, &DockAppLayout::onDragLeave);
    connect(this, &DockAppLayout::dragEntered, this, &DockAppLayout::onDragEnter);

    m_dem = new DBusEntryManager(this);
}

QSize DockAppLayout::sizeHint() const
{
    QSize size;
    int w = 0;
    int h = 0;
    switch (direction()) {
    case QBoxLayout::LeftToRight:
    case QBoxLayout::RightToLeft:
        size.setHeight(DockModeData::instance()->getItemHeight());
        for (QWidget * widget : widgets()) {
            w += widget->width();
        }
        if (dragingWidget() && isDraging()) {
            w += dragingWidget()->width();
        }
        else if (getDragFromOutside()){
            w += DockModeData::instance()->getNormalItemWidth();
        }
        size.setWidth(w + getLayoutSpacing() * widgets().count());
        break;
    case QBoxLayout::TopToBottom:
    case QBoxLayout::BottomToTop:
        size.setWidth(DockModeData::instance()->getNormalItemWidth());
        for (QWidget * widget : widgets()) {
            h += widget->height();
        }
        if (dragingWidget()) {
            h += dragingWidget()->height();
        }
        else if (getDragFromOutside()) {
            h += DockModeData::instance()->getItemHeight();
        }
        size.setHeight(h + getLayoutSpacing() * widgets().count());
        break;
    }

    return size;
}

void DockAppLayout::initEntries() const
{
    m_appManager->initEntries();
}

void DockAppLayout::updateWindowIconGeometries()
{
//    qDebug() << "update window icon geometries.";

    for (QWidget *w : widgets()) {
        DockAppItem * item = qobject_cast<DockAppItem *>(w);
        if (item) {
            item->setWindowIconGeometries();
        }
    }
}

void DockAppLayout::enterEvent(QEvent *e)
{
    Q_UNUSED(e)

    setDragable(true);
}

void DockAppLayout::resizeEvent(QResizeEvent *e)
{
    MovableLayout::resizeEvent(e);
    updateItemWidths();
}

//bool DockAppLayout::eventFilter(QObject *obj, QEvent *e)
//{
//    if (e->type() == QEvent::Move) {
//        QMoveEvent *me = (QMoveEvent *)e;
//        QRect r(0, 0, width(), height());
//        if (me && isDraging() && !r.contains(mapFromGlobal(QCursor::pos()))) {
//            //show mask to catch draging widget
//            //fixme
//            m_mask->move(QCursor::pos().x() - 15, QCursor::pos().y() - 15); //15,拖动时的鼠标位移
//            m_mask->show();
//        }
//    }

//    return QWidget::eventFilter(obj, e);
//}

//void DockAppLayout::initDropMask()
//{
//    m_mask = new DropMask;
//    connect(m_mask, &DropMask::droped, this, [=] {
//        setIsDraging(false);
//        emit requestSpacingItemsDestroy(false);
//        setFixedSize(sizeHint());
//    });
//    connect(m_mask, &DropMask::invalidDroped, this, &DockAppLayout::restoreDragingWidget);
//    connect(this, &DockAppLayout::dragEntered, m_mask, &DropMask::hide);
//    connect(this, &DockAppLayout::startDrag, this, [=](QDrag* drag) {
//        setIsDraging(true);

//        if (DockModeData::instance()->getDockMode() == Dock::FashionMode) {
//            DockAppItem *item = qobject_cast<DockAppItem *>(dragingWidget());
//            if (item) {
//                drag->setPixmap(item->iconPixmap());
//            }
//        }

//        emit itemHoverableChange(false);
//    });
//}

void DockAppLayout::onDrop(QDropEvent *event)
{
//    m_mask ->hide();
    setIsDraging(false);
    setDragFromOutside(false);

    if (event->source() == this) {  //from itself
        m_dem->Reorder(appIds());
        event->accept();
    } else if (event->mimeData()->formats().indexOf("RequestDock") != -1){    //from launcher
        QJsonObject dataObj = QJsonDocument::fromJson(event->mimeData()->data("RequestDock")).object();
        if (dataObj.isEmpty() || m_ddam->IsDocked(dataObj.value("appKey").toString())) {
            emit requestSpacingItemsDestroy(true);
        }
        else {
            m_ddam->RequestDock(dataObj.value("appKey").toString(), "", "", "");
            m_appManager->setDockingItemId(dataObj.value("appKey").toString());

            qDebug() << "App drop to dock: " << dataObj.value("appKey").toString();
        }
    } else {  //from desktop file
        QList<QUrl> urls = event->mimeData()->urls();
        if (!urls.isEmpty()) {
            QStringList normals;
            QStringList desktops;
            separateFiles(urls, normals, desktops);
            if (desktops.length() > 0) {
                for (QString path : desktops) {
                    if (!isDesktopFileDocked(path)) {
                        QString appKey = getAppKeyByPath(path);
                        m_ddam->RequestDock(appKey, "", "", "");
                        m_appManager->setDockingItemId(appKey);

                        qDebug() << "Desktop file drop to dock: " << appKey;
                    }
                }
            } else {
                //just normal files, try to open files by the target app
                int index = getHoverIndextByPos(mapFromGlobal(QCursor::pos()));
                if (index != -1) {
                    DockAppItem *item = qobject_cast<DockAppItem *>(widget(index));
                    item->openFiles(normals);
                }
            }
        }
    }

    // interval 0 stands for timeout will be triggered on idle.
    QTimer::singleShot(0, this, &DockAppLayout::updateWindowIconGeometries);
}

void DockAppLayout::onDragLeave(QDragLeaveEvent *event)
{
    Q_UNUSED(event)

    setDragFromOutside(false);
}

void DockAppLayout::onDragEnter(QDragEnterEvent *event)
{
    if (event->source() == this) {
        setDragFromOutside(true);
        return;
    }
    else if (event->mimeData()->formats().indexOf("RequestDock") != -1) {
        QJsonObject dataObj = QJsonDocument::fromJson(event->mimeData()->data("RequestDock")).object();
        if (dataObj.isEmpty() || m_ddam->IsDocked(dataObj.value("appKey").toString())) {
            setDragable(false);
            setDragFromOutside(false);
            emit requestSpacingItemsDestroy(false);
        }
        else {
            setDragFromOutside(true);
        }
    }
    else {  //from desktop file
        QList<QUrl> urls = event->mimeData()->urls();
        if (!urls.isEmpty()) {
            QStringList normals;
            QStringList desktops;
            separateFiles(urls, normals, desktops);
            if (desktops.length() > 0) {
                for (QString path : desktops) {
                    //多个文件中只要存在一个有效并且未docked的desktop文件，都可以做拖入dock的操作
                    if (!isDesktopFileDocked(path)) {
                        setDragFromOutside(true);
                        return;
                    }
                }
            }

            setDragable(false);
            setDragFromOutside(false);
            emit requestSpacingItemsDestroy(false);
        }
        else {
            setDragFromOutside(false);
        }
    }
}

void DockAppLayout::initAppManager()
{
    m_appManager = new DockAppManager(this);
    connect(m_appManager, &DockAppManager::entryAdded, this, &DockAppLayout::onAppItemAdd);
    connect(m_appManager, &DockAppManager::entryAppend, this, &DockAppLayout::onAppAppend);
    connect(m_appManager, &DockAppManager::entryRemoved, this, &DockAppLayout::onAppItemRemove);
    connect(m_appManager, &DockAppManager::requestSort, this, [=] {
        m_dem->Reorder(appIds());
    });
}

void DockAppLayout::onAppItemRemove(const QString &id)
{
    QList<DockItem *> tmpList = this->widgets();
    for (DockItem * item : tmpList) {
        DockAppItem *tmpItem = qobject_cast<DockAppItem *>(item);
        if (tmpItem && tmpItem->getItemId() == id) {
            removeWidget(item);
            tmpItem->setVisible(false);
            tmpItem->deleteLater();
            return;
        }
    }
    updateItemWidths();
}

void DockAppLayout::onAppItemAdd(DockAppItem *item)
{
    const int index = hoverIndex();

    if (index == -1)
        return onAppAppend(item);

    insertWidget(index, item);
    createConnections(item);
    updateItemWidths();
}

void DockAppLayout::onAppAppend(DockAppItem *item)
{


    addWidget(item);
    createConnections(item);
    updateItemWidths();
}

QStringList DockAppLayout::appIds()
{
    QStringList ids;
    for (QWidget *w : widgets()) {
        DockAppItem * item = qobject_cast<DockAppItem *>(w);
        if (item) {
            ids << item->getItemId();
        }
    }

    return ids;
}

void DockAppLayout::createConnections(DockAppItem *item)
{
    connect(item, &DockAppItem::needPreviewShow, this, [=](QPoint pos) {
        DockAppItem * s = qobject_cast<DockAppItem *>(sender());
        if (s) {
            emit needPreviewShow(s, pos);
        }
    });
    connect(item, &DockAppItem::activatedChanged, this, [this](bool) {
        updateItemWidths();
    });
    connect(item, &DockAppItem::needPreviewHide, this, &DockAppLayout::needPreviewHide);
    connect(item, &DockAppItem::needPreviewUpdate, this, &DockAppLayout::needPreviewUpdate);
    connect(this, &DockAppLayout::itemHoverableChange, item, &DockAppItem::setHoverable);
}

// update the width of all items according to the spaces that left.
void DockAppLayout::updateItemWidths()
{
    DockModeData *dmd = DockModeData::instance();

    if (dmd->getDockMode() == Dock::ClassicMode) {
        QList<DockAppItem *> activeItems;
        QList<DockAppItem *> inactiveItems;

        QList<DockItem *> tmpList = this->widgets();
        for (DockItem * item : tmpList) {
            DockAppItem *tmpItem = qobject_cast<DockAppItem *>(item);
            if (tmpItem->actived()) {
                activeItems << tmpItem;
            } else {
                inactiveItems << tmpItem;
            }
        }

        int itemsCount = activeItems.length() + inactiveItems.length();
        bool cannotHoldAllAsInactives = itemsCount * dmd->getNormalItemWidth() + itemsCount * getLayoutSpacing() > this->width();

        if (cannotHoldAllAsInactives || activeItems.length() == 0) {
            // TODO: should scale the items or something.
            for (DockAppItem *item : activeItems) {
                item->setFixedSize(dmd->getActivedItemWidth(), dmd->getItemHeight());
                item->refreshUI();
            }

            for (DockAppItem *item: inactiveItems) {
                item->setFixedSize(dmd->getNormalItemWidth(), dmd->getItemHeight());
                item->refreshUI();
            }
        } else {
            uint normalItemWidth = dmd->getNormalItemWidth();
            int widthLeft = this->width() - inactiveItems.length() * normalItemWidth - itemsCount * getLayoutSpacing();
            int averageItemWidthLeft = widthLeft / activeItems.length();
            uint activeItemWidth = qMin(dmd->getActivedItemWidth(), averageItemWidthLeft);

            for (DockAppItem *item : activeItems) {
                item->setFixedSize(activeItemWidth, dmd->getItemHeight());
                item->refreshUI();
            }

            for (DockAppItem *item: inactiveItems) {
                item->setFixedSize(normalItemWidth, dmd->getItemHeight());
                item->refreshUI();
            }
        }
    } else {
        const int itemCount = widgets().size();
        const bool canHold = itemCount * dmd->getActivedItemWidth() < width();
        const QList<DockItem *> tmpList = this->widgets();

        const int itemWidth = canHold ? dmd->getActivedItemWidth() : (width() - itemCount * dmd->getAppItemSpacing()) / itemCount;

        for (DockItem * item : tmpList) {
            DockAppItem *tmpItem = qobject_cast<DockAppItem *>(item);
            if (tmpItem) {
                if (tmpItem->actived()) {
                    tmpItem->setFixedSize(itemWidth, dmd->getDockHeight());
                } else {
                    tmpItem->setFixedSize(itemWidth, dmd->getDockHeight());
                }

                tmpItem->refreshUI();
            }
        }
    }
}

bool DockAppLayout::getDragFromOutside() const
{
    return m_dragFromOutside;
}

void DockAppLayout::setDragFromOutside(bool dragFromOutside)
{
    m_dragFromOutside = dragFromOutside;
}

bool DockAppLayout::isDraging() const
{
    return m_isDraging;
}

void DockAppLayout::setIsDraging(bool isDraging)
{
    m_isDraging = isDraging;

    emit itemHoverableChange(!isDraging);
}

bool DockAppLayout::isDesktopFileDocked(const QString &path)
{
    QSettings ds(path, QSettings::IniFormat);
    ds.setIniCodec(QTextCodec::codecForName("UTF-8"));
    ds.beginGroup("Desktop Entry");
    QString appKey = ds.value("X-Deepin-AppID").toString();
    ds.endGroup();

    return m_ddam->IsDocked(appKey).value();
}

QString DockAppLayout::getAppKeyByPath(const QString &path)
{
    QSettings ds(path, QSettings::IniFormat);
    ds.setIniCodec(QTextCodec::codecForName("UTF-8"));
    ds.beginGroup("Desktop Entry");
    QString appKey = ds.value("X-Deepin-AppID").toString();
    ds.endGroup();

    return appKey;
}

void DockAppLayout::separateFiles(const QList<QUrl> &urls, QStringList &normals, QStringList &desktopes)
{
    for (QUrl url : urls) {
        if (url.fileName().endsWith(".desktop")) {
            QSettings ds(url.path(), QSettings::IniFormat);
            ds.setIniCodec(QTextCodec::codecForName("UTF-8"));
            ds.beginGroup("Desktop Entry");
            QString appKey = ds.value("X-Deepin-AppID").toString();
            if (!appKey.isEmpty()) {
                desktopes << url.path();
            }
            else {
                normals << url.path();
            }
            ds.endGroup();
        }
        else {
            normals << url.path();
        }
    }
}

