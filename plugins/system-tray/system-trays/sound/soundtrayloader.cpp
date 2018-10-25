/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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
