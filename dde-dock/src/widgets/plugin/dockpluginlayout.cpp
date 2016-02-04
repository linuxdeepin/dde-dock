/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "dockpluginlayout.h"
#include "../../panel/panelmenu.h"
#include "../../controller/dockmodedata.h"

DockPluginLayout::DockPluginLayout(QWidget *parent) : MovableLayout(parent)
{
    setAcceptDrops(false);
    initPluginManager();
}

QSize DockPluginLayout::sizeHint() const
{
    QSize size;
    int w = 0;
    int h = 0;
    switch (direction()) {
    case QBoxLayout::LeftToRight:
    case QBoxLayout::RightToLeft:
        size.setHeight(DockModeData::instance()->getDockHeight());
        for (QWidget * widget : widgets()) {
            w += widget->sizeHint().width();
        }
        size.setWidth(w + getLayoutSpacing() * widgets().count());
        break;
    case QBoxLayout::TopToBottom:
    case QBoxLayout::BottomToTop:
        size.setWidth(DockModeData::instance()->getAppletsItemWidth());
        for (QWidget * widget : widgets()) {
            h += widget->height();
        }
        size.setHeight(h + getLayoutSpacing() * widgets().count());
        break;
    }

    return size;
}

void DockPluginLayout::initAllPlugins()
{
    QTimer::singleShot(500, m_pluginManager, SLOT(initAll()));
}

void DockPluginLayout::initPluginManager()
{
    m_pluginManager = new DockPluginsManager(this);

    connect(m_pluginManager, &DockPluginsManager::itemAppend, [=](DockItem *targetItem){
        this->insertWidget(0, targetItem);
        connect(targetItem, &DockItem::needPreviewShow, [=](QPoint pos) {
            DockItem *s = qobject_cast<DockItem *>(sender());
            if (s)
                emit needPreviewShow(s, pos);
        });
        connect(targetItem, &DockItem::needPreviewHide, this, &DockPluginLayout::needPreviewHide);
        connect(targetItem, &DockItem::needPreviewUpdate, this, &DockPluginLayout::needPreviewUpdate);
        connect(this, &DockPluginLayout::itemHoverableChange, targetItem, &DockItem::setHoverable);
    });
    connect(m_pluginManager, &DockPluginsManager::itemInsert, [=](DockItem *baseItem, DockItem *targetItem){
        int index = indexOf(baseItem);
        insertWidget(index != -1 ? index : count(), targetItem);
        connect(targetItem, &DockItem::needPreviewShow, this, [=](QPoint pos) {
            DockItem *s = qobject_cast<DockItem *>(sender());
            if (s)
                emit needPreviewShow(s, pos);
        });
        connect(targetItem, &DockItem::needPreviewHide, this, &DockPluginLayout::needPreviewHide);
        connect(targetItem, &DockItem::needPreviewUpdate, this, &DockPluginLayout::needPreviewUpdate);
        connect(this, &DockPluginLayout::itemHoverableChange, targetItem, &DockItem::setHoverable);
    });
    connect(m_pluginManager, &DockPluginsManager::itemRemoved, [=](DockItem* item) {
        removeWidget(item);
    });
    connect(PanelMenu::instance(), &PanelMenu::settingPlugin, [=]{
        m_pluginManager->onPluginsSetting(getScreenRect().height - parentWidget()->height());
    });
}

DisplayRect DockPluginLayout::getScreenRect()
{
    DBusDisplay d;
    return d.primaryRect();
}

