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
    , m_itemManager(DockItemManager::instance(this))
    , m_dockInter(dockInter)
    , m_settingsMenu(new QMenu)
    , m_fashionModeAct(new QAction(tr("Fashion Mode"), this))
    , m_efficientModeAct(new QAction(tr("Efficient Mode"), this))
    , m_topPosAct(new QAction(tr("Top"), this))
    , m_bottomPosAct(new QAction(tr("Bottom"), this))
    , m_leftPosAct(new QAction(tr("Left"), this))
    , m_rightPosAct(new QAction(tr("Right"), this))
    , m_keepShownAct(new QAction(tr("Keep Shown"), this))
    , m_keepHiddenAct(new QAction(tr("Keep Hidden"), this))
    , m_smartHideAct(new QAction(tr("Smart Hide"), this))
    , m_menuEnable(true)
    , m_autoHide(true)
    , m_trashPluginShow(true)
    , m_opacity(0.4)
{
    initMember();
    initUI();
    initConnection();

    QTimer::singleShot(0, this, [=] {onGSettingsChanged("enable");});
}

MenuWorker::~MenuWorker()
{
    delete m_settingsMenu;

    QList<QString> keys = m_settingsModuleMap.keys();
    foreach (QString key, keys) {
        QGSettings *settings = m_settingsModuleMap[key];
        m_settingsModuleMap.remove(key);

        delete settings;
        settings = nullptr;
    }
}

void MenuWorker::initMember()
{
    m_settingsMenu->setAccessibleName("settingsmenu");

    m_fashionModeAct->setCheckable(true);
    m_efficientModeAct->setCheckable(true);
    m_topPosAct->setCheckable(true);
    m_bottomPosAct->setCheckable(true);
    m_leftPosAct->setCheckable(true);
    m_rightPosAct->setCheckable(true);
    m_keepShownAct->setCheckable(true);
    m_keepHiddenAct->setCheckable(true);
    m_smartHideAct->setCheckable(true);
}

void MenuWorker::initUI()
{
    QMenu *modeSubMenu = new QMenu(m_settingsMenu);
    modeSubMenu->setAccessibleName("modesubmenu");
    modeSubMenu->addAction(m_fashionModeAct);
    modeSubMenu->addAction(m_efficientModeAct);
    m_modeSubMenuAct = new QAction(tr("Mode"), this);
    m_modeSubMenuAct->setMenu(modeSubMenu);

    QMenu *locationSubMenu = new QMenu(m_settingsMenu);
    locationSubMenu->setAccessibleName("locationsubmenu");
    locationSubMenu->addAction(m_topPosAct);
    locationSubMenu->addAction(m_bottomPosAct);
    locationSubMenu->addAction(m_leftPosAct);
    locationSubMenu->addAction(m_rightPosAct);
    m_locationSubMenuAct = new QAction(tr("Location"), this);
    m_locationSubMenuAct->setMenu(locationSubMenu);

    QMenu *statusSubMenu = new QMenu(m_settingsMenu);
    statusSubMenu->setAccessibleName("statussubmenu");
    statusSubMenu->addAction(m_keepShownAct);
    statusSubMenu->addAction(m_keepHiddenAct);
    statusSubMenu->addAction(m_smartHideAct);
    m_statusSubMenuAct = new QAction(tr("Status"), this);
    m_statusSubMenuAct->setMenu(statusSubMenu);

    m_hideSubMenu = new QMenu(m_settingsMenu);
    m_hideSubMenu->setAccessibleName("pluginsmenu");
    m_hideSubMenuAct = new QAction(tr("Plugins"), this);
    m_hideSubMenuAct->setMenu(m_hideSubMenu);

    setSettingsMenu();
    m_settingsMenu->setTitle("Settings Menu");
}

void MenuWorker::initConnection()
{
    connect(m_settingsMenu, &QMenu::triggered, this, &MenuWorker::menuActionClicked);

    connect(settingsModule("menu"), &QGSettings::changed, this, &MenuWorker::onGSettingsChanged);
    connect(settingsModule("trash"), &QGSettings::changed, this, &MenuWorker::onTrashGSettingsChanged);

    connect(m_itemManager, &DockItemManager::trayVisableCountChanged, this, &MenuWorker::trayVisableCountChanged, Qt::QueuedConnection);

    DApplication *app = qobject_cast<DApplication *>(qApp);
    if (app) {
        connect(app, &DApplication::iconThemeChanged, this, &MenuWorker::gtkIconThemeChanged);
    }
}

const QGSettings *MenuWorker::settingsModule(const QString &module)
{
    if (!m_settingsModuleMap.contains(module))
        m_settingsModuleMap.insert(module, new QGSettings(QString("com.deepin.dde.dock.module." + module).toUtf8()));

    return m_settingsModuleMap[module];
}

void MenuWorker::setSettingsMenu()
{
    for (auto act : m_settingsMenu->actions())
        m_settingsMenu->removeAction(act);

    if (settingsModule("menu")->get("modeVisible").toBool())
        m_settingsMenu->addAction(m_modeSubMenuAct);

    if (settingsModule("menu")->get("locationVisible").toBool())
        m_settingsMenu->addAction(m_locationSubMenuAct);

    if (settingsModule("menu")->get("statusVisible").toBool())
        m_settingsMenu->addAction(m_statusSubMenuAct);

    if (settingsModule("menu")->get("hideVisible").toBool())
        m_settingsMenu->addAction(m_hideSubMenuAct);
}

void MenuWorker::showDockSettingsMenu()
{
    QTimer::singleShot(0, this, [=] {
        onGSettingsChanged("enable");
    });

    setAutoHide(false);
    setSettingsMenu();

    bool hasComposite = DWindowManagerHelper::instance()->hasComposite();

    // create actions
    QList<QAction *> actions;
    for (auto *p : m_itemManager->pluginList()) {
        if (!p->pluginIsAllowDisable())
            continue;

        const bool enable = !p->pluginIsDisable();
        const QString &name = p->pluginName();
        const QString &display = p->pluginDisplayName();

        if (!m_trashPluginShow && name == "trash") {
            continue;
        }

        if (name == "multitasking" && !hasComposite) {
            continue;
        }

        if (name == "deepin-screen-recorder-plugin") {
            continue;
        }
        QAction *act = new QAction(display, this);
        act->setCheckable(true);
        act->setChecked(enable);
        act->setData(name);

        if (!settingsModule(name)->keys().contains("visible") || settingsModule(name)->get("visible").toBool())
            actions << act;
    }

    // sort by name
    std::sort(actions.begin(), actions.end(), [](QAction * a, QAction * b) -> bool {
        return a->data().toString() > b->data().toString();
    });

    // add actions
    qDeleteAll(m_hideSubMenu->actions());
    for (auto act : actions)
        m_hideSubMenu->addAction(act);

    const DisplayMode displayMode = static_cast<DisplayMode>(m_dockInter->displayMode());
    const Position position = static_cast<Position>(m_dockInter->position());
    const HideMode hideMode = static_cast<HideMode>(m_dockInter->hideMode());

    m_fashionModeAct->setChecked(displayMode == Fashion);
    m_efficientModeAct->setChecked(displayMode == Efficient);
    m_topPosAct->setChecked(position == Top);
    m_bottomPosAct->setChecked(position == Bottom);
    m_leftPosAct->setChecked(position == Left);
    m_rightPosAct->setChecked(position == Right);
    m_keepShownAct->setChecked(hideMode == KeepShowing);
    m_keepHiddenAct->setChecked(hideMode == KeepHidden);
    m_smartHideAct->setChecked(hideMode == SmartHide);

    m_settingsMenu->exec(QCursor::pos());

    setAutoHide(true);
}

void MenuWorker::onGSettingsChanged(const QString &key)
{
    if (key != "enable") {
        return;
    }

    const QGSettings *setting = settingsModule("menu");
    m_menuEnable = setting->keys().contains("enable") ? setting->get("enable").toBool() : m_menuEnable;
}

void MenuWorker::onTrashGSettingsChanged(const QString &key)
{
    if (key != "enable") {
        return ;
    }

    const QGSettings *setting = settingsModule("trash");
    m_trashPluginShow = setting->keys().contains("enable") ? setting->get("enable").toBool() : m_trashPluginShow;
}

void MenuWorker::menuActionClicked(QAction *action)
{
    Q_ASSERT(action);

    if (action == m_fashionModeAct)
        return m_dockInter->setDisplayMode(DisplayMode::Fashion);
    if (action == m_efficientModeAct)
        return m_dockInter->setDisplayMode(Efficient);

    if (action == m_topPosAct)
        return m_dockInter->setPosition(Top);
    if (action == m_bottomPosAct)
        return m_dockInter->setPosition(Bottom);
    if (action == m_leftPosAct)
        return m_dockInter->setPosition(Left);
    if (action == m_rightPosAct)
        return m_dockInter->setPosition(Right);

    if (action == m_keepShownAct)
        return m_dockInter->setHideMode(KeepShowing);
    if (action == m_keepHiddenAct)
        return m_dockInter->setHideMode(KeepHidden);
    if (action == m_smartHideAct)
        return m_dockInter->setHideMode(SmartHide);

    // check plugin hide menu.
    const QString &data = action->data().toString();
    if (data.isEmpty())
        return;
    for (auto *p : m_itemManager->pluginList()) {
        if (p->pluginName() == data)
            return p->pluginStateSwitched();
    }
}

void MenuWorker::trayVisableCountChanged(const int &count)
{
    Q_UNUSED(count);

    emit trayCountChanged();
}

void MenuWorker::gtkIconThemeChanged()
{
    m_itemManager->refershItemsIcon();
}

void MenuWorker::setAutoHide(const bool autoHide)
{
    if (m_autoHide == autoHide)
        return;

    m_autoHide = autoHide;
    emit autoHideChanged(m_autoHide);
}
