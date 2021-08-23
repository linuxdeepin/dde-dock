/*
 * Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
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

#include "displaymanager.h"
#include "utils.h"

#include <QScreen>
#include <QApplication>
#include <QDBusConnection>
#include <QDesktopWidget>

DisplayManager::DisplayManager(QObject *parent)
    : QObject(parent)
    , m_gsettings(Utils::SettingsPtr("com.deepin.dde.dock.mainwindow", "/com/deepin/dde/dock/mainwindow/", this))
    , m_onlyInPrimary(Utils::SettingValue("com.deepin.dde.dock.mainwindow", "/com/deepin/dde/dock/mainwindow/", "onlyShowPrimary", false).toBool())
{
    connect(qApp, &QApplication::primaryScreenChanged, this, &DisplayManager::primaryScreenChanged);
    connect(qApp, &QGuiApplication::screenAdded, this, &DisplayManager::screenCountChanged);
    connect(qApp, &QGuiApplication::screenRemoved, this, &DisplayManager::screenCountChanged);

    if (m_gsettings)
        connect(m_gsettings, &QGSettings::changed, this, &DisplayManager::onGSettingsChanged);

    screenCountChanged();

    QTimer::singleShot(0, this, &DisplayManager::screenInfoChanged);
}

/**
 * @brief DisplayManager::screens
 * @return 返回当前可用的QScreen指针列表
 */
QList<QScreen *> DisplayManager::screens() const
{
    return m_screens;
}

/**
 * @brief DisplayManager::screen
 * @param screenName
 * @return 根据screenName参数找到对应的QScreen指针返回，否则返回nullptr
 */
QScreen *DisplayManager::screen(const QString &screenName) const
{
    for (auto s : m_screens) {
        if (s->name() == screenName)
            return s;
    }

    return nullptr;
}

/**
 * @brief DisplayManager::primary
 * @return 主屏幕名称
 */
QString DisplayManager::primary() const
{
    return qApp->primaryScreen() ? qApp->primaryScreen()->name() : QString();
}

/**
 * @brief Display::screenWidth
 * @return 所有屏幕逻辑宽度之和
 */
int DisplayManager::screenRawWidth() const
{
    int width = 0;
    for (auto s : m_screens) {
        width = qMax(width, s->geometry().x() + int(s->geometry().width() * s->devicePixelRatio()));
    }

    return width;
}
/**
 * @brief Display::screenHeight
 * @return 所有屏幕逻辑高度之和
 */

int DisplayManager::screenRawHeight() const
{
    int height = 0;
    for (auto s : m_screens) {
        height = qMax(height, s->geometry().y() + int(s->geometry().height() * s->devicePixelRatio()));
    }

    return height;
}

/**
 * @brief DisplayManager::canDock
 * @param s QScreen指针
 * @param pos 任务栏位置（上下左右）
 * @return 判断当前s屏幕上pos位置任务栏位置是否允许停靠
 */
bool DisplayManager::canDock(QScreen *s, Position pos) const
{
    return s ? m_screenPositionMap[s].value(pos) : false;
}

/**判断屏幕是否为复制模式的依据，第一个屏幕的X和Y值是否和其他的屏幕的X和Y值相等
 * 对于复制模式，这两个值肯定是相等的，如果不是复制模式，这两个值肯定不等，目前支持双屏
 * @brief DisplayManager::isCopyMode
 * @return
 */
bool DisplayManager::isCopyMode()
{
    QList<QScreen *> screens = this->screens();
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

/**
 * @brief DisplayManager::updateScreenDockInfo
 * 更新屏幕停靠信息
 */
void DisplayManager::updateScreenDockInfo()
{
    // TODO 目前仅仅支持双屏，如果超过双屏，会出现异常，这里可以考虑做成通用的处理规则

    // 先清除原先的数据，然后再更新
    m_screenPositionMap.clear();

    if (m_screens.isEmpty())
        return;

    // reset map
    for (auto s : m_screens) {
        QMap <Position, bool> map;
        map.insert(Position::Top, true);
        map.insert(Position::Bottom, true);
        map.insert(Position::Left, true);
        map.insert(Position::Right, true);
        m_screenPositionMap.insert(s, map);
    }

    // 仅显示在主屏时的处理
    if (m_onlyInPrimary) {
        for (auto s : m_screens) {
            if (s != qApp->primaryScreen()) {
                QMap <Position, bool> map;
                map.insert(Position::Top, false);
                map.insert(Position::Bottom, false);
                map.insert(Position::Left, false);
                map.insert(Position::Right, false);
                m_screenPositionMap.insert(s, map);
            }
        }
        return;
    }

    if (m_screens.size() == 1) {
        return;
    }

    // 最多支持双屏,这里只计算双屏,单屏默认四边均可停靠任务栏
    if (m_screens.size() == 2) {
        QRect s0 = m_screens.at(0)->geometry();
        s0.setSize(s0.size() * m_screens.at(0)->devicePixelRatio());
        QRect s1 = m_screens.at(1)->geometry();
        s1.setSize(s1.size() * m_screens.at(1)->devicePixelRatio());

        qInfo() << "monitor info changed" << m_screens.at(0)->name() << s0 << m_screens.at(1)->name() << s1;

        int s0top = s0.y();
        int s0bottom = s0.y() + s0.height();
        int s0left = s0.x();
        int s0Right = s0.x() + s0.width();

        int s1top = s1.y();
        int s1bottom = s1.y() + s1.height();
        int s1left = s1.x();
        int s1Right = s1.x() + s1.width();

        QPoint s0topLeft = QPoint(s0.x(), s0.y());
        QPoint s0topRight = QPoint(s0.x()+ s0.width(), s0.y());
        QPoint s0bottomRight = QPoint(s0.x()+ s0.width(), s0.y() + s0.height());
        QPoint s0bottomLeft = QPoint(s0.x(), s0.y() + s0.height());

        QPoint s1topLeft = QPoint(s1.x(), s1.y());
        QPoint s1topRight = QPoint(s1.x() + s1.width(), s1.y());
        QPoint s1bottomRight = QPoint(s1.x() + s1.width(), s1.y() + s1.height());
        QPoint s1bottomLeft = QPoint(s1.x(), s1.y() + s1.height());

        /*
         * 对角拼接，重置，默认均可停靠
---------                       ---------
|       |                       |       |
|   s0  |                       |   s1  |
|       |                       |       |
-----------------               -----------------
        |       |                       |       |
        |   s1  |                           s0  |
        |       |                       |       |
        ---------                       ---------
*/
        if (s0bottomRight == s1topLeft
                || s0topLeft == s1bottomRight) {
            return;
        }

        /*
         * 左右拼接，s0左，s1右
---------------------               -------------
|           |       |               |           |
|           |  s1   |               |           |--------
|     s0    |       |               |     s0    |       |
|           |--------               |           |  s1   |
|           |                       |           |       |
-------------                       ---------------------
*/
        if (s0Right == s1left
                && (s0topRight == s1topLeft || s0bottomRight == s1bottomLeft)) {
            m_screenPositionMap[m_screens.at(0)].insert(Position::Right, false);
            m_screenPositionMap[m_screens.at(1)].insert(Position::Left, false);
            return;
        }

        /*
         * 左右拼接，s1左，s0右
---------------------               -------------
|           |       |               |           |
|           |  s0   |               |           |--------
|     s1    |       |               |     s1    |       |
|           |--------               |           |  s0   |
|           |                       |           |       |
-------------                       ---------------------
*/
        if (s0left== s1Right
                && (s0topLeft == s1topRight || s0bottomLeft == s1bottomRight)) {
            m_screenPositionMap[m_screens.at(0)].insert(Position::Left, false);
            m_screenPositionMap[m_screens.at(1)].insert(Position::Right, false);
            return;
        }

        /*
         * 上下拼接，s0上，s1下
---------                           ---------
|       |                           |       |
|   s0  |                           |   s0  |
|       |                           |       |
-------------                   -------------
|           |                   |           |
|           |                   |           |
|     s1    |                   |     s1    |
|           |                   |           |
|           |                   |           |
-------------                   -------------
*/
        if (s0bottom == s1top
                && (s0bottomLeft == s1topLeft || s0bottomRight == s1topRight)) {
            m_screenPositionMap[m_screens.at(0)].insert(Position::Bottom, false);
            m_screenPositionMap[m_screens.at(1)].insert(Position::Top, false);
            return;
        }

        /*
         * 上下拼接，s1上，s0下
---------                   ---------
|       |                   |       |
|   s1  |                   |   s1  |
|       |                   |       |
-------------           -------------
|           |           |           |
|           |           |           |
|     s0    |           |     s0    |
|           |           |           |
|           |           |           |
-------------           -------------
*/
        if (s0top == s1bottom
                && (s0topLeft == s1bottomLeft || s0topRight == s1bottomRight)) {
            m_screenPositionMap[m_screens.at(0)].insert(Position::Top, false);
            m_screenPositionMap[m_screens.at(1)].insert(Position::Bottom, false);
            return;
        }
        return;
    }

    // 目前不支持链接超过两个以上的显示器
    Q_UNREACHABLE();
}

/**
 * @brief DisplayManager::screenCountChanged
 * 屏幕数量发生变化时，此函数应被调用，更新屏幕相关信息
 * @note 除初始化时需要手动调用一次外，其他时间会自动被调用
 */
void DisplayManager::screenCountChanged()
{
    // 找到过期的screen指针
    QList<QScreen *> to_remove_list;
    for (auto s : m_screens) {
        if (!qApp->screens().contains(s))
            to_remove_list.append(s);
    }

    // 找出新增的screen指针
    QList<QScreen *> to_add_list;
    for (auto s : qApp->screens()) {
        if (!m_screens.contains(s)) {
            to_add_list.append(s);
        }
    }

    // 取消关联
    for (auto s : to_remove_list) {
        disconnect(s);
        m_screens.removeOne(s);
    }

    // 创建关联
    for (auto s : to_add_list) {
        s->setOrientationUpdateMask(Qt::PrimaryOrientation
                                    | Qt::LandscapeOrientation
                                    | Qt::PortraitOrientation
                                    | Qt::InvertedLandscapeOrientation
                                    | Qt::InvertedPortraitOrientation);

        // 显示器信息发生任何变化时，都应该重新刷新一次任务栏的显示位置
        connect(s, &QScreen::geometryChanged, this, &DisplayManager::dockInfoChanged);
        connect(s, &QScreen::availableGeometryChanged, this, &DisplayManager::dockInfoChanged);
        connect(s, &QScreen::physicalSizeChanged, this, &DisplayManager::dockInfoChanged);
        connect(s, &QScreen::physicalDotsPerInchChanged, this, &DisplayManager::dockInfoChanged);
        connect(s, &QScreen::logicalDotsPerInchChanged, this, &DisplayManager::dockInfoChanged);
        connect(s, &QScreen::virtualGeometryChanged, this, &DisplayManager::dockInfoChanged);
        connect(s, &QScreen::primaryOrientationChanged, this, &DisplayManager::dockInfoChanged);
        connect(s, &QScreen::orientationChanged, this, &DisplayManager::dockInfoChanged);
        connect(s, &QScreen::refreshRateChanged, this, &DisplayManager::dockInfoChanged);

        m_screens.append(s);
    }

    // 屏幕数量发生变化，应该刷新一下任务栏的显示
    dockInfoChanged();
}

void DisplayManager::dockInfoChanged()
{
    updateScreenDockInfo();

#ifdef QT_DEBUG
    qInfo() << m_screenPositionMap;
#endif

    Q_EMIT screenInfoChanged();
}

/**
 * @brief DisplayManager::onGSettingsChanged
 * @param key
 * 监听onlyShowPrimary配置的变化，此时有变化时应该刷新一下任务栏的显示信息
 */
void DisplayManager::onGSettingsChanged(const QString &key)
{
    if (key == "onlyShowPrimary") {
        m_onlyInPrimary = Utils::SettingValue("com.deepin.dde.dock.mainwindow", "/com/deepin/dde/dock/mainwindow/", "onlyShowPrimary", false).toBool();

        dockInfoChanged();
    }
}
