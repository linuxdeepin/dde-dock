// Copyright (C) 2018 ~ 2020 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MENUWORKER_H
#define MENUWORKER_H

#include "constants.h"
#include "dbusutil.h"

#include <QObject>

class QMenu;
class QGSettings;
/**
 * @brief The MenuWorker class  此类用于处理任务栏右键菜单的逻辑
 */
class MenuWorker : public QObject
{
    Q_OBJECT
public:
    explicit MenuWorker(QObject *parent = nullptr);

    void exec();

private:
    void createMenu(QMenu *settingsMenu);

private slots:
    void onDockSettingsTriggered();

private:
    Dock::DisplayMode m_displaymode;
    Dock::Position m_position;
    Dock::HideMode m_hideMode;
};

#endif // MENUWORKER_H
