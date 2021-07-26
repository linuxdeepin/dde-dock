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

#include "indicatortraywidget.h"
#include "util/utils.h"

#include <QLabel>
#include <QBoxLayout>

#include <QDBusConnection>
#include <QDBusInterface>

IndicatorTrayWidget::IndicatorTrayWidget(const QString &indicatorName, QWidget *parent, Qt::WindowFlags f)
    : AbstractTrayWidget(parent, f)
    , m_indicatorName(indicatorName)
    , m_gsettings(Utils::ModuleSettingsPtr("keyboard", QByteArray(), this))
    , m_enableClick(true)
{
    setAttribute(Qt::WA_TranslucentBackground);

    auto layout = new QVBoxLayout(this);
    layout->setContentsMargins(0, 0, 0, 0);
    m_label = new QLabel(this);

    QPalette p = m_label->palette();
    p.setColor(QPalette::Foreground, Qt::white);
    p.setColor(QPalette::Background, Qt::transparent);
    m_label->setPalette(p);

    m_label->setAttribute(Qt::WA_TranslucentBackground);

    layout->addWidget(m_label, 0, Qt::AlignCenter);
    setLayout(layout);

    // register dbus
    auto path = QString("/com/deepin/dde/Dock/Indicator/") + m_indicatorName;
    auto interface =  QString("com.deepin.dde.Dock.Indicator.") + m_indicatorName;
    auto sessionBus = QDBusConnection::sessionBus();
    sessionBus.registerObject(path,
                              interface,
                              this,
                              QDBusConnection::ExportScriptableSlots);

    if (m_gsettings) {
        // 显示键盘布局时更新label的状态
        if (m_gsettings->keys().contains("itemEnable"))
            enableLabel(m_gsettings->get("itemEnable").toBool());

        connect(m_gsettings, &QGSettings::changed, this, &IndicatorTrayWidget::onGSettingsChanged);
    }
}

IndicatorTrayWidget::~IndicatorTrayWidget()
{
}

QString IndicatorTrayWidget::itemKeyForConfig()
{
    return toIndicatorKey(m_indicatorName);
}

void IndicatorTrayWidget::updateIcon()
{

}

void IndicatorTrayWidget::sendClick(uint8_t buttonIndex, int x, int y)
{
    if (m_enableClick)
        Q_EMIT clicked(buttonIndex, x, y);
}

void IndicatorTrayWidget::enableLabel(bool enable)
{
    QPalette p = m_label->palette();
    if (!enable) {
        m_enableClick = false;
        p.setColor(QPalette::Disabled, QPalette::Foreground, Qt::lightGray);
        p.setColor(QPalette::Disabled, QPalette::Background, Qt::transparent);
        m_label->setEnabled(enable);
    } else {
        m_enableClick = true;
        p.setColor(QPalette::Normal, QPalette::BrightText, Qt::white);
        p.setColor(QPalette::Normal, QPalette::Background, Qt::transparent);
        m_label->setEnabled(enable);
    }

    m_label->setPalette(p);
    m_label->update();
}

void IndicatorTrayWidget::resizeEvent(QResizeEvent *event)
{
    return QWidget::resizeEvent(event);
}

void IndicatorTrayWidget::setPixmapData(const QByteArray &data)
{
    auto rawPixmap = QPixmap::fromImage(QImage::fromData(data));
    rawPixmap.setDevicePixelRatio(devicePixelRatioF());
    m_label->setPixmap(rawPixmap);
}

void IndicatorTrayWidget::setText(const QString &text)
{
    m_label->setText(text);
}

void IndicatorTrayWidget::onGSettingsChanged(const QString &key)
{
    Q_UNUSED(key);

    if (m_gsettings && m_gsettings->keys().contains("itemEnable")) {
        const bool itemEnable = m_gsettings->get("itemEnable").toBool();
        enableLabel(itemEnable);
    }
}

