/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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

#include "dockinterface.h"

#ifdef USE_AM
// 因为 types/dockrect.h 文件中定义了DockRect类，而在此处也定义了DockRect，
// 所以在此处先加上DOCKRECT_H宏(types/dockrect.h文件中定义的宏)来禁止包含types/dockrect.h头文件
// 否则会出现重复定义的错误
#define DOCKRECT_H
#include <com_deepin_dde_daemon_dock.h>

DockRect::DockRect()
    : x(0)
    , y(0)
    , w(0)
    , h(0)
{
}

QDebug operator<<(QDebug debug, const DockRect &rect)
{
    debug << QString("DockRect(%1, %2, %3, %4)").arg(rect.x)
                                                .arg(rect.y)
                                                .arg(rect.w)
                                                .arg(rect.h);

    return debug;
}

DockRect::operator QRect() const
{
    return QRect(x, y, w, h);
}

QDBusArgument &operator<<(QDBusArgument &arg, const DockRect &rect)
{
    arg.beginStructure();
    arg << rect.x << rect.y << rect.w << rect.h;
    arg.endStructure();

    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, DockRect &rect)
{
    arg.beginStructure();
    arg >> rect.x >> rect.y >> rect.w >> rect.h;
    arg.endStructure();

    return arg;
}

void registerDockRectMetaType()
{
    qRegisterMetaType<DockRect>("DockRect");
    qDBusRegisterMetaType<DockRect>();
}

/*
 * Implementation of interface class __Dock
 */

class DockPrivate
{
public:
   DockPrivate() = default;

    // begin member variables
    int DisplayMode;
    QStringList DockedApps;
    QList<QDBusObjectPath> Entries;
    DockRect FrontendWindowRect;
    int HideMode;
    int HideState;
    uint HideTimeout;
    uint IconSize;
    double Opacity;
    int Position;
    uint ShowTimeout;
    uint WindowSize;
    uint WindowSizeEfficient;
    uint WindowSizeFashion;

public:
    QMap<QString, QDBusPendingCallWatcher *> m_processingCalls;
    QMap<QString, QList<QVariant>> m_waittingCalls;
};

// 窗管中提供的ActiveWindow接口，MinimizeWindow目前还在开发过程中，因此，关于这两个接口暂时使用v23的后端接口
// 等窗管完成了这几个接口后，删除此处v20的接口，改成v23提供的新接口即可
using DockInter = com::deepin::dde::daemon::Dock;
/**
 * @brief 任务栏的部分DBUS接口是通过窗管获取的，由于AM后端并未提供窗管的相关接口，因此，
 * 此处先将窗管的接口集成进来，作为私有类，只提供任务栏接口使用
 */
class WM : public QDBusAbstractInterface
{
public:
    static inline const char *staticInterfaceName()
    { return "com.deepin.wm"; }

public:
    explicit WM(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent = Q_NULLPTR);
    ~WM();

public Q_SLOTS: // METHODS

    inline QDBusPendingReply<> ActivateWindow(uint in0)
    {
        return m_dockInter->ActivateWindow(in0);
    }

    QDBusPendingReply<> MinimizeWindow(uint in0)
    {
        return m_dockInter->MinimizeWindow(in0);
    }

    inline QDBusPendingReply<> CancelPreviewWindow()
    {
        return asyncCallWithArgumentList(QStringLiteral("CancelPreviewWindow"), QList<QVariant>());
    }

    inline QDBusPendingReply<> PreviewWindow(uint in0)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0);
        return asyncCallWithArgumentList(QStringLiteral("CancelPreviewWindow"), argumentList);
    }

private:
    DockInter *m_dockInter;
};

WM::WM(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent)
    : QDBusAbstractInterface(service, path, staticInterfaceName(), connection, parent)
    , m_dockInter(new DockInter("com.deepin.dde.daemon.Dock", "/com/deepin/dde/daemon/Dock", QDBusConnection::sessionBus(), this))
{
}

WM::~WM()
{
}

Dde_Dock::Dde_Dock(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent)
    : QDBusAbstractInterface(service, path, staticInterfaceName(), connection, parent)
    , d_ptr(new DockPrivate)
    , m_wm(new WM("com.deepin.wm", "/com/deepin/wm", QDBusConnection::sessionBus(), this))
{
    QDBusConnection::sessionBus().connect(this->service(), this->path(),
                                          "org.freedesktop.DBus.Properties",
                                          "PropertiesChanged","sa{sv}as",
                                          this,
                                          SLOT(onPropertyChanged(const QDBusMessage &)));

    if (QMetaType::type("DockRect") == QMetaType::UnknownType)
        registerDockRectMetaType();
}

Dde_Dock::~Dde_Dock()
{
    qDeleteAll(d_ptr->m_processingCalls.values());
    delete d_ptr;
}

void Dde_Dock::onPropertyChanged(const QDBusMessage& msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QString interfaceName = msg.arguments().at(0).toString();
    if (interfaceName != staticInterfaceName())
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    QStringList keys = changedProps.keys();
    foreach(const QString &prop, keys) {
    const QMetaObject* self = metaObject();
        for (int i=self->propertyOffset(); i < self->propertyCount(); ++i) {
            QMetaProperty p = self->property(i);
            if (p.name() == prop)
                Q_EMIT p.notifySignal().invoke(this);
        }
    }
}

int Dde_Dock::displayMode()
{
    return qvariant_cast<int>(property("DisplayMode"));
}

void Dde_Dock::setDisplayMode(int value)
{
    setProperty("DisplayMode", QVariant::fromValue(value));
}

QStringList Dde_Dock::dockedApps()
{
    return qvariant_cast<QStringList>(property("DockedApps"));
}

QList<QDBusObjectPath> Dde_Dock::entries()
{
    return qvariant_cast<QList<QDBusObjectPath>>(property("Entries"));
}

DockRect Dde_Dock::frontendWindowRect()
{
    return qvariant_cast<DockRect>(property("FrontendWindowRect"));
}

int Dde_Dock::hideMode()
{
    return qvariant_cast<int>(property("HideMode"));
}

void Dde_Dock::setHideMode(int value)
{
   internalPropSet("HideMode", QVariant::fromValue(value));
}

int Dde_Dock::hideState()
{
    return qvariant_cast<int>(property("HideState"));
}

uint Dde_Dock::hideTimeout()
{
    return qvariant_cast<uint>(property("HideTimeout"));
}

void Dde_Dock::setHideTimeout(uint value)
{
   setProperty("HideTimeout", QVariant::fromValue(value));
}

uint Dde_Dock::iconSize()
{
    return qvariant_cast<uint>(property("IconSize"));
}

void Dde_Dock::setIconSize(uint value)
{
   setProperty("IconSize", QVariant::fromValue(value));
}

double Dde_Dock::opacity()
{
    return qvariant_cast<double>(property("Opacity"));
}

void Dde_Dock::setOpacity(double value)
{
   setProperty("Opacity", QVariant::fromValue(value));
}

int Dde_Dock::position()
{
    return qvariant_cast<int>(property("Position"));
}

void Dde_Dock::setPosition(int value)
{
   setProperty("Position", QVariant::fromValue(value));
}

uint Dde_Dock::showTimeout()
{
    return qvariant_cast<uint>(property("ShowTimeout"));
}

void Dde_Dock::setShowTimeout(uint value)
{
   setProperty("ShowTimeout", QVariant::fromValue(value));
}

uint Dde_Dock::windowSize()
{
    return qvariant_cast<uint>(property("WindowSize"));
}

void Dde_Dock::setWindowSize(uint value)
{
   setProperty("WindowSize", QVariant::fromValue(value));
}

uint Dde_Dock::windowSizeEfficient()
{
    return qvariant_cast<uint>(property("WindowSizeEfficient"));
}

void Dde_Dock::setWindowSizeEfficient(uint value)
{
   setProperty("WindowSizeEfficient", QVariant::fromValue(value));
}

uint Dde_Dock::windowSizeFashion()
{
    return qvariant_cast<uint>(property("WindowSizeFashion"));
}

void Dde_Dock::setWindowSizeFashion(uint value)
{
    setProperty("WindowSizeFashion", QVariant::fromValue(value));
}

QDBusPendingReply<> Dde_Dock::ActivateWindow(uint in0)
{
    return m_wm->ActivateWindow(in0);
}

QDBusPendingReply<> Dde_Dock::PreviewWindow(uint in0)
{
    return m_wm->PreviewWindow(in0);
}

QDBusPendingReply<> Dde_Dock::CancelPreviewWindow()
{
    return m_wm->CancelPreviewWindow();
}

QDBusPendingReply<> Dde_Dock::MinimizeWindow(uint in0)
{
    return m_wm->MinimizeWindow(in0);
}

void Dde_Dock::CallQueued(const QString &callName, const QList<QVariant> &args)
{
    if (d_ptr->m_waittingCalls.contains(callName)) {
        d_ptr->m_waittingCalls[callName] = args;
        return;
    }

    if (d_ptr->m_processingCalls.contains(callName)) {
        d_ptr->m_waittingCalls.insert(callName, args);
    } else {
        QDBusPendingCallWatcher *watcher = new QDBusPendingCallWatcher(asyncCallWithArgumentList(callName, args));
        connect(watcher, &QDBusPendingCallWatcher::finished, this, &Dde_Dock::onPendingCallFinished);
        d_ptr->m_processingCalls.insert(callName, watcher);
    }
}

void Dde_Dock::onPendingCallFinished(QDBusPendingCallWatcher *w)
{
    w->deleteLater();
    const auto callName = d_ptr->m_processingCalls.key(w);
    Q_ASSERT(!callName.isEmpty());
    if (callName.isEmpty())
        return;

    d_ptr->m_processingCalls.remove(callName);
    if (!d_ptr->m_waittingCalls.contains(callName))
        return;

    const auto args = d_ptr->m_waittingCalls.take(callName);
    CallQueued(callName, args);
}

#endif
