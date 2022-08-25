/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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
#include "dragwidget.h"
#include "utils.h"
#include "constants.h"

#include <QCoreApplication>
#include <QMouseEvent>

DragWidget::DragWidget(QWidget *parent)
    : QWidget(parent)
    , m_dragStatus(false)
{
    setObjectName("DragWidget");
}

bool DragWidget::isDraging() const
{
    return m_dragStatus;
}

void DragWidget::onTouchMove(double scaleX, double scaleY)
{
    Q_UNUSED(scaleX);
    Q_UNUSED(scaleY);

    static QPoint lastPos;
    QPoint curPos = QCursor::pos();
    if (lastPos == curPos) {
        return;
    }
    lastPos = curPos;
    qApp->postEvent(this, new QMouseEvent(QEvent::MouseMove, mapFromGlobal(curPos)
                                                  , QPoint(), curPos, Qt::LeftButton, Qt::LeftButton
                                          , Qt::NoModifier, Qt::MouseEventSynthesizedByApplication));
}

void DragWidget::mousePressEvent(QMouseEvent *event)
{
    // qt转发的触屏按下信号不进行响应
    if (event->source() == Qt::MouseEventSynthesizedByQt)
        return;

    if (event->button() == Qt::LeftButton) {
        m_resizePoint = event->globalPos();
        m_dragStatus = true;
        this->grabMouse();
    }
}

void DragWidget::mouseMoveEvent(QMouseEvent *)
{
    if (m_dragStatus) {
        QPoint offset = QPoint(QCursor::pos() - m_resizePoint);
        emit dragPointOffset(offset);
    }
}

void DragWidget::mouseReleaseEvent(QMouseEvent *)
{
    if (!m_dragStatus)
        return;

    m_dragStatus =  false;
    releaseMouse();
    emit dragFinished();
}

void DragWidget::enterEvent(QEvent *)
{
    if (Utils::IS_WAYLAND_DISPLAY)
        updateCursor();
    else
        QApplication::setOverrideCursor(cursor());
}

void DragWidget::leaveEvent(QEvent *)
{
    QApplication::setOverrideCursor(Qt::ArrowCursor);
}

void DragWidget::updateCursor()
{
    QString theme = Utils::SettingValue("com.deepin.xsettings", "/com/deepin/xsettings/", "gtk-cursor-theme-name", "bloom").toString();
    int cursorSize = Utils::SettingValue("com.deepin.xsettings", "/com/deepin/xsettings/", "gtk-cursor-theme-size", 24).toInt();
    Dock::Position position = static_cast<Dock::Position>(qApp->property("position").toInt());

    static QString lastTheme;
    static int lastPosition = -1;
    static int lastCursorSize = -1;
    if (theme != lastTheme || position != lastPosition || cursorSize != lastCursorSize) {
        lastTheme = theme;
        lastPosition = position;
        lastCursorSize = cursorSize;
        const char* cursorName = (position == Dock::Position::Bottom || position == Dock::Position::Top) ? "v_double_arrow" : "h_double_arrow";
        QCursor *newCursor = ImageUtil::loadQCursorFromX11Cursor(theme.toStdString().c_str(), cursorName, cursorSize);
        if (!newCursor)
            return;

        setCursor(*newCursor);
        static QCursor *lastCursor = nullptr;
        if (lastCursor)
            delete lastCursor;

        lastCursor = newCursor;
    }
}
