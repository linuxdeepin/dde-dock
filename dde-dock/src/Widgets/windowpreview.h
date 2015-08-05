#ifndef WINDOWPREVIEW_H
#define WINDOWPREVIEW_H

#include <QWidget>
#include <QImage>
#include <QByteArray>

class QPaintEvent;
class WindowPreview : public QWidget
{
    Q_OBJECT

public:
    WindowPreview(WId sourceWindow, QWidget *parent = 0);
    ~WindowPreview();

    void onTimeout();

    QByteArray imageData;

protected:
    void paintEvent(QPaintEvent * event);

private:
    WId m_sourceWindow;
};

#endif // WINDOWPREVIEW_H
