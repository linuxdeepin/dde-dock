/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QTimer>
#include <QHBoxLayout>
#include <QPushButton>

#include "dockpanel.h"
#include "controller/dockmodedata.h"
#include "controller/stylemanager.h"

const int REFLECTION_Y = 25;
const int REFLECTION_HEIGHT = 40;
const int FASHION_PANEL_LPADDING = 21;
const int FASHION_PANEL_RPADDING = 21;
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

    initHideStateManager();
    initGlobalPreview();
    initShowHideAnimation();
    initPluginLayout();
    initAppLayout();
    initMainLayout();
    initReflection();

    setMinimumHeight(m_dockModeData->getDockHeight());  //set height for border-image calculate
    reloadStyleSheet();

    connect(m_dockModeData, &DockModeData::dockModeChanged, this, &DockPanel::onDockModeChanged);
    connect(this, &DockPanel::sizeChanged, this, &DockPanel::updateReflection);
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
//    QStateMachine * machine = new QStateMachine(this);

//    QState * showState = new QState(machine);
//    showState->assignProperty(this,"y", 0);
//    QState * hideState = new QState(machine);
    //y should change with DockMode changed
//    connect(this, &DockPanel::startHide, [&]{
//        hideState->assignProperty(this,"y", m_dockModeData->getDockHeight());
//    });
//    machine->setInitialState(showState);

    QPropertyAnimation *showAnimation = new QPropertyAnimation(this, "y");
    showAnimation->setDuration(SHOW_ANIMATION_DURATION);
    showAnimation->setEasingCurve(SHOW_EASINGCURVE);
    connect(showAnimation,&QPropertyAnimation::finished,this,&DockPanel::onShowPanelFinished);
    connect(showAnimation, &QPropertyAnimation::stateChanged, this, &DockPanel::changeItemHoverable);

    QPropertyAnimation *hideAnimation = new QPropertyAnimation(this, "y");
    hideAnimation->setDuration(HIDE_ANIMATION_DURATION);
    hideAnimation->setEasingCurve(HIDE_EASINGCURVE);
    connect(hideAnimation,&QPropertyAnimation::finished,this,&DockPanel::onHidePanelFinished);
    connect(hideAnimation, &QPropertyAnimation::stateChanged, this, &DockPanel::changeItemHoverable);

    connect(this, &DockPanel::startHide, [=] {
        hideAnimation->setStartValue(0);
        hideAnimation->setEndValue(m_dockModeData->getDockHeight());
        hideAnimation->start();
    });

    connect(this, &DockPanel::startShow, [=] {
        showAnimation->setStartValue(m_dockModeData->getDockHeight());
        showAnimation->setEndValue(0);
        showAnimation->start();
    });

//    connect(this, &DockPanel::startHide, [hideAnimation] {hideAnimation->start();});

//    QSignalTransition *st = showState->addTransition(this,SIGNAL(startHide()), hideState);
//    st->addAnimation(hideAnimation);
//    connect(this, &DockPanel::startShow, [this] {qDebug() << "Ccc";});
//    QSignalTransition *ht = hideState->addTransition(this,SIGNAL(startShow()), showState);
//    ht->addAnimation(showAnimation);

//    machine->start();
}

void DockPanel::initHideStateManager()
{
    m_HSManager = new DBusHideStateManager(this);
    connect(m_HSManager,&DBusHideStateManager::ChangeState, this, &DockPanel::onHideStateChanged);

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
    connect(m_pluginLayout, &DockPluginLayout::pluginsInitDone, this, &DockPanel::pluginsInitDone, Qt::QueuedConnection);
}

void DockPanel::initAppLayout()
{
    m_appLayout = new DockAppLayout(this);
    m_appLayout->setAutoResize(m_dockModeData->getDockMode() == Dock::FashionMode);
    m_appLayout->resize(0, m_dockModeData->getItemHeight());
    m_appLayout->setLayoutSpacing(m_dockModeData->getAppItemSpacing());

    connect(m_appLayout, &DockAppLayout::sizeChanged, this, &DockPanel::onContentsSizeChanged);
    connect(m_appLayout, &DockAppLayout::needPreviewShow, this, &DockPanel::onNeedPreviewShow);
    connect(m_appLayout, &DockAppLayout::needPreviewHide, this, &DockPanel::onNeedPreviewHide);
    connect(m_appLayout, &DockAppLayout::needPreviewUpdate, this, &DockPanel::onNeedPreviewUpdate);
}

void DockPanel::initMainLayout()
{
    QHBoxLayout *mLayout = new QHBoxLayout(this);
    mLayout->setSpacing(0);
    mLayout->setContentsMargins(0, 0, 0, 0);
    m_launcherItem = new DockLauncherItem();
    mLayout->addWidget(m_launcherItem, 0, Qt::AlignTop);
    mLayout->addSpacing(m_dockModeData->getAppItemSpacing());
    mLayout->addWidget(m_appLayout, 0, Qt::AlignTop);
    mLayout->addSpacing(8);
    mLayout->addWidget(m_pluginLayout, 0, Qt::AlignTop);

    //for init
    onDockModeChanged(m_dockModeData->getDockMode(), m_dockModeData->getDockMode());
}

void DockPanel::initReflection()
{
    m_launcherReflection = new ReflectionEffect(m_launcherItem, this);
    m_pluginReflection = new ReflectionEffect(m_pluginLayout, this);
    m_appReflection = new ReflectionEffect(m_appLayout, this);
}

void DockPanel::initGlobalPreview()
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
        m_HSManager->UpdateState();
        if (lastContent) {
            AppPreviewsContainer *tmpFrame = qobject_cast<AppPreviewsContainer *>(lastContent);
            if (tmpFrame)
                tmpFrame->clearUpPreview();
        }
    });
    connect(m_globalPreview, &PreviewWindow::previewFrameHided, m_HSManager, &DBusHideStateManager::UpdateState);
}

void DockPanel::onDockModeChanged(Dock::DockMode, Dock::DockMode)
{
    reloadStyleSheet();

    m_pluginLayout->setLayoutSpacing(m_dockModeData->getAppletsItemSpacing());
    m_pluginLayout->setFixedHeight(m_dockModeData->getItemHeight());
    QHBoxLayout *mLayout = qobject_cast<QHBoxLayout *>(layout());
    if (m_dockModeData->getDockMode() == Dock::FashionMode) {
        mLayout->setAlignment(m_pluginLayout, Qt::AlignTop);
        m_pluginLayout->setAlignment(Qt::AlignTop);
    }
    else {
        mLayout->setAlignment(m_pluginLayout, Qt::AlignVCenter);
        m_pluginLayout->setAlignment(Qt::AlignVCenter);
    }

    // interval 0 stands for timeout will be triggered on idle.
    QTimer::singleShot(0, m_appLayout, &DockAppLayout::updateItemWidths);
    QTimer::singleShot(0, m_appLayout, &DockAppLayout::updateWindowIconGeometries);
}

void DockPanel::onHideStateChanged(int dockState)
{
    bool containsMouse = parentWidget()->geometry().contains(QCursor::pos());

    if (dockState == Dock::HideStateShowing) {
        emit startShow();
    } else if (dockState == Dock::HideStateHiding && !containsMouse && !m_globalPreview->isVisible()) {
        emit startHide();
    }
}

void DockPanel::onShowPanelFinished()
{
    //dbus的ToggleShow接口会在判断时把HideStateShown对应的切换到HideStateShowing导致一直没法再切换
    m_dockModeData->setHideState(Dock::HideStateHiding);
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

void DockPanel::onNeedPreviewShow(DockItem *item, const QPoint &pos)
{
    if (item && item->getApplet()) {
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

        m_appLayout->setFixedSize(rec.width - m_pluginLayout->width() - m_launcherItem->width(), m_dockModeData->getItemHeight());
    }
    m_appLayout->updateWindowIconGeometries();

    setFixedSize(sizeHint().width(), m_dockModeData->getDockHeight());
    emit sizeChanged();
}

void DockPanel::changeItemHoverable(QAbstractAnimation::State state)
{
    bool v = true;
    switch (state) {
    case QAbstractAnimation::Running:
        v = false;
        break;
    case QAbstractAnimation::Paused:
    case QAbstractAnimation::Stopped:
        v = true;
        break;
    default:
        break;
    }

    m_launcherItem->setHoverable(v);
    m_appLayout->itemHoverableChange(v);
    m_pluginLayout->itemHoverableChange(v);
}

void DockPanel::reloadStyleSheet()
{
    m_isFashionMode = m_dockModeData->getDockMode() == Dock::FashionMode;

    // INFO: 这里用 unpolish/polish 会出现只有部分 qss 起作用的情况，必须重新加载所有 qss
    StyleManager::instance()->initStyleSheet();

//    style()->unpolish(this);
//    style()->polish(this);  // force a stylesheet recomputation
}

void DockPanel::showPanelMenu()
{
    QPoint tmpPos = QCursor::pos();

    PanelMenu::instance()->showMenu(tmpPos.x(),tmpPos.y());
}

void DockPanel::updateReflection()
{
    if (m_dockModeData->getDockMode() == Dock::FashionMode) {
        m_launcherReflection->setVisible(true);
        m_launcherReflection->setFixedSize(m_launcherItem->width(), REFLECTION_HEIGHT);
        m_launcherReflection->move(m_launcherItem->x(), REFLECTION_Y);
        m_launcherReflection->updateReflection();

        m_appReflection->setVisible(true);
        m_appReflection->setFixedSize(m_appLayout->width(), REFLECTION_HEIGHT);
        m_appReflection->move(m_appLayout->x(), REFLECTION_Y);
        m_appReflection->updateReflection();

        m_pluginReflection->setVisible(true);
        m_pluginReflection->setFixedSize(m_pluginLayout->width(), REFLECTION_HEIGHT);
        m_pluginReflection->move(m_pluginLayout->x(), REFLECTION_Y);
        m_pluginReflection->updateReflection();
    }
    else {
        m_launcherReflection->setVisible(false);
        m_pluginReflection->setVisible(false);
        m_appReflection->setVisible(false);
    }
}

void DockPanel::loadResources()
{
    m_appLayout->initEntries();
    m_pluginLayout->initAllPlugins();
}

QSize DockPanel::sizeHint() const
{
    int w = m_appLayout->width() + m_pluginLayout->width() + m_launcherItem->width();
    int h = m_appLayout->height() + m_pluginLayout->height() + m_launcherItem->height();
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
