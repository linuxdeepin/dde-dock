/*
 * Copyright (C) 2020 ~ 2021 Deepin Technology Co., Ltd.
 *
 * Author:     songwentao <songwentao@uniontech.com>
 *
 * Maintainer: songwentao <songwentao@uniontech.com>
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
 *
 * This program aims to cache the the icon and name of apps to the hash table,
 * which can decrease the repeated resource consumption of loading the app info in the
 * running time.
 */

#include "menudialog.h"

#include <QApplication>
#include <QEvent>
#include <QMouseEvent>

Menu::Menu(QWidget *dockItem, QWidget *parent)
    : QMenu(parent)
    , m_dockItem(dockItem)
    , m_eventInter(new XEventMonitor("com.deepin.api.XEventMonitor", "/com/deepin/api/XEventMonitor", QDBusConnection::sessionBus(), this))
    , m_dockInter(new DBusDock("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this))
{
    setWindowFlags(Qt::WindowStaysOnTopHint | Qt::X11BypassWindowManagerHint | Qt::Dialog);
    setObjectName("rightMenu");
    qApp->installEventFilter(this);

    if (m_dockItem)
        m_dockItem->installEventFilter(this);

    // 按下任务栏以外区域释放鼠标时，关闭右键菜单，否则会导致点击菜单项后无响应
    connect(m_eventInter, &XEventMonitor::ButtonRelease, this, &Menu::onButtonPress);
}

void Menu::onButtonPress()
{
    if (!QRect(m_dockInter->frontendWindowRect()).contains(QCursor::pos()))
        this->hide();
}

/** 右键菜单显示后在很多场景下都需要隐藏，为避免在各个控件中分别做隐藏右键菜单窗口的处理，
 *  因此这里统一做了处理。
 * @brief Menu::eventFilter
 * @param watched 过滤器监听对象
 * @param event 过滤器事件对象
 * @return 返回true, 事件不再向下传递，返回false，事件向下传递
 */
bool Menu::eventFilter(QObject *watched, QEvent *event)
{
    // 存在rightMenu和rightMenuWindow的对象名
    if (!watched->objectName().startsWith("rightMenu")) {
        if (event->type() == QEvent::MouseButtonPress) {
            if (watched != m_dockItem && watched != this && this->isVisible()) {
                // 鼠标点击时，除当前菜单外，其他显示的菜单都要隐藏
                hide();
            }
        } else if (event->type() == QEvent::DragMove || event->type() == QEvent::Wheel || event->type() == QEvent::Move) {
            // 按下应用拖动，按下应用从菜单上方移动时，鼠标滚轮滚动时，隐藏右键菜单
            hide();
        }
    }

    return QMenu::eventFilter(watched, event);
}
