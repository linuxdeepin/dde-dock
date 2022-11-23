/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "pluginchildpage.h"

#include <QLabel>
#include <QVBoxLayout>
#include <QEvent>

PluginChildPage::PluginChildPage(QWidget *parent)
    : QWidget(parent)
    , m_headerWidget(new QWidget(this))
    , m_back(new QLabel(m_headerWidget))
    , m_title(new QLabel(m_headerWidget))
    , m_container(new QWidget(this))
    , m_topWidget(nullptr)
    , m_containerLayout(new QVBoxLayout(m_container))
    , m_isBack(false)
{
    initUi();
    m_back->installEventFilter(this);
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
    widget->installEventFilter(this);
    m_containerLayout->addWidget(widget);
    widget->show();
    m_isBack = false;
    QMetaObject::invokeMethod(this, &PluginChildPage::resetHeight, Qt::QueuedConnection);
}

void PluginChildPage::setTitle(const QString &text)
{
    m_title->setText(text);
}

void PluginChildPage::setCanBack(bool canBack)
{
    m_back->setVisible(canBack);
}

bool PluginChildPage::isBack()
{
    return m_isBack;
}

void PluginChildPage::initUi()
{
    m_back->setFixedWidth(24);
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

bool PluginChildPage::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_back && event->type() == QEvent::MouseButtonRelease) {
        m_isBack = true;
        Q_EMIT back();
        return true;
    }
    if (watched == m_topWidget) {
        if (event->type() == QEvent::Hide) {
            Q_EMIT closeSelf();
            return true;
        }
        if (event->type() == QEvent::Resize)
            resetHeight();
    }
    return QWidget::eventFilter(watched, event);
}

void PluginChildPage::resetHeight()
{
    QMargins m = m_containerLayout->contentsMargins();
    m_container->setFixedHeight(m.top() + m.bottom() + m_topWidget->height());
    setFixedHeight(m_headerWidget->height() + m_container->height());
}
