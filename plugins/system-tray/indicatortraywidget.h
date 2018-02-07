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

#pragma once

#include <QScopedPointer>
#include "abstracttraywidget.h"

class IndicatorTrayWidgetPrivate;
class IndicatorTrayWidget: public AbstractTrayWidget
{
    Q_OBJECT
public:
    explicit IndicatorTrayWidget(const QString &itemKey, QWidget *parent = Q_NULLPTR, Qt::WindowFlags f = Qt::WindowFlags());
    ~IndicatorTrayWidget();

    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;
    void sendClick(uint8_t, int, int) Q_DECL_OVERRIDE;

    QSize sizeHint() const Q_DECL_OVERRIDE;

    static QString toTrayWidgetId(const QString &indicatorKey) { return QString("indicator:%1").arg(indicatorKey); }
    static QString toIndicatorId(QString itemKey) { return itemKey.remove("indicator:"); }
    static bool isIndicatorKey(const QString &itemKey) { return itemKey.startsWith("indicator:"); }

public Q_SLOTS:
    Q_SCRIPTABLE void setPixmapData(const QByteArray &data);
    Q_SCRIPTABLE void setPixmapPath(const QString &text);
    Q_SCRIPTABLE void setText(const QString &text);

public Q_SLOTS:
    void iconPropertyChanged(const QDBusMessage &msg);
    void textPropertyChanged(const QDBusMessage &msg);

Q_SIGNALS:
    void delayLoaded();
    void clicked(uint8_t, int, int);

private:
    QScopedPointer<IndicatorTrayWidgetPrivate> d_ptr;
    Q_DECLARE_PRIVATE_D(qGetPtrHelper(d_ptr), IndicatorTrayWidget)
};

