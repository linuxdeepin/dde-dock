/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
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
#include "../widgets/tipswidget.h"
#include "../frame/util/imageutil.h"
#include "../frame/util/utils.h"

#include <DApplication>
#include <DDBusSender>
#include <DGuiApplicationHelper>

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>
#include <QApplication>
#include <QDBusInterface>

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

// menu actions
#define MUTE     "mute"
#define SETTINGS "settings"

using namespace Dock;

SoundItem::SoundItem(QWidget *parent)
    : QWidget(parent)
    , m_tipsLabel(new TipsWidget(this))
    , m_applet(new SoundApplet(this))
    , m_sinkInter(nullptr)
{
    m_tipsLabel->setAccessibleName("soundtips");
    m_tipsLabel->setVisible(false);

    m_applet->setVisible(false);

    connect(m_applet, &SoundApplet::defaultSinkChanged, this, &SoundItem::sinkChanged);
    connect(m_applet, &SoundApplet::volumeChanged, this, &SoundItem::refresh, Qt::QueuedConnection);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        refreshIcon();
    });
}

QWidget *SoundItem::tipsWidget()
{
    if (m_sinkInter)
        refreshTips(std::min(150, qRound(m_sinkInter->volume() * 100.0)), true);
    else
        refreshTips(m_applet->volumeValue(), true);

    m_tipsLabel->resize(m_tipsLabel->sizeHint().width() + 10,
                        m_tipsLabel->sizeHint().height());

    return m_tipsLabel;
}

QWidget *SoundItem::popupApplet()
{
    return m_applet;
}

const QString SoundItem::contextMenu()
{
    QList<QVariant> items;
    items.reserve(2);

    QMap<QString, QVariant> open;
    open["itemId"] = MUTE;
    if (m_sinkInter->mute()) {
        open["itemText"] = tr("Unmute");
        if (!m_applet->existActiveOutputDevice())
            open["isActive"] = false;
        else
            open["isActive"] = true;
    } else {
        open["itemText"] = tr("Mute");
        open["isActive"] = true;
    }
    items.push_back(open);

    if (!QFile::exists(ICBC_CONF_FILE)) {
        QMap<QString, QVariant> settings;
        settings["itemId"] = SETTINGS;
        settings["itemText"] = tr("Sound settings");
        settings["isActive"] = true;
        items.push_back(settings);
#ifdef QT_DEBUG
        qInfo() << "----------icbc sound setting.";
#endif
    }

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
        m_sinkInter->SetMuteQueued(!m_sinkInter->mute());
    else if (menuId == SETTINGS)
        DDBusSender()
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method(QString("ShowModule"))
        .arg(QString("sound"))
        .call();
}

void SoundItem::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    // 保持横纵比
    if (position == Dock::Bottom || position == Dock::Top) {
        setMaximumWidth(height());
        setMaximumHeight(QWIDGETSIZE_MAX);
    } else {
        setMaximumHeight(width());
        setMaximumWidth(QWIDGETSIZE_MAX);
    }

    refreshIcon();
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
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_iconPixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / m_iconPixmap.devicePixelRatioF(), m_iconPixmap);
}

void SoundItem::refresh(const int volume)
{
    refreshIcon();
    refreshTips(volume, false);
}

void SoundItem::refreshIcon()
{
    if (!m_sinkInter)
        return;

    const double volmue = m_applet->volumeValue();
    const double maxVolmue = m_applet->maxVolumeValue();
    const bool mute = m_sinkInter->mute();
    const Dock::DisplayMode displayMode = Dock::DisplayMode::Efficient;

    QString iconString;
    if (displayMode == Dock::Fashion) {
        QString volumeString;
        if (volmue >= 1000)
            volumeString = "100";
        else
            volumeString = QString("0") + ('0' + int(volmue / 100)) + "0";

        iconString = "audio-volume-" + volumeString;

        if (mute)
            iconString += "-muted";
    } else {
        QString volumeString;
        if (mute)
            volumeString = "muted";
        else if (int(volmue) == 0)
            volumeString = "off";
        else if (volmue / maxVolmue > double(2) / 3)
            volumeString = "high";
        else if (volmue / maxVolmue > double(1) / 3)
            volumeString = "medium";
        else
            volumeString = "low";

        iconString = QString("audio-volume-%1-symbolic").arg(volumeString);
    }

    const auto ratio = devicePixelRatioF();
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconString.append(PLUGIN_MIN_ICON_NAME);

    m_iconPixmap = ImageUtil::loadSvg(iconString, ":/", iconSize, ratio);

    update();
}

void SoundItem::refreshTips(const int volume, const bool force)
{
    if (!force && !m_tipsLabel->isVisible())
        return;

    QString value;
    if (m_sinkInter->mute()) {
        m_tipsLabel->setText(QString(tr("Mute")));
    } else {
        m_tipsLabel->setText(QString(tr("Volume %1").arg(QString::number(volume) + '%')));
    }
}

void SoundItem::sinkChanged(DBusSink *sink)
{
    m_sinkInter = sink;

    if (m_sinkInter)
        refresh(std::min(150, qRound(m_sinkInter->volume() * 100.0)));
    else
        refresh(m_applet->volumeValue());
}
