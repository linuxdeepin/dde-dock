#ifndef VOLUMNMODEL_H
#define VOLUMNMODEL_H

#include <QObject>

#include <com_deepin_daemon_audio.h>
#include <com_deepin_daemon_audio_sink.h>

using DBusAudio = com::deepin::daemon::Audio;
using DBusSink = com::deepin::daemon::audio::Sink;

class QDBusMessage;
class AudioSink;
class AudioPorts;

class VolumeModel : public QObject
{
    Q_OBJECT

public:
    explicit VolumeModel(QObject *parent);
    ~VolumeModel();

    void setActivePort(AudioPorts *port);

    QList<AudioSink *> sinks() const;
    QList<AudioPorts *> ports() const;

    AudioSink *defaultSink() const;

    void setVolume(int volume);
    void setMute(bool value);

    int volume();
    bool isMute();
    bool existActiveOutputDevice();

Q_SIGNALS:
    void defaultSinkChanged(AudioSink *);
    void volumeChanged(int);
    void muteChanged(bool);
    void checkPortChanged();

private Q_SLOTS:
    void onDefaultSinkChanged(const QDBusObjectPath & value);

private:
    void reloadSinks();
    void reloadPorts();
    void clearSinks();
    void clearPorts();

private:
    QList<AudioSink *> m_sinks;
    QList<AudioPorts *> m_ports;

    DBusAudio *m_audio;
};

class AudioSink : public QObject
{
    Q_OBJECT

    friend class VolumeModel;

Q_SIGNALS:
    void volumeChanged(int);
    void muteChanged(bool);

public:
    bool isDefault();
    bool isHeadPhone();

    void setBalance(double value, bool isPlay);
    void setFade(double value);
    void setMute(bool mute);
    void setPort(QString name);
    void setVolume(int value, bool isPlay);

    bool isMute();
    bool supportBalance();
    bool suoportFade();
    double balance();
    double baseVolume();
    double fade();
    int volume();
    QString description();
    QString name();
    uint cardId();

protected:
    explicit AudioSink(QString path, bool isDefault, QObject *parent = nullptr);
    ~AudioSink();
    void setDefault(bool isDefaultSink);

private:
    QString m_devicePath;
    DBusSink *m_sink;
    bool m_isDefault;
};

class AudioPorts : public QObject
{
    Q_OBJECT

    friend class VolumeModel;

public:
    uint cardId() const;
    QString cardName() const;
    QString name() const;
    QString description() const;
    int direction() const;
    bool isChecked() const;
    bool isHeadPhone() const;

protected:
    AudioPorts(uint cardId, QString cardName);
    ~AudioPorts();
    void setName(const QString &name);
    void setDescription(const QString &desc);
    void setDirection(int dir);
    void setIsChecked(bool isChecked);

private:
    uint m_cardId;
    QString m_cardName;
    QString m_portName;
    QString m_description;
    int m_direction;
    bool m_isCheck;
    bool m_isHeadPhone;
};

#endif // VOLUMNMODEL_H
