#include <QHBoxLayout>
#include <QTimer>

#include "dockpanel.h"
#include "controller/dockmodedata.h"
#include "controller/plugins/dockpluginproxy.h"
#include "controller/plugins/dockpluginsmanager.h"

DockPanel::DockPanel(QWidget *parent)
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

    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &DockPanel::onDockModeChanged);
    connect(PanelMenu::instance(), &PanelMenu::menuItemInvoked, [=] {
        //To ensure that dock will not hide at changing the hide-mode to keepshowing
        m_menuItemInvoked = true;
    });
}

bool DockPanel::isFashionMode()
{
    return m_isFashionMode;
}

void DockPanel::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::RightButton)
        showPanelMenu();
}

void DockPanel::initShowHideAnimation()
{
    QStateMachine * machine = new QStateMachine(this);

    QState * showState = new QState(machine);
    showState->assignProperty(this,"y", 0);
    QState * hideState = new QState(machine);
    //y should change with DockMode changed
    connect(this, &DockPanel::startHide, [=]{
        hideState->assignProperty(this,"y", m_dockModeData->getDockHeight());
    });
    machine->setInitialState(showState);

    QPropertyAnimation *showAnimation = new QPropertyAnimation(this, "y");
    showAnimation->setDuration(SHOW_ANIMATION_DURATION);
    showAnimation->setEasingCurve(SHOW_EASINGCURVE);
    connect(showAnimation,&QPropertyAnimation::finished,this,&DockPanel::onShowPanelFinished);

    QPropertyAnimation *hideAnimation = new QPropertyAnimation(this, "y");
    hideAnimation->setDuration(HIDE_ANIMATION_DURATION);
    hideAnimation->setEasingCurve(HIDE_EASINGCURVE);
    connect(hideAnimation,&QPropertyAnimation::finished,this,&DockPanel::onHidePanelFinished);

    QSignalTransition *st = showState->addTransition(this,SIGNAL(startHide()), hideState);
    st->addAnimation(hideAnimation);
    QSignalTransition *ht = hideState->addTransition(this,SIGNAL(startShow()),showState);
    ht->addAnimation(showAnimation);

    machine->start();
}

void DockPanel::initHideStateManager()
{
    m_HSManager = new DBusHideStateManager(this);
    connect(m_HSManager,&DBusHideStateManager::ChangeState, this,&DockPanel::onHideStateChanged);

    //for initialization
    m_HSManager->UpdateState();
}

void DockPanel::initWidthAnimation()
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

void DockPanel::initPluginManager()
{
    m_pluginManager = new DockPluginsManager(this);

    connect(m_dockModeData, &DockModeData::dockModeChanged, m_pluginManager, &DockPluginsManager::onDockModeChanged);
    connect(m_pluginManager, &DockPluginsManager::itemAppend, [=](DockItem *targetItem){
        m_pluginLayout->insertWidget(0, targetItem);
        connect(targetItem, &DockItem::needPreviewShow, this, &DockPanel::onNeedPreviewShow);
        connect(targetItem, &DockItem::needPreviewHide, this, &DockPanel::onNeedPreviewHide);
        connect(targetItem, &DockItem::needPreviewUpdate, this, &DockPanel::onNeedPreviewUpdate);
    });
    connect(m_pluginManager, &DockPluginsManager::itemInsert, [=](DockItem *baseItem, DockItem *targetItem){
        int index = m_pluginLayout->indexOf(baseItem);
        m_pluginLayout->insertWidget(index != -1 ? index : m_pluginLayout->count(), targetItem);
        connect(targetItem, &DockItem::needPreviewShow, this, &DockPanel::onNeedPreviewShow);
        connect(targetItem, &DockItem::needPreviewHide, this, &DockPanel::onNeedPreviewHide);
        connect(targetItem, &DockItem::needPreviewUpdate, this, &DockPanel::onNeedPreviewUpdate);
    });
    connect(m_pluginManager, &DockPluginsManager::itemRemoved, [=](DockItem* item) {
        m_pluginLayout->removeWidget(item);
    });
    connect(PanelMenu::instance(), &PanelMenu::settingPlugin, [=]{
        m_pluginManager->onPluginsSetting(getScreenRect().height - height());
    });
}

void DockPanel::initPluginLayout()
{
    m_pluginLayout = new DockPluginLayout(this);
    m_pluginLayout->setLayoutSpacing(m_dockModeData->getAppletsItemSpacing());
    m_pluginLayout->resize(0, m_dockModeData->getItemHeight());
    connect(m_pluginLayout, &DockAppLayout::sizeChanged, this, &DockPanel::resizeWithContent);
}

void DockPanel::initAppLayout()
{
    m_appLayout = new DockAppLayout(this);
    m_appLayout->setStyleSheet("DockMovableLayout {background-color: red;}");
    m_appLayout->resize(0, m_dockModeData->getItemHeight());
    m_appLayout->setLayoutSpacing(m_dockModeData->getAppItemSpacing());
    m_appLayout->move(0, 1);
    connect(m_appLayout, &DockAppLayout::sizeChanged, this, &DockPanel::resizeWithContent);
}

void DockPanel::initAppManager()
{
    m_appManager = new DockAppManager(this);
    connect(m_appManager, &DockAppManager::entryAdded, this, &DockPanel::onAppItemAdd);
    connect(m_appManager, &DockAppManager::entryRemoved, this, &DockPanel::onAppItemRemove);

    //Make sure the item which was dragged to the dock can be show at once
//    connect(m_appLayout, &DockLayout::itemDocking, m_appManager, &AppManager::setDockingItemId);
}

void DockPanel::initReflection()
{
    if (m_appLayout)
    {
        m_appReflection = new ReflectionEffect(m_appLayout, this);
        connect(m_appLayout, &DockAppLayout::sizeChanged, this, &DockPanel::updateLeftReflection);
        connect(m_dockModeData, &DockModeData::dockModeChanged, this, &DockPanel::updateLeftReflection);
        updateLeftReflection();
    }

    if (m_pluginLayout)
    {
        m_pluginReflection = new ReflectionEffect(m_pluginLayout, this);
        connect(m_dockModeData, &DockModeData::dockModeChanged, this, &DockPanel::updateRightReflection);
        updateRightReflection();
    }
}

void DockPanel::initScreenMask()
{
    m_maskWidget = new ScreenMask();
    m_maskWidget->hide();
//    connect(m_maskWidget, &ScreenMask::itemEntered, m_appLayout, &DockLayout::removeSpacingItem);
//    connect(m_maskWidget, &ScreenMask::itemMissing, m_appLayout, &DockLayout::restoreTmpItem);
}

void DockPanel::initGlobalPreview()
{
    m_globalPreview = new PreviewWindow(DArrowRectangle::ArrowBottom);

    //make sure all app-preview will be destroy to save resources
    connect(m_globalPreview, &PreviewWindow::showFinish, [=] (QWidget *lastContent) {
        m_previewShown = true;
        if (lastContent) {
            AppPreviewsContainer *tmpFrame = qobject_cast<AppPreviewsContainer *>(lastContent);
            if (tmpFrame)
                tmpFrame->clearUpPreview();
        }
    });
    connect(m_globalPreview, &PreviewWindow::hideFinish, [=] (QWidget *lastContent) {
        m_previewShown = false;
        m_HSManager->UpdateState();
        if (lastContent) {
            AppPreviewsContainer *tmpFrame = qobject_cast<AppPreviewsContainer *>(lastContent);
            if (tmpFrame)
                tmpFrame->clearUpPreview();
        }
    });
}

void DockPanel::onItemDragStarted()
{
    m_maskWidget->show();
}

void DockPanel::resizeWithContent()
{
    if (m_dockModeData->getDockMode() == Dock::FashionMode)
    {
        m_appLayout->resize(m_appLayout->sizeHint().width() + m_dockModeData->getAppItemSpacing(), m_dockModeData->getItemHeight());
        m_pluginLayout->resize(m_pluginLayout->sizeHint().width(), m_dockModeData->getAppletsItemHeight());

        int targetWidth = FASHION_PANEL_LPADDING
                + FASHION_PANEL_RPADDING
                + m_appLayout->sizeHint().width()
                + m_pluginLayout->sizeHint().width();

        m_widthAnimation->setStartValue(width());
        m_widthAnimation->setEndValue(targetWidth);
        m_widthAnimation->start();

    }
    else
    {
        DisplayRect rec = getScreenRect();
        m_pluginLayout->resize(m_pluginLayout->sizeHint().width(), m_dockModeData->getItemHeight());
        m_pluginLayout->move(rec.width - m_pluginLayout->width(), 1);

        m_appLayout->move(0, 1);
        m_appLayout->setFixedSize(rec.width - m_pluginLayout->width(), m_dockModeData->getItemHeight());

        setFixedSize(m_appLayout->width() + m_pluginLayout->width(), m_dockModeData->getDockHeight());

        emit sizeChanged();
    }
}

void DockPanel::onAppItemAdd(DockAppItem *item, bool delayShow)
{
    m_appLayout->addWidget(item);
    connect(item, &DockAppItem::needPreviewShow, this, &DockPanel::onNeedPreviewShow);
    connect(item, &DockAppItem::needPreviewHide, this, &DockPanel::onNeedPreviewHide);
    connect(item, &DockAppItem::needPreviewUpdate, this, &DockPanel::onNeedPreviewUpdate);
}

void DockPanel::onAppItemRemove(const QString &id)
{
    QList<QWidget *> tmpList = m_appLayout->widgets();
    for (QWidget * item : tmpList)
    {
        DockAppItem *tmpItem = qobject_cast<DockAppItem *>(item);
        if (tmpItem && tmpItem->getItemId() == id)
        {
            m_appLayout->removeWidget(item);
            tmpItem->setVisible(false);
            tmpItem->deleteLater();
            return;
        }
    }
}

void DockPanel::onDockModeChanged(Dock::DockMode newMode, Dock::DockMode)
{
    m_appLayout->setLayoutSpacing(m_dockModeData->getAppItemSpacing());
    m_pluginLayout->setLayoutSpacing(m_dockModeData->getAppletsItemSpacing());

    reanchorsLayout(newMode);

    reloadStyleSheet();
}

void DockPanel::onHideStateChanged(int dockState)
{
    bool containsMouse = m_parentWidget->geometry().contains(QCursor::pos());
    if (dockState == Dock::HideStateShowing) {
        emit startShow();
    }
    else if (dockState == Dock::HideStateHiding && !containsMouse && !m_menuItemInvoked && !m_previewShown) {
        emit startHide();
    }
    else {
        m_menuItemInvoked = false;
    }
}

void DockPanel::onShowPanelFinished()
{
    m_dockModeData->setHideState(Dock::HideStateShown);
    emit panelHasShown();
}

void DockPanel::onHidePanelFinished()
{
    m_dockModeData->setHideState(Dock::HideStateHidden);
    emit panelHasHidden();
}

void DockPanel::onNeedPreviewHide(bool immediately)
{
    int interval = immediately ? 0 : DELAY_HIDE_PREVIEW_INTERVAL;
    m_globalPreview->hidePreview(interval);
}

void DockPanel::onNeedPreviewShow(QPoint pos)
{
    DockItem *item = qobject_cast<DockItem *>(sender());
    if (item && item->getApplet()) {
        m_previewShown = true;
        m_lastPreviewPos = pos;
        m_globalPreview->setArrowX(-1);//reset x to move arrow to horizontal-center
        m_globalPreview->setContent(item->getApplet());
        m_globalPreview->showPreview(pos.x(),
                                     pos.y() + m_globalPreview->shadowBlurRadius() + m_globalPreview->shadowDistance(),
                                     DELAY_SHOW_PREVIEW_INTERVAL);
    }
}

void DockPanel::onNeedPreviewUpdate()
{
    if (!m_globalPreview->isVisible())
        return;
    m_globalPreview->resizeWithContent();
    m_globalPreview->showPreview(m_lastPreviewPos.x(),
                                 m_lastPreviewPos.y() + m_globalPreview->shadowBlurRadius() + m_globalPreview->shadowDistance(),
                                 DELAY_SHOW_PREVIEW_INTERVAL);
}

void DockPanel::reanchorsLayout(Dock::DockMode mode)
{
    if (mode == Dock::FashionMode)
    {
//        m_appLayout->resize(m_appLayout->getContentsWidth() + m_dockModeData->getAppItemSpacing(),m_dockModeData->getItemHeight());
        m_appLayout->resize(m_appLayout->sizeHint().width() + m_dockModeData->getAppItemSpacing(),m_dockModeData->getItemHeight());
        m_pluginLayout->resize(m_pluginLayout->sizeHint().width(),m_dockModeData->getAppletsItemHeight());
        this->setFixedSize(FASHION_PANEL_LPADDING
                           + FASHION_PANEL_RPADDING
//                           + m_appLayout->getContentsWidth()
                           + m_appLayout->sizeHint().width()
                           + m_pluginLayout->sizeHint().width()
                           ,m_dockModeData->getDockHeight());
        m_appLayout->move(FASHION_PANEL_LPADDING,1);

        m_pluginLayout->move(m_appLayout->x() + m_appLayout->width() - m_dockModeData->getAppItemSpacing(),1);

        emit sizeChanged();
    }
    else
    {
        DisplayRect rec = getScreenRect();
        m_pluginLayout->resize(m_pluginLayout->sizeHint().width(), m_dockModeData->getItemHeight());
        m_pluginLayout->move(rec.width - m_pluginLayout->width(),1);

        m_appLayout->move(0,1);
        m_appLayout->resize(rec.width - m_pluginLayout->width() ,m_dockModeData->getItemHeight());

        this->setFixedSize(m_appLayout->width() + m_pluginLayout->width(), m_dockModeData->getDockHeight());

        emit sizeChanged();
    }
}

void DockPanel::updateRightReflection()
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

void DockPanel::updateLeftReflection()
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

void DockPanel::showPluginLayoutMask()
{
//    if (!m_pluginLayoutMask){
//        m_pluginLayoutMask = new LayoutDropMask(this);
//        connect(m_pluginLayoutMask, &LayoutDropMask::itemDrop, [=]{
//            m_pluginLayoutMask->hide();
//            m_appLayout->restoreTmpItem();
//        });
//        connect(m_pluginLayoutMask, &LayoutDropMask::itemMove, [=]{
//            //readjust position and size
//            m_pluginLayoutMask->setFixedSize(m_pluginLayout->size());
//            m_pluginLayoutMask->move(m_pluginLayout->pos());
//        });
//        connect(m_pluginLayoutMask, &LayoutDropMask::itemEnter, m_appLayout, &DockLayout::removeSpacingItem);
//    }
//    m_pluginLayoutMask->setFixedSize(m_pluginLayout->size());
//    m_pluginLayoutMask->move(m_pluginLayout->pos());
//    m_pluginLayoutMask->raise();
//    m_pluginLayoutMask->show();
}

void DockPanel::hidePluginLayoutMask()
{
//    if (m_pluginLayoutMask)
//        m_pluginLayoutMask->hide();
}


void DockPanel::reloadStyleSheet()
{
    m_isFashionMode = m_dockModeData->getDockMode() == Dock::FashionMode;

    style()->unpolish(this);
    style()->polish(this);  // force a stylesheet recomputation
}

void DockPanel::showPanelMenu()
{
    QPoint tmpPos = QCursor::pos();

    PanelMenu::instance()->showMenu(tmpPos.x(),tmpPos.y());

//    m_appLayout->itemHoverableChange(false);
//    m_pluginLayout->itemHoverableChange(false);
}

void DockPanel::loadResources()
{
    m_appManager->initEntries();
//    m_appLayout->setaddItemDelayInterval(500);

    QTimer::singleShot(500, m_pluginManager, SLOT(initAll()));
}

void DockPanel::setY(int value)
{
    move(x(), value);
}

DisplayRect DockPanel::getScreenRect()
{
    DBusDisplay d;
    return d.primaryRect();
}

DockPanel::~DockPanel()
{

}


//class LayoutDropMask : public  QFrame
//{
//    Q_OBJECT
//public:
//    LayoutDropMask(QWidget *parent = 0);

//signals:
//    void itemMove();
//    void itemEnter();
//    void itemDrop();

//protected:
//    void dragEnterEvent(QDragEnterEvent *event);
//    void dragMoveEvent(QDragMoveEvent *event);
//    void dropEvent(QDropEvent *event);

//};

//LayoutDropMask::LayoutDropMask(QWidget *parent) : QFrame(parent)
//{
//    setAcceptDrops(true);
//}

//void LayoutDropMask::dragEnterEvent(QDragEnterEvent *event)
//{
//    emit itemEnter();
//    event->setDropAction(Qt::MoveAction);
//    event->accept();
//}

//void LayoutDropMask::dragMoveEvent(QDragMoveEvent *event)
//{
//    Q_UNUSED(event);
//    emit itemMove();
//}

//void LayoutDropMask::dropEvent(QDropEvent *event)
//{
//    emit itemDrop();
//    event->accept();
//}
