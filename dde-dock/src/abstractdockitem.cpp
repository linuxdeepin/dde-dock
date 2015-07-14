#include <QWidget>
#include <QFrame>
#include <QLabel>
#include "Widgets/arrowrectangle.h"

#include "abstractdockitem.h"

AbstractDockItem::AbstractDockItem(QWidget * parent) :
    QFrame(parent)
{

}

AbstractDockItem::~AbstractDockItem()
{

}

QString AbstractDockItem::getTitle()
{
    return "";
}

QWidget * AbstractDockItem::getApplet()
{
    return NULL;
}

bool AbstractDockItem::moveable()
{
    return m_moveable;
}

bool AbstractDockItem::actived()
{
    return m_isActived;
}

void AbstractDockItem::resize(int width,int height){
    QFrame::resize(width,height);

    emit widthChanged();
}

void AbstractDockItem::resize(const QSize &size){
    QFrame::resize(size);

    emit widthChanged();
}

QPoint AbstractDockItem::getNextPos()
{
    return m_itemNextPos;
}

void AbstractDockItem::setNextPos(const QPoint &value)
{
    m_itemNextPos = value;
}

void AbstractDockItem::setNextPos(int x, int y)
{
    m_itemNextPos.setX(x); m_itemNextPos.setY(y);
}

int AbstractDockItem::globalX()
{
    return mapToGlobal(QPoint(0,0)).x();
}

int AbstractDockItem::globalY()
{
    return mapToGlobal(QPoint(0,0)).y();
}

QPoint AbstractDockItem::globalPos()
{
    return mapToGlobal(QPoint(0,0));
}

void AbstractDockItem::showPreview()
{
    if (!m_previewAR->isHidden())
    {
        m_previewAR->resizeWithContent();
        return;
    }
    QWidget *tmpContent = getApplet();
    if (tmpContent == NULL) {
        QString title = getTitle();
        // TODO: memory management
        tmpContent = new QLabel(title);
        tmpContent->setStyleSheet("QLabel { color: white }");
        tmpContent->setFixedSize(100, 20);
    }

    m_previewAR->setArrorDirection(ArrowRectangle::ArrowBottom);
    m_previewAR->setContent(tmpContent);
    m_previewAR->showAtBottom(globalX() + width() / 2,globalY() - 5);
}

void AbstractDockItem::hidePreview(int interval)
{
    m_previewAR->delayHide(interval);
}

void AbstractDockItem::cancelHide()
{
    m_previewAR->cancelHide();
}

void AbstractDockItem::resizePreview()
{
    m_previewAR->resizeWithContent();
    m_previewAR->showAtBottom(globalX() + width() / 2,globalY() - 5);
}
