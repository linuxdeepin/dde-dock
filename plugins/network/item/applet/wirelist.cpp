#include "wirelist.h"

#include <QTimer>
#include <QDebug>

using namespace dde::network;

const int Width = 300;
const int ItemHeight = 30;

WireList::WireList(WiredDevice *device, QWidget *parent)
    : QScrollArea(parent)
    , m_device(device)
    , m_updateAPTimer(new QTimer(this))
    , m_deviceName(new QLabel(this))
    , m_switchBtn(new DSwitchButton(this))
{

    m_updateAPTimer->setSingleShot(true);
    m_updateAPTimer->setInterval(100);

    auto centralWidget = new QWidget(this);
    centralWidget->setFixedWidth(Width);
    m_centralLayout = new QVBoxLayout;
    m_centralLayout->setMargin(0);
    m_centralLayout->setSpacing(0);
    centralWidget->setLayout(m_centralLayout);

    auto controlPanel = new QWidget(this);
    controlPanel->setFixedWidth(Width);
    controlPanel->setFixedHeight(ItemHeight);
    auto controlPanelLayout = new QHBoxLayout;
    controlPanelLayout->setMargin(0);
    controlPanelLayout->setSpacing(0);
    controlPanelLayout->addSpacing(5);
    controlPanelLayout->addWidget(m_deviceName);
    controlPanelLayout->addStretch();
    controlPanelLayout->addWidget(m_switchBtn);
    controlPanelLayout->addSpacing(5);
    controlPanel->setLayout(controlPanelLayout);

    m_centralLayout->addWidget(controlPanel);

    setWidget(centralWidget);
    setFixedWidth(Width);
    setFrameShape(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);

    connect(m_device, &WiredDevice::connectionsChanged, this, &WireList::changeConnections);
    connect(m_device, &WiredDevice::activeWiredConnectionInfoChanged, this, &WireList::changeActiveWiredConnectionInfo);
    connect(m_device, &WiredDevice::activeConnectionsChanged, this, &WireList::changeActiveConnections);
    connect(m_device, &WiredDevice::activeConnectionsInfoChanged, this, &WireList::changeActiveConnectionsInfo);

    connect(m_switchBtn, &DSwitchButton::checkedChanged, this, &WireList::deviceEnabled);

    connect(m_updateAPTimer, &QTimer::timeout, this, &WireList::updateConnectionList);

    QMetaObject::invokeMethod(this, "loadConnectionList", Qt::QueuedConnection);
}

void WireList::changeConnections(const QList<QJsonObject> &connections)
{
    for (auto object : connections) {
        qDebug() << object;
    }
}

void WireList::changeActiveWiredConnectionInfo(const QJsonObject &connInfo)
{
    for (auto object : connInfo) {
        qDebug() << object;
    }
}

void WireList::changeActiveConnections(const QList<QJsonObject> &activeConns)
{
    for (auto object : activeConns) {
        qDebug() << object;
    }
}

void WireList::changeActiveConnectionsInfo(const QList<QJsonObject> &activeConnInfoList)
{
    for (auto object : activeConnInfoList) {
        qDebug() << object;
    }
}
