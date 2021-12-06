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
#include "utils.h"

#include <widgets/comboxwidget.h>
#include <widgets/titledslideritem.h>
#include <widgets/dccslider.h>
#include <widgets/titlelabel.h>

#include <DSlider>
#include <DListView>
#include <DTipLabel>

#include <QApplication>
#include <QScreen>
#include <QLabel>
#include <QVBoxLayout>
#include <QDBusConnection>
#include <QDBusInterface>
#include <QDBusError>
#include <QMap>
#include <QScrollArea>
#include <QScroller>
#include <QComboBox>
#include <QTimer>

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
    , m_modeComboxWidget(new ComboxWidget(this))
    , m_positionComboxWidget(new ComboxWidget(this))
    , m_stateComboxWidget(new ComboxWidget(this))
    , m_sizeSlider(new TitledSliderItem(tr("Size"), this))
    , m_screenSettingTitle(new TitleLabel(tr("Multiple Displays"), this))
    , m_screenSettingComboxWidget(new ComboxWidget(this))
    , m_pluginAreaTitle(new TitleLabel(tr("Plugin Area"), this))
    , m_pluginTips(new DTipLabel(tr("Select which icons appear in the Dock"), this))
    , m_pluginView(new DListView(this))
    , m_pluginModel(new QStandardItemModel(this))
    , m_daemonDockInter(new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this))
    , m_dockInter(new DBusInter("com.deepin.dde.Dock", "/com/deepin/dde/Dock", QDBusConnection::sessionBus(), this))
    , m_dconfigWatcher(new ConfigWatcher("dde.dock.plugin.dconfig", this))
    , m_sliderPressed(false)
{
    // 异步，否则频繁调用可能会导致卡顿
    m_daemonDockInter->setSync(false);
    initUI();

    connect(m_dockInter, &DBusInter::pluginVisibleChanged, this, &ModuleWidget::updateItemCheckStatus);
}

ModuleWidget::~ModuleWidget()
{
}

void ModuleWidget::initUI()
{
    setBackgroundRole(QPalette::Base);
    setFrameShape(QFrame::NoFrame);
    setWidgetResizable(true);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);

    QVBoxLayout *layout = new QVBoxLayout;
    layout->setContentsMargins(10, 10, 10, 10);
    layout->setSpacing(10);

    static QMap<QString, int> g_modeMap = {{tr("Fashion mode"), Fashion}
                                           , {tr("Efficient mode"), Efficient}};
    // 模式
    if (Utils::SettingValue("com.deepin.dde.dock.module.menu", QByteArray(), "modeVisible", true).toBool()) {
        m_modeComboxWidget->setTitle(tr("Mode"));
        m_modeComboxWidget->addBackground();
        m_modeComboxWidget->setComboxOption(QStringList() << tr("Fashion mode") << tr("Efficient mode"));
        m_modeComboxWidget->setCurrentText(g_modeMap.key(m_daemonDockInter->displayMode()));
        connect(m_modeComboxWidget, &ComboxWidget::onSelectChanged, this, [ = ] (const QString &text) {
            m_daemonDockInter->setDisplayMode(g_modeMap.value(text));
        });
        connect(m_daemonDockInter, &DBusDock::DisplayModeChanged, this, [ = ] (int value) {
            DisplayMode mode = static_cast<DisplayMode>(value);
            if (g_modeMap.key(mode) == m_modeComboxWidget->comboBox()->currentText())
                return;

            m_modeComboxWidget->setCurrentText(g_modeMap.key(mode));
        });
        layout->addWidget(m_modeComboxWidget);
        m_dconfigWatcher->bind("Control-Center_Dock_Model", m_modeComboxWidget);
    } else {
        m_modeComboxWidget->setVisible(false);
    }

    if (Utils::SettingValue("com.deepin.dde.dock.module.menu", QByteArray(), "locationVisible", true).toBool()) {
        // 位置
        static QMap<QString, int> g_positionMap = {{tr("Top"), Top}
                                                   , {tr("Bottom"), Bottom}
                                                   , {tr("Left"), Left}
                                                   , {tr("Right"), Right}};

        m_positionComboxWidget->setTitle(tr("Location"));
        m_positionComboxWidget->addBackground();
        m_positionComboxWidget->setComboxOption(QStringList() << tr("Top") << tr("Bottom") << tr("Left") << tr("Right"));
        m_positionComboxWidget->setCurrentText(g_positionMap.key(m_daemonDockInter->position()));
        connect(m_positionComboxWidget, &ComboxWidget::onSelectChanged, this, [ = ] (const QString &text) {
            m_daemonDockInter->setPosition(g_positionMap.value(text));
        });
        connect(m_daemonDockInter, &DBusDock::PositionChanged, this, [ = ] (int position) {
            Position pos = static_cast<Position>(position);
            if (g_positionMap.key(pos) == m_positionComboxWidget->comboBox()->currentText())
                return;

            m_positionComboxWidget->setCurrentText(g_positionMap.key(pos));
        });
        layout->addWidget(m_positionComboxWidget);
        m_dconfigWatcher->bind("Control-Center_Dock_Location", m_positionComboxWidget);
    } else {
        m_positionComboxWidget->setVisible(false);
    }

    // 状态
    if (Utils::SettingValue("com.deepin.dde.dock.module.menu", QByteArray(), "statusVisible", true).toBool()) {
        static QMap<QString, int> g_stateMap = {{tr("Keep shown"), KeepShowing}
                                                , {tr("Keep hidden"), KeepHidden}
                                                , {tr("Smart hide"), SmartHide}};

        m_stateComboxWidget->setTitle(tr("Status"));
        m_stateComboxWidget->addBackground();
        m_stateComboxWidget->setComboxOption(QStringList() << tr("Keep shown") << tr("Keep hidden") << tr("Smart hide"));
        m_stateComboxWidget->setCurrentText(g_stateMap.key(m_daemonDockInter->hideMode()));
        connect(m_stateComboxWidget, &ComboxWidget::onSelectChanged, this, [ = ] (const QString &text) {
            m_daemonDockInter->setHideMode(g_stateMap.value(text));
        });
        connect(m_daemonDockInter, &DBusDock::HideModeChanged, this, [ = ] (int value) {
            HideMode mode = static_cast<HideMode>(value);
            if (g_stateMap.key(mode) == m_stateComboxWidget->comboBox()->currentText())
                return;

            m_stateComboxWidget->setCurrentText(g_stateMap.key(mode));
        });
        layout->addWidget(m_stateComboxWidget);
        m_dconfigWatcher->bind("Control-Center_Dock_State", m_stateComboxWidget);
    } else {
        m_stateComboxWidget->setVisible(false);
    }

    // 高度调整控件
    m_sizeSlider->addBackground();
    m_sizeSlider->slider()->setRange(40, 100);
    QStringList ranges;
    ranges << tr("Small") << "" << tr("Large");
    m_sizeSlider->setAnnotations(ranges);
    connect(m_daemonDockInter, &DBusDock::DisplayModeChanged, this, &ModuleWidget::updateSliderValue);
    connect(m_daemonDockInter, &DBusDock::WindowSizeFashionChanged, this, &ModuleWidget::updateSliderValue);
    connect(m_daemonDockInter, &DBusDock::WindowSizeEfficientChanged, this, &ModuleWidget::updateSliderValue);
    connect(m_sizeSlider->slider(), &DSlider::sliderMoved, m_sizeSlider->slider(), &DSlider::valueChanged);
    connect(m_sizeSlider->slider(), &DSlider::valueChanged, this, [ = ] (int value) {
        m_dockInter->resizeDock(value, true);
    });
    connect(m_sizeSlider->slider(), &DSlider::sliderPressed, m_dockInter, [ = ] {
        m_daemonDockInter->blockSignals(true);
        m_sliderPressed = true;
    });
    connect(m_sizeSlider->slider(), &DSlider::sliderReleased, m_dockInter, [ = ] {
        m_daemonDockInter->blockSignals(false);
        m_sliderPressed = false;

        // 松开手后通知dock拖拽状态接触
        QTimer::singleShot(0, this, [ = ] {
            int offset = m_sizeSlider->slider()->value();
            m_dockInter->resizeDock(offset, false);
        });
    });

    updateSliderValue();
    m_dconfigWatcher->bind("Control-Center_Dock_Size", m_sizeSlider);

    layout->addWidget(m_sizeSlider);

    // 多屏显示设置
    if (QDBusConnection::sessionBus().interface()->isServiceRegistered("com.deepin.dde.Dock")
            && QApplication::screens().size() > 1
            && !isCopyMode()
            && Utils::SettingValue("com.deepin.dde.dock.module.menu", QByteArray(), "multiscreenVisible", true).toBool()) {
        static QMap<QString, bool> g_screenSettingMap = {{tr("On screen where the cursor is"), false}
                                                         , {tr("Only on main screen"), true}};

        layout->addSpacing(10);
        layout->addWidget(m_screenSettingTitle);
        m_screenSettingComboxWidget->setTitle(tr("Show Dock"));
        m_screenSettingComboxWidget->addBackground();
        m_screenSettingComboxWidget->setComboxOption(QStringList() << tr("On screen where the cursor is") << tr("Only on main screen"));
        m_screenSettingComboxWidget->setCurrentText(g_screenSettingMap.key(m_dockInter->showInPrimary()));
        connect(m_screenSettingComboxWidget, &ComboxWidget::onSelectChanged, this, [ = ] (const QString &text) {
            m_dockInter->setShowInPrimary(g_screenSettingMap.value(text));
        });
        // 这里不会生效，但实际场景中也不存在有其他可配置的地方，可以不用处理
        connect(m_dockInter, &DBusInter::ShowInPrimaryChanged, this, [ = ] (bool showInPrimary) {
            if (m_screenSettingComboxWidget->comboBox()->currentText() == g_screenSettingMap.key(showInPrimary))
                return;

            m_screenSettingComboxWidget->blockSignals(true);
            m_screenSettingComboxWidget->setCurrentText(g_screenSettingMap.key(showInPrimary));
            m_screenSettingComboxWidget->blockSignals(false);
        });
        layout->addWidget(m_screenSettingComboxWidget);
        m_dconfigWatcher->bind("Control-Center_Dock_Multi-screen", m_screenSettingTitle);
        m_dconfigWatcher->bind("Control-Center_Dock_Multi-screen", m_screenSettingComboxWidget);
    } else {
        m_screenSettingTitle->setVisible(false);
        m_screenSettingComboxWidget->setVisible(false);
    }

    // 插件区域
    QDBusPendingReply<QStringList> reply = m_dockInter->GetLoadedPlugins();
    QStringList plugins = reply.value();
    if (reply.error().type() != QDBusError::ErrorType::NoError
            || !Utils::SettingValue("com.deepin.dde.dock.module.menu", QByteArray(), "hideVisible", true).toBool()) {
        m_pluginAreaTitle->setVisible(false);
        m_pluginTips->setVisible(false);
        m_pluginView->setVisible(false);
        qWarning() << "dbus call failed, method: 'GetLoadedPlugins()'";
    } else {
        const QMap<QString, QString> &pluginIconMap = {{"AiAssistant",      "dcc_dock_assistant"}
                                                       , {"show-desktop",   "dcc_dock_desktop"}
                                                       , {"onboard",        "dcc_dock_keyboard"}
                                                       , {"notifications",  "dcc_dock_notify"}
                                                       , {"shutdown",       "dcc_dock_power"}
                                                       , {"multitasking",   "dcc_dock_task"}
                                                       , {"datetime",       "dcc_dock_time"}
                                                       , {"trash",          "dcc_dock_trash"}};
        if (plugins.size() != 0) {
            layout->addSpacing(10);
            layout->addWidget(m_pluginAreaTitle);
            m_dconfigWatcher->bind("Control-Center_Dock_Plugins", m_pluginAreaTitle);

            DFontSizeManager::instance()->bind(m_pluginTips, DFontSizeManager::T8);
            m_pluginTips->adjustSize();
            m_pluginTips->setWordWrap(true);
            m_pluginTips->setContentsMargins(10, 5, 10, 5);
            m_pluginTips->setAlignment(Qt::AlignLeft);
            layout->addWidget(m_pluginTips);
            m_dconfigWatcher->bind("Control-Center_Dock_Plugins", m_pluginTips);

            m_pluginView->setAccessibleName("pluginList");
            m_pluginView->setBackgroundType(DStyledItemDelegate::BackgroundType::ClipCornerBackground);
            m_pluginView->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
            m_pluginView->setSelectionMode(QListView::SelectionMode::NoSelection);
            m_pluginView->setEditTriggers(DListView::NoEditTriggers);
            m_pluginView->setFrameShape(DListView::NoFrame);
            m_pluginView->setViewportMargins(0, 0, 0, 0);
            m_pluginView->setItemSpacing(1);

            QMargins itemMargins(m_pluginView->itemMargins());
            itemMargins.setLeft(14);
            m_pluginView->setItemMargins(itemMargins);

            m_pluginView->setVerticalScrollMode(QAbstractItemView::ScrollPerPixel);
            QScroller *scroller = QScroller::scroller(m_pluginView->viewport());
            QScrollerProperties sp;
            sp.setScrollMetric(QScrollerProperties::VerticalOvershootPolicy, QScrollerProperties::OvershootAlwaysOff);
            scroller->setScrollerProperties(sp);

            m_pluginView->setModel(m_pluginModel);

            layout->addWidget(m_pluginView);
            m_dconfigWatcher->bind("Control-Center_Dock_Plugins", m_pluginView);

            for (auto name : plugins) {
                DStandardItem *item = new DStandardItem(name);
                item->setFontSize(DFontSizeManager::T8);
                QSize size(16, 16);

                // 插件图标
                auto leftAction = new DViewItemAction(Qt::AlignVCenter, size, size, true);
                leftAction->setIcon(QIcon::fromTheme(pluginIconMap.value(m_dockInter->getPluginKey(name), "dcc_dock_plug_in")));
                item->setActionList(Qt::Edge::LeftEdge, {leftAction});

                auto rightAction = new DViewItemAction(Qt::AlignVCenter, size, size, true);
                bool visible = m_dockInter->getPluginVisible(name);
                auto checkstatus = visible ? DStyle::SP_IndicatorChecked : DStyle::SP_IndicatorUnchecked ;
                auto checkIcon = qobject_cast<DStyle *>(style())->standardIcon(checkstatus);
                rightAction->setIcon(checkIcon);
                item->setActionList(Qt::Edge::RightEdge, {rightAction});
                m_pluginModel->appendRow(item);

                connect(rightAction, &DViewItemAction::triggered, this, [ = ] {
                    bool checked = m_dockInter->getPluginVisible(name);
                    m_dockInter->setPluginVisible(name, !checked);
                    updateItemCheckStatus(name, !checked);
                });
            }
            // 固定大小,防止滚动
            int lineHeight = m_pluginView->visualRect(m_pluginView->indexAt(QPoint(0, 0))).height();
            m_pluginView->setMinimumHeight(lineHeight * plugins.size() + 10);
        } else {
            m_pluginAreaTitle->setVisible(false);
            m_pluginTips->setVisible(false);
            m_pluginView->setVisible(false);
        }
    }

    // 保持内容正常铺满
    layout->addStretch();

    // 界面内容过多时可滚动查看
    QWidget *widget = new QWidget;
    widget->setLayout(layout);
    setWidget(widget);
}

/**判断屏幕是否为复制模式的依据，第一个屏幕的X和Y值是否和其他的屏幕的X和Y值相等
 * 对于复制模式，这两个值肯定是相等的，如果不是复制模式，这两个值肯定不等，目前支持双屏
 * @brief DisplayManager::isCopyMode
 * @return
 */
bool ModuleWidget::isCopyMode()
{
    QList<QScreen *> screens = qApp->screens();
    if (screens.size() < 2)
        return false;

    // 在多个屏幕的情况下，如果所有屏幕的位置的X和Y值都相等，则认为是复制模式
    QRect screenRect = screens[0]->availableGeometry();
    for (int i = 1; i < screens.size(); i++) {
        QRect rect = screens[i]->availableGeometry();
        if (screenRect.x() != rect.x() || screenRect.y() != rect.y())
            return false;
    }

    return true;
}

void ModuleWidget::updateSliderValue()
{
    auto displayMode = m_daemonDockInter->displayMode();

    m_sizeSlider->slider()->blockSignals(true);
    if (displayMode == DisplayMode::Fashion) {
        if (int(m_daemonDockInter->windowSizeFashion()) != m_sizeSlider->slider()->value())
            m_sizeSlider->slider()->setValue(int(m_daemonDockInter->windowSizeFashion()));
    } else if (displayMode == DisplayMode::Efficient) {
        if (int(m_daemonDockInter->windowSizeEfficient()) != m_sizeSlider->slider()->value())
            m_sizeSlider->slider()->setValue(int(m_daemonDockInter->windowSizeEfficient()));
    }
    m_sizeSlider->slider()->blockSignals(false);
}

void ModuleWidget::updateItemCheckStatus(const QString &name, bool visible)
{
    for (int i = 0; i < m_pluginModel->rowCount(); ++i) {
        auto item = static_cast<DStandardItem *>(m_pluginModel->item(i));
        if (item->text() != name || item->actionList(Qt::Edge::RightEdge).size() < 1)
            continue;

        auto action = item->actionList(Qt::Edge::RightEdge).first();
        auto checkstatus = visible ? DStyle::SP_IndicatorChecked : DStyle::SP_IndicatorUnchecked ;
        auto icon = qobject_cast<DStyle *>(style())->standardIcon(checkstatus);
        action->setIcon(icon);
        m_pluginView->update(item->index());
        break;
    }
}
