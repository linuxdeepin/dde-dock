/*
 * This file was generated by qdbusxml2cpp version 0.8
 * Command line was: qdbusxml2cpp ./dde-dock/plugins/power/dbus/org.deepin.dde.Power1.xml -a ./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.Power1Adaptor -i ./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.Power1.h
 *
 * qdbusxml2cpp is Copyright (C) 2017 The Qt Company Ltd.
 *
 * This is an auto-generated file.
 * Do not edit! All changes made to it will be lost.
 */

#include "./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.Power1Adaptor.h"
#include <QtCore/QMetaObject>
#include <QtCore/QByteArray>
#include <QtCore/QList>
#include <QtCore/QMap>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QVariant>

/*
 * Implementation of adaptor class Power1Adaptor
 */

Power1Adaptor::Power1Adaptor(QObject *parent)
    : QDBusAbstractAdaptor(parent)
{
    // constructor
    setAutoRelaySignals(true);
}

Power1Adaptor::~Power1Adaptor()
{
    // destructor
}

BatteryPercentageMap Power1Adaptor::batteryPercentage() const
{
    // get the value of property BatteryPercentage
    return qvariant_cast< BatteryPercentageMap >(parent()->property("BatteryPercentage"));
}

