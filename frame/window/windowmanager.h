// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef WINDOWMANAGER_H
#define WINDOWMANAGER_H

#include "constants.h"
#include "statusnotifierwatcher_interface.h"

#include <QObject>

namespace Dtk {namespace Gui { class DWindowManagerHelper; }}

class MainWindowBase;
class MainWindow;
class TrayMainWindow;
class MultiScreenWorker;
class MenuWorker;
class QDBusConnectionInterface;

using namespace Dtk::Gui;

class WindowManager : public QObject
{
    Q_OBJECT

public:
    explicit WindowManager(MultiScreenWorker *multiScreenWorker, QObject *parent = nullptr);
    ~WindowManager() override;

    void addWindow(MainWindowBase *window);
    void launch();
    void sendNotifications();
    void callShow();
    void resizeDock(int offset, bool dragging);
    QRect geometry() const;

Q_SIGNALS:
    void panelGeometryChanged();

private:
    void initConnection();
    void initSNIHost();
    void initMember();
    void updateMainGeometry(const Dock::HideState &hideState);
    QParallelAnimationGroup *createAnimationGroup(const Dock::AniAction &aniAction, const QString &screenName, const Dock::Position &position) const;

    void showAniFinish();
    void animationFinish(bool showOrHide);
    void hideAniFinish();
    QRect getDockGeometry(bool withoutScale = false) const;         // 计算左右侧加起来的区域大小

    void RegisterDdeSession();
    void updateDockGeometry(const QRect &rect);

private Q_SLOTS:
    void onUpdateDockGeometry(const Dock::HideMode &hideMode);
    void onPositionChanged(const Dock::Position &position);
    void onDisplayModeChanged(const Dock::DisplayMode &displayMode);

    void onDbusNameOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner);

    void onPlayAnimation(const QString &screenName, const Dock::Position &pos, Dock::AniAction act, bool containMouse = false, bool updatePos = false);
    void onChangeDockPosition(QString fromScreen, QString toScreen, const Dock::Position &fromPos, const Dock::Position &toPos);
    void onRequestUpdateFrontendGeometry();
    void onRequestNotifyWindowManager();

    void onServiceRestart();

private:
    MultiScreenWorker *m_multiScreenWorker;
    QString m_sniHostService;

    Dock::DisplayMode m_displayMode;
    Dock::Position m_position;

    QDBusConnectionInterface *m_dbusDaemonInterface;
    org::kde::StatusNotifierWatcher *m_sniWatcher;      // DBUS状态通知
    QList<MainWindowBase *> m_topWindows;
};

#endif // WINDOWMANAGER_H
