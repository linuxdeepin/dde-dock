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

#include "wirelessapplet.h"
#include "accesspointwidget.h"

#include <QJsonDocument>
#include <QScreen>
#include <QDebug>
#include <QGuiApplication>

#include <dinputdialog.h>
#include <QScrollBar>

DWIDGET_USE_NAMESPACE

#define WIDTH           300
#define MAX_HEIGHT      300
#define ITEM_HEIGHT     30

WirelessList::WirelessList(const QSet<NetworkDevice>::const_iterator &deviceIter, QWidget *parent)
    : QScrollArea(parent),

      m_device(*deviceIter),
      m_activeAP(),

      m_updateAPTimer(new QTimer(this)),
      m_pwdDialog(new DInputDialog(nullptr)),
      m_autoConnBox(new QCheckBox),

      m_centralLayout(new QVBoxLayout),
      m_centralWidget(new QWidget),
      m_controlPanel(new DeviceControlWidget),
      m_networkInter(new DBusNetwork(this))
{
    setFixedHeight(WIDTH);

    m_currentClickAPW = nullptr;

    m_autoConnBox->setText(tr("Auto-connect"));

    const auto ratio = qApp->devicePixelRatio();
    QPixmap iconPix = QIcon::fromTheme("notification-network-wireless-full").pixmap(QSize(48, 48) * ratio);
    iconPix.setDevicePixelRatio(ratio);

    m_pwdDialog->setTextEchoMode(QLineEdit::Password);
    m_pwdDialog->setWindowFlags(Qt::WindowStaysOnTopHint | Qt::FramelessWindowHint | Qt::Dialog);
    m_pwdDialog->setTextEchoMode(DLineEdit::Password);
    m_pwdDialog->setIcon(iconPix);
    m_pwdDialog->addSpacing(10);
    m_pwdDialog->addContent(m_autoConnBox, Qt::AlignLeft);
    m_pwdDialog->setOkButtonText(tr("Connect"));
    m_pwdDialog->setCancelButtonText(tr("Cancel"));

    m_updateAPTimer->setSingleShot(true);
    m_updateAPTimer->setInterval(100);

    m_centralWidget->setFixedWidth(WIDTH);
    m_centralWidget->setLayout(m_centralLayout);

//    m_centralLayout->addWidget(m_controlPanel);
    m_centralLayout->setSpacing(0);
    m_centralLayout->setMargin(0);

    // initialization state.
    m_deviceEnabled = m_networkInter->IsDeviceEnabled(m_device.dbusPath());

    setWidget(m_centralWidget);
    setFrameStyle(QFrame::NoFrame);
    setFixedWidth(300);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setStyleSheet("background-color:transparent;");

    m_indicator = new DPictureSequenceView(this);
    m_indicator->setPictureSequence(":/wireless/indicator/resources/wireless/spinner14/Spinner%1.png", QPair<int, int>(1, 91), 2);
    m_indicator->setFixedSize(QSize(14, 14) * ratio);
    m_indicator->setVisible(false);

    connect(m_networkInter, &DBusNetwork::AccessPointAdded, this, &WirelessList::APAdded);
    connect(m_networkInter, &DBusNetwork::AccessPointRemoved, this, &WirelessList::APRemoved);
    connect(m_networkInter, &DBusNetwork::AccessPointPropertiesChanged, this, &WirelessList::APPropertiesChanged);
    connect(m_networkInter, &DBusNetwork::DevicesChanged, this, &WirelessList::deviceStateChanged);
    connect(m_networkInter, &DBusNetwork::NeedSecrets, this, &WirelessList::needSecrets);
    connect(m_networkInter, &DBusNetwork::DeviceEnabled, this, &WirelessList::deviceEnabled);

    connect(m_controlPanel, &DeviceControlWidget::deviceEnableChanged, this, &WirelessList::deviceEnableChanged);
    connect(m_controlPanel, &DeviceControlWidget::requestRefresh, m_networkInter, &DBusNetwork::RequestWirelessScan);

    connect(m_updateAPTimer, &QTimer::timeout, this, &WirelessList::updateAPList);

    connect(this, &WirelessList::activeAPChanged, m_updateAPTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(this, &WirelessList::wirelessStateChanged, m_updateAPTimer, static_cast<void (QTimer::*)()>(&QTimer::start));

    connect(m_networkInter, &DBusNetwork::NeedSecretsFinished, m_pwdDialog, &DInputDialog::close);
    connect(m_pwdDialog, &DInputDialog::textValueChanged, this, &WirelessList::onPwdDialogTextChanged);
    connect(m_pwdDialog, &DInputDialog::cancelButtonClicked, this, &WirelessList::pwdDialogCanceled);
    connect(m_pwdDialog, &DInputDialog::accepted, this, &WirelessList::pwdDialogAccepted);

    connect(this->verticalScrollBar(), &QScrollBar::valueChanged, this, [=] {
        if (!m_currentClickAPW) return;

        const int h = -(m_currentClickAPW->height() - m_indicator->height()) / 2;
        m_indicator->move(m_currentClickAPW->mapTo(this, m_currentClickAPW->rect().topRight()) - QPoint(35, h));
    });

    QMetaObject::invokeMethod(this, "init", Qt::QueuedConnection);
}

WirelessList::~WirelessList()
{
    m_pwdDialog->deleteLater();
}

NetworkDevice::NetworkState WirelessList::wirelessState() const
{
    return m_device.state();
}

int WirelessList::activeAPStrgength() const
{
    return m_activeAP.strength();
}

QWidget *WirelessList::controlPanel()
{
    return m_controlPanel;
}

void WirelessList::init()
{
    loadAPList();
    onActiveAPChanged();
    deviceStateChanged();
}

void WirelessList::APAdded(const QString &devPath, const QString &info)
{
    if (devPath != m_device.path())
        return;

    AccessPoint ap(info);
    if (m_apList.contains(ap))
        return;

    m_apList.append(ap);
    m_updateAPTimer->start();
}

void WirelessList::APRemoved(const QString &devPath, const QString &info)
{
    if (devPath != m_device.path())
        return;

    AccessPoint ap(info);
    if (ap.ssid() == m_activeAP.ssid())
        return;

//    m_apList.removeOne(ap);
//    m_updateAPTimer->start();

    // NOTE: if one ap removed, prehaps another ap has same ssid, so we need to refersh ap list instead of remove it
    m_apList.clear();
    loadAPList();
}

void WirelessList::setDeviceInfo(const int index)
{
    // set device enable state
    m_controlPanel->setDeviceEnabled(m_deviceEnabled);

    // set device name
    if (index == -1)
        m_controlPanel->setDeviceName(tr("Wireless Network"));
    else
        m_controlPanel->setDeviceName(tr("Wireless Network %1").arg(index));
}

void WirelessList::loadAPList()
{
    const QString data = m_networkInter->GetAccessPoints(m_device.dbusPath());
    const QJsonDocument doc = QJsonDocument::fromJson(data.toUtf8());
    Q_ASSERT(doc.isArray());

    for (auto item : doc.array())
    {
        Q_ASSERT(item.isObject());

        AccessPoint ap(item.toObject());
        if (!m_apList.contains(ap))
            m_apList.append(ap);
    }

    m_updateAPTimer->start();
}

void WirelessList::APPropertiesChanged(const QString &devPath, const QString &info)
{
    if (devPath != m_device.path())
        return;

    QJsonDocument doc = QJsonDocument::fromJson(info.toUtf8());
    Q_ASSERT(doc.isObject());
    const AccessPoint ap(doc.object());

    auto it = std::find_if(m_apList.begin(), m_apList.end(),
                           [&] (const AccessPoint &a) {return a == ap;});

    if (it == m_apList.end())
        return;

    *it = ap;
    if (m_activeAP.path() == ap.path())
    {
        m_activeAP = ap;
        emit activeAPChanged();
    }

//    if (*it > ap)
//    {
//        *it = ap;
//        m_activeAP = ap;
//        m_updateAPTimer->start();

//        emit activeAPChanged();
//    }
}

void WirelessList::updateAPList()
{
    Q_ASSERT(sender() == m_updateAPTimer);

    int avaliableAPCount = 0;

    if (m_networkInter->IsDeviceEnabled(m_device.dbusPath()))
    {
        m_currentClickAPW = nullptr;
        // sort ap list by strength
        // std::sort(m_apList.begin(), m_apList.end(), std::greater<AccessPoint>());
        //        const bool wirelessActived = m_device.state() == NetworkDevice::Activated;

        // NOTE: Keep the amount consistent
        if(m_apList.size() > m_apwList.size()) {
            int i = m_apList.size() - m_apwList.size();
            for (int index = 0; index != i; index++) {
                AccessPointWidget *apw = new AccessPointWidget;
                apw->setFixedHeight(ITEM_HEIGHT);
                m_apwList << apw;
                m_centralLayout->addWidget(apw);

                connect(apw, &AccessPointWidget::requestActiveAP, this, &WirelessList::activateAP);
                connect(apw, &AccessPointWidget::requestDeactiveAP, this, &WirelessList::deactiveAP);
            }
        } else if (m_apList.size() < m_apwList.size()) {
            if (!m_apwList.isEmpty()) {
                int i = m_apwList.size() - m_apList.size();
                for (int index = 0; index != i; index++) {
                    AccessPointWidget *apw = m_apwList.last();
                    m_apwList.removeLast();
                    m_centralLayout->removeWidget(apw);
                    disconnect(apw, &AccessPointWidget::clicked, this, &WirelessList::updateIndicatorPos);
                    apw->deleteLater();
                }
            }
        }

        std::sort(m_apList.begin(), m_apList.end(), [&] (const AccessPoint &ap1, const AccessPoint &ap2) {
            if (ap1 == m_activeAP)
                return true;

            if (ap2 == m_activeAP)
                return false;

            return ap1.strength() > ap2.strength();
        });

        for (int i = 0; i != m_apList.size(); i++) {
            m_apwList[i]->updateAP(m_apList[i]);
            ++avaliableAPCount;
            connect(m_apwList[i], &AccessPointWidget::clicked, this, &WirelessList::updateIndicatorPos, Qt::UniqueConnection);
        }

        // update active AP state
        AccessPointWidget *apw = m_apwList.first();
        apw->setActiveState(m_device.state());

        // If the order of item changes
        if (m_apList.contains(m_currentClickAP) && m_indicator->isVisible()) {
            m_currentClickAPW = m_apwList.at(m_apList.indexOf(m_currentClickAP));
            const int h = -(m_currentClickAPW->height() - m_indicator->height()) / 2;
            m_indicator->move(m_currentClickAPW->mapTo(this, m_currentClickAPW->rect().topRight()) - QPoint(35, h));
        }

        if (m_device.state() == NetworkDevice::Activated ||
                m_device.state() == NetworkDevice::Failed ||
                m_device.state() == NetworkDevice::Unknow) {
            m_indicator->stop();
            m_indicator->hide();
        }
    }

    const int contentHeight = avaliableAPCount * ITEM_HEIGHT;
    m_centralWidget->setFixedHeight(contentHeight);
    setFixedHeight(std::min(contentHeight, MAX_HEIGHT));
}

void WirelessList::deviceEnableChanged(const bool enable)
{
    m_networkInter->EnableDevice(m_device.dbusPath(), enable);
    m_updateAPTimer->start();
}

void WirelessList::deviceStateChanged()
{
    const QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->devices().toUtf8());
    Q_ASSERT(doc.isObject());
    const QJsonObject obj = doc.object();

    for (auto infoList(obj.constBegin()); infoList != obj.constEnd(); ++infoList)
    {
        Q_ASSERT(infoList.value().isArray());

        if (infoList.key() != "wireless")
            continue;

        const auto list = infoList.value().toArray();
        for (auto i(0); i != list.size(); ++i)
        {
            const QJsonObject info = list[i].toObject();
            if (info.value("Path") == m_device.path())
            {
                const NetworkDevice prevInfo = m_device;
                m_device = NetworkDevice(NetworkDevice::Wireless, info);

                setDeviceInfo(list.size() == 1 ? -1 : i + 1);

                if (prevInfo.state() != m_device.state())
                    emit wirelessStateChanged(m_device.state());
                if (prevInfo.activeAp() != m_device.activeAp())
                    onActiveAPChanged();

                return;
            }
        }
    }
}

void WirelessList::onActiveAPChanged()
{
    const QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->GetAccessPoints(m_device.dbusPath()).value().toUtf8());
    Q_ASSERT(doc.isArray());

    for (auto dev : doc.array())
    {
        Q_ASSERT(dev.isObject());
        const QJsonObject obj = dev.toObject();

        if (obj.value("Path").toString() != m_device.activeAp())
            continue;

        m_activeAP = AccessPoint(obj);
        break;
    }

    emit activeAPChanged();
}

void WirelessList::pwdDialogAccepted()
{
    if (m_pwdDialog->textValue().isEmpty())
        return m_pwdDialog->setTextAlert(true);
    m_networkInter->FeedSecret(m_lastConnPath, m_lastConnSecurity, m_pwdDialog->textValue(), m_autoConnBox->isChecked());
}

void WirelessList::pwdDialogCanceled()
{
    m_networkInter->CancelSecret(m_lastConnPath, m_lastConnSecurity);
    m_pwdDialog->close();
}

void WirelessList::onPwdDialogTextChanged(const QString &text)
{
    m_pwdDialog->setTextAlert(false);

    do {
        if (text.isEmpty())
            break;
        const int len = text.length();

        // in wpa, password length must >= 8
        if (len < 8 && m_lastConnSecurityType.startsWith("wifi-wpa"))
            break;
        if (!(len == 5 || len == 13 || len == 16) && m_lastConnSecurityType.startsWith("wifi-wep"))
            break;

        return m_pwdDialog->setOkButtonEnabled(true);
    } while (false);

    m_pwdDialog->setOkButtonEnabled(false);
}

void WirelessList::deviceEnabled(const QString &devPath, const bool enable)
{
    if (devPath != m_device.path())
        return;

    if (m_deviceEnabled != enable) {
        m_deviceEnabled = enable;
        m_controlPanel->setDeviceEnabled(enable);
        m_updateAPTimer->start();
    }
}

void WirelessList::activateAP(const QDBusObjectPath &apPath, const QString &ssid)
{
    QString uuid;

    const QJsonDocument doc = QJsonDocument::fromJson(m_networkInter->connections().toUtf8());
    for (auto it : doc.object().value("wireless").toArray())
    {
        const QJsonObject obj = it.toObject();
        if (obj.value("Ssid").toString() != ssid)
            continue;
        if (obj.value("HwAddress").toString() != m_device.usingHwAddr())
            continue;

        uuid = obj.value("Uuid").toString();
        if (!uuid.isEmpty())
            break;
    }

    m_networkInter->ActivateAccessPoint(uuid, apPath, m_device.dbusPath()).waitForFinished();
}

void WirelessList::deactiveAP()
{
    m_activeAP = AccessPoint();
    m_networkInter->DisconnectDevice(QDBusObjectPath(m_device.path()));
}

void WirelessList::needSecrets(const QString &info)
{
    const QJsonObject infoObject = QJsonDocument::fromJson(info.toUtf8()).object();
    const QString connPath = infoObject.value("ConnectionPath").toString();
    const QString security = infoObject.value("SettingName").toString();
    const QString securityType = infoObject.value("KeyType").toString();
    const QString ssid = infoObject.value("ConnectionId").toString();
    const bool defaultAutoConnect = infoObject.value("AutoConnect").toBool();

    // check is our device' ap
    QString connHwAddr;
    QJsonDocument conns = QJsonDocument::fromJson(m_networkInter->connections().toUtf8());
    for (auto item : conns.object()["wireless"].toArray())
    {
        auto info = item.toObject();
        if (info["Path"].toString() != connPath)
            continue;
        connHwAddr = info["HwAddress"].toString();
        break;
    }
    if (connHwAddr != m_device.usingHwAddr())
        return;

    m_lastConnPath = connPath;
    m_lastConnSecurity = security;
    m_lastConnSecurityType = securityType;

    m_autoConnBox->setChecked(defaultAutoConnect);
    m_pwdDialog->setTitle(tr("Password required to connect to <font color=\"#faca57\">%1</font>").arg(ssid));

    // clear old config
    m_pwdDialog->setTextValue(QString());
    m_pwdDialog->setTextAlert(true);
    m_pwdDialog->setOkButtonEnabled(false);

    // check if controlcenter handle this request
//    QDBusInterface iface("com.deepin.dde.ControlCenter",
//                         "/com/deepin/dde/ControlCenter/Network",
//                         "com.deepin.dde.ControlCenter.Network");
//    if (iface.isValid() && iface.call("active").arguments().first().toBool())
//        return m_pwdDialog->hide();

    if (!m_pwdDialog->isVisible())
        m_pwdDialog->show();
}

void WirelessList::updateIndicatorPos()
{
    m_currentClickAPW = static_cast<AccessPointWidget*>(sender());

    if (m_currentClickAPW->active()) return;

    m_currentClickAP = m_currentClickAPW->ap();

    const int h = -(m_currentClickAPW->height() - m_indicator->height()) / 2;
    m_indicator->move(m_currentClickAPW->mapTo(this, m_currentClickAPW->rect().topRight()) - QPoint(35, h));
    m_indicator->show();
    m_indicator->play();
}
