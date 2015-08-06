#ifndef WINDOWPREVIEW_H
#define WINDOWPREVIEW_H

#include <QWidget>
#include <QImage>
#include <QByteArray>

class Monitor;
class QPaintEvent;
class WindowPreview : public QWidget
{
    Q_OBJECT
public:
    friend class Monitor;

    WindowPreview(WId sourceWindow, QWidget *parent = 0);
    ~WindowPreview();

    QByteArray imageData;

protected:
    void paintEvent(QPaintEvent * event);

private:
    WId m_sourceWindow;

    Monitor * m_monitor;

    void installMonitor();
    void removeMonitor();

    void prepareRepaint();
};

#endif // WINDOWPREVIEW_H
