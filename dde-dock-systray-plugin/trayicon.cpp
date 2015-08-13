#include <QPainter>
#include <QBitmap>

#include "dockconstants.h"

#include "trayicon.h"

TrayIcon::TrayIcon(WId winId, QWidget *parent) :
    QWidget(parent)
{
    initItemMask();
    resize(Dock::APPLET_CLASSIC_ICON_SIZE,
           Dock::APPLET_CLASSIC_ICON_SIZE);

    m_win = QWindow::fromWinId(winId);

    QWidget * winItem = QWidget::createWindowContainer(m_win, this);
    winItem->resize(size());
}


void TrayIcon::initItemMask()
{
    m_itemMask = QPixmap(Dock::APPLET_CLASSIC_ICON_SIZE,
                         Dock::APPLET_CLASSIC_ICON_SIZE);
    m_itemMask.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&m_itemMask);
    // turn off  antialiasing.
    painter.setRenderHint(QPainter::Antialiasing, false);

    painter.setBrush(Qt::black);
    painter.drawRoundedRect(m_itemMask.rect(),
                            m_itemMask.width() / 2,
                            m_itemMask.height() / 2);

    painter.end();
}

void TrayIcon::maskOn()
{
    m_win->setMask(m_itemMask.mask());
}

void TrayIcon::maskOff()
{
    m_win->setMask(QRegion(0, 0, m_itemMask.width(), m_itemMask.height()));
}
