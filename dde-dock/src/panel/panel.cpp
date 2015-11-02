#include <QHBoxLayout>
#include <QTimer>

#include "panel.h"
#include "controller/dockmodedata.h"
#include "controller/plugins/dockpluginproxy.h"
#include "controller/plugins/dockpluginmanager.h"

Panel::Panel(QWidget *parent)
    : QLabel(parent),m_parentWidget(parent)
{
    setObjectName("Panel");

    initGlobalPreview();
    initShowHideAnimation();
    initHideStateManager();
    initWidthAnimation();
    initPluginLayout();
    initAppLayout();
    initPluginManager();
    initAppManager();
    initReflection();
    initScreenMask();

    setMinimumHeight(m_dockModeData->getDockHeight());  //set height for border-image calculate
    reloadStyleSheet();

    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &Panel::onDockModeChanged);
}

bool Panel::isFashionMode()
{
    return m_isFashionMode;
}

void Panel::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::RightButton)
        showPanelMenu();
}

void Panel::mouseReleaseEvent(QMouseEvent *)
{

}

void Panel::initShowHideAnimation()
{
    QStateMachine * machine = new QStateMachine(this);
    QState * showState = new QState(machine);
    showState->assignProperty(this,"y", 0);
    QState * hideState = new QState(machine);
    hideState->assignProperty(this,"y", height());
    machine->setInitialState(showState);

    QPropertyAnimation *showAnimation = new QPropertyAnimation(this, "y");
    showAnimation->setDuration(SHOW_HIDE_ANIMATION_DURATION);
    showAnimation->setEasingCurve(SHOW_HIDE_EASINGCURVE);
    connect(showAnimation,&QPropertyAnimation::finished,this,&Panel::onShowPanelFinished);
    QPropertyAnimation *hideAnimation = new QPropertyAnimation(this, "y");
    hideAnimation->setDuration(SHOW_HIDE_ANIMATION_DURATION);
    hideAnimation->setEasingCurve(SHOW_HIDE_EASINGCURVE);
    connect(hideAnimation,&QPropertyAnimation::finished,this,&Panel::onHidePanelFinished);

    QSignalTransition *ts1 = showState->addTransition(this,SIGNAL(startHide()), hideState);
    ts1->addAnimation(hideAnimation);
    connect(ts1,&QSignalTransition::triggered,[=]{m_HSManager->SetState(2);});
    QSignalTransition *ts2 = hideState->addTransition(this,SIGNAL(startShow()),showState);
    ts2->addAnimation(showAnimation);
    connect(ts2,&QSignalTransition::triggered,[=]{m_HSManager->SetState(0);});

    machine->start();
}

void Panel::initHideStateManager()
{
    m_HSManager = new DBusHideStateManager(this);
    connect(m_HSManager,&DBusHideStateManager::ChangeState,this,&Panel::onHideStateChanged);

    //for initialization
    m_HSManager->UpdateState();
}

void Panel::initWidthAnimation()
{
    m_widthAnimation = new QPropertyAnimation(this, "width", this);
    m_widthAnimation->setDuration(WIDTH_ANIMATION_DURATION);
    connect(m_widthAnimation, &QPropertyAnimation::valueChanged, [=]{
        m_appLayout->move(FASHION_PANEL_LPADDING, 1);
        m_pluginLayout->move(width() - m_pluginLayout->width() - FASHION_PANEL_RPADDING, 1);
        updateRightReflection();

        emit sizeChanged();
    });
}

void Panel::initPluginManager()
{
    m_pluginManager = new DockPluginManager(this);

    connect(m_dockModeData, &DockModeData::dockModeChanged, m_pluginManager, &DockPluginManager::onDockModeChanged);
    connect(m_pluginManager, &DockPluginManager::itemMove, [=](AbstractDockItem *baseItem, AbstractDockItem *targetItem){
        m_pluginLayout->moveItem(m_pluginLayout->indexOf(targetItem), m_pluginLayout->indexOf(baseItem));
    });
    connect(m_pluginManager, &DockPluginManager::itemAppend, [=](AbstractDockItem *targetItem){
        m_pluginLayout->addItem(targetItem);
        connect(targetItem, &AbstractDockItem::needPreviewShow, this, &Panel::onNeedPreviewShow);
        connect(targetItem, &AbstractDockItem::needPreviewHide, this, &Panel::onNeedPreviewHide);
        connect(targetItem, &AbstractDockItem::needPreviewImmediatelyHide, this, &Panel::onNeedPreviewImmediatelyHide);
        connect(targetItem, &AbstractDockItem::needPreviewUpdate, m_globalPreview, &PreviewWindow::resizeWithContent);
    });
    connect(m_pluginManager, &DockPluginManager::itemInsert, [=](AbstractDockItem *baseItem, AbstractDockItem *targetItem){
        m_pluginLayout->insertItem(targetItem, m_pluginLayout->indexOf(baseItem));
        connect(targetItem, &AbstractDockItem::needPreviewShow, this, &Panel::onNeedPreviewShow);
        connect(targetItem, &AbstractDockItem::needPreviewHide, this, &Panel::onNeedPreviewHide);
        connect(targetItem, &AbstractDockItem::needPreviewImmediatelyHide, this, &Panel::onNeedPreviewImmediatelyHide);
        connect(targetItem, &AbstractDockItem::needPreviewUpdate, m_globalPreview, &PreviewWindow::resizeWithContent);
    });
    connect(m_pluginManager, &DockPluginManager::itemRemoved, [=](AbstractDockItem* item) {
        item->setVisible(false);
        item->deleteLater();
        m_pluginLayout->removeItem(item);
    });
    connect(PanelMenu::instance(), &PanelMenu::settingPlugin, [=]{
        QRect rec = QApplication::desktop()->screenGeometry();
        m_pluginManager->onPluginsSetting(rec.height() - height());
    });
}

void Panel::initPluginLayout()
{
    m_pluginLayout = new DockLayout(this);
    m_pluginLayout->setSpacing(m_dockModeData->getAppletsItemSpacing());
    m_pluginLayout->resize(0, m_dockModeData->getItemHeight());
    connect(m_pluginLayout, &DockLayout::contentsWidthChange, this, &Panel::onLayoutContentsWidthChanged);
}

void Panel::initAppLayout()
{
    m_appLayout = new DockLayout(this);
    m_appLayout->setaddItemDelayInterval(0);
    m_appLayout->setAcceptDrops(true);
    m_appLayout->setSpacing(m_dockModeData->getAppItemSpacing());
    m_appLayout->move(0, 1);

    connect(m_appLayout, &DockLayout::startDrag, this, &Panel::onItemDragStarted);
    connect(m_appLayout, &DockLayout::itemDropped, this, &Panel::onItemDropped);
    connect(m_appLayout, &DockLayout::contentsWidthChange, this, &Panel::onLayoutContentsWidthChanged);

    //for plugin layout mask
    connect(m_appLayout, &DockLayout::startDrag, this, &Panel::showPluginLayoutMask);
    connect(m_appLayout, &DockLayout::itemDropped, this, &Panel::hidePluginLayoutMask);
}

void Panel::initAppManager()
{
    m_appManager = new AppManager(this);
    connect(m_appManager, &AppManager::entryAdded, this, &Panel::onAppItemAdd);
    connect(m_appManager, &AppManager::entryRemoved, this, &Panel::onAppItemRemove);
}

void Panel::initReflection()
{
    if (m_appLayout)
    {
        m_appReflection = new ReflectionEffect(m_appLayout, this);
        connect(m_appLayout, &DockLayout::contentsWidthChange, this, &Panel::updateLeftReflection);
        connect(m_dockModeData, &DockModeData::dockModeChanged, this, &Panel::updateLeftReflection);
        updateLeftReflection();
    }

    if (m_pluginLayout)
    {
        m_pluginReflection = new ReflectionEffect(m_pluginLayout, this);
        connect(m_appLayout, &DockLayout::contentsWidthChange, this, &Panel::updateRightReflection);
        connect(m_pluginLayout, &DockLayout::contentsWidthChange, this, &Panel::updateRightReflection);
        connect(m_dockModeData, &DockModeData::dockModeChanged, this, &Panel::updateRightReflection);
        updateRightReflection();
    }
}

void Panel::initScreenMask()
{
    m_maskWidget = new ScreenMask();
    m_maskWidget->hide();
    connect(m_maskWidget, &ScreenMask::itemDropped, this, &Panel::onItemDropped);
    connect(m_maskWidget, &ScreenMask::itemEntered, m_appLayout, &DockLayout::removeSpacingItem);
    connect(m_maskWidget, &ScreenMask::itemMissing, m_appLayout, &DockLayout::restoreTmpItem);

    //for plugin layout mask
    connect(m_maskWidget, &ScreenMask::itemDropped, this, &Panel::hidePluginLayoutMask);
    connect(m_maskWidget, &ScreenMask::itemMissing, this, &Panel::hidePluginLayoutMask);
}

void Panel::initGlobalPreview()
{
    m_globalPreview = new PreviewWindow(DArrowRectangle::ArrowBottom);

    //make sure all app-preview will be destroy to save resources
    connect(m_globalPreview, &PreviewWindow::showFinish, [=] (QWidget *lastContent) {
        if (lastContent) {
            AppPreviewsContainer *tmpFrame = qobject_cast<AppPreviewsContainer *>(lastContent);
            if (tmpFrame)
                tmpFrame->clearUpPreview();
        }
    });
    connect(m_globalPreview, &PreviewWindow::hideFinish, [=] (QWidget *lastContent) {
        if (lastContent) {
            AppPreviewsContainer *tmpFrame = qobject_cast<AppPreviewsContainer *>(lastContent);
            if (tmpFrame)
                tmpFrame->clearUpPreview();
        }
    });
}

void Panel::onItemDropped()
{
    m_maskWidget->hide();
    AppItem *item = qobject_cast<AppItem *>(m_appLayout->getDraggingItem());
    if (item)
        m_ddam->RequestUndock(item->itemData().id);
    m_appLayout->clearTmpItem();
    m_appLayout->relayout();
}

void Panel::onItemDragStarted()
{
    m_maskWidget->show();
}

void Panel::onLayoutContentsWidthChanged()
{
    if (m_dockModeData->getDockMode() == Dock::FashionMode)
    {
        m_appLayout->resize(m_appLayout->getContentsWidth() + m_dockModeData->getAppItemSpacing(),m_dockModeData->getItemHeight());
        m_pluginLayout->resize(m_pluginLayout->getContentsWidth(),m_dockModeData->getAppletsItemHeight());

        int targetWidth = FASHION_PANEL_LPADDING
                + FASHION_PANEL_RPADDING
                + m_appLayout->getContentsWidth()
                + m_pluginLayout->getContentsWidth();

        m_widthAnimation->setStartValue(width());
        m_widthAnimation->setEndValue(targetWidth);
        m_widthAnimation->start();

    }
    else
    {
        QRect rec = QApplication::desktop()->screenGeometry();
        m_pluginLayout->resize(m_pluginLayout->getContentsWidth(),m_dockModeData->getItemHeight());
        m_pluginLayout->move(rec.width() - m_pluginLayout->width(),1);

        m_appLayout->move(0,1);
        m_appLayout->resize(rec.width() - m_pluginLayout->width() ,m_dockModeData->getItemHeight());

        this->setFixedSize(m_appLayout->width() + m_pluginLayout->width(),m_dockModeData->getDockHeight());

        emit sizeChanged();
    }
}

void Panel::onAppItemAdd(AbstractDockItem *item)
{
    m_appLayout->addItem(item);
    connect(item, &AbstractDockItem::needPreviewShow, this, &Panel::onNeedPreviewShow);
    connect(item, &AbstractDockItem::needPreviewHide, this, &Panel::onNeedPreviewHide);
    connect(item, &AbstractDockItem::needPreviewImmediatelyHide, this, &Panel::onNeedPreviewImmediatelyHide);
    connect(item, &AbstractDockItem::needPreviewUpdate, m_globalPreview, &PreviewWindow::resizeWithContent);
}

void Panel::onAppItemRemove(const QString &id)
{
    QList<AbstractDockItem *> tmpList = m_appLayout->getItemList();
    for (int i = 0; i < tmpList.count(); i ++)
    {
        AppItem *tmpItem = qobject_cast<AppItem *>(tmpList.at(i));
        if (tmpItem && tmpItem->getItemId() == id)
        {
            tmpItem->setVisible(false);
            tmpItem->deleteLater();
            m_appLayout->removeItem(i);
            return;
        }
    }
}

void Panel::onDockModeChanged(Dock::DockMode newMode, Dock::DockMode)
{
    m_appLayout->setSpacing(m_dockModeData->getAppItemSpacing());
    m_appLayout->relayout();
    m_pluginLayout->setSpacing(m_dockModeData->getAppletsItemSpacing());
    m_pluginLayout->relayout();

    reanchorsLayout(newMode);

    reloadStyleSheet();
}

void Panel::onHideStateChanged(int dockState)
{
    bool containsMouse = m_parentWidget->geometry().contains(QCursor::pos());
    if (dockState == 0)
        emit startShow();
    else if (dockState == 1 && !containsMouse)
        emit startHide();
}

void Panel::onShowPanelFinished()
{
    m_HSManager->SetState(1);
    emit panelHasShown();
}

void Panel::onHidePanelFinished()
{
    m_HSManager->SetState(3);
    emit panelHasHidden();
}

void Panel::onNeedPreviewHide()
{
    m_globalPreview->hidePreview(DELAY_HIDE_PREVIEW_INTERVAL);
}

void Panel::onNeedPreviewShow(QPoint pos)
{
    AbstractDockItem *item = qobject_cast<AbstractDockItem *>(sender());
    if (item && item->getApplet()) {
        m_globalPreview->setArrowX(-1);//reset x to move arrow to horizontal-center
        m_globalPreview->setContent(item->getApplet());
        m_globalPreview->showPreview(pos.x(),
                                     pos.y() + m_globalPreview->shadowBlurRadius() + m_globalPreview->shadowDistance(),
                                     DELAY_SHOW_PREVIEW_INTERVAL);
    }
}

void Panel::onNeedPreviewImmediatelyHide()
{
    m_globalPreview->hidePreview(0);
}

void Panel::reanchorsLayout(Dock::DockMode mode)
{
    if (mode == Dock::FashionMode)
    {
        m_appLayout->resize(m_appLayout->getContentsWidth() + m_dockModeData->getAppItemSpacing(),m_dockModeData->getItemHeight());
        m_pluginLayout->resize(m_pluginLayout->getContentsWidth(),m_dockModeData->getAppletsItemHeight());
        this->setFixedSize(FASHION_PANEL_LPADDING
                           + FASHION_PANEL_RPADDING
                           + m_appLayout->getContentsWidth()
                           + m_pluginLayout->getContentsWidth()
                           ,m_dockModeData->getDockHeight());
        m_appLayout->move(FASHION_PANEL_LPADDING,1);

        m_pluginLayout->move(m_appLayout->x() + m_appLayout->width() - m_dockModeData->getAppItemSpacing(),1);

        emit sizeChanged();
    }
    else
    {
        QRect rec = QApplication::desktop()->screenGeometry();
        m_pluginLayout->resize(m_pluginLayout->getContentsWidth(), m_dockModeData->getItemHeight());
        m_pluginLayout->move(rec.width() - m_pluginLayout->width(),1);

        m_appLayout->move(0,1);
        m_appLayout->resize(rec.width() - m_pluginLayout->width() ,m_dockModeData->getItemHeight());

        this->setFixedSize(m_appLayout->width() + m_pluginLayout->width(), m_dockModeData->getDockHeight());

        emit sizeChanged();
    }
}

void Panel::updateRightReflection()
{
    if (!m_pluginReflection)
        return;
    if (m_dockModeData->getDockMode() == Dock::FashionMode)
    {
        m_pluginReflection->setFixedSize(m_pluginLayout->width(), REFLECTION_HEIGHT);
        m_pluginReflection->move(m_pluginLayout->x(), m_pluginLayout->y() + m_pluginLayout->height());
        m_pluginReflection->updateReflection();
    }
    else
        m_pluginReflection->setFixedSize(m_pluginLayout->width(), 0);
}

void Panel::updateLeftReflection()
{
    if (!m_appReflection)
        return;
    if (m_dockModeData->getDockMode() == Dock::FashionMode){
        m_appReflection->setFixedSize(m_appLayout->width(), 40);
        m_appReflection->move(m_appLayout->x(), m_appLayout->y() + 25);
        m_appReflection->updateReflection();
    }
    else
        m_appReflection->setFixedSize(m_appLayout->width(), 0);
}

void Panel::showPluginLayoutMask()
{
    if (!m_pluginLayoutMask){
        m_pluginLayoutMask = new LayoutDropMask(this);
        connect(m_pluginLayoutMask, &LayoutDropMask::itemDrop, [=]{
            m_pluginLayoutMask->hide();
            m_appLayout->restoreTmpItem();
        });
        connect(m_pluginLayoutMask, &LayoutDropMask::itemMove, [=]{
            //readjust position and size
            m_pluginLayoutMask->setFixedSize(m_pluginLayout->size());
            m_pluginLayoutMask->move(m_pluginLayout->pos());
        });
        connect(m_pluginLayoutMask, &LayoutDropMask::itemEnter, m_appLayout, &DockLayout::removeSpacingItem);
    }
    m_pluginLayoutMask->setFixedSize(m_pluginLayout->size());
    m_pluginLayoutMask->move(m_pluginLayout->pos());
    m_pluginLayoutMask->raise();
    m_pluginLayoutMask->show();
}

void Panel::hidePluginLayoutMask()
{
    if (m_pluginLayoutMask)
        m_pluginLayoutMask->hide();
}


void Panel::reloadStyleSheet()
{
    m_isFashionMode = m_dockModeData->getDockMode() == Dock::FashionMode;

    style()->unpolish(this);
    style()->polish(this);  // force a stylesheet recomputation
}

void Panel::showPanelMenu()
{
    QPoint tmpPos = QCursor::pos();

    PanelMenu::instance()->showMenu(tmpPos.x(),tmpPos.y());

    m_appLayout->itemHoverableChange(false);
    m_pluginLayout->itemHoverableChange(false);
}

void Panel::loadResources()
{
    m_appManager->initEntries();
    m_appLayout->setaddItemDelayInterval(500);

    QTimer::singleShot(500, m_pluginManager, SLOT(initAll()));
}

void Panel::setY(int value)
{
    move(x(), value);
}

Panel::~Panel()
{

}



LayoutDropMask::LayoutDropMask(QWidget *parent) : QFrame(parent)
{
    setAcceptDrops(true);
}

void LayoutDropMask::dragEnterEvent(QDragEnterEvent *event)
{
    emit itemEnter();
    event->setDropAction(Qt::MoveAction);
    event->accept();
}

void LayoutDropMask::dragMoveEvent(QDragMoveEvent *event)
{
    Q_UNUSED(event);
    emit itemMove();
}

void LayoutDropMask::dropEvent(QDropEvent *event)
{
    emit itemDrop();
    event->accept();
}
