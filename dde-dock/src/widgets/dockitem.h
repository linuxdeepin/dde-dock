/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DOCKITEM_H
#define DOCKITEM_H

#include <QDebug>
#include <QLabel>
#include <QJsonValue>
#include <QJsonObject>

#include "previewwindow.h"
#include "highlighteffect.h"
#include "dbus/dbusmenu.h"
#include "dbus/dbusmenumanager.h"
#include "interfaces/dockconstants.h"

class DockItemTitle;
class DockItem : public QFrame
{
    Q_OBJECT
public:
    explicit DockItem(QWidget *parent = 0);
    virtual ~DockItem();

    virtual QString getTitle() = 0;
    virtual QString getItemId() = 0;
    virtual QWidget * getApplet() = 0;
    virtual QString getMenuContent();
    virtual void invokeMenuItem(QString menuItemId, bool checked);
    virtual QPixmap grab(const QRect &rectangle = QRect(QPoint(0, 0), QSize(-1, -1)));

    void showMenu(const QPoint &menuPos = QPoint(0, 0));
    void showPreview(const QPoint &previewPos = QPoint(0, 0));
    void hidePreview(bool immediately = false);
    void setFixedSize(int width, int height);

    int globalX();
    int globalY();
    QPoint globalPos();

    bool hoverable() const;
    void setHoverable(bool hoverable);

signals:
    void needPreviewHide(bool immediately);
    void needPreviewShow(QPoint pos);
    void needPreviewUpdate();
    //signals for hightlight
    void mouseEnter();
    void mouseLeave();
    void mousePress();
    void mouseRelease();

protected:
    bool m_hoverable = true;
    HighlightEffect * m_highlight;
    PreviewWindow *m_titlePreview;
    DockItemTitle *m_titleLabel;
    DBusMenu * m_dbusMenu;
    DBusMenuManager * m_dbusMenuManager;

    void resizeEvent(QResizeEvent *) Q_DECL_OVERRIDE;
    void moveEvent(QMoveEvent *) Q_DECL_OVERRIDE;

private:
    void initHighlight();

};

#endif // DOCKITEM_H
