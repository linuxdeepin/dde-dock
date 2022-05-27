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

#include <QWidget>

#include <com_deepin_daemon_timedate.h>

namespace Dtk { namespace Gui { class DRegionMonitor; };
                namespace Widget { class DBlurEffectWidget; } }

using namespace Dtk::Widget;

using Timedate = com::deepin::daemon::Timedate;

class QuickPluginWindow;
class QBoxLayout;
class TrayGridView;
class TrayModel;
class TrayDelegate;
class SystemPluginWindow;
class QLabel;
class QDropEvent;
class DateTimeDisplayer;

class TrayManagerWindow : public QWidget
{
    Q_OBJECT

public:
    explicit TrayManagerWindow(QWidget *parent = nullptr);
    ~TrayManagerWindow() override;

    void updateLayout();
    void setPositon(Dock::Position position);
    QSize suitableSize();

Q_SIGNALS:
    void sizeChanged();

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

    int appDatetimeSize();
    QPainterPath roundedPaths();

private:
    QWidget *m_appPluginDatetimeWidget;
    SystemPluginWindow *m_systemPluginWidget;
    QWidget *m_appPluginWidget;
    QuickPluginWindow *m_quickIconWidget;
    DateTimeDisplayer *m_dateTimeWidget;
    QBoxLayout *m_appPluginLayout;
    QBoxLayout *m_appDatetimeLayout;
    QBoxLayout *m_mainLayout;
    TrayGridView *m_trayView;
    TrayModel *m_model;
    TrayDelegate *m_delegate;
    Dock::Position m_postion;
    QLabel *m_splitLine;
};

#endif // PLUGINWINDOW_H
