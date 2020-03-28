#include "tipswidget.h"

#include <QApplication>
#include <QPainter>

TipsWidget::TipsWidget(QWidget *parent) : QFrame(parent)
{
    connect(qApp, &QApplication::fontChanged, this, [=] {
         setText(m_text);
     });
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
    painter.setPen(QPen(palette().brightText(), 1));

    QTextOption option;
    option.setAlignment(Qt::AlignCenter);
    painter.drawText(rect(), m_text, option);
}
