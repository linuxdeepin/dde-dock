#include "datetimewidget.h"
#include "constants.h"

#include <QApplication>
#include <QPainter>
#include <QDebug>

DatetimeWidget::DatetimeWidget(QWidget *parent)
    : QWidget(parent)
{

}

QSize DatetimeWidget::sizeHint() const
{
    QFontMetrics fm(qApp->font());

    return fm.boundingRect("88:88").size() + QSize(20, 10);
}

void DatetimeWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    QPainter painter(this);

    if (displayMode == Dock::Efficient)
    {
        painter.setPen(Qt::white);
        painter.drawText(rect(), Qt::AlignCenter, "88:88");
    }
    else
    {

    }
}
