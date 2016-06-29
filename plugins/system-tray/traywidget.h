#ifndef TRAYWIDGET_H
#define TRAYWIDGET_H

#include <QWidget>
#include <QTimer>

class TrayWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TrayWidget(quint32 winId, QWidget *parent = 0);
    ~TrayWidget();

    const QImage trayImage() const;
    void sendClick(uint8_t mouseButton, int x, int y);

private:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);

    void wrapWindow();
    void updateIcon();
    void hideIcon();
    QImage getImageNonComposite() const;

private:
    WId m_windowId;
    WId m_containerWid;
    QImage m_image;

    QTimer *m_updateTimer;
    QPoint m_pressPoint;
};

#endif // TRAYWIDGET_H
