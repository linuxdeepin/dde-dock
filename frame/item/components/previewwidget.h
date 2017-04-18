#ifndef PREVIEWWIDGET_H
#define PREVIEWWIDGET_H

#include <QWidget>

class PreviewWidget : public QWidget
{
    Q_OBJECT
public:
    explicit PreviewWidget(const WId wid, QWidget *parent = 0);

signals:
    void requestActivateWindow(const WId wid) const;

private slots:
    void refershImage();

private:
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);

private:
    const WId m_wid;
    QImage m_image;
};

#endif // PREVIEWWIDGET_H
