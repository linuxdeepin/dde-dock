#ifndef DMOVABLEDIALOG_H
#define DMOVABLEDIALOG_H


#include <QDialog>
#include <QPoint>
class QMouseEvent;
class QPushButton;
class QResizeEvent;

class DMovabelDialog:public QDialog
{
public:
    DMovabelDialog(QWidget *parent = 0);
    ~DMovabelDialog();

    QPushButton* getCloseButton();

public slots:
    void setMovableHeight(int height);
    void moveCenter();

protected:
    void mouseMoveEvent(QMouseEvent *event);
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);
    void resizeEvent(QResizeEvent* event);

private:
    QPoint m_dragPosition;
    int m_movableHeight = 30;
    QPushButton* m_closeButton;
};

#endif // DMOVABLEDIALOG_H
