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
#ifndef TRAYMANAGERWINDOW_H
#define TRAYMANAGERWINDOW_H

#include "constants.h"
#include "dbusutil.h"

#include "org_deepin_dde_timedate1.h"

#include <QWidget>

namespace Dtk { namespace Gui { class DRegionMonitor; };
                namespace Widget { class DBlurEffectWidget; } }

using namespace Dtk::Widget;

using Timedate = org::deepin::dde::Timedate1;

class QuickPluginWindow;
class QBoxLayout;
class TrayGridView;
class TrayModel;
class TrayDelegate;
class SystemPluginWindow;
class QLabel;
class QDropEvent;
class DateTimeDisplayer;
class QPainterPath;

class TrayManagerWindow : public QWidget
{
    Q_OBJECT

public:
    explicit TrayManagerWindow(QWidget *parent = nullptr);
    ~TrayManagerWindow() override;

    void updateBorderRadius(int borderRadius);
    void updateLayout();
    void setPositon(Dock::Position position);
    void setDisplayMode(Dock::DisplayMode displayMode);
    QSize suitableSize() const;
    QSize suitableSize(const Dock::Position &position) const;

Q_SIGNALS:
    void requestUpdate();

protected:
    void resizeEvent(QResizeEvent *event) override;
    void dragEnterEvent(QDragEnterEvent *e) override;
    void dragMoveEvent(QDragMoveEvent *e) override;
    void dropEvent(QDropEvent *e) override;
    void dragLeaveEvent(QDragLeaveEvent *event) override;
    void paintEvent(QPaintEvent *event) override;

private:
    void initUi();
    void initConnection();

    void resetChildWidgetSize();
    void resetMultiDirection();
    void resetSingleDirection();
    QColor maskColor(uint8_t alpha) const;

    int appDatetimeSize(const Dock::Position &position) const;
    QPainterPath roundedPaths();
    void updateItemLayout(int dockSize);

private Q_SLOTS:
    void onTrayCountChanged();

private:
    QWidget *m_appPluginDatetimeWidget;
    DockInter *m_dockInter;
    SystemPluginWindow *m_systemPluginWidget;
    QWidget *m_appPluginWidget;
    QuickPluginWindow *m_quickIconWidget;
    DateTimeDisplayer *m_dateTimeWidget;
    QBoxLayout *m_appPluginLayout;
    QBoxLayout *m_mainLayout;
    TrayGridView *m_trayView;
    TrayModel *m_model;
    TrayDelegate *m_delegate;
    Dock::Position m_position;
    Dock::DisplayMode m_displayMode;
    QLabel *m_splitLine;
    bool m_singleShow;                              // 用于记录当前日期时间和插件区域是显示一行还是显示多行
    int m_borderRadius;                             // 圆角的值
};

#endif // PLUGINWINDOW_H
