#include "panel.h"
#include "dockpluginproxy.h"
#include "dockpluginmanager.h"
#include "Controller/dockmodedata.h"
#include <QHBoxLayout>

Panel::Panel(QWidget *parent)
    : QLabel(parent),m_parentWidget(parent)
{
    this->setObjectName("Panel");

    m_pluginLayout = new DockLayout(this);
    m_pluginLayout->setSortDirection(DockLayout::RightToLeft);
    m_pluginLayout->setSpacing(m_dockModeData->getAppletsItemSpacing());
    m_pluginLayout->resize(0,m_dockModeData->getItemHeight());

    m_appLayout = new DockLayout(this);
    m_appLayout->setAcceptDrops(true);
    m_appLayout->setSpacing(m_dockModeData->getAppItemSpacing());
    m_appLayout->resize(this->width() - m_pluginLayout->width(),m_dockModeData->getItemHeight());
    m_appLayout->move(0,1);

    connect(m_appLayout, &DockLayout::dragStarted, this, &Panel::slotDragStarted);
    connect(m_appLayout, &DockLayout::itemDropped, this, &Panel::slotItemDropped);
    connect(m_appLayout, &DockLayout::contentsWidthChange, this, &Panel::slotLayoutContentsWidthChanged);
    connect(m_pluginLayout, &DockLayout::contentsWidthChange, this, &Panel::slotLayoutContentsWidthChanged);

    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &Panel::changeDockMode);

    DockPluginManager *pluginManager = new DockPluginManager(this);

    connect(m_dockModeData, &DockModeData::dockModeChanged, pluginManager, &DockPluginManager::onDockModeChanged);
    connect(pluginManager, &DockPluginManager::itemAdded, [=](AbstractDockItem* item) {
        m_pluginLayout->addItem(item);
    });
    connect(pluginManager, &DockPluginManager::itemRemoved, [=](AbstractDockItem* item) {
        int index = m_pluginLayout->indexOf(item);
        if (index != -1) {
            m_pluginLayout->removeItem(index);
        }
    });

    pluginManager->initAll();

    initAppManager();
    initHSManager();
    initState();
    initReflection();
    initScreenMask();

    updateBackground();
}

void Panel::showScreenMask()
{
//    qWarning() << "[Info:]" << "Show Screen Mask.";
    m_maskWidget->show();
}

bool Panel::isFashionMode()
{
    return m_isFashionMode;
}

void Panel::hideScreenMask()
{
//    qWarning() << "[Info:]" << "Hide Screen Mask.";
    m_maskWidget->hide();
}

void Panel::setContainMouse(bool value)
{
    m_containMouse = value;
}

void Panel::slotDragStarted()
{
    showScreenMask();
}

void Panel::slotItemDropped()
{
    hideScreenMask();
    m_appLayout->clearTmpItem();
    m_appLayout->relayout();
}

void Panel::slotEnteredMask()
{
    m_appLayout->relayout();
}

void Panel::slotExitedMask()
{
    m_appLayout->relayout();
}

void Panel::changeDockMode(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    updateBackground();

    m_appLayout->relayout();
    m_pluginLayout->relayout();

    reanchorsLayout(newMode);

    qWarning() << "AppCount:********" << m_appLayout->getItemCount();
}

void Panel::slotLayoutContentsWidthChanged()
{
    reanchorsLayout(m_dockModeData->getDockMode());
}

void Panel::slotAddAppItem(AbstractDockItem *item)
{
    m_appLayout->addItem(item);
}

void Panel::slotRemoveAppItem(const QString &id)
{
    QList<AbstractDockItem *> tmpList = m_appLayout->getItemList();
    for (int i = 0; i < tmpList.count(); i ++)
    {
        AppItem *tmpItem = qobject_cast<AppItem *>(tmpList.at(i));
        if (tmpItem && tmpItem->getItemId() == id)
        {
            m_appLayout->removeItem(i);
            return;
        }
    }
}

void Panel::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::RightButton)
        showMenu();
}

void Panel::mouseReleaseEvent(QMouseEvent *event)
{

}

void Panel::reanchorsLayout(Dock::DockMode mode)
{
    if (mode == Dock::FashionMode)
    {
        m_appLayout->resize(m_appLayout->getContentsWidth() + m_dockModeData->getAppItemSpacing(),m_dockModeData->getItemHeight());
        m_pluginLayout->setSortDirection(DockLayout::LeftToRight);
        m_pluginLayout->resize(m_pluginLayout->getContentsWidth(),m_dockModeData->getItemHeight());
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
        m_pluginLayout->setSortDirection(DockLayout::RightToLeft);
        m_pluginLayout->resize(m_pluginLayout->getContentsWidth(),m_dockModeData->getItemHeight());
        m_pluginLayout->move(m_parentWidget->width() - m_pluginLayout->width(),1);

        m_appLayout->move(0,1);
        m_appLayout->resize(m_parentWidget->width() - m_pluginLayout->width() ,m_dockModeData->getItemHeight());

        this->setFixedSize(m_appLayout->width() + m_pluginLayout->width(),m_dockModeData->getDockHeight());
        this->move((m_parentWidget->width() - m_appLayout->width() - m_pluginLayout->width()) / 2,0);
    }
}

void Panel::showMenu()
{
    QPoint tmpPos = QCursor::pos();

    PanelMenu::instance()->showMenu(tmpPos.x(),tmpPos.y());
}

void Panel::updateBackground()
{
    m_isFashionMode = m_dockModeData->getDockMode() == Dock::FashionMode;

    style()->unpolish(this);
    style()->polish(this);// force a stylesheet recomputation
}

void Panel::initAppManager()
{
    m_appManager = new AppManager(this);
    connect(m_appManager,SIGNAL(entryAdded(AbstractDockItem*)),this, SLOT(slotAddAppItem(AbstractDockItem*)));
    connect(m_appManager, SIGNAL(entryRemoved(QString)),this, SLOT(slotRemoveAppItem(QString)));
    m_appManager->updateEntries();
}

void Panel::hasShown()
{
    m_HSManager->SetState(1);
    emit panelHasShown();
}

void Panel::hasHidden()
{
    m_HSManager->SetState(3);
    emit panelHasHidden();
}

void Panel::hideStateChanged(int value)
{
    if (value == 0)
    {
        emit startShow();
    }
    else if (value == 1 && !m_containMouse)
    {
        emit startHide();
    }
}

void Panel::initHSManager()
{
    m_HSManager = new DBusHideStateManager(this);
    connect(m_HSManager,&DBusHideStateManager::ChangeState,this,&Panel::hideStateChanged);
}

void Panel::initState()
{
    QStateMachine * machine = new QStateMachine(this);
    QState * showState = new QState(machine);
    showState->assignProperty(this,"y", 0);
    QState * hideState = new QState(machine);
    hideState->assignProperty(this,"y", height());
    machine->setInitialState(showState);

    QPropertyAnimation *sa = new QPropertyAnimation(this, "y");
    sa->setDuration(SHOW_HIDE_DURATION);
    sa->setEasingCurve(SHOW_HIDE_EASINGCURVE);
    connect(sa,&QPropertyAnimation::finished,this,&Panel::hasShown);
    QPropertyAnimation *ha = new QPropertyAnimation(this, "y");
    ha->setDuration(SHOW_HIDE_DURATION);
    ha->setEasingCurve(SHOW_HIDE_EASINGCURVE);
    connect(ha,&QPropertyAnimation::finished,this,&Panel::hasHidden);

    QSignalTransition *ts1 = showState->addTransition(this,SIGNAL(startHide()), hideState);
    ts1->addAnimation(ha);
    connect(ts1,&QSignalTransition::triggered,[=]{
        m_HSManager->SetState(2);
    });
    QSignalTransition *ts2 = hideState->addTransition(this,SIGNAL(startShow()),showState);
    ts2->addAnimation(sa);
    connect(ts2,&QSignalTransition::triggered,[=]{
        m_HSManager->SetState(0);
    });

    machine->start();
}

void Panel::initReflection()
{
    if (m_appLayout)
    {
        m_appReflection = new ReflectionEffect(m_appLayout, this);
        connect(m_appLayout, &DockLayout::frameUpdate, this, &Panel::updateLeftReflection);
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
    connect(m_maskWidget,SIGNAL(itemDropped(QPoint)),this,SLOT(slotItemDropped()));
    connect(m_maskWidget,SIGNAL(itemEntered()),this,SLOT(slotEnteredMask()));
    connect(m_maskWidget,SIGNAL(itemExited()),this,SLOT(slotExitedMask()));
    connect(m_maskWidget, &ScreenMask::itemMissing, m_appLayout, &DockLayout::restoreTmpItem);
}

void Panel::updateLeftReflection()
{
    if (m_dockModeData->getDockMode() == Dock::FashionMode){
        m_appReflection->setFixedSize(m_appLayout->width(), REFLECTION_HEIGHT);
        m_appReflection->move(m_appLayout->x(), m_appLayout->y() + m_appLayout->height());
        m_appReflection->updateReflection();
    }
    else
        m_appReflection->setFixedSize(m_appLayout->width(), 0);
}

void Panel::updateRightReflection()
{
    if (m_dockModeData->getDockMode() == Dock::FashionMode)
    {
        m_pluginReflection->setFixedSize(m_pluginLayout->width(), REFLECTION_HEIGHT);
        m_pluginReflection->move(m_pluginLayout->x(), m_pluginLayout->y() + m_pluginLayout->height());
        m_pluginReflection->updateReflection();
    }
    else
        m_pluginReflection->setFixedSize(m_pluginLayout->width(), 0);
}

void Panel::setY(int value)
{
    move(x(), value);
}

Panel::~Panel()
{

}

