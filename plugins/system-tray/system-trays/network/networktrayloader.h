#ifndef NETWORKTRAYLOADER_H
#define NETWORKTRAYLOADER_H

#include "../abstracttrayloader.h"
#include "item/abstractnetworktraywidget.h"

#include <QObject>

#include <NetworkWorker>
#include <NetworkModel>

class NetworkTrayLoader : public AbstractTrayLoader
{
    Q_OBJECT
public:
    explicit NetworkTrayLoader(QObject *parent = nullptr);

public Q_SLOTS:
    void load() Q_DECL_OVERRIDE;

private:
    AbstractNetworkTrayWidget *trayWidgetByPath(const QString &path);

private Q_SLOTS:
    void onDeviceListChanged(const QList<dde::network::NetworkDevice *> devices);
    void refreshWiredItemVisible();

private:
    dde::network::NetworkModel *m_networkModel;
    dde::network::NetworkWorker *m_networkWorker;

    QMap<QString, AbstractNetworkTrayWidget *> m_trayWidgetsMap;
    QTimer *m_delayRefreshTimer;
};

#endif // NETWORKTRAYLOADER_H
