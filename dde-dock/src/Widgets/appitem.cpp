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
    connect(m_dockModeData, &DockModeData::dockModeChanged,this, &AppItem::slotDockModeChanged);

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
}

QWidget *AppItem::getApplet()
{
    if (!m_preview)
        return NULL;

    QJsonArray tmpArray = QJsonDocument::fromJson(m_itemData.xidsJsonString.toUtf8()).array();
    if (m_itemData.isActived && !tmpArray.isEmpty())
    {
        foreach (QJsonValue v, tmpArray) {
            QString title = v.toObject().value("Title").toString();
            int xid = v.toObject().value("Xid").toInt();
            m_preview->addItem(title,xid);
        }
    } else {
        m_titleLabel = new ItemTitleLabel;

        m_titleLabel->setTitle(m_itemData.title);
        m_preview->setTitleLabel(m_titleLabel);
    }

    return m_preview;
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

QString AppItem::getItemId()
{
    return m_itemData.id;
}

AppItemData AppItem::itemData() const
{
    return m_itemData;
}

void AppItem::slotMousePress(QMouseEvent *event)
{
    //qWarning() << "mouse press...";
    emit mousePress(event->globalX(), event->globalY());
    hidePreview();
}

void AppItem::slotMouseRelease(QMouseEvent *event)
{
    //qWarning() << "mouse release...";
    emit mouseRelease(event->globalX(), event->globalY());

    if (event->button() == Qt::LeftButton)
        m_entryProxyer->Activate(event->globalX(),event->globalY());
    else if (event->button() == Qt::RightButton)
        showMenu();
}

void AppItem::slotMouseEnter()
{
    emit mouseEntered();
    m_appBackground->setIsHovered(true);
    showPreview();
}

void AppItem::slotMouseLeave()
{
    emit mouseExited();
    m_appBackground->setIsHovered(false);
    hidePreview();
}

void AppItem::slotDockModeChanged(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    if (newMode == Dock::FashionMode)
    {
//        m_appBackground->setVisible(false);
    }
    else
    {
        m_appBackground->setVisible(true);
    }

    setActived(actived());
    resizeResources();
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

void AppItem::resizeBackground()
{
    m_appBackground->resize(width(),height());
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
        m_appBackground->setIsCurrentOpened(true);
    }
    else
    {
        m_itemData.currentOpened = false;
        m_appBackground->setIsCurrentOpened(false);
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
        updateIcon();
    }

    if (m_appBackground != NULL)
    {
        resizeBackground();
        m_appBackground->move(0,0);
    }

    updateTitle();
}

void AppItem::initBackground()
{
    m_appBackground = new AppBackground(this);
    m_appBackground->move(0,0);
    connect(this, &AppItem::mousePress, m_appBackground, &AppBackground::slotMousePress);
    connect(this, SIGNAL(widthChanged()),this, SLOT(resizeBackground()));

    if (m_dockModeData->getDockMode() == Dock::FashionMode)
    {
//        m_appBackground->setVisible(false);
    }
    else
    {
        m_appBackground->setVisible(true);
    }
}

void AppItem::initTitle()
{
    m_appTitle = new QLabel(this);
    m_appTitle->setObjectName("ClassicModeTitle");
    m_appTitle->setAlignment(Qt::AlignVCenter | Qt::AlignLeft);
}

void AppItem::initAppIcon()
{
    m_appIcon = new AppIcon(this);
    connect(m_appIcon, &AppIcon::mousePress, this, &AppItem::slotMousePress);
    connect(m_appIcon, &AppIcon::mouseRelease, this, &AppItem::slotMouseRelease);
    connect(m_appIcon, &AppIcon::mouseEnter, this, &AppItem::slotMouseEnter);
    connect(m_appIcon, &AppIcon::mouseLeave, this, &AppItem::slotMouseLeave);
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
        resize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());
    else
        resize(m_dockModeData->getActivedItemWidth(), m_dockModeData->getItemHeight());

    m_appBackground->setIsActived(value);
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

void AppItem::initMenu()
{
    m_menuManager = new DBusMenuManager(this);
}

void AppItem::initPreview()
{
    m_preview = new AppPreviews();
    connect(m_preview,&AppPreviews::sizeChanged,this,&AppItem::resizePreview);
    connect(this, &AppItem::previewHidden, m_preview, &AppPreviews::clearUpPreview);
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

void AppItem::mouseMoveEvent(QMouseEvent *event)
{
    //this event will only execp onec then handle by Drag
    emit dragStart();

    Qt::MouseButtons btn = event->buttons();
    if(btn == Qt::LeftButton)
    {
        QDrag* drag = new QDrag(this);
        QMimeData* data = new QMimeData();
        QImage dataImg = m_appIcon->grab().toImage();
        data->setImageData(QVariant(dataImg));
        drag->setMimeData(data);

        QPixmap pixmap = m_appIcon->grab();
        drag->setPixmap(pixmap.scaled(m_dockModeData->getAppIconSize(), m_dockModeData->getAppIconSize()));

        drag->setHotSpot(QPoint(15,15));

        drag->exec(Qt::CopyAction | Qt::MoveAction, Qt::MoveAction);
    }
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
}

AppItem::~AppItem()
{
    if (m_preview)
        m_preview->deleteLater();
}

