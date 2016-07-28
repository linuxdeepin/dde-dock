#include "datetimewidget.h"
#include "constants.h"

#include <QApplication>
#include <QPainter>
#include <QDebug>
#include <QSvgRenderer>

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
    const QRect r = rect();

    // draw background
    const QPixmap background = loadSvg(":/icons/resources/icons/background.svg", QSize(perfectIconSize, perfectIconSize));
    const QPoint backgroundOffset = r.center() - background.rect().center();
    painter.drawPixmap(backgroundOffset, background);

    const int bigNumHeight = perfectIconSize / 2.5;
    const int bigNumWidth = double(bigNumHeight) * 8 / 18;
    const int smallNumHeight = bigNumHeight / 2;
    const int smallNumWidth = double(smallNumHeight) * 5 / 9;

    // draw big num 1
    const QString bigNum1Path = QString(":/icons/resources/icons/big%1.svg").arg(currentTimeString[0]);
    const QPixmap bigNum1 = loadSvg(bigNum1Path, QSize(bigNumWidth, bigNumHeight));
    const QPoint bigNum1Offset = backgroundOffset + QPoint(perfectIconSize / 2 - bigNumWidth * 2 + 1, perfectIconSize / 2 - bigNumHeight / 2);
    painter.drawPixmap(bigNum1Offset, bigNum1);

    // draw big num 2
    const QString bigNum2Path = QString(":/icons/resources/icons/big%1.svg").arg(currentTimeString[1]);
    const QPixmap bigNum2 = loadSvg(bigNum2Path, QSize(bigNumWidth, bigNumHeight));
    const QPoint bigNum2Offset = bigNum1Offset + QPoint(bigNumWidth + 1, 0);
    painter.drawPixmap(bigNum2Offset, bigNum2);

    // draw small num 1
    const QString smallNum1Path = QString(":/icons/resources/icons/small%1.svg").arg(currentTimeString[2]);
    const QPixmap smallNum1 = loadSvg(smallNum1Path, QSize(smallNumWidth, smallNumHeight));
    const QPoint smallNum1Offset = bigNum2Offset + QPoint(bigNumWidth + 2, smallNumHeight);
    painter.drawPixmap(smallNum1Offset, smallNum1);

    // draw small num 2
    const QString smallNum2Path = QString(":/icons/resources/icons/small%1.svg").arg(currentTimeString[3]);
    const QPixmap smallNum2 = loadSvg(smallNum2Path, QSize(smallNumWidth, smallNumHeight));
    const QPoint smallNum2Offset = smallNum1Offset + QPoint(smallNumWidth + 1, 0);
    painter.drawPixmap(smallNum2Offset, smallNum2);
}

const QPixmap DatetimeWidget::loadSvg(const QString &fileName, const QSize size)
{
    QPixmap pixmap(size);
    QSvgRenderer renderer(fileName);
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();

    return pixmap;
}
