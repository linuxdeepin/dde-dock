// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "sounddeviceswidget.h"
#include "settingdelegate.h"
#include "imageutil.h"
#include "slidercontainer.h"
#include "sounddeviceport.h"

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
#define ROWSPACE 10
#define DESCRIPTIONHEIGHT 15
#define ITEMRADIUS 12
#define SLIDERHEIGHT 36

#define AUDIOPORT 0
#define AUDIOSETTING 1

const int sortRole = itemFlagRole + 1;

SoundDevicesWidget::SoundDevicesWidget(QWidget *parent)
    : QWidget(parent)
    , m_sliderParent(new QWidget(this))
    , m_sliderContainer(new SliderContainer(m_sliderParent))
    , m_descriptionLabel(new QLabel(tr("Output Device"), this))
    , m_deviceList(new DListView(this))
    , m_soundInter(new DBusAudio("org.deepin.dde.Audio1", "/org/deepin/dde/Audio1", QDBusConnection::sessionBus(), this))
    , m_sinkInter(new DBusSink("org.deepin.dde.Audio1", m_soundInter->defaultSink().path(), QDBusConnection::sessionBus(), this))
    , m_model(new QStandardItemModel(this))
    , m_delegate(new SettingDelegate(m_deviceList))
{
    initUi();
    initConnection();
    onAudioDevicesChanged();

    QMetaObject::invokeMethod(this, [ this ] {
        deviceEnabled(m_ports.size() > 0);
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

    m_sliderParent->setFixedHeight(SLIDERHEIGHT);

    QHBoxLayout *sliderLayout = new QHBoxLayout(m_sliderParent);
    sliderLayout->setContentsMargins(11, 0, 11, 0);
    sliderLayout->setSpacing(0);

    QPixmap leftPixmap = ImageUtil::loadSvg(leftIcon(), QSize(24, 24));
    m_sliderContainer->setIcon(SliderContainer::IconPosition::LeftIcon, leftPixmap, QSize(), 5);
    QPixmap rightPixmap = ImageUtil::loadSvg(rightIcon(), QSize(24, 24));
    m_sliderContainer->setIcon(SliderContainer::IconPosition::RightIcon, rightPixmap, QSize(), 7);

    SliderProxyStyle *proxy = new SliderProxyStyle(SliderProxyStyle::Normal);
    m_sliderContainer->setSliderProxyStyle(proxy);
    m_sliderContainer->setRange(0, std::round(m_soundInter->maxUIVolume() * 100.00));
    m_sliderContainer->setPageStep(2);
    sliderLayout->addWidget(m_sliderContainer);

    QHBoxLayout *topLayout = new QHBoxLayout(this);
    topLayout->setContentsMargins(10, 0, 10, 0);
    topLayout->setSpacing(0);
    topLayout->addWidget(m_sliderParent);

    layout->addLayout(topLayout);
    layout->addWidget(m_descriptionLabel);

    m_deviceList->setModel(m_model);
    m_deviceList->setViewMode(QListView::ListMode);
    m_deviceList->setMovement(QListView::Free);
    m_deviceList->setItemRadius(ITEMRADIUS);
    m_deviceList->setWordWrap(false);
    m_deviceList->verticalScrollBar()->setVisible(false);
    m_deviceList->horizontalScrollBar()->setVisible(false);
    m_deviceList->setOrientation(QListView::Flow::TopToBottom, false);
    layout->addWidget(m_deviceList);
    m_deviceList->setSpacing(ROWSPACE);

    m_deviceList->setItemDelegate(m_delegate);
    m_model->setSortRole(sortRole);
    m_descriptionLabel->setFixedHeight(DESCRIPTIONHEIGHT);

    // 增加音量设置
    DStandardItem *settingItem = new DStandardItem;
    settingItem->setText(tr("Sound settings"));
    settingItem->setFlags(Qt::NoItemFlags);
    settingItem->setData(false, itemCheckRole);
    settingItem->setData(AUDIOSETTING, itemFlagRole);
    m_model->appendRow(settingItem);

    m_sliderParent->installEventFilter(this);
}

void SoundDevicesWidget::onAudioDevicesChanged()
{
    QMap<uint, QStringList> tmpCardIds;
    const QString cards = m_soundInter->cardsWithoutUnavailable();
    QJsonDocument doc = QJsonDocument::fromJson(cards.toUtf8());
    QJsonArray jCards = doc.array();
    for (QJsonValue cV : jCards) {
        QJsonObject jCard = cV.toObject();
        const uint cardId = jCard["Id"].toInt();
        const QString cardName = jCard["Name"].toString();
        QJsonArray jPorts = jCard["Ports"].toArray();

        QStringList tmpPorts;

        for (QJsonValue pV : jPorts) {
            QJsonObject jPort = pV.toObject();
            const double portAvai = jPort["Available"].toDouble();
            if (portAvai != 2 && portAvai != 0 )
                continue;

            const QString portId = jPort["Name"].toString();
            const QString portName = jPort["Description"].toString();

            SoundDevicePort *port = findPort(portId, cardId);
            bool includePort = (port != nullptr);
            if (!port)
                port = new SoundDevicePort(m_model);

            port->setId(portId);
            port->setName(portName);
            port->setDirection(SoundDevicePort::Direction(jPort["Direction"].toDouble()));
            port->setCardId(cardId);
            port->setCardName(cardName);

            if (!includePort)
                startAddPort(port);

            tmpPorts << portId;
        }
        tmpCardIds.insert(cardId, tmpPorts);
    }

    // 重新获取切换的设备信息
    onDefaultSinkChanged(m_soundInter->defaultSink());

    for (SoundDevicePort *port : m_ports) {
        // 只要有一个设备在控制中心被禁用后，在任务栏声音设备列表中该设备会被移除，
        if (!m_soundInter->IsPortEnabled(port->cardId(), port->id()))
            removeDisabledDevice(port->id(), port->cardId());

        // 判断端口是否在最新的设备列表中
        if (!tmpCardIds.contains(port->cardId()) || !tmpCardIds[port->cardId()].contains(port->id()))
            startRemovePort(port->id(), port->cardId());
    }
}

void SoundDevicesWidget::initConnection()
{
    connect(m_sinkInter, &DBusSink::VolumeChanged, this, [ = ](double value) { m_sliderContainer->updateSliderValue(value * 100); });
    connect(m_sinkInter, &DBusSink::MuteChanged, this, [ = ] { m_sliderContainer->updateSliderValue(m_sinkInter->volume() * 100); });
    connect(m_soundInter, &DBusAudio::DefaultSinkChanged, this, &SoundDevicesWidget::onDefaultSinkChanged);
    connect(m_delegate, &SettingDelegate::selectIndexChanged, this, &SoundDevicesWidget::onSelectIndexChanged);
    connect(m_soundInter, &DBusAudio::PortEnabledChanged, this, &SoundDevicesWidget::onAudioDevicesChanged);
    connect(m_soundInter, &DBusAudio::CardsWithoutUnavailableChanged, this, &SoundDevicesWidget::onAudioDevicesChanged);
    connect(m_soundInter, &DBusAudio::MaxUIVolumeChanged, this, [ = ] (double maxValue) {
        m_sliderContainer->setRange(0, std::round(maxValue * 100.00));
    });
    connect(m_sliderContainer, &SliderContainer::sliderValueChanged, this, [ this ](int value) {
        m_sinkInter->SetVolume(value * 0.01, true);
    });
}

/**
 * @brief SoundApplet::startAddPort 添加端口前判断
 * @param port 端口
 */
void SoundDevicesWidget::startAddPort(SoundDevicePort *port)
{
    if (findPort(port->id(), port->cardId()))
        return;

    if (port->direction() == SoundDevicePort::Out) {
        m_ports.append(port);
        addPort(port);
    }
}

void SoundDevicesWidget::startRemovePort(const QString &portId, const uint &cardId)
{
    SoundDevicePort *port = findPort(portId, cardId);
    if (port) {
        m_ports.removeOne(port);
        port->deleteLater();
        removePort(portId, cardId);
    }
}

void SoundDevicesWidget::addPort(const SoundDevicePort *port)
{
    DStandardItem *portItem = new DStandardItem;
    QString deviceName = port->name() + "(" + port->cardName() + ")";
    portItem->setIcon(QIcon(soundIconFile()));
    portItem->setText(deviceName);
    portItem->setTextColorRole(QPalette::BrightText);
    portItem->setData(QVariant::fromValue<const SoundDevicePort *>(port), itemDataRole);
    portItem->setData(port->isActive(), itemCheckRole);
    portItem->setData(AUDIOPORT, itemFlagRole);

    connect(port, &SoundDevicePort::nameChanged, this, [ = ](const QString &str) {
        QString devName = str + "(" + port->cardName() + ")";
        portItem->setText(devName);
    });
    connect(port, &SoundDevicePort::cardNameChanged, this, [ = ](const QString &str) {
        QString devName = port->name() + "(" + str + ")";
        portItem->setText(devName);
    });
    connect(port, &SoundDevicePort::isActiveChanged, this, [ = ](bool isActive) {
        portItem->setCheckState(isActive ? Qt::CheckState::Checked : Qt::CheckState::Unchecked);
    });

    if (port->isActive()) {
        portItem->setCheckState(Qt::CheckState::Checked);
    }

    m_model->appendRow(portItem);
    // 遍历列表，依次对不同的设备排序
    int row = 0;
    int rowCount = m_model->rowCount();
    for (int i = 0; i < rowCount; i++) {
        QStandardItem *item = m_model->item(i);
        if (item->data(itemFlagRole).toInt() == AUDIOSETTING) {
            item->setData(rowCount - 1, sortRole);
        } else {
            item->setData(row, sortRole);
            row++;
        }
    }

    m_model->sort(0);
    if (m_ports.size() == 1)
        deviceEnabled(true);

    resizeHeight();
}

void SoundDevicesWidget::removePort(const QString &portId, const uint &cardId)
{
    int removeRow = -1;
    for (int i = 0; i < m_model->rowCount(); i++) {
        QStandardItem *item = m_model->item(i);
        if (item->data(itemFlagRole).toInt() != AUDIOPORT)
            continue;

        const SoundDevicePort *port = item->data(itemDataRole).value<const SoundDevicePort *>();
        if (port && port->id() == portId && cardId == port->cardId()) {
            removeRow = i;
            break;
        }
    }

    if (removeRow >= 0)
        m_model->removeRow(removeRow);

    if (m_ports.size() == 0)
        deviceEnabled(false);

    resizeHeight();
}

void SoundDevicesWidget::activePort(const QString &portId, const uint &cardId)
{
    for (SoundDevicePort *it : m_ports)
        it->setIsActive(it->id() == portId && it->cardId() == cardId);
}

void SoundDevicesWidget::removeDisabledDevice(QString portId, unsigned int cardId)
{
    startRemovePort(portId, cardId);
    if (m_sinkInter->activePort().name == portId && m_sinkInter->card() == cardId) {
        for (SoundDevicePort *port : m_ports)
            port->setIsActive(false);
    }
}

void SoundDevicesWidget::deviceEnabled(bool enable)
{
    m_sliderContainer->setEnabled(enable);
    Q_EMIT enableChanged(enable);
}

QString SoundDevicesWidget::leftIcon()
{
    QString iconLeft = QString(":/icons/resources/audio-volume-%1").arg(m_sinkInter->mute() ? "muted" : "low");
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

const QString SoundDevicesWidget::soundIconFile() const
{
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        return QString(":/icons/resources/ICON_Device_Laptop_dark.svg");

    return QString(":/icons/resources/ICON_Device_Laptop.svg");
}

void SoundDevicesWidget::resizeHeight()
{
    int deviceListHeight = 0;
    for (int i = 0; i < m_model->rowCount(); i++) {
        QModelIndex index = m_model->index(i, 0);
        deviceListHeight += m_deviceList->visualRect(index).height() + m_deviceList->spacing() * 2;
        if (i >= 9)  // 最大显示12个条目，包含标题、滑块和设备标题、设备列表，因此，设备列表最多显示9个条目
            break;
    }
    m_deviceList->setFixedHeight(deviceListHeight);
    QMargins m = layout()->contentsMargins();
    int height = m.top() + m.bottom() + HEADERHEIGHT + m_sliderContainer->height() + ITEMSPACE
            + m_descriptionLabel->height() + m_deviceList->height();

    setFixedHeight(height);
}

void SoundDevicesWidget::resetVolumeInfo()
{
    //无声卡状态下，会有伪sink设备，显示音量为0
    m_sliderContainer->updateSliderValue(findPort(m_sinkInter->activePort().name, m_sinkInter->card()) != nullptr ? m_sinkInter->volume() * 100 : 0);
}

uint SoundDevicesWidget::audioPortCardId(const AudioPort &audioport) const
{
    QString cards = m_soundInter->cardsWithoutUnavailable();
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

SoundDevicePort *SoundDevicesWidget::findPort(const QString &portId, const uint &cardId) const
{
    auto it = std::find_if(m_ports.begin(), m_ports.end(), [ = ] (SoundDevicePort *p) {
        return (p->id() == portId && p->cardId() == cardId);
    });

    if (it != m_ports.end()) {
        return *it;
    }

    return nullptr;
}

void SoundDevicesWidget::onSelectIndexChanged(const QModelIndex &index)
{
    int flag = index.data(itemFlagRole).toInt();
    if (flag == AUDIOPORT) {
        const SoundDevicePort *port = m_model->data(index, itemDataRole).value<const SoundDevicePort *>();
        if (port) {
            m_soundInter->SetPort(port->cardId(), port->id(), int(port->direction()));
            //手动勾选时启用设备
            m_soundInter->SetPortEnabled(port->cardId(), port->id(), true);
            m_deviceList->update();
        }
    } else {
        // 如果是点击声音设置，则打开控制中心的声音模块
        DDBusSender().service("org.deepin.dde.ControlCenter1")
                .path("/org/deepin/dde/ControlCenter1")
                .interface("org.deepin.dde.ControlCenter1")
                .method("ShowPage").arg(QString("sound")).call();
        emit requestHide();
    }
}

void SoundDevicesWidget::onDefaultSinkChanged(const QDBusObjectPath &value)
{
    delete m_sinkInter;
    m_sinkInter = new DBusSink("org.deepin.dde.Audio1", m_soundInter->defaultSink().path(), QDBusConnection::sessionBus(), this);
    connect(m_sinkInter, &DBusSink::VolumeChanged, this, [ = ](double value) { m_sliderContainer->updateSliderValue(value * 100); });
    connect(m_sinkInter, &DBusSink::MuteChanged, this, [ = ] { m_sliderContainer->updateSliderValue(m_sinkInter->volume() * 100); });

    QString portId = m_sinkInter->activePort().name;
    uint cardId = m_sinkInter->card();
    activePort(portId, cardId);

    for (int i = 0; i < m_model->rowCount() ; i++) {
        QStandardItem *item = m_model->item(i);
        if (item->data(itemFlagRole).toInt() != AUDIOPORT)
            continue;

        const SoundDevicePort *soundPort = item->data(itemDataRole).value<const SoundDevicePort *>();
        item->setData((soundPort && soundPort->id() == portId && soundPort->cardId() == cardId), itemCheckRole);
    }

    resetVolumeInfo();
    m_deviceList->update();
}
