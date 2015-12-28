#include "dockappitem.h"

const int INVALID_MOVE_RADIUS = 10;
const QEasingCurve MOVE_ANIMATION_CURVE = QEasingCurve::OutCubic;

DockAppItem::DockAppItem(QWidget *parent) :
    DockItem(parent)
{
    m_appManager = new DBusDockedAppManager(this);
    m_dockModeData = DockModeData::instance();
    connect(m_dockModeData, &DockModeData::dockModeChanged,this, &DockAppItem::onDockModeChanged);

//    setAcceptDrops(true);
    setFixedSize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());

    initAppIcon();
    initBackground();
    initTitle();
    m_appIcon->raise();
    initClientManager();
    initPreviewContainer();
}

DockAppItemData DockAppItem::itemData() const
{
    return m_itemData;
}

QWidget *DockAppItem::getApplet()
{
    if (m_itemData.isActived && !m_itemData.xidTitleMap.isEmpty()) {
        m_previewContainer->clearUpPreview();
        //Returns a list containing all the keys in the map in ascending order.
        QList<int> xids = m_itemData.xidTitleMap.keys();
        for (int xid : xids) {
            m_previewContainer->addItem(m_itemData.xidTitleMap[xid], xid);
        }

        return m_previewContainer;
    }
    else {
        return NULL;    //use getTitle() to show title by DockItem
    }
}

QString DockAppItem::getItemId()
{
    return m_itemData.id;
}

QString DockAppItem::getTitle()
{
    return m_itemData.title;
}

void DockAppItem::setEntryProxyer(DBusDockEntry *entryProxyer)
{
    m_entryProxyer = entryProxyer;
    m_entryProxyer->setParent(this);
    connect(m_entryProxyer, &DBusDockEntry::DataChanged, this, &DockAppItem::onDbusDataChanged);

    initData();
}

void DockAppItem::mousePressEvent(QMouseEvent *event)
{
    //mouse event accept area are app-icon-area in FashionMode
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        onMousePress(event);
    else
        QFrame::mousePressEvent(event);
}

void DockAppItem::mouseReleaseEvent(QMouseEvent *event)
{
    //mouse event accept area are app-icon-area in FashionMode
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        onMouseRelease(event);
    else
        QFrame::mouseReleaseEvent(event);
}

void DockAppItem::dropEvent(QDropEvent *event)
{
    QList<QUrl> urls = event->mimeData()->urls();
    if (!urls.isEmpty()) {
        for (QUrl url : urls) {
            QString us = url.toString();
            if (!us.endsWith(".desktop")) {
                m_entryProxyer->HandleDragDrop(0, 0, us, 0);
                qDebug() << "Try to open file:" << us;
            }
        }
    }
}

void DockAppItem::enterEvent(QEvent *)
{
    //mouse event accept area are app-icon-area in FashionMode
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        onMouseEnter();
}

void DockAppItem::leaveEvent(QEvent *)
{
    //mouse event accept area are app-icon-area in FashionMode
    if (m_dockModeData->getDockMode() != Dock::FashionMode)
        onMouseLeave();
}

void DockAppItem::resizeEvent(QResizeEvent *e)
{
    Q_UNUSED(e)

    resizeResources();
}

void DockAppItem::initClientManager()
{
    m_clientManager = new DBusClientManager(this);
    connect(m_clientManager, &DBusClientManager::ActiveWindowChanged, this, &DockAppItem::setCurrentOpened);
}

void DockAppItem::initBackground()
{
    m_appBG = new DockAppBG(this);
    m_appBG->move(0, 0);
}

void DockAppItem::initPreviewContainer()
{
    m_previewContainer = new AppPreviewsContainer();
    connect(m_previewContainer,&AppPreviewsContainer::requestHide, [=]{hidePreview();});
    connect(m_previewContainer,&AppPreviewsContainer::sizeChanged, this, &DockAppItem::needPreviewUpdate);
}

void DockAppItem::initAppIcon()
{
    m_appIcon = new DockAppIcon(this);
    connect(m_appIcon, &DockAppIcon::mousePress, this, &DockAppItem::onMousePress);
    connect(m_appIcon, &DockAppIcon::mouseRelease, this, &DockAppItem::onMouseRelease);
    connect(m_appIcon, &DockAppIcon::mouseEnter, this, &DockAppItem::onMouseEnter);
    connect(m_appIcon, &DockAppIcon::mouseLeave, this, &DockAppItem::onMouseLeave);
}

void DockAppItem::initTitle()
{
    m_appTitle = new QLabel(this);
    m_appTitle->setObjectName("ClassicModeTitle");
    m_appTitle->setAlignment(Qt::AlignVCenter | Qt::AlignLeft);
}

void DockAppItem::initData()
{
    StringMap dataMap = m_entryProxyer->data();
    m_itemData.title = dataMap.value("title");
    m_itemData.iconPath = dataMap.value("icon");
    m_itemData.menuJsonString = dataMap.value("menu");
    updateXidTitleMap();
    m_itemData.isActived = dataMap.value("app-status") == "active";
    m_itemData.currentOpened = m_itemData.xidTitleMap.keys().indexOf(m_clientManager->CurrentActiveWindow().value()) != -1;
    m_itemData.id = m_entryProxyer->id();

    setActived(m_itemData.isActived);
    setCurrentOpened(m_clientManager->CurrentActiveWindow());
    updateIcon();
    updateTitle();
}

void DockAppItem::updateIcon()
{
    m_appIcon->setFixedSize(m_dockModeData->getAppIconSize(), m_dockModeData->getAppIconSize());
    m_appIcon->setIcon(m_itemData.iconPath);

    reanchorIcon();
}

void DockAppItem::updateTitle()
{
    m_itemData.title = m_entryProxyer->data().value("title");

    switch (m_dockModeData->getDockMode()) {
    case Dock::FashionMode:
    case Dock::EfficientMode:
        m_appTitle->setFixedSize(0, 0);
        m_appTitle->setVisible(false);
        return;
    case Dock::ClassicMode:
        m_appIcon->setVisible(true);
        m_appTitle->setFixedSize((m_actived ? (width() - m_appIcon->width()) : 0), m_appIcon->height());
        m_appTitle->move((m_appIcon->x() + m_appIcon->width()), m_appIcon->y());
        m_appTitle->show();
        break;
    default:
        break;
    }

    QFontMetrics fm(m_appTitle->font());
    m_appTitle->setText(fm.elidedText(m_itemData.title, Qt::ElideRight, (width() - m_appIcon->width() - 10)));

}

void DockAppItem::updateState()
{
    m_itemData.isActived = m_entryProxyer->data().value("app-status") == "active";
    setActived(m_itemData.isActived);
}

void DockAppItem::updateXidTitleMap()
{
    m_itemData.xidTitleMap.clear();
    QJsonArray array = QJsonDocument::fromJson(m_entryProxyer->data().value("app-xids").toUtf8()).array();
    for (QJsonValue v : array) {
        QJsonObject obj = v.toObject();
        m_itemData.xidTitleMap.insert(obj.value("Xid").toInt(), obj.value("Title").toString());
    }
}

void DockAppItem::updateMenuJsonString()
{
    m_itemData.menuJsonString = m_entryProxyer->data().value("menu");
}

void DockAppItem::onDbusDataChanged(const QString &, const QString &)
{
    updateTitle();
    updateState();
    updateXidTitleMap();
    updateMenuJsonString();
}

void DockAppItem::onDockModeChanged(Dock::DockMode, Dock::DockMode)
{
    setActived(actived());
    resizeResources();
}

void DockAppItem::onMousePress(QMouseEvent *event)
{
    Q_UNUSED(event)

    hidePreview(true);
}

void DockAppItem::onMouseRelease(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton) {
        m_entryProxyer->Activate(event->globalX(), event->globalY(), event->timestamp());
        m_appBG->showActivatingAnimation();
    }
    else if (event->button() == Qt::RightButton)
        showMenu();
}

void DockAppItem::onMouseEnter()
{
    if (m_dockModeData->getHideState() != Dock::HideStateShown || !hoverable())
        return;

    m_appBG->setIsHovered(true);
    showPreview();
}

void DockAppItem::onMouseLeave()
{
    m_appBG->setIsHovered(false);
    hidePreview(false);
}

void DockAppItem::resizeBackground()
{
    if (m_appBG)
        m_appBG->setFixedSize(width(), height());
}

void DockAppItem::resizeResources()
{
    if (m_appIcon != NULL)
        updateIcon();

    if (m_appBG != NULL) {
        resizeBackground();
        m_appBG->move(0, 0);
    }

    updateTitle();
}

void DockAppItem::reanchorIcon()
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
        break;
    default:
        break;
    }
}

void DockAppItem::setCurrentOpened(uint value)
{
    if (m_itemData.xidTitleMap.keys().indexOf(value) != -1) {
        m_itemData.currentOpened = true;
    }
    else {
        m_itemData.currentOpened = false;
    }

    m_appBG->setIsCurrentOpened(m_itemData.currentOpened);
}

void DockAppItem::setActived(bool value)
{
    m_actived = value;
    if (!value)
        setFixedSize(m_dockModeData->getNormalItemWidth(), m_dockModeData->getItemHeight());
    else
        setFixedSize(m_dockModeData->getActivedItemWidth(), m_dockModeData->getItemHeight());

    m_appBG->setIsActived(value);
    reanchorIcon();
}

void DockAppItem::invokeMenuItem(QString id, bool)
{
    m_entryProxyer->HandleMenuItem(id, 0);
}

QString DockAppItem::getMenuContent()
{
    return m_itemData.menuJsonString;
}

bool DockAppItem::actived() const
{
    return m_actived;
}

DockAppItem::~DockAppItem()
{

}

