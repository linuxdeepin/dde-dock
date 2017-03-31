#include "datetimewidget.h"
#include "constants.h"

#include <QApplication>
#include <QPainter>
#include <QDebug>
#include <QSvgRenderer>
#include <QMouseEvent>

DatetimeWidget::DatetimeWidget(QWidget *parent)
    : QWidget(parent),

    m_settings("deepin", "dde-dock-datetime"),

    m_24HourFormat(m_settings.value("24HourFormat").toBool())
{

}

void DatetimeWidget::toggleHourFormat()
{
    m_24HourFormat = !m_24HourFormat;

    m_settings.setValue("24HourFormat", m_24HourFormat);

    m_cachedTime.clear();
    update();
}

QSize DatetimeWidget::sizeHint() const
{
    QFontMetrics fm(qApp->font());

    return fm.boundingRect("88:88").size() + QSize(20, 10);
}

void DatetimeWidget::resizeEvent(QResizeEvent *e)
{
    m_cachedTime.clear();

    QWidget::resizeEvent(e);
}

void DatetimeWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const QDateTime current = QDateTime::currentDateTime();

    QPainter painter(this);

    if (displayMode == Dock::Efficient)
    {
        const QString text = current.toString(m_24HourFormat ? "hh:mm" : "hh:mm A");

        painter.setPen(Qt::white);
        painter.drawText(rect(), Qt::AlignCenter, text.left(5));
        return;
    }

    const QString currentTimeString = current.toString(m_24HourFormat ? "hhmm" : "hhmma");
    // check cache valid
    if (m_cachedTime != currentTimeString)
    {
        m_cachedTime = currentTimeString;

        // draw new pixmap
        m_cachedTime = currentTimeString;
        m_cachedIcon = QPixmap(size());
        m_cachedIcon.fill(Qt::transparent);
        QPainter p(&m_cachedIcon);

        // draw fashion mode datetime plugin
        const int perfectIconSize = qMin(width(), height()) * 0.8;
        const QRect r = rect();

        // draw background
        const QPixmap background = loadSvg(":/icons/resources/icons/background.svg", QSize(perfectIconSize, perfectIconSize));
        const QPoint backgroundOffset = r.center() - background.rect().center();
        p.drawPixmap(backgroundOffset, background);

        const int bigNumHeight = perfectIconSize / 2.5;
        const int bigNumWidth = double(bigNumHeight) * 8 / 18;
        const int smallNumHeight = bigNumHeight / 2;
        const int smallNumWidth = double(smallNumHeight) * 5 / 9;

        // draw big num 1
        const QString bigNum1Path = QString(":/icons/resources/icons/big%1.svg").arg(currentTimeString[0]);
        const QPixmap bigNum1 = loadSvg(bigNum1Path, QSize(bigNumWidth, bigNumHeight));
        const QPoint bigNum1Offset = backgroundOffset + QPoint(perfectIconSize / 2 - bigNumWidth * 2 + 1, perfectIconSize / 2 - bigNumHeight / 2);
        p.drawPixmap(bigNum1Offset, bigNum1);

        // draw big num 2
        const QString bigNum2Path = QString(":/icons/resources/icons/big%1.svg").arg(currentTimeString[1]);
        const QPixmap bigNum2 = loadSvg(bigNum2Path, QSize(bigNumWidth, bigNumHeight));
        const QPoint bigNum2Offset = bigNum1Offset + QPoint(bigNumWidth + 1, 0);
        p.drawPixmap(bigNum2Offset, bigNum2);

        // draw small num 1
        const QString smallNum1Path = QString(":/icons/resources/icons/small%1.svg").arg(currentTimeString[2]);
        const QPixmap smallNum1 = loadSvg(smallNum1Path, QSize(smallNumWidth, smallNumHeight));
        const QPoint smallNum1Offset = bigNum2Offset + QPoint(bigNumWidth + 2, smallNumHeight);
        p.drawPixmap(smallNum1Offset, smallNum1);

        // draw small num 2
        const QString smallNum2Path = QString(":/icons/resources/icons/small%1.svg").arg(currentTimeString[3]);
        const QPixmap smallNum2 = loadSvg(smallNum2Path, QSize(smallNumWidth, smallNumHeight));
        const QPoint smallNum2Offset = smallNum1Offset + QPoint(smallNumWidth + 1, 0);
        p.drawPixmap(smallNum2Offset, smallNum2);
    }

    // draw cached fashion mode time item
    painter.drawPixmap(rect().center() - m_cachedIcon.rect().center(), m_cachedIcon);
}

void DatetimeWidget::mousePressEvent(QMouseEvent *e)
{
    if (e->button() != Qt::RightButton)
        return QWidget::mousePressEvent(e);

    const QPoint p(e->pos() - rect().center());
    if (p.manhattanLength() < std::min(width(), height()) * 0.8 * 0.5)
    {
        emit requestContextMenu();
        return;
    }

    QWidget::mousePressEvent(e);
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
