#ifndef REFRESHBUTTON_H
#define REFRESHBUTTON_H

#include <QWidget>

class QTimer;

class RefreshButton : public QWidget
{
    Q_OBJECT
public:
    explicit RefreshButton(QWidget *parent = nullptr);
    void setRotateIcon(QString path);
    void startRotate();
    void stopRotate();

signals:
    void clicked();

protected:
    void paintEvent(QPaintEvent *e) override;
    void mousePressEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;

private:
    void initConnect();

    QTimer *m_refreshTimer;
    QPixmap m_pixmap;
    QPoint m_pressPos;
    int m_rotateAngle;
};

#endif // REFRESHBUTTON_H
