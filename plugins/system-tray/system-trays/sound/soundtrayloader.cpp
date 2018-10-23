#include "soundtrayloader.h"
#include "soundtraywidget.h"

#define SoundItemKey "system-tray-sound"
#define SoundService "com.deepin.daemon.Audio"

SoundTrayLoader::SoundTrayLoader(QObject *parent) : AbstractTrayLoader(SoundService, parent)
{
}

void SoundTrayLoader::load()
{
    emit systemTrayAdded(SoundItemKey, new SoundTrayWidget);
}
