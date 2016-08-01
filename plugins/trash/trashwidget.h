#ifndef TRASHWIDGET_H
#define TRASHWIDGET_H

#include <QWidget>

class TrashWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TrashWidget(QWidget *parent = 0);

protected:
    void paintEvent(QPaintEvent *e);
};

#endif // TRASHWIDGET_H
