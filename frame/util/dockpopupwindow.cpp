/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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
    if (o != getContent() || e->type() != QEvent::Resize)
        return false;

    // FIXME: ensure position move after global mouse release event
    if (isVisible()) {
        QTimer::singleShot(10, this, [=] {
            // NOTE(sbw): double check is necessary, in this time, the popup maybe already hided.
            if (isVisible())
                show(m_lastPoint, m_model);
        });
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

    QRect popupRect(pos(), size());
    if (popupRect.contains(x, y))
        return;

    emit accept();
    hide();
}
