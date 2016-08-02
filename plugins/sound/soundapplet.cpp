#include "soundapplet.h"

#define WIDTH       200

SoundApplet::SoundApplet(QWidget *parent)
    : QScrollArea(parent),

      m_centeralWidget(new QWidget(this)),

      m_audioInter(new DBusAudio(this)),
      m_defSinkInter(nullptr)
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

    QMetaObject::invokeMethod(this, "defaultSinkChanged", Qt::QueuedConnection);
}

void SoundApplet::defaultSinkChanged()
{
    delete m_defSinkInter;

    const QDBusObjectPath defSinkPath = m_audioInter->GetDefaultSink();
    m_defSinkInter = new DBusSink(defSinkPath.path(), this);

    emit defaultSinkChanged(m_defSinkInter);
}
