#ifndef WINDOWPREVIEW_H
#define WINDOWPREVIEW_H

#include <QWidget>
#include <QWindow>

class QImage;
class QTimer;

class WindowPreview : public QWidget
{
    Q_OBJECT

public:
    WindowPreview(WId sourceWindow, QWidget *parent = 0);
    ~WindowPreview();

protected:
    void paintEvent(QPaintEvent * event) Q_DECL_OVERRIDE;

private:
    WId m_sourceWindow;
    QImage *m_cache;
    QTimer *m_timer;

    void clearCache();
    void updateCache();
};

#endif // WINDOWPREVIEW_H
