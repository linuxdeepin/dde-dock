// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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

    bool isLoaded();

signals:
    void delayLoaded();
    void removed();

private slots:
    void textPropertyChanged(const QDBusMessage &message);
    void iconPropertyChanged(const QDBusMessage &message);

private:
    QScopedPointer<IndicatorPluginPrivate> d_ptr;
    bool m_isLoaded;
    Q_DECLARE_PRIVATE_D(qGetPtrHelper(d_ptr), IndicatorPlugin)
};
