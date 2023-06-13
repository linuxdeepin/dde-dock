// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#pragma once

#include "taskmanager/entry.h"
#include <QtCore/QObject>
#include <QtCore/QMetaObject>
#include <QtCore/QVariant>
#include <QtDBus/QtDBus>
#include <QtCore/QByteArray>
#include <QtCore/QList>
#include <QtCore/QMap>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QVariant>
#include <QDBusObjectPath>
#include <QRect>

/*
 * Adaptor class for interface org.deepin.dde.daemon.Dock1
 */

class Entry;
class DockDaemonDBusAdaptor: public QDBusAbstractAdaptor
{
    Q_OBJECT
    Q_CLASSINFO("D-Bus Interface", "org.deepin.dde.daemon.Dock1")
    Q_CLASSINFO("D-Bus Introspection", ""
                                       "  <interface name=\"org.deepin.dde.daemon.Dock1\">\n"
                                       "    <method name=\"CloseWindow\">\n"
                                       "      <arg direction=\"in\" type=\"u\" name=\"win\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"GetEntryIDs\">\n"
                                       "      <arg direction=\"out\" type=\"as\" name=\"list\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"IsDocked\">\n"
                                       "      <arg direction=\"in\" type=\"s\" name=\"desktopFile\"/>\n"
                                       "      <arg direction=\"out\" type=\"b\" name=\"value\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"IsOnDock\">\n"
                                       "      <arg direction=\"in\" type=\"s\" name=\"desktopFile\"/>\n"
                                       "      <arg direction=\"out\" type=\"b\" name=\"value\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"MoveEntry\">\n"
                                       "      <arg direction=\"in\" type=\"i\" name=\"index\"/>\n"
                                       "      <arg direction=\"in\" type=\"i\" name=\"newIndex\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"QueryWindowIdentifyMethod\">\n"
                                       "      <arg direction=\"in\" type=\"u\" name=\"win\"/>\n"
                                       "      <arg direction=\"out\" type=\"s\" name=\"identifyMethod\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"GetDockedAppsDesktopFiles\">\n"
                                       "      <arg direction=\"out\" type=\"as\" name=\"desktopFiles\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"GetPluginSettings\">\n"
                                       "      <arg direction=\"out\" type=\"s\" name=\"jsonStr\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"SetPluginSettings\">\n"
                                       "      <arg direction=\"in\" type=\"s\" name=\"jsonStr\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"MergePluginSettings\">\n"
                                       "      <arg direction=\"in\" type=\"s\" name=\"jsonStr\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"RemovePluginSettings\">\n"
                                       "      <arg direction=\"in\" type=\"s\" name=\"key1\"/>\n"
                                       "      <arg direction=\"in\" type=\"as\" name=\"key2List\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"RequestDock\">\n"
                                       "      <arg direction=\"in\" type=\"s\" name=\"desktopFile\"/>\n"
                                       "      <arg direction=\"in\" type=\"i\" name=\"index\"/>\n"
                                       "      <arg direction=\"out\" type=\"b\" name=\"ok\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"RequestUndock\">\n"
                                       "      <arg direction=\"in\" type=\"s\" name=\"desktopFile\"/>\n"
                                       "      <arg direction=\"out\" type=\"b\" name=\"ok\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"SetShowRecent\">\n"
                                       "      <arg direction=\"in\" type=\"b\" name=\"visible\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"SetShowMultiWindow\">\n"
                                       "      <arg direction=\"in\" type=\"b\" name=\"visible\"/>\n"
                                       "    </method>\n"
                                       "    <method name=\"SetFrontendWindowRect\">\n"
                                       "      <arg direction=\"in\" type=\"i\" name=\"x\"/>\n"
                                       "      <arg direction=\"in\" type=\"i\" name=\"y\"/>\n"
                                       "      <arg direction=\"in\" type=\"u\" name=\"width\"/>\n"
                                       "      <arg direction=\"in\" type=\"u\" name=\"height\"/>\n"
                                       "    </method>\n"
                                       "    <signal name=\"ServiceRestarted\"/>\n"
                                       "    <signal name=\"EntryAdded\">\n"
                                       "      <arg type=\"o\" name=\"path\"/>\n"
                                       "      <arg type=\"i\" name=\"index\"/>\n"
                                       "    </signal>\n"
                                       "    <signal name=\"EntryRemoved\">\n"
                                       "      <arg type=\"s\" name=\"entryId\"/>\n"
                                       "    </signal>\n"
                                       "    <property access=\"readwrite\" type=\"u\" name=\"ShowTimeout\"/>\n"
                                       "    <property access=\"readwrite\" type=\"u\" name=\"HideTimeout\"/>\n"
                                       "    <property access=\"readwrite\" type=\"u\" name=\"WindowSizeEfficient\"/>\n"
                                       "    <property access=\"readwrite\" type=\"u\" name=\"WindowSizeFashion\"/>\n"
                                       "    <property access=\"read\" type=\"(iiii)\" name=\"FrontendWindowRect\"/>\n"
                                       "    <property access=\"readwrite\" type=\"i\" name=\"HideMode\"/>\n"
                                       "    <property access=\"readwrite\" type=\"i\" name=\"DisplayMode\"/>\n"
                                       "    <property access=\"read\" type=\"i\" name=\"HideState\"/>\n"
                                       "    <property access=\"readwrite\" type=\"i\" name=\"Position\"/>\n"
                                       "    <property access=\"readwrite\" type=\"u\" name=\"IconSize\"/>\n"
                                       "    <property access=\"read\" type=\"as\" name=\"DockedApps\"/>\n"
                                       "    <property access=\"read\" type=\"b\" name=\"ShowRecent\"/>\n"
                                       "    <property access=\"read\" type=\"b\" name=\"ShowMultiWindow\"/>\n"
                                       "  </interface>\n"
                                       "")
public:
    DockDaemonDBusAdaptor(QObject *parent);
    virtual ~DockDaemonDBusAdaptor();

public: // PROPERTIES
    Q_PROPERTY(int DisplayMode READ displayMode WRITE setDisplayMode NOTIFY DisplayModeChanged)
    int displayMode() const;
    void setDisplayMode(int value);

    Q_PROPERTY(QStringList DockedApps READ dockedApps NOTIFY DockedAppsChanged)
    QStringList dockedApps() const;

    Q_PROPERTY(int HideMode READ hideMode WRITE setHideMode NOTIFY HideModeChanged)
    int hideMode() const;
    void setHideMode(int value);

    Q_PROPERTY(int HideState READ hideState NOTIFY HideStateChanged)
    int hideState() const;

    Q_PROPERTY(uint HideTimeout READ hideTimeout WRITE setHideTimeout NOTIFY HideTimeoutChanged)
    uint hideTimeout() const;
    void setHideTimeout(uint value);

    Q_PROPERTY(uint WindowSizeEfficient READ windowSizeEfficient WRITE setWindowSizeEfficient NOTIFY WindowSizeEfficientChanged)
    uint windowSizeEfficient() const;
    void setWindowSizeEfficient(uint value);

    Q_PROPERTY(uint WindowSizeFashion READ windowSizeFashion WRITE setWindowSizeFashion NOTIFY WindowSizeFashionChanged)
    uint windowSizeFashion() const;
    void setWindowSizeFashion(uint value);

    Q_PROPERTY(QRect FrontendWindowRect READ frontendWindowRect NOTIFY FrontendWindowRectChanged)
    QRect frontendWindowRect() const;

    Q_PROPERTY(uint IconSize READ iconSize WRITE setIconSize NOTIFY IconSizeChanged)
    uint iconSize() const;
    void setIconSize(uint value);

    Q_PROPERTY(int Position READ position WRITE setPosition NOTIFY PositionChanged)
    int position() const;
    void setPosition(int value);

    Q_PROPERTY(uint ShowTimeout READ showTimeout WRITE setShowTimeout NOTIFY ShowTimeoutChanged)
    uint showTimeout() const;
    void setShowTimeout(uint value);

    Q_PROPERTY(bool ShowRecent READ showRecent NOTIFY showRecentChanged)
    bool showRecent() const;

    Q_PROPERTY(bool ShowMultiWindow READ showMultiWindow NOTIFY ShowMultiWindowChanged)
    bool showMultiWindow() const;

public Q_SLOTS: // METHODS
    void CloseWindow(uint win);
    QStringList GetEntryIDs();
    bool IsDocked(const QString &desktopFile);
    bool IsOnDock(const QString &desktopFile);
    void MoveEntry(int index, int newIndex);
    QString QueryWindowIdentifyMethod(uint win);
    QStringList GetDockedAppsDesktopFiles();
    QString GetPluginSettings();
    void SetPluginSettings(QString jsonStr);
    void MergePluginSettings(QString jsonStr);
    void RemovePluginSettings(QString key1, QStringList key2List);
    bool RequestDock(const QString &desktopFile, int index);
    bool RequestUndock(const QString &desktopFile);
    void SetShowRecent(bool visible);
    void SetShowMultiWindow(bool showMultiWindow);
    void SetFrontendWindowRect(int x, int y, uint width, uint height);

Q_SIGNALS: // SIGNALS
    void ServiceRestarted();
    void EntryAdded(const Entry *entry, int index);
    void EntryRemoved(const QString &entryId);

    void DisplayModeChanged(int value) const;
    void DockedAppsChanged(const QStringList &value) const;
    void EntriesChanged(const QList<QDBusObjectPath> &value) const;
    void FrontendWindowRectChanged(const QRect &dockRect) const;
    void HideModeChanged(int value) const;
    void HideStateChanged(int value) const;
    void HideTimeoutChanged(uint value) const;
    void IconSizeChanged(uint value) const;
    void PositionChanged(int value) const;
    void ShowTimeoutChanged(uint value) const;
    void WindowSizeEfficientChanged(uint value) const;
    void WindowSizeFashionChanged(uint value) const;
    void showRecentChanged(bool) const;
    void ShowMultiWindowChanged(bool) const;
};