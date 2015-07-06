#include "arrowrectangle.h"

ArrowRectangle::ArrowRectangle(QWidget * parent) :
    QWidget(parent)
{
    this->setWindowFlags(Qt::FramelessWindowHint | Qt::ToolTip | Qt::WindowStaysOnTopHint);
    this->setAttribute(Qt::WA_TranslucentBackground);
}

void ArrowRectangle::showAtLeft(int x, int y)
{
    this->move(x,y);
    if (this->isHidden())
    {
        this->show();
    }

    this->arrowDirection = ArrowRectangle::arrowLeft;
    this->repaint();
}

void ArrowRectangle::showAtRight(int x, int y)
{
    this->move(x,y);
    if (this->isHidden())
    {
        this->show();
    }

    this->arrowDirection = ArrowRectangle::arrowRight;
    this->repaint();
}

void ArrowRectangle::showAtTop(int x, int y)
{
    this->move(x,y);
    if (this->isHidden())
    {
        this->show();
    }

    this->arrowDirection = ArrowRectangle::arrowTop;
    this->repaint();
}

void ArrowRectangle::showAtBottom(int x, int y)
{
    this->move(x,y);
    if (this->isHidden())
    {
        this->show();
    }

    this->arrowDirection = ArrowRectangle::arrowBottom;
    this->repaint();
}

void ArrowRectangle::setContent(QWidget *content)
{
    content->setParent(this);
    content->move((width() - content->width()) / 2,(height() - content->height()) / 2);
}

void ArrowRectangle::move(int x, int y)
{
    switch (arrowDirection)
    {
    case arrowLeft:
        QWidget::move(x,y - height() / 2);
        break;
    case arrowRight:
        QWidget::move(x - width(),y - height() / 2);
        break;
    case arrowTop:
        QWidget::move(x - width() / 2,y);
        break;
    case arrowBottom:
        QWidget::move(x - width() / 2,y - height());
        break;
    default:
        QWidget::move(x,y);
        break;
    }

}

// override methods
void ArrowRectangle::paintEvent(QPaintEvent *)
{
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing);

    QPainterPath border;
    QRectF textRec;

    switch (arrowDirection)
    {
    case ArrowRectangle::arrowLeft:
        border = getLeftCornerPath();
        textRec = QRectF(arrowHeight,0,width() - arrowHeight, height());
        break;
    case ArrowRectangle::arrowRight:
        border = getRightCornerPath();
        textRec = QRectF(0,0,width() - arrowHeight, height());
        break;
    case ArrowRectangle::arrowTop:
        border = getTopCornerPath();
        textRec = QRectF(0,arrowHeight,width(), height() - arrowHeight);
        break;
    case ArrowRectangle::arrowBottom:
        border = getBottomCornerPath();
        textRec = QRectF(0,0,width(), height() - arrowHeight);
        break;
    default:
        border = getRightCornerPath();
        textRec = QRectF(0,0,width() - arrowHeight, height());
    }

    QPen strokePen;
    strokePen.setColor(strokeColor);
    strokePen.setWidth(strokeWidth);
    painter.strokePath(border, strokePen);
    painter.fillPath(border, QBrush(backgroundColor == "" ? QColor(0,0,0,150) : QColor(backgroundColor)));
}

int ArrowRectangle::getRadius()
{
    return this->radius;
}

int ArrowRectangle::getArrowHeight()
{
    return this->arrowHeight;
}

int ArrowRectangle::getArrowWidth()
{
    return this->arrowWidth;
}

QString ArrowRectangle::getBackgroundColor()
{
    return this->backgroundColor;
}

void ArrowRectangle::setWidth(int value)
{
    this->setMinimumWidth(value);
    this->setMaximumWidth(value);
}

void ArrowRectangle::setHeight(int value)
{
    this->setMinimumHeight(value);
    this->setMaximumHeight(value);
}

void ArrowRectangle::setRadius(int value)
{
    this->radius = value;
}

void ArrowRectangle::setArrowHeight(int value)
{
    this->arrowHeight = value;
}

void ArrowRectangle::setArrowWidth(int value)
{
    this->arrowWidth = value;
}

void ArrowRectangle::setBackgroundColor(QString value)
{
    this->backgroundColor = value;
}

QPainterPath ArrowRectangle::getLeftCornerPath()
{
    QRect rect = this->rect().marginsRemoved(QMargins(shadowWidth,shadowWidth,shadowWidth,shadowWidth));

    QPoint cornerPoint(rect.x(), rect.y() + rect.height() / 2);
    QPoint topLeft(rect.x() + arrowHeight, rect.y());
    QPoint topRight(rect.x() + rect.width(), rect.y());
    QPoint bottomRight(rect.x() + rect.width(), rect.y() + rect.height());
    QPoint bottomLeft(rect.x() + arrowHeight, rect.y() + rect.height());
    int radius = this->radius > (rect.height() / 2) ? rect.height() / 2 : this->radius;

    QPainterPath border;
    border.moveTo(topLeft);
    border.lineTo(topRight.x() - radius, topRight.y());
    border.arcTo(topRight.x() - 2 * radius, topRight.y(), 2 * radius, 2 * radius, 90, -90);
    border.lineTo(bottomRight.x(), bottomRight.y() - radius);
    border.arcTo(bottomRight.x() - 2 * radius, bottomRight.y() - 2 * radius, 2 * radius, 2 * radius, 0, -90);
    border.lineTo(bottomLeft);
    border.lineTo(cornerPoint);
    border.lineTo(topLeft);

    return border;
}

QPainterPath ArrowRectangle::getRightCornerPath()
{
    QRect rect = this->rect().marginsRemoved(QMargins(shadowWidth,shadowWidth,shadowWidth,shadowWidth));

    QPoint cornerPoint(rect.x() + rect.width(), rect.y() + rect.height() / 2);
    QPoint topLeft(rect.x(), rect.y());
    QPoint topRight(rect.x() + rect.width() - arrowHeight, rect.y());
    QPoint bottomRight(rect.x() + rect.width() - arrowHeight, rect.y() + rect.height());
    QPoint bottomLeft(rect.x(), rect.y() + rect.height());
    int radius = this->radius > (rect.height() / 2) ? rect.height() / 2 : this->radius;

    QPainterPath border;
    border.moveTo(topLeft.x() + radius, topLeft.y());
    border.lineTo(topRight);
    border.lineTo(cornerPoint);
    border.lineTo(bottomRight);
    border.lineTo(bottomLeft.x() + radius, bottomLeft.y());
    border.arcTo(bottomLeft.x(), bottomLeft.y() - 2 * radius, 2 * radius, 2 * radius, -90, -90);
    border.lineTo(topLeft.x(), topLeft.y() + radius);
    border.arcTo(topLeft.x(), topLeft.y(), 2 * radius, 2 * radius, 180, -90);

    return border;
}

QPainterPath ArrowRectangle::getTopCornerPath()
{
    QRect rect = this->rect().marginsRemoved(QMargins(shadowWidth,shadowWidth,shadowWidth,shadowWidth));

    QPoint cornerPoint(rect.x() + rect.width() / 2, rect.y());
    QPoint topLeft(rect.x(), rect.y() + arrowHeight);
    QPoint topRight(rect.x() + rect.width(), rect.y() + arrowHeight);
    QPoint bottomRight(rect.x() + rect.width(), rect.y() + rect.height());
    QPoint bottomLeft(rect.x(), rect.y() + rect.height());
    int radius = this->radius > (rect.height() / 2 - arrowHeight) ? rect.height() / 2 -arrowHeight : this->radius;

    QPainterPath border;
    border.moveTo(topLeft.x() + radius, topLeft.y());
    border.lineTo(cornerPoint.x() - arrowWidth / 2, cornerPoint.y() + arrowHeight);
    border.lineTo(cornerPoint);
    border.lineTo(cornerPoint.x() + arrowWidth / 2, cornerPoint.y() + arrowHeight);
    border.lineTo(topRight.x() - radius, topRight.y());
    border.arcTo(topRight.x() - 2 * radius, topRight.y(), 2 * radius, 2 * radius, 90, -90);
    border.lineTo(bottomRight.x(), bottomRight.y() - radius);
    border.arcTo(bottomRight.x() - 2 * radius, bottomRight.y() - 2 * radius, 2 * radius, 2 * radius, 0, -90);
    border.lineTo(bottomLeft.x() + radius, bottomLeft.y());
    border.arcTo(bottomLeft.x(), bottomLeft.y() - 2 * radius, 2 * radius, 2 * radius, - 90, -90);
    border.lineTo(topLeft.x(), topLeft.y() + radius);
    border.arcTo(topLeft.x(), topLeft.y(), 2 * radius, 2 * radius, 180, -90);

    return border;
}

QPainterPath ArrowRectangle::getBottomCornerPath()
{
    QRect rect = this->rect().marginsRemoved(QMargins(shadowWidth,shadowWidth,shadowWidth,shadowWidth));

    QPoint cornerPoint(rect.x() + rect.width() / 2, rect.y()  + rect.height());
    QPoint topLeft(rect.x(), rect.y());
    QPoint topRight(rect.x() + rect.width(), rect.y());
    QPoint bottomRight(rect.x() + rect.width(), rect.y() + rect.height() - arrowHeight);
    QPoint bottomLeft(rect.x(), rect.y() + rect.height() - arrowHeight);
    int radius = this->radius > (rect.height() / 2 - arrowHeight) ? rect.height() / 2 -arrowHeight : this->radius;

    QPainterPath border;
    border.moveTo(topLeft.x() + radius, topLeft.y());
    border.lineTo(topRight.x() - radius, topRight.y());
    border.arcTo(topRight.x() - 2 * radius, topRight.y(), 2 * radius, 2 * radius, 90, -90);
    border.lineTo(bottomRight.x(), bottomRight.y() - radius);
    border.arcTo(bottomRight.x() - 2 * radius, bottomRight.y() - 2 * radius, 2 * radius, 2 * radius, 0, -90);
    border.lineTo(cornerPoint.x() + arrowWidth / 2, cornerPoint.y() - arrowHeight);
    border.lineTo(cornerPoint);
    border.lineTo(cornerPoint.x() - arrowWidth / 2, cornerPoint.y() - arrowHeight);
    border.lineTo(bottomLeft.x() + radius, bottomLeft.y());
    border.arcTo(bottomLeft.x(), bottomLeft.y() - 2 * radius, 2 * radius, 2 * radius, -90, -90);
    border.lineTo(topLeft.x(), topLeft.y() + radius);
    border.arcTo(topLeft.x(), topLeft.y(), 2 * radius, 2 * radius, 180, -90);

    return border;
}

ArrowRectangle::~ArrowRectangle()
{

}


