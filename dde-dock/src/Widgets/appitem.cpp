#include "appitem.h"

AppItem::AppItem(QWidget *parent) :
    AbstractDockItem(parent)
{
    setAcceptDrops(true);
    resize(dockCons->getNormalItemWidth(), dockCons->getItemHeight());
    initBackground();
    connect(dockCons, &DockConstants::dockModeChanged,this, &AppItem::slotDockModeChanged);
}

AppItem::AppItem(QString title, QWidget *parent):
    AbstractDockItem(parent)
{
    m_itemTitle = title;

    setAcceptDrops(true);
    resize(dockCons->getNormalItemWidth(), dockCons->getItemHeight());
    initBackground();
    connect(dockCons, &DockConstants::dockModeChanged,this, &AppItem::slotDockModeChanged);
}

AppItem::AppItem(QString title, QString iconPath, QWidget *parent) :
    AbstractDockItem(parent)
{
    m_itemTitle = title;
    m_itemIconPath = iconPath;

    setAcceptDrops(true);
    resize(dockCons->getNormalItemWidth(), dockCons->getItemHeight());
    initBackground();
    setIcon(m_itemIconPath,dockCons->getAppIconSize());
    connect(dockCons, &DockConstants::dockModeChanged,this, &AppItem::slotDockModeChanged);
}

void AppItem::setIcon(const QString &iconPath, int size)
{
    m_appIcon = new AppIcon(iconPath, this);
    m_appIcon->resize(size, size);

    reanchorIcon();
}

void AppItem::setActived(bool value)
{
    m_isActived = value;
    if (!value)
        resize(dockCons->getNormalItemWidth(), dockCons->getItemHeight());
    else
        resize(dockCons->getActivedItemWidth(), dockCons->getItemHeight());
}

void AppItem::setCurrentOpened(bool value)
{
    m_isCurrentOpened = value;
}

bool AppItem::currentOpened()
{
    return m_isCurrentOpened;
}

void AppItem::slotDockModeChanged(DockConstants::DockMode newMode, DockConstants::DockMode oldMode)
{
    if (newMode == DockConstants::FashionMode)
    {
        appBackground->setVisible(false);
    }
    else
    {
        appBackground->setVisible(true);
    }

    setActived(actived());
    resizeResources();
}

void AppItem::reanchorIcon()
{
    switch (dockCons->getDockMode()) {
    case DockConstants::FashionMode:
        m_appIcon->move((width() - m_appIcon->width()) / 2, 0);
        break;
    case DockConstants::EfficientMode:
        m_appIcon->move((width() - m_appIcon->width()) / 2, (height() - m_appIcon->height()) / 2);
        break;
    case DockConstants::ClassicMode:
        m_appIcon->move((height() - m_appIcon->height()) / 2, (height() - m_appIcon->height()) / 2);
    default:
        break;
    }
}

void AppItem::resizeBackground()
{
    appBackground->resize(width(),height());
}

void AppItem::resizeResources()
{
    if (m_appIcon != NULL)
    {
        m_appIcon->resize(dockCons->getAppIconSize(),dockCons->getAppIconSize());
        reanchorIcon();
    }

    if (appBackground != NULL)
    {
        resizeBackground();
        appBackground->move(0,0);
    }
}

void AppItem::initBackground()
{
    appBackground = new AppBackground(this);
    appBackground->move(0,0);
    connect(this, SIGNAL(widthChanged()),this, SLOT(resizeBackground()));
}

void AppItem::mousePressEvent(QMouseEvent * event)
{
    //qWarning() << "mouse press...";
    emit mousePress(event->globalX(), event->globalY());
    ////////////FOR TEST ONLY/////////////////////
    appBackground->setIsActived(!appBackground->getIsActived());
    setActived(!actived());
}

void AppItem::mouseReleaseEvent(QMouseEvent * event)
{
//    qWarning() << "mouse release...";
    emit mouseRelease(event->globalX(), event->globalY());
}

void AppItem::mouseDoubleClickEvent(QMouseEvent * event)
{
    emit mouseDoubleClick();
    ////////////FOR TEST ONLY/////////////////////
    appBackground->setIsCurrentOpened(!appBackground->getIsCurrentOpened());
}

void AppItem::mouseMoveEvent(QMouseEvent *event)
{
    //this event will only execp onec then handle by Drag
    emit dragStart();

    Qt::MouseButtons btn = event->buttons();
    if(btn == Qt::LeftButton)
    {
        QDrag* drag = new QDrag(this);
        QMimeData* data = new QMimeData();
        QImage dataImg(m_itemIconPath);
        data->setImageData(QVariant(dataImg));
        drag->setMimeData(data);

        QPixmap pixmap(m_itemIconPath);
        drag->setPixmap(pixmap);

        drag->setHotSpot(QPoint(15,15));

        drag->exec(Qt::CopyAction | Qt::MoveAction, Qt::MoveAction);
    }
}

void AppItem::enterEvent(QEvent *event)
{
    emit mouseEntered();
    appBackground->setIsHovered(true);
}

void AppItem::leaveEvent(QEvent *event)
{
    emit mouseExited();
    appBackground->setIsHovered(false);
}

void AppItem::dragEnterEvent(QDragEnterEvent *event)
{
    emit dragEntered(event);

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
    emit dragExited(event);
}

void AppItem::dropEvent(QDropEvent *event)
{
    qWarning() << "Item get drop:" << event->pos();
    emit drop(event);
}

AppItem::~AppItem()
{

}

