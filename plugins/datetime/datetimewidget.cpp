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
    const QDateTime current = QDateTime::currentDateTime();

    QPainter painter(this);

    if (displayMode == Dock::Efficient)
    {
        painter.setPen(Qt::white);
        painter.drawText(rect(), Qt::AlignCenter, current.toString(tr("HH:mm")));
        return;
    }

    // draw fashion mode datetime plugin
    const int perfectIconSize = qMin(width(), height()) * 0.8;
    const QString currentTimeString = current.toString("HHmmss");

    QPixmap pixmap(":/icons/resources/icons/panel.png");
    QPainter itemPainter(&pixmap);
    itemPainter.drawPixmap(8, 15, QPixmap(QString(":/icons/resources/icons/big%1.png").arg(currentTimeString[0])));
    itemPainter.drawPixmap(18, 15, QPixmap(QString(":/icons/resources/icons/big%1.png").arg(currentTimeString[1])));
    itemPainter.drawPixmap(28, 24, QPixmap(QString(":/icons/resources/icons/small%1.png").arg(currentTimeString[2])));
    itemPainter.drawPixmap(34, 24, QPixmap(QString(":/icons/resources/icons/small%1.png").arg(currentTimeString[3])));
    itemPainter.end();

    pixmap = pixmap.scaled(perfectIconSize, perfectIconSize, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    painter.drawPixmap(rect().center() - pixmap.rect().center(), pixmap, pixmap.rect());
}
