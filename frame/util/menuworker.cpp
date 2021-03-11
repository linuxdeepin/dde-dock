/*
 * Copyright (C) 2018 ~ 2028 Deepin Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng_cm@deepin.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng_cm@deepin.com>
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
#include "menuworker.h"
#include "dockitemmanager.h"

#include <QAction>
#include <QMenu>
#include <QGSettings>

#include <DApplication>

MenuWorker::MenuWorker(DBusDock *dockInter,QWidget *parent)
    : QObject (parent)
    , m_dockInter(dockInter)
    , m_autoHide(true)
{
    DApplication *app = qobject_cast<DApplication *>(qApp);
    if (app) {
        connect(app, &DApplication::iconThemeChanged, this, &MenuWorker::gtkIconThemeChanged);
    }
}

MenuWorker::~MenuWorker()
{
}

const QGSettings *MenuWorker::SettingsPtr(const QString &module)
{
    return QGSettings::isSchemaInstalled(QString("com.deepin.dde.dock.module." + module).toUtf8())
            ? new QGSettings(QString("com.deepin.dde.dock.module." + module).toUtf8(), QByteArray(), this) // 自动销毁
            : nullptr;
}

void MenuWorker::showDockSettingsMenu()
{
    // 菜单功能被禁用
    const QGSettings *setting = SettingsPtr("menu");
    if (setting && setting->keys().contains("enable") && !setting->get("enable").toBool())
        return;

    // 菜单将要被打开
    setAutoHide(false);

    QMenu settingsMenu;
    settingsMenu.setAccessibleName("settingsmenu");

    // 模式
    if (SettingsPtr("menu") && SettingsPtr("menu")->get("modeVisible").toBool()) {
        const DisplayMode displayMode = static_cast<DisplayMode>(m_dockInter->displayMode());

        QMenu *modeSubMenu = new QMenu(&settingsMenu);
        modeSubMenu->setAccessibleName("modesubmenu");

        QAction *fashionModeAct = new QAction(tr("Fashion Mode"), this);
        QAction *efficientModeAct = new QAction(tr("Efficient Mode"), this);

        fashionModeAct->setCheckable(true);
        efficientModeAct->setCheckable(true);

        fashionModeAct->setChecked(displayMode == Fashion);
        efficientModeAct->setChecked(displayMode == Efficient);

        connect(fashionModeAct, &QAction::triggered, this, [ = ]{ m_dockInter->setDisplayMode(DisplayMode::Fashion); });
        connect(efficientModeAct, &QAction::triggered, this, [ = ]{ m_dockInter->setDisplayMode(DisplayMode::Efficient); });

        modeSubMenu->addAction(fashionModeAct);
        modeSubMenu->addAction(efficientModeAct);

        QAction *act = new QAction(tr("Mode"), this);
        act->setMenu(modeSubMenu);

        settingsMenu.addAction(act);
    }

    // 位置
    if (SettingsPtr("menu") && SettingsPtr("menu")->get("locationVisible").toBool()) {
        const Position position = static_cast<Position>(m_dockInter->position());

        QMenu *locationSubMenu = new QMenu(&settingsMenu);
        locationSubMenu->setAccessibleName("locationsubmenu");

        QAction *topPosAct = new QAction(tr("Top"), this);
        QAction *bottomPosAct = new QAction(tr("Bottom"), this);
        QAction *leftPosAct = new QAction(tr("Left"), this);
        QAction *rightPosAct = new QAction(tr("Right"), this);

        topPosAct->setCheckable(true);
        bottomPosAct->setCheckable(true);
        leftPosAct->setCheckable(true);
        rightPosAct->setCheckable(true);

        topPosAct->setChecked(position == Top);
        bottomPosAct->setChecked(position == Bottom);
        leftPosAct->setChecked(position == Left);
        rightPosAct->setChecked(position == Right);

        connect(topPosAct, &QAction::triggered, this, [ = ]{ m_dockInter->setPosition(Top); });
        connect(bottomPosAct, &QAction::triggered, this, [ = ]{ m_dockInter->setPosition(Bottom); });
        connect(leftPosAct, &QAction::triggered, this, [ = ]{ m_dockInter->setPosition(Left); });
        connect(rightPosAct, &QAction::triggered, this, [ = ]{ m_dockInter->setPosition(Right); });

        locationSubMenu->addAction(topPosAct);
        locationSubMenu->addAction(bottomPosAct);
        locationSubMenu->addAction(leftPosAct);
        locationSubMenu->addAction(rightPosAct);

        QAction *act = new QAction(tr("Location"), this);
        act->setMenu(locationSubMenu);

        settingsMenu.addAction(act);
    }

    // 状态
    if (SettingsPtr("menu") && SettingsPtr("menu")->get("statusVisible").toBool()) {
        const HideMode hideMode = static_cast<HideMode>(m_dockInter->hideMode());

        QMenu *statusSubMenu = new QMenu(&settingsMenu);
        statusSubMenu->setAccessibleName("statussubmenu");

        QAction *keepShownAct = new QAction(tr("Keep Shown"), this);
        QAction *keepHiddenAct = new QAction(tr("Keep Hidden"), this);
        QAction *smartHideAct = new QAction(tr("Smart Hide"), this);

        keepShownAct->setCheckable(true);
        keepHiddenAct->setCheckable(true);
        smartHideAct->setCheckable(true);

        keepShownAct->setChecked(hideMode == KeepShowing);
        keepHiddenAct->setChecked(hideMode == KeepHidden);
        smartHideAct->setChecked(hideMode == SmartHide);

        connect(keepShownAct, &QAction::triggered, this, [ = ]{ m_dockInter->setHideMode(KeepShowing); });
        connect(keepHiddenAct, &QAction::triggered, this, [ = ]{ m_dockInter->setHideMode(KeepHidden); });
        connect(smartHideAct, &QAction::triggered, this, [ = ]{ m_dockInter->setHideMode(SmartHide); });

        statusSubMenu->addAction(keepShownAct);
        statusSubMenu->addAction(keepHiddenAct);
        statusSubMenu->addAction(smartHideAct);

        QAction *act = new QAction(tr("Status"), this);
        act->setMenu(statusSubMenu);

        settingsMenu.addAction(act);
    }

    // 插件
    if (SettingsPtr("menu") && SettingsPtr("menu")->get("hideVisible").toBool()) {
        QMenu *hideSubMenu = new QMenu(&settingsMenu);
        hideSubMenu->setAccessibleName("pluginsmenu");

        QAction *hideSubMenuAct = new QAction(tr("Plugins"), this);
        hideSubMenuAct->setMenu(hideSubMenu);

        // create actions
        QList<QAction *> actions;
        for (auto *p : DockItemManager::instance()->pluginList()) {
            if (!p->pluginIsAllowDisable())
                continue;

            const bool enable = !p->pluginIsDisable();
            const QString &name = p->pluginName();
            const QString &display = p->pluginDisplayName();

            // 模块和菜单均需要响应enable配置的变化
            const QGSettings *setting = SettingsPtr(name);
            if (setting && setting->keys().contains("enable") && !setting->get("enable").toBool()) {
                continue;
            }

            // 未开启窗口特效时，同样不显示多任务视图插件
            if (name == "multitasking" && !DWindowManagerHelper::instance()->hasComposite()) {
                continue;
            }

            // TODO 记得让录屏那边加一个enable的配置项，默认值设置成false,就不用针对这个插件特殊处理了
            if (name == "deepin-screen-recorder-plugin") {
                continue;
            }

            QAction *act = new QAction(display, this);
            act->setCheckable(true);
            act->setChecked(enable);
            act->setData(name);

            connect(act, &QAction::triggered, this, [ p ]{p->pluginStateSwitched();});

            // check plugin hide menu.
            if (SettingsPtr(name) && (!SettingsPtr(name)->keys().contains("visible") || SettingsPtr(name)->get("visible").toBool()))
                actions << act;
        }

        // sort by name
        std::sort(actions.begin(), actions.end(), [](QAction * a, QAction * b) -> bool {
            return a->data().toString() > b->data().toString();
        });

        // add plugins actions
        qDeleteAll(hideSubMenu->actions());
        for (auto act : actions)
            hideSubMenu->addAction(act);

        // add plugins menu
        settingsMenu.addAction(hideSubMenuAct);
    }

    settingsMenu.setTitle("Settings Menu");
    settingsMenu.exec(QCursor::pos());

    // 菜单已经关闭
    setAutoHide(true);
}

void MenuWorker::gtkIconThemeChanged()
{
    DockItemManager::instance()->refershItemsIcon();
}

void MenuWorker::setAutoHide(const bool autoHide)
{
    if (m_autoHide == autoHide)
        return;

    m_autoHide = autoHide;
    emit autoHideChanged(m_autoHide);
}
