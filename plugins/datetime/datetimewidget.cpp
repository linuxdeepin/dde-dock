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

#include "datetimewidget.h"
#include "constants.h"

#include <QApplication>
#include <QPainter>
#include <QDebug>
#include <QSvgRenderer>
#include <QMouseEvent>

#define PLUGIN_STATE_KEY    "enable"
#define SHOW_DATE_MIN_HEIGHT 45

DatetimeWidget::DatetimeWidget(QWidget *parent)
    : QWidget(parent)
{

}

void DatetimeWidget::set24HourFormat(const bool value)
{
    if (m_24HourFormat == value) {
        return;
    }

    m_24HourFormat = value;

    update();

    if (isVisible()) {
        emit requestUpdateGeometry();
    }
}

QSize DatetimeWidget::sizeHint() const
{
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    QFontMetrics fm(qApp->font());

    QString timeString;
    if (m_24HourFormat)
        timeString = "88:88";
    else
        timeString = "88:88 A.A.";

    if (displayMode == Dock::Fashion && rect().height() > SHOW_DATE_MIN_HEIGHT)
        timeString += "\n8888-88-88";

    return fm.size(Qt::TextExpandTabs, timeString) + QSize(20, 20);
}

void DatetimeWidget::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);
}

void DatetimeWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const auto ratio = devicePixelRatioF();
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    const QDateTime current = QDateTime::currentDateTime();

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);

    QString format;
    if (m_24HourFormat)
        format = "hh:mm";
    else
        format = "hh:mm AP";

    if (displayMode == Dock::Fashion && rect().height() > SHOW_DATE_MIN_HEIGHT)
        format += "\nyyyy-MM-dd";

    painter.setPen(Qt::white);
    painter.drawText(rect(), Qt::AlignCenter, current.toString(format));
}

const QPixmap DatetimeWidget::loadSvg(const QString &fileName, const QSize size)
{
    const auto ratio = devicePixelRatioF();

    QPixmap pixmap(size * ratio);
    QSvgRenderer renderer(fileName);
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();

    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}
