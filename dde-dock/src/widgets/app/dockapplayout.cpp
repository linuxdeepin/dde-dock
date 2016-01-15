#include <QDrag>
#include "dockapplayout.h"
#include "../../controller/dockmodedata.h"

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
            connect(group, &QPropertyAnimation::finished, [=]{
                emit droped();
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

DockAppLayout::DockAppLayout(QWidget *parent) :
    MovableLayout(parent), m_isDraging(false)
{
    initDropMask();
    initAppManager();

    qApp->installEventFilter(this);
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

void DockAppLayout::initEntries() const
{
    m_appManager->initEntries();
}

bool DockAppLayout::eventFilter(QObject *obj, QEvent *e)
{
    if (e->type() == QEvent::Move) {
        QMoveEvent *me = (QMoveEvent *)e;
        if (me && isDraging() && !geometry().contains(mapFromGlobal(QCursor::pos()))) {
            //show mask to catch draging widget
            //fixme
            m_mask->move(QCursor::pos().x() - 15, QCursor::pos().y() - 15); //15,拖动时的鼠标位移
            m_mask->show();
        }
    }

    return QWidget::eventFilter(obj, e);
}

void DockAppLayout::initDropMask()
{
    m_mask = new DropMask;
    connect(m_mask, &DropMask::droped, this, [=] {
        setIsDraging(false);
        emit requestSpacingItemsDestroy();
    });
    connect(m_mask, &DropMask::invalidDroped, this, &DockAppLayout::restoreDragingWidget);
    connect(this, &DockAppLayout::dragEntered, m_mask, &DropMask::hide);
    connect(this, &DockAppLayout::startDrag, this, [=](QDrag* drag) {
        setIsDraging(true);

        if (DockModeData::instance()->getDockMode() == Dock::FashionMode) {
            DockAppItem *item = qobject_cast<DockAppItem *>(dragingWidget());
            if (item) {
                drag->setPixmap(item->iconPixmap());
            }
        }

        emit itemHoverableChange(false);
    });
}

void DockAppLayout::onDrop(QDropEvent *event)
{
    m_mask ->hide();
    setIsDraging(false);

    if (event->source() == this) {  //from itself
        m_ddam->Sort(appIds());
        event->accept();
    }
    else if (event->mimeData()->formats().indexOf("RequestDock") != -1){    //from launcher
        QJsonObject dataObj = QJsonDocument::fromJson(event->mimeData()->data("RequestDock")).object();
        if (dataObj.isEmpty() || m_ddam->IsDocked(dataObj.value("appKey").toString())) {
            emit requestSpacingItemsDestroy();
        }
        else {
            m_ddam->ReqeustDock(dataObj.value("appKey").toString(), "", "", "");
            m_appManager->setDockingItemId(dataObj.value("appKey").toString());

            qDebug() << "App drop to dock: " << dataObj.value("appKey").toString();
        }
    }
    else {  //from desktop file
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
            removeWidget(item);
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
    connect(this, &DockAppLayout::itemHoverableChange, item, &DockAppItem::setHoverable);
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
    connect(this, &DockAppLayout::itemHoverableChange, item, &DockAppItem::setHoverable);
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

bool DockAppLayout::isDraging() const
{
    return m_isDraging;
}

void DockAppLayout::setIsDraging(bool isDraging)
{
    m_isDraging = isDraging;

    emit itemHoverableChange(!isDraging);
}

