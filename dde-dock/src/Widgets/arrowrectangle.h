#ifndef ARROWRECTANGLE_H
#define ARROWRECTANGLE_H

#include <QDesktopWidget>
#include <QWidget>
#include <QLabel>
#include <QTextLine>
#include <QHBoxLayout>
#include <QVBoxLayout>
#include <QPainter>
#include <QTimer>
#include <QDebug>

class ArrowRectangle : public QWidget
{
    Q_OBJECT
public:
    enum ArrowDirection {
        ArrowLeft,
        ArrowRight,
        ArrowTop,
        ArrowBottom
    };

    explicit ArrowRectangle(QWidget * parent = 0);
    ~ArrowRectangle();

    int getRadius() const;
    int getArrowHeight() const;
    int getArrowWidth() const;
    int getArrowX() const;
    int getArrowY() const;
    int getMargin() const;
    QString getBackgroundColor();

    void setArrorDirection(ArrowDirection value);
    void setWidth(int value);
    void setHeight(int value);
    void setRadius(int value);
    void setArrowHeight(int value);
    void setArrowWidth(int value);
    void setArrowX(int value);
    void setArrowY(int value);
    void setMargin(int value);
    void setBackgroundColor(QString value);

    virtual void show(ArrowDirection direction, int x,int y);

    void setContent(QWidget *content);
    void resizeWithContent();
    QSize getFixedSize();
    void move(int x,int y);

protected:
    void paintEvent(QPaintEvent *);

private:
    int radius = 3;
    int arrowHeight = 8;
    int arrowWidth = 12;
    int m_margin = 5;
    int m_arrowX = 0;
    int m_arrowY = 0;
    QString backgroundColor;

    int strokeWidth = 1;
    QColor strokeColor = QColor(255,255,255,130);
    int shadowWidth = 2;
    QColor shadowColor = Qt::black;

    ArrowDirection arrowDirection = ArrowRectangle::ArrowRight;

    QWidget *m_content = NULL;

    QPoint m_lastPos = QPoint(0, 0);
private:
    QPainterPath getLeftCornerPath();
    QPainterPath getRightCornerPath();
    QPainterPath getTopCornerPath();
    QPainterPath getBottomCornerPath();

};

#endif // ARROWRECTANGLE_H
