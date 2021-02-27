/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef SOUNDAPPLET_H
#define SOUNDAPPLET_H

#include "componments/volumeslider.h"

#include <com_deepin_daemon_audio.h>
#include <com_deepin_daemon_audio_sink.h>

#include <QScrollArea>
#include <QVBoxLayout>
#include <QLabel>
#include <QSlider>

#include <DIconButton>
#include <DListView>

DWIDGET_USE_NAMESPACE

using DBusAudio = com::deepin::daemon::Audio;
using DBusSink = com::deepin::daemon::audio::Sink;

class HorizontalSeparator;
class QGSettings;

namespace Dock{
class TipsWidget;
}
class Port : public QObject
{
    Q_OBJECT
public:
    enum Direction {
        Out = 1,
        In = 2
    };

    explicit Port(QObject *parent = nullptr);
    virtual ~Port() {}

    inline QString id() const { return m_id; }
    void setId(const QString &id);

    inline QString name() const { return m_name; }
    void setName(const QString &name);

    inline QString cardName() const { return m_cardName; }
    void setCardName(const QString &cardName);

    inline bool isActive() const { return m_isActive; }
    void setIsActive(bool isActive);

    inline Direction direction() const { return m_direction; }
    void setDirection(const Direction &direction);

    inline uint cardId() const { return m_cardId; }
    void setCardId(const uint &cardId);

Q_SIGNALS:
    void idChanged(QString id) const;
    void nameChanged(QString name) const;
    void cardNameChanged(QString name) const;
    void isActiveChanged(bool ative) const;
    void directionChanged(Direction direction) const;
    void cardIdChanged(uint cardId) const;

private:
    QString m_id;
    QString m_name;
    uint m_cardId;
    QString m_cardName;
    bool m_isActive;
    Direction m_direction;
};

class SoundApplet : public QScrollArea
{
    Q_OBJECT

public:
    explicit SoundApplet(QWidget *parent = 0);

    int volumeValue() const;
    int maxVolumeValue() const;
    VolumeSlider *mainSlider();
    void startAddPort(Port *port);
    void startRemovePort(const QString &portId, const uint &cardId);
    bool containsPort(const Port *port);
    Port *findPort(const QString &portId, const uint &cardId) const;
    void setUnchecked(DStandardItem *pi);
    void initUi();

signals:
    void volumeChanged(const int value) const;
    void defaultSinkChanged(DBusSink *sink) const;

private slots:
    void defaultSinkChanged();
    void onVolumeChanged(double volume);
    void volumeSliderValueChanged();
    void sinkInputsChanged();
    void toggleMute();
    void onPlaySoundEffect();
    void increaseVolumeChanged();
    void cardsChanged(const QString &cards);
    void removePort(const QString &portId, const uint &cardId);
    void addPort(const Port *port);
    void activePort(const QString &portId,const uint &cardId);
    void haldleDbusSignal(const QDBusMessage &msg);
    void updateListHeight();
    void portEnableChange(unsigned int cardId, QString portId);

private:
    void refreshIcon();
    void updateCradsInfo();
    void enableDevice(bool flag);
    void disableAllDevice();//禁用所有设备
    void removeLastDevice();//移除最后一个设备
    void removeDisabledDevice(QString portId, unsigned int cardId);
    void updateVolumeSliderStatus(const QString &status);

private:
    QWidget *m_centralWidget;
    DIconButton *m_volumeBtn;
    QLabel *m_volumeIconMax;
    VolumeSlider *m_volumeSlider;
    Dock::TipsWidget *m_soundShow;
    QVBoxLayout *m_centralLayout;
    HorizontalSeparator *m_separator;
    Dock::TipsWidget *m_deviceLabel;

    DBusAudio *m_audioInter;
    DBusSink *m_defSinkInter;
    DTK_WIDGET_NAMESPACE::DListView  *m_listView;
    QStandardItemModel *m_model;
    QList<Port *> m_ports;
    QString m_deviceInfo;
    QPointer<Port> m_lastPort;//最后一个因为只有一个设备而被直接移除的设备
    QGSettings *m_gsettings;
};

#endif // SOUNDAPPLET_H
