#ifndef WINDOWPREVIEW_H
#define WINDOWPREVIEW_H

#include <QWidget>
#include <QFrame>
#include <QImage>
#include <QStyle>
#include <QByteArray>

class Monitor;
class QPaintEvent;
class WindowPreview : public QFrame
{
    Q_OBJECT
    Q_PROPERTY(bool isHover READ isHover WRITE setIsHover)
public:
    friend class Monitor;

    WindowPreview(WId sourceWindow, QWidget *parent = 0);
    ~WindowPreview();

    QByteArray imageData;

    bool isHover() const;
    void setIsHover(bool isHover);

protected:
    void paintEvent(QPaintEvent * event);

private:
    WId m_sourceWindow;

    Monitor * m_monitor;
    int m_borderWidth = 3;
    bool m_isHover = false;

    void installMonitor();
    void removeMonitor();

    void prepareRepaint();
};

#endif // WINDOWPREVIEW_H
