#ifndef WIREDITEM_H
#define WIREDITEM_H

#include <QWidget>

class WiredItem : public QWidget
{
    Q_OBJECT

public:
    explicit WiredItem(QWidget *parent = 0);

protected:
    QSize sizeHint() const;
};

#endif // WIREDITEM_H
