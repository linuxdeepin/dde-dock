// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKSCREEN_H
#define DOCKSCREEN_H

#include <QString>

/**
 * @brief The DockScreen class
 * 保存任务栏的屏幕信息
 */
class DockScreen
{
public:
    static DockScreen *instance();

    const QString &current() const;
    const QString &last() const;
    const QString &primary() const;
    void updateDockedScreen(const QString &screenName);
    void updatePrimary(const QString &primary);

private:
    explicit DockScreen();

private:
    QString m_primary;
    QString m_currentScreen;
    QString m_lastScreen;
};

#endif // DOCKSCREEN_H
