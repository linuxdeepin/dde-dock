/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "sounddeviceswidget.h"
#include "brightnessmodel.h"
#include "settingdelegate.h"
#include "imageutil.h"
#include "slidercontainer.h"

#include <DListView>
#include <DPushButton>
#include <DLabel>
#include <DGuiApplicationHelper>
#include <DDBusSender>
#include <DBlurEffectWidget>
#include <DPaletteHelper>

#include <QVBoxLayout>
#include <QScrollBar>
#include <QEvent>
#include <QProcess>
#include <QDBusInterface>
#include <QDBusConnection>

DWIDGET_USE_NAMESPACE

#define HEADERHEIGHT 30
#define ITEMSPACE 16

#define AUDIOPORT 0
#define AUDIOSETTING 1

const int cardIdRole = itemFlagRole + 1;

SoundDevicesWidget::SoundDevicesWidget(QWidget *parent)
    : QWidget(parent)
    , m_sliderParent(new QWidget(this))
    , m_sliderContainer(new SliderContainer(m_sliderParent))
    , m_descriptionLabel(new QLabel(tr("Output Device"), this))
    , m_deviceList(new DListView(this))
    , m_volumeModel(new DBusAudio("org.deepin.daemon.Audio1", "/org/deepin/daemon/Audio1", QDBusConnection::sessionBus(), this))
    , m_audioSink(new DBusSink("org.deepin.daemon.Audio1", m_volumeModel->defaultSink().path(), QDBusConnection::sessionBus(), this))
    , m_model(new QStandardItemModel(this))
    , m_delegate(new SettingDelegate(m_deviceList))
{
    initUi();
    initConnection();
    onAudioDevicesChanged();
    m_sliderParent->installEventFilter(this);

    QMetaObject::invokeMethod(this, [ this ] {
        resetVolumeInfo();
        resizeHeight();
    }, Qt::QueuedConnection);
}

SoundDevicesWidget::~SoundDevicesWidget()
{
}

bool SoundDevicesWidget::eventFilter(QObject *watcher, QEvent *event)
{
    if ((watcher == m_sliderParent) && (event->type() == QEvent::Paint)) {
        QPainter painter(m_sliderParent);
        painter.setRenderHint(QPainter::Antialiasing); // 抗锯齿
        painter.setPen(Qt::NoPen);

        DPalette dpa = DPaletteHelper::instance()->palette(m_sliderParent);
        painter.setBrush(dpa.brush(DPalette::ColorRole::Midlight));
        painter.drawRoundedRect(m_sliderParent->rect(), 10, 10);
    }

    return QWidget::eventFilter(watcher, event);
}

void SoundDevicesWidget::initUi()
{
    QVBoxLayout *layout = new QVBoxLayout(this);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->setSpacing(6);

    m_sliderParent->setFixedHeight(36);

    QHBoxLayout *sliderLayout = new QHBoxLayout(m_sliderParent);
    sliderLayout->setContentsMargins(11, 0, 11, 0);
    sliderLayout->setSpacing(0);

    QPixmap leftPixmap = ImageUtil::loadSvg(leftIcon(), QSize(24, 24));
    m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon, leftPixmap, QSize(), 5);
    QPixmap rightPixmap = ImageUtil::loadSvg(rightIcon(), QSize(24, 24));
    m_sliderContainer->setIcon(SliderContainer::IconPosition::RightIcon, rightPixmap, QSize(), 7);

    SliderProxyStyle *proxy = new SliderProxyStyle(SliderProxyStyle::Normal);
    m_sliderContainer->setSliderProxyStyle(proxy);
    sliderLayout->addWidget(m_sliderContainer);

    QHBoxLayout *topLayout = new QHBoxLayout(this);
    topLayout->setContentsMargins(10, 0, 10, 0);
    topLayout->setSpacing(0);
    topLayout->addWidget(m_sliderParent);

    layout->addLayout(topLayout);
    layout->addSpacing(4);
    layout->addWidget(m_descriptionLabel);

    m_deviceList->setModel(m_model);
    m_deviceList->setViewMode(QListView::ListMode);
    m_deviceList->setMovement(QListView::Free);
    m_deviceList->setItemRadius(12);
    m_deviceList->setWordWrap(false);
    m_deviceList->verticalScrollBar()->setVisible(false);
    m_deviceList->horizontalScrollBar()->setVisible(false);
    m_deviceList->setOrientation(QListView::Flow::TopToBottom, false);
    layout->addWidget(m_deviceList);
    m_deviceList->setSpacing(10);

    m_deviceList->setItemDelegate(m_delegate);
}

void SoundDevicesWidget::onAudioDevicesChanged()
{
    QList<AudioPort> ports = m_audioSink->ports();
    for (AudioPort port : ports) {
        if (port.availability != 0 && port.availability != 2)
            continue;

        uint cardId = audioPortCardId(port);
        if (!m_volumeModel->IsPortEnabled(cardId, port.name))
            continue;

        DStandardItem *item = new DStandardItem;
        item->setText(QString("%1(%2)").arg(port.description).arg(port.name));
        item->setIcon(QIcon(soundIconFile(port)));
        item->setFlags(Qt::NoItemFlags);
        item->setData(port.availability == 2, itemCheckRole);
        item->setData(QVariant::fromValue(port), itemDataRole);
        item->setData(AUDIOPORT, itemFlagRole);
        item->setData(cardId, cardIdRole);
        m_model->appendRow(item);
        if (port.availability == 2)
            m_deviceList->setCurrentIndex(m_model->indexFromItem(item));
    }

    DStandardItem *settingItem = new DStandardItem;
    settingItem->setText(tr("Sound settings"));
    settingItem->setFlags(Qt::NoItemFlags);
    settingItem->setData(false, itemCheckRole);
    settingItem->setData(AUDIOSETTING, itemFlagRole);
    m_model->appendRow(settingItem);
}

void SoundDevicesWidget::initConnection()
{
    connect(m_audioSink, &DBusSink::VolumeChanged, m_sliderContainer, &SliderContainer::updateSliderValue);
    connect(m_volumeModel, &DBusAudio::DefaultSinkChanged, this, &SoundDevicesWidget::onDefaultSinkChanged);
    connect(m_delegate, &SettingDelegate::selectIndexChanged, this, &SoundDevicesWidget::onSelectIndexChanged);
    connect(m_volumeModel, &DBusAudio::PortEnabledChanged, this, &SoundDevicesWidget::onAudioDevicesChanged);
    connect(m_volumeModel, &DBusAudio::CardsWithoutUnavailableChanged, this, &SoundDevicesWidget::onAudioDevicesChanged);

    connect(m_sliderContainer, &SliderContainer::sliderValueChanged, this, [ this ](int value) {
        m_audioSink->SetVolume(value, true);
    });
}

QString SoundDevicesWidget::leftIcon()
{
    QString iconLeft = QString(":/icons/resources/audio-volume-%1").arg(m_audioSink->mute() ? "muted" : "low");
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconLeft.append("-dark");

    return iconLeft;
}

QString SoundDevicesWidget::rightIcon()
{
    QString iconRight = QString(":/icons/resources/audio-volume-high");
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconRight.append("-dark");

    return iconRight;
}

const QString SoundDevicesWidget::soundIconFile(const AudioPort port) const
{
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        return QString(":/icons/resources/ICON_Device_Laptop_dark.svg");

    return QString(":/icons/resources/ICON_Device_Laptop.svg");
}

void SoundDevicesWidget::resizeHeight()
{
    m_deviceList->adjustSize();
    QMargins m = layout()->contentsMargins();
    int height = m.top() + m.bottom() + HEADERHEIGHT + m_sliderContainer->height() + ITEMSPACE
            + m_descriptionLabel->height() + m_deviceList->height();

    setFixedHeight(height);
}

void SoundDevicesWidget::resetVolumeInfo()
{
    m_sliderContainer->updateSliderValue(m_audioSink->volume() * 100);
}

uint SoundDevicesWidget::audioPortCardId(const AudioPort &audioport) const
{
    QString cards = m_volumeModel->cardsWithoutUnavailable();
    QJsonParseError error;
    QJsonDocument json = QJsonDocument::fromJson(cards.toLocal8Bit(), &error);
    if (error.error != QJsonParseError::NoError)
        return -1;

    QJsonArray array = json.array();
    for (const QJsonValue value : array) {
        QJsonObject cardObject = value.toObject();
        uint cardId = static_cast<uint>(cardObject.value("Id").toInt());
        QJsonArray jPorts = cardObject.value("Ports").toArray();
        for (const QJsonValue jPortValue : jPorts) {
             QJsonObject jPort = jPortValue.toObject();
             if (!jPort.value("Enabled").toBool())
                 continue;

             int direction = jPort.value("Direction").toInt();
             if (direction != 1)
                 continue;

             if (jPort.value("Name").toString() == audioport.name)
                 return cardId;
        }
    }

    return -1;
}

void SoundDevicesWidget::onSelectIndexChanged(const QModelIndex &index)
{
    int flag = index.data(itemFlagRole).toInt();
    if (flag == AUDIOPORT) {
        // 如果是点击具体的声音设备
        AudioPort port = index.data(itemDataRole).value<AudioPort>();
        uint cardId = index.data(cardIdRole).toUInt();
        if (cardId >= 0) {
            m_volumeModel->SetPort(cardId, port.name, 1);
            m_deviceList->update();
        }
    } else {
        // 如果是点击声音设置，则打开控制中心的声音模块
        DDBusSender().service("org.deepin.dde.ControlCenter1")
                .path("/org/deepin/dde/ControlCenter1")
                .interface("org.deepin.dde.ControlCenter1")
                .method("ShowPage").arg(QString("sound")).call();
        hide();
    }
}

void SoundDevicesWidget::onDefaultSinkChanged(const QDBusObjectPath &value)
{
    delete m_audioSink;
    m_audioSink = new DBusSink("org.deepin.daemon.Audio1", m_volumeModel->defaultSink().path(), QDBusConnection::sessionBus(), this);
    connect(m_audioSink, &DBusSink::VolumeChanged, m_sliderContainer, &SliderContainer::updateSliderValue);

    resetVolumeInfo();
    m_deviceList->update();
}
