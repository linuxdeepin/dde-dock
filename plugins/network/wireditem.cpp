#include "constants.h"
#include "wireditem.h"
#include "util/imageutil.h"

#include <QPainter>

WiredItem::WiredItem(QWidget *parent)
    : QWidget(parent),

      m_networkManager(NetworkManager::instance(this))
{

}

void WiredItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_icon.rect().center(), m_icon);
}

void WiredItem::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    reloadIcon();
}

QSize WiredItem::sizeHint() const
{
    return QSize(24, 24);
}

void WiredItem::reloadIcon()
{
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const bool connect = m_networkManager->states().testFlag(NetworkManager::WiredConnection);

    if (displayMode == Dock::Fashion)
    {
        const int size = std::min(width(), height()) * 0.8;

        if (connect)
            m_icon = ImageUtil::loadSvg(":/wired/resources/wired/wired-connected.svg", size);
        else
            m_icon = ImageUtil::loadSvg(":/wired/resources/wired/wired-disconnected.svg", size);
    } else {
        if (connect)
            m_icon = ImageUtil::loadSvg(":/wired/resources/wired/wired-connected-small.svg", 16);
        else
            m_icon = ImageUtil::loadSvg(":/wired/resources/wired/wired-disconnected-small.svg", 16);
    }
}
