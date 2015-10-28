#ifndef SIGNALMANAGER_H
#define SIGNALMANAGER_H

#include <QObject>

class SignalManager : public QObject
{
    Q_OBJECT
public:
    static SignalManager *instance();

signals:
    void requestAppIconUpdate();

private:
    explicit SignalManager(QObject *parent = 0);
    static SignalManager *m_signalManager;
};

#endif // SIGNALMANAGER_H
