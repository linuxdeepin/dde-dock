// Copyright (C) 2011 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef BRIGHTNESSMODEL_H
#define BRIGHTNESSMODEL_H

#include <QDBusObjectPath>
#include <QObject>

class BrightMonitor;
class QDBusMessage;
class QScreen;

class BrightnessModel : public QObject
{
    Q_OBJECT

public:
    explicit BrightnessModel(QObject *parent = nullptr);
    ~BrightnessModel();

    QList<BrightMonitor *> monitors();
    BrightMonitor *primaryMonitor() const;

Q_SIGNALS:
    void primaryChanged(BrightMonitor *);
    void screenVisibleChanged(bool);
    void monitorLightChanged();

protected Q_SLOTS:
    void primaryScreenChanged(QScreen *screen);
    void onPropertyChanged(const QDBusMessage &msg);

private:
    QList<BrightMonitor *> readMonitors(const QList<QDBusObjectPath> &paths);

private:
    QList<BrightMonitor *> m_monitor;
    QString m_primaryScreenName;
};

class BrightMonitor : public QObject
{
    Q_OBJECT

public:
    explicit BrightMonitor(QString path, QObject *parent);
    ~BrightMonitor();

Q_SIGNALS:
    void brightnessChanged(int);
    void nameChanged(QString);
    void enabledChanged(bool);

public:
    void setPrimary(bool primary);
    int brightness();
    bool enabled();
    QString name();
    bool isPrimary();

public slots:
    void setBrightness(int value);
    void onPropertyChanged(const QDBusMessage &msg);

private:
    QDBusMessage callMethod(const QString &methodName, const QList<QVariant> &argument);

private:
    QString m_path;
    QString m_name;
    int m_brightness;
    bool m_enabled;
    bool m_isPrimary;
};

#endif // DISPLAYMODEL_H
