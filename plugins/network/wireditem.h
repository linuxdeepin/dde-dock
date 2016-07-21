#ifndef WIREDITEM_H
#define WIREDITEM_H

#include "networkmanager.h"

#include <QWidget>

class WiredItem : public QWidget
{
    Q_OBJECT

public:
    explicit WiredItem(QWidget *parent = 0);

protected:
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    QSize sizeHint() const;

private:
    void reloadIcon();

private:
    NetworkManager *m_networkManager;

    QPixmap m_icon;
};

#endif // WIREDITEM_H
