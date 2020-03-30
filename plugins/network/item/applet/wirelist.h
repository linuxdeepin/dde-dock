#ifndef WIRELIST_H
#define WIRELIST_H

#include <QScrollArea>
#include <QPointer>
#include <QVBoxLayout>
#include <QLabel>

#include <WiredDevice>
#include <dpicturesequenceview.h>
#include <DSwitchButton>

DWIDGET_USE_NAMESPACE

class WireList : public QScrollArea
{
    Q_OBJECT
public:
    WireList(dde::network::WiredDevice *device, QWidget *parent = nullptr);

public slots:
    void changeConnections(const QList<QJsonObject> &connections);
    void changeActiveWiredConnectionInfo(const QJsonObject &connInfo);
    void changeActiveConnections(const QList<QJsonObject> &activeConns);
    void changeActiveConnectionsInfo(const QList<QJsonObject> &activeConnInfoList);
    void deviceEnabled(bool enabled);
    void updateConnectionList();

private slots:
    void loadConnectionList();

private:
    QPointer<dde::network::WiredDevice> m_device;

    QTimer *m_updateAPTimer;

    QLabel *m_deviceName;
    DSwitchButton *m_switchBtn;

    QVBoxLayout *m_centralLayout;
};

#endif // WIRELIST_H
