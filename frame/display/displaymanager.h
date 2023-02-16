// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DISPLAYMANAGER_H
#define DISPLAYMANAGER_H

#include <QObject>

#include "singleton.h"
#include "constants.h"
#include "org_deepin_dde_display1.h"

using DisplayInter = org::deepin::dde::Display1;
using namespace Dock;

class QScreen;
class QTimer;
class QGSettings;
class DisplayManager: public QObject, public Singleton<DisplayManager>
{
    Q_OBJECT
    friend class Singleton<DisplayManager>;

public:
    explicit DisplayManager(QObject *parent = Q_NULLPTR);

    QList<QScreen *> screens() const;
    QScreen *screen(const QString &screenName) const;
    QScreen *screenAt(const QPoint &pos) const;
    QString primary() const;
    int screenRawWidth() const;
    int screenRawHeight() const;
    bool canDock(QScreen *s, Position pos) const;
    bool isCopyMode();

private:
    void updateScreenDockInfo();

private Q_SLOTS:
    void screenCountChanged();
    void dockInfoChanged();
    void onGSettingsChanged(const QString &key);

Q_SIGNALS:
    void primaryScreenChanged();
    void screenInfoChanged();       // 屏幕信息发生变化，需要调整任务栏显示，只需要这一个信号，其他的都不要，简化流程

private:
    QList<QScreen *> m_screens;
    QMap<QScreen *, QMap<Position, bool>> m_screenPositionMap;
    const QGSettings *m_gsettings;              // 多屏配置控制
    bool m_onlyInPrimary;
};

#endif // DISPLAYMANAGER_H
