#include "appitem.h"

AppItem::AppItem(QWidget *parent) :
    QFrame(parent)
{
    this->setParent(parent);

    this->initBackground();
}

AppItem::AppItem(QString title, QWidget *parent):
    QFrame(parent)
{
    this->setParent(parent);
    this->itemTitle = title;

    this->initBackground();
}

AppItem::AppItem(QString title, QString iconPath, QWidget *parent) :
    QFrame(parent)
{
    this->setParent(parent);
    this->itemTitle = title;
    this->itemIconPath = iconPath;

    this->initBackground();
    this->setIcon(itemIconPath);
}

void AppItem::setTitle(const QString &title)
{
    this->itemTitle = title;
}

void AppItem::setIcon(const QString &iconPath, int size)
{
    appIcon = new AppIcon(iconPath,this);
    appIcon->resize(size,size);
//    appIcon->setIcon(iconPath);
    appIcon->move(this->width() / 2, this->height() / 2);
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

void AppItem::setMoveable(bool value)
{
    this->itemMoveable = value;
}

bool AppItem::getMoveable()
{
    return this->itemMoveable;
}

void AppItem::setIndex(int value)
{
    this->itemIndex = value;
}

int AppItem::getIndex()
{
    return this->itemIndex;
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
    qWarning() << "press...";
    emit mousePress(event->globalX(), event->globalY(),this);
}

void AppItem::mouseReleaseEvent(QMouseEvent * event)
{
    emit mouseRelease(event->globalX(), event->globalY(),this);
}

void AppItem::mouseMoveEvent(QMouseEvent * event)
{
    emit mouseMove(event->globalX(), event->globalY(),this);
}

void AppItem::mouseDoubleClickEvent(QMouseEvent * event)
{
    emit mouseDoubleClick(this);
}

AppItem::~AppItem()
{

}

