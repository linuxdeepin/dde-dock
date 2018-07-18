#include "tipswidget.h"

#include <QPainter>

TipsWidget::TipsWidget(QWidget *parent) : QFrame(parent)
{

}

void TipsWidget::setText(const QString &text)
{
    m_text = text;

    setFixedSize(fontMetrics().width(text) + 6, fontMetrics().height());

    update();
}

void TipsWidget::paintEvent(QPaintEvent *event)
{
    QFrame::paintEvent(event);

    QPainter painter(this);

    QPen pen(Qt::white);
    painter.setPen(pen);

    QTextOption option;
    option.setAlignment(Qt::AlignCenter);
    painter.drawText(rect(), m_text, option);
}
