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

#ifndef PLUGINSITEM_H
#define PLUGINSITEM_H

#include "dockitem.h"
#include "pluginsiteminterface.h"

class PluginsItem : public DockItem
{
    Q_OBJECT

public:
    explicit PluginsItem(PluginsItemInterface* const pluginInter, const QString &itemKey, QWidget *parent = 0);
    ~PluginsItem();

    int itemSortKey() const;
    void setItemSortKey(const int order) const;
    void detachPluginWidget();

    bool allowContainer() const;
    bool isInContainer() const;
    void setInContainer(const bool container);

    using DockItem::showContextMenu;
    using DockItem::hidePopup;

    inline ItemType itemType() const override {return Plugins;}
    QSize sizeHint() const override;

public slots:
    void refershIcon() override;

private:
    void mousePressEvent(QMouseEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;
    bool eventFilter(QObject *o, QEvent *e) override;

    void invokedMenuItem(const QString &itemId, const bool checked) override;
    void showPopupWindow(QWidget * const content, const bool model = false) override;
    const QString contextMenu() const override;
    QWidget *popupTips() override;

private:
    void startDrag();
    void mouseClicked();

private:
    PluginsItemInterface * const m_pluginInter;
    QWidget *m_centralWidget;
    const QString m_itemKey;
    bool m_draging;

    static QPoint MousePressPoint;
};

#endif // PLUGINSITEM_H
