#ifndef HORIZONTALSEPERATOR_H
#define HORIZONTALSEPERATOR_H

#include <QWidget>

class HorizontalSeperator : public QWidget
{
    Q_OBJECT

public:
    explicit HorizontalSeperator(QWidget *parent = 0);

protected:
    void paintEvent(QPaintEvent *e);
};

#endif // HORIZONTALSEPERATOR_H
