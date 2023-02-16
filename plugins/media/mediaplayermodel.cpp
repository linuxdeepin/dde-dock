// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "mediaplayermodel.h"

#include <QDBusConnectionInterface>
#include <QDBusInterface>
#include <QDBusPendingCall>
#include <QDBusReply>
#include <QDebug>
#include <QJsonDocument>
#include <QJsonObject>
#include <QMetaMethod>
#include <QDBusAbstractInterface>

MediaPlayerModel::MediaPlayerModel(QObject *parent)
    : QObject(parent)
    , m_isActived(false)
    , m_mediaInter(nullptr)
{
    initMediaPlayer();
}

MediaPlayerModel::~MediaPlayerModel()
{
}

bool MediaPlayerModel::isActived()
{
    return m_isActived;
}

bool MediaPlayerModel::canGoNext()
{
    return m_mediaInter ? m_mediaInter->canGoNext() : false;
}

bool MediaPlayerModel::canGoPrevious()
{
    return m_mediaInter ? m_mediaInter->canGoPrevious() : false;
}

bool MediaPlayerModel::canPause()
{
    return m_mediaInter ? m_mediaInter->canPause() : false;
}

MediaPlayerModel::PlayStatus MediaPlayerModel::status()
{
    if (!m_isActived || !m_mediaInter)
        return PlayStatus::Stop;

    return convertStatus(m_mediaInter->playbackStatus());
}

const QString MediaPlayerModel::name()
{
    if (m_mediaInter) {
        Dict data = m_mediaInter->metadata();
        return data["xesam:title"].toString();
    }

    return QString();
}

const QString MediaPlayerModel::iconUrl()
{
    if (m_mediaInter) {
        Dict data = m_mediaInter->metadata();
        return data["mpris:artUrl"].toString();
    }

    return QString();
}

const QString MediaPlayerModel::album()
{
    if (m_mediaInter) {
        Dict data = m_mediaInter->metadata();
        return data["xesam:album"].toString();
    }

    return QString();
}

const QString MediaPlayerModel::artist()
{
    if (m_mediaInter) {
        Dict data = m_mediaInter->metadata();
        return data["xesam:artist"].toString();
    }

    return QString();
}

void MediaPlayerModel::setStatus(const MediaPlayerModel::PlayStatus &stat)
{
    if (!m_mediaInter)
        return;

    switch (stat) {
    case MediaPlayerModel::PlayStatus::Play: {
        m_mediaInter->Play();
        break;
    }
    case MediaPlayerModel::PlayStatus::Stop: {
        m_mediaInter->Stop();
        break;
    }
    case MediaPlayerModel::PlayStatus::Pause: {
        m_mediaInter->Pause();
        break;
    }
    default: break;
    }
}

void MediaPlayerModel::playNext()
{
    if (m_mediaInter)
        m_mediaInter->Next();
}

void MediaPlayerModel::initMediaPlayer()
{
    QDBusInterface dbusInter("org.freedesktop.DBus", "/", "org.freedesktop.DBus", QDBusConnection::sessionBus(), this);
    QDBusPendingCall call = dbusInter.asyncCall("ListNames");
    QDBusPendingCallWatcher *watcher = new QDBusPendingCallWatcher(call, this);
    connect(watcher, &QDBusPendingCallWatcher::finished, [ = ] {
        m_serviceName.clear();
        if (call.isError())
            return;

        QDBusReply<QStringList> reply = call.reply();
        const QStringList &serviceList = reply.value();

        for (const QString &serv : serviceList) {
            if (!serv.startsWith("org.mpris.MediaPlayer2"))
                continue;

            QDBusInterface serviceInterface(serv, "/org/mpris/MediaPlayer2",
                                            "org.mpris.MediaPlayer2.Player", QDBusConnection::sessionBus(), this);
            // 如果开启了谷歌浏览器的后台服务(org.mpris.MediaPlayer2.chromium.instance17352)
            // 也符合名称要求，但是它不是音乐服务，此时需要判断是否存在这个属性
            QVariant v = serviceInterface.property("CanPlay");
            if (!v.isValid() || !v.value<bool>())
                continue;

            m_serviceName = serv;
            break;
        }

        if (!m_serviceName.isEmpty()) {
            m_isActived = true;

            m_mediaInter = new MediaPlayerInterface(m_serviceName, "/org/mpris/MediaPlayer2", QDBusConnection::sessionBus(), this);
            connect(m_mediaInter, &MediaPlayerInterface::PlaybackStatusChanged, this, [ this ] {
                Q_EMIT statusChanged(convertStatus(m_mediaInter->playbackStatus()));
            });
            connect(m_mediaInter, &MediaPlayerInterface::MetadataChanged, this, &MediaPlayerModel::metadataChanged);
            Dict v = m_mediaInter->metadata();
            m_name = v.value("xesam:title").toString();
            m_icon = v.value("mpris:artUrl").toString();
            m_album = v.value("xesam:album").toString();
            m_artist = v.value("xesam:artist").toString();
            Q_EMIT startStop(m_isActived);
            return;
        }

        QDBusConnectionInterface *dbusInterface = QDBusConnection::sessionBus().interface();
        connect(dbusInterface, &QDBusConnectionInterface::serviceOwnerChanged, this,
                [ = ](const QString &name, const QString &, const QString &newOwner) {
            if (name.startsWith("org.mpris.MediaPlayer2")) {
                // 启动了音乐播放
                m_isActived = !newOwner.isEmpty();
                if (m_isActived) {
                    m_serviceName = name;
                    m_mediaInter = new MediaPlayerInterface(m_serviceName, "/org/mpris/MediaPlayer2", QDBusConnection::sessionBus(), this);
                    connect(m_mediaInter, &MediaPlayerInterface::PlaybackStatusChanged, this, [ this ] {
                        Q_EMIT statusChanged(convertStatus(m_mediaInter->playbackStatus()));
                    });
                    connect(m_mediaInter, &MediaPlayerInterface::MetadataChanged, this, &MediaPlayerModel::metadataChanged);
                    Dict v = m_mediaInter->metadata();
                    m_name = v.value("xesam:title").toString();
                    m_icon = v.value("mpris:artUrl").toString();
                    m_album = v.value("xesam:album").toString();
                    m_artist = v.value("xesam:artist").toString();
                } else {
                    if (!m_serviceName.isEmpty()) {
                        delete m_mediaInter;
                        m_mediaInter = nullptr;
                    }
                    m_serviceName.clear();
                }
                Q_EMIT startStop(m_isActived);
            }
        });
        connect(dbusInterface, &QDBusConnectionInterface::serviceUnregistered, this,
                [ = ](const QString &service) {
            if (service.startsWith("org.mpris.MediaPlayer2")) {
                // 启动了音乐播放
                m_serviceName.clear();
                m_isActived = false;
                Q_EMIT startStop(m_isActived);
            }
        });
    });
    connect(watcher, &QDBusPendingCallWatcher::finished, watcher, &QDBusPendingCallWatcher::deleteLater);
}

MediaPlayerModel::PlayStatus MediaPlayerModel::convertStatus(const QString &stat)
{
    if (stat == "Paused")
        return PlayStatus::Pause;
    if (stat == "Playing")
        return PlayStatus::Play;
    if (stat == "Stopped")
        return PlayStatus::Stop;

    return PlayStatus::Unknow;
}

MediaPlayerInterface::MediaPlayerInterface(const QString &service, const QString &path, const QDBusConnection &connection, QObject *parent)
    : QDBusAbstractInterface(service, path, "org.mpris.MediaPlayer2.Player", connection, parent)
{
    QDBusConnection::sessionBus().connect(this->service(), this->path(), "org.freedesktop.DBus.Properties",  "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(const QDBusMessage &)));
}

MediaPlayerInterface::~MediaPlayerInterface()
{
    QDBusConnection::sessionBus().disconnect(this->service(), this->path(), "org.freedesktop.DBus.Properties",  "PropertiesChanged", "sa{sv}as", this, SLOT(onPropertyChanged(const QDBusMessage &)));
}

void MediaPlayerInterface::onPropertyChanged(const QDBusMessage &msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QString interfaceName = msg.arguments().at(0).toString();
    if (interfaceName !="org.mpris.MediaPlayer2.Player")
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    QStringList keys = changedProps.keys();
    foreach(const QString &prop, keys) {
    const QMetaObject* self = metaObject();
        for (int i = self->propertyOffset(); i < self->propertyCount(); ++i) {
            QMetaProperty p = self->property(i);
            if (p.name() == prop) {
                Q_EMIT p.notifySignal().invoke(this);
            }
        }
    }
}
