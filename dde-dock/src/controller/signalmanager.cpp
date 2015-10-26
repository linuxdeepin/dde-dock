#include "signalmanager.h"

SignalManager *SignalManager::m_signalManager = NULL;
SignalManager *SignalManager::instance()
{
    if (!m_signalManager)
        m_signalManager = new SignalManager;
    return m_signalManager;
}

SignalManager::SignalManager(QObject *parent) : QObject(parent)
{

}

