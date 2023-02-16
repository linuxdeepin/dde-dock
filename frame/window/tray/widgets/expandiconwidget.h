// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef EXPANDICONWIDGET_H
#define EXPANDICONWIDGET_H

#include "constants.h"
#include "basetraywidget.h"
#include "dbusutil.h"

class TrayGridView;
class TrayModel;
class TrayDelegate;
class TrayGridWidget;

namespace Dtk { namespace Gui { class DRegionMonitor; } }

class ExpandIconWidget : public BaseTrayWidget
{
    Q_OBJECT

public:
    explicit ExpandIconWidget(QWidget *parent = Q_NULLPTR, Qt::WindowFlags f = Qt::WindowFlags());
    ~ExpandIconWidget() override;
    void setPositon(Dock::Position position);

    void sendClick(uint8_t mouseButton, int x, int y) override;
    void setTrayPanelVisible(bool visible);
    QString itemKeyForConfig() override { return "Expand"; }
    void updateIcon() override {}
    QPixmap icon() override;
    static TrayGridWidget *popupTrayView();

protected:
    void paintEvent(QPaintEvent *event) override;
    void moveEvent(QMoveEvent *event) override;
    const QString dropIconFile() const;

private:
    Dock::Position m_position;
};

// 绘制圆角窗体
class TrayGridWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TrayGridWidget(QWidget *parent);

    static void setPosition(const Dock::Position &position);
    void setTrayGridView(TrayGridView *trayView);
    void setReferGridView(TrayGridView *trayView);
    TrayGridView *trayView() const;
    void resetPosition();

protected:
    void paintEvent(QPaintEvent *event) override;
    void showEvent(QShowEvent *event) override;
    void hideEvent(QHideEvent *event) override;

private:
    void initMember();
    QColor maskColor() const;
    ExpandIconWidget *expandWidget() const;

private:
    DockInter *m_dockInter;
    TrayGridView *m_trayGridView;
    TrayGridView *m_referGridView;
    Dtk::Gui::DRegionMonitor *m_regionInter;
    static Dock::Position m_position;
};

#endif // EXPANDICONWIDGET_H
