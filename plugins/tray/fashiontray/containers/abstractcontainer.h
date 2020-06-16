/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:
 *
 * Maintainer:  zhaolong <zhaolong@uniontech.com>
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

#ifndef ABSTRACTCONTAINER_H
#define ABSTRACTCONTAINER_H

#include "constants.h"
#include "../../trayplugin.h"
#include "../fashiontraywidgetwrapper.h"

#include <QWidget>

class AbstractContainer : public QWidget
{
    Q_OBJECT
public:
    explicit AbstractContainer(TrayPlugin *trayPlugin, QWidget *parent = nullptr);

    virtual bool acceptWrapper(FashionTrayWidgetWrapper *wrapper) = 0;
    virtual void refreshVisible();

    virtual void addWrapper(FashionTrayWidgetWrapper *wrapper);
    virtual bool removeWrapper(FashionTrayWidgetWrapper *wrapper);
    virtual bool removeWrapperByTrayWidget(AbstractTrayWidget *trayWidget);
    virtual FashionTrayWidgetWrapper *takeWrapper(FashionTrayWidgetWrapper *wrapper);
    virtual void setDockPosition(const Dock::Position pos);
    virtual void setExpand(const bool expand);
    virtual QSize totalSize() const;
    virtual int itemCount();

    int itemSize() {return m_itemSize;}
    void setItemSize(int itemSize);
    void clearWrapper();
    void saveCurrentOrderToConfig();
    bool isEmpty();
    bool containsWrapper(FashionTrayWidgetWrapper *wrapper);
    bool containsWrapperByTrayWidget(AbstractTrayWidget *trayWidget);
    FashionTrayWidgetWrapper *wrapperByTrayWidget(AbstractTrayWidget *trayWidget);

    void addDraggingWrapper(FashionTrayWidgetWrapper *wrapper);
    FashionTrayWidgetWrapper *takeDraggingWrapper();

Q_SIGNALS:
    void attentionChanged(FashionTrayWidgetWrapper *wrapper, const bool attention);
    void requestDraggingWrapper();
    void draggingStateChanged(FashionTrayWidgetWrapper *wrapper, const bool dragging);

protected:
    virtual int whereToInsert(FashionTrayWidgetWrapper *wrapper);

    TrayPlugin *trayPlugin() const;
    QList<QPointer<FashionTrayWidgetWrapper>> wrapperList() const;
    QBoxLayout *wrapperLayout() const;
    void setWrapperLayout(QBoxLayout *layout);
    bool expand() const;
    Dock::Position dockPosition() const;
//    QSize wrapperSize() const;

protected:
    void dragEnterEvent(QDragEnterEvent *event) override;
    void paintEvent(QPaintEvent *event) override;

private Q_SLOTS:
    void onWrapperAttentionhChanged(const bool attention);
    void onWrapperDragStart();
    void onWrapperDragStop();
    void onWrapperRequestSwapWithDragging();

private:
    TrayPlugin *m_trayPlugin;
    QBoxLayout *m_wrapperLayout;

    QPointer<FashionTrayWidgetWrapper> m_currentDraggingWrapper;
    QList<QPointer<FashionTrayWidgetWrapper>> m_wrapperList;

    bool m_expand;
    Dock::Position m_dockPosition;

//    QSize m_wrapperSize;
    int m_itemSize = 40;
};

#endif // ABSTRACTCONTAINER_H
