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
#include "fashiontray/fashiontrayconstants.h"
#include "system-trays/systemtrayitem.h"

#include <QDebug>
#include <QResizeEvent>

#define ExpandedKey "fashion-tray-expanded"

int FashionTrayItem::TrayWidgetWidth = TrayWidgetWidthMin;
int FashionTrayItem::TrayWidgetHeight = TrayWidgetHeightMin;

FashionTrayItem::FashionTrayItem(TrayPlugin *trayPlugin, QWidget *parent)
    : QWidget(parent),
      m_mainBoxLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight)),
      m_leftSpliter(new QLabel),
      m_rightSpliter(new QLabel),
      m_attentionDelayTimer(new QTimer(this)),
      m_trayPlugin(trayPlugin),
      m_controlWidget(new FashionTrayControlWidget(trayPlugin->dockPosition())),
      m_currentDraggingTray(nullptr),
      m_normalContainer(new NormalContainer(m_trayPlugin)),
      m_attentionContainer(new AttentionContainer(m_trayPlugin))
{
    setAcceptDrops(true);

    m_leftSpliter->setStyleSheet("background-color: rgba(255, 255, 255, 0.1);");
    m_rightSpliter->setStyleSheet("background-color: rgba(255, 255, 255, 0.1);");

    m_controlWidget->setFixedSize(QSize(TrayWidgetWidth, TrayWidgetHeight));
    m_normalContainer->setVisible(false);
    m_attentionContainer->setVisible(false);

    m_mainBoxLayout->setMargin(0);
    m_mainBoxLayout->setContentsMargins(0, 0, 0, 0);
    m_mainBoxLayout->setSpacing(TraySpace);

    m_mainBoxLayout->addWidget(m_leftSpliter);
    m_mainBoxLayout->addWidget(m_normalContainer);
    m_mainBoxLayout->addWidget(m_controlWidget);
    m_mainBoxLayout->addWidget(m_attentionContainer);
    m_mainBoxLayout->addWidget(m_rightSpliter);

    m_mainBoxLayout->setAlignment(Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_leftSpliter, Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_normalContainer, Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_controlWidget, Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_attentionContainer, Qt::AlignCenter);
    m_mainBoxLayout->setAlignment(m_rightSpliter, Qt::AlignCenter);

    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
    setLayout(m_mainBoxLayout);

    m_attentionDelayTimer->setInterval(3000);
    m_attentionDelayTimer->setSingleShot(true);

    connect(m_controlWidget, &FashionTrayControlWidget::expandChanged, this, &FashionTrayItem::onExpandChanged);
    connect(m_normalContainer, &NormalContainer::attentionChanged, this, &FashionTrayItem::onWrapperAttentionChanged);
    connect(m_attentionContainer, &NormalContainer::attentionChanged, this, &FashionTrayItem::onWrapperAttentionChanged);

    // do not call init immediately the TrayPlugin has not be constructed for now
    QTimer::singleShot(0, this, &FashionTrayItem::init);
}

void FashionTrayItem::setTrayWidgets(const QMap<QString, AbstractTrayWidget *> &itemTrayMap)
{
    clearTrayWidgets();

    for (auto it = itemTrayMap.constBegin(); it != itemTrayMap.constEnd(); ++it) {
        trayWidgetAdded(it.key(), it.value());
    }
}

void FashionTrayItem::trayWidgetAdded(const QString &itemKey, AbstractTrayWidget *trayWidget)
{
    if (m_normalContainer->containsWrapperByTrayWidget(trayWidget)) {
        qDebug() << "Reject! want to insert duplicate trayWidget:" << itemKey << trayWidget;
        return;
    }

    FashionTrayWidgetWrapper *wrapper = new FashionTrayWidgetWrapper(itemKey, trayWidget);

    do {
        if (m_normalContainer->acceptWrapper(wrapper)) {
            m_normalContainer->addWrapper(wrapper);
            break;
        }
    } while (false);

    requestResize();
}

void FashionTrayItem::trayWidgetRemoved(AbstractTrayWidget *trayWidget)
{
    bool deleted = false;

    do {
        if (m_normalContainer->removeWrapperByTrayWidget(trayWidget)) {
            deleted = true;
            break;
        }
        if (m_attentionContainer->removeWrapperByTrayWidget(trayWidget)) {
            deleted = true;
            break;
        }
    } while (false);

    if (!deleted) {
        qDebug() << "Error! can not find the tray widget in fashion tray list" << trayWidget;
    }

    requestResize();
}

void FashionTrayItem::clearTrayWidgets()
{
    m_normalContainer->clearWrapper();
    m_attentionContainer->clearWrapper();

    requestResize();
}

void FashionTrayItem::setDockPosition(Dock::Position pos)
{
    m_controlWidget->setDockPostion(pos);
    SystemTrayItem::setDockPostion(pos);

    m_normalContainer->setDockPosition(pos);
    m_attentionContainer->setDockPosition(pos);

    if (pos == Dock::Position::Top || pos == Dock::Position::Bottom) {
        m_mainBoxLayout->setDirection(QBoxLayout::Direction::LeftToRight);
    } else{
        m_mainBoxLayout->setDirection(QBoxLayout::Direction::TopToBottom);
    }

    requestResize();
}

void FashionTrayItem::onExpandChanged(const bool expand)
{
    m_trayPlugin->saveValue(ExpandedKey, expand);

    refreshHoldContainerPosition();

    if (expand) {
        m_normalContainer->setExpand(expand);
    } else {
        // hide all tray immediately if Dock is in maxed size
        // the property "DockIsMaxiedSize" of qApp is set by DockSettings class
        if (qApp->property("DockIsMaxiedSize").toBool()) {
            m_normalContainer->setExpand(expand);
        } else {
            // hide all tray widget delay for fold animation
            QTimer::singleShot(350, this, [=] {
                m_normalContainer->setExpand(expand);
            });
        }
    }

    m_attentionContainer->setExpand(expand);

    m_attentionDelayTimer->start();

    attentionWrapperToNormalWrapper();

    requestResize();
}

// used by QMetaObject::invokeMethod in TrayPluginItem / MainPanel class
void FashionTrayItem::setSuggestIconSize(QSize size)
{
    size = size * 0.6;

    int length = qMin(size.width(), size.height());
    // 设置最小值
//    length = qMax(length, TrayWidgetWidthMin);

    if (length == TrayWidgetWidth || length == TrayWidgetHeight) {
        return;
    }

    TrayWidgetWidth = length;
    TrayWidgetHeight = length;

    QSize newSize(length, length);

    m_controlWidget->setFixedSize(newSize);

    m_normalContainer->setWrapperSize(newSize);
    m_attentionContainer->setWrapperSize(newSize);

    requestResize();
}

void FashionTrayItem::setRightSplitVisible(const bool visible)
{
    if (visible) {
        m_rightSpliter->setStyleSheet("background-color: rgba(255, 255, 255, 0.1);");
    } else {
        m_rightSpliter->setStyleSheet("background-color: transparent;");
    }
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
    const Dock::Position dockPosition = m_trayPlugin->dockPosition();

    if (dockPosition == Dock::Position::Top || dockPosition == Dock::Position::Bottom) {
        m_leftSpliter->setFixedSize(SpliterSize, mSize.height() * 0.8);
        m_rightSpliter->setFixedSize(SpliterSize, mSize.height() * 0.8);
    } else{
        m_leftSpliter->setFixedSize(mSize.width() * 0.8, SpliterSize);
        m_rightSpliter->setFixedSize(mSize.width() * 0.8, SpliterSize);
    }

    QWidget::resizeEvent(event);
}

void FashionTrayItem::dragEnterEvent(QDragEnterEvent *event)
{
    // accept but do not handle the trays drag event
    // in order to avoid the for forbidden label displayed on the mouse
    if (event->mimeData()->hasFormat(TRAY_ITEM_DRAG_MIMEDATA)) {
        event->accept();
        return;
    }

    QWidget::dragEnterEvent(event);
}

QSize FashionTrayItem::sizeHint() const
{
    return wantedTotalSize();
}

void FashionTrayItem::init()
{
    qDebug() << "init Fashion mode tray plugin item";
    m_controlWidget->setExpanded(m_trayPlugin->getValue(ExpandedKey, true).toBool());
    setDockPosition(m_trayPlugin->dockPosition());
    onExpandChanged(m_controlWidget->expanded());
}

QSize FashionTrayItem::wantedTotalSize() const
{
    QSize size;
    const Dock::Position dockPosition = m_trayPlugin->dockPosition();

    if (m_controlWidget->expanded()) {
        if (dockPosition == Dock::Position::Top || dockPosition == Dock::Position::Bottom) {
            size.setWidth(
                        SpliterSize * 2 // 两个分隔条
                        + TraySpace * 2 // 两个分隔条旁边的 space
                        + TrayWidgetWidth // 控制按钮
                        + m_normalContainer->sizeHint().width() // 普通区域
//                        + m_holdContainer->sizeHint().width() // 保留区域
                        + m_attentionContainer->sizeHint().width() // 活动区域
                        );
            size.setHeight(height());
        } else {
            size.setWidth(width());
            size.setHeight(
                        SpliterSize * 2 // 两个分隔条
                        + TraySpace * 2 // 两个分隔条旁边的 space
                        + TrayWidgetWidth // 控制按钮
                        + m_normalContainer->sizeHint().height() // 普通区域
//                        + m_holdContainer->sizeHint().height() // 保留区域
                        + m_attentionContainer->sizeHint().height() // 活动区域
                        );
        }
    } else {
        if (dockPosition == Dock::Position::Top || dockPosition == Dock::Position::Bottom) {
            size.setWidth(
                        SpliterSize * 2 // 两个分隔条
                        + TraySpace * 2 // 两个分隔条旁边的 space
                        + TrayWidgetWidth // 控制按钮
//                        + m_holdContainer->sizeHint().width() // 保留区域
                        + m_attentionContainer->sizeHint().width() // 活动区域
                        );
            size.setHeight(height());
        } else {
            size.setWidth(width());
            size.setHeight(
                        SpliterSize * 2 // 两个分隔条
                        + TraySpace * 2 // 两个分隔条旁边的 space
                        + TrayWidgetWidth // 控制按钮
//                        + m_holdContainer->sizeHint().height() // 保留区域
                        + m_attentionContainer->sizeHint().height() // 活动区域
                        );
        }
    }

    return size;
}

void FashionTrayItem::onWrapperAttentionChanged(FashionTrayWidgetWrapper *wrapper, const bool attention)
{
    if (m_controlWidget->expanded()) {
        return;
    }

    // 在timer处于Active状态期间不设置新的活动图标
    if (attention && m_attentionDelayTimer->isActive()) {
        return;
    }

    if (attention) {
        // ignore the attention which is come from AttentionContainer
        if (m_attentionContainer->containsWrapper(wrapper)) {
            return;
        }
        // move previous attention wrapper from AttentionContainer to NormalContainer
        attentionWrapperToNormalWrapper();
        // move current attention wrapper from NormalContainer to AttentionContainer
        normalWrapperToAttentionWrapper(wrapper);
    } else {
        // only focus the disattention from AttentionContainer
        if (m_attentionContainer->containsWrapper(wrapper)) {
            attentionWrapperToNormalWrapper();
        }
    }

    m_attentionDelayTimer->start();

    requestResize();
}

void FashionTrayItem::attentionWrapperToNormalWrapper()
{
    FashionTrayWidgetWrapper *preAttentionWrapper = m_attentionContainer->takeAttentionWrapper();
    if (preAttentionWrapper) {
        m_normalContainer->addWrapper(preAttentionWrapper);
    }
}

void FashionTrayItem::normalWrapperToAttentionWrapper(FashionTrayWidgetWrapper *wrapper)
{
    FashionTrayWidgetWrapper *attentionWrapper = m_normalContainer->takeWrapper(wrapper);
    if (attentionWrapper) {
        m_attentionContainer->addWrapper(attentionWrapper);
    } else {
        qDebug() << "Warnning: not find the attention wrapper in NormalContainer";
    }
}

void FashionTrayItem::requestResize()
{
    // reset property "FashionTraySize" to notify dock resize
    // DockPluginsController will watch this property
    setProperty("FashionTraySize", sizeHint());
}

void FashionTrayItem::refreshHoldContainerPosition()
{
//    const int destIndex = m_mainBoxLayout->indexOf(m_controlWidget)
//            + (m_controlWidget->expanded() ? 0 : 1);

//    m_mainBoxLayout->insertWidget(destIndex, m_holdContainer);
}
