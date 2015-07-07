#include "appitem.h"

AppItem::AppItem(QWidget *parent) :
    AbstractDockItem(parent)
{
    setAcceptDrops(true);
    resize(dockCons->getNormalItemWidth(), dockCons->getItemHeight());
    initBackground();
    initClientManager();
    connect(dockCons, &DockModeData::dockModeChanged,this, &AppItem::slotDockModeChanged);

    initMenu();
    initPreviewAR();
}

QWidget *AppItem::getContents()
{
    AppPreviews *preview = new AppPreviews();
    QJsonArray tmpArray = QJsonDocument::fromJson(m_itemData.xidsJsonString.toUtf8()).array();
    if (m_itemData.isActived && !tmpArray.isEmpty())
    {
        foreach (QJsonValue v, tmpArray) {
            QString title = v.toObject().value("Title").toString();
            int xid = v.toObject().value("Xid").toInt();
            preview->addItem(title,xid);
        }
    }
    else
    {
        preview->setTitle(m_itemData.title);
    }
    return preview;
}

void AppItem::setEntryProxyer(DBusEntryProxyer *entryProxyer)
{
    m_entryProxyer = entryProxyer;
    m_entryProxyer->setParent(this);
    connect(m_entryProxyer, SIGNAL(DataChanged(QString,QString)),this, SLOT(dbusDataChanged(QString,QString)));

    initData();
}

void AppItem::destroyItem(const QString &id)
{

}

QString AppItem::itemId() const
{
    return m_itemData.id;
}

AppItemData AppItem::itemData() const
{
    return m_itemData;
}

void AppItem::slotDockModeChanged(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (newMode == Dock::FashionMode)
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
    case Dock::FashionMode:
        m_appIcon->move((width() - m_appIcon->width()) / 2, 0);
        break;
    case Dock::EfficientMode:
        m_appIcon->move((width() - m_appIcon->width()) / 2, (height() - m_appIcon->height()) / 2);
        break;
    case Dock::ClassicMode:
        m_appIcon->move((height() - m_appIcon->height()) / 2, (height() - m_appIcon->height()) / 2);
    default:
        break;
    }
}

void AppItem::resizeBackground()
{
    appBackground->resize(width(),height());
}

void AppItem::dbusDataChanged(const QString &key, const QString &value)
{
    updateTitle();
    updateState();
    updateXids();
    updateMenuJsonString();
}

void AppItem::setCurrentOpened(uint value)
{
    if (m_itemData.xidsJsonString.indexOf(QString::number(value)) != -1)
    {
        m_itemData.currentOpened = true;
        appBackground->setIsCurrentOpened(true);
    }
    else
    {
        m_itemData.currentOpened = false;
        appBackground->setIsCurrentOpened(false);
    }
}

void AppItem::menuItemInvoked(QString id, bool)
{
    m_entryProxyer->HandleMenuItem(id);
    m_menuManager->UnregisterMenu(m_menuInterfacePath);
}

void AppItem::resizeResources()
{
    if (m_appIcon != NULL)
    {
        m_appIcon->resize(dockCons->getAppIconSize(),dockCons->getAppIconSize());
        updateIcon();
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

void AppItem::initClientManager()
{
    m_clientmanager = new DBusClientManager(this);
    connect(m_clientmanager, SIGNAL(ActiveWindowChanged(uint)),this, SLOT(setCurrentOpened(uint)));
}

void AppItem::setActived(bool value)
{
    m_isActived = value;
    if (!value)
        resize(dockCons->getNormalItemWidth(), dockCons->getItemHeight());
    else
        resize(dockCons->getActivedItemWidth(), dockCons->getItemHeight());

    appBackground->setIsActived(value);
}

void AppItem::initData()
{
    StringMap dataMap = m_entryProxyer->data();
    m_itemData.title = dataMap.value("title");
    m_itemData.iconPath = dataMap.value("icon");
    m_itemData.menuJsonString = dataMap.value("menu");
    m_itemData.xidsJsonString = dataMap.value("app-xids");
    m_itemData.isActived = dataMap.value("app-status") == "active";
    m_itemData.currentOpened = m_itemData.xidsJsonString.indexOf(QString::number(m_clientmanager->CurrentActiveWindow())) != -1;
    m_itemData.id = m_entryProxyer->id();

    setActived(m_itemData.isActived);
    setCurrentOpened(m_clientmanager->CurrentActiveWindow());
    updateIcon();
}

void AppItem::updateIcon()
{
    if (m_appIcon == NULL)
    {
        m_appIcon = new AppIcon(this);
    }
    m_appIcon->resize(height(), height());
    m_appIcon->setIcon(m_itemData.iconPath);

    reanchorIcon();
}

void AppItem::updateTitle()
{
    m_itemData.title = m_entryProxyer->data().value("title");
    //TODO,update view label
}

void AppItem::updateState()
{
    m_itemData.isActived = m_entryProxyer->data().value("app-status") == "active";
    setActived(m_itemData.isActived);
    appBackground->setIsActived(m_itemData.isActived);
}

void AppItem::updateXids()
{
    m_itemData.xidsJsonString = m_entryProxyer->data().value("app-xids");
}

void AppItem::updateMenuJsonString()
{
    m_itemData.menuJsonString = m_entryProxyer->data().value("menu");
}

void AppItem::initMenu()
{
    m_menuManager = new DBusMenuManager(this);
}

void AppItem::initPreviewAR()
{
    m_previewAR = new ArrowRectangle();
    m_previewAR->setHeight(130);
    m_previewAR->setWidth(200);
}

void AppItem::showMenu()
{
    if (m_menuManager->isValid()){
        QDBusPendingReply<QDBusObjectPath> pr = m_menuManager->RegisterMenu();
        if (pr.count() == 1){
            QDBusObjectPath op = pr.argumentAt(0).value<QDBusObjectPath>();
            m_menuInterfacePath = op.path();
            DBusMenu *m_menu = new DBusMenu(m_menuInterfacePath,this);
            connect(m_menu,SIGNAL(MenuUnregistered()),m_menu,SLOT(deleteLater()));
            connect(m_menu,SIGNAL(ItemInvoked(QString,bool)),this,SLOT(menuItemInvoked(QString,bool)));

            QJsonObject targetObj;
            targetObj.insert("x",QJsonValue(globalX() + width() / 2));
            targetObj.insert("y",QJsonValue(globalY() - 5));
            targetObj.insert("isDockMenu",QJsonValue(true));
            targetObj.insert("menuJsonContent",QJsonValue(m_itemData.menuJsonString));

            m_menu->ShowMenu(QString(QJsonDocument(targetObj).toJson()));
        }
    }
}

void AppItem::showPreview()
{
    QWidget *tmpContent = getContents();
    m_previewAR->setMinimumSize(tmpContent->width() + Dock::APP_PREVIEW_MARGIN * 2,tmpContent->height() + Dock::APP_PREVIEW_MARGIN * 2);
    m_previewAR->resize(tmpContent->width() + Dock::APP_PREVIEW_MARGIN * 2,tmpContent->height() + Dock::APP_PREVIEW_MARGIN * 2);
    m_previewAR->setContent(getContents());
    m_previewAR->showAtBottom(globalX() + width() / 2,globalY() - 5);
}

void AppItem::hidePreview()
{
    m_previewAR->delayHide();
}

void AppItem::mousePressEvent(QMouseEvent * event)
{
    //qWarning() << "mouse press...";
    emit mousePress(event->globalX(), event->globalY());
    hidePreview();
}

void AppItem::mouseReleaseEvent(QMouseEvent * event)
{
//    qWarning() << "mouse release...";
    emit mouseRelease(event->globalX(), event->globalY());

    if (event->button() == Qt::LeftButton)
        m_entryProxyer->Activate(event->globalX(),event->globalY());
    else if (event->button() == Qt::RightButton)
        showMenu();
}

void AppItem::mouseDoubleClickEvent(QMouseEvent * event)
{
    emit mouseDoubleClick();
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
        QImage dataImg(m_appIcon->getIconPath());
        data->setImageData(QVariant(dataImg));
        drag->setMimeData(data);

        QPixmap pixmap(m_appIcon->getIconPath());
        drag->setPixmap(pixmap);

        drag->setHotSpot(QPoint(15,15));

        drag->exec(Qt::CopyAction | Qt::MoveAction, Qt::MoveAction);
    }
}

void AppItem::enterEvent(QEvent *event)
{
    emit mouseEntered();
    appBackground->setIsHovered(true);
    showPreview();
}

void AppItem::leaveEvent(QEvent *event)
{
    emit mouseExited();
    appBackground->setIsHovered(false);
    hidePreview();
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

