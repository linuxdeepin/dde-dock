#ifndef CLOSEBUTTON_H
#define CLOSEBUTTON_H

#include <QWidget>
#include <QLabel>

class CloseButton : public QLabel
{
    Q_OBJECT
public:
    explicit CloseButton(QWidget *parent = 0);

    void mousePressEvent(QMouseEvent *ev);
    void mouseReleaseEvent(QMouseEvent *ev);
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);
signals:
    void clicked();
    void hovered();
    void exited();
    void pressed();
    void released();

public slots:
private:
    void setIcon(const QString &path);

private:
    bool isPressed = false;
    const QString ICON_NORMAL_PATH = "://Resources/images/close_normal.png";
    const QString ICON_HOVER_PATH = "://Resources/images/close_hover.png";
    const QString ICON_PRESS_PATH = "://Resources/images/close_press.png";
};

#endif // CLOSEBUTTON_H
