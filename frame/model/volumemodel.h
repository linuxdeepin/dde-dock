#ifndef VOLUMNMODEL_H
#define VOLUMNMODEL_H

#include <QObject>

class QDBusMessage;
class AudioSink;
class AudioPorts;

class VolumeModel : public QObject
{
    Q_OBJECT

Q_SIGNALS:
    void defaultSinkChanged(AudioSink *);
    void volumeChanged(int);
    void muteChanged(bool);
    void checkPortChanged();

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

private Q_SLOTS:
    void onPropertyChanged(const QDBusMessage& msg);

private:
    void initService();
    void reloadSinks();
    void reloadPorts();
    void clearSinks();
    void clearPorts();
    void updateDefaultSink(AudioSink *audioSink);

private:
    QDBusMessage callMethod(const QString &methodName, const QList<QVariant> &argument);
    template<typename T>
    T properties(const QString &propName);

private:
    QList<AudioSink *> m_sinks;
    QList<AudioPorts *> m_ports;
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
    void setVolume(double value, bool isPlay);

    bool isMute();
    bool supportBalance();
    bool suoportFade();
    double balance();
    double baseVolume();
    double fade();
    int volume();
    QString description();
    QString name();
    int cardId();

private Q_SLOTS:
    void onPropertyChanged(const QDBusMessage& msg);

protected:
    explicit AudioSink(QString path, QObject *parent = nullptr);
    ~AudioSink();

private:
    QDBusMessage callMethod(const QString &methodName, const QList<QVariant> &argument);
    template<typename T>
    T getProperties(const QString &propName);

    QList<QVariant> getPropertiesByFreeDesktop(const QString &propName);

private:
    QString m_devicePath;
};

class AudioPorts : public QObject
{
    Q_OBJECT

    friend class VolumeModel;

public:
    int cardId() const;
    QString cardName() const;
    QString name() const;
    QString description() const;
    int direction() const;
    bool isChecked() const;
    bool isHeadPhone() const;

protected:
    AudioPorts(int cardId, QString cardName);
    ~AudioPorts();
    void setName(const QString &name);
    void setDescription(const QString &desc);
    void setDirection(int dir);
    void setIsChecked(bool isChecked);

private:
    int m_cardId;
    QString m_cardName;
    QString m_portName;
    QString m_description;
    int m_direction;
    bool m_isCheck;
    bool m_isHeadPhone;
};

#endif // VOLUMNMODEL_H
