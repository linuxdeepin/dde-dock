/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "sounditem.h"
#include "constants.h"

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>
#include <QApplication>

// menu actions
#define MUTE    "mute"
#define SETTINS "settings"

SoundItem::SoundItem(QWidget *parent)
    : QWidget(parent),

      m_tipsLabel(new QLabel(this)),
      m_applet(new SoundApplet(this)),
      m_sinkInter(nullptr)
{
//    QIcon::setThemeName("deepin");

    m_tipsLabel->setObjectName("sound");
    m_tipsLabel->setVisible(false);
    //    m_tipsLabel->setFixedWidth(145);
    m_tipsLabel->setAlignment(Qt::AlignCenter);
    m_tipsLabel->setStyleSheet("color:white;"
                               "padding: 0 3px;");

    m_applet->setVisible(false);

    connect(m_applet, static_cast<void (SoundApplet::*)(DBusSink*) const>(&SoundApplet::defaultSinkChanged), this, &SoundItem::sinkChanged);
    connect(m_applet, &SoundApplet::volumeChanged, this, &SoundItem::refershTips, Qt::QueuedConnection);
}

QWidget *SoundItem::tipsWidget()
{
    refershTips(true);

    m_tipsLabel->resize(m_tipsLabel->sizeHint().width() + 10,
                        m_tipsLabel->sizeHint().height());

    return m_tipsLabel;
}

QWidget *SoundItem::popupApplet()
{
    return m_applet;
}

const QString SoundItem::contextMenu() const
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

void SoundItem::invokeMenuItem(const QString menuId, const bool checked)
{
    Q_UNUSED(checked);

    if (menuId == MUTE)
        m_sinkInter->SetMute(!m_sinkInter->mute());
    else if (menuId == SETTINS)
        QProcess::startDetached("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:sound\"");
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
    {
        e->accept();
        emit requestContextMenu();
        return;
    }

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
    painter.drawPixmap(rect().center() - m_iconPixmap.rect().center() / qApp->devicePixelRatio(), m_iconPixmap);
}

void SoundItem::refershIcon()
{
    if (!m_sinkInter)
        return;

    const double volmue = m_applet->volumeValue();
    const bool mute = m_sinkInter->mute();
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    QString iconString;
    if (displayMode == Dock::Fashion)
    {
        QString volumeString;
        if (volmue >= 1000)
            volumeString = "100";
        else
            volumeString = QString("0") + ('0' + int(volmue / 100)) + "0";

        iconString = ":/icons/image/audio-volume-" + volumeString;

        if (mute)
            iconString += "-muted";
    } else {
        QString volumeString;
        if (mute)
            volumeString = "muted";
        else if (volmue / 1000.0f >= double(2)/3)
            volumeString = "high";
        else if (volmue / 1000.0f >= double(1)/3)
            volumeString = "medium";
        else
            volumeString = "low";

        iconString = QString(":/icons/image/audio-volume-%1-symbolic").arg(volumeString);
    }

    const auto ratio = qApp->devicePixelRatio();
    const int iconSize = displayMode == Dock::Fashion ? std::min(width(), height()) * 0.8 : 16;
    const QIcon icon = QIcon::fromTheme(iconString);
    m_iconPixmap = icon.pixmap(iconSize * ratio, iconSize * ratio);
    m_iconPixmap.setDevicePixelRatio(ratio);

    update();
}

void SoundItem::refershTips(const bool force)
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

void SoundItem::sinkChanged(DBusSink *sink)
{
    m_sinkInter = sink;

    connect(m_sinkInter, &DBusSink::MuteChanged, this, &SoundItem::refershIcon);
    connect(m_sinkInter, &DBusSink::VolumeChanged, this, &SoundItem::refershIcon);
    refershIcon();
}
