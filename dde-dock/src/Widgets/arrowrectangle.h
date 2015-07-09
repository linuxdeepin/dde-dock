#ifndef ARROWRECTANGLE_H
#define ARROWRECTANGLE_H

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

    int getRadius();
    int getArrowHeight();
    int getArrowWidth();
    int getMargin();
    QString getBackgroundColor();

    void setArrorDirection(ArrowDirection value);
    void setWidth(int value);
    void setHeight(int value);
    void setRadius(int value);
    void setArrowHeight(int value);
    void setArrowWidth(int value);
    void setMargin(int value);
    void setBackgroundColor(QString value);

    void show(int x,int y);
    void showAtLeft(int x,int y);
    void showAtRight(int x,int y);
    void showAtTop(int x,int y);
    void showAtBottom(int x,int y);

    void delayHide(int interval = 500);
    void setContent(QWidget *content);
    void destroyContent();
    void move(int x,int y);

public slots:
    void slotHide();
    void slotCancelHide();
protected:
    virtual void paintEvent(QPaintEvent *);

private:
    int radius = 3;
    int arrowHeight = 8;
    int arrowWidth = 20;
    int m_margin = 5;
    QString backgroundColor;

    int strokeWidth = 1;
    QColor strokeColor = QColor(255,255,255,130);
    int shadowWidth = 2;
    QColor shadowColor = Qt::black;

    ArrowDirection arrowDirection = ArrowRectangle::ArrowRight;

    QWidget *m_content = NULL;
    QTimer *m_destroyTimer = NULL;
private:
    QPainterPath getLeftCornerPath();
    QPainterPath getRightCornerPath();
    QPainterPath getTopCornerPath();
    QPainterPath getBottomCornerPath();

};

#endif // ARROWRECTANGLE_H
