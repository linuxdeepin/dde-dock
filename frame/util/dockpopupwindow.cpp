// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dockpopupwindow.h"
#include "imageutil.h"
#include "utils.h"
#include "dbusutil.h"
#include "dockscreen.h"
#include "displaymanager.h"

#include <QScreen>
#include <QApplication>
#include <QDesktopWidget>
#include <QAccessible>
#include <QAccessibleEvent>
#include <QCursor>
#include <QGSettings>

DWIDGET_USE_NAMESPACE

#define DOCK_SCREEN DockScreen::instance()
#define DIS_INS DisplayManager::instance()

DockPopupWindow::DockPopupWindow(QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_model(false)
    , m_eventMonitor(new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus(), this))
    , m_enableMouseRelease(true)
    , m_extendWidget(nullptr)
    , m_lastWidget(nullptr)
{
    setContentsMargins(0, 0, 0, 0);
    m_wmHelper = DWindowManagerHelper::instance();

    setWindowFlags(Qt::ToolTip | Qt::WindowStaysOnTopHint);
    if (Utils::IS_WAYLAND_DISPLAY) {
        setAttribute(Qt::WA_NativeWindow);
        windowHandle()->setProperty("_d_dwayland_window-type", "override");
    } else {
        setAttribute(Qt::WA_InputMethodEnabled, false);
    }

    connect(m_eventMonitor, &XEventMonitor::ButtonPress, this, &DockPopupWindow::onButtonPress);

    if (Utils::IS_WAYLAND_DISPLAY)
        QDBusConnection::sessionBus().connect("com.deepin.dde.lockFront", "/com/deepin/dde/lockFront", "com.deepin.dde.lockFront", "Visible", "b", this, SLOT(hide()));
}

DockPopupWindow::~DockPopupWindow()
{
}

bool DockPopupWindow::model() const
{
    return isVisible() && m_model;
}

QWidget *DockPopupWindow::getContent()
{
    return m_lastWidget;
}

void DockPopupWindow::setContent(QWidget *content)
{
    if (m_lastWidget)
        m_lastWidget->removeEventFilter(this);
    content->installEventFilter(this);

    QAccessibleEvent event(this, QAccessible::NameChanged);
    QAccessible::updateAccessibility(&event);

    if (!content->objectName().trimmed().isEmpty())
        setAccessibleName(content->objectName() + "-popup");

    m_lastWidget = content;
    content->setParent(this);
    content->show();
    resize(content->sizeHint());
}

void DockPopupWindow::setExtendWidget(QWidget *widget)
{
    m_extendWidget = widget;
    connect(widget, &QWidget::destroyed, this, [ this ] { m_extendWidget = nullptr; }, Qt::UniqueConnection);
}

void DockPopupWindow::setPosition(Dock::Position position)
{
    m_position = position;
}

QWidget *DockPopupWindow::extendWidget() const
{
    return m_extendWidget;
}

void DockPopupWindow::show(const QPoint &pos, const bool model)
{
    m_model = model;
    show(pos.x(), pos.y());

    if (!m_eventKey.isEmpty()) {
        m_eventMonitor->UnregisterArea(m_eventKey);
        m_eventKey.clear();
    }

    if (m_model) {
        m_eventKey = m_eventMonitor->RegisterFullScreen();
    }

    blockButtonRelease();
}

void DockPopupWindow::show(const int x, const int y)
{
    QPoint displayPoint;
    m_lastPoint = QPoint(x, y);
    switch (m_position) {
        case Dock::Position::Left:
            displayPoint = m_lastPoint + QPoint(0, -m_lastWidget->height() / 2);
            break;
        case Dock::Position::Right:
            displayPoint = m_lastPoint + QPoint(-m_lastWidget->width(), -m_lastWidget->height() / 2);
            break;
        case Dock::Position::Top:
            displayPoint = m_lastPoint + QPoint(-m_lastWidget->width() / 2, 0);
            break;
        case Dock::Position::Bottom:
            displayPoint = m_lastPoint + QPoint(-m_lastWidget->width() / 2, -m_lastWidget->height());
            break;
    }
    blockButtonRelease();
    QScreen *screen = DIS_INS->screen(DOCK_SCREEN->current());
    if (!screen)
        return;
    QRect screenRect = screen->geometry();
    if (getContent()->width() <= screenRect.width()) {
        displayPoint.setX(qMax(screenRect.x(), displayPoint.x()));
        displayPoint.setX(qMin(screenRect.x() + screenRect.width() - getContent()->width(), displayPoint.x()));
    }
    move(displayPoint);
    resize(m_lastWidget->size());
    DBlurEffectWidget::show();
    activateWindow();
}

void DockPopupWindow::blockButtonRelease()
{
    // 短暂的不处理鼠标release事件，防止出现刚显示又被隐藏的情况
    m_enableMouseRelease = false;
    QTimer::singleShot(10, this, [this] {
        m_enableMouseRelease = true;
    });
}

void DockPopupWindow::hide()
{
    if (!m_eventKey.isEmpty()) {
        m_eventMonitor->UnregisterArea(m_eventKey);
        m_eventKey.clear();
    }

    DBlurEffectWidget::hide();
}

void DockPopupWindow::showEvent(QShowEvent *e)
{
    DBlurEffectWidget::showEvent(e);
    if (Utils::IS_WAYLAND_DISPLAY) {
        Utils::updateCursor(this);
    }

    QTimer::singleShot(1, this, &DockPopupWindow::ensureRaised);
}

void DockPopupWindow::hideEvent(QHideEvent *event)
{
    m_extendWidget = nullptr;
    Dtk::Widget::DBlurEffectWidget::hideEvent(event);
}

void DockPopupWindow::enterEvent(QEvent *e)
{
    DBlurEffectWidget::enterEvent(e);
    if (Utils::IS_WAYLAND_DISPLAY) {
        Utils::updateCursor(this);
    }

    QTimer::singleShot(1, this, &DockPopupWindow::ensureRaised);
}

bool DockPopupWindow::eventFilter(QObject *o, QEvent *e)
{
    if (o != getContent())
        return false;

    switch(e->type()) {
    case QEvent::Resize: {
        // FIXME: ensure position move after global mouse release event
        if (isVisible()) {
            QTimer::singleShot(10, this, [=] {
                // NOTE(sbw): double check is necessary, in this time, the popup maybe already hided.
                if (isVisible())
                    show(m_lastPoint, m_model);
            });
        }
        break;
    }
    case QEvent::WindowDeactivate:
    case QEvent::Hide: {
        this->hide();
        break;
    }
    default:
        break;
    }

    return false;
}

void DockPopupWindow::ensureRaised()
{
    if (isVisible())
        raise();
}

void DockPopupWindow::onButtonPress(int type, int x, int y, const QString &key)
{
    if (!m_enableMouseRelease)
        return;
    QScreen *screen = DIS_INS->screen(DOCK_SCREEN->current());
    if (!screen)
        return;
    QRect screenRect = screen->geometry();
    QRect popupRect(((pos() - screenRect.topLeft()) * qApp->devicePixelRatio() + screenRect.topLeft()), size() * qApp->devicePixelRatio());
    if (popupRect.contains(QPoint(x, y)))
        return;

    if (m_extendWidget) {
        // 计算额外添加的区域，如果鼠标的点击点在额外的区域内，也无需隐藏
        QPoint extendPoint = m_extendWidget->mapToGlobal(QPoint(0, 0));
        QRect extendRect(((extendPoint - screenRect.topLeft()) * qApp->devicePixelRatio() + screenRect.topLeft()), m_extendWidget->size() * qApp->devicePixelRatio());
        if (extendRect.contains(QPoint(x, y)))
            return;
    }

    // if there is something focus on widget, return
    if (auto focus = qApp->focusWidget()) {
        auto className = QString(focus->metaObject()->className());
        qDebug() << "Find focused widget, focus className is" << className;
        if (className == "QLineEdit") {
            qDebug() << "PopupWindow window will not be hidden";
            return;
        }
    }

    emit accept();
    hide();
}

PopupSwitchWidget::PopupSwitchWidget(QWidget *parent)
    : QWidget(parent)
    , m_containerLayout(new QVBoxLayout(this))
    , m_topWidget(nullptr)
{
    m_containerLayout->setContentsMargins(0, 0, 0, 0);
    m_containerLayout->setSpacing(0);
}

PopupSwitchWidget::~PopupSwitchWidget()
{
}

void PopupSwitchWidget::pushWidget(QWidget *widget)
{
    // 首先将界面其他的窗体移除
    for (int i = m_containerLayout->count() - 1; i >= 0; i--) {
        QLayoutItem *item = m_containerLayout->itemAt(i);
        item->widget()->removeEventFilter(this);
        item->widget()->hide();
        m_containerLayout->removeItem(item);
    }
    m_topWidget = widget;
    setFixedSize(widget->size());
    widget->installEventFilter(this);
    m_containerLayout->addWidget(widget);
    widget->show();
}

bool PopupSwitchWidget::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_topWidget && event->type() == QEvent::Resize) {
        setFixedSize(m_topWidget->size());
    }

    return QWidget::eventFilter(watched, event);
}
