#ifndef ARROWRECTANGLE_H
#define ARROWRECTANGLE_H

#include <QWidget>
#include <QLabel>
#include <QTextLine>
#include <QHBoxLayout>
#include <QVBoxLayout>
#include <QPainter>

class ArrowRectangle : public QWidget
{
    Q_OBJECT
public:
    explicit ArrowRectangle(QWidget * parent = 0);
    ~ArrowRectangle();

    int getRadius();
    int getArrowHeight();
    int getArrowWidth();
    QString getBackgroundColor();

    void setWidth(int value);
    void setHeight(int value);
    void setRadius(int value);
    void setArrowHeight(int value);
    void setArrowWidth(int value);
    void setBackgroundColor(QString value);

    void showAtLeft(int x,int y);
    void showAtRight(int x,int y);
    void showAtTop(int x,int y);
    void showAtBottom(int x,int y);

    void setContent(QWidget *content);
    void move(int x,int y);
protected:
    virtual void paintEvent(QPaintEvent *);

private:
    enum ArrowDirection {
        arrowLeft,
        arrowRight,
        arrowTop,
        arrowBottom
    };

    int radius = 3;
    int arrowHeight = 8;
    int arrowWidth = 20;
    QString backgroundColor;

    int strokeWidth = 1;
    QColor strokeColor = QColor(255,255,255,130);
    int shadowWidth = 2;
    QColor shadowColor = Qt::black;

    ArrowDirection arrowDirection = ArrowRectangle::arrowRight;

private:
    QPainterPath getLeftCornerPath();
    QPainterPath getRightCornerPath();
    QPainterPath getTopCornerPath();
    QPainterPath getBottomCornerPath();

};

#endif // ARROWRECTANGLE_H
