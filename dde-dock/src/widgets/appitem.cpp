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
    connect(this, &AppItem::requestAnimationFinish, animation, &QPropertyAnimation::stop);
    connect(this, &AppItem::requestAnimationFinish, animation, &QPropertyAnimation::deleteLater);
}

AppItemData AppItem::itemData() const
{
    return m_itemData;
}

QWidget *AppItem::getApplet()
{
    if (!m_preview)
        initPreview();

    if (m_itemData.isActived && !m_itemData.xidTitleMap.isEmpty())
    {
        m_preview->clearUpPreview();
        //Returns a list containing all the keys in the map in ascending order.
        QList<int> xids = m_itemData.xidTitleMap.keys();
        foreach (int xid, xids) {
            m_preview->addItem(m_itemData.xidTitleMap[xid], xid);
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

    AppItem *tmpItem = qobject_cast<AppItem *>(event->source());
    if (tmpItem){    //from brother item
        event->ignore();
        emit dragEntered(event);
    }
    else if (event->mimeData()->formats().indexOf("RequestDock") != -1){    //from desktop or launcher
        QJsonObject dataObj = QJsonDocument::fromJson(event->mimeData()->data("RequestDock")).object();
        if (!dataObj.isEmpty() && !m_ddam->IsDocked(dataObj.value("appKey").toString())){
            event->ignore();
            emit dragEntered(event);
        }
    }
    else    //other files
    {
        event->setDropAction(Qt::CopyAction);
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

    m_lastPressPos = event->pos();
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
    QRect moveRect(QPoint(m_lastPressPos.x() - INVALID_MOVE_RADIUS, m_lastPressPos.y() - INVALID_MOVE_RADIUS),
                   QPoint(m_lastPressPos.x() + INVALID_MOVE_RADIUS, m_lastPressPos.y() + INVALID_MOVE_RADIUS));
    if (!moveRect.contains(event->pos())) {
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
    connect(m_preview,&AppPreviews::sizeChanged, this, &AppItem::needPreviewUpdate);
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

void AppItem::initData()
{
    StringMap dataMap = m_entryProxyer->data();
    m_itemData.title = dataMap.value("title");
    m_itemData.iconPath = dataMap.value("icon");
    m_itemData.menuJsonString = dataMap.value("menu");
    updateXidTitleMap();
    m_itemData.isActived = dataMap.value("app-status") == "active";
    m_itemData.currentOpened = m_itemData.xidTitleMap.keys().indexOf(m_clientmanager->CurrentActiveWindow().value()) != -1;
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

void AppItem::updateXidTitleMap()
{
    m_itemData.xidTitleMap.clear();
    QJsonArray nArray = QJsonDocument::fromJson(m_entryProxyer->data().value("app-xids").toUtf8()).array();
    foreach (QJsonValue value, nArray) {
        QJsonObject obj = value.toObject();
        m_itemData.xidTitleMap.insert(obj.value("Xid").toInt(), obj.value("Title").toString());
    }
}

void AppItem::updateMenuJsonString()
{
    m_itemData.menuJsonString = m_entryProxyer->data().value("menu");
}

void AppItem::onDbusDataChanged(const QString &, const QString &)
{
    updateTitle();
    updateState();
    updateXidTitleMap();
    updateMenuJsonString();
}

void AppItem::onDockModeChanged(Dock::DockMode, Dock::DockMode)
{
    setActived(actived());
    resizeResources();
}

void AppItem::onMousePress(QMouseEvent *event)
{
    //qWarning() << "mouse press...";
    emit mousePress(event);

    hidePreview(true);

    m_lastPressPos = event->pos();
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
    if (!hoverable())
        return;

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
        if (itemData().isActived)
            m_appIcon->move((height() - m_appIcon->height()) / 2, (height() - m_appIcon->height()) / 2);
        else
            m_appIcon->move((width() - m_appIcon->width()) / 2, (height() - m_appIcon->height()) / 2);
    default:
        break;
    }
}

void AppItem::setCurrentOpened(uint value)
{
    if (m_itemData.xidTitleMap.keys().indexOf(value) != -1)
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
    reanchorIcon();
}

void AppItem::invokeMenuItem(QString id, bool)
{
    m_entryProxyer->HandleMenuItem(id);
}

QString AppItem::getMenuContent()
{
    return m_itemData.menuJsonString;
}

AppItem::~AppItem()
{

}

