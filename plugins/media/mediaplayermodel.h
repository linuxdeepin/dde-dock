// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MEDIAPLAYERMODEL_H
#define MEDIAPLAYERMODEL_H

#include <QObject>
#include <QDBusAbstractInterface>
#include <QDBusPendingReply>

typedef QMap<QString, QVariant> Dict;
Q_DECLARE_METATYPE(Dict)

class QDBusMessage;
class QDBusConnection;
class MediaPlayerInterface;

class MediaPlayerModel : public QObject
{
    Q_OBJECT

public:
    enum PlayStatus {
        Unknow = 0,
        Play,
        Pause,
        Stop
    };

public:
    explicit MediaPlayerModel(QObject *parent = nullptr);
    ~MediaPlayerModel();

    bool isActived();
    bool canGoNext();
    bool canGoPrevious();
    bool canPause();

    PlayStatus status();
    const QString name();
    const QString iconUrl();
    const QString album();
    const QString artist();

    void setStatus(const PlayStatus &stat);
    void playNext();

Q_SIGNALS:
    void startStop(bool);
    void statusChanged(const PlayStatus &);
    void metadataChanged();

private:
    void initMediaPlayer();
    PlayStatus convertStatus(const QString &stat);

private:
    bool m_isActived;
    QString m_serviceName;
    QString m_name;
    QString m_icon;
    QString m_album;
    QString m_artist;
    MediaPlayerInterface *m_mediaInter;
};

class MediaPlayerInterface : public QDBusAbstractInterface
{
    Q_OBJECT

public:
    MediaPlayerInterface(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent = nullptr);
    ~MediaPlayerInterface();

public:
    inline QDBusPendingReply<> Play() {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("Play"), argumentList);
    }

    inline QDBusPendingReply<> Stop() {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("Stop"), argumentList);
    }

    inline QDBusPendingReply<> Pause() {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("Pause"), argumentList);
    }

    inline QDBusPendingReply<> Next() {
        QList<QVariant> argumentList;
        return asyncCallWithArgumentList(QStringLiteral("Next"), argumentList);
    }

    Q_PROPERTY(Dict Metadata READ metadata NOTIFY MetadataChanged)
    inline Dict metadata() const
    { return qvariant_cast<Dict>(property("Metadata")); }

    Q_PROPERTY(bool CanGoNext READ canGoNext NOTIFY CanGoNextChanged)
    inline bool canGoNext() const
    { return qvariant_cast< bool >(property("CanGoNext")); }

    Q_PROPERTY(bool CanGoPrevious READ canGoPrevious NOTIFY CanGoPreviousChanged)
    inline bool canGoPrevious() const
    { return qvariant_cast< bool >(property("CanGoPrevious")); }

    Q_PROPERTY(bool CanPause READ canPause NOTIFY CanPauseChanged)
    inline bool canPause() const
    { return qvariant_cast< bool >(property("CanPause")); }

    Q_PROPERTY(QString PlaybackStatus READ playbackStatus NOTIFY PlaybackStatusChanged)
    inline QString playbackStatus() const
    { return qvariant_cast< QString >(property("PlaybackStatus")); }

Q_SIGNALS:
    void MetadataChanged();
    void CanGoNextChanged();
    void CanGoPreviousChanged();
    void CanPauseChanged();
    void PlaybackStatusChanged();

private Q_SLOTS:
    void onPropertyChanged(const QDBusMessage &msg);
};

#endif // MEDIAPLAYERLISTENER_H
