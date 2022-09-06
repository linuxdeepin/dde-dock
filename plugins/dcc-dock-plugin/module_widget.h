// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MODULE_WIDGET_H
#define MODULE_WIDGET_H

#include <QScrollArea>

#include <dtkwidget_global.h>

#include <com_deepin_dde_daemon_dock.h>

#include "com_deepin_dde_dock.h"
#include "config_watcher.h"

namespace dcc {
namespace widgets {
class ComboxWidget;
class TitledSliderItem;
}
}

DWIDGET_BEGIN_NAMESPACE
class DListView;
class DTipLabel;
DWIDGET_END_NAMESPACE

class TitleLabel;
class QStandardItemModel;

using namespace dcc::widgets;
using namespace dcc_dock_plugin;
using DBusDock = com::deepin::dde::daemon::Dock;
using DBusInter = com::deepin::dde::Dock;

class ModuleWidget : public QScrollArea
{
    Q_OBJECT
public:
    explicit ModuleWidget(QWidget *parent = nullptr);
    ~ ModuleWidget();

private:
    void initUI();
    bool isCopyMode();

private Q_SLOTS:
    void updateSliderValue();
    void updateItemCheckStatus(const QString &name, bool visible);

private:
    ComboxWidget *m_modeComboxWidget;
    ComboxWidget *m_positionComboxWidget;
    ComboxWidget *m_stateComboxWidget;

    TitledSliderItem *m_sizeSlider;

    TitleLabel *m_screenSettingTitle;
    ComboxWidget *m_screenSettingComboxWidget;

    TitleLabel *m_pluginAreaTitle;
    DTK_WIDGET_NAMESPACE::DTipLabel *m_pluginTips;
    DTK_WIDGET_NAMESPACE::DListView *m_pluginView;
    QStandardItemModel *m_pluginModel;

    DBusDock *m_daemonDockInter;
    DBusInter *m_dockInter;
    ConfigWatcher *m_dconfigWatcher;

    bool m_sliderPressed;
};

#endif // MODULE_WIDGET_H
