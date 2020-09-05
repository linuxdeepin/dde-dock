/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     rekols <rekols@foxmail.com>
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

#include "dbusadaptors.h"
#include <DDBusSender>
#include <QDebug>
#include <QWindow>
#include <QCursor>
#include <QGSettings>

#include <X11/Xcursor/Xcursor.h>


QCursor* DBusAdaptors::kbtLoadQCursorFromX11Cursor(const char* theme, const char* cursorName, int cursorSize)
{
    if (theme == nullptr || cursorName == nullptr || cursorSize <= 0)
        return nullptr;

   XcursorImages *images = XcursorLibraryLoadImages(cursorName, theme, cursorSize);
    if (images == nullptr || images->images[0] == nullptr) {
        qWarning() << "loadCursorFalied, theme =" << theme << ", cursorName=" << cursorName;
        return nullptr;
    }
    const int imgW = images->images[0]->width;
    const int imgH = images->images[0]->height;
    QImage img((const uchar*)images->images[0]->pixels, imgW, imgH, QImage::Format_ARGB32);
    QPixmap pixmap = QPixmap::fromImage(img);
    QCursor *cursor = new QCursor(pixmap, images->images[0]->xhot, images->images[0]->yhot);
    delete images;
    return cursor;
}

void DBusAdaptors::kbtLoadQCursorForUpdateMenu(QWidget *menu_win)
{
    static QCursor *lastArrowCursor = nullptr;
    static QString  lastCursorTheme;
    int lastCursorSize = 0;
    QGSettings gsetting("com.deepin.xsettings", "/com/deepin/xsettings/");
    QString theme = gsetting.get("gtk-cursor-theme-name").toString();
    int cursorSize = gsetting.get("gtk-cursor-theme-size").toInt();
    if (theme != lastCursorTheme || cursorSize != lastCursorSize)
    {
        qDebug() << QString("Menu Update Cursor (theme=%1,%2 ; size=%3,%4)...").arg(lastCursorTheme).arg(theme).arg(lastCursorSize).arg(cursorSize);
        QCursor *cursor = kbtLoadQCursorFromX11Cursor(theme.toStdString().c_str(), "left_ptr", cursorSize);
        lastCursorTheme = theme;
        lastCursorSize = cursorSize;
        if(menu_win)
            menu_win->setCursor(*cursor);
        if (lastArrowCursor != nullptr)
            delete lastArrowCursor;

        lastArrowCursor = cursor;
    }
}

DBusAdaptors::DBusAdaptors(QObject *parent)
    : QDBusAbstractAdaptor(parent),
      m_keyboard(new Keyboard("com.deepin.daemon.InputDevices",
                              "/com/deepin/daemon/InputDevice/Keyboard",
                              QDBusConnection::sessionBus(), this)),
      m_menu(new QMenu())
{
    m_keyboard->setSync(false);

    //qt自定义插件包增加menu属性,供QMenu使用显示属性。
    m_menu->setAttribute(Qt::WA_NativeWindow);
    m_menu->windowHandle()->setProperty("_d_dwayland_window-type" , "menu");
    //修复右键菜单光标变大问题，临时解决方案
    kbtLoadQCursorForUpdateMenu(m_menu);

    connect(m_keyboard, &Keyboard::CurrentLayoutChanged, this, &DBusAdaptors::onCurrentLayoutChanged);
    connect(m_keyboard, &Keyboard::UserLayoutListChanged, this, &DBusAdaptors::onUserLayoutListChanged);

    connect(m_menu, &QMenu::triggered, this, &DBusAdaptors::handleActionTriggered);

    // init data
    initAllLayoutList();
    onCurrentLayoutChanged(m_keyboard->currentLayout());
    onUserLayoutListChanged(m_keyboard->userLayoutList());
}

DBusAdaptors::~DBusAdaptors()
{
}

QString DBusAdaptors::layout() const
{
    if (m_userLayoutList.size() < 2) {
        // do NOT show keyboard indicator
        return QString();
    }

    if (m_currentLayout.isEmpty()) {
        // refetch data
        QTimer::singleShot(1000, m_keyboard, &Keyboard::currentLayout);
        qWarning() << Q_FUNC_INFO << "currentLayout is Empty!!";
    }

    return m_currentLayout;
}

void DBusAdaptors::onClicked(int button, int x, int y)
{
//    button value means(XCB_BUTTON_INDEX):
//    0, Any of the following (or none)
//    1, The left mouse button.
//    2, The right mouse button.
//    3, The middle mouse button.
//    4, Scroll wheel. TODO: direction?
//    5, Scroll wheel. TODO: direction?

    Q_UNUSED(button);

    if (m_menu && m_userLayoutList.size() >= 2) {
        m_menu->exec(QPoint(x, y));
    }
}

void DBusAdaptors::onCurrentLayoutChanged(const QString &value)
{
    m_currentLayoutRaw = value;
    m_currentLayout = value.split(';').first();

    refreshMenuSelection();

    emit layoutChanged(layout());
}

void DBusAdaptors::onUserLayoutListChanged(const QStringList &value)
{
    m_userLayoutList = value;

    initAllLayoutList();

    emit layoutChanged(layout());
}

void DBusAdaptors::initAllLayoutList()
{
    QDBusPendingCall call = m_keyboard->LayoutList();
    QDBusPendingCallWatcher *watcher = new QDBusPendingCallWatcher(call, this);
    connect(watcher, &QDBusPendingCallWatcher::finished, this, [=] {
        if (call.isError()) {
            qWarning() << "failed to get all keyboard list: " << call.error().message();
        } else {
            QDBusReply<KeyboardLayoutList> reply = call.reply();
            m_allLayoutList = reply.value();
            refreshMenu();
        }

        watcher->deleteLater();
    });
}

void DBusAdaptors::refreshMenu()
{
    if (!m_menu || m_userLayoutList.size() < 2) {
        return;
    }

    // all action object will be deleted
    m_menu->clear();

    for (const QString &layoutRawName : m_userLayoutList) {
        const QString layoutName = duplicateCheck(layoutRawName);
        const QString layoutLocalizedName = m_allLayoutList.value(layoutRawName);
        const QString text = QString("%1 (%2)").arg(layoutLocalizedName).arg(layoutName);

        QAction *action = new QAction(text, m_menu);
        action->setObjectName(layoutRawName);
        action->setCheckable(true);
        action->setChecked(layoutRawName == m_currentLayoutRaw);
        m_menu->addAction(action);
    }

    m_menu->addSeparator();

    // will be deleted after QMenu->clear() above
    m_addLayoutAction = new QAction(tr("Add keyboard layout"), m_menu);

    m_menu->addAction(m_addLayoutAction);
}

void DBusAdaptors::refreshMenuSelection()
{
    for (QAction *action : m_menu->actions()) {
        action->setChecked(action->objectName() == m_currentLayoutRaw);
    }
}

void DBusAdaptors::handleActionTriggered(QAction *action)
{
    if (action == m_addLayoutAction) {
        DDBusSender()
                .service("com.deepin.dde.ControlCenter")
                .interface("com.deepin.dde.ControlCenter")
                .path("/com/deepin/dde/ControlCenter")
                .method("ShowPage")
                .arg(QString("keyboard"))
                .arg(QString("Keyboard Layout/Add Keyboard Layout"))
                .call();
    }

    const QString layout = action->objectName();
    if (m_userLayoutList.contains(layout)) {
        m_keyboard->setCurrentLayout(layout);
    }
}

QString DBusAdaptors::duplicateCheck(const QString &kb)
{
    QStringList list;
    const QString kbFirst = kb.split(";").first();
    for (const QString &data : m_userLayoutList) {
        if (data.split(";").first() == kbFirst) {
            list << data;
        }
    }

    const QString kblayout = kb.split(";").first().mid(0, 2);

    return kblayout + (list.count() > 1 ? QString::number(list.indexOf(kb) + 1) : "");
}
