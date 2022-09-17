// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SOUNDAPPLET_H
#define SOUNDAPPLET_H

#include "componments/volumeslider.h"

#include <com_deepin_daemon_audio.h>
#include <com_deepin_daemon_audio_sink.h>

#include <DIconButton>
#include <DListView>
#include <DApplicationHelper>

#include <QScrollArea>
#include <QVBoxLayout>
#include <QLabel>
#include <QSlider>

DWIDGET_USE_NAMESPACE

using DBusAudio = com::deepin::daemon::Audio;
using DBusSink = com::deepin::daemon::audio::Sink;

class HorizontalSeperator;
class QGSettings;

namespace Dock {
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

class BackgroundWidget : public QWidget
{
public:
    explicit BackgroundWidget(QWidget *parent = nullptr)
        : QWidget(parent) {}

protected:
    void paintEvent(QPaintEvent *event)
    {
        QPainter painter(this);
        painter.setPen(Qt::NoPen);
        if (DApplicationHelper::instance()->themeType() == DApplicationHelper::LightType) {
            painter.setBrush(QColor(0, 0, 0, 0.03 * 255));
        } else {
            painter.setBrush(QColor(255, 255, 255, 0.03 * 255));
        }
        painter.drawRect(rect());

        return QWidget::paintEvent(event);
    }

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

    bool existActiveOutputDevice();

signals:
    void volumeChanged(const int value) const;
    void defaultSinkChanged(DBusSink *sink) const;

private slots:
    void onDefaultSinkChanged();
    void onVolumeChanged(double volume);
    void volumeSliderValueChanged();
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

protected:
    bool eventFilter(QObject *watcher, QEvent *event) override;

private:
    QWidget *m_centralWidget;
    QLabel *m_volumeIconMin;
    QLabel *m_volumeIconMax;
    VolumeSlider *m_volumeSlider;
    QLabel *m_soundShow;
    QLabel *m_deviceLabel;
    QVBoxLayout *m_centralLayout;
    HorizontalSeperator *m_seperator;
    HorizontalSeperator *m_secondSeperator;

    DBusAudio *m_audioInter;
    DBusSink *m_defSinkInter;
    DTK_WIDGET_NAMESPACE::DListView  *m_listView;
    QStandardItemModel *m_model;
    QList<Port *> m_ports;
    QString m_deviceInfo;
    QPointer<Port> m_lastPort;//最后一个因为只有一个设备而被直接移除的设备
    const QGSettings *m_gsettings;
};

#endif // SOUNDAPPLET_H
