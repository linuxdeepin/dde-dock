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
#include "volumedeviceswidget.h"
#include "brightnessmodel.h"
#include "volumemodel.h"
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

VolumeDevicesWidget::VolumeDevicesWidget(VolumeModel *model, QWidget *parent)
    : QWidget(parent)
    , m_sliderParent(new QWidget(this))
    , m_sliderContainer(new SliderContainer(m_sliderParent))
    , m_descriptionLabel(new QLabel(tr("Output Device"), this))
    , m_deviceList(new DListView(this))
    , m_volumeModel(model)
    , m_audioSink(nullptr)
    , m_model(new QStandardItemModel(this))
    , m_delegate(new SettingDelegate(m_deviceList))
{
    initUi();
    initConnection();
    reloadAudioDevices();
    m_sliderParent->installEventFilter(this);

    QMetaObject::invokeMethod(this, [ this ] {
        resetVolumeInfo();
        resizeHeight();
    }, Qt::QueuedConnection);
}

VolumeDevicesWidget::~VolumeDevicesWidget()
{
}

bool VolumeDevicesWidget::eventFilter(QObject *watcher, QEvent *event)
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

void VolumeDevicesWidget::initUi()
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

void VolumeDevicesWidget::reloadAudioDevices()
{
    QList<AudioPorts *> ports = m_volumeModel->ports();
    for (AudioPorts *port : ports) {
        DStandardItem *item = new DStandardItem;
        item->setText(QString("%1(%2)").arg(port->description()).arg(port->cardName()));
        item->setIcon(QIcon(soundIconFile(port)));
        item->setFlags(Qt::NoItemFlags);
        item->setData(port->isChecked(), itemCheckRole);
        item->setData(QVariant::fromValue(port), itemDataRole);
        m_model->appendRow(item);
        if (port->isChecked())
            m_deviceList->setCurrentIndex(m_model->indexFromItem(item));
    }

    DStandardItem *settingItem = new DStandardItem;
    settingItem->setText(tr("Sound settings"));
    settingItem->setFlags(Qt::NoItemFlags);
    settingItem->setData(false, itemCheckRole);
    m_model->appendRow(settingItem);
}

void VolumeDevicesWidget::initConnection()
{
    m_audioSink = m_volumeModel->defaultSink();

    if (m_audioSink)
        connect(m_audioSink, &AudioSink::volumeChanged, m_sliderContainer, &SliderContainer::updateSliderValue);
    connect(m_volumeModel, &VolumeModel::defaultSinkChanged, this, [ this ](AudioSink *sink) {
        if (m_audioSink)
            disconnect(m_audioSink);

        m_audioSink = sink;
        if (sink)
            connect(sink, &AudioSink::volumeChanged, m_sliderContainer, &SliderContainer::updateSliderValue);

        resetVolumeInfo();
        m_deviceList->update();
    });

    connect(m_sliderContainer, &SliderContainer::sliderValueChanged, this, [ this ](int value) {
        AudioSink *defSink = m_volumeModel->defaultSink();
        if (!defSink)
            return;

        defSink->setVolume(value, true);
    });
    connect(m_volumeModel, &VolumeModel::checkPortChanged, this, [ this ] {
        for (int i = 0; i < m_model->rowCount(); i++) {
            QModelIndex index = m_model->index(i, 0);
            AudioPorts *port = index.data(itemDataRole).value<AudioPorts *>();
            if (port)
                m_model->setData(index, port->isChecked(), itemCheckRole);
        }
        m_deviceList->update();
    });

    connect(m_delegate, &SettingDelegate::selectIndexChanged, this, [ this ](const QModelIndex &index) {
        AudioPorts *port = index.data(itemDataRole).value<AudioPorts *>();
        if (port) {
            m_volumeModel->setActivePort(port);
            m_deviceList->update();
        } else {
            // 打开控制中心的声音模块
            DDBusSender().service("org.deepin.dde.ControlCenter")
                    .path("/org/deepin/dde/ControlCenter")
                    .interface("org.deepin.dde.ControlCenter")
                    .method("ShowPage").arg(QString("sound")).call();
            hide();
        }
    });
}

QString VolumeDevicesWidget::leftIcon()
{
    QString iconLeft = QString(":/icons/resources/audio-volume-%1").arg(m_volumeModel->isMute() ? "muted" : "low");
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconLeft.append("-dark");

    return iconLeft;
}

QString VolumeDevicesWidget::rightIcon()
{
    QString iconRight = QString(":/icons/resources/audio-volume-high");
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconRight.append("-dark");

    return iconRight;
}

const QString VolumeDevicesWidget::soundIconFile(AudioPorts *port) const
{
    if (!port)
        return QString();

    if (port->isHeadPhone()) {
        if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
            return QString(":/icons/resources/ICON_Device_Headphone_dark.svg");

        return QString(":/icons/resources/ICON_Device_Headphone.svg");
    }

    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        return QString(":/icons/resources/ICON_Device_Laptop_dark.svg");

    return QString(":/icons/resources/ICON_Device_Laptop.svg");
}

void VolumeDevicesWidget::resizeHeight()
{
    m_deviceList->adjustSize();
    QMargins m = layout()->contentsMargins();
    int height = m.top() + m.bottom() + HEADERHEIGHT + m_sliderContainer->height() + ITEMSPACE
            + m_descriptionLabel->height() + m_deviceList->height();

    setFixedHeight(height);
}

void VolumeDevicesWidget::resetVolumeInfo()
{
    AudioSink *defaultSink = m_volumeModel->defaultSink();
    if (!defaultSink)
        return;

    m_sliderContainer->updateSliderValue(defaultSink->volume());
}
