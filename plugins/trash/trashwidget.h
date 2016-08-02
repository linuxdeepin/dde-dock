#ifndef TRASHWIDGET_H
#define TRASHWIDGET_H

#include "popupcontrolwidget.h"

#include <QWidget>

class TrashWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TrashWidget(QWidget *parent = 0);

    QWidget *popupApplet();

protected:
    void paintEvent(QPaintEvent *e);

private:
    PopupControlWidget *m_popupApplet;
};

#endif // TRASHWIDGET_H
