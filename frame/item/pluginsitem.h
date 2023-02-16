// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLUGINSITEM_H
#define PLUGINSITEM_H

#include "dockitem.h"
#include "pluginsiteminterface.h"

class QGSettings;

class PluginsItem : public DockItem
{
    Q_OBJECT

    friend class QuickSettingController;

public:
    ~PluginsItem() override;

    int itemSortKey() const;
    void setItemSortKey(const int order) const;
    void detachPluginWidget();

    QString pluginName() const;
    PluginsItemInterface::PluginSizePolicy pluginSizePolicy() const;

    using DockItem::showContextMenu;
    using DockItem::hidePopup;

    ItemType itemType() const override;
    QSize sizeHint() const override;

    QWidget *centralWidget() const;

    virtual void setDraging(bool bDrag) override;

    PluginsItemInterface *pluginItem() const;

public slots:
    void refreshIcon() override;

protected:
    explicit PluginsItem(PluginsItemInterface *const pluginInter, const QString &itemKey, const QJsonObject &jsonData, QWidget *parent = nullptr);

private slots:
    void onGSettingsChanged(const QString &key);

protected:
    void mousePressEvent(QMouseEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;
    bool eventFilter(QObject *watched, QEvent *event) override;
    void showEvent(QShowEvent *event) override;

    void invokedMenuItem(const QString &itemId, const bool checked) override;
    void showPopupWindow(QWidget *const content, const bool model = false) override;
    const QString contextMenu() const override;
    QWidget *popupTips() override;
    void resizeEvent(QResizeEvent *event) override;

private:
    void startDrag();
    void mouseClicked();
    bool checkGSettingsControl() const;
    QString pluginApi() const;

private:
    PluginsItemInterface *const m_pluginInter;
    QWidget *m_centralWidget;

    QJsonObject m_jsonData;
    const QString m_itemKey;
    bool m_dragging;

    static QPoint MousePressPoint;
    const QGSettings *m_gsettings;
};

#endif // PLUGINSITEM_H
