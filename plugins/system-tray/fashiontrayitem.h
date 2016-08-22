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
    void resizeEvent(QResizeEvent *e);
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);

    const QPixmap loadSvg(const QString &fileName, const int size) const;

private:
    bool m_enableMouseEvent;

    TrayWidget *m_activeTray;

    QPixmap m_backgroundPixmap;
    QPoint m_pressPoint;
};

#endif // FASHIONTRAYITEM_H
