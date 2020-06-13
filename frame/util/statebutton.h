#ifndef STATEBUTTON_H
#define STATEBUTTON_H

#include <QWidget>

class StateButton : public QWidget
{
    Q_OBJECT

public:
    enum Type {
        Check,
        Fork
    };

public:
    explicit StateButton(QWidget *parent = nullptr);
    void setType(Type type);

signals:
    void click();

protected:
    void paintEvent(QPaintEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;

private:
    void drawCheck(QPainter &painter, QPen &pen, int radius);
    void drawFork(QPainter &painter, QPen &pen, int radius);

private:
    Type m_type;
};

#endif // STATEBUTTON_H
