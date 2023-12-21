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
#include <QMetaObject>
#include <QUrl>

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
    return m_name;
}

const QString MediaPlayerModel::iconUrl()
{
    return m_icon;
}

const QString MediaPlayerModel::album()
{
    return m_album;
}

const QString MediaPlayerModel::artist()
{
    return m_artist;
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

void MediaPlayerModel::updateMetadata()
{
    if (!m_mediaInter)
        return;

    Dict v = m_mediaInter->metadata();
    m_name = v.value("xesam:title").toString();
    m_icon = v.value("mpris:artUrl").toString();
    m_album = v.value("xesam:album").toString();
    m_artist = v.value("xesam:artist").toString();

    auto getName = [&v](const QString &service){
        if (service.contains("vlc", Qt::CaseInsensitive)) {
            const QString &url = v.value("xesam:url").toString();
            if (!url.isEmpty())
                return QUrl(url).fileName();
        }

        return tr("Unknown");
    };
    auto getIcon = [](const QString &service){
        QMap<QString, QString> serv2Icon = {
            {"vlc", "vlc"},
            {"chromium", "chrome"},
            {"firefox", "firefox"},
            {"movie", "video"},
            {"music", "music"}
        };

        for(auto k : serv2Icon.keys())
            if (service.contains(k, Qt::CaseInsensitive))
                return serv2Icon.value(k);

        return QString("music");
    };

    const QString &service = m_mediaInter->service();
    if (m_name.isEmpty())
        m_name = getName(service);

    if (m_icon.isEmpty())
        m_icon = getIcon(service);

    if (m_album.isEmpty())
        m_album = tr("Unknown");
    if (m_artist.isEmpty())
        m_artist = tr("Unknown");

    Q_EMIT MediaPlayerModel::metadataChanged();
}

void MediaPlayerModel::onServiceDiscovered(const QString &service)
{
    auto mediaInter = new MediaPlayerInterface(service, "/org/mpris/MediaPlayer2",
                                               QDBusConnection::sessionBus(), this);
    // 影院不太希望被控制。。canShowInUI:false
    if (!mediaInter->canControl() || !mediaInter->canShowInUI()) {
        delete mediaInter;
        return;
    }

    if (!m_mprisServices.contains(service))
        m_mprisServices << service;

    m_isActived = !m_mprisServices.isEmpty();

    if (m_mediaInter && m_mediaInter->service() == service) {
        Q_EMIT startStop(m_isActived);
        delete mediaInter;
        return;
    }

    if (m_mediaInter)
        delete m_mediaInter;

    m_mediaInter = mediaInter;

    updateMetadata();

    connect(m_mediaInter, &MediaPlayerInterface::PlaybackStatusChanged, this, [ this ] {
        Q_EMIT statusChanged(convertStatus(m_mediaInter->playbackStatus()));
    });
    connect(m_mediaInter, &MediaPlayerInterface::MetadataChanged, this, &MediaPlayerModel::updateMetadata);

    Q_EMIT startStop(m_isActived);
}

void MediaPlayerModel::onServiceDisappears(const QString &service)
{
    if (!m_mprisServices.contains(service))
        return;

    m_mprisServices.removeAll(service);
    m_isActived = !m_mprisServices.isEmpty();

    if (m_mediaInter && m_mediaInter->service() == service) {
        delete m_mediaInter;
        m_mediaInter = nullptr;
    }

    // 退出当前播放器后，继续控制上一个播放器
    if (m_isActived)
        return onServiceDiscovered(m_mprisServices.last());

    Q_EMIT startStop(m_isActived);
}

void MediaPlayerModel::initMediaPlayer()
{
    QDBusInterface dbusInter("org.freedesktop.DBus", "/", "org.freedesktop.DBus", QDBusConnection::sessionBus(), this);
    QDBusPendingCall call = dbusInter.asyncCall("ListNames");
    QDBusPendingCallWatcher *watcher = new QDBusPendingCallWatcher(call, this);
    connect(watcher, &QDBusPendingCallWatcher::finished, [ = ] {
        m_mprisServices.clear();
        if (call.isError())
            return;

        QDBusReply<QStringList> reply = call.reply();
        const QStringList &serviceList = reply.value();
        for (const QString &serv : serviceList) {
            if (!serv.startsWith("org.mpris.MediaPlayer2"))
                continue;

            onServiceDiscovered(serv);
//            break;
        }

        QDBusConnectionInterface *dbusInterface = QDBusConnection::sessionBus().interface();
        connect(dbusInterface, &QDBusConnectionInterface::serviceOwnerChanged, this,
                [ = ](const QString &name, const QString &, const QString &newOwner) {
            if (!name.startsWith("org.mpris.MediaPlayer2"))
                return;

            if (newOwner.isEmpty()) {
                onServiceDisappears(name);
            } else {
                onServiceDiscovered(name);
            }
        });
        connect(dbusInterface, &QDBusConnectionInterface::serviceUnregistered, this,
                [ = ](const QString &service) {
            if (service.startsWith("org.mpris.MediaPlayer2")) {
                onServiceDiscovered(service);
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
#if QT_VERSION < QT_VERSION_CHECK(6, 0, 0)
            QGenericArgument value(QMetaType::typeName(p.type()), const_cast<void*>(changedProps[prop].constData()));
#else
            QGenericArgument value{p.metaType().name(), const_cast<void*>(changedProps[prop].constData())};
#endif
            if (p.name() == prop) {
                Q_EMIT p.notifySignal().invoke(this, value);
            }
        }
    }
}
