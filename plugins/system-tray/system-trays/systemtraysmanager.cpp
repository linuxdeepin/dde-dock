#include "systemtraysmanager.h"
#include "sound/soundtrayloader.h"
#include "shutdown/shutdowntrayloader.h"
#include "shutdown/powertrayloader.h"
#include "network/networktrayloader.h"

SystemTraysManager::SystemTraysManager(QObject *parent)
    : QObject(parent)
{
    AbstractTrayLoader *soundLoader = new SoundTrayLoader(this);
    AbstractTrayLoader *shutdownLoader = new ShutdownTrayLoader(this);
    AbstractTrayLoader *powerLoader = new PowerTrayLoader(this);
    AbstractTrayLoader *networkLoader = new NetworkTrayLoader(this);

    m_loaderList.append(soundLoader);
    m_loaderList.append(shutdownLoader);
    m_loaderList.append(powerLoader);
    m_loaderList.append(networkLoader);

    connect(soundLoader, &AbstractTrayLoader::systemTrayAdded, this, &SystemTraysManager::systemTrayWidgetAdded);
    connect(soundLoader, &AbstractTrayLoader::systemTrayRemoved, this, &SystemTraysManager::systemTrayWidgetRemoved);
    connect(shutdownLoader, &AbstractTrayLoader::systemTrayAdded, this, &SystemTraysManager::systemTrayWidgetAdded);
    connect(shutdownLoader, &AbstractTrayLoader::systemTrayRemoved, this, &SystemTraysManager::systemTrayWidgetRemoved);
    connect(powerLoader, &AbstractTrayLoader::systemTrayAdded, this, &SystemTraysManager::systemTrayWidgetAdded);
    connect(powerLoader, &AbstractTrayLoader::systemTrayRemoved, this, &SystemTraysManager::systemTrayWidgetRemoved);
    connect(networkLoader, &AbstractTrayLoader::systemTrayAdded, this, &SystemTraysManager::systemTrayWidgetAdded);
    connect(networkLoader, &AbstractTrayLoader::systemTrayRemoved, this, &SystemTraysManager::systemTrayWidgetRemoved);
}

void SystemTraysManager::startLoad()
{
    for (auto loader : m_loaderList) {
        if (loader->waitService() && !loader->serviceExist()) {
            loader->waitServiceForLoad();
            continue;
        }
        loader->load();
    }
}
