// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "pluginchildpage.h"

#include <DStyle>
#include <DGuiApplicationHelper>

#include <QLabel>
#include <QVBoxLayout>
#include <QEvent>
#include <QPainterPath>
#include <QPushButton>

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

PluginChildPage::PluginChildPage(QWidget *parent)
    : QWidget(parent)
    , m_headerWidget(new QWidget(this))
    , m_back(new QPushButton(m_headerWidget))
    , m_title(new QLabel(m_headerWidget))
    , m_container(new QWidget(this))
    , m_topWidget(nullptr)
    , m_containerLayout(new QVBoxLayout(m_container))
{
    initUi();
    initConnection();
}

PluginChildPage::~PluginChildPage()
{
}

void PluginChildPage::pushWidget(QWidget *widget)
{
    // 首先将界面其他的窗体移除
    for (int i = m_containerLayout->count() - 1; i >= 0; i--) {
        QLayoutItem *item = m_containerLayout->itemAt(i);
        item->widget()->removeEventFilter(this);
        item->widget()->hide();
        m_containerLayout->removeItem(item);
    }
    m_topWidget = widget;
    if (widget) {
        widget->installEventFilter(this);
        m_containerLayout->addWidget(widget);
        widget->show();
    }
    QMetaObject::invokeMethod(this, &PluginChildPage::resetHeight, Qt::QueuedConnection);
}

void PluginChildPage::setTitle(const QString &text)
{
    m_title->setText(text);
}

void PluginChildPage::initUi()
{
    m_back->setIcon(backPixmap());
    m_back->setFixedWidth(24);
    m_back->setFlat(true);
    m_title->setAlignment(Qt::AlignCenter);
    QHBoxLayout *headerLayout = new QHBoxLayout(m_headerWidget);
    headerLayout->setContentsMargins(11, 12, 24 + 11, 12);
    headerLayout->setSpacing(0);
    headerLayout->addWidget(m_back);
    headerLayout->addWidget(m_title);
    m_headerWidget->setFixedHeight(48);

    QVBoxLayout *mainLayout = new QVBoxLayout(this);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->setSpacing(0);

    mainLayout->addWidget(m_headerWidget);
    mainLayout->addWidget(m_container);
    m_containerLayout->setContentsMargins(11, 0, 11, 0);
    m_containerLayout->setSpacing(0);
}

void PluginChildPage::initConnection()
{
    connect(m_back, &QPushButton::clicked, this, &PluginChildPage::back);
}

bool PluginChildPage::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_topWidget && event->type() == QEvent::Resize) {
        resetHeight();
    }
    return QWidget::eventFilter(watched, event);
}

void PluginChildPage::resetHeight()
{
    QMargins m = m_containerLayout->contentsMargins();
    m_container->setFixedHeight(m.top() + m.bottom() + (m_topWidget ? m_topWidget->height() : 0));
    setFixedHeight(m_headerWidget->height() + m_container->height());
}

QPixmap PluginChildPage::backPixmap() const
{
    QPixmap pixmap(16, 16);
    pixmap.fill(Qt::transparent);

    // 设置背景色
    QColor backColor;
    if (DGuiApplicationHelper::ColorType::DarkType == DGuiApplicationHelper::instance()->themeType()) {
        backColor = Qt::black;
        backColor.setAlphaF(0.05);
    } else {
        backColor = Qt::white;
        backColor.setAlphaF(0.2);
    }
    QPainter painter(&pixmap);
    painter.setRenderHints(QPainter::SmoothPixmapTransform | QPainter::Antialiasing);
    painter.setPen(Qt::NoPen);
    QPainterPath path;
    path.addEllipse(pixmap.rect());
    painter.fillPath(path, backColor);
    QSize arrowSize(10, 10);
    QPixmap arrowPixmap = DStyle::standardIcon(style(), DStyle::SP_ArrowLeave).pixmap(arrowSize);
    // 让其居中显示
    QSize backSize = pixmap.size();
    painter.drawPixmap((backSize.width() - arrowSize.width()) / 2, (backSize.height() - arrowSize.height()) / 2,
                       arrowSize.width(), arrowSize.height(), arrowPixmap);
    painter.end();
    return pixmap;
}
