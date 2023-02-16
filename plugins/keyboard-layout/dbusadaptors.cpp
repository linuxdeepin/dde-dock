// Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "utils.h"
#include "dbusadaptors.h"

#include <DDBusSender>
#include <QDebug>

DBusAdaptors::DBusAdaptors(QObject *parent)
    : QDBusAbstractAdaptor(parent)
    , m_keyboard(new Keyboard("org.deepin.dde.InputDevices1",
                              "/org/deepin/dde/InputDevice1/Keyboard",
                              QDBusConnection::sessionBus(), this))
    , m_menu(new QMenu)
    , m_gsettings(Utils::ModuleSettingsPtr("keyboard", QByteArray(), this))
{
    m_keyboard->setSync(false);

    connect(m_keyboard, &Keyboard::CurrentLayoutChanged, this, &DBusAdaptors::onCurrentLayoutChanged);
    connect(m_keyboard, &Keyboard::UserLayoutListChanged, this, &DBusAdaptors::onUserLayoutListChanged);

    connect(m_menu, &QMenu::triggered, this, &DBusAdaptors::handleActionTriggered);

    // init data
    initAllLayoutList();
    onCurrentLayoutChanged(m_keyboard->currentLayout());
    onUserLayoutListChanged(m_keyboard->userLayoutList());

    if (m_gsettings)
        connect(m_gsettings, &QGSettings::changed, this, &DBusAdaptors::onGSettingsChanged);
}

DBusAdaptors::~DBusAdaptors()
{
    delete m_menu;
}

QString DBusAdaptors::layout() const
{
    if (m_gsettings && m_gsettings->keys().contains("enable") && !m_gsettings->get("enable").toBool())
        return QString();

    if (m_userLayoutList.size() < 2) {
        // do NOT show keyboard indicator
        return QString();
    }

    if (m_currentLayout.isEmpty()) {
        // refetch data
        QTimer::singleShot(1000, m_keyboard, &Keyboard::currentLayout);
        qDebug() << Q_FUNC_INFO << "currentLayout is Empty!!";
    }

    return m_currentLayout;
}

void DBusAdaptors::setLayout(const QString &str)
{
    m_currentLayout = str;
    emit layoutChanged(str);
}

Keyboard *DBusAdaptors::getCurrentKeyboard()
{
    return m_keyboard;
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
            qDebug() << "failed to get all keyboard list: " << call.error().message();
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
                .service("org.deepin.dde.ControlCenter1")
                .interface("org.deepin.dde.ControlCenter1")
                .path("/org/deepin/dde/ControlCenter1")
                .method("ShowPage")
                .arg(QString("keyboard/keyboardGeneral/keyboardLayout"))
                .call();
    }

    const QString layout = action->objectName();
    if (m_userLayoutList.contains(layout)) {
        m_keyboard->setCurrentLayout(layout);
    }
}

void DBusAdaptors::onGSettingsChanged(const QString &key)
{
    Q_UNUSED(key);

    // 键盘布局插件处显示的内容就是QLabel中的内容，有文字了就显示，没有文字就不显示了
    if (m_gsettings && m_gsettings->keys().contains("enable")) {
        const bool enable = m_gsettings->get("enable").toBool();
        QString layoutStr = getCurrentKeyboard()->currentLayout().split(';').first();
        setLayout(enable ? layoutStr : "");
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
