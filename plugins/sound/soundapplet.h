// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SOUNDAPPLET_H
#define SOUNDAPPLET_H

#include "componments/volumeslider.h"

#include "org_deepin_dde_audio1.h"
#include "org_deepin_dde_audio1_sink.h"

#include <DIconButton>
#include <DListView>
#include <DGuiApplicationHelper>

#include <QScrollArea>
#include <QVBoxLayout>
#include <QLabel>
#include <QSlider>

DWIDGET_USE_NAMESPACE

using DBusAudio = org::deepin::dde::Audio1;
using DBusSink = org::deepin::dde::audio1::Sink;
using DTK_NAMESPACE::Gui::DGuiApplicationHelper;

class HorizontalSeperator;
class QGSettings;
class SoundDevicePort;

namespace Dock {
class TipsWidget;
}

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
        if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
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
    void startAddPort(SoundDevicePort *port);
    void startRemovePort(const QString &portId, const uint &cardId);
    bool containsPort(const SoundDevicePort *port);
    SoundDevicePort *findPort(const QString &portId, const uint &cardId) const;
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
    void addPort(const SoundDevicePort *port);
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
    QList<SoundDevicePort *> m_ports;
    QString m_deviceInfo;
    QPointer<SoundDevicePort> m_lastPort;//最后一个因为只有一个设备而被直接移除的设备
    const QGSettings *m_gsettings;
};

#endif // SOUNDAPPLET_H
