#ifndef SOUNDAPPLET_H
#define SOUNDAPPLET_H

#include "dbus/dbusaudio.h"
#include "dbus/dbussink.h"

#include <QScrollArea>
#include <QVBoxLayout>

class SoundApplet : public QScrollArea
{
    Q_OBJECT

public:
    explicit SoundApplet(QWidget *parent = 0);

signals:
    void defaultSinkChanged(DBusSink *sink) const;

private slots:
    void defaultSinkChanged();

private:
    QWidget *m_centeralWidget;
    QVBoxLayout *m_centeralLayout;

    DBusAudio *m_audioInter;
    DBusSink *m_defSinkInter;
};

#endif // SOUNDAPPLET_H
