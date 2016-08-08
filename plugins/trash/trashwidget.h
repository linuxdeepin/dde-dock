#ifndef TRASHWIDGET_H
#define TRASHWIDGET_H

#include "popupcontrolwidget.h"

#include <QWidget>
#include <QPixmap>

class TrashWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TrashWidget(QWidget *parent = 0);

    QWidget *popupApplet();

    QSize sizeHint() const;

protected:
    void dragEnterEvent(QDragEnterEvent *e);
    void dropEvent(QDropEvent *e);
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);

private slots:
    void updateIcon();
    void moveToTrash(const QUrl &url);

private:
    PopupControlWidget *m_popupApplet;

    QPixmap m_icon;
};

#endif // TRASHWIDGET_H
