#ifndef TRAYMONITOR_H
#define TRAYMONITOR_H

#include <QObject>

#include "dbustraymanager.h"
#include "statusnotifierwatcher_interface.h"

using namespace org::kde;
class TrayMonitor : public QObject
{
    Q_OBJECT

public:
    explicit TrayMonitor(QObject *parent = nullptr);

public Q_SLOTS:
    void onTrayIconsChanged();
    void onSniItemsChanged();

    void startLoadIndicators();

Q_SIGNALS:
    void requestUpdateIcon(quint32);
    void xEmbedTrayAdded(quint32);
    void xEmbedTrayRemoved(quint32);

    void sniTrayAdded(const QString &);
    void sniTrayRemoved(const QString &);

    void indicatorFounded(const QString &);

private:
    DBusTrayManager *m_trayInter;
    StatusNotifierWatcher *m_sniWatcher;

    QList<quint32> m_trayWids;
    QStringList m_sniServices;
};

#endif // TRAYMONITOR_H
