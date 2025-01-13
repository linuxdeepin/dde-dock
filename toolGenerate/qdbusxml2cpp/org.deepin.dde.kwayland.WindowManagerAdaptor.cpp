/*
 * This file was generated by qdbusxml2cpp version 0.8
 * Command line was: qdbusxml2cpp ./dde-dock/frame/dbusinterface/xml/org.deepin.dde.kwayland.WindowManager.xml -a ./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.kwayland.WindowManagerAdaptor -i ./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.kwayland.WindowManager.h
 *
 * qdbusxml2cpp is Copyright (C) 2017 The Qt Company Ltd.
 *
 * This is an auto-generated file.
 * Do not edit! All changes made to it will be lost.
 */

#include "./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.kwayland.WindowManagerAdaptor.h"
#include <QtCore/QMetaObject>
#include <QtCore/QByteArray>
#include <QtCore/QList>
#include <QtCore/QMap>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QVariant>

/*
 * Implementation of adaptor class WindowManagerAdaptor
 */

WindowManagerAdaptor::WindowManagerAdaptor(QObject *parent)
    : QDBusAbstractAdaptor(parent)
{
    // constructor
    setAutoRelaySignals(true);
}

WindowManagerAdaptor::~WindowManagerAdaptor()
{
    // destructor
}

uint WindowManagerAdaptor::ActiveWindow()
{
    // handle method call org.deepin.dde.KWayland1.WindowManager.ActiveWindow
    uint out0;
    QMetaObject::invokeMethod(parent(), "ActiveWindow", Q_RETURN_ARG(uint, out0));
    return out0;
}

void WindowManagerAdaptor::HideDesktop()
{
    // handle method call org.deepin.dde.KWayland1.WindowManager.HideDesktop
    QMetaObject::invokeMethod(parent(), "HideDesktop");
}

bool WindowManagerAdaptor::IsShowingDesktop()
{
    // handle method call org.deepin.dde.KWayland1.WindowManager.IsShowingDesktop
    bool out0;
    QMetaObject::invokeMethod(parent(), "IsShowingDesktop", Q_RETURN_ARG(bool, out0));
    return out0;
}

bool WindowManagerAdaptor::IsValid()
{
    // handle method call org.deepin.dde.KWayland1.WindowManager.IsValid
    bool out0;
    QMetaObject::invokeMethod(parent(), "IsValid", Q_RETURN_ARG(bool, out0));
    return out0;
}

void WindowManagerAdaptor::SetShowingDesktop(bool show)
{
    // handle method call org.deepin.dde.KWayland1.WindowManager.SetShowingDesktop
    QMetaObject::invokeMethod(parent(), "SetShowingDesktop", Q_ARG(bool, show));
}

void WindowManagerAdaptor::ShowDesktop()
{
    // handle method call org.deepin.dde.KWayland1.WindowManager.ShowDesktop
    QMetaObject::invokeMethod(parent(), "ShowDesktop");
}

QVariantList WindowManagerAdaptor::Windows()
{
    // handle method call org.deepin.dde.KWayland1.WindowManager.Windows
    QVariantList out0;
    QMetaObject::invokeMethod(parent(), "Windows", Q_RETURN_ARG(QVariantList, out0));
    return out0;
}
