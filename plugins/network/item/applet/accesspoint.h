/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#ifndef ACCESSPOINT_H
#define ACCESSPOINT_H

#include <QObject>
#include <QJsonObject>

class AccessPoint : public QObject
{
    Q_OBJECT

public:
    explicit AccessPoint();
    explicit AccessPoint(const QJsonObject &apInfo);
    AccessPoint(const AccessPoint &ap);
    bool operator==(const AccessPoint &ap) const;
    bool operator>(const AccessPoint &ap) const;
    AccessPoint &operator=(const AccessPoint &ap);
    const QString ssid() const;
    const QString path() const;
    const QString uuid() const;
    int strength() const;
    bool secured() const;

    bool isEmpty() const;
    void updateApInfo(const QJsonObject &apInfo);

private:
    //信号强度
    int m_strength;
    //是否需要密码
    bool m_secured;
    bool m_securedInEap;
    //apPath
    QString m_path;
    //apSsid
    QString m_ssid;
    //apUuid
    QString m_uuid;
};

#endif // ACCESSPOINT_H
