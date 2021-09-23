/*
 * Copyright (C) 2011 ~ 2021 Uniontech Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng@uniontech.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng@uniontech.com>
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
#include "module_widget.h"
#include "gsetting_watcher.h"

#include <widgets/comboxwidget.h>
#include <widgets/titledslideritem.h>
#include <widgets/dccslider.h>
#include <widgets/titlelabel.h>
#include <widgets/switchwidget.h>

#include <DSlider>

#include <QLabel>
#include <QVBoxLayout>
#include <QDBusConnection>
#include <QDBusInterface>
#include <QDBusError>
#include <QMap>
#include <QScrollArea>

DWIDGET_USE_NAMESPACE

enum DisplayMode {
    Fashion     = 0,
    Efficient   = 1,
};

enum HideMode {
    KeepShowing     = 0,
    KeepHidden      = 1,
    SmartHide       = 3,
};

enum Position {
    Top         = 0,
    Right       = 1,
    Bottom      = 2,
    Left        = 3,
};

ModuleWidget::ModuleWidget(QWidget *parent)
    : QScrollArea(parent)
    , m_modeComboxWidget(new ComboxWidget)
    , m_positionComboxWidget(new ComboxWidget)
    , m_stateComboxWidget(new ComboxWidget)
    , m_sizeSlider(new TitledSliderItem(tr("Size")))
    , m_screenSettingTitle(new TitleLabel(tr("Multi screen config")))
    , m_screenSettingComboxWidget(new ComboxWidget)
    , m_pluginAreaTitle(new TitleLabel(tr("Plugin area")))
    , m_daemonDockInter(new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this))
    , m_dockInter(new DBusInter("com.deepin.dde.Dock", "/com/deepin/dde/Dock", QDBusConnection::sessionBus(), this))
    , m_gsettingsWatcher(new GSettingWatcher("com.deepin.dde.control-center", "personalization", this))
{
    initUI();
}

ModuleWidget::~ModuleWidget()
{
    delete m_modeComboxWidget;
    delete m_positionComboxWidget;
    delete m_stateComboxWidget;
    delete m_sizeSlider;
    delete m_screenSettingTitle;
    delete m_screenSettingComboxWidget;
    delete m_pluginAreaTitle;
}

void ModuleWidget::initUI()
{
    setFrameShape(QFrame::NoFrame);
    setWidgetResizable(true);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);

    QVBoxLayout *layout = new QVBoxLayout;
    layout->setContentsMargins(10, 10, 10, 10);
    layout->setSpacing(10);

    static QMap<QString, int> g_modeMap = {{tr("Fashion mode"), Fashion}
                                           , {tr("Efficient mode"), Efficient}};
    // 模式
    m_modeComboxWidget->setTitle(tr("Mode"));
    m_modeComboxWidget->addBackground();
    m_modeComboxWidget->setComboxOption(QStringList() << tr("Fashion mode") << tr("Efficient mode"));
    m_modeComboxWidget->setCurrentText(g_modeMap.key(m_daemonDockInter->displayMode()));
    connect(m_modeComboxWidget, &ComboxWidget::onSelectChanged, this, [ = ] (const QString &text) {
        m_daemonDockInter->setDisplayMode(g_modeMap.value(text));
    });
    layout->addWidget(m_modeComboxWidget);
    m_gsettingsWatcher->bind("displayMode", m_modeComboxWidget);// 转换settingName?

    static QMap<QString, int> g_positionMap = {{tr("Top"), Top}
                                               , {tr("Bottom"), Bottom}
                                               , {tr("Left"), Left}
                                               , {tr("Right"), Right}};
    // 位置
    m_positionComboxWidget->setTitle(tr("Position"));
    m_positionComboxWidget->addBackground();
    m_positionComboxWidget->setComboxOption(QStringList() << tr("Top") << tr("Bottom") << tr("Left") << tr("Right"));
    m_positionComboxWidget->setCurrentText(g_positionMap.key(m_daemonDockInter->position()));
    connect(m_positionComboxWidget, &ComboxWidget::onSelectChanged, this, [ = ] (const QString &text) {
        m_daemonDockInter->setPosition(g_positionMap.value(text));
    });
    layout->addWidget(m_positionComboxWidget);
    m_gsettingsWatcher->bind("position", m_positionComboxWidget);

    static QMap<QString, int> g_stateMap = {{tr("Always show"), KeepShowing}
                                            , {tr("Always hide"), KeepHidden}
                                            , {tr("Smart hide"), SmartHide}};
    // 状态
    m_stateComboxWidget->setTitle(tr("State"));
    m_stateComboxWidget->addBackground();
    m_stateComboxWidget->setComboxOption(QStringList() << tr("Always show") << tr("Always hide") << tr("Smart hide"));
    m_stateComboxWidget->setCurrentText(g_stateMap.key(m_daemonDockInter->hideMode()));
    connect(m_stateComboxWidget, &ComboxWidget::onSelectChanged, this, [ = ] (const QString &text) {
        m_daemonDockInter->setHideMode(g_stateMap.value(text));
    });
    layout->addWidget(m_stateComboxWidget);
    m_gsettingsWatcher->bind("hideMode", m_stateComboxWidget);

    // 高度调整控件
    m_sizeSlider->addBackground();
    m_sizeSlider->slider()->setRange(40, 100);
    QStringList ranges;
    ranges << tr("Small") << tr("Big");
    m_sizeSlider->setAnnotations(ranges);
    connect(m_daemonDockInter, &DBusDock::DisplayModeChanged, this, &ModuleWidget::updateSliderValue);
    connect(m_sizeSlider->slider(), &DSlider::valueChanged, this, [ = ] (int value) {
        if (m_daemonDockInter->displayMode() == DisplayMode::Fashion) {
            m_daemonDockInter->setWindowSizeFashion(uint(value));
        } else if (m_daemonDockInter->displayMode() == DisplayMode::Efficient) {
            m_daemonDockInter->setWindowSizeEfficient(uint(value));
        }
    });
    updateSliderValue(m_daemonDockInter->displayMode());
    m_gsettingsWatcher->bind("sizeSlider", m_sizeSlider);

    layout->addWidget(m_sizeSlider);

    // 多屏显示设置
    if (QDBusConnection::sessionBus().interface()->isServiceRegistered("com.deepin.dde.Dock")) {
        static QMap<QString, bool> g_screenSettingMap = {{tr("Follow the mouse"), false}
                                                        , {tr("Only show in primary"), true}};

        layout->addSpacing(10);
        layout->addWidget(m_screenSettingTitle);
        m_screenSettingComboxWidget->setTitle(tr("Dock position"));
        m_screenSettingComboxWidget->addBackground();
        m_screenSettingComboxWidget->setComboxOption(QStringList() << tr("Follow the mouse") << tr("Only show in primary"));
        m_screenSettingComboxWidget->setCurrentText(g_screenSettingMap.key(m_dockInter->showInPrimary()));
        connect(m_screenSettingComboxWidget, &ComboxWidget::onSelectChanged, this, [ = ] (const QString &text) {
            m_dockInter->setShowInPrimary(g_screenSettingMap.value(text));
        });
        connect(m_dockInter, &DBusInter::ShowInPrimaryChanged, m_screenSettingComboxWidget, &ComboxWidget::setCurrentIndex);
        layout->addWidget(m_screenSettingComboxWidget);
        m_gsettingsWatcher->bind("multiScreenArea", m_screenSettingTitle);
        m_gsettingsWatcher->bind("multiScreenArea", m_screenSettingComboxWidget);
    }

    // 插件区域
    QDBusPendingReply<QStringList> reply = m_dockInter->GetLoadedPlugins();
    QStringList plugins = reply.value();
    if (reply.error().type() != QDBusError::ErrorType::NoError) {
        qWarning() << "dbus call failed, method: 'GetLoadedPlugins()'";
    } else {
        if (plugins.size() != 0) {
            layout->addSpacing(10);
            layout->addWidget(m_pluginAreaTitle);
            m_gsettingsWatcher->bind("pluginArea", m_pluginAreaTitle);
            for (auto name : plugins) {
                SwitchWidget *widget = new SwitchWidget(this);
                widget->setTitle(name);
                widget->addBackground();
                widget->setChecked(m_dockInter->getPluginVisible(name));
                connect(widget, &SwitchWidget::checkedChanged, this, [ = ] (const bool checked) {
                    m_dockInter->setPluginVisible(widget->title(), checked);
                });
                connect(m_dockInter, &DBusInter::pluginVisibleChanged, this, [ = ] (const QString &pluginName, bool visible) {
                    if (pluginName == widget->title()) {
                        widget->setChecked(visible);
                    }
                });

                layout->addWidget(widget);
                m_gsettingsWatcher->bind("pluginArea", widget);
                m_pluginWidgetList.append(widget);
            }
        }
    }

    // 保持内容正常铺满
    layout->addStretch();

    // 界面内容过多时可滚动查看
    QWidget *widget = new QWidget;
    widget->setLayout(layout);
    setWidget(widget);
}

void ModuleWidget::updateSliderValue(int displayMode)
{
    if (displayMode == DisplayMode::Fashion) {
        m_sizeSlider->setValueLiteral(QString::number(m_daemonDockInter->windowSizeFashion()));
    } else if (displayMode == DisplayMode::Efficient) {
        m_sizeSlider->setValueLiteral(QString::number(m_daemonDockInter->windowSizeEfficient()));
    } else {
        Q_ASSERT_X(false, __FILE__, "not supported");
    }
}
