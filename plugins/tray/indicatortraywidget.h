// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#pragma once

#include <QScopedPointer>
#include <QLabel>

#include "abstracttraywidget.h"

class QGSettings;

class IndicatorTrayWidget: public AbstractTrayWidget
{
    Q_OBJECT
public:
    explicit IndicatorTrayWidget(const QString &indicatorName, QWidget *parent = Q_NULLPTR, Qt::WindowFlags f = Qt::WindowFlags());
    ~IndicatorTrayWidget();

    QString itemKeyForConfig() override;
    void updateIcon() override;
    void sendClick(uint8_t, int, int) override;
    void enableLabel(bool enable);
    static QString toIndicatorKey(const QString &indicatorName) { return QString("indicator:%1").arg(indicatorName); }
    static bool isIndicatorKey(const QString &itemKey) { return itemKey.startsWith("indicator:"); }

protected:
    void resizeEvent(QResizeEvent *event) override;

public Q_SLOTS:
    Q_SCRIPTABLE void setPixmapData(const QByteArray &data);
    Q_SCRIPTABLE void setText(const QString &text);

private slots:
    void onGSettingsChanged(const QString &key);

Q_SIGNALS:
    void clicked(uint8_t, int, int);

private:
    QLabel *m_label;

    QString m_indicatorName;
    const QGSettings *m_gsettings;
    bool m_enableClick;              // 置灰时设置为false，不触发click信号
};

