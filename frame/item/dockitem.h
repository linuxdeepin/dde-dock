// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKITEM_H
#define DOCKITEM_H

#include "constants.h"
#include "dockpopupwindow.h"

#include <QFrame>
#include <QPointer>
#include <QGestureEvent>

#include <memory>

using namespace Dock;

class QMenu;

class DockItem : public QWidget
{
    Q_OBJECT

public:
    enum ItemType {
        Launcher,               // 启动器
        App,                    // 任务栏区域的应用
        Plugins,                // 插件区域图标
        FixedPlugin,            // 固定区域图标，例如多任务试图
        Placeholder,
        TrayPlugin,             // 托盘插件
        QuickSettingPlugin,     // 快捷设置区域插件
        StretchPlugin,          // 时尚模式下的固定在最右侧的插件，例如开关机插件
        AppMultiWindow          // APP的多开应用的窗口
    };

public:
    explicit DockItem(QWidget *parent = nullptr);
    ~DockItem() override;

    static void setDockPosition(const Position side);
    static void setDockDisplayMode(const DisplayMode mode);

    inline virtual ItemType itemType() const = 0;

    QSize sizeHint() const override;
    virtual QString accessibleName();

public slots:
    virtual void refreshIcon() {}

    void showPopupApplet(QWidget *const applet);
    void hidePopup();
    virtual void setDraging(bool bDrag);
    virtual void checkEntry() {}

    bool isDragging();
signals:
    void dragStarted() const;
    void itemDropped(QObject *destination, const QPoint &dropPoint) const;
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefreshWindowVisible() const;

protected:
    bool event(QEvent *event) override;
    void paintEvent(QPaintEvent *e) override;
    void mousePressEvent(QMouseEvent *e) override;
    void enterEvent(QEvent *e) override;
    void leaveEvent(QEvent *e) override;

    const QRect perfectIconRect() const;
    const QPoint popupMarkPoint() ;
    const QPoint topleftPoint() const;

    void hideNonModel();
    void popupWindowAccept();
    virtual void showPopupWindow(QWidget *const content, const bool model = false);
    virtual void showHoverTips();
    virtual void invokedMenuItem(const QString &itemId, const bool checked);
    virtual const QString contextMenu() const;
    virtual QWidget *popupTips();

    bool checkAndResetTapHoldGestureState();
    virtual void gestureEvent(QGestureEvent *event);

protected slots:
    void showContextMenu();
    void onContextMenuAccepted();

private:
    void updatePopupPosition();
    void menuActionClicked(QAction *action);

protected:
    bool m_hover;
    bool m_popupShown;
    bool m_tapAndHold;
    bool m_draging;
    QMenu *m_contextMenu;

    QPointer<QWidget> m_lastPopupWidget;

    QTimer *m_popupTipsDelayTimer;
    QTimer *m_popupAdjustDelayTimer;

    static Position DockPosition;
    static DisplayMode DockDisplayMode;
    static QPointer<DockPopupWindow> PopupWindow;
};

#endif // DOCKITEM_H
