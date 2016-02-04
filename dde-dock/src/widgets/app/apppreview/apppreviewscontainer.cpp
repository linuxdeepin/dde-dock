/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QPushButton>
#include <QHBoxLayout>
#include <QApplication>
#include <QDesktopWidget>
#include <QDebug>

#include "apppreviewscontainer.h"
#include "apppreviewloaderframe.h"
#include "interfaces/dockconstants.h"

const int PREVIEW_PADDING = 5;
const int BUTTON_SIZE = Dock::APP_PREVIEW_CLOSEBUTTON_SIZE;
const int PREVIEW_HEIGHT = Dock::APP_PREVIEW_HEIGHT;
const int PREVIEW_WIDTH = Dock::APP_PREVIEW_WIDTH;
const int PREVIEW_MARGIN = Dock::APP_PREVIEW_MARGIN;

AppPreviewsContainer::AppPreviewsContainer(QWidget *parent) :
    QWidget(parent), m_isClosing(false)
{
    m_clientManager  = new DBusClientManager(this);


    m_mainLayout = new QHBoxLayout(this);
    m_mainLayout->setAlignment(Qt::AlignCenter);
    m_mainLayout->setContentsMargins(Dock::APP_PREVIEW_CLOSEBUTTON_SIZE / 2, 0, 0, 0);
    m_mainLayout->setSpacing(PREVIEW_MARGIN - Dock::APP_PREVIEW_CLOSEBUTTON_SIZE / 2);

    resize(PREVIEW_WIDTH,PREVIEW_HEIGHT);
}

void AppPreviewsContainer::addItem(const QString &title, int xid)
{
    if (m_previewMap.keys().indexOf(xid) != -1)
        return;

    AppPreviewLoaderFrame *f = new AppPreviewLoaderFrame(title, xid, this);
    connect(f, &AppPreviewLoaderFrame::requestPreviewClose, this, &AppPreviewsContainer::removePreview);
    connect(f, &AppPreviewLoaderFrame::requestPreviewActive, this, &AppPreviewsContainer::activatePreview);

    m_mainLayout->addWidget(f);

    m_previewMap.insert(xid, f);

    setItemCount(m_previewMap.count());
}

void AppPreviewsContainer::leaveEvent(QEvent *)
{
    if (m_isClosing)
        m_isClosing = false;
}

void AppPreviewsContainer::removePreview(int xid)
{
    m_isClosing = true;

    m_previewMap.remove(xid);
    m_mainLayout->removeWidget(qobject_cast<AppPreviewLoaderFrame *>(sender()));
    sender()->deleteLater();
    m_clientManager->CloseWindow(xid);

    if (m_mainLayout->count() <= 0)
    {
        emit requestHide();
        return;
    }

    setItemCount(m_previewMap.count());

    emit sizeChanged();
}

void AppPreviewsContainer::activatePreview(int xid)
{
    m_clientManager->ActiveWindow(xid);

    emit requestHide();
}

void AppPreviewsContainer::clearUpPreview()
{
    QLayoutItem *child;
    while ((child = m_mainLayout->takeAt(0)) != 0){
        child->widget()->deleteLater();
        delete child;
    }

    m_previewMap.clear();
}

QSize AppPreviewsContainer::getNormalContentSize()
{

    int contentWidth = m_mainLayout->count() * (PREVIEW_WIDTH + PREVIEW_MARGIN)
            + PREVIEW_MARGIN - PREVIEW_PADDING * 2;
    int contentHeight = PREVIEW_HEIGHT + PREVIEW_MARGIN*2 - PREVIEW_PADDING * 2;

    return QSize(contentWidth, contentHeight);
}

void AppPreviewsContainer::setItemCount(int count)
{
    QSize frameSize(PREVIEW_WIDTH + BUTTON_SIZE / 2, PREVIEW_HEIGHT + BUTTON_SIZE);
    QRect dr = QApplication::desktop()->geometry();
    bool outOfScreen = getNormalContentSize().width() > dr.width();

    //if the total width larger than screen width,scale the preview frame size
    if (outOfScreen) {
        int w = (dr.width() - (count + 1) * PREVIEW_MARGIN) / count + BUTTON_SIZE / 2;
        int h = w * (PREVIEW_HEIGHT + BUTTON_SIZE) / (PREVIEW_WIDTH + BUTTON_SIZE / 2);
        frameSize = QSize(w, h);
    }

    foreach (AppPreviewLoaderFrame *frame, m_previewMap.values()) {
        frame->shrink(frameSize, outOfScreen);
    }

    int contentWidth = count * (frameSize.width() - BUTTON_SIZE / 2) + (count + 1) * PREVIEW_MARGIN;
    int contentHeight = PREVIEW_HEIGHT + PREVIEW_MARGIN * 2;

    resize(contentWidth, contentHeight);
}

AppPreviewsContainer::~AppPreviewsContainer()
{

}





