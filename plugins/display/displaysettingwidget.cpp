// Copyright (C) 2021 ~ 2022 Uniontech Software Technology Co.,Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "displaysettingwidget.h"
#include "brightnessadjwidget.h"
#include "devcollaborationwidget.h"

#include <QPushButton>
#include <QVBoxLayout>

#include <DDBusSender>

const int ItemSpacing = 10;

DisplaySettingWidget::DisplaySettingWidget(QWidget *parent)
    : QWidget(parent)
    , m_brightnessAdjWidget(new BrightnessAdjWidget(this))
    , m_collaborationWidget(new DevCollaborationWidget(this))
    , m_settingBtn(new QPushButton(tr("Multi-Screen Collaboration"), this))
{
    initUI();

    connect(m_settingBtn, &QPushButton::clicked, this, [ this ](){
        DDBusSender().service("org.deepin.dde.ControlCenter1")
                .path("/org/deepin/dde/ControlCenter1")
                .interface("org.deepin.dde.ControlCenter1")
                .method("ShowPage").arg(QString("display")).call();
        Q_EMIT requestHide();
    });
}

void DisplaySettingWidget::initUI()
{
    setContentsMargins(0, 10, 0, 30);
    QVBoxLayout *mainLayout = new QVBoxLayout();
    mainLayout->setMargin(0);
    mainLayout->setSpacing(ItemSpacing);

    mainLayout->addWidget(m_brightnessAdjWidget);
    mainLayout->addWidget(m_collaborationWidget);
    mainLayout->addWidget(m_settingBtn);
    mainLayout->addStretch();

    setLayout(mainLayout);

    resizeWidgetHeight();
    connect(m_collaborationWidget, &DevCollaborationWidget::sizeChanged,
            this, &DisplaySettingWidget::resizeWidgetHeight);
}

void DisplaySettingWidget::resizeWidgetHeight()
{
    QMargins margins = this->contentsMargins();
    setFixedHeight(margins.top() + margins.bottom() + m_brightnessAdjWidget->height() +
                   m_collaborationWidget->height() + m_settingBtn->height() + ItemSpacing * 2);
}
