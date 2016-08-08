#ifndef CONTAINERWIDGET_H
#define CONTAINERWIDGET_H

#include <QWidget>

class ContainerWidget : public QWidget
{
    Q_OBJECT

public:
    explicit ContainerWidget(QWidget *parent = 0);

    QSize sizeHint() const;
};

#endif // CONTAINERWIDGET_H
