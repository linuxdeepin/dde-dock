#include <QHBoxLayout>

#include "panel.h"
#include "controller/dockmodedata.h"
#include "controller/plugins/dockpluginproxy.h"
#include "controller/plugins/dockpluginmanager.h"

Panel::Panel(QWidget *parent)
    : QLabel(parent),m_parentWidget(parent)
{
    setObjectName("Panel");

    initShowHideAnimation();
    initHideStateManager();
    initWidthAnimation();
    initPluginLayout();
    initAppLayout();
    initPluginManager();
    initAppManager();
    initReflection();
    initScreenMask();

    reloadStyleSheet();

    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &Panel::onDockModeChanged);
}

void Panel::setContainMouse(bool value)
{
    m_containMouse = value;
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
}

void Panel::initWidthAnimation()
{
    m_widthAnimation = new QPropertyAnimation(this, "width", this);
    m_widthAnimation->setDuration(WIDTH_ANIMATION_DURATION);
    connect(m_widthAnimation, &QPropertyAnimation::valueChanged, [=]{
        m_appLayout->move(FASHION_PANEL_LPADDING, 1);
        m_pluginLayout->move(width() - m_pluginLayout->width() - FASHION_PANEL_RPADDING, 1);
        updateRightReflection();

        this->move((m_parentWidget->width() - width()) / 2,0);
    });
}

void Panel::initPluginManager()
{
    DockPluginManager *pluginManager = new DockPluginManager(this);

    connect(m_dockModeData, &DockModeData::dockModeChanged, pluginManager, &DockPluginManager::onDockModeChanged);
    connect(pluginManager, &DockPluginManager::itemAppend, m_pluginLayout, &DockLayout::addItem);
    connect(pluginManager, &DockPluginManager::itemMove, [=](AbstractDockItem *baseItem, AbstractDockItem *targetItem){
        m_pluginLayout->moveItem(m_pluginLayout->indexOf(targetItem), m_pluginLayout->indexOf(baseItem));
    });
    connect(pluginManager, &DockPluginManager::itemInsert, [=](AbstractDockItem *baseItem, AbstractDockItem *targetItem){
        m_pluginLayout->insertItem(targetItem, m_pluginLayout->indexOf(baseItem));
    });
    connect(pluginManager, &DockPluginManager::itemRemoved, [=](AbstractDockItem* item) {
        m_pluginLayout->removeItem(item);
    });
    connect(PanelMenu::instance(), &PanelMenu::settingPlugin, [=]{
        QRect rec = QApplication::desktop()->screenGeometry();
        pluginManager->onPluginsSetting(rec.height() - height());
    });

    pluginManager->initAll();
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
    m_appLayout->setAcceptDrops(true);
    m_appLayout->setSpacing(m_dockModeData->getAppItemSpacing());
    m_appLayout->move(0, 1);

    connect(m_appLayout, &DockLayout::startDrag, this, &Panel::onItemDragStarted);
    connect(m_appLayout, &DockLayout::itemDropped, this, &Panel::onItemDropped);
    connect(m_appLayout, &DockLayout::contentsWidthChange, this, &Panel::onLayoutContentsWidthChanged);
}

void Panel::initAppManager()
{
    m_appManager = new AppManager(this);
    connect(m_appManager, &AppManager::entryAdded, this, &Panel::onAppItemAdd);
    connect(m_appManager, &AppManager::entryRemoved, this, &Panel::onAppItemRemove);
    m_appManager->updateEntries();
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
}

void Panel::onItemDropped()
{
    m_maskWidget->hide();
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
        m_pluginLayout->resize(m_pluginLayout->getContentsWidth(),m_dockModeData->getItemHeight());
        m_pluginLayout->move(m_parentWidget->width() - m_pluginLayout->width(),1);

        m_appLayout->move(0,1);
        m_appLayout->resize(m_parentWidget->width() - m_pluginLayout->width() ,m_dockModeData->getItemHeight());

        this->setFixedSize(m_appLayout->width() + m_pluginLayout->width(),m_dockModeData->getDockHeight());
        this->move((m_parentWidget->width() - m_appLayout->width() - m_pluginLayout->width()) / 2,0);
    }
}

void Panel::onAppItemAdd(AbstractDockItem *item)
{
    m_appLayout->addItem(item);
}

void Panel::onAppItemRemove(const QString &id)
{
    QList<AbstractDockItem *> tmpList = m_appLayout->getItemList();
    for (int i = 0; i < tmpList.count(); i ++)
    {
        AppItem *tmpItem = qobject_cast<AppItem *>(tmpList.at(i));
        if (tmpItem && tmpItem->getItemId() == id)
        {
            m_appLayout->removeItem(i);
            tmpItem->deleteLater();
            return;
        }
    }
}

void Panel::onDockModeChanged(Dock::DockMode newMode, Dock::DockMode)
{
    m_appLayout->relayout();
    m_pluginLayout->relayout();

    reanchorsLayout(newMode);

    reloadStyleSheet();
}

void Panel::onHideStateChanged(int dockState)
{
    if (dockState == 0)
        emit startShow();
    else if (dockState == 1 && !m_containMouse)
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
        this->move((m_parentWidget->width() - width()) / 2,0);
    }
    else
    {
        m_pluginLayout->resize(m_pluginLayout->getContentsWidth(), m_dockModeData->getItemHeight());
        m_pluginLayout->move(m_parentWidget->width() - m_pluginLayout->width(),1);

        m_appLayout->move(0,1);
        m_appLayout->resize(m_parentWidget->width() - m_pluginLayout->width() ,m_dockModeData->getItemHeight());

        this->setFixedSize(m_appLayout->width() + m_pluginLayout->width(), m_dockModeData->getDockHeight());
        this->move((m_parentWidget->width() - m_appLayout->width() - m_pluginLayout->width()) / 2,0);
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
}

void Panel::setY(int value)
{
    move(x(), value);
}

Panel::~Panel()
{

}

