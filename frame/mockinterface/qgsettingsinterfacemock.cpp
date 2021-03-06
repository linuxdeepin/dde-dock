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
#include <QGSettings>
#include <QVariant>

#include "qgsettingsinterfacemock.h"

QGSettingsInterfaceMock::QGSettingsInterfaceMock(const QByteArray &schema_id, const QByteArray &path, QObject *parent)
{

}

QGSettingsInterfaceMock::~QGSettingsInterfaceMock()
{

}

QGSettingsInterface::Type QGSettingsInterfaceMock::type()
{
    return Type::FAKE;
}

QGSettings *QGSettingsInterfaceMock::gsettings()
{
    return nullptr;
}

QVariant QGSettingsInterfaceMock::get(const QString &key) const
{
    return QVariant();
}

void QGSettingsInterfaceMock::set(const QString &key, const QVariant &value)
{

}

bool QGSettingsInterfaceMock::trySet(const QString &key, const QVariant &value)
{
    return false;
}

QStringList QGSettingsInterfaceMock::keys() const
{
    return QStringList();
}

QVariantList QGSettingsInterfaceMock::choices(const QString &key) const
{
    return QVariantList();
}

void QGSettingsInterfaceMock::reset(const QString &key)
{

}

bool QGSettingsInterfaceMock::isSchemaInstalled(const QByteArray &schema_id)
{
    return false;
}
