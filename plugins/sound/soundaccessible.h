// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SOUNDACCESSIBLE_H
#define SOUNDACCESSIBLE_H
#include "accessibledefine.h"

#include "sounditem.h"
#include "soundapplet.h"
#include "./componments/volumeslider.h"

SET_BUTTON_ACCESSIBLE(SoundItem, "plugin-sounditem")
SET_FORM_ACCESSIBLE(SoundApplet, "soundapplet")
SET_SLIDER_ACCESSIBLE(VolumeSlider, "volumeslider")

QAccessibleInterface *soundAccessibleFactory(const QString &classname, QObject *object)
{
    QAccessibleInterface *interface = nullptr;

    USE_ACCESSIBLE(classname, SoundItem)
    ELSE_USE_ACCESSIBLE(classname, SoundApplet)
    ELSE_USE_ACCESSIBLE(classname, VolumeSlider);

    return interface;
}

#endif // SOUNDACCESSIBLE_H
