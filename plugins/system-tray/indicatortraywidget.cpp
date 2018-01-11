#include "indicatortraywidget.h"

#include <QDebug>
#include <QLabel>
#include <QBoxLayout>

#include <QFile>
#include <QTimer>
#include <QJsonDocument>
#include <QJsonObject>

#include <QVariantMap>
#include <QDBusConnection>
#include <QDBusInterface>
#include <QDBusArgument>
#include <QDBusReply>
#include <QDBusPendingCall>
#include <QMetaProperty>

class IndicatorTrayWidgetPrivate
{
public:
    IndicatorTrayWidgetPrivate(IndicatorTrayWidget *parent) : q_ptr(parent) {}

    void updateContent();

    void initDBus(const QString &indicatorKey);

    template<typename Func>
    void featData(const QString &key,
                  const QJsonObject &data,
                  const char *propertyChangedSlot,
                  Func const &callback)
    {
        Q_Q(IndicatorTrayWidget);
        auto dataConfig = data.value(key).toObject();
        auto dbusService = dataConfig.value("dbus_service").toString();
        auto dbusPath = dataConfig.value("dbus_path").toString();
        auto dbusInterface = dataConfig.value("dbus_interface").toString();
        auto isSystemBus = dataConfig.value("system_dbus").toBool(false);
        auto bus = isSystemBus ? QDBusConnection::systemBus() : QDBusConnection::sessionBus();

        QDBusInterface interface(dbusService, dbusPath, dbusInterface, bus, q);

        if (dataConfig.contains("dbus_method")) {
            QString methodName = dataConfig.value("dbus_method").toString();
            auto ratio = q->devicePixelRatioF();
            QDBusReply<QByteArray> reply = interface.call(methodName.toStdString().c_str(), ratio);
            callback(reply.value());
        }

        if (dataConfig.contains("dbus_properties")) {
            auto propertyName = dataConfig.value("dbus_properties").toString().toStdString();
            propertyInterfaceNames.insert(key, dbusInterface);
            propertyNames.insert(key, QString::fromStdString(propertyName));
            callback(interface.property(propertyName.c_str()));
            QDBusConnection::sessionBus().connect(dbusService,
                                                  dbusPath,
                                                  "org.freedesktop.DBus.Properties",
                                                  "PropertiesChanged",
                                                  "sa{sv}as",
                                                  q,
                                                  propertyChangedSlot);
        }
    }

    template<typename Func>
    void propertyChanged(const QString &key, const QDBusMessage &msg, Func const &callback)
    {
        QList<QVariant> arguments = msg.arguments();
        if (3 != arguments.count()) {
            qWarning() << "arguments count must be 3";
            return;
        }
        QString interfaceName = msg.arguments().at(0).toString();
        if (interfaceName != propertyInterfaceNames.value(key)) {
            qWarning() << "interfaceName mismatch" << interfaceName << propertyInterfaceNames.value(key) << key;
            return;
        }
        QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
        if (changedProps.contains(propertyNames.value(key))) {
            callback(changedProps.value(propertyNames.value(key)));
        }
    }

    QLabel                  *label = Q_NULLPTR;
    QMap<QString, QString>  propertyNames;
    QMap<QString, QString>  propertyInterfaceNames;

    IndicatorTrayWidget *q_ptr;
    Q_DECLARE_PUBLIC(IndicatorTrayWidget)
};

IndicatorTrayWidget::IndicatorTrayWidget(const QString &indicatorKey, QWidget *parent, Qt::WindowFlags f) :
    AbstractTrayWidget(parent, f),
    d_ptr(new IndicatorTrayWidgetPrivate(this))
{
    Q_D(IndicatorTrayWidget);

    setAttribute(Qt::WA_TranslucentBackground);

    auto layout = new QVBoxLayout(this);
    layout->setContentsMargins(0, 0, 0, 0);
    d->label = new QLabel(this);

    QPalette p = palette();
    p.setColor(QPalette::Foreground, Qt::white);
    p.setColor(QPalette::Background, Qt::red);
    d->label->setPalette(p);

    layout->addWidget(d->label, 0, Qt::AlignCenter);
    setLayout(layout);

    // register dbus
    auto path = QString("/com/deepin/dde/Dock/Indicator/") + indicatorKey;
    auto interface =  QString("com.deepin.dde.Dock.Indicator.") + indicatorKey;
    auto sessionBus = QDBusConnection::sessionBus();
    sessionBus.registerObject(path,
                              interface,
                              this,
                              QDBusConnection::ExportScriptableSlots);

    d->initDBus(indicatorKey);
}

IndicatorTrayWidget::~IndicatorTrayWidget()
{
}

void IndicatorTrayWidget::setActive(const bool)
{

}

void IndicatorTrayWidget::updateIcon()
{

}

const QImage IndicatorTrayWidget::trayImage()
{
    return grab().toImage();
}

void IndicatorTrayWidget::sendClick(uint8_t buttonIndex, int x, int y)
{
    Q_EMIT clicked(buttonIndex, x, y);
}

QSize IndicatorTrayWidget::sizeHint() const
{
    auto sz = AbstractTrayWidget::sizeHint();
    sz.setHeight(26);
    sz.setWidth(26);
    return sz;
}

void IndicatorTrayWidget::setPixmapData(const QByteArray &data)
{
    Q_D(IndicatorTrayWidget);
    auto rawPixmap = QPixmap::fromImage(QImage::fromData(data));
    rawPixmap.setDevicePixelRatio(devicePixelRatioF());
    d->label->setPixmap(rawPixmap);
    d->updateContent();
}

void IndicatorTrayWidget::setPixmapPath(const QString &text)
{
    Q_D(IndicatorTrayWidget);
    d->label->setPixmap(QPixmap(text));
    d->updateContent();
}

void IndicatorTrayWidget::setText(const QString &text)
{
    Q_D(IndicatorTrayWidget);
    d->label->setText(text);
    d->updateContent();
}

void IndicatorTrayWidget::iconPropertyChanged(const QDBusMessage &msg)
{
    Q_D(IndicatorTrayWidget);
    d->propertyChanged("icon", msg, [ = ](QVariant v) {
        setPixmapData(v.toByteArray());
    });
}

void IndicatorTrayWidget::textPropertyChanged(const QDBusMessage &msg)
{
    Q_D(IndicatorTrayWidget);
    d->propertyChanged("text", msg, [ = ](QVariant v) {
        setText(v.toString());
    });
}

void IndicatorTrayWidgetPrivate::updateContent()
{
    Q_Q(IndicatorTrayWidget);
    q->update();
    Q_EMIT q->iconChanged();
}

void IndicatorTrayWidgetPrivate::initDBus(const QString &indicatorKey)
{
    Q_Q(IndicatorTrayWidget);

    QString filepath = QString("/etc/dde-dock/indicator/%1.json").arg(indicatorKey);
    QFile confFile(filepath);
    if (!confFile.open(QIODevice::ReadOnly)) {
        qCritical() << "read indicator config Error";
    }

    QJsonDocument doc = QJsonDocument::fromJson(confFile.readAll());
    confFile.close();
    auto config = doc.object();

    auto delay = config.value("delay").toInt(0);

    qDebug() << "delay load" << delay << indicatorKey << q;

    q->hide();
    QTimer::singleShot(delay, [ = ]() {
        auto data = config.value("data").toObject();

        if (data.contains("text")) {
            featData("text", data, SLOT(textPropertyChanged(QDBusMessage)), [ = ](QVariant v) {
                q->setText(v.toString());
            });
        }
        if (data.contains("icon")) {
            featData("icon", data, SLOT(iconPropertyChanged(QDBusMessage)), [ = ](QVariant v) {
                q->setPixmapData(v.toByteArray());
            });
        }

        auto action = config.value("action").toObject();
        q->connect(q, &IndicatorTrayWidget::clicked, q, [ = ](uint8_t /*button_index*/, int /*x*/, int /*y*/) {
            auto triggerConfig = action.value("trigger").toObject();
            auto dbusService = triggerConfig.value("dbus_service").toString();
            auto dbusPath = triggerConfig.value("dbus_path").toString();
            auto dbusInterface = triggerConfig.value("dbus_interface").toString();
            auto methodName = triggerConfig.value("dbus_method").toString();
            auto isSystemBus = triggerConfig.value("system_dbus").toBool(false);
            auto bus = isSystemBus ? QDBusConnection::systemBus() : QDBusConnection::sessionBus();

            QDBusInterface interface(dbusService, dbusPath, dbusInterface, bus, q);
            interface.asyncCall(methodName);
        });

        q->show();
        Q_EMIT q->delayLoaded();
    });
}
