#include "soundapplet.h"

#define WIDTH       200

SoundApplet::SoundApplet(QWidget *parent)
    : QScrollArea(parent),

      m_centeralWidget(new QWidget(this))
{
    m_centeralLayout = new QVBoxLayout;
    m_centeralWidget->setLayout(m_centeralLayout);

    setFixedWidth(WIDTH);
    setFrameStyle(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setStyleSheet("background-color:transparent;");
}
