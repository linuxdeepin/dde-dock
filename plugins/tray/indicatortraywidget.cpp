// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "indicatortraywidget.h"
#include "constants.h"
#include "util/utils.h"

#include <QLabel>
#include <QBoxLayout>
#include <QResizeEvent>

#include <QDBusConnection>
#include <QDBusInterface>
#include <DGuiApplicationHelper>
DGUI_USE_NAMESPACE

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

    updateLabelColor();
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &IndicatorTrayWidget::updateLabelColor);
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
    m_enableClick = enable;
    m_label->setEnabled(enable);
    updateLabelColor();
}

void IndicatorTrayWidget::resizeEvent(QResizeEvent *event)
{
    QWidget::resizeEvent(event);

    if (event->size().height() <= PLUGIN_BACKGROUND_MIN_SIZE || event->oldSize().height() <= PLUGIN_BACKGROUND_MIN_SIZE)
        updateLabelColor();
}

void IndicatorTrayWidget::updateLabelColor()
{
    QPalette p = m_label->palette();
    p.setColor(QPalette::Foreground, m_label->isEnabled() ? Qt::white : Qt::lightGray);

    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        p.setColor(QPalette::Foreground, m_label->isEnabled() ? Qt::black : Qt::darkGray);
    }

    m_label->setPalette(p);
    m_label->update();
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

