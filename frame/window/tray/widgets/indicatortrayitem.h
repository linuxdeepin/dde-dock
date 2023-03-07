// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#pragma once

#include <QScopedPointer>
#include <QLabel>
#include <QPaintEvent>

#include "basetraywidget.h"

class QGSettings;

class IndicatorTrayItem: public BaseTrayWidget
{
    Q_OBJECT

public:
    explicit IndicatorTrayItem(const QString &indicatorName, QWidget *parent = Q_NULLPTR, Qt::WindowFlags f = Qt::WindowFlags());
    ~IndicatorTrayItem() override;

    QString itemKeyForConfig() override;
    void updateIcon() override;
    void sendClick(uint8_t, int, int) override;
    void enableLabel(bool enable);
    static QString toIndicatorKey(const QString &indicatorName) { return QString("indicator:%1").arg(indicatorName); }
    static bool isIndicatorKey(const QString &itemKey) { return itemKey.startsWith("indicator:"); }
    QPixmap icon() override;
    const QByteArray &pixmapData() const;
    const QString text() const;
    bool containsPoint(const QPoint &mouse) override { return false; }

private:
    void paintEvent(QPaintEvent *) override;

public Q_SLOTS:
    Q_SCRIPTABLE void setPixmapData(const QByteArray &data);
    Q_SCRIPTABLE void setText(const QString &text);

Q_SIGNALS:
    void clicked(uint8_t, int, int);
    void textChanged(const QString &text);

private:
    QString m_indicatorName;
    bool m_enableClick;              // 置灰时设置为false，不触发click信号
    QByteArray m_pixmapData;
    QString m_text;
};

