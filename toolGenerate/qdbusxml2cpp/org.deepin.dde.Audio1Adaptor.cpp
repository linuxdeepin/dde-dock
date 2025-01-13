/*
 * This file was generated by qdbusxml2cpp version 0.8
 * Command line was: qdbusxml2cpp ./dde-dock/frame/dbusinterface/xml/org.deepin.dde.Audio1.xml -a ./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.Audio1Adaptor -i ./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.Audio1.h
 *
 * qdbusxml2cpp is Copyright (C) 2017 The Qt Company Ltd.
 *
 * This is an auto-generated file.
 * Do not edit! All changes made to it will be lost.
 */

#include "./dde-dock/toolGenerate/qdbusxml2cpp/org.deepin.dde.Audio1Adaptor.h"
#include <QtCore/QMetaObject>
#include <QtCore/QByteArray>
#include <QtCore/QList>
#include <QtCore/QMap>
#include <QtCore/QString>
#include <QtCore/QStringList>
#include <QtCore/QVariant>

/*
 * Implementation of adaptor class Audio1Adaptor
 */

Audio1Adaptor::Audio1Adaptor(QObject *parent)
    : QDBusAbstractAdaptor(parent)
{
    // constructor
    setAutoRelaySignals(true);
}

Audio1Adaptor::~Audio1Adaptor()
{
    // destructor
}

QString Audio1Adaptor::bluetoothAudioMode() const
{
    // get the value of property BluetoothAudioMode
    return qvariant_cast< QString >(parent()->property("BluetoothAudioMode"));
}

QStringList Audio1Adaptor::bluetoothAudioModeOpts() const
{
    // get the value of property BluetoothAudioModeOpts
    return qvariant_cast< QStringList >(parent()->property("BluetoothAudioModeOpts"));
}

QString Audio1Adaptor::cards() const
{
    // get the value of property Cards
    return qvariant_cast< QString >(parent()->property("Cards"));
}

QString Audio1Adaptor::cardsWithoutUnavailable() const
{
    // get the value of property CardsWithoutUnavailable
    return qvariant_cast< QString >(parent()->property("CardsWithoutUnavailable"));
}

QDBusObjectPath Audio1Adaptor::defaultSink() const
{
    // get the value of property DefaultSink
    return qvariant_cast< QDBusObjectPath >(parent()->property("DefaultSink"));
}

QDBusObjectPath Audio1Adaptor::defaultSource() const
{
    // get the value of property DefaultSource
    return qvariant_cast< QDBusObjectPath >(parent()->property("DefaultSource"));
}

bool Audio1Adaptor::increaseVolume() const
{
    // get the value of property IncreaseVolume
    return qvariant_cast< bool >(parent()->property("IncreaseVolume"));
}

void Audio1Adaptor::setIncreaseVolume(bool value)
{
    // set the value of property IncreaseVolume
    parent()->setProperty("IncreaseVolume", QVariant::fromValue(value));
}

double Audio1Adaptor::maxUIVolume() const
{
    // get the value of property MaxUIVolume
    return qvariant_cast< double >(parent()->property("MaxUIVolume"));
}

bool Audio1Adaptor::reduceNoise() const
{
    // get the value of property ReduceNoise
    return qvariant_cast< bool >(parent()->property("ReduceNoise"));
}

void Audio1Adaptor::setReduceNoise(bool value)
{
    // set the value of property ReduceNoise
    parent()->setProperty("ReduceNoise", QVariant::fromValue(value));
}

QList<QDBusObjectPath> Audio1Adaptor::sinkInputs() const
{
    // get the value of property SinkInputs
    return qvariant_cast< QList<QDBusObjectPath> >(parent()->property("SinkInputs"));
}

QList<QDBusObjectPath> Audio1Adaptor::sinks() const
{
    // get the value of property Sinks
    return qvariant_cast< QList<QDBusObjectPath> >(parent()->property("Sinks"));
}

QList<QDBusObjectPath> Audio1Adaptor::sources() const
{
    // get the value of property Sources
    return qvariant_cast< QList<QDBusObjectPath> >(parent()->property("Sources"));
}

bool Audio1Adaptor::IsPortEnabled(uint in0, const QString &in1)
{
    // handle method call org.deepin.dde.Audio1.IsPortEnabled
    bool out0;
    QMetaObject::invokeMethod(parent(), "IsPortEnabled", Q_RETURN_ARG(bool, out0), Q_ARG(uint, in0), Q_ARG(QString, in1));
    return out0;
}

void Audio1Adaptor::Reset()
{
    // handle method call org.deepin.dde.Audio1.Reset
    QMetaObject::invokeMethod(parent(), "Reset");
}

void Audio1Adaptor::SetBluetoothAudioMode(const QString &in0)
{
    // handle method call org.deepin.dde.Audio1.SetBluetoothAudioMode
    QMetaObject::invokeMethod(parent(), "SetBluetoothAudioMode", Q_ARG(QString, in0));
}

void Audio1Adaptor::SetDefaultSink(const QString &in0)
{
    // handle method call org.deepin.dde.Audio1.SetDefaultSink
    QMetaObject::invokeMethod(parent(), "SetDefaultSink", Q_ARG(QString, in0));
}

void Audio1Adaptor::SetDefaultSource(const QString &in0)
{
    // handle method call org.deepin.dde.Audio1.SetDefaultSource
    QMetaObject::invokeMethod(parent(), "SetDefaultSource", Q_ARG(QString, in0));
}

void Audio1Adaptor::SetPort(uint in0, const QString &in1, int in2)
{
    // handle method call org.deepin.dde.Audio1.SetPort
    QMetaObject::invokeMethod(parent(), "SetPort", Q_ARG(uint, in0), Q_ARG(QString, in1), Q_ARG(int, in2));
}

void Audio1Adaptor::SetPortEnabled(uint in0, const QString &in1, bool in2)
{
    // handle method call org.deepin.dde.Audio1.SetPortEnabled
    QMetaObject::invokeMethod(parent(), "SetPortEnabled", Q_ARG(uint, in0), Q_ARG(QString, in1), Q_ARG(bool, in2));
}

