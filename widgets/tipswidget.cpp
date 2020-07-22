/*
 * Copyright (C) 2011 ~ 2018 uniontech Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
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
#include "tipswidget.h"

#include <QPainter>
#include <QAccessible>
#include <QTextDocument>
namespace Dock{
TipsWidget::TipsWidget(QWidget *parent)
    : QFrame(parent)
    , m_width(0)
    , m_type(SingleLine)
{

}

void TipsWidget::setText(const QString &text)
{
    m_type = TipsWidget::SingleLine;
    // 如果传递的是富文本，获取富文本中的纯文本内容进行显示
    QTextDocument document;
    document.setHtml(text);
    // 同时去掉两边的空白信息，例如qBittorrent的提示
    m_text = document.toPlainText().simplified();

    setFixedSize(fontMetrics().width(m_text) + 6, fontMetrics().height());

    update();

#ifndef QT_NO_ACCESSIBILITY
    if (accessibleName().isEmpty()) {
        QAccessibleEvent event(this, QAccessible::NameChanged);
        QAccessible::updateAccessibility(&event);
    }
#endif
}

void TipsWidget::setTextList(const QStringList &textList)
{
    m_type = TipsWidget::MultiLine;
    m_textList = textList;

    int maxLength = 0;
    int k = fontMetrics().height() * m_textList.size();
    setFixedHeight(k);
    for (QString text : m_textList) {
        int fontLength = fontMetrics().width(text) + 6;
        maxLength = maxLength > fontLength ? maxLength : fontLength;
    }
    m_width = maxLength;
    setFixedWidth(maxLength);

    update();
}

void TipsWidget::paintEvent(QPaintEvent *event)
{
    QFrame::paintEvent(event);

    QPainter painter(this);
    painter.setPen(QPen(palette().brightText(), 1));
    QTextOption option;
    int fontHeight = fontMetrics().height();
    option.setAlignment(Qt::AlignCenter);

    switch (m_type) {
    case SingleLine: {
        painter.drawText(rect(), m_text, option);
    }
        break;
    case MultiLine: {
        int y = 0;
        if (m_textList.size() != 1)
            option.setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
        for (QString text : m_textList) {
            painter.drawText(QRect(0, y, m_width, fontHeight), text, option);
            y += fontHeight;
        }
    }
        break;
    }
}

bool TipsWidget::event(QEvent *event)
{
    if (event->type() == QEvent::FontChange) {
        if (m_type == SingleLine) {
            if (!m_text.trimmed().isEmpty()) {
                 setFixedSize(fontMetrics().width(m_text) + 6, fontMetrics().height());
                 update();
            }
        } else {
            if (m_textList.size() > 0) {
                int maxLength = 0;
                setFixedHeight(fontMetrics().height() * m_textList.size());
                for (QString text : m_textList) {
                    int fontLength = fontMetrics().width(text) + 6;
                    maxLength = qMax(maxLength,fontLength);
                }
                m_width = maxLength;
                setFixedWidth(maxLength);
                update();
            }
        }
    }
    return QFrame::event(event);
}
}
