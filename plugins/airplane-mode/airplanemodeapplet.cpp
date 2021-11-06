/*
 * Copyright (C) 2020 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     weizhixiang <weizhixiang@uniontech.com>
 *
 * Maintainer: weizhixiang <weizhixiang@uniontech.com>
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

#include "airplanemodeapplet.h"
#include "constants.h"

#include <DSwitchButton>

#include <QLabel>
#include <QHBoxLayout>

DWIDGET_USE_NAMESPACE

AirplaneModeApplet::AirplaneModeApplet(QWidget *parent)
    : QWidget(parent)
    , m_switchBtn(new DSwitchButton(this))
{
    setMinimumWidth(PLUGIN_ITEM_WIDTH);
    setFixedHeight(30);
    QLabel *title = new QLabel(this);
    title->setText(tr("Airplane Mode"));
    QHBoxLayout *appletlayout = new QHBoxLayout;
    appletlayout->setMargin(0);
    appletlayout->setSpacing(0);
    appletlayout->addSpacing(0);
    appletlayout->addWidget(title);
    appletlayout->addStretch();
    appletlayout->addWidget(m_switchBtn);
    appletlayout->addSpacing(0);
    setLayout(appletlayout);

    connect(m_switchBtn, &DSwitchButton::checkedChanged, this, &AirplaneModeApplet::enableChanged);
}

void AirplaneModeApplet::setEnabled(bool enable)
{
    m_switchBtn->blockSignals(true);
    m_switchBtn->setChecked(enable);
    m_switchBtn->blockSignals(false);
}
