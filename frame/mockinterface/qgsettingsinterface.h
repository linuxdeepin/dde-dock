/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng@uniontech.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng@uniontech.com>
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

#ifndef QGSETTINGSINTERFACE_H
#define QGSETTINGSINTERFACE_H

#include <QVariant>
#include <QStringList>

class QGSettings;
class QGSettingsInterface
{
public:
    enum Type {
        REAL,   // 持有真正的QGSettings指针
        FAKE    // Mock类
    };

    virtual ~QGSettingsInterface() {}

    virtual Type type() = 0;
    virtual QGSettings *gsettings() = 0;
    virtual QVariant get(const QString &key) const = 0;
    virtual void set(const QString &key, const QVariant &value) = 0;
    virtual bool trySet(const QString &key, const QVariant &value) = 0;
    virtual QStringList keys() const = 0;
    virtual QVariantList choices(const QString &key) const = 0;
    virtual void reset(const QString &key) = 0;
    static bool isSchemaInstalled(const QByteArray &schema_id) {Q_UNUSED(schema_id); return false;}

};
#endif // QGSETTINGSINTERFACE_H
