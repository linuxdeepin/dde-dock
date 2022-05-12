#include "volumemodel.h"

#include <QDBusArgument>
#include <QDBusConnection>
#include <QDBusConnectionInterface>
#include <QDBusInterface>
#include <QDBusMessage>
#include <QDBusObjectPath>
#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QVariantMap>
#include <QDebug>
#include <QDBusMetaType>

/**
 * @brief 声音控制的类
 * @param parent
 */

static const QString serviceName = QString("com.deepin.daemon.Audio");
static const QString servicePath = QString("/com/deepin/daemon/Audio");
static const QString interfaceName = QString("com.deepin.daemon.Audio");
static const QString propertiesInterface = QString("org.freedesktop.DBus.Properties");

VolumeModel::VolumeModel(QObject *parent)
    : QObject(parent)
{
    initService();
}

VolumeModel::~VolumeModel()
{
    clearPorts();
    clearSinks();
}

void VolumeModel::setActivePort(AudioPorts *port)
{
    callMethod("SetPort", { port->cardId(), port->name(), port->direction() });
}

QList<AudioSink *> VolumeModel::sinks() const
{
    return m_sinks;
}

QList<AudioPorts *> VolumeModel::ports() const
{
    return m_ports;
}

AudioSink *VolumeModel::defaultSink() const
{
    for (AudioSink *sink : m_sinks) {
        if (sink->isDefault())
            return sink;
    }

    return nullptr;
}

void VolumeModel::setVolume(int volumn)
{
    for (AudioSink *audiosink : m_sinks) {
        if (audiosink->isDefault()) {
            audiosink->setVolume(volumn, true);
            break;
        }
    }
}

void VolumeModel::setMute(bool value)
{
    for (AudioSink *audiosink : m_sinks) {
        if (audiosink->isDefault()) {
            audiosink->setMute(value);
            break;
        }
    }
}

int VolumeModel::volume()
{
    for (AudioSink *audiosink : m_sinks) {
        if (audiosink->isDefault())
            return audiosink->volume();
    }

    return 0;
}

bool VolumeModel::isMute()
{
    for (AudioSink *audiosink : m_sinks) {
        if (audiosink->isDefault())
            return audiosink->isMute();
    }

    return false;
}

bool VolumeModel::existActiveOutputDevice()
{
    for (AudioPorts *port : m_ports) {
        if (port->direction() == 1)
            return true;
    }

    return false;
}

void VolumeModel::onPropertyChanged(const QDBusMessage &msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QString interName = msg.arguments().at(0).toString();
    if (interName != interfaceName)
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    if (changedProps.contains("DefaultSink")) {
        QVariant defaultSink = changedProps.value("DefaultSink");
        QString defaultSinkPath = defaultSink.value<QDBusObjectPath>().path();
        for (AudioSink *audioSink : m_sinks) {
            if (audioSink->m_devicePath == defaultSinkPath) {
                updateDefaultSink(audioSink);
                Q_EMIT defaultSinkChanged(audioSink);
                break;
            }
        }
    }
}

void VolumeModel::initService()
{
    QDBusConnection::sessionBus().connect(serviceName, servicePath, propertiesInterface,
       "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(const QDBusMessage &)));

    reloadSinks();
    reloadPorts();

    QDBusConnectionInterface *dbusInterface = QDBusConnection::sessionBus().interface();
    connect(dbusInterface, &QDBusConnectionInterface::serviceOwnerChanged, this,
            [ = ](const QString &name, const QString &, const QString &newOwner) {
        if (name == serviceName) {
            if (newOwner.isEmpty()) {
                QDBusConnection::sessionBus().disconnect(serviceName, servicePath, propertiesInterface,
                    "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(const QDBusMessage &)));
                clearSinks();
            } else {
                QDBusConnection::sessionBus().connect(serviceName, servicePath, propertiesInterface,
                   "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(const QDBusMessage &)));
                reloadSinks();
                reloadPorts();
            }
        }
    });
}

void VolumeModel::reloadSinks()
{
    clearSinks();
    QList<QDBusObjectPath> sinkPaths = properties<QList<QDBusObjectPath>>("Sinks");
    for (const QDBusObjectPath &sinkPath : sinkPaths) {
        AudioSink *sink = new AudioSink(sinkPath.path(), this);
        connect(sink, &AudioSink::volumeChanged, this, [ = ](int volume) {
            if (sink->isDefault())
                Q_EMIT volumeChanged(volume);
        });
        connect(sink, &AudioSink::muteChanged, this, [ = ](bool isMute) {
            if (sink->isDefault())
                Q_EMIT muteChanged(isMute);
        });

        m_sinks << sink;
    }
}

void VolumeModel::reloadPorts()
{
    clearPorts();
    QString cards = properties<QString>("CardsWithoutUnavailable");
    QJsonParseError error;
    QJsonDocument json = QJsonDocument::fromJson(cards.toLocal8Bit(), &error);
    if (error.error != QJsonParseError::NoError)
        return;

    int sinkCardId = -1;
    QString sinkCardName;
    AudioSink *sink = defaultSink();
    if (sink) {
        sinkCardId = sink->cardId();
        sinkCardName = sink->name();
    }

    QJsonArray array = json.array();
    for (const QJsonValue value : array) {
        QJsonObject cardObject = value.toObject();
        int cardId = cardObject.value("Id").toInt();
        QString cardName = cardObject.value("Name").toString();
        QJsonArray jPorts = cardObject.value("Ports").toArray();
        for (const QJsonValue jPortValue : jPorts) {
            QJsonObject jPort = jPortValue.toObject();
            if (!jPort.value("Enabled").toBool())
                continue;

            int direction = jPort.value("Direction").toInt();
            if (direction != 1)
                continue;

            AudioPorts *port = new AudioPorts(cardId, cardName);
            port->setName(jPort.value("Name").toString());
            port->setDescription(jPort.value("Description").toString());
            port->setDirection(direction);
            if (port->cardId() == sinkCardId && port->name() == sinkCardName)
                port->setIsChecked(true);

            m_ports << port;
        }
    }
}

void VolumeModel::updateDefaultSink(AudioSink *audioSink)
{
    bool checkChanged = false;
    for (AudioPorts *port : m_ports) {
        bool oldChecked = port->isChecked();
        port->setIsChecked(port->cardId() == audioSink->cardId()
                            && port->name() == audioSink->name());

        if (oldChecked != port->isChecked() && port->isChecked())
            checkChanged = true;
    }

    if (checkChanged)
        Q_EMIT checkPortChanged();
}

void VolumeModel::clearSinks()
{
    for (AudioSink *sink : m_sinks)
        delete sink;

    m_sinks.clear();
}

void VolumeModel::clearPorts()
{
    for (AudioPorts *port : m_ports)
        delete port;

    m_ports.clear();
}

QDBusMessage VolumeModel::callMethod(const QString &methodName, const QList<QVariant> &argument)
{
    QDBusInterface dbusInter(serviceName, servicePath, interfaceName, QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        QDBusPendingCall reply = dbusInter.asyncCallWithArgumentList(methodName, argument);
        reply.waitForFinished();
        return reply.reply();
    }
    return QDBusMessage();
}

template<typename T>
T VolumeModel::properties(const QString &propName)
{
    QDBusInterface dbusInter(serviceName, servicePath, interfaceName, QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        QByteArray ba = propName.toLatin1();
        const char *prop = ba.data();
        return dbusInter.property(prop).value<T>();
    }

    return T();
}

/**
 * @brief 具体的声音设备
 * @param parent
 */

AudioSink::AudioSink(QString path, QObject *parent)
    : QObject(parent)
    , m_devicePath(path)
{
    QDBusConnection::sessionBus().connect(serviceName, path, propertiesInterface,
        "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(const QDBusMessage&)));
}

AudioSink::~AudioSink()
{
}

bool AudioSink::isDefault()
{
    QDBusInterface dbusInter(serviceName, servicePath, interfaceName, QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        QString defaultSink = dbusInter.property("DefaultSink").value<QDBusObjectPath>().path();
        return defaultSink == m_devicePath;
    }

    return false;
}

bool AudioSink::isHeadPhone()
{
    return false;
}

void AudioSink::setBalance(double value, bool isPlay)
{
    callMethod("SetBalance", { value, isPlay });
}

void AudioSink::setFade(double value)
{
    callMethod("SetFade", { value });
}

void AudioSink::setMute(bool mute)
{
    callMethod("SetMute", { mute });
}

void AudioSink::setPort(QString name)
{
    callMethod("SetPort", { name });
}

void AudioSink::setVolume(double value, bool isPlay)
{
    callMethod("SetVolume", { value * 0.01, isPlay });
}

bool AudioSink::isMute()
{
    return getProperties<bool>("Mute");
}

bool AudioSink::supportBalance()
{
    return getProperties<bool>("SupportBalance");
}

bool AudioSink::suoportFade()
{
    return getProperties<bool>("SupportFade");
}

double AudioSink::balance()
{
    return getProperties<double>("Balance");
}

double AudioSink::baseVolume()
{
    return getProperties<double>("BaseVolume");
}

double AudioSink::fade()
{
    return getProperties<double>("Fade");
}

int AudioSink::volume()
{
    return static_cast<int>(getProperties<double>("Volume") * 100);
}

QString AudioSink::description()
{
    QVariantList value = getPropertiesByFreeDesktop("ActivePort");
    if (value.size() >= 2)
        return value[1].toString();

    return getProperties<QString>("Description");
}

QString AudioSink::name()
{
    QVariantList value = getPropertiesByFreeDesktop("ActivePort");
    if (value.size() >= 2)
        return value[0].toString();

    return getProperties<QString>("Name");
}

int AudioSink::cardId()
{
    return getProperties<int>("Card");
}

void AudioSink::onPropertyChanged(const QDBusMessage &msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QString interfaceName = msg.arguments().at(0).toString();
    if (interfaceName != propertiesInterface)
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    if (changedProps.contains("Volume"))
        Q_EMIT volumeChanged(static_cast<int>(changedProps.value("Volume").toDouble() * 100));

    if (changedProps.contains("Mute"))
        Q_EMIT muteChanged(changedProps.value("Mute").toBool());
}

template<typename T>
T AudioSink::getProperties(const QString &propName)
{
    QDBusInterface dbusInter(serviceName, m_devicePath, interfaceName + QString(".Sink"), QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        QByteArray ba = propName.toLatin1();
        const char *prop = ba.data();
        return dbusInter.property(prop).value<T>();
    }

    return T();
}

QDBusMessage AudioSink::callMethod(const QString &methodName, const QList<QVariant> &argument)
{
    QDBusInterface dbusInter(serviceName, m_devicePath, interfaceName + QString(".Sink"), QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        QDBusPendingCall reply = dbusInter.asyncCallWithArgumentList(methodName, argument);
        reply.waitForFinished();
        return reply.reply();
    }

    return QDBusMessage();
}

static QVariantList argToString(const QDBusArgument &busArg)
{
    QVariantList out;
    QString busSig = busArg.currentSignature();
    bool doIterate = false;
    QDBusArgument::ElementType elementType = busArg.currentType();

    switch (elementType) {
    case QDBusArgument::BasicType:
    case QDBusArgument::VariantType: {
        QVariant value = busArg.asVariant();
        switch (value.type()) {
        case QVariant::Bool:
            out << value.toBool();
            break;
        case QVariant::Int:
            out << value.toInt();
            break;
        case QVariant::String:
            out << value.toString();
            break;
        case QVariant::UInt:
            out << value.toUInt();
            break;
        case QVariant::ULongLong:
            out << value.toULongLong();
            break;
        case QMetaType::UChar:
            out << value.toUInt();
            break;
        default:
            out << QVariant();
            break;
        }
        out += busArg.asVariant().toString();
        break;
    }
    case QDBusArgument::StructureType:
        busArg.beginStructure();
        doIterate = true;
        break;
    case QDBusArgument::ArrayType:
        busArg.beginArray();
        doIterate = true;
        break;
    case QDBusArgument::MapType:
        busArg.beginMap();
        doIterate = true;
        break;
    case QDBusArgument::UnknownType:
    default:
        out << QVariant();
        return out;
    }

    if (doIterate && !busArg.atEnd()) {
        while (!busArg.atEnd()) {
            out << argToString(busArg);
            if (out.isEmpty())
                break;
        }
    }

    return out;
}

QVariantList AudioSink::getPropertiesByFreeDesktop(const QString &propName)
{
    QDBusInterface dbusInter(serviceName, m_devicePath, "org.freedesktop.DBus.Properties", QDBusConnection::sessionBus());
    if (dbusInter.isValid()) {
        QDBusPendingCall reply = dbusInter.asyncCallWithArgumentList("Get", { interfaceName + ".Sink", propName });
        reply.waitForFinished();
        QVariantList lists = reply.reply().arguments();
        for (QVariantList::ConstIterator it = lists.begin(); it != lists.end(); ++it) {
            QVariant arg = (*it);
            const QVariant v = qvariant_cast<QDBusVariant>(arg).variant();
            return argToString(v.value<QDBusArgument>());
        }
    }

    return QVariantList();
}

AudioPorts::AudioPorts(int cardId, QString cardName)
    : m_cardId(cardId)
    , m_cardName(cardName)
    , m_isCheck(false)
    , m_isHeadPhone(false)
{
}

AudioPorts::~AudioPorts()
{
}

int AudioPorts::cardId() const
{
    return m_cardId;
}

QString AudioPorts::cardName() const
{
    return m_cardName;
}

void AudioPorts::setName(const QString &name)
{
    m_portName = name;
}

QString AudioPorts::name() const
{
    return m_portName;
}

void AudioPorts::setDescription(const QString &desc)
{
    m_description = desc;
}

QString AudioPorts::description() const
{
    return m_description;
}

void AudioPorts::setDirection(int dir)
{
    m_direction = dir;
}

void AudioPorts::setIsChecked(bool isChecked)
{
    m_isCheck = isChecked;
}

int AudioPorts::direction() const
{
    return m_direction;
}

bool AudioPorts::isChecked() const
{
    return m_isCheck;
}

bool AudioPorts::isHeadPhone() const
{
    return m_isHeadPhone;
}
