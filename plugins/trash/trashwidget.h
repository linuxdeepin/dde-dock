#ifndef TRASHWIDGET_H
#define TRASHWIDGET_H

#include "popupcontrolwidget.h"

#include <QWidget>
#include <QPixmap>

#include <DMenu>
#include <DAction>

class TrashWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TrashWidget(QWidget *parent = 0);

    QWidget *popupApplet();

    QSize sizeHint() const;

signals:
    void requestRefershWindowVisible() const;

protected:
    void dragEnterEvent(QDragEnterEvent *e);
    void dropEvent(QDropEvent *e);
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);

private:
    const QPoint popupMarkPoint();

private slots:
    void updateIcon();
    void showMenu();
    void removeApp(const QString &appKey);
    void menuTriggered(Dtk::Widget::DAction *action);
    void moveToTrash(const QUrl &url);

private:
    PopupControlWidget *m_popupApplet;

    QPixmap m_icon;

    Dtk::Widget::DAction m_openAct;
    Dtk::Widget::DAction m_clearAct;
};

#endif // TRASHWIDGET_H
