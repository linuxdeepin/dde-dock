/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#ifndef SYSTEMPLUGINWINDOW_H
#define SYSTEMPLUGINWINDOW_H

#include "constants.h"
#include "dockitem.h"

#include <QWidget>

class FixedPluginController;
class StretchPluginsItem;
class QBoxLayout;
class PluginsItemInterface;

namespace Dtk { namespace Widget { class DListView; } }

DWIDGET_USE_NAMESPACE

class SystemPluginWindow : public QWidget
{
    Q_OBJECT

public:
    explicit SystemPluginWindow(QWidget *parent = nullptr);
    ~SystemPluginWindow() override;
    void setPositon(Dock::Position position);
    QSize suitableSize() const;
    QSize suitableSize(const Dock::Position &position) const;

Q_SIGNALS:
    void itemChanged();

private:
    void initUi();
    bool pluginExist(StretchPluginsItem *pluginItem);

private Q_SLOTS:
    void onPluginItemAdded(StretchPluginsItem *pluginItem);
    void onPluginItemRemoved(StretchPluginsItem *pluginItem);
    void onPluginItemUpdated(StretchPluginsItem *pluginItem);

private:
    FixedPluginController *m_pluginController;
    DListView *m_listView;
    Dock::Position m_position;
    QBoxLayout *m_mainLayout;
};

class StretchPluginsItem : public DockItem
{
    Q_OBJECT

public:
    StretchPluginsItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~StretchPluginsItem() override;
    void setPosition(Dock::Position position);
    PluginsItemInterface *pluginInter() const;
    QString itemKey() const;
    QSize suitableSize() const;
    QSize suitableSize(const Dock::Position &position) const;

    inline ItemType itemType() const override { return DockItem::StretchPlugin; }

protected:
    void paintEvent(QPaintEvent *event) override;
    void mousePressEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;

    const QString contextMenu() const override;
    void invokedMenuItem(const QString &itemId, const bool checked) override;

private:
    void mouseClick();
    QFont textFont() const;
    QFont textFont(const Dock::Position &position) const;
    bool needShowText() const;

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
    Dock::Position m_position;
    QPoint m_mousePressPoint;
};

#endif // SYSTEMPLUGINWINDOW_H
