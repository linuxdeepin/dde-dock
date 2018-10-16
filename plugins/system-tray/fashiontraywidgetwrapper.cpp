#include "fashiontraywidgetwrapper.h"

#include <QPainter>
#include <QDebug>

FashionTrayWidgetWrapper::FashionTrayWidgetWrapper(AbstractTrayWidget *absTrayWidget, QWidget *parent)
    : QWidget(parent),
      m_absTrayWidget(absTrayWidget),
      m_layout(new QVBoxLayout(this)),
      m_attention(false)

{
    m_layout->setSpacing(0);
    m_layout->setMargin(0);
    m_layout->setContentsMargins(0, 0, 0, 0);

    m_layout->addWidget(m_absTrayWidget);

    setLayout(m_layout);

    connect(m_absTrayWidget, &AbstractTrayWidget::iconChanged, this, &FashionTrayWidgetWrapper::onTrayWidgetIconChanged);
    connect(m_absTrayWidget, &AbstractTrayWidget::clicked, this, &FashionTrayWidgetWrapper::onTrayWidgetClicked);
}

AbstractTrayWidget *FashionTrayWidgetWrapper::absTrayWidget() const
{
    return m_absTrayWidget;
}

void FashionTrayWidgetWrapper::paintEvent(QPaintEvent *event)
{
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);
    painter.setPen(QColor(QColor::fromRgb(40, 40, 40)));
    painter.setBrush(QColor(QColor::fromRgb(40, 40, 40)));
    painter.setOpacity(0.5);

    painter.drawRoundRect(rect());
}

void FashionTrayWidgetWrapper::onTrayWidgetIconChanged()
{
    setAttention(true);
}

void FashionTrayWidgetWrapper::onTrayWidgetClicked()
{
    setAttention(false);
}

bool FashionTrayWidgetWrapper::attention() const
{
    return m_attention;
}

void FashionTrayWidgetWrapper::setAttention(bool attention)
{
    m_attention = attention;

    Q_EMIT attentionChanged(m_attention);
}
