#include "wirelessapplet.h"

WirelessApplet::WirelessApplet(const QString &devicePath, QWidget *parent)
    : QScrollArea(parent),
      m_devicePath(devicePath)
{
    setFixedHeight(300);

    setFrameStyle(QFrame::NoFrame);
    setFixedWidth(300);
    setStyleSheet("background-color:transparent;");
}
