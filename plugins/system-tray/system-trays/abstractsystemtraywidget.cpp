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

#include "abstractsystemtraywidget.h"
#include "dbus/dbusmenu.h"

#include <QProcess>
#include <QDebug>

#include <xcb/xproto.h>

Dock::Position AbstractSystemTrayWidget::DockPosition = Dock::Position::Top;
QPointer<DockPopupWindow> AbstractSystemTrayWidget::PopupWindow = nullptr;

AbstractSystemTrayWidget::AbstractSystemTrayWidget(QWidget *parent)
    : AbstractTrayWidget(parent),
      m_popupShown(false),
      m_popupTipsDelayTimer(new QTimer(this)),
      m_popupAdjustDelayTimer(new QTimer(this)),
      m_menuManagerInter(new DBusMenuManager(this))
{
    if (PopupWindow.isNull())
    {
        DockPopupWindow *arrowRectangle = new DockPopupWindow(nullptr);
        arrowRectangle->setShadowBlurRadius(20);
        arrowRectangle->setRadius(6);
        arrowRectangle->setShadowYOffset(2);
        arrowRectangle->setShadowXOffset(0);
        arrowRectangle->setArrowWidth(18);
        arrowRectangle->setArrowHeight(10);
        PopupWindow = arrowRectangle;
    }

    m_popupTipsDelayTimer->setInterval(500);
    m_popupTipsDelayTimer->setSingleShot(true);

    m_popupAdjustDelayTimer->setInterval(10);
    m_popupAdjustDelayTimer->setSingleShot(true);

    connect(m_popupTipsDelayTimer, &QTimer::timeout, this, &AbstractSystemTrayWidget::showHoverTips);
    connect(m_popupAdjustDelayTimer, &QTimer::timeout, this, &AbstractSystemTrayWidget::updatePopupPosition, Qt::QueuedConnection);
}

AbstractSystemTrayWidget::~AbstractSystemTrayWidget()
{
    if (m_popupShown)
        popupWindowAccept();
}

void AbstractSystemTrayWidget::sendClick(uint8_t mouseButton, int x, int y)
{
    switch (mouseButton) {
    case XCB_BUTTON_INDEX_1: {
        showPopupApplet(trayPopupApplet());
        QProcess::startDetached(trayClickCommand());
        break;
    }
    case XCB_BUTTON_INDEX_2:
        break;
    case XCB_BUTTON_INDEX_3: {
        showContextMenu();
        break;
    }
    default:
        qDebug() << "unknown mouse button key";
        break;
    }
}

bool AbstractSystemTrayWidget::event(QEvent *event)
{
    if (m_popupShown)
    {
        switch (event->type())
        {
        case QEvent::Paint:
            if (!m_popupAdjustDelayTimer->isActive())
                m_popupAdjustDelayTimer->start();
            break;
        default:;
        }
    }

    return AbstractTrayWidget::event(event);
}

void AbstractSystemTrayWidget::enterEvent(QEvent *event)
{
    m_popupTipsDelayTimer->start();
    update();

    AbstractTrayWidget::enterEvent(event);
}

void AbstractSystemTrayWidget::leaveEvent(QEvent *event)
{
    m_popupTipsDelayTimer->stop();

    // auto hide if popup is not model window
    if (m_popupShown && !PopupWindow->model())
        hidePopup();

    update();

    AbstractTrayWidget::leaveEvent(event);
}

void AbstractSystemTrayWidget::mousePressEvent(QMouseEvent *event)
{
    m_popupTipsDelayTimer->stop();
    hideNonModel();

    AbstractTrayWidget::mousePressEvent(event);
}

const QPoint AbstractSystemTrayWidget::popupMarkPoint() const
{
    QPoint p(topleftPoint());

    const QRect r = rect();
    const int offset = 2;
    switch (DockPosition)
    {
    case Dock::Position::Top:       p += QPoint(r.width() / 2, r.height() + offset);      break;
    case Dock::Position::Bottom:    p += QPoint(r.width() / 2, 0 - offset);               break;
    case Dock::Position::Left:      p += QPoint(r.width() + offset, r.height() / 2);      break;
    case Dock::Position::Right:     p += QPoint(0 - offset, r.height() / 2);              break;
    }

    return p;
}

const QPoint AbstractSystemTrayWidget::topleftPoint() const
{
    QPoint p;
    const QWidget *w = this;
    do {
        p += w->pos();
        w = qobject_cast<QWidget *>(w->parent());
    } while (w);

    return p;
}

void AbstractSystemTrayWidget::hidePopup()
{
    m_popupTipsDelayTimer->stop();
    m_popupAdjustDelayTimer->stop();
    m_popupShown = false;
    PopupWindow->hide();

    emit PopupWindow->accept();
    emit requestWindowAutoHide(true);
}

void AbstractSystemTrayWidget::hideNonModel()
{
    // auto hide if popup is not model window
    if (m_popupShown && !PopupWindow->model())
        hidePopup();
}

void AbstractSystemTrayWidget::popupWindowAccept()
{
    if (!PopupWindow->isVisible())
        return;

    disconnect(PopupWindow.data(), &DockPopupWindow::accept, this, &AbstractSystemTrayWidget::popupWindowAccept);

    hidePopup();
}

void AbstractSystemTrayWidget::showPopupApplet(QWidget * const applet)
{
    // another model popup window already exists
    if (PopupWindow->model())
        return;

    if (!applet) {
        return;
    }

    showPopupWindow(applet, true);
}

void AbstractSystemTrayWidget::showPopupWindow(QWidget * const content, const bool model)
{
    m_popupShown = true;
    m_lastPopupWidget = content;

    if (model)
        emit requestWindowAutoHide(false);

    DockPopupWindow *popup = PopupWindow.data();
    QWidget *lastContent = popup->getContent();
    if (lastContent)
        lastContent->setVisible(false);

    switch (DockPosition)
    {
    case Dock::Position::Top:   popup->setArrowDirection(DockPopupWindow::ArrowTop);     break;
    case Dock::Position::Bottom:popup->setArrowDirection(DockPopupWindow::ArrowBottom);  break;
    case Dock::Position::Left:  popup->setArrowDirection(DockPopupWindow::ArrowLeft);    break;
    case Dock::Position::Right: popup->setArrowDirection(DockPopupWindow::ArrowRight);   break;
    }
    popup->resize(content->sizeHint());
    popup->setContent(content);

    const QPoint p = popupMarkPoint();
    if (!popup->isVisible())
        QMetaObject::invokeMethod(popup, "show", Qt::QueuedConnection, Q_ARG(QPoint, p), Q_ARG(bool, model));
    else
        popup->show(p, model);

    connect(popup, &DockPopupWindow::accept, this, &AbstractSystemTrayWidget::popupWindowAccept, Qt::UniqueConnection);
}

void AbstractSystemTrayWidget::showHoverTips()
{
    // another model popup window already exists
    if (PopupWindow->model())
        return;

    // if not in geometry area
    const QRect r(topleftPoint(), size());
    if (!r.contains(QCursor::pos()))
        return;

    QWidget * const content = trayTipsWidget();
    if (!content)
        return;

    showPopupWindow(content);
}

void AbstractSystemTrayWidget::showContextMenu()
{
    const QString menuJson = contextMenu();
    if (menuJson.isEmpty())
        return;

    QDBusPendingReply<QDBusObjectPath> result = m_menuManagerInter->RegisterMenu();

    result.waitForFinished();
    if (result.isError())
    {
        qWarning() << result.error();
        return;
    }

    const QPoint p = popupMarkPoint();

    QJsonObject menuObject;
    menuObject.insert("x", QJsonValue(p.x()));
    menuObject.insert("y", QJsonValue(p.y()));
    menuObject.insert("isDockMenu", QJsonValue(true));
    menuObject.insert("menuJsonContent", QJsonValue(menuJson));

    switch (DockPosition)
    {
    case Dock::Position::Top:       menuObject.insert("direction", "top");      break;
    case Dock::Position::Bottom:    menuObject.insert("direction", "bottom");   break;
    case Dock::Position::Left:      menuObject.insert("direction", "left");     break;
    case Dock::Position::Right:     menuObject.insert("direction", "right");    break;
    }

    const QDBusObjectPath path = result.argumentAt(0).value<QDBusObjectPath>();
    DBusMenu *menuInter = new DBusMenu(path.path(), this);

    connect(menuInter, &DBusMenu::ItemInvoked, this, &AbstractSystemTrayWidget::invokedMenuItem);
    connect(menuInter, &DBusMenu::ItemInvoked, menuInter, &DBusMenu::deleteLater);
    connect(menuInter, &DBusMenu::MenuUnregistered, this, &AbstractSystemTrayWidget::onContextMenuAccepted, Qt::QueuedConnection);

    menuInter->ShowMenu(QString(QJsonDocument(menuObject).toJson()));

    hidePopup();
    emit requestWindowAutoHide(false);
}

void AbstractSystemTrayWidget::onContextMenuAccepted()
{
    emit requestRefershWindowVisible();
    emit requestWindowAutoHide(true);
}

void AbstractSystemTrayWidget::updatePopupPosition()
{
    Q_ASSERT(sender() == m_popupAdjustDelayTimer);

    if (!m_popupShown || !PopupWindow->model())
        return;

    if (PopupWindow->getContent() != m_lastPopupWidget.data())
        return popupWindowAccept();

    const QPoint p = popupMarkPoint();
    PopupWindow->show(p, PopupWindow->model());
}
