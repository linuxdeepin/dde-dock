/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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
#ifndef SCREENSPLITER_WAYLAND_H
#define SCREENSPLITER_WAYLAND_H

#include "screenspliter.h"

#include <QWidget>

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
    explicit ScreenSpliter_Wayland(AppItem *appItem, DockEntryInter *entryInter, QObject *parent);
    ~ScreenSpliter_Wayland() override;

    void startSplit(const QRect &rect) override;
    bool split(SplitDirection direction) override;
    bool suportSplitScreen() override;
    bool releaseSplit() override;

private:
    void setMaskVisible(const QRect &rect, bool visible);
    QString splitUuid() const;
    bool windowSupportSplit(const QString &uuid) const;
    QString firstWindowUuid() const;

private Q_SLOTS:
    void onSplitStateChange(const char* uuid, int splitable);

private:
    static SplitWindowManager *m_splitManager;
    QRect m_splitRect;
    bool m_checkedNotSupport;
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

Q_SIGNALS:
    void splitStateChange(const char *, int);

private Q_SLOTS:
    void onConnectionFinished();

private:
    ClientManagement *m_clientManagement;
    QThread *m_connectionThread;
    ConnectionThread *m_connectionThreadObject;
};

#endif // SCREENSPLITER_WAYLAND_H
