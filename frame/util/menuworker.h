// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MENUWORKER_H
#define MENUWORKER_H
#include <QObject>

#include "constants.h"

#include <com_deepin_dde_daemon_dock.h>

using DBusDock = com::deepin::dde::daemon::Dock;
class QMenu;
class QGSettings;
/**
 * @brief The MenuWorker class  此类用于处理任务栏右键菜单的逻辑
 */
class MenuWorker : public QObject
{
    Q_OBJECT
public:
    explicit MenuWorker(DBusDock *dockInter,QWidget *parent = nullptr);

    void showDockSettingsMenu(QMenu *menu);

signals:
    void autoHideChanged(const bool autoHide) const;

public slots:
    void setAutoHide(const bool autoHide);
    void onNotifyDaemonInterfaceUpdate(DBusDock *dockInter);

private:
    QMenu *createMenu(QMenu *settingsMenu);

private slots:
    void onDockSettingsTriggered();

private:
    DBusDock *m_dockInter;
    bool m_autoHide;
};

#endif // MENUWORKER_H
