/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
 *
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

#include "switchitem.h"
#include "bluetoothconstants.h"

#include <DHiDPIHelper>
#include <DApplicationHelper>

#include <QHBoxLayout>
#include <QFontMetrics>
#include <QLabel>
#include <QEvent>

extern void initFontColor(QWidget *widget);

SwitchItem::SwitchItem(QWidget *parent)
    : QWidget(parent)
    , m_title(new QLabel(this))
    , m_switchBtn(new DSwitchButton(this))
    , m_spinner(new DSpinner (this))
    , m_default(false)
    , m_timer (new QTimer (this))
{
    initFontColor(m_title);
    m_timer->setSingleShot(true);
    m_timer->setInterval(1000);

    m_switchBtn->setFixedWidth(SWITCHBUTTONWIDTH);
    m_spinner->setFixedSize(24, 24);
    m_spinner->start();
    m_spinner->setVisible(false);

    const QPixmap pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh_dark.svg");

    m_loadingIndicator = new DLoadingIndicator;
    m_loadingIndicator->setSmooth(true);
    m_loadingIndicator->setAniDuration(500);
    m_loadingIndicator->setAniEasingCurve(QEasingCurve::InOutCirc);
    m_loadingIndicator->installEventFilter(this);
    m_loadingIndicator->setFixedSize(pixmap.size() / devicePixelRatioF());
    m_loadingIndicator->viewport()->setAutoFillBackground(false);
    m_loadingIndicator->setFrameShape(QFrame::NoFrame);
    m_loadingIndicator->installEventFilter(this);

    auto themeChanged = [&](DApplicationHelper::ColorType themeType){
        Q_UNUSED(themeType)
        setLoadIndicatorIcon();
    };
    themeChanged(DApplicationHelper::instance()->themeType());

    setFixedHeight(CONTROLHEIGHT);
    auto switchLayout = new QHBoxLayout;
    switchLayout->setSpacing(0);
    switchLayout->setMargin(0);
    switchLayout->addSpacing(MARGIN);
    switchLayout->addWidget(m_title);
    switchLayout->addStretch();
    switchLayout->addWidget(m_loadingIndicator);
    switchLayout->addSpacing(MARGIN);
    switchLayout->addWidget(m_switchBtn);
    switchLayout->addWidget(m_spinner);
    switchLayout->addSpacing(MARGIN);
    setLayout(switchLayout);

    connect(m_switchBtn, &DSwitchButton::toggled, [&](bool change) {
        m_checkState = change;
        emit checkedChanged(change);
        m_switchBtn->setEnabled(false);
        m_timer->start();
    });
    connect(m_timer, &QTimer::timeout, [&] (){
       m_switchBtn->setEnabled(true);
    });
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, themeChanged);
}

void SwitchItem::setChecked(const bool checked,bool notify)
{
    m_checkState = checked;
    if(!notify) {   // 防止收到蓝牙开启或关闭信号后再触发一次打开或关闭
        m_switchBtn->blockSignals(true);
        m_switchBtn->setChecked(checked);
        m_switchBtn->blockSignals(false);
        emit justUpdateView(checked);
    }
    else {
        m_switchBtn->setChecked(checked);
    }
}

void SwitchItem::setTitle(const QString &title)
{
    int width = POPUPWIDTH - MARGIN * 2 - m_switchBtn->width() - 3;
    QString strTitle = QFontMetrics(m_title->font()).elidedText(title, Qt::ElideRight, width);
    m_title->setText(strTitle);
}

bool SwitchItem::eventFilter(QObject *obj, QEvent *event)
{
    if (obj == m_loadingIndicator) {
        if (event->type() == QEvent::MouseButtonPress) {
            if(!m_loadingIndicator->loading())
                Q_EMIT refresh();
        }
    }
    return false;
}

void SwitchItem::setLoading(const bool bloading)
{
    m_loadingIndicator->setLoading(bloading);
}

void SwitchItem::setLoadIndicatorIcon()
{
    QString filePath =  ":/wireless/resources/wireless/refresh.svg";
    if(DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        filePath = ":/wireless/resources/wireless/refresh_dark.svg";
    const QPixmap pixmap = DHiDPIHelper::loadNxPixmap(filePath);
    m_loadingIndicator->setImageSource(pixmap);
}

void SwitchItem::loadStatusChange(bool isLoad)
{
    if (isLoad) {
        m_switchBtn->hide();
        m_spinner->show();
    } else {
        m_spinner->hide();
        m_switchBtn->show();
    }
}

//void SwitchItem::mousePressEvent(QMouseEvent *event)
//{
//    emit clicked(m_adapterId);
//    QWidget::mousePressEvent(event);
//}
