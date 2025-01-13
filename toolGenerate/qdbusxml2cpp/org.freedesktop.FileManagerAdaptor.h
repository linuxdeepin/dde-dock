/*
 * This file was generated by qdbusxml2cpp version 0.8
 * Command line was: qdbusxml2cpp ./dde-dock/frame/dbusinterface/xml/org.freedesktop.FileManager.xml -a ./dde-dock/toolGenerate/qdbusxml2cpp/org.freedesktop.FileManagerAdaptor -i ./dde-dock/toolGenerate/qdbusxml2cpp/org.freedesktop.FileManager.h
 *
 * qdbusxml2cpp is Copyright (C) 2017 The Qt Company Ltd.
 *
 * This is an auto-generated file.
 * This file may have been hand-edited. Look for HAND-EDIT comments
 * before re-generating it.
 */

#ifndef ORG_FREEDESKTOP_FILEMANAGERADAPTOR_H
#define ORG_FREEDESKTOP_FILEMANAGERADAPTOR_H

#include <QtCore/QObject>
#include <QtDBus/QtDBus>
#include "./dde-dock/toolGenerate/qdbusxml2cpp/org.freedesktop.FileManager.h"
QT_BEGIN_NAMESPACE
class QByteArray;
template<class T> class QList;
template<class Key, class Value> class QMap;
class QString;
class QStringList;
class QVariant;
QT_END_NAMESPACE

/*
 * Adaptor class for interface org.freedesktop.FileManager1
 */
class FileManager1Adaptor: public QDBusAbstractAdaptor
{
    Q_OBJECT
    Q_CLASSINFO("D-Bus Interface", "org.freedesktop.FileManager1")
    Q_CLASSINFO("D-Bus Introspection", ""
"  <interface name=\"org.freedesktop.FileManager1\">\n"
"    <method name=\"ShowFolders\">\n"
"      <arg direction=\"in\" type=\"as\" name=\"URIs\"/>\n"
"      <arg direction=\"in\" type=\"s\" name=\"StartupId\"/>\n"
"    </method>\n"
"    <method name=\"ShowItems\">\n"
"      <arg direction=\"in\" type=\"as\" name=\"URIs\"/>\n"
"      <arg direction=\"in\" type=\"s\" name=\"StartupId\"/>\n"
"    </method>\n"
"    <method name=\"ShowItemProperties\">\n"
"      <arg direction=\"in\" type=\"as\" name=\"URIs\"/>\n"
"      <arg direction=\"in\" type=\"s\" name=\"StartupId\"/>\n"
"    </method>\n"
"    <method name=\"Trash\">\n"
"      <arg direction=\"in\" type=\"as\" name=\"URIs\"/>\n"
"    </method>\n"
"  </interface>\n"
        "")
public:
    FileManager1Adaptor(QObject *parent);
    virtual ~FileManager1Adaptor();

public: // PROPERTIES
public Q_SLOTS: // METHODS
    void ShowFolders(const QStringList &URIs, const QString &StartupId);
    void ShowItemProperties(const QStringList &URIs, const QString &StartupId);
    void ShowItems(const QStringList &URIs, const QString &StartupId);
    void Trash(const QStringList &URIs);
Q_SIGNALS: // SIGNALS
};

#endif