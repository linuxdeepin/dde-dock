#include "soundtraywidget.h"
#include "sounditem.h"
#include "constants.h"

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>
#include <QApplication>
#include <DApplication>
#include <DDBusSender>
#include "../widgets/tipswidget.h"

// menu actions
#define MUTE    "mute"
#define SETTINS "settings"

DWIDGET_USE_NAMESPACE

SoundTrayWidget::SoundTrayWidget(QWidget *parent)
    : AbstractSystemTrayWidget(parent),
      m_tipsLabel(new TipsWidget(this)),
      m_applet(new SoundApplet(this)),
      m_sinkInter(nullptr)
{
    m_tipsLabel->setObjectName("sound");
    m_tipsLabel->setVisible(false);

    m_applet->setVisible(false);

    connect(m_applet, static_cast<void (SoundApplet::*)(DBusSink*) const>(&SoundApplet::defaultSinkChanged), this, &SoundTrayWidget::sinkChanged);
    connect(m_applet, &SoundApplet::volumeChanged, this, &SoundTrayWidget::refreshTips, Qt::QueuedConnection);
    connect(static_cast<DApplication*>(qApp), &DApplication::iconThemeChanged, this, &SoundTrayWidget::updateIcon);
}

void SoundTrayWidget::setActive(const bool active)
{

}

void SoundTrayWidget::updateIcon()
{
    if (!m_sinkInter)
        return;

    const double volmue = m_applet->volumeValue();
    const bool mute = m_sinkInter->mute();

    QString iconString;
    QString volumeString;

    if (mute)
        volumeString = "muted";
    else if (volmue / 1000.0f >= double(2)/3)
        volumeString = "high";
    else if (volmue / 1000.0f >= double(1)/3)
        volumeString = "medium";
    else
        volumeString = "low";

    iconString = QString("audio-volume-%1-symbolic").arg(volumeString);

    const auto ratio = qApp->devicePixelRatio();
    const int iconSize = 16;
    const QIcon icon = QIcon::fromTheme(iconString);
    m_iconPixmap = icon.pixmap(iconSize * ratio, iconSize * ratio);
    m_iconPixmap.setDevicePixelRatio(ratio);

    update();
}

const QImage SoundTrayWidget::trayImage()
{
    return m_iconPixmap.toImage();
}

QWidget *SoundTrayWidget::trayTipsWidget()
{
    refreshTips(true);

    m_tipsLabel->resize(m_tipsLabel->sizeHint().width() + 10,
                        m_tipsLabel->sizeHint().height());

    return m_tipsLabel;
}

QWidget *SoundTrayWidget::trayPopupApplet()
{
    return m_applet;
}

const QString SoundTrayWidget::contextMenu() const
{
    QList<QVariant> items;
    items.reserve(2);

    QMap<QString, QVariant> open;
    open["itemId"] = MUTE;
    if (m_sinkInter->mute())
        open["itemText"] = tr("Unmute");
    else
        open["itemText"] = tr("Mute");
    open["isActive"] = true;
    items.push_back(open);

    QMap<QString, QVariant> settings;
    settings["itemId"] = SETTINS;
    settings["itemText"] = tr("Audio Settings");
    settings["isActive"] = true;
    items.push_back(settings);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void SoundTrayWidget::invokeMenuItem(const QString menuId, const bool checked)
{
    Q_UNUSED(checked);

    if (menuId == MUTE)
        m_sinkInter->SetMuteQueued(!m_sinkInter->mute());
    else if (menuId == SETTINS)
        DDBusSender()
            .service("com.deepin.dde.ControlCenter")
            .interface("com.deepin.dde.ControlCenter")
            .path("/com/deepin/dde/ControlCenter")
            .method(QString("ShowModule"))
            .arg(QString("sound"))
            .call();
}

QSize SoundTrayWidget::sizeHint() const
{
    return QSize(26, 26);
}

void SoundTrayWidget::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    updateIcon();
}

void SoundTrayWidget::wheelEvent(QWheelEvent *e)
{
    QWheelEvent *event = new QWheelEvent(e->pos(), e->delta(), e->buttons(), e->modifiers());
    qApp->postEvent(m_applet->mainSlider(), event);

    e->accept();
}

void SoundTrayWidget::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_iconPixmap.rect().center() / m_iconPixmap.devicePixelRatioF(), m_iconPixmap);
}

void SoundTrayWidget::refreshTips(const bool force)
{
    if (!force && !m_tipsLabel->isVisible())
        return;

    if(!m_sinkInter)
        return;

    QString value;
    if (m_sinkInter->mute()) {
        value = QString("0") + '%';
    } else {
        if (m_sinkInter->volume() * 1000 < m_applet->volumeValue())
            value = QString::number(m_applet->volumeValue() / 10) + '%';
        else
            value = QString::number(int(m_sinkInter->volume() * 100)) + '%';
    }
    m_tipsLabel->setText(QString(tr("Current Volume %1").arg(value)));
}

void SoundTrayWidget::sinkChanged(DBusSink *sink)
{
    m_sinkInter = sink;

    connect(m_sinkInter, &DBusSink::MuteChanged, this, &SoundTrayWidget::updateIcon);
    connect(m_sinkInter, &DBusSink::VolumeChanged, this, &SoundTrayWidget::updateIcon);
    updateIcon();
}
