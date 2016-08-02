#include "soundapplet.h"

#define WIDTH       200

SoundApplet::SoundApplet(QWidget *parent)
    : QScrollArea(parent),

      m_centeralWidget(new QWidget(this)),

      m_audioInter(new DBusAudio(this))
{
    m_centeralLayout = new QVBoxLayout;
    m_centeralWidget->setLayout(m_centeralLayout);
    m_centeralWidget->setFixedWidth(WIDTH);

    setFixedWidth(WIDTH);
    setWidget(m_centeralWidget);
    setFrameStyle(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setStyleSheet("background-color:transparent;");

    for (auto sp : m_audioInter->sinks())
    {
        DBusSink *sink = new DBusSink(sp.path(), this);
        qDebug() << sink->name();
    }
}
