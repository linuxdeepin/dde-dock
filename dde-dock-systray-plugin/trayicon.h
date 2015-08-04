#ifndef TRAYICON_H
#define TRAYICON_H

#include <QWindow>
#include <QWidget>

class TrayIcon : public QWidget
{
    Q_OBJECT
public:
    explicit TrayIcon(WId winId, QWidget *parent = 0);

    void maskOn();
    void maskOff();

private:
    QWindow * m_win;
    QPixmap m_itemMask;

    void initItemMask();
};

#endif // TRAYICON_H
