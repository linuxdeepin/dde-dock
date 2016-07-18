#ifndef FASHIONTRAYITEM_H
#define FASHIONTRAYITEM_H

#include <QWidget>

#include <traywidget.h>

class FashionTrayItem : public QWidget
{
    Q_OBJECT

public:
    explicit FashionTrayItem(QWidget *parent = 0);

    TrayWidget *activeTray();

    void setMouseEnable(const bool enable);

public slots:
    void setActiveTray(TrayWidget *tray);

private:
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);

private:
    bool m_enableMouseEvent;

    TrayWidget *m_activeTray;

    QPoint m_pressPoint;
};

#endif // FASHIONTRAYITEM_H
