#ifndef SOUNDITEM_H
#define SOUNDITEM_H

#include "soundapplet.h"
#include "dbus/dbussink.h"

#include <QWidget>

class SoundItem : public QWidget
{
    Q_OBJECT

public:
    explicit SoundItem(QWidget *parent = 0);

    QWidget *popupApplet();

protected:
    QSize sizeHint() const;
    void resizeEvent(QResizeEvent *e);
    void paintEvent(QPaintEvent *e);

private slots:
    void refershIcon();
    void sinkChanged(DBusSink *sink);

private:
    SoundApplet *m_applet;
    DBusSink *m_sinkInter;
    QPixmap m_iconPixmap;
};

#endif // SOUNDITEM_H
