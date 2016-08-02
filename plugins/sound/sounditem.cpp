#include "sounditem.h"

#include <QPainter>
#include <QIcon>

SoundItem::SoundItem(QWidget *parent)
    : QWidget(parent),

      m_applet(new SoundApplet(this))
{
    QIcon::setThemeName("deepin");

    m_applet->setVisible(false);
}

QWidget *SoundItem::popupApplet()
{
    return m_applet;
}

QSize SoundItem::sizeHint() const
{
    return QSize(24, 24);
}

void SoundItem::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    refershIcon();
}

void SoundItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_iconPixmap.rect().center(), m_iconPixmap);
}

void SoundItem::refershIcon()
{
    const int iconSize = std::min(width(), height()) * 0.8;
    const QIcon icon = QIcon::fromTheme("audio-volume-080");
    m_iconPixmap = icon.pixmap(iconSize, iconSize);

    update();
}
