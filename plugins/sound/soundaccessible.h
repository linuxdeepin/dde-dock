#ifndef SOUNDACCESSIBLE_H
#define SOUNDACCESSIBLE_H
#include "../frame/window/accessibledefine.h"

#include "sounditem.h"
#include "soundapplet.h"
#include "sinkinputwidget.h"
#include "./componments/volumeslider.h"
#include "./componments/horizontalseparator.h"

SET_BUTTON_ACCESSIBLE(SoundItem, "plugin-sounditem")
SET_FORM_ACCESSIBLE(SoundApplet, "soundapplet")
SET_FORM_ACCESSIBLE(SinkInputWidget, "sinkinputwidget")
SET_SLIDER_ACCESSIBLE(VolumeSlider, "volumeslider")
SET_FORM_ACCESSIBLE(HorizontalSeparator, "horizontalseparator")

QAccessibleInterface *soundAccessibleFactory(const QString &classname, QObject *object)
{
    QAccessibleInterface *interface = nullptr;

    USE_ACCESSIBLE(classname, SoundItem);
    USE_ACCESSIBLE(classname, SoundApplet);
    USE_ACCESSIBLE(classname, SinkInputWidget);
    USE_ACCESSIBLE(classname, VolumeSlider);
    USE_ACCESSIBLE(classname, HorizontalSeparator);

    return interface;
}

#endif // SOUNDACCESSIBLE_H
