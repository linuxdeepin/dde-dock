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
#ifndef QGSETTINGSMOCK_H
#define QGSETTINGSMOCK_H
#include "qgsettingsinterface.h"
#include <gmock/gmock.h>

#include <QVariant>
#include <QString>

class QGSettingsMock : public QGSettingsInterface
{
public:
    virtual ~QGSettingsMock() {}

    MOCK_METHOD0(type, Type(void));
    MOCK_METHOD0(gsettings, QGSettings *(void));
    MOCK_CONST_METHOD1(get, QVariant(const QString &key));
    MOCK_METHOD2(set, void(const QString &key, const QVariant &value));
    MOCK_METHOD2(trySet, bool (const QString &key, const QVariant &value));
    MOCK_CONST_METHOD0(keys, QStringList(void));
    MOCK_CONST_METHOD1(choices, QVariantList(const QString &key));
    MOCK_METHOD1(reset, void(const QString &key));

    static bool isSchemaInstalled(const QByteArray &schema_id) {Q_UNUSED(schema_id); return true;}
};
#endif // QGSETTINGSMOCK_H
