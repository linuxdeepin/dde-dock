#include "appitem.h"

AppItem::AppItem(QWidget *parent) :
    AbstractDockItem(parent)
{
    setAcceptDrops(true);

    resize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());

    initAppIcon();
    initBackground();
    initHighlight();
    initTitle();
    m_appIcon->raise();
    initClientManager();
    connect(m_dockModeData, &DockModeData::dockModeChanged,this, &AppItem::onDockModeChanged);

    initMenu();
    initPreview();
}

void AppItem::moveWithAnimation(QPoint targetPos, int duration)
{
    setNextPos(targetPos);

    QPropertyAnimation *animation = new QPropertyAnimation(this, "pos");
    animation->setStartValue(pos());
    animation->setEndValue(getNextPos());
    animation->setDuration(duration);
    animation->setEasingCurve(MOVE_ANIMATION_CURVE);
    animation->start();
    connect(animation, &QPropertyAnimation::finished, this, &AppItem::moveAnimationFinished);
    connect(animation, &QPropertyAnimation::finished, animation, &QPropertyAnimation::deleteLater);
}

AppItemData AppItem::itemData() const
{
    return m_itemData;
}

QWidget *AppItem::getApplet()
{
    if (!m_preview)
        initPreview();

    QJsonArray tmpArray = QJsonDocument::fromJson(m_itemData.xidsJsonString.toUtf8()).array();
    if (m_itemData.isActived && !tmpArray.isEmpty())
    {
        foreach (QJsonValue v, tmpArray) {
            QString title = v.toObject().value("Title").toString();
            int xid = v.toObject().value("Xid").toInt();
            m_preview->addItem(title,xid);
        }
    } else {
        return NULL;    //use getTitle() to show title by abstractdockitem
    }

    return m_preview;
}

QString AppItem::getItemId()
{
    return m_itemData.id;
}

QString AppItem::getTitle()
{
    return m_itemData.title;
}

void AppItem::setEntryProxyer(DBusEntryProxyer *entryProxyer)
{
    m_entryProxyer = entryProxyer;
    m_entryProxyer->setParent(this);
    connect(m_entryProxyer, &DBusEntryProxyer::DataChanged, this, &AppItem::onDbusDataChanged);

    initData();
}

void AppItem::dragEnterEvent(QDragEnterEvent *event)
{
    onMouseLeave(); //enterEvent may be active along with the dragEnterEvent, so hide preview here

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

void AppItem::mousePressEvent(QMouseEvent *event)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        onMousePress(event);
    else
        QFrame::mousePressEvent(event);
}

void AppItem::mouseReleaseEvent(QMouseEvent *event)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        onMouseRelease(event);
    else
        QFrame::mouseReleaseEvent(event);
}

void AppItem::mouseMoveEvent(QMouseEvent *event)
{
    //this event will only execp onec then handle by Drag
    emit dragStart();

    Qt::MouseButtons btn = event->buttons();
    if(btn == Qt::LeftButton)
    {
        //drag and mimeData object will delete automatically
        QDrag* drag = new QDrag(this);
        QMimeData* mimeData = new QMimeData();
        QImage dataImg = m_appIcon->grab().toImage();
        mimeData->setImageData(QVariant(dataImg));
        drag->setMimeData(mimeData);
        drag->setHotSpot(QPoint(15,15));

        if (m_dockModeData->getDockMode() == Dock::FashionMode){
            QPixmap pixmap = m_appIcon->grab();
            drag->setPixmap(pixmap.scaled(m_dockModeData->getAppIconSize(), m_dockModeData->getAppIconSize()));
        }
        else{
            QPixmap pixmap = this->grab();
            drag->setPixmap(pixmap.scaled(this->size()));
        }

        drag->exec(Qt::CopyAction | Qt::MoveAction, Qt::MoveAction);
    }
}

void AppItem::dropEvent(QDropEvent *event)
{
    qWarning() << "Item get drop:" << event->pos();
}

void AppItem::enterEvent(QEvent *)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        onMouseEnter();
}

void AppItem::leaveEvent(QEvent *)
{
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        onMouseLeave();
}

void AppItem::initClientManager()
{
    m_clientmanager = new DBusClientManager(this);
    connect(m_clientmanager, &DBusClientManager::ActiveWindowChanged, this, &AppItem::setCurrentOpened);
}

void AppItem::initBackground()
{
    m_appBackground = new AppBackground(this);
    m_appBackground->move(0,0);
    connect(this, &AppItem::mouseRelease, m_appBackground, &AppBackground::slotMouseRelease);
    connect(this, &AppItem::widthChanged, this, &AppItem::resizeBackground);
}

void AppItem::initPreview()
{
    m_preview = new AppPreviews();
    connect(m_preview,&AppPreviews::requestHide, [=]{hidePreview();});
    connect(m_preview,&AppPreviews::sizeChanged, this, &AppItem::resizePreview);
    connect(this, &AppItem::previewHidden, m_preview, &AppPreviews::clearUpPreview);
}

void AppItem::initAppIcon()
{
    m_appIcon = new AppIcon(this);
    connect(m_appIcon, &AppIcon::mousePress, this, &AppItem::onMousePress);
    connect(m_appIcon, &AppIcon::mouseRelease, this, &AppItem::onMouseRelease);
    connect(m_appIcon, &AppIcon::mouseEnter, this, &AppItem::onMouseEnter);
    connect(m_appIcon, &AppIcon::mouseLeave, this, &AppItem::onMouseLeave);
}

void AppItem::initTitle()
{
    m_appTitle = new QLabel(this);
    m_appTitle->setObjectName("ClassicModeTitle");
    m_appTitle->setAlignment(Qt::AlignVCenter | Qt::AlignLeft);
}

void AppItem::initMenu()
{
    m_menuManager = new DBusMenuManager(this);
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
    updateTitle();
}

void AppItem::updateIcon()
{
    m_appIcon->resize(m_dockModeData->getAppIconSize(),m_dockModeData->getAppIconSize());
    m_appIcon->setIcon(m_itemData.iconPath);

    reanchorIcon();
}

void AppItem::updateTitle()
{
    m_itemData.title = m_entryProxyer->data().value("title");

    switch (m_dockModeData->getDockMode()) {
    case Dock::FashionMode:
    case Dock::EfficientMode:
        m_appTitle->resize(0,0);
        m_appTitle->setVisible(false);
        return;
    case Dock::ClassicMode:
        m_appIcon->setVisible(true);
        m_appTitle->resize(m_isActived ? (width() - m_appIcon->width()) : 0,m_appIcon->height());
        m_appTitle->move(m_appIcon->x() + m_appIcon->width(), m_appIcon->y());
        m_appTitle->show();
        break;
    default:
        break;
    }

    QFontMetrics fm(m_appTitle->font());
    m_appTitle->setText(fm.elidedText(m_itemData.title,Qt::ElideRight,width() - m_appIcon->width() - 10));

}

void AppItem::updateState()
{
    m_itemData.isActived = m_entryProxyer->data().value("app-status") == "active";
    setActived(m_itemData.isActived);
    m_appBackground->setIsActived(m_itemData.isActived);
}

void AppItem::updateXids()
{
    m_itemData.xidsJsonString = m_entryProxyer->data().value("app-xids");
}

void AppItem::updateMenuJsonString()
{
    m_itemData.menuJsonString = m_entryProxyer->data().value("menu");
}

void AppItem::onDbusDataChanged(const QString &, const QString &)
{
    updateTitle();
    updateState();
    updateXids();
    updateMenuJsonString();
}

void AppItem::onDockModeChanged(Dock::DockMode, Dock::DockMode)
{
    setActived(actived());
    resizeResources();
}

void AppItem::onMenuItemInvoked(QString id, bool)
{
    m_entryProxyer->HandleMenuItem(id);
    m_menuManager->UnregisterMenu(m_menuInterfacePath);
}

void AppItem::onMousePress(QMouseEvent *event)
{
    //qWarning() << "mouse press...";
    emit mousePress(event);
    hidePreview(0);
}

void AppItem::onMouseRelease(QMouseEvent *event)
{
    //qWarning() << "mouse release...";
    emit mouseRelease(event);

    if (event->button() == Qt::LeftButton)
        m_entryProxyer->Activate(event->globalX(),event->globalY());
    else if (event->button() == Qt::RightButton)
        showMenu();
}

void AppItem::onMouseEnter()
{
    emit mouseEntered();
    m_appBackground->setIsHovered(true);
    showPreview();
}

void AppItem::onMouseLeave()
{
    emit mouseExited();
    m_appBackground->setIsHovered(false);
    hidePreview();
}

void AppItem::resizeBackground()
{
    m_appBackground->resize(width(),height());
}

void AppItem::resizeResources()
{
    if (m_appIcon != NULL)
        updateIcon();

    if (m_appBackground != NULL)
    {
        resizeBackground();
        m_appBackground->move(0,0);
    }

    updateTitle();
}

void AppItem::reanchorIcon()
{
    switch (m_dockModeData->getDockMode()) {
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

void AppItem::setCurrentOpened(uint value)
{
    if (m_itemData.xidsJsonString.indexOf(QString::number(value)) != -1)
    {
        m_itemData.currentOpened = true;
        m_appBackground->setIsCurrentOpened(true);
    }
    else
    {
        m_itemData.currentOpened = false;
        m_appBackground->setIsCurrentOpened(false);
    }
}

void AppItem::setActived(bool value)
{
    m_isActived = value;
    if (!value)
        resize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());
    else
        resize(m_dockModeData->getActivedItemWidth(), m_dockModeData->getItemHeight());

    m_appBackground->setIsActived(value);
}

void AppItem::showMenu()
{
    if (m_menuManager->isValid()){
        QDBusPendingReply<QDBusObjectPath> pr = m_menuManager->RegisterMenu();
        if (pr.count() == 1){
            QDBusObjectPath op = pr.argumentAt(0).value<QDBusObjectPath>();
            m_menuInterfacePath = op.path();
            DBusMenu *m_menu = new DBusMenu(m_menuInterfacePath,this);
            connect(m_menu, &DBusMenu::MenuUnregistered, m_menu, &DBusMenu::deleteLater);
            connect(m_menu, &DBusMenu::ItemInvoked, this, &AppItem::onMenuItemInvoked);

            QJsonObject targetObj;
            targetObj.insert("x",QJsonValue(globalX() + width() / 2));
            targetObj.insert("y",QJsonValue(globalY() - 5));
            targetObj.insert("isDockMenu",QJsonValue(true));
            targetObj.insert("menuJsonContent",QJsonValue(m_itemData.menuJsonString));

            m_menu->ShowMenu(QString(QJsonDocument(targetObj).toJson()));
        }
    }
}

AppItem::~AppItem()
{
    if (m_preview)
        m_preview->deleteLater();
}

