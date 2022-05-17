#include "volumemodel.h"

#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QVariantMap>
#include <QDebug>

/**
 * @brief 声音控制的类
 * @param parent
 */

static const QString serviceName = QString("com.deepin.daemon.Audio");
static const QString servicePath = QString("/com/deepin/daemon/Audio");

VolumeModel::VolumeModel(QObject *parent)
    : QObject(parent)
    , m_audio(new DBusAudio(serviceName, servicePath, QDBusConnection::sessionBus(), this))
{
    reloadSinks();
    reloadPorts();
    connect(m_audio, &DBusAudio::DefaultSinkChanged, this, &VolumeModel::onDefaultSinkChanged);
}

VolumeModel::~VolumeModel()
{
    clearPorts();
    clearSinks();
}

void VolumeModel::setActivePort(AudioPorts *port)
{
    m_audio->SetPort(port->cardId(), port->name(), port->direction());
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

void VolumeModel::reloadSinks()
{
    clearSinks();
    const QString defaultSinkPath = m_audio->defaultSink().path();
    QList<QDBusObjectPath> sinkPaths = m_audio->sinks();
    for (const QDBusObjectPath &sinkPath : sinkPaths) {
        const QString path = sinkPath.path();
        AudioSink *sink = new AudioSink(path, (path == defaultSinkPath), this);
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
    QString cards = m_audio->cardsWithoutUnavailable();
    QJsonParseError error;
    QJsonDocument json = QJsonDocument::fromJson(cards.toLocal8Bit(), &error);
    if (error.error != QJsonParseError::NoError)
        return;

    uint sinkCardId = 0;
    QString sinkCardName;
    AudioSink *sink = defaultSink();
    if (sink) {
        sinkCardId = sink->cardId();
        sinkCardName = sink->name();
    }

    QJsonArray array = json.array();
    for (const QJsonValue value : array) {
        QJsonObject cardObject = value.toObject();
        uint cardId = static_cast<uint>(cardObject.value("Id").toInt());
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

void VolumeModel::onDefaultSinkChanged(const QDBusObjectPath &value)
{
    AudioSink *audioSink = nullptr;
    const QString defaultPath = value.path();
    for (AudioSink *sink : m_sinks) {
        sink->setDefault(defaultPath == sink->m_devicePath);
        if (sink->isDefault())
            audioSink = sink;
    }

    if (!audioSink)
        return;

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

    Q_EMIT defaultSinkChanged(audioSink);
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

/**
 * @brief 具体的声音设备
 * @param parent
 */

AudioSink::AudioSink(QString path, bool isDefault, QObject *parent)
    : QObject(parent)
    , m_devicePath(path)
    , m_sink(new DBusSink(serviceName, path, QDBusConnection::sessionBus(), this))
    , m_isDefault(isDefault)
{
    connect(m_sink, &DBusSink::MuteChanged, this, &AudioSink::muteChanged);
    connect(m_sink, &DBusSink::VolumeChanged, this, [ this ](double value) {
        Q_EMIT this->volumeChanged(static_cast<int>(value * 100));
    });
}

AudioSink::~AudioSink()
{
}

void AudioSink::setDefault(bool isDefaultSink)
{
    m_isDefault = isDefaultSink;
}

bool AudioSink::isDefault()
{
    return m_isDefault;
}

bool AudioSink::isHeadPhone()
{
    return false;
}

void AudioSink::setBalance(double value, bool isPlay)
{
    m_sink->SetBalance(value, isPlay);
}

void AudioSink::setFade(double value)
{
    m_sink->SetFade(value);
}

void AudioSink::setMute(bool mute)
{
    m_sink->SetMute(mute);
}

void AudioSink::setPort(QString name)
{
    m_sink->SetPort(name);
}

void AudioSink::setVolume(int value, bool isPlay)
{
    m_sink->SetVolume((value * 0.01), isPlay);
}

bool AudioSink::isMute()
{
    return m_sink->mute();
}

bool AudioSink::supportBalance()
{
    return m_sink->supportBalance();
}

bool AudioSink::suoportFade()
{
    return m_sink->supportFade();
}

double AudioSink::balance()
{
    return m_sink->balance();
}

double AudioSink::baseVolume()
{
    return m_sink->baseVolume();
}

double AudioSink::fade()
{
    return m_sink->fade();
}

int AudioSink::volume()
{
    return static_cast<int>(m_sink->volume() * 100);
}

QString AudioSink::description()
{
    return m_sink->activePort().description;
}

QString AudioSink::name()
{
    return m_sink->activePort().name;
}

uint AudioSink::cardId()
{
    return m_sink->card();
}

AudioPorts::AudioPorts(uint cardId, QString cardName)
    : m_cardId(cardId)
    , m_cardName(cardName)
    , m_direction(0)
    , m_isCheck(false)
    , m_isHeadPhone(false)
{
}

AudioPorts::~AudioPorts()
{
}

uint AudioPorts::cardId() const
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
