// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SYSTEMPLUGINWINDOW_H
#define SYSTEMPLUGINWINDOW_H

#include "constants.h"
#include "dockitem.h"
#include "dbusutil.h"

#include <QWidget>

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
    void setDisplayMode(const Dock::DisplayMode &displayMode);
    void setPositon(Dock::Position position);
    QSize suitableSize() const;
    QSize suitableSize(const Dock::Position &position) const;

Q_SIGNALS:
    void itemChanged();
    void requestDrop(QDropEvent *dropEvent);
    void requestDrawBackground(const QRect &rect);

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    void initUi();
    void initConnection();
    StretchPluginsItem *findPluginItemWidget(PluginsItemInterface *pluginItem);
    void pluginAdded(PluginsItemInterface *plugin);
    QList<StretchPluginsItem *> stretchItems() const;

private Q_SLOTS:
    void onPluginItemRemoved(PluginsItemInterface *pluginItem);
    void onPluginItemUpdated(PluginsItemInterface *pluginItem);

private:
    DListView *m_listView;
    Dock::DisplayMode m_displayMode;
    Dock::Position m_position;
    QBoxLayout *m_mainLayout;
};

class StretchPluginsItem : public DockItem
{
    Q_OBJECT

public:
    StretchPluginsItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~StretchPluginsItem() override;
    void setDisplayMode(const Dock::DisplayMode &displayMode);
    static void setPosition(Dock::Position position);
    PluginsItemInterface *pluginInter() const;
    QString itemKey() const;
    QSize suitableSize() const;
    QSize suitableSize(const Dock::Position &position) const;

    inline ItemType itemType() const override { return DockItem::StretchPlugin; }

protected:
    void paintEvent(QPaintEvent *event) override;
    void mousePressEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;

    const QString contextMenu() const override;
    void invokedMenuItem(const QString &itemId, const bool checked) override;

    QWidget *popupTips() override;

private:
    void mouseClick();
    QFont textFont() const;
    QFont textFont(const Dock::Position &position) const;
    bool needShowText() const;

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
    Dock::DisplayMode m_displayMode;
    static Dock::Position m_position;
    QPoint m_mousePressPoint;
    uint m_windowSizeFashion;
};

#endif // SYSTEMPLUGINWINDOW_H
