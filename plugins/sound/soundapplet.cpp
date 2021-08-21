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

#include "soundapplet.h"
#include "util/horizontalseperator.h"
#include "../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"
#include "util/utils.h"

#include <DGuiApplicationHelper>
#include <DApplication>
#include <DStandardItem>
#include <DFontSizeManager>
#include <DApplicationHelper>

#include <QLabel>
#include <QIcon>
#include <QScrollBar>
#include <QPainter>
#include <QListIterator>

#define SEPARATOR_HEIGHT 2
#define WIDTH       260
#define MAX_HEIGHT  300
#define ICON_SIZE   24
#define ITEM_HEIGHT 24
#define ITEM_SPACING 5
#define DEVICE_SPACING 10
#define SLIDER_HIGHT 70
#define TITLE_HEIGHT 46
#define GSETTING_SOUND_OUTPUT_SLIDER "soundOutputSlider"

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE
using namespace Dock;

Q_DECLARE_METATYPE(const Port *)

Port::Port(QObject *parent)
    : QObject(parent)
    , m_isActive(false)
    , m_direction(Out)
{

}

void Port::setId(const QString &id)
{
    if (id != m_id) {
        m_id = id;
        Q_EMIT idChanged(id);
    }
}

void Port::setName(const QString &name)
{
    if (name != m_name) {
        m_name = name;
        Q_EMIT nameChanged(name);
    }
}

void Port::setCardName(const QString &cardName)
{
    if (cardName != m_cardName) {
        m_cardName = cardName;
        Q_EMIT cardNameChanged(cardName);
    }
}

void Port::setIsActive(bool isActive)
{
    if (isActive != m_isActive) {
        m_isActive = isActive;
        Q_EMIT isActiveChanged(isActive);
    }
}

void Port::setDirection(const Direction &direction)
{
    if (direction != m_direction) {
        m_direction = direction;
        Q_EMIT directionChanged(direction);
    }
}

void Port::setCardId(const uint &cardId)
{
    if (cardId != m_cardId) {
        m_cardId = cardId;
        Q_EMIT cardIdChanged(cardId);
    }
}

SoundApplet::SoundApplet(QWidget *parent)
    : QScrollArea(parent)
    , m_centralWidget(new QWidget(this))
    , m_volumeIconMin(new QLabel(this))
    , m_volumeIconMax(new QLabel(this))
    , m_volumeSlider(new VolumeSlider(this))
    , m_soundShow(new QLabel(this))
    , m_deviceLabel(new QLabel(this))
    , m_seperator(new HorizontalSeperator(this))
    , m_secondSeperator(new HorizontalSeperator(this))
    , m_audioInter(new DBusAudio("com.deepin.daemon.Audio", "/com/deepin/daemon/Audio", QDBusConnection::sessionBus(), this))
    , m_defSinkInter(nullptr)
    , m_listView(new DListView(this))
    , m_model(new QStandardItemModel(m_listView))
    , m_deviceInfo("")
    , m_lastPort(nullptr)
    , m_gsettings(Utils::ModuleSettingsPtr("sound", QByteArray(), this))
{
    initUi();

    m_volumeIconMin->installEventFilter(this);
}

void SoundApplet::initUi()
{
    m_listView->setFrameShape(QFrame::NoFrame);
    m_listView->setEditTriggers(DListView::NoEditTriggers);
    m_listView->setSelectionMode(QAbstractItemView::NoSelection);
    m_listView->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_listView->setBackgroundType(DStyledItemDelegate::NoBackground);
    m_listView->setItemRadius(0);
    m_listView->setSizeAdjustPolicy(QAbstractScrollArea::AdjustToContents);
    m_listView->setFixedHeight(0);
    m_listView->setSizePolicy(QSizePolicy::Preferred, QSizePolicy::Preferred);
    m_listView->setModel(m_model);

    m_centralWidget->setAccessibleName("volumn-centralwidget");
    m_volumeIconMin->setAccessibleName("volume-button");
    m_volumeIconMax->setAccessibleName("volume-iconmax");
    m_volumeSlider->setAccessibleName("volume-slider");
    m_soundShow->setAccessibleName("volume-soundtips");
    horizontalScrollBar()->setAccessibleName("volume-horizontalscrollbar");
    verticalScrollBar()->setAccessibleName("volume-verticalscrollbar");

    m_volumeIconMin->setFixedSize(ICON_SIZE, ICON_SIZE);
    m_volumeIconMax->setFixedSize(ICON_SIZE, ICON_SIZE);

    m_soundShow->setText(QString("%1%").arg(0));
    m_soundShow->setFixedHeight(TITLE_HEIGHT);
    m_soundShow->setForegroundRole(QPalette::BrightText);
    DFontSizeManager::instance()->bind(m_soundShow, DFontSizeManager::T8, QFont::Medium);

    m_deviceLabel->setText(tr("Device"));
    m_deviceLabel->setFixedHeight(TITLE_HEIGHT);
    m_deviceLabel->setForegroundRole(QPalette::BrightText);
    DFontSizeManager::instance()->bind(m_deviceLabel, DFontSizeManager::T4, QFont::Medium);

    m_volumeSlider->setFixedHeight(SLIDER_HIGHT);
    m_volumeSlider->setMinimum(0);
    m_volumeSlider->setMaximum(m_audioInter->maxUIVolume() * 100.0f);

    // 标题部分
    QHBoxLayout *deviceLayout = new QHBoxLayout;
    deviceLayout->setSpacing(0);
    deviceLayout->setMargin(0);
    deviceLayout->setContentsMargins(20, 0, 10, 0);
    deviceLayout->addWidget(m_deviceLabel, 0, Qt::AlignLeft);
    deviceLayout->addWidget(m_soundShow, 0, Qt::AlignRight);

    BackgroundWidget *deviceWidget = new BackgroundWidget(this);
    deviceWidget->setLayout(deviceLayout);

    // 音量滑动条
    QHBoxLayout *volumeCtrlLayout = new QHBoxLayout;
    volumeCtrlLayout->setSpacing(0);
    volumeCtrlLayout->setMargin(0);
    volumeCtrlLayout->setContentsMargins(12, 0, 12, 0);
    volumeCtrlLayout->addWidget(m_volumeIconMin);
    volumeCtrlLayout->addWidget(m_volumeSlider);
    volumeCtrlLayout->addWidget(m_volumeIconMax);

    BackgroundWidget *volumnWidget = new BackgroundWidget(this);
    volumnWidget->setLayout(volumeCtrlLayout);

    m_centralLayout = new QVBoxLayout(this);
    m_centralLayout->setContentsMargins(0, 0, 0, 0);
    m_centralLayout->setMargin(0);
    m_centralLayout->setSpacing(0);
    m_centralLayout->addWidget(deviceWidget);
    m_centralLayout->addWidget(m_seperator);
    m_centralLayout->addWidget(volumnWidget);
    // 需要判断是否有声音端口
    m_centralLayout->addWidget(m_secondSeperator);

    m_centralLayout->addWidget(m_listView);
    m_centralWidget->setLayout(m_centralLayout);
    m_centralWidget->setFixedWidth(WIDTH);
    m_centralWidget->setSizePolicy(QSizePolicy::Preferred, QSizePolicy::Preferred);

    setFixedWidth(WIDTH);
    setWidget(m_centralWidget);
    setContentsMargins(0, 0, 0, 0);
    setFrameShape(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);
    m_listView->setItemDelegate(new DStyledItemDelegate(m_listView));

    m_secondSeperator->setVisible(m_model->rowCount() > 1);

    updateVolumeSliderStatus(Utils::SettingValue("com.deepin.dde.dock.module.sound", QByteArray(), "Enabled").toString());
    connect(m_gsettings, &QGSettings::changed, [ = ] (const QString &key) {
        if (key == GSETTING_SOUND_OUTPUT_SLIDER) {
            updateVolumeSliderStatus(m_gsettings->get(GSETTING_SOUND_OUTPUT_SLIDER).toString());
        }
    });
    connect(qApp, &QGuiApplication::fontChanged, this, &SoundApplet::updateListHeight);
    connect(m_volumeSlider, &VolumeSlider::valueChanged, this, &SoundApplet::volumeSliderValueChanged);
    connect(m_audioInter, &DBusAudio::DefaultSinkChanged, this, &SoundApplet::onDefaultSinkChanged);
    connect(m_audioInter, &DBusAudio::IncreaseVolumeChanged, this, &SoundApplet::increaseVolumeChanged);
    connect(m_audioInter, &DBusAudio::PortEnabledChanged, [this](uint cardId, QString portId) {
        portEnableChange(cardId, portId);
    });;
    connect(m_listView, &DListView::clicked, this, [this](const QModelIndex & idx) {
        const Port * port = m_listView->model()->data(idx, Qt::WhatsThisPropertyRole).value<const Port *>();
        if (port) {
            m_audioInter->SetPort(port->cardId(), port->id(), int(port->direction()));
            //手动勾选时启用设备
            m_audioInter->SetPortEnabled(port->cardId(), port->id(), true);
        }

    });
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &SoundApplet::refreshIcon);
    connect(qApp, &DApplication::iconThemeChanged, this, &SoundApplet::refreshIcon);
    QDBusConnection::sessionBus().connect("com.deepin.daemon.Audio", "/com/deepin/daemon/Audio", "org.freedesktop.DBus.Properties"
                                          ,"PropertiesChanged", "sa{sv}as", this, SLOT(haldleDbusSignal(QDBusMessage)));

    QMetaObject::invokeMethod(this, "onDefaultSinkChanged", Qt::QueuedConnection);

    refreshIcon();

    updateCradsInfo();
}

int SoundApplet::volumeValue() const
{
    return m_volumeSlider->value();
}

int SoundApplet::maxVolumeValue() const
{
    return m_volumeSlider->maximum();
}

VolumeSlider *SoundApplet::mainSlider()
{
    return m_volumeSlider;
}

void SoundApplet::onDefaultSinkChanged()
{
    //防止手动切换设备，与后端交互时，获取到多个信号，设备切换多次，造成混乱
    QThread::msleep(200);

    if (m_defSinkInter) {
        delete m_defSinkInter;
        m_defSinkInter = nullptr;
    }

    const QDBusObjectPath defSinkPath = m_audioInter->defaultSink();
    m_defSinkInter = new DBusSink("com.deepin.daemon.Audio", defSinkPath.path(), QDBusConnection::sessionBus(), this);

    connect(m_defSinkInter, &DBusSink::VolumeChanged, this, &SoundApplet::onVolumeChanged);
    connect(m_defSinkInter, &DBusSink::MuteChanged, this, [ = ] {
        onVolumeChanged(m_defSinkInter->volume());
    });

    QString portId = m_defSinkInter->activePort().name;
    uint cardId = m_defSinkInter->card();
    //最后一个设备会被移除，但是当在控制中心选中此设备后需要添加，并勾选
    if (!m_lastPort.isNull() && m_lastPort->cardId() == cardId && m_lastPort->id() == portId) {
        startAddPort(m_lastPort);
    }
    activePort(portId,cardId);

    //无声卡状态下，会有伪sink设备，显示音量为0
    onVolumeChanged(findPort(portId, cardId) != nullptr ? m_defSinkInter->volume() : 0);
    emit defaultSinkChanged(m_defSinkInter);
}

void SoundApplet::onVolumeChanged(double volume)
{
    m_volumeSlider->setValue(std::min(150, qRound(volume * 100.0)));
    m_soundShow->setText(QString::number(volume * 100) + '%');
    emit volumeChanged(m_volumeSlider->value());
    refreshIcon();
}

void SoundApplet::volumeSliderValueChanged()
{
    m_defSinkInter->SetVolume(m_volumeSlider->value() / 100.0f, true);
    if (m_defSinkInter->mute())
        m_defSinkInter->SetMuteQueued(false);
}

void SoundApplet::cardsChanged(const QString &cards)
{
    QMap<uint, QStringList> tmpCardIds;

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
            if (portAvai == 2 || portAvai == 0 ) { // 0 Unknow 1 Not available 2 Available
                const QString portId = jPort["Name"].toString();
                const QString portName = jPort["Description"].toString();

                Port *port = findPort(portId, cardId);
                const bool include = port != nullptr;
                if (!include) { port = new Port(m_model); }

                port->setId(portId);
                port->setName(portName);
                port->setDirection(Port::Direction(jPort["Direction"].toDouble()));
                port->setCardId(cardId);
                port->setCardName(cardName);

                if (!include) {
                    startAddPort(port);
                }

                tmpPorts << portId;
            }
        }
        tmpCardIds.insert(cardId, tmpPorts);
    }

    onDefaultSinkChanged();//重新获取切换的设备信息

    // 判断是否存在激活的输出设备
    enableDevice(existActiveOutputDevice());

    for (Port *port : m_ports) {
        //只要有一个设备在控制中心被禁用后，在任务栏声音设备列表中该设备会被移除，
        if (!m_audioInter->IsPortEnabled(port->cardId(), port->id())) {
            removeDisabledDevice(port->id(), port->cardId());
        }
        //判断端口是否在最新的设备列表中
        if (tmpCardIds.contains(port->cardId())) {
            if (!tmpCardIds[port->cardId()].contains(port->id())) {
                startRemovePort(port->id(), port->cardId());
            }
        }
        else {
            startRemovePort(port->id(), port->cardId());
        }
    }
    //当只有一个设备剩余时，该设备也需要移除
    removeLastDevice();
    updateListHeight();
}

void SoundApplet::increaseVolumeChanged()
{
    // 触发VolumeSlider::valueChanged
    m_volumeSlider->setMaximum(m_audioInter->maxUIVolume() * 100.0f);

    // onDefaultSinkChanged()比increaseVolumeChanged()先执行，设置完最大值后需要重新设置当前值
    if (m_defSinkInter) {
        //　规避因开启声音增强，声音值超过100后，关闭声音增强后读取到数据为旧数据的问题
        if (m_audioInter->increaseVolume())
            m_volumeSlider->setValue(std::min(150, qRound(m_defSinkInter->volume() * 100.0)));
    }
}

void SoundApplet::refreshIcon()
{
    if (!m_defSinkInter)
        return;

    const bool mute = existActiveOutputDevice() ? m_defSinkInter->mute() : true;

    QString volumeString;

    if (mute) {
        volumeString = "muted";
    } else {
        volumeString = "off";
    }

    QString iconLeft = QString("audio-volume-%1-symbolic").arg(volumeString);
    QString iconRight = QString("audio-volume-high-symbolic");

    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        iconLeft.append("-dark");
        iconRight.append("-dark");
    }

    const auto ratio = devicePixelRatioF();
    QPixmap ret = ImageUtil::loadSvg(iconRight, ":/", ICON_SIZE, ratio);
    m_volumeIconMax->setPixmap(ret);

    ret = ImageUtil::loadSvg(iconLeft, ":/", ICON_SIZE, ratio);
    m_volumeIconMin->setPixmap(ret);
}

/**
 * @brief SoundApplet::startAddPort 添加端口前判断
 * @param port 端口
 */
void SoundApplet::startAddPort(Port *port)
{
    if (!containsPort(port) && port->direction() == Port::Out) {
        m_ports.append(port);
        addPort(port);
    }
}

/**
 * @brief SoundApplet::startRemovePort 移除端口前判断
 * @param portId 端口
 * @param cardId 声卡
 */
void SoundApplet::startRemovePort(const QString &portId, const uint &cardId)
{
    Port *port = findPort(portId, cardId);
    if (port) {
        m_ports.removeOne(port);
        port->deleteLater();
        removePort(portId, cardId);
    }
}

bool SoundApplet::containsPort(const Port *port)
{
    return findPort(port->id(), port->cardId()) != nullptr;
}

Port *SoundApplet::findPort(const QString &portId, const uint &cardId) const
{
    auto it = std::find_if(m_ports.begin(), m_ports.end(), [ = ] (Port *p) {
        return (p->id() == portId && p->cardId() == cardId);
    });

    if (it != m_ports.end()) {
        return *it;
    }

    return nullptr;
}

void SoundApplet::addPort(const Port *port)
{
    DStandardItem *pi = new DStandardItem;
    QString deviceName = port->name() + "(" + port->cardName() + ")";
    pi->setText(deviceName);
    pi->setTextColorRole(QPalette::BrightText);
    pi->setData(QVariant::fromValue<const Port *>(port), Qt::WhatsThisPropertyRole);

    connect(port, &Port::nameChanged, this, [ = ](const QString &str) {
        QString devName = str + "(" + port->cardName() + ")";
        pi->setText(devName);
    });
    connect(port, &Port::cardNameChanged, this, [ = ](const QString &str) {
        QString devName = port->name() + "(" + str + ")";
        pi->setText(devName);
    });
    connect(port, &Port::isActiveChanged, this, [ = ](bool isActive) {
        pi->setCheckState(isActive ? Qt::CheckState::Checked : Qt::CheckState::Unchecked);
    });

    if (port->isActive()) {
        pi->setCheckState(Qt::CheckState::Checked);
    }

    m_model->appendRow(pi);
    m_model->sort(0);
    m_secondSeperator->setVisible(m_model->rowCount() > 1);
    updateListHeight();
}

void SoundApplet::removePort(const QString &portId, const uint &cardId)
{
    auto rmFunc = [ = ](QStandardItemModel * model) {
        for (int i = 0; i < model->rowCount();) {
            auto item = model->item(i);
            auto port = item->data(Qt::WhatsThisPropertyRole).value<const Port *>();
            if (port->id() == portId && cardId == port->cardId()) {
                model->removeRow(i);
                break;
            } else {
                ++i;
            }
        }
    };

    rmFunc(m_model);
    m_secondSeperator->setVisible(m_model->rowCount() > 1);
    updateListHeight();
}

/**
 * @brief SoundApplet::activePort 激活某一指定端口
 * @param portId 端口
 * @param cardId 声卡
 */
void SoundApplet::activePort(const QString &portId, const uint &cardId)
{
    for (Port *it : m_ports) {
        if (it->id() == portId && it->cardId() == cardId) {
            it->setIsActive(true);
            enableDevice(true);
        }
        else {
            it->setIsActive(false);
        }
    }
}

void SoundApplet::updateCradsInfo()
{
    QString info = m_audioInter->property("CardsWithoutUnavailable").toString();
    if(m_deviceInfo != info){
        cardsChanged(info);
        m_deviceInfo = info;
    }
}

void SoundApplet::enableDevice(bool flag)
{
    QString status = m_gsettings ? m_gsettings->get(GSETTING_SOUND_OUTPUT_SLIDER).toString() : "Enabled";
    if ("Disabled" == status ) {
        m_volumeSlider->setEnabled(false);
    } else if ("Enabled" == status) {
        m_volumeSlider->setEnabled(flag);
    }
    m_volumeIconMin->setEnabled(flag);
    m_soundShow->setEnabled(flag);
    m_volumeIconMax->setEnabled(flag);
    m_deviceLabel->setEnabled(flag);
}

void SoundApplet::disableAllDevice()
{
    for (Port *port : m_ports) {
        port->setIsActive(false);
    }
}

/**
 * @brief SoundApplet::removeLastDevice
 * 移除最后一个设备
 */
void SoundApplet::removeLastDevice()
{
    if (m_ports.count() == 1 && m_ports.at(0)) {
        m_lastPort = new Port(m_model);
        m_lastPort->setId(m_ports.at(0)->id());
        m_lastPort->setName(m_ports.at(0)->name());
        m_lastPort->setDirection(m_ports.at(0)->direction());
        m_lastPort->setCardId(m_ports.at(0)->cardId());
        m_lastPort->setCardName(m_ports.at(0)->cardName());
        startRemovePort(m_ports.at(0)->id(), m_ports.at(0)->cardId());
        qDebug() << "remove last output device";
    }
}

/**
 * @brief SoundApplet::removeDisabledDevice 移除禁用设备
 * @param portId
 * @param cardId
 */
void SoundApplet::removeDisabledDevice(QString portId, unsigned int cardId)
{
    startRemovePort(portId, cardId);
    if (m_defSinkInter->activePort().name == portId && m_defSinkInter->card() == cardId) {
        enableDevice(false);
        disableAllDevice();
    }
    qDebug() << "remove disabled output device";
}

void SoundApplet::updateVolumeSliderStatus(const QString &status)
{
    bool flag = true;
    if ("Enabled" == status) {
        flag = true;
    } else if ("Disabled" == status) {
        flag = false;
    }
    m_volumeSlider->setEnabled(flag);
    m_volumeIconMin->setEnabled(flag);
    m_volumeIconMax->setEnabled(flag);

    flag = "Hiden" != status;
    m_volumeSlider->setVisible(flag);
    m_volumeIconMin->setVisible(flag);
    m_volumeIconMax->setVisible(flag);
}

/** 判断是否存在未禁用的声音输出设备
 * @brief SoundApplet::existActiveOutputDevice
 * @return 存在返回true,否则返回false
 */
bool SoundApplet::existActiveOutputDevice()
{
    QString info = m_audioInter->property("CardsWithoutUnavailable").toString();

    QJsonDocument doc = QJsonDocument::fromJson(info.toUtf8());
    QJsonArray jCards = doc.array();
    for (QJsonValue cV : jCards) {
        QJsonObject jCard = cV.toObject();
        QJsonArray jPorts = jCard["Ports"].toArray();

        for (QJsonValue pV : jPorts) {
            QJsonObject jPort = pV.toObject();
            if (jPort["Direction"].toInt() == 1 && jPort["Enabled"].toBool())
                return true;
        }
    }

    return false;
}

bool SoundApplet::eventFilter(QObject *watcher, QEvent *event)
{
    // 当控制中心禁用所有输出设备时，静音按钮置灰，其他情况正常．
    if (watcher == m_volumeIconMin && event->type() == QEvent::MouseButtonRelease) {
        m_defSinkInter->SetMuteQueued(!m_defSinkInter->mute());
    }
    return false;
}

void SoundApplet::haldleDbusSignal(const QDBusMessage &msg)
{
    Q_UNUSED(msg)

    updateCradsInfo();
}

void SoundApplet::updateListHeight()
{
    //设备数多于10个时显示滚动条,固定高度
    int count = m_model->rowCount() == 1 ? 0 : m_model->rowCount();

    if (m_model->rowCount() > 10) {
        count = 10;
        m_listView->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOn);
    } else {
        m_listView->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    }

    int visualHeight = 0;
    for (int i = 0; i < count; i++)
        visualHeight += m_listView->visualRect(m_model->index(i, 0)).height();

    int listMargin = m_listView->contentsMargins().top() + m_listView->contentsMargins().bottom();
    //显示声音设备列表高度 = 设备的高度 + 间隔 + 边距
    int viewHeight = visualHeight + m_listView->spacing() * (count - 1) + listMargin;
    // 设备信息高度 = 设备标签 + 分隔线 + 滚动条 + 间隔
    int labelHeight = m_deviceLabel->height() > m_soundShow->height() ? m_deviceLabel->height() : m_soundShow->height();
    int infoHeight = labelHeight + m_seperator->height() * 2 + m_volumeSlider->height();
    int margain = m_centralLayout->contentsMargins().top() + m_centralLayout->contentsMargins().bottom();
    //整个界面高度 = 显示声音设备列表高度 + 设备信息高度 + 边距
    int totalHeight = viewHeight + infoHeight + margain;
    //加上分割线占用的高度，否则显示界面高度不够显示，会造成音频列表item最后一项比其它项的高度小
    m_listView->setFixedHeight(viewHeight);
    setFixedHeight(totalHeight);
    m_centralWidget->setFixedHeight(totalHeight);
}

void SoundApplet::portEnableChange(unsigned int cardId, QString portId)
{
    Q_UNUSED(cardId)
    Q_UNUSED(portId)
    m_deviceInfo = "";
    updateCradsInfo();
}
