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
#include "system-trays/systemtrayitem.h"

#include <QDebug>
#include <QResizeEvent>

#define SpliterSize 2
#define TraySpace 10
#define TrayWidgetWidthMin 24
#define TrayWidgetHeightMin 24

int FashionTrayItem::TrayWidgetWidth = TrayWidgetWidthMin;
int FashionTrayItem::TrayWidgetHeight = TrayWidgetHeightMin;

FashionTrayItem::FashionTrayItem(TrayPlugin *trayPlugin, QWidget *parent)
    : QWidget(parent),
      m_mainBoxLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight)),
      m_trayBoxLayout(new QBoxLayout(QBoxLayout::Direction::LeftToRight)),
      m_leftSpliter(new QLabel),
      m_rightSpliter(new QLabel),
      m_attentionDelayTimer(new QTimer(this)),
      m_dockPosistion(trayPlugin->dockPosition()),
      m_trayPlugin(trayPlugin),
      m_controlWidget(new FashionTrayControlWidget(m_dockPosistion)),
      m_currentAttentionTray(nullptr)
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

    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
    setLayout(m_mainBoxLayout);

    m_attentionDelayTimer->setInterval(3000);
    m_attentionDelayTimer->setSingleShot(true);

    setDockPostion(m_dockPosistion);
    onTrayListExpandChanged(m_controlWidget->expanded());

    connect(m_controlWidget, &FashionTrayControlWidget::expandChanged, this, &FashionTrayItem::onTrayListExpandChanged);
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
    for (auto w : m_wrapperList) {
        if (w->absTrayWidget() == trayWidget) {
            qDebug() << "Reject! want to isert duplicate trayWidget:" << itemKey << trayWidget;
            return;
        }
    }

    FashionTrayWidgetWrapper *wrapper = new FashionTrayWidgetWrapper(itemKey, trayWidget);
    wrapper->setFixedSize(QSize(TrayWidgetWidth, TrayWidgetHeight));

    m_wrapperList.append(wrapper);
    m_trayBoxLayout->insertWidget(whereToInsert(wrapper), wrapper);

    wrapper->setVisible(m_controlWidget->expanded());

    if (wrapper->attention()) {
        setCurrentAttentionTray(wrapper);
    }

    connect(wrapper, &FashionTrayWidgetWrapper::attentionChanged, this, &FashionTrayItem::onTrayAttentionChanged, static_cast<Qt::ConnectionType>(Qt::QueuedConnection | Qt::UniqueConnection));

    if (trayWidget->trayTyep() == AbstractTrayWidget::TrayType::SystemTray) {
        SystemTrayItem * sysTrayWidget = static_cast<SystemTrayItem *>(trayWidget);
        connect(sysTrayWidget, &SystemTrayItem::requestWindowAutoHide, this, &FashionTrayItem::requestWindowAutoHide, Qt::UniqueConnection);
        connect(sysTrayWidget, &SystemTrayItem::requestRefershWindowVisible, this, &FashionTrayItem::requestRefershWindowVisible, Qt::UniqueConnection);
    }

    requestResize();
}

void FashionTrayItem::trayWidgetRemoved(AbstractTrayWidget *trayWidget)
{
    bool founded = false;

    for (auto wrapper : m_wrapperList) {
        // found the removed tray
        if (wrapper->absTrayWidget() == trayWidget) {
            // the removed tray is a attention tray
            if (m_currentAttentionTray == wrapper) {
                if (m_controlWidget->expanded()) {
                    m_trayBoxLayout->removeWidget(m_currentAttentionTray);
                } else {
                    m_mainBoxLayout->removeWidget(m_currentAttentionTray);
                }
                m_currentAttentionTray = nullptr;
            } else {
                m_trayBoxLayout->removeWidget(wrapper);
            }
            // do not delete real tray object, just delete it's wrapper object
            // the real tray object should be deleted in TrayPlugin class
            trayWidget->setParent(nullptr);
            wrapper->deleteLater();
            m_wrapperList.removeAll(wrapper);
            founded = true;
            break;
        }
    }

    if (!founded) {
        qDebug() << "Error! can not find the tray widget in fashion tray list" << trayWidget;
    }

    requestResize();
}

void FashionTrayItem::clearTrayWidgets()
{
    QList<QPointer<FashionTrayWidgetWrapper>> mList = m_wrapperList;

    for (auto wrapper : mList) {
        trayWidgetRemoved(wrapper->absTrayWidget());
    }

    m_wrapperList.clear();

    requestResize();
}

void FashionTrayItem::setDockPostion(Dock::Position pos)
{
    m_dockPosistion = pos;

    m_controlWidget->setDockPostion(m_dockPosistion);
    SystemTrayItem::setDockPostion(m_dockPosistion);

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
    if (!isVisible())
        return;

    if (expand) {
        refreshTraysVisible();
    } else {
        // hide all tray widget delay for fold animation
        QTimer::singleShot(350, this, [=] {refreshTraysVisible();});
        requestResize();
    }
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

    for (auto wrapper : m_wrapperList) {
        wrapper->setFixedSize(newSize);
    }

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
            size.setWidth(m_wrapperList.size() * TrayWidgetWidth // 所有插件
                          + TrayWidgetWidth // 控制按钮
                          + SpliterSize * 2 // 两个分隔条
                          + 3 * TraySpace // MainBoxLayout所有space
                          + (m_wrapperList.size() - 1) * TraySpace); // TrayBoxLayout所有space
            size.setHeight(height());
        } else {
            size.setWidth(width());
            size.setHeight(m_wrapperList.size() * TrayWidgetHeight // 所有插件
                          + TrayWidgetHeight // 控制按钮
                          + SpliterSize * 2 // 两个分隔条
                          + 3 * TraySpace // MainBoxLayout所有space
                          + (m_wrapperList.size() - 1) * TraySpace); // TrayBoxLayout所有space
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

int FashionTrayItem::whereToInsert(FashionTrayWidgetWrapper *wrapper) const
{
    int insertIndex = m_trayPlugin->itemSortKey(wrapper->itemKey());
    int firstSystemTrayIndex = 0;

    for (int i = 0; i < m_wrapperList.size(); ++i) {
        if (m_wrapperList.at(i)->absTrayWidget()->trayTyep() != AbstractTrayWidget::TrayType::SystemTray) {
            continue;
        }
        firstSystemTrayIndex = i;
        break;
    }

    if (insertIndex < -1) {
        insertIndex = 0;
    }

    switch (wrapper->absTrayWidget()->trayTyep()) {
    case AbstractTrayWidget::TrayType::ApplicationTray:
        // 应用图标的位置不允许超过第一个系统图标
        if (insertIndex == -1 || insertIndex > firstSystemTrayIndex) {
            insertIndex = firstSystemTrayIndex;
        }
        break;
    case AbstractTrayWidget::TrayType::SystemTray:
        if (insertIndex == -1 || insertIndex > m_wrapperList.size()) {
            insertIndex = m_wrapperList.size();
        } else if (insertIndex == 0) {
            insertIndex = firstSystemTrayIndex;
        } else {
            insertIndex += firstSystemTrayIndex;
        }
        break;
    default:
        insertIndex = 0;
        break;
    }

    return insertIndex;
}

void FashionTrayItem::onTrayAttentionChanged(const bool attention)
{
    // 设置attention为false之后，启动timer，在timer处于Active状态期间不重设attention为true
    if (!attention) {
        m_attentionDelayTimer->start();
    } else if (attention && m_attentionDelayTimer->isActive()) {
        return;
    }

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
    // reset property "FashionTraySize" to notify dock resize
    // DockPluginsController will watch this property
    setProperty("FashionTraySize", sizeHint());
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
    m_trayBoxLayout->insertWidget(whereToInsert(m_currentAttentionTray), m_currentAttentionTray);
    m_currentAttentionTray->setVisible(false);
    m_currentAttentionTray->setAttention(false);
}

void FashionTrayItem::switchAttionTray(FashionTrayWidgetWrapper *attentionWrapper)
{
    if (!m_currentAttentionTray || !attentionWrapper) {
        return;
    }

    m_mainBoxLayout->replaceWidget(m_currentAttentionTray, attentionWrapper);
    m_trayBoxLayout->removeWidget(attentionWrapper);
    m_trayBoxLayout->insertWidget(whereToInsert(m_currentAttentionTray), m_currentAttentionTray);

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

void FashionTrayItem::refreshTraysVisible()
{
    const bool expand = m_controlWidget->expanded();

    if (m_currentAttentionTray) {
        if (expand) {
            m_mainBoxLayout->removeWidget(m_currentAttentionTray);
            m_trayBoxLayout->insertWidget(whereToInsert(m_currentAttentionTray), m_currentAttentionTray);
        }

        m_currentAttentionTray = nullptr;
    }

    for (auto wrapper : m_wrapperList) {
        wrapper->setVisible(expand);
        // reset all tray item attention state
        wrapper->setAttention(false);
    }

    m_attentionDelayTimer->start();

    requestResize();
}
