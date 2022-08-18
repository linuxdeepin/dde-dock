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
#ifndef MULTIWINDOWHELPER_H
#define MULTIWINDOWHELPER_H

#include "constants.h"

#include <QObject>

class AppMultiItem;

class MultiWindowHelper : public QObject
{
    Q_OBJECT

public:
    explicit MultiWindowHelper(QWidget *appWidget, QWidget *multiWindowWidget, QObject *parent = nullptr);

    void setDisplayMode(Dock::DisplayMode displayMode);
    void addMultiWindow(int, AppMultiItem *item);
    void removeMultiWindow(AppMultiItem *item);

Q_SIGNALS:
    void requestUpdate();

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    int itemIndex(AppMultiItem *item);
    void insertChildWidget(QWidget *parentWidget, int index, AppMultiItem *item);
    void resetMultiItemPosition();

private:
    QWidget *m_appWidget;
    QWidget *m_multiWindowWidget;
    Dock::DisplayMode m_displayMode;
};

#endif // MULTIWINDOWHELPER_H
