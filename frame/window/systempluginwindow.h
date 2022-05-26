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
#include "dockpluginscontroller.h"

#include <QWidget>

class FixedPluginController;
class StretchPluginsItem;
class QBoxLayout;

namespace Dtk { namespace Widget { class DListView; } }

DWIDGET_USE_NAMESPACE

class SystemPluginWindow : public QWidget
{
    Q_OBJECT

public:
    explicit SystemPluginWindow(QWidget *parent = nullptr);
    ~SystemPluginWindow() override;
    void setPositon(Dock::Position position);
    QSize suitableSize();

Q_SIGNALS:
    void sizeChanged();

protected:
    void resizeEvent(QResizeEvent *event) override;

private:
    void initUi();

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

class FixedPluginController : public AbstractPluginsController
{
    Q_OBJECT

public:
    explicit FixedPluginController(QObject *parent);
    void startLoader();

Q_SIGNALS:
    void pluginItemInserted(StretchPluginsItem *);
    void pluginItemRemoved(StretchPluginsItem *);
    void pluginItemUpdated(StretchPluginsItem *);

protected:
    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemUpdate(PluginsItemInterface * const itemInter, const QString &) override;
    void itemRemoved(PluginsItemInterface * const itemInter, const QString &) override;
    bool needLoad(PluginsItemInterface *itemInter) override;

    void requestWindowAutoHide(PluginsItemInterface * const, const QString &, const bool) override {}
    void requestRefreshWindowVisible(PluginsItemInterface * const, const QString &) override {}
    void requestSetAppletVisible(PluginsItemInterface * const, const QString &, const bool) override {}

private:
    QList<StretchPluginsItem *> m_pluginItems;
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

protected:
    void paintEvent(QPaintEvent *event) override;
    void mousePressEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;

    const QString contextMenu() const override;

private:
    void mouseClick();
    QFont textFont() const;
    bool needShowText() const;

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
    Dock::Position m_position;
    QPoint m_mousePressPoint;
};

#endif // SYSTEMPLUGINWINDOW_H
