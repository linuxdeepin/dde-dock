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
