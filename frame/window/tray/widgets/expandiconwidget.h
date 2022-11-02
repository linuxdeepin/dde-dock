/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
    const QString dropIconFile() const;

private:
    Dtk::Gui::DRegionMonitor *m_regionInter;
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
    void setOwnerWidget(QWidget *widget);
    TrayGridView *trayView() const;
    void resetPosition();

protected:
    void paintEvent(QPaintEvent *event) override;

private:
    QColor maskColor() const;

private:
    DockInter *m_dockInter;
    TrayGridView *m_trayGridView;
    static Dock::Position m_position;
    QWidget *m_ownerWidget;
};

#endif // EXPANDICONWIDGET_H
