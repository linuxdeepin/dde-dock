// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MAINWINDOWBASE_H
#define MAINWINDOWBASE_H

#include "constants.h"
#include "dbusutil.h"

#include <DBlurEffectWidget>
#include <DPlatformWindowHandle>
#include <DGuiApplicationHelper>

#include <QEvent>
#include <QMouseEvent>
#include <utils.h>

class DragWidget;
class MultiScreenWorker;

DWIDGET_USE_NAMESPACE

class MainWindowBase : public DBlurEffectWidget
{
    Q_OBJECT

public:
    enum class DockWindowType {
        MainWindow,     // 主窗口
        TrayWindow      // 主窗口之外的其他窗口
    };

public:
    explicit MainWindowBase(MultiScreenWorker *multiScreenWorker, QWidget *parent = Q_NULLPTR);
    virtual ~MainWindowBase();

    void setOrder(int order);                                   // 窗体展示的顺序，按照左到右和上到下
    int order() const;

    virtual DockWindowType windowType() const = 0;
    virtual void setDisplayMode(const Dock::DisplayMode &displayMode);
    virtual void setPosition(const Dock::Position &position);
    // 用来更新子区域的位置，一般用于在执行动画的过程中，根据当前的位置来更新里面panel的大小
    virtual void updateParentGeometry(const Dock::Position &pos, const QRect &rect) = 0;
    virtual QRect getDockGeometry(QScreen *screen, const Dock::Position &pos, const Dock::DisplayMode &displaymode, const Dock::HideState &hideState, bool withoutScale = false) const;
    QVariantAnimation *createAnimation(QScreen *screen, const Dock::Position &pos, const Dock::AniAction &act);
    virtual void resetPanelGeometry() {}                        // 重置内部区域，为了让内部区域和当前区域始终保持一致
    virtual int dockSpace() const;                              // 与后面窗体之间的间隔
    virtual void serviceRestart() {}                            // 服务重新启动后的操作
    virtual void animationFinished(bool showOrHide) {}

Q_SIGNALS:
    void requestUpdate();

protected:
    void resizeEvent(QResizeEvent *event) override;
    void moveEvent(QMoveEvent *) override;
    void enterEvent(QEvent *e) override;
    void mousePressEvent(QMouseEvent *event) override;
    void showEvent(QShowEvent *event) override;

    Dock::DisplayMode displayMode() const;
    Dock::Position position() const;

    int windowSize() const;

    bool isDraging() const;

    virtual void updateRadius(int borderRadius) {}
    virtual QSize suitableSize(const Dock::Position &pos, const int &screenSize, const double &deviceRatio) const = 0;

private:
    void initUi();
    void initAttribute();
    void initConnection();
    void initMember();
    void updateDragGeometry();

    int getBorderRadius() const;
    QRect getAnimationRect(const QRect &sourceRect, const Dock::Position &pos) const;

private Q_SLOTS:
    void onMainWindowSizeChanged(QPoint offset);
    void resetDragWindow();
    void touchRequestResizeDock();
    void adjustShadowMask();
    void onCompositeChanged();
    void onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType);

protected:
    DPlatformWindowHandle m_platformWindowHandle;

private:
    Dock::DisplayMode m_displayMode;
    Dock::Position m_position;
    DragWidget *m_dragWidget;
    MultiScreenWorker *m_multiScreenWorker;
    QTimer *m_updateDragAreaTimer;
    QTimer *m_shadowMaskOptimizeTimer;
    bool m_isShow;
    int m_order;
};

#endif // MAINWINDOWBASE_H
