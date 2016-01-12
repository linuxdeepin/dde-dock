#include "dockapplayout.h"
#include "../../controller/dockmodedata.h"

DockAppLayout::DockAppLayout(QWidget *parent) : MovableLayout(parent)
{
    initAppManager();

    m_ddam = new DBusDockedAppManager(this);

    connect(this, &DockAppLayout::drop, this, &DockAppLayout::onDrop);
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
        size.setWidth(w + getLayoutSpacing() * widgets().count());
        break;
    case QBoxLayout::TopToBottom:
    case QBoxLayout::BottomToTop:
        size.setWidth(DockModeData::instance()->getNormalItemWidth());
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

void DockAppLayout::onDrop(QDropEvent *event)
{
    //form itself
    if (event->source() == this) {
        m_ddam->Sort(appIds());
        event->accept();
    }
    //from launcher
    else if (event->mimeData()->formats().indexOf("RequestDock") != -1){
        QJsonObject dataObj = QJsonDocument::fromJson(event->mimeData()->data("RequestDock")).object();
        if (dataObj.isEmpty() || m_ddam->IsDocked(dataObj.value("appKey").toString()))
            emit spacingItemAdded();
        else {
            m_ddam->ReqeustDock(dataObj.value("appKey").toString(), "", "", "");
            m_appManager->setDockingItemId(dataObj.value("appKey").toString());

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
                        m_appManager->setDockingItemId(appKey);

                        qDebug() << "Desktop file drop to dock: " << appKey;
                    }
                }
            }
        }
    }

}

void DockAppLayout::initAppManager()
{
    m_appManager = new DockAppManager(this);
    connect(m_appManager, &DockAppManager::entryAdded, this, &DockAppLayout::onAppItemAdd);
    connect(m_appManager, &DockAppManager::entryAppend, this, &DockAppLayout::onAppAppend);
    connect(m_appManager, &DockAppManager::entryRemoved, this, &DockAppLayout::onAppItemRemove);
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

void DockAppLayout::onAppItemAdd(DockAppItem *item)
{
    insertWidget(hoverIndex(), item);
    connect(item, &DockAppItem::needPreviewShow, this, [=](QPoint pos) {
        DockAppItem * s = qobject_cast<DockAppItem *>(sender());
        if (s) {
            emit needPreviewShow(s, pos);
        }
    });
    connect(item, &DockAppItem::needPreviewHide, this, &DockAppLayout::needPreviewHide);
    connect(item, &DockAppItem::needPreviewUpdate, this, &DockAppLayout::needPreviewUpdate);
}

void DockAppLayout::onAppAppend(DockAppItem *item)
{
    addWidget(item);
    connect(item, &DockAppItem::needPreviewShow, this, [=](QPoint pos) {
        DockAppItem * s = qobject_cast<DockAppItem *>(sender());
        if (s) {
            emit needPreviewShow(s, pos);
        }
    });
    connect(item, &DockAppItem::needPreviewHide, this, &DockAppLayout::needPreviewHide);
    connect(item, &DockAppItem::needPreviewUpdate, this, &DockAppLayout::needPreviewUpdate);
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

