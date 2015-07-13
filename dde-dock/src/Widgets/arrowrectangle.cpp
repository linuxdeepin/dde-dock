#include "arrowrectangle.h"

ArrowRectangle::ArrowRectangle(QWidget * parent) :
    QWidget(parent)
{
    this->setWindowFlags(Qt::FramelessWindowHint | Qt::ToolTip | Qt::WindowStaysOnTopHint);
    this->setAttribute(Qt::WA_TranslucentBackground);
}

void ArrowRectangle::show(int x, int y)
{
    if (m_destroyTimer)
        m_destroyTimer->stop();
    this->move(x,y);
    if (this->isHidden())
    {
        QWidget::show();
    }

    this->repaint();
}

void ArrowRectangle::showAtLeft(int x, int y)
{
    if (m_destroyTimer)
        m_destroyTimer->stop();
    this->arrowDirection = ArrowRectangle::ArrowLeft;
    this->move(x,y);
    if (this->isHidden())
    {
        QWidget::show();
    }

    this->repaint();
}

void ArrowRectangle::showAtRight(int x, int y)
{
    if (m_destroyTimer)
        m_destroyTimer->stop();
    this->arrowDirection = ArrowRectangle::ArrowRight;
    this->move(x,y);
    if (this->isHidden())
    {
        QWidget::show();
    }

    this->repaint();
}

void ArrowRectangle::showAtTop(int x, int y)
{
    if (m_destroyTimer)
        m_destroyTimer->stop();
    this->arrowDirection = ArrowRectangle::ArrowTop;
    this->move(x,y);
    if (this->isHidden())
    {
        QWidget::show();
    }

    this->repaint();
}

void ArrowRectangle::showAtBottom(int x, int y)
{
    if (m_destroyTimer)
        m_destroyTimer->stop();
    this->arrowDirection = ArrowRectangle::ArrowBottom;
    this->move(x,y);
    if (this->isHidden())
    {
        QWidget::show();
    }

    this->repaint();
}

void ArrowRectangle::delayHide(int interval)
{
    if (!m_destroyTimer)
    {
        m_destroyTimer = new QTimer(this);
        connect(m_destroyTimer,&QTimer::timeout,this,&ArrowRectangle::slotHide);
        connect(m_destroyTimer,&QTimer::timeout,m_destroyTimer,&QTimer::stop);
    }
    m_destroyTimer->stop();
    m_destroyTimer->start(interval);
}

void ArrowRectangle::cancelHide()
{
    if (m_destroyTimer)
        m_destroyTimer->stop();
}

void ArrowRectangle::setContent(QWidget *content)
{
    if (!content)
    {
        return;
    }
    if (m_content)
    {
        content->deleteLater();
        return;
    }
    m_content = content;
    m_content->setParent(this);

    resizeWithContent();
    switch(arrowDirection)
    {
    case ArrowLeft:
        m_content->move(arrowHeight + m_margin,m_margin);
        break;
    case ArrowRight:
        m_content->move(m_margin,m_margin);
        break;
    case ArrowTop:
        m_content->move(m_margin,m_margin + arrowHeight);
        break;
    case ArrowBottom:
        m_content->move(m_margin,m_margin);
        break;
    }
}

void ArrowRectangle::resizeWithContent()
{
    if (m_content)
    {
        switch(arrowDirection)
        {
        case ArrowLeft:
        case ArrowRight:
            resize(m_content->width() + m_margin * 2 + arrowHeight,m_content->height() + m_margin * 2);
            break;
        case ArrowTop:
        case ArrowBottom:
            resize(m_content->width() + m_margin * 2,m_content->height() + m_margin * 2 + arrowHeight);
            break;
        }
    }

    repaint();
}

void ArrowRectangle::destroyContent()
{
    if (m_content)
    {
        delete m_content;
        m_content = NULL;
    }
}

void ArrowRectangle::move(int x, int y)
{
    QDesktopWidget dw;
    QRect rec = dw.screenGeometry();
    int xLeftValue = x - width() / 2;
    int xRightValue = x + width() / 2 - rec.width();
    int yTopValue = y - height() / 2;
    int yBottomValue = y + height() / 2 - rec.height();
    switch (arrowDirection)
    {
    case ArrowLeft:
        if (yTopValue < rec.y())
        {
            setArrowY(height() / 2 + yTopValue);
            yTopValue = rec.y();
        }
        else if (yBottomValue > 0)
        {
            setArrowY(height() / 2 + yBottomValue);
            yTopValue = rec.height() - height();
        }
        QWidget::move(x,yTopValue);
        break;
    case ArrowRight:
        if (yTopValue < rec.y())
        {
            setArrowY(height() / 2 + yTopValue);
            yTopValue = rec.y();
        }
        else if (yBottomValue > 0)
        {
            setArrowY(height() / 2 + yBottomValue);
            yTopValue = rec.height() - height();
        }
        QWidget::move(x - width(),yTopValue);
        break;
    case ArrowTop:
        if (xLeftValue < rec.x())//out of screen in left side
        {
            setArrowX(width() / 2 + xLeftValue);
            xLeftValue = rec.x();
        }
        else if(xRightValue > 0)//out of screen in right side
        {
            setArrowX(width() / 2 + xRightValue);
            xLeftValue = rec.width() - width();
        }
        QWidget::move(xLeftValue,y);
        break;
    case ArrowBottom:
        if (xLeftValue < rec.x())//out of screen in left side
        {
            setArrowX(width() / 2 + xLeftValue);
            xLeftValue = rec.x();
        }
        else if(xRightValue > 0)//out of screen in right side
        {
            setArrowX(width() / 2 + xRightValue);
            xLeftValue = rec.width() - width();
        }
        QWidget::move(xLeftValue,y - height());
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
    case ArrowRectangle::ArrowLeft:
        border = getLeftCornerPath();
        textRec = QRectF(arrowHeight,0,width() - arrowHeight, height());
        break;
    case ArrowRectangle::ArrowRight:
        border = getRightCornerPath();
        textRec = QRectF(0,0,width() - arrowHeight, height());
        break;
    case ArrowRectangle::ArrowTop:
        border = getTopCornerPath();
        textRec = QRectF(0,arrowHeight,width(), height() - arrowHeight);
        break;
    case ArrowRectangle::ArrowBottom:
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
    painter.fillPath(border, QBrush(backgroundColor == "" ? QColor(0,0,0,200) : QColor(backgroundColor)));
}

void ArrowRectangle::slotHide()
{
    destroyContent();
    hide();
}

void ArrowRectangle::slotCancelHide()
{
    if (m_destroyTimer)
        m_destroyTimer->stop();
}

int ArrowRectangle::getRadius() const
{
    return this->radius;
}

int ArrowRectangle::getArrowHeight() const
{
    return this->arrowHeight;
}

int ArrowRectangle::getArrowWidth() const
{
    return this->arrowWidth;
}

int ArrowRectangle::getArrowX() const
{
    return this->m_arrowX;
}

int ArrowRectangle::getArrowY() const
{
    return this->m_arrowY;
}

int ArrowRectangle::getMargin() const
{
    return this->m_margin;
}

QString ArrowRectangle::getBackgroundColor()
{
    return this->backgroundColor;
}

void ArrowRectangle::setArrorDirection(ArrowDirection value)
{
    arrowDirection = value;
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

void ArrowRectangle::setArrowX(int value)
{
    if (value < arrowWidth / 2)
        this->m_arrowX = arrowWidth / 2;
    else if (value > (width() - arrowWidth / 2))
        this->m_arrowX = width() - arrowWidth / 2;
    else
        this->m_arrowX = value;
}

void ArrowRectangle::setArrowY(int value)
{
    if (value < arrowWidth / 2)
        this->m_arrowY = arrowWidth / 2;
    else if (value > (height() - arrowWidth / 2))
        this->m_arrowY = height() - arrowWidth / 2;
    else
        this->m_arrowY = value;

}

void ArrowRectangle::setMargin(int value)
{
    this->m_margin = value;
}

void ArrowRectangle::setBackgroundColor(QString value)
{
    this->backgroundColor = value;
}

QPainterPath ArrowRectangle::getLeftCornerPath()
{
    QRect rect = this->rect().marginsRemoved(QMargins(shadowWidth,shadowWidth,shadowWidth,shadowWidth));

    QPoint cornerPoint(rect.x(), rect.y() + (m_arrowY > 0 ? m_arrowY : rect.height() / 2));
    QPoint topLeft(rect.x() + arrowHeight, rect.y());
    QPoint topRight(rect.x() + rect.width(), rect.y());
    QPoint bottomRight(rect.x() + rect.width(), rect.y() + rect.height());
    QPoint bottomLeft(rect.x() + arrowHeight, rect.y() + rect.height());
    int radius = this->radius > (rect.height() / 2) ? rect.height() / 2 : this->radius;

    QPainterPath border;
    border.moveTo(topLeft.x() - radius,topLeft.y());
    border.lineTo(topRight.x() - radius, topRight.y());
    border.arcTo(topRight.x() - 2 * radius, topRight.y(), 2 * radius, 2 * radius, 90, -90);
    border.lineTo(bottomRight.x(), bottomRight.y() - radius);
    border.arcTo(bottomRight.x() - 2 * radius, bottomRight.y() - 2 * radius, 2 * radius, 2 * radius, 0, -90);
    border.lineTo(bottomLeft.x() - radius,bottomLeft.y());
    border.arcTo(bottomLeft.x(),bottomLeft.y() - 2 * radius,2 * radius,2 * radius,-90,-90);
    border.lineTo(cornerPoint.x() + arrowHeight,cornerPoint.y() + arrowWidth / 2);
    border.lineTo(cornerPoint);
    border.lineTo(cornerPoint.x() + arrowHeight,cornerPoint.y() - arrowWidth / 2);
    border.lineTo(topLeft.x(),topLeft.y() + radius);
    border.arcTo(topLeft.x(),topLeft.y(),2 * radius,2 * radius,-180,-90);
    border.lineTo(topLeft.x() - radius,topLeft.y());

    return border;
}

QPainterPath ArrowRectangle::getRightCornerPath()
{
    QRect rect = this->rect().marginsRemoved(QMargins(shadowWidth,shadowWidth,shadowWidth,shadowWidth));

    QPoint cornerPoint(rect.x() + rect.width(), rect.y() + (m_arrowY > 0 ? m_arrowY : rect.height() / 2));
    QPoint topLeft(rect.x(), rect.y());
    QPoint topRight(rect.x() + rect.width() - arrowHeight, rect.y());
    QPoint bottomRight(rect.x() + rect.width() - arrowHeight, rect.y() + rect.height());
    QPoint bottomLeft(rect.x(), rect.y() + rect.height());
    int radius = this->radius > (rect.height() / 2) ? rect.height() / 2 : this->radius;

    QPainterPath border;
    border.moveTo(topLeft.x() + radius, topLeft.y());
    border.lineTo(topRight.x() - radius,topRight.y());
    border.arcTo(topRight.x() - 2 * radius,topRight.y(),2 * radius,2 * radius,90,-90);
    border.lineTo(cornerPoint.x() - arrowHeight,cornerPoint.y() - arrowWidth / 2);
    border.lineTo(cornerPoint);
    border.lineTo(cornerPoint.x() - arrowHeight,cornerPoint.y() + arrowWidth / 2);
    border.lineTo(bottomRight.x(),bottomRight.y() - radius);
    border.arcTo(bottomRight.x() - 2 * radius,bottomRight.y() - 2 * radius,2 * radius,2 * radius,0,-90);
    border.lineTo(bottomLeft.x() + radius, bottomLeft.y());
    border.arcTo(bottomLeft.x(), bottomLeft.y() - 2 * radius, 2 * radius, 2 * radius, -90, -90);
    border.lineTo(topLeft.x(), topLeft.y() + radius);
    border.arcTo(topLeft.x(), topLeft.y(), 2 * radius, 2 * radius, 180, -90);

    return border;
}

QPainterPath ArrowRectangle::getTopCornerPath()
{
    QRect rect = this->rect().marginsRemoved(QMargins(shadowWidth,shadowWidth,shadowWidth,shadowWidth));

    QPoint cornerPoint(rect.x() + (m_arrowX > 0 ? m_arrowX : rect.width() / 2), rect.y());
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

    QPoint cornerPoint(rect.x() + (m_arrowX > 0 ? m_arrowX : rect.width() / 2), rect.y()  + rect.height());
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


