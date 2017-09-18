/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include <QFrame>
#include <QTimer>
#include <QBoxLayout>

#include <DBlurEffectWidget>
#include <DWindowManagerHelper>

#define xstr(s) str(s)
#define str(s) #s
#define PANEL_BORDER    0
#define PANEL_PADDING   6

DWIDGET_USE_NAMESPACE

class MainPanel : public QFrame
{
    Q_OBJECT
    Q_PROPERTY(int displayMode READ displayMode DESIGNABLE true)
    Q_PROPERTY(int position READ position DESIGNABLE true)

public:
    explicit MainPanel(QWidget *parent = 0);

    void updateDockPosition(const Position dockPosition);
    void updateDockDisplayMode(const Dock::DisplayMode displayMode);
    int displayMode() const;
    int position() const;

signals:
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefershWindowVisible() const;
    void geometryChanged();

private:
    void moveEvent(QMoveEvent *e);
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void dragMoveEvent(QDragMoveEvent *e);
    void dragLeaveEvent(QDragLeaveEvent *e);
    void dropEvent(QDropEvent *e);

    void manageItem(DockItem *item);
    DockItem *itemAt(const QPoint &point);

    void updateBlurEffect() const;

private slots:
    void adjustItemSize();
    void itemInserted(const int index, DockItem *item);
    void itemRemoved(DockItem *item);
    void itemMoved(DockItem *item, const int index);
    void itemDragStarted();
    void itemDropped(QObject *destnation);

private:
    Position m_position;
    DisplayMode m_displayMode;
    QBoxLayout *m_itemLayout;

    QTimer *m_itemAdjustTimer;
    DockItemController *m_itemController;

    QTimer *m_updateEffectTimer;
    DBlurEffectWidget *m_effectWidget;
    DWindowManagerHelper *m_wmHelper;

    static DockItem *DragingItem;
    static PlaceholderItem *RequestDockItem;
};

#endif // MAINPANEL_H
