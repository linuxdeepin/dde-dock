/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             listenerri <listenerri@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef MAINPANEL_H
#define MAINPANEL_H

#include "controller/dockitemcontroller.h"
#include "util/docksettings.h"
#include "item/showdesktopitem.h"

#include <QFrame>
#include <QTimer>
#include <QBoxLayout>

#include <DBlurEffectWidget>
#include <DWindowManagerHelper>

#define xstr(s) str(s)
#define str(s) #s
#define PANEL_BORDER    0
#define PANEL_PADDING   0
#define PANEL_MARGIN    1
#define WINDOW_OVERFLOW 4

DWIDGET_USE_NAMESPACE

class MainPanel : public DBlurEffectWidget
{
    Q_OBJECT
    Q_PROPERTY(int displayMode READ displayMode DESIGNABLE true)
    Q_PROPERTY(int position READ position DESIGNABLE true)

public:
    explicit MainPanel(QWidget *parent = 0);
    virtual ~MainPanel();

    void updateDockPosition(const Position dockPosition);
    void updateDockDisplayMode(const Dock::DisplayMode displayMode);
    int displayMode() const;
    int position() const;

    void setEffectEnabled(const bool enabled);

    bool eventFilter(QObject *watched, QEvent *event);

    void setFixedSize(const QSize &size);
    void setComposite(const bool hasComposite);

signals:
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefershWindowVisible() const;
    void geometryChanged();

private:
    void moveEvent(QMoveEvent *e);
    void resizeEvent(QResizeEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void dragMoveEvent(QDragMoveEvent *e);
    void dragLeaveEvent(QDragLeaveEvent *e);
    void dropEvent(QDropEvent *e);

    void manageItem(DockItem *item);
    DockItem *itemAt(const QPoint &point);

private slots:
    void adjustItemSize();
    void itemInserted(const int index, DockItem *item);
    void itemRemoved(DockItem *item);
    void itemMoved(DockItem *item, const int index);
    void itemDragStarted();
    void itemDropped(QObject *destnation);
    void handleDragMove(QDragMoveEvent *e, bool isFilter);
    void checkMouseReallyLeave();

private:
    Position m_position;
    DisplayMode m_displayMode;

    QBoxLayout *m_itemLayout;
    QTimer *m_itemAdjustTimer;
    QTimer *m_checkMouseLeaveTimer;
    QWidget *m_appDragWidget;
    QVariantAnimation *m_sizeChangeAni;

    ShowDesktopItem *m_showDesktopItem;
    DockItemController *m_itemController;

    QString m_draggingMimeKey;
    QSize m_destSize;
};

#endif // MAINPANEL_H
