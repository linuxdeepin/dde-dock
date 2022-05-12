/*
 * Copyright (C) 2011 ~ 2022 Deepin Technology Co., Ltd.
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

#ifndef BRIGHTNESSMODEL_H
#define BRIGHTNESSMODEL_H

#include <QObject>

class BrightMonitor;
class QDBusMessage;

class BrightnessModel : public QObject
{
    Q_OBJECT

Q_SIGNALS:
    void brightnessChanged(BrightMonitor *);

public:
    explicit BrightnessModel(QObject *parent = nullptr);
    ~BrightnessModel();

    QList<BrightMonitor *> monitors();
    void setBrightness(BrightMonitor *monitor, int brightness);
    void setBrightness(QString name, int brightness);

protected:
    QDBusMessage callMethod(const QString &methodName, const QList<QVariant> &argument);

protected Q_SLOTS:
    void onPropertyChanged(const QDBusMessage &msg);

private:
    QList<BrightMonitor *> m_monitor;
};

class BrightMonitor : public QObject
{
    Q_OBJECT

    friend class BrightnessModel;

Q_SIGNALS:
    void brightnessChanged(int);
    void nameChanged(QString);
    void enabledChanged(bool);

public:
    int brihtness();
    bool enabled();
    QString name();
    //bool isDefault();

protected:
    explicit BrightMonitor(QString path, QObject *parent);
    ~BrightMonitor();

protected Q_SLOTS:
    void onPropertyChanged(const QDBusMessage &msg);

private:
    QString m_path;
    QString m_name;
    int m_brightness;
    bool m_enabled;
};

#endif // DISPLAYMODEL_H
