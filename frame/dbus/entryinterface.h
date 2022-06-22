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

#ifndef DOCK_ENTRY_H
#define DOCK_ENTRY_H

#ifdef USE_AM

#define WINDOWLIST_H
#define WINDOWINFOLIST_H

#include <QObject>
#include <QByteArray>
#include <QList>
#include <QMap>
#include <QString>
#include <QStringList>
#include <QVariant>
#include <DBusExtendedAbstractInterface>
#include <QtDBus>

typedef QList<quint32> WindowList;

void registerWindowListMetaType();

class WindowInfo
{
public:
    friend QDebug operator<<(QDebug argument, const WindowInfo &info);
    friend QDBusArgument &operator<<(QDBusArgument &argument, const WindowInfo &info);
    friend const QDBusArgument &operator>>(const QDBusArgument &argument, WindowInfo &info);

    bool operator==(const WindowInfo &rhs) const;

public:
    bool attention;
    QString title;
};

Q_DECLARE_METATYPE(WindowInfo)

typedef QMap<quint32, WindowInfo> WindowInfoMap;
Q_DECLARE_METATYPE(WindowInfoMap)

void registerWindowInfoMetaType();
void registerWindowInfoMapMetaType();

/*
 * Proxy class for interface com.deepin.dde.daemon.Dock.Entry
 */
class EntryPrivate;

class Dock_Entry : public QDBusAbstractInterface
{
    Q_OBJECT

public:
    static inline const char *staticInterfaceName()
    { return "org.deepin.dde.daemon.Dock1.Entry"; }

public:
    explicit Dock_Entry(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent = 0);

    ~Dock_Entry();

    Q_PROPERTY(uint CurrentWindow READ currentWindow NOTIFY CurrentWindowChanged)
    uint currentWindow();

    Q_PROPERTY(QString DesktopFile READ desktopFile NOTIFY DesktopFileChanged)
    QString desktopFile();

    Q_PROPERTY(QString Icon READ icon NOTIFY IconChanged)
    QString icon();

    Q_PROPERTY(QString Id READ id NOTIFY IdChanged)
    QString id();

    Q_PROPERTY(bool IsActive READ isActive NOTIFY IsActiveChanged)
    bool isActive();

    Q_PROPERTY(bool IsDocked READ isDocked NOTIFY IsDockedChanged)
    bool isDocked();

    Q_PROPERTY(QString Menu READ menu NOTIFY MenuChanged)
    QString menu();

    Q_PROPERTY(QString Name READ name NOTIFY NameChanged)
    QString name();

    Q_PROPERTY(WindowInfoMap WindowInfos READ windowInfos NOTIFY WindowInfosChanged)
    WindowInfoMap windowInfos();

public Q_SLOTS: // METHODS
    inline QDBusPendingReply<> Activate(uint in0)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0);
        return asyncCallWithArgumentList(QStringLiteral("Activate"), argumentList);
    }

    inline void ActivateQueued(uint in0)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0);
        CallQueued(QStringLiteral("Activate"), argumentList);
    }

    inline QDBusPendingReply<> Check()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("Check"), argumentList);
    }

    inline void CheckQueued()
    {
        QList<QVariant> argumentList;
        CallQueued(QStringLiteral("Check"), argumentList);
    }

    inline QDBusPendingReply<> ForceQuit()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("ForceQuit"), argumentList);
    }

    inline void ForceQuitQueued()
    {
        QList<QVariant> argumentList;
        CallQueued(QStringLiteral("ForceQuit"), argumentList);
    }

    inline QDBusPendingReply<WindowList> GetAllowedCloseWindows()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("GetAllowedCloseWindows"), argumentList);
    }

    inline QDBusPendingReply<> HandleDragDrop(uint in0, const QStringList &in1)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0) << QVariant::fromValue(in1);
        return asyncCallWithArgumentList(QStringLiteral("HandleDragDrop"), argumentList);
    }

    inline void HandleDragDropQueued(uint in0, const QStringList &in1)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0) << QVariant::fromValue(in1);
        CallQueued(QStringLiteral("HandleDragDrop"), argumentList);
    }

    inline QDBusPendingReply<> HandleMenuItem(uint in0, const QString &in1)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0) << QVariant::fromValue(in1);
        return asyncCallWithArgumentList(QStringLiteral("HandleMenuItem"), argumentList);
    }

    inline void HandleMenuItemQueued(uint in0, const QString &in1)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0) << QVariant::fromValue(in1);
        CallQueued(QStringLiteral("HandleMenuItem"), argumentList);
    }

    inline QDBusPendingReply<> NewInstance(uint in0)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0);
        return asyncCallWithArgumentList(QStringLiteral("NewInstance"), argumentList);
    }

    inline void NewInstanceQueued(uint in0)
    {
        QList<QVariant> argumentList;
        argumentList << QVariant::fromValue(in0);
        CallQueued(QStringLiteral("NewInstance"), argumentList);
    }

    inline QDBusPendingReply<> PresentWindows()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("PresentWindows"), argumentList);
    }

    inline void PresentWindowsQueued()
    {
        QList<QVariant> argumentList;
        CallQueued(QStringLiteral("PresentWindows"), argumentList);
    }

    inline QDBusPendingReply<> RequestDock()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("RequestDock"), argumentList);
    }

    inline void RequestDockQueued()
    {
        QList<QVariant> argumentList;
        CallQueued(QStringLiteral("RequestDock"), argumentList);
    }

    inline QDBusPendingReply<> RequestUndock()
    {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("RequestUndock"), argumentList);
    }

    inline void RequestUndockQueued()
    {
        QList<QVariant> argumentList;
        CallQueued(QStringLiteral("RequestUndock"), argumentList);
    }

Q_SIGNALS: // SIGNALS
    // begin property changed signals
    void IsActiveChanged(bool value) const;
    void IsDockedChanged(bool value) const;
    void MenuChanged(const QString &value) const;
    void IconChanged(const QString &value) const;
    void NameChanged(const QString &value) const;
    void DesktopFileChanged(const QString &value) const;
    void CurrentWindowChanged(uint32_t value) const;

    void WindowInfosChanged(WindowInfoMap value) const;
    void IdChanged(const QString &value) const;

private:
    QVariant asyncProperty(const QString &propertyName);

public Q_SLOTS:
    void CallQueued(const QString &callName, const QList<QVariant> &args);

private Q_SLOTS:
    void onPendingCallFinished(QDBusPendingCallWatcher *w);
    void onPropertyChanged(const QString &propName, const QVariant &value);

private:
    EntryPrivate *d_ptr;
};

namespace org {
  namespace deepin {
    namespace dde {
      namespace daemon {
        namespace dock {
          typedef ::Dock_Entry DockEntry;
        }
      }
    }
  }
}

#endif  // USE_AM

#endif  // DOCK_ENTRY_H
