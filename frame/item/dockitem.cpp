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

#include "dockitem.h"
#include "dbus/dbusmenu.h"
#include "dbus/dbusmenumanager.h"
#include "components/hoverhighlighteffect.h"

#include <QMouseEvent>
#include <QJsonObject>

Position DockItem::DockPosition = Position::Top;
DisplayMode DockItem::DockDisplayMode = DisplayMode::Efficient;
QPointer<DockPopupWindow> DockItem::PopupWindow(nullptr);

DockItem::DockItem(QWidget *parent)
    : QWidget(parent),
      m_hover(false),
      m_popupShown(false),

      m_hoverEffect(new HoverHighlightEffect(this)),

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

    setGraphicsEffect(m_hoverEffect);

    connect(m_popupTipsDelayTimer, &QTimer::timeout, this, &DockItem::showHoverTips);
    connect(m_popupAdjustDelayTimer, &QTimer::timeout, this, &DockItem::updatePopupPosition, Qt::QueuedConnection);
}

DockItem::~DockItem()
{
    if (m_popupShown)
        popupWindowAccept();
}

void DockItem::setDockPosition(const Position side)
{
    DockPosition = side;
}

void DockItem::setDockDisplayMode(const DisplayMode mode)
{
    DockDisplayMode = mode;
}

bool DockItem::event(QEvent *event)
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

    return QWidget::event(event);
}

void DockItem::updatePopupPosition()
{
    Q_ASSERT(sender() == m_popupAdjustDelayTimer);

    if (!m_popupShown || !PopupWindow->isVisible())
        return;

    const QPoint p = popupMarkPoint();
    PopupWindow->show(p, PopupWindow->model());
}

void DockItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);
}

void DockItem::mousePressEvent(QMouseEvent *e)
{
    m_popupTipsDelayTimer->stop();
    hideNonModel();

    if (e->button() == Qt::RightButton)
        return showContextMenu();
}

void DockItem::enterEvent(QEvent *e)
{
    m_hover = true;
    m_hoverEffect->setHighlighting(true);
    m_popupTipsDelayTimer->start();

    update();

    return QWidget::enterEvent(e);
}

void DockItem::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    m_hover = false;
    m_hoverEffect->setHighlighting(false);
    m_popupTipsDelayTimer->stop();

    // auto hide if popup is not model window
    if (m_popupShown && !PopupWindow->model())
        hidePopup();

    update();
}

const QRect DockItem::perfectIconRect() const
{
    const QRect itemRect = rect();
    const int iconSize = std::min(itemRect.width(), itemRect.height()) * 0.8;

    QRect iconRect;
    iconRect.setWidth(iconSize);
    iconRect.setHeight(iconSize);
    iconRect.moveTopLeft(itemRect.center() - iconRect.center());

    return iconRect;
}

void DockItem::showContextMenu()
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
    case Top:       menuObject.insert("direction", "top");      break;
    case Bottom:    menuObject.insert("direction", "bottom");   break;
    case Left:      menuObject.insert("direction", "left");     break;
    case Right:     menuObject.insert("direction", "right");    break;
    }

    const QDBusObjectPath path = result.argumentAt(0).value<QDBusObjectPath>();
    DBusMenu *menuInter = new DBusMenu(path.path(), this);

    connect(menuInter, &DBusMenu::ItemInvoked, this, &DockItem::invokedMenuItem);
    connect(menuInter, &DBusMenu::MenuUnregistered, this, [=] {
        emit requestRefershWindowVisible();
        emit requestWindowAutoHide(true);
        menuInter->deleteLater();
    });

    menuInter->ShowMenu(QString(QJsonDocument(menuObject).toJson()));

    hidePopup();
    emit requestWindowAutoHide(false);
}

void DockItem::showHoverTips()
{
    // another model popup window is alread exists
    if (PopupWindow->isVisible() && PopupWindow->model())
        return;

    // if not in geometry area
    const QRect r(topleftPoint(), size());
    if (!r.contains(QCursor::pos()))
        return;

    QWidget * const content = popupTips();
    if (!content)
        return;

    showPopupWindow(content);
}

void DockItem::showPopupWindow(QWidget * const content, const bool model)
{
    m_popupShown = true;

    if (model)
        emit requestWindowAutoHide(false);

    DockPopupWindow *popup = PopupWindow.data();
    QWidget *lastContent = popup->getContent();
    if (lastContent)
        lastContent->setVisible(false);

    switch (DockPosition)
    {
    case Top:   popup->setArrowDirection(DockPopupWindow::ArrowTop);     break;
    case Bottom:popup->setArrowDirection(DockPopupWindow::ArrowBottom);  break;
    case Left:  popup->setArrowDirection(DockPopupWindow::ArrowLeft);    break;
    case Right: popup->setArrowDirection(DockPopupWindow::ArrowRight);   break;
    }
    popup->resize(content->sizeHint());
    popup->setContent(content);

    const QPoint p = popupMarkPoint();
    if (!popup->isVisible())
        QMetaObject::invokeMethod(popup, "show", Qt::QueuedConnection, Q_ARG(QPoint, p), Q_ARG(bool, model));
    else
        popup->show(p, model);

    connect(popup, &DockPopupWindow::accept, this, &DockItem::popupWindowAccept, Qt::UniqueConnection);
}

void DockItem::popupWindowAccept()
{
    if (!PopupWindow->isVisible())
        return;

    disconnect(PopupWindow.data(), &DockPopupWindow::accept, this, &DockItem::popupWindowAccept);

    hidePopup();
}

void DockItem::showPopupApplet(QWidget * const applet)
{
    // another model popup window is alread exists
    if (PopupWindow->isVisible() && PopupWindow->model())
        return;

    showPopupWindow(applet, true);
}

void DockItem::invokedMenuItem(const QString &itemId, const bool checked)
{
    Q_UNUSED(itemId)
    Q_UNUSED(checked)
}

const QString DockItem::contextMenu() const
{
    return QString();
}

QWidget *DockItem::popupTips()
{
    return nullptr;
}

const QPoint DockItem::popupMarkPoint() const
{
    QPoint p(topleftPoint());

    const QRect r = rect();
    const int offset = 2;
    switch (DockPosition)
    {
    case Top:       p += QPoint(r.width() / 2, r.height() + offset);      break;
    case Bottom:    p += QPoint(r.width() / 2, 0 - offset);               break;
    case Left:      p += QPoint(r.width() + offset, r.height() / 2);      break;
    case Right:     p += QPoint(0 - offset, r.height() / 2);              break;
    }

    return p;
}

const QPoint DockItem::topleftPoint() const
{
    QPoint p;
    const QWidget *w = this;
    do {
        p += w->pos();
        w = qobject_cast<QWidget *>(w->parent());
    } while (w);

    return p;
}

void DockItem::hidePopup()
{
    m_popupTipsDelayTimer->stop();
    m_popupShown = false;
    PopupWindow->hide();

    emit PopupWindow->accept();
    emit requestWindowAutoHide(true);
}

void DockItem::hideNonModel()
{
    // auto hide if popup is not model window
    if (m_popupShown && !PopupWindow->model())
        hidePopup();
}
