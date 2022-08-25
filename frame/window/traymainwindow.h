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
#ifndef TRAYMAINWINDOW_H
#define TRAYMAINWINDOW_H

#include "constants.h"
#include "mainwindowbase.h"

#include <DBlurEffectWidget>

class TrayManagerWindow;
class MultiScreenWorker;

DWIDGET_USE_NAMESPACE

class TrayMainWindow : public MainWindowBase
{
    Q_OBJECT

public:
    TrayMainWindow(MultiScreenWorker *multiScreenWorker, QWidget *parent = Q_NULLPTR);
    void setPosition(const Dock::Position &position) override;
    TrayManagerWindow *trayManagerWindow() const;

    void setDisplayMode(const Dock::DisplayMode &displayMode) override;
    DockWindowType windowType() const override;
    void updateParentGeometry(const Dock::Position &position, const QRect &rect) override;
    QSize suitableSize(const Dock::Position &pos, const int &, const double &) const override;
    QSize suitableSize() const;
    void resetPanelGeometry() override;

protected:
    int dockSpace() const override;
    void updateRadius(int borderRadius) override;

private:
    void initUI();
    void initConnection();

private Q_SLOTS:
    void onRequestUpdate();

private:
    TrayManagerWindow *m_trayManager;
    MultiScreenWorker *m_multiScreenWorker;
};

#endif // TRAYMAINWINDOW_H
