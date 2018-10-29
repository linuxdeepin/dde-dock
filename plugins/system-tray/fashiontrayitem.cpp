/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#include "fashiontrayitem.h"
#include "system-trays/abstractsystemtraywidget.h"

#include <QDebug>
#include <QResizeEvent>

#define SpliterSize 2
#define TraySpace 10
#define TrayWidgetWidth 28
#define TrayWidgetHeight 28

FashionTrayItem::FashionTrayItem(Dock::Position pos, QWidget *parent)
    : QWidget(parent),
      m_mainBoxLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this)),
      m_trayBoxLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight, this)),
      m_leftSpliter(new QLabel(this)),
      m_rightSpliter(new QLabel(this)),
      m_controlWidget(new FashionTrayControlWidget(pos, this)),
      m_currentAttentionTray(nullptr),
      m_dockPosistion(pos)
{
    m_leftSpliter->setStyleSheet("background-color: rgba(255, 255, 255, 0.1);");
    m_rightSpliter->setStyleSheet("background-color: rgba(255, 255, 255, 0.1);");

    m_controlWidget->setFixedSize(QSize(TrayWidgetWidth, TrayWidgetHeight));

    m_mainBoxLayout->setMargin(0);
    m_mainBoxLayout->setContentsMargins(0, 0, 0, 0);
    m_mainBoxLayout->setSpacing(TraySpace);

    m_trayBoxLayout->setMargin(0);
    m_trayBoxLayout->setContentsMargins(0, 0, 0, 0);
    m_trayBoxLayout->setSpacing(TraySpace);

    m_mainBoxLayout->addWidget(m_leftSpliter);
    m_mainBoxLayout->addLayout(m_trayBoxLayout);
    m_mainBoxLayout->addWidget(m_controlWidget);
    m_mainBoxLayout->addWidget(m_rightSpliter);

    m_mainBoxLayout->setAlignment(Qt::AlignCenter);
    m_trayBoxLayout->setAlignment(Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_leftSpliter, Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_controlWidget, Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_rightSpliter, Qt::AlignCenter);

    setLayout(m_mainBoxLayout);

    setDockPostion(pos);
    onTrayListExpandChanged(m_controlWidget->expanded());

    connect(m_controlWidget, &FashionTrayControlWidget::expandChanged, this, &FashionTrayItem::onTrayListExpandChanged);
}

void FashionTrayItem::setTrayWidgets(const QList<AbstractTrayWidget *> &trayWidgetList)
{
    for (auto widget : trayWidgetList) {
        trayWidgetAdded(widget);
    }
}

void FashionTrayItem::trayWidgetAdded(AbstractTrayWidget *trayWidget)
{
    if (m_trayWidgetWrapperMap.keys().contains(trayWidget)) {
        return;
    }

    FashionTrayWidgetWrapper *wrapper = new FashionTrayWidgetWrapper(trayWidget);
    wrapper->setFixedSize(QSize(TrayWidgetWidth, TrayWidgetHeight));

    m_trayWidgetWrapperMap.insert(trayWidget, wrapper);
    m_trayBoxLayout->addWidget(wrapper);
    wrapper->setVisible(m_controlWidget->expanded());

    if (wrapper->attention()) {
        setCurrentAttentionTray(wrapper);
    }

    connect(wrapper, &FashionTrayWidgetWrapper::attentionChanged, this, &FashionTrayItem::onTrayAttentionChanged, Qt::UniqueConnection);

    if (trayWidget->trayTyep() == AbstractTrayWidget::TrayType::SystemTray) {
        AbstractSystemTrayWidget * sysTrayWidget = static_cast<AbstractSystemTrayWidget *>(trayWidget);
        connect(sysTrayWidget, &AbstractSystemTrayWidget::requestWindowAutoHide, this, &FashionTrayItem::requestWindowAutoHide, Qt::UniqueConnection);
        connect(sysTrayWidget, &AbstractSystemTrayWidget::requestRefershWindowVisible, this, &FashionTrayItem::requestRefershWindowVisible, Qt::UniqueConnection);
    }

    requestResize();
}

void FashionTrayItem::trayWidgetRemoved(AbstractTrayWidget *trayWidget)
{
    auto it = m_trayWidgetWrapperMap.constBegin();

    for (; it != m_trayWidgetWrapperMap.constEnd(); ++it) {
        if (it.key() == trayWidget) {
            // removing the attention tray
            if (m_currentAttentionTray == it.value()) {
                if (m_controlWidget->expanded()) {
                    m_trayBoxLayout->removeWidget(m_currentAttentionTray);
                } else {
                    m_mainBoxLayout->removeWidget(m_currentAttentionTray);
                }
                m_currentAttentionTray = nullptr;
            } else {
                m_trayBoxLayout->removeWidget(it.value());
            }
            it.value()->deleteLater();
            m_trayWidgetWrapperMap.remove(it.key());
            break;
        }
    }

    if (it == m_trayWidgetWrapperMap.constEnd()) {
        qDebug() << "can not find the tray widget in fashion tray list:" << trayWidget;
        return;
    }

    requestResize();
}

void FashionTrayItem::clearTrayWidgets()
{
    if (m_currentAttentionTray) {
        m_mainBoxLayout->removeWidget(m_currentAttentionTray);
        m_currentAttentionTray = nullptr;
    }

    for (auto wrapper : m_trayWidgetWrapperMap.values()) {
        m_trayBoxLayout->removeWidget(wrapper);
        wrapper->deleteLater();
    }

    m_trayWidgetWrapperMap.clear();

    requestResize();
}

void FashionTrayItem::setDockPostion(Dock::Position pos)
{
    m_dockPosistion = pos;

    m_controlWidget->setDockPostion(m_dockPosistion);
    AbstractSystemTrayWidget::setDockPostion(m_dockPosistion);

    if (pos == Dock::Position::Top || pos == Dock::Position::Bottom) {
        m_mainBoxLayout->setDirection(QBoxLayout::Direction::LeftToRight);
        m_trayBoxLayout->setDirection(QBoxLayout::Direction::LeftToRight);
    } else{
        m_mainBoxLayout->setDirection(QBoxLayout::Direction::TopToBottom);
        m_trayBoxLayout->setDirection(QBoxLayout::Direction::TopToBottom);
    }

    requestResize();
}

void FashionTrayItem::onTrayListExpandChanged(const bool expand)
{
    if (m_currentAttentionTray) {
        if (expand) {
            m_mainBoxLayout->removeWidget(m_currentAttentionTray);
            m_trayBoxLayout->addWidget(m_currentAttentionTray);
        } else {
            m_trayBoxLayout->removeWidget(m_currentAttentionTray);
            m_mainBoxLayout->insertWidget(m_mainBoxLayout->indexOf(m_controlWidget) + 1, m_currentAttentionTray);
        }
    }

    m_mainBoxLayout->setAlignment(m_currentAttentionTray, Qt::AlignCenter);

    for (auto i = m_trayWidgetWrapperMap.begin(); i != m_trayWidgetWrapperMap.end(); ++i) {
        if (i.value() == m_currentAttentionTray) {
            continue;
        }
        i.value()->setVisible(expand);
    }

    requestResize();
}

void FashionTrayItem::showEvent(QShowEvent *event)
{
    requestResize();

    QWidget::showEvent(event);
}

void FashionTrayItem::hideEvent(QHideEvent *event)
{
    requestResize();

    QWidget::hideEvent(event);
}

void FashionTrayItem::resizeEvent(QResizeEvent *event)
{
    const QSize &mSize = event->size();

    if (m_dockPosistion == Dock::Position::Top || m_dockPosistion == Dock::Position::Bottom) {
        m_leftSpliter->setFixedSize(SpliterSize, mSize.height() * 0.8);
        m_rightSpliter->setFixedSize(SpliterSize, mSize.height() * 0.8);
    } else{
        m_leftSpliter->setFixedSize(mSize.width() * 0.8, SpliterSize);
        m_rightSpliter->setFixedSize(mSize.width() * 0.8, SpliterSize);
    }

    QWidget::resizeEvent(event);
}

QSize FashionTrayItem::sizeHint() const
{
    return wantedTotalSize();
}

QSize FashionTrayItem::wantedTotalSize() const
{
    QSize size;

    if (m_controlWidget->expanded()) {
        if (m_dockPosistion == Dock::Position::Top || m_dockPosistion == Dock::Position::Bottom) {
            size.setWidth(m_trayWidgetWrapperMap.size() * TrayWidgetWidth // 所有插件
                          + TrayWidgetWidth // 控制按钮
                          + SpliterSize * 2 // 两个分隔条
                          + 3 * TraySpace // MainBoxLayout所有space
                          + (m_trayWidgetWrapperMap.size() - 1) * TraySpace); // TrayBoxLayout所有space
            size.setHeight(height());
        } else {
            size.setWidth(width());
            size.setHeight(m_trayWidgetWrapperMap.size() * TrayWidgetHeight // 所有插件
                          + TrayWidgetHeight // 控制按钮
                          + SpliterSize * 2 // 两个分隔条
                          + 3 * TraySpace // MainBoxLayout所有space
                          + (m_trayWidgetWrapperMap.size() - 1) * TraySpace); // TrayBoxLayout所有space
        }
    } else {
        if (m_dockPosistion == Dock::Position::Top || m_dockPosistion == Dock::Position::Bottom) {
            size.setWidth(TrayWidgetWidth // 控制按钮
                          + (m_currentAttentionTray ? TrayWidgetWidth : 0) // 活动状态的tray
                          + SpliterSize * 2 // 两个分隔条
                          + 3 * TraySpace); // MainBoxLayout所有space
            size.setHeight(height());
        } else {
            size.setWidth(width());
            size.setHeight(TrayWidgetHeight // 控制按钮
                          + (m_currentAttentionTray ? TrayWidgetHeight : 0) // 活动状态的tray
                          + SpliterSize * 2 // 两个分隔条
                          + 3 * TraySpace); // MainBoxLayout所有space
        }
    }

    return size;
}

void FashionTrayItem::onTrayAttentionChanged(const bool attention)
{
    FashionTrayWidgetWrapper *wrapper = static_cast<FashionTrayWidgetWrapper *>(sender());

    Q_ASSERT(wrapper);

    if (attention) {
        setCurrentAttentionTray(wrapper);
    } else {
        if (m_currentAttentionTray != wrapper) {
            return;
        }

        if (m_controlWidget->expanded()) {
            m_currentAttentionTray = nullptr;
        } else {
            moveInAttionTray();
            m_currentAttentionTray = nullptr;
            requestResize();
        }
    }
}

void FashionTrayItem::setCurrentAttentionTray(FashionTrayWidgetWrapper *attentionWrapper)
{
    if (!attentionWrapper) {
        return;
    }

    if (m_controlWidget->expanded()) {
        m_currentAttentionTray = attentionWrapper;
    } else {
        if (m_currentAttentionTray == attentionWrapper) {
            return;
        }
        moveInAttionTray();
        bool sizeChanged = !m_currentAttentionTray;
        m_currentAttentionTray = attentionWrapper;
        moveOutAttionTray();
        if (sizeChanged) {
            requestResize();
        }
    }

    m_mainBoxLayout->setAlignment(m_currentAttentionTray, Qt::AlignCenter);
}

void FashionTrayItem::requestResize()
{
    // reset property "FashionSystemTraySize" to notify dock resize
    // DockPluginsController will watch this property

    setProperty("FashionSystemTraySize", isVisible() ? wantedTotalSize() : QSize(0, 0));
}

void FashionTrayItem::moveOutAttionTray()
{
    if (!m_currentAttentionTray) {
        return;
    }

    m_trayBoxLayout->removeWidget(m_currentAttentionTray);
    m_mainBoxLayout->insertWidget(m_mainBoxLayout->indexOf(m_rightSpliter), m_currentAttentionTray);
    m_currentAttentionTray->setVisible(true);
}

void FashionTrayItem::moveInAttionTray()
{
    if (!m_currentAttentionTray) {
        return;
    }

    m_mainBoxLayout->removeWidget(m_currentAttentionTray);
    m_trayBoxLayout->addWidget(m_currentAttentionTray);
    m_currentAttentionTray->setVisible(m_controlWidget->expanded());
}

void FashionTrayItem::switchAttionTray(FashionTrayWidgetWrapper *attentionWrapper)
{
    if (!m_currentAttentionTray || !attentionWrapper) {
        return;
    }

    m_mainBoxLayout->replaceWidget(m_currentAttentionTray, attentionWrapper);
    m_trayBoxLayout->removeWidget(attentionWrapper);
    m_trayBoxLayout->addWidget(m_currentAttentionTray);

    attentionWrapper->setVisible(true);
    m_currentAttentionTray->setVisible(m_controlWidget->expanded());

    m_currentAttentionTray = attentionWrapper;
}

void FashionTrayItem::requestWindowAutoHide(const bool autoHide)
{
    // reset property "RequestWindowAutoHide" to EMIT the signal of DockItem
    // TODO: 考虑新增插件接口

    setProperty("RequestWindowAutoHide", autoHide);
}

void FashionTrayItem::requestRefershWindowVisible()
{
    // reset property "RequestRefershWindowVisible" to EMIT the signal of DockItem
    // TODO: 考虑新增插件接口

    setProperty("RequestRefershWindowVisible", !property("RequestRefershWindowVisible").toBool());
}
