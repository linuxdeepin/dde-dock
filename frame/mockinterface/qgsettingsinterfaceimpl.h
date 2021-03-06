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
#ifndef QGSETTINGSINTERFACEIMPL_H
#define QGSETTINGSINTERFACEIMPL_H
#include <QObject>

#include "qgsettingsinterface.h"

class QGSettings;
class QGSettingsInterfaceImpl : public QGSettingsInterface
{
public:
    QGSettingsInterfaceImpl(const QByteArray &schema_id, const QByteArray &path = QByteArray(), QObject *parent = nullptr);
    ~QGSettingsInterfaceImpl() override;

    virtual Type type() override;
    virtual QGSettings *gsettings() override;
    virtual QVariant get(const QString &key) const override;
    virtual void set(const QString &key, const QVariant &value) override;
    virtual bool trySet(const QString &key, const QVariant &value) override;
    virtual QStringList keys() const override;
    virtual QVariantList choices(const QString &key) const override;
    virtual void reset(const QString &key) override;
    static bool isSchemaInstalled(const QByteArray &schema_id);

private:
    QGSettings *m_gsettings;
};

#endif // QGSETTINGSINTERFACEIMPL_H
