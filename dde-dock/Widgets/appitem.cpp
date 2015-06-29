#include "appitem.h"

AppItem::AppItem(QWidget *parent) :
    DockItem(parent)
{
    setParent(parent);

    initBackground();
    setAcceptDrops(true);
}

AppItem::AppItem(QString title, QWidget *parent):
    DockItem(parent)
{
    this->setParent(parent);
    this->itemTitle = title;

    this->initBackground();
}

AppItem::AppItem(QString title, QString iconPath, QWidget *parent) :
    DockItem(parent)
{
    this->setParent(parent);
    this->itemTitle = title;
    this->itemIconPath = iconPath;

    this->initBackground();
    this->setIcon(itemIconPath);
}

void AppItem::resize(const QSize &size)
{
    QFrame::resize(size);
    resizeResources();
}

void AppItem::resize(int width, int height)
{
    QFrame::resize(width,height);
    resizeResources();
}

QPoint AppItem::getNextPos()
{
    return this->nextPos;
}

void AppItem::setNextPos(const QPoint &value)
{
    this->nextPos = value;
}

void AppItem::setNextPos(int x, int y)
{
    this->nextPos.setX(x);
    this->nextPos.setY(y);
}

void AppItem::resizeResources()
{
    if (appIcon != NULL)
    {
        appIcon->resize(DockConstants::getInstants()->getIconSize(),DockConstants::getInstants()->getIconSize());
        appIcon->move(this->width() / 2 - appIcon->width() / 2, this->height() / 2 - appIcon->height() / 2);
    }

    if (appBackground != NULL)
    {
        appBackground->resize(this->width(),this->height());
        appBackground->move(0,0);
    }
}

void AppItem::initBackground()
{
    appBackground = new AppBackground(this);
    appBackground->resize(this->width(),this->height());
    appBackground->move(0,0);
}

void AppItem::mousePressEvent(QMouseEvent * event)
{
    //qWarning() << "mouse press...";
    emit mousePress(event->globalX(), event->globalY(),this);
}

void AppItem::mouseReleaseEvent(QMouseEvent * event)
{
//    qWarning() << "mouse release...";
    emit mouseRelease(event->globalX(), event->globalY(),this);
}

void AppItem::mouseDoubleClickEvent(QMouseEvent * event)
{
    emit mouseDoubleClick(this);
}

void AppItem::mouseMoveEvent(QMouseEvent *event)
{
    //this event will only execp onec then handle by Drag
    emit dragStart(this);

    Qt::MouseButtons btn = event->buttons();
    if(btn == Qt::LeftButton)
    {
        QDrag* drag = new QDrag(this);
        QMimeData* data = new QMimeData();
        QImage dataImg(this->itemIconPath);
        data->setImageData(QVariant(dataImg));
        drag->setMimeData(data);

        QPixmap pixmap(this->itemIconPath);
        drag->setPixmap(pixmap);

        drag->setHotSpot(QPoint(15,15));

        drag->exec(Qt::CopyAction | Qt::MoveAction, Qt::MoveAction);
    }
}

void AppItem::enterEvent(QEvent *event)
{
    emit mouseEntered(this);
}

void AppItem::leaveEvent(QEvent *event)
{
    emit mouseExited(this);
}

void AppItem::dragEnterEvent(QDragEnterEvent *event)
{
    emit dragEntered(event,this);

    AppItem *tmpItem = NULL;
    tmpItem = dynamic_cast<AppItem *>(event->source());
    if (tmpItem)
    {
//        qWarning()<< "[Info:]" << "Brother Item.";
    }
    else
    {
        event->setDropAction(Qt::MoveAction);
        event->accept();
    }
}

void AppItem::dragLeaveEvent(QDragLeaveEvent *event)
{
    emit dragExited(event,this);
}

void AppItem::dropEvent(QDropEvent *event)
{
    qWarning() << "Item get drop:" << event->pos();
    emit drop(event,this);
}

AppItem::~AppItem()
{

}

