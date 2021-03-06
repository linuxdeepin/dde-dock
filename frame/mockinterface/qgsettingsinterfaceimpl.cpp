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

#include "qgsettingsinterfaceimpl.h"

QGSettingsInterfaceImpl::QGSettingsInterfaceImpl(const QByteArray &schema_id, const QByteArray &path, QObject *parent)
    : m_gsettings(new QGSettings(schema_id, path, parent))
{

}

QGSettingsInterfaceImpl::~QGSettingsInterfaceImpl()
{

}

QGSettingsInterface::Type QGSettingsInterfaceImpl::type()
{
     return Type::REAL;
}

QGSettings *QGSettingsInterfaceImpl::gsettings()
{
    return m_gsettings;
}

QVariant QGSettingsInterfaceImpl::get(const QString &key) const
{
    return m_gsettings->get(key);
}

void QGSettingsInterfaceImpl::set(const QString &key, const QVariant &value)
{
    return m_gsettings->set(key, value);
}

bool QGSettingsInterfaceImpl::trySet(const QString &key, const QVariant &value)
{
    return m_gsettings->trySet(key, value);
}

QStringList QGSettingsInterfaceImpl::keys() const
{
    return m_gsettings->keys();
}

QVariantList QGSettingsInterfaceImpl::choices(const QString &key) const
{
    return m_gsettings->choices(key);
}

void QGSettingsInterfaceImpl::reset(const QString &key)
{
    return m_gsettings->reset(key);
}

bool QGSettingsInterfaceImpl::isSchemaInstalled(const QByteArray &schema_id)
{
    return QGSettings::isSchemaInstalled(schema_id);
}
