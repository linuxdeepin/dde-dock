// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SCREENSPLITER_WAYLAND_H
#define SCREENSPLITER_WAYLAND_H

#include "screenspliter.h"

#include <QWidget>
#include <QMap>

namespace KWayland {
namespace Client {
class Registry;
class DDEShell;
class DDEShellSurface;
class EventQueue;
class Compositor;
class Surface;
class ClientManagement;
class ConnectionThread;
}
}

class AppItem;
class QWindow;
class QThread;
class SplitWindowManager;

class WindowInfo;
typedef QMap<quint32, WindowInfo> WindowInfoMap;

using namespace KWayland::Client;

class ScreenSpliter_Wayland : public ScreenSpliter
{
    Q_OBJECT

public:
    explicit ScreenSpliter_Wayland(AppItem *appItem, QObject *parent);
    ~ScreenSpliter_Wayland() override;

    void startSplit(const QRect &rect) override;
    bool split(SplitDirection direction) override;
    bool suportSplitScreen() override;
    bool releaseSplit() override;

private:
    void setMaskVisible(const QRect &rect, bool visible);
    bool windowSupportSplit(const QString &uuid) const;

private:
    static SplitWindowManager *m_splitManager;
    QRect m_splitRect;
};

class SplitWindowManager : public QObject
{
    Q_OBJECT

    friend class ScreenSpliter_Wayland;

protected:
    explicit SplitWindowManager(QObject *parent = Q_NULLPTR);
    ~SplitWindowManager() override;

    bool canSplit(const QString &uuid) const;
    void requestSplitWindow(const char *uuid, const ScreenSpliter::SplitDirection &direction);

private Q_SLOTS:
    void onConnectionFinished();

private:
    ClientManagement *m_clientManagement;
    QThread *m_connectionThread;
    ConnectionThread *m_connectionThreadObject;
    QMap<QString, int> m_splitWindowState;
};

#endif // SCREENSPLITER_WAYLAND_H
