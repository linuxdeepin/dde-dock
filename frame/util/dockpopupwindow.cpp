// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dockpopupwindow.h"
#include "imageutil.h"
#include "utils.h"
#include "dbusutil.h"

#include <QScreen>
#include <QApplication>
#include <QDesktopWidget>
#include <QAccessible>
#include <QAccessibleEvent>
#include <QCursor>
#include <QGSettings>

DWIDGET_USE_NAMESPACE

DockPopupWindow::DockPopupWindow(QWidget *parent)
    : DArrowRectangle(ArrowBottom, parent)
    , m_model(false)
    , m_eventMonitor(new XEventMonitor(xEventMonitorService, xEventMonitorPath, QDBusConnection::sessionBus(), this))
    , m_enableMouseRelease(true)
    , m_extendWidget(nullptr)
{
    setMargin(0);
    m_wmHelper = DWindowManagerHelper::instance();

    compositeChanged();

    setWindowFlags(Qt::X11BypassWindowManagerHint | Qt::WindowStaysOnTopHint | Qt::WindowDoesNotAcceptFocus);
    if (Utils::IS_WAYLAND_DISPLAY) {
        setAttribute(Qt::WA_NativeWindow);
        windowHandle()->setProperty("_d_dwayland_window-type", "override");
    } else {
        setAttribute(Qt::WA_InputMethodEnabled, false);
    }

    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, this, &DockPopupWindow::compositeChanged);
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

void DockPopupWindow::setContent(QWidget *content)
{
    QWidget *lastWidget = getContent();
    if (lastWidget)
        lastWidget->removeEventFilter(this);
    content->installEventFilter(this);

    QAccessibleEvent event(this, QAccessible::NameChanged);
    QAccessible::updateAccessibility(&event);

    if (!content->objectName().trimmed().isEmpty())
        setAccessibleName(content->objectName() + "-popup");

    DArrowRectangle::setContent(content);
}

void DockPopupWindow::setExtendWidget(QWidget *widget)
{
    m_extendWidget = widget;
    connect(widget, &QWidget::destroyed, this, [ this ] { m_extendWidget = nullptr; }, Qt::UniqueConnection);
}

QWidget *DockPopupWindow::extengWidget() const
{
    return m_extendWidget;
}

void DockPopupWindow::show(const QPoint &pos, const bool model)
{
    m_model = model;
    m_lastPoint = pos;

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
    m_lastPoint = QPoint(x, y);
    blockButtonRelease();

    DArrowRectangle::show(x, y);
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

    DArrowRectangle::hide();
}

void DockPopupWindow::showEvent(QShowEvent *e)
{
    DArrowRectangle::showEvent(e);
    if (Utils::IS_WAYLAND_DISPLAY) {
        Utils::updateCursor(this);
    }

    QTimer::singleShot(1, this, &DockPopupWindow::ensureRaised);
}

void DockPopupWindow::hideEvent(QHideEvent *event)
{
    m_extendWidget = nullptr;
    Dtk::Widget::DArrowRectangle::hideEvent(event);
}

void DockPopupWindow::enterEvent(QEvent *e)
{
    DArrowRectangle::enterEvent(e);
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
    case QEvent::Hide: {
        this->hide();
        break;
    }
    default:
        break;
    }

    return false;
}

void DockPopupWindow::compositeChanged()
{
    if (m_wmHelper->hasComposite())
        setBorderColor(QColor(255, 255, 255, 255 * 0.05));
    else
        setBorderColor(QColor("#2C3238"));
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

    QRect popupRect(pos() * qApp->devicePixelRatio(), size() * qApp->devicePixelRatio()) ;
    if (popupRect.contains(x, y))
        return;

    if (m_extendWidget) {
        // 计算额外添加的区域，如果鼠标的点击点在额外的区域内，也无需隐藏
        QPoint extendPoint = m_extendWidget->mapToGlobal(QPoint(0, 0));
        QRect extendRect(extendPoint * qApp->devicePixelRatio(), m_extendWidget->size() * qApp->devicePixelRatio());
        if (extendRect.contains(QPoint(x, y)))
            return;
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
