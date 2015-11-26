#ifndef TRAYICON_H
#define TRAYICON_H

#include <QWindow>
#include <QFrame>

class QTimer;
class QPaintEvent;
class QMouseEvent;
class TrayIcon : public QFrame
{
    Q_OBJECT
public:
    explicit TrayIcon(WId winId, QWidget *parent = 0);
    ~TrayIcon();

    void maskOn();
    void maskOff();

protected:
    void paintEvent(QPaintEvent *);
//    void mousePressEvent(QMouseEvent *);

private:
    WId m_windowId;
    WId m_containerWid;
    bool m_masked;
    QTimer * m_updateTimer;

    void wrapWindow();
    void updateWindow();
    QImage getImageNonComposite();
    void sendClick(uint8_t, int, int);
};

#endif // TRAYICON_H
