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
    void clicked(uint8_t, int, int);

private:
    QScopedPointer<IndicatorTrayWidgetPrivate> d_ptr;
    Q_DECLARE_PRIVATE_D(qGetPtrHelper(d_ptr), IndicatorTrayWidget)
};

