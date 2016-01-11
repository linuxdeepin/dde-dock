#include <QTimer>
#include <QHBoxLayout>
#include <QPushButton>

#include "dockpanel.h"
#include "controller/dockmodedata.h"
#include "controller/old/pluginproxy.h"

//const int REFLECTION_HEIGHT = 15;
const int FASHION_PANEL_LPADDING = 21;
const int FASHION_PANEL_RPADDING = 21;
//const int WIDTH_ANIMATION_DURATION = 200;
const int SHOW_ANIMATION_DURATION = 300;
const int HIDE_ANIMATION_DURATION = 300;
const int DELAY_HIDE_PREVIEW_INTERVAL = 200;
const int DELAY_SHOW_PREVIEW_INTERVAL = 200;
const QEasingCurve SHOW_EASINGCURVE = QEasingCurve::OutCubic;
const QEasingCurve HIDE_EASINGCURVE = QEasingCurve::Linear;

DockPanel::DockPanel(QWidget *parent)
    : QLabel(parent)
{
    setObjectName("Panel");

    initGlobalPreview();
    initShowHideAnimation();
    initHideStateManager();
    initPluginLayout();
    initAppLayout();

    initMainLayout();

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

void DockPanel::initPluginLayout()
{
    m_pluginLayout = new DockPluginLayout(this);
    m_pluginLayout->setAutoResize(true);
    m_pluginLayout->resize(0, m_dockModeData->getAppletsItemHeight());
    m_pluginLayout->setLayoutSpacing(m_dockModeData->getAppletsItemSpacing());

    connect(m_pluginLayout, &DockPluginLayout::sizeChanged, this, &DockPanel::onContentsSizeChanged);
    connect(m_pluginLayout, &DockPluginLayout::needPreviewShow, this, &DockPanel::onNeedPreviewShow);
    connect(m_pluginLayout, &DockPluginLayout::needPreviewHide, this, &DockPanel::onNeedPreviewHide);
    connect(m_pluginLayout, &DockPluginLayout::needPreviewUpdate, this, &DockPanel::onNeedPreviewUpdate);
}

void DockPanel::initAppLayout()
{
    m_appLayout = new DockAppLayout(this);
    m_appLayout->setAutoResize(m_dockModeData->getDockMode() == Dock::FashionMode);
    m_appLayout->resize(0, m_dockModeData->getItemHeight());
    m_appLayout->setLayoutSpacing(m_dockModeData->getAppItemSpacing());

    connect(m_appLayout, &DockPluginLayout::sizeChanged, this, &DockPanel::onContentsSizeChanged);
    connect(m_appLayout, &DockAppLayout::needPreviewShow, this, &DockPanel::onNeedPreviewShow);
    connect(m_appLayout, &DockAppLayout::needPreviewHide, this, &DockPanel::onNeedPreviewHide);
    connect(m_appLayout, &DockAppLayout::needPreviewUpdate, this, &DockPanel::onNeedPreviewUpdate);
}

void DockPanel::initMainLayout()
{
    QHBoxLayout *mLayout = new QHBoxLayout(this);
    mLayout->setSpacing(0);
    mLayout->setContentsMargins(0, 0, 0, 0);
    mLayout->addWidget(m_appLayout, 0, Qt::AlignTop);
    mLayout->addWidget(m_pluginLayout, 0, Qt::AlignTop);

    //for init
    onDockModeChanged(m_dockModeData->getDockMode(), m_dockModeData->getDockMode());
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

void DockPanel::onDockModeChanged(Dock::DockMode, Dock::DockMode)
{
    reloadStyleSheet();

    m_pluginLayout->setLayoutSpacing(m_dockModeData->getAppletsItemSpacing());
    m_pluginLayout->setFixedHeight(m_pluginLayout->sizeHint().height());
    QHBoxLayout *mLayout = qobject_cast<QHBoxLayout *>(layout());
    if (m_dockModeData->getDockMode() == Dock::FashionMode) {
        mLayout->setAlignment(m_pluginLayout, Qt::AlignTop);
    }
    else {
        mLayout->setAlignment(m_pluginLayout, Qt::AlignVCenter);
    }
}

void DockPanel::onHideStateChanged(int dockState)
{
    bool containsMouse = parentWidget()->geometry().contains(QCursor::pos());
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

void DockPanel::onContentsSizeChanged()
{
    if (m_dockModeData->getDockMode() == Dock::FashionMode) {
        m_appLayout->setAutoResize(true);
        m_appLayout->update();
    }
    else {
        DisplayRect rec = getScreenRect();
        m_appLayout->setAutoResize(false);
        m_appLayout->setFixedSize(rec.width - m_pluginLayout->width() , m_dockModeData->getItemHeight());
    }

    emit sizeChanged();
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
    m_appLayout->initEntries();
//    m_appLayout->setaddItemDelayInterval(500);
    m_pluginLayout->initAllPlugins();
}

QSize DockPanel::sizeHint() const
{
    int w = m_appLayout->width() + m_pluginLayout->width() ;
    int h = m_appLayout->height() + m_pluginLayout->height();
    if (m_dockModeData->getDockMode() == Dock::FashionMode) {
        w = w + FASHION_PANEL_LPADDING + FASHION_PANEL_RPADDING;
    }

    return QSize(w, h);
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
