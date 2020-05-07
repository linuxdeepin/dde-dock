#include "tipswidget.h"

#include <QPainter>

TipsWidget::TipsWidget(QWidget *parent) : QFrame(parent)
{

}

void TipsWidget::setText(const QString &text)
{
    m_type = TipsWidget::SingleLine;
    m_text = text;

    setFixedSize(fontMetrics().width(text) + 6, fontMetrics().height());

    update();
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
