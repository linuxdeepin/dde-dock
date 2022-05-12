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
