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
    connect(qApp, &QApplication::primaryScreenChanged, this, &DisplayManager::dockInfoChanged);
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

    // 适配多个屏幕的情况
    for(auto s : m_screens) {
        QList<QScreen *> otherScreens = m_screens;
        otherScreens.removeAll(s);
        for (auto other : otherScreens) {
            QRect ourRect = QRect(s->geometry().topLeft(), s->geometry().size() * s->devicePixelRatio());
            int ourBottom = ourRect.top() + ourRect.height();
            int ourTop = ourRect.top();
            int ourLeft = ourRect.left();
            int ourRight = ourRect.left() + ourRect.width();
            QPoint ourLeftBottom = QPoint(ourLeft, ourBottom);
            QPoint ourRightBottom = QPoint(ourRight, ourBottom);
            QPoint ourRightTop = QPoint(ourRight, ourTop);

            QRect otherRect = QRect(other->geometry().topLeft(), other->geometry().size() * other->devicePixelRatio());
            int otherBottom = otherRect.top() + otherRect.height();
            int otherTop = otherRect.top();
            int otherLeft = otherRect.left();
            int otherRight = otherRect.left() + otherRect.width();
            QPoint otherLeftBottom = QPoint(otherLeft, otherBottom);
            QPoint otherLeftTop = QPoint(otherLeft, otherTop);
            QPoint otherRightTop = QPoint(otherRight, otherTop);

            /*
                     * 上下拼接，our屏幕左右移动。
                     * our屏幕从other屏幕对角的左上侧向右移动，至other屏幕对角的右上侧位置
              ---------                                   ---------
              |       |                                   |       |
              |  our  |           ======>>>>>             |  our  |
              |       |                                   |       |
              ---------------                       ---------------
                    |       |                       |       |
                    | other |                       | other |
                    |       |                       |       |
                    ---------                       ---------
            */
            // 上下拼接
            if (ourBottom == otherTop
                    && (ourRight >= otherLeft)
                    && (ourLeft <= otherRight)) {
                // 排除对角排列
                if (ourLeftBottom == otherRightTop
                        || ourRightBottom == otherLeftTop)
                    continue;
                m_screenPositionMap[s][Position::Bottom] = false;
                m_screenPositionMap[other][Position::Top] = false;
            }

            /*
                     * 左右拼接，our屏幕上下移动。
                     * our屏幕从other屏幕对角的左上侧向下移动，至other屏幕对角的最左下侧位置
              ---------                                          ---------
              |       |                                          |       |
              |  our  |              ======>>>>>        ---------| other |
              |       |--------                         |        |       |
              --------|       |                         |  our   |--------
                      | other |                         |        |
                      |       |                         ----------
                      ---------
            */
            // 左右拼接
            if (otherLeft == ourRight
                    && (ourTop <= otherBottom)
                    && (ourBottom >= otherTop)) {
                // 排除对角排列
                if (ourRightTop == otherLeftBottom
                        || ourRightBottom == otherLeftTop)
                    continue;
                m_screenPositionMap[s][Position::Right] = false;
                m_screenPositionMap[other][Position::Left] = false;
            }
        }
    }
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
