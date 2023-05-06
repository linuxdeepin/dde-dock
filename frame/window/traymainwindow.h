// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
