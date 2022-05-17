/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#pragma once

#include "indicatortrayitem.h"

#include <QObject>
#include <QScopedPointer>

class IndicatorPluginPrivate;
class IndicatorPlugin : public QObject
{
    Q_OBJECT
public:
    explicit IndicatorPlugin(const QString &indicatorName, QObject *parent = nullptr);
    ~IndicatorPlugin();

    IndicatorTrayItem *widget();

    void removeWidget();

signals:
    void delayLoaded();
    void removed();

private slots:
    void textPropertyChanged(const QDBusMessage &message);
    void iconPropertyChanged(const QDBusMessage &message);

private:
    QScopedPointer<IndicatorPluginPrivate> d_ptr;
    Q_DECLARE_PRIVATE_D(qGetPtrHelper(d_ptr), IndicatorPlugin)
};
