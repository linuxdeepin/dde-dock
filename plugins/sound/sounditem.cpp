#include "sounditem.h"
#include "constants.h"

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>

SoundItem::SoundItem(QWidget *parent)
    : QWidget(parent),

      m_tipsLabel(new QLabel(this)),
      m_applet(new SoundApplet(this)),
      m_sinkInter(nullptr)
{
    QIcon::setThemeName("deepin");

    m_tipsLabel->setVisible(false);
    m_tipsLabel->setFixedWidth(60);
    m_tipsLabel->setAlignment(Qt::AlignCenter);
    m_tipsLabel->setStyleSheet("color:white;"
                               "padding:5px 10px;");

    m_applet->setVisible(false);

    connect(m_applet, static_cast<void (SoundApplet::*)(DBusSink*) const>(&SoundApplet::defaultSinkChanged), this, &SoundItem::sinkChanged);
    connect(m_applet, &SoundApplet::volumeChanegd, this, &SoundItem::refershTips);
}

QWidget *SoundItem::tipsWidget()
{
    refershTips(true);

    return m_tipsLabel;
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

void SoundItem::wheelEvent(QWheelEvent *e)
{
    QWheelEvent *event = new QWheelEvent(e->pos(), e->delta(), e->buttons(), e->modifiers());
    qApp->postEvent(m_applet->mainSlider(), event);

    e->accept();
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
    const bool mute = m_sinkInter->mute();
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

void SoundItem::refershTips(const bool force)
{
    if (!force && !m_tipsLabel->isVisible())
        return;

    const QString value = QString::number(m_applet->volumeValue() / 10) + '%';
    m_tipsLabel->setText(QString(tr("Current Volume %1").arg(value)));
}

void SoundItem::sinkChanged(DBusSink *sink)
{
    m_sinkInter = sink;

    connect(m_sinkInter, &DBusSink::MuteChanged, this, &SoundItem::refershIcon);
    connect(m_sinkInter, &DBusSink::VolumeChanged, this, &SoundItem::refershIcon);
    refershIcon();
}
