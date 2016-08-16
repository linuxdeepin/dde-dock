#include "sounditem.h"
#include "constants.h"

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>

SoundItem::SoundItem(QWidget *parent)
    : QWidget(parent),

      m_applet(new SoundApplet(this)),
      m_sinkInter(nullptr)
{
    QIcon::setThemeName("deepin");

    m_applet->setVisible(false);

    connect(m_applet, static_cast<void (SoundApplet::*)(DBusSink*) const>(&SoundApplet::defaultSinkChanged), this, &SoundItem::sinkChanged);
}

QWidget *SoundItem::popupApplet()
{
    return m_applet;
}

QSize SoundItem::sizeHint() const
{
    return QSize(26, 26);
}

void SoundItem::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    refershIcon();
}

void SoundItem::mousePressEvent(QMouseEvent *e)
{
    if (e->button() != Qt::RightButton)
        return QWidget::mousePressEvent(e);

    const QPoint p(e->pos() - rect().center());
    if (p.manhattanLength() < std::min(width(), height()) * 0.8 * 0.5)
        return;

    return QWidget::mousePressEvent(e);
}

void SoundItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_iconPixmap.rect().center(), m_iconPixmap);
}

void SoundItem::refershIcon()
{
    if (!m_sinkInter)
        return;

    const double volmue = m_sinkInter->volume();
    const bool mute = m_sinkInter->mute() || volmue < 0.001;
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    QString iconString;
    if (displayMode == Dock::Fashion)
    {
        QString volumeString;
        if (volmue >= 1.0)
            volumeString = "100";
        else
            volumeString = QString("0") + ('0' + int(volmue * 10)) + "0";

        iconString = "audio-volume-" + volumeString + (mute ? "-muted" : "");
    } else {
        QString volumeString;
        if (mute)
            volumeString = "muted";
        else if (volmue >= double(2)/3)
            volumeString = "high";
        else if (volmue >= double(1)/3)
            volumeString = "medium";
        else
            volumeString = "low";

        iconString = QString("audio-volume-%1-symbolic").arg(volumeString);
    }

    const int iconSize = displayMode == Dock::Fashion ? std::min(width(), height()) * 0.8 : 16;
    const QIcon icon = QIcon::fromTheme(iconString);
    m_iconPixmap = icon.pixmap(iconSize, iconSize);

    update();
}

void SoundItem::sinkChanged(DBusSink *sink)
{
    m_sinkInter = sink;

    connect(m_sinkInter, &DBusSink::MuteChanged, this, &SoundItem::refershIcon);
    connect(m_sinkInter, &DBusSink::VolumeChanged, this, &SoundItem::refershIcon);
    refershIcon();
}
