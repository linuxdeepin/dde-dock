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

#include "QHBoxLayout"

SwitchItem::SwitchItem(QWidget *parent)
    : QWidget(parent)
    , m_title(new QLabel(this))
    , m_switchBtn(new DSwitchButton(this))
    , m_default(false)
{
    setFixedHeight(35);
    auto switchLayout = new QHBoxLayout(this);
    switchLayout->setSpacing(0);
    switchLayout->setMargin(0);
    switchLayout->addSpacing(5);
    switchLayout->addWidget(m_title);
    switchLayout->addStretch();
    switchLayout->addWidget(m_switchBtn);
    switchLayout->addSpacing(5);
    setLayout(switchLayout);

    connect(m_switchBtn, &DSwitchButton::toggled, [&](bool change) {
        emit checkedChanged(change);
    });
}

void SwitchItem::setChecked(const bool checked)
{
    m_switchBtn->setChecked(checked);
}

void SwitchItem::setTitle(const QString &title)
{
    m_title->setText(title);
}

//void SwitchItem::mousePressEvent(QMouseEvent *event)
//{
//    emit clicked(m_adapterId);
//    QWidget::mousePressEvent(event);
//}
