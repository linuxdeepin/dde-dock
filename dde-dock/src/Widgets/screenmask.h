#ifndef SCREENMASK_H
#define SCREENMASK_H

#include <QApplication>
#include <QDesktopWidget>
#include <QWidget>
#include <QLabel>
#include <QPixmap>
#include <QTransform>
#include <QPropertyAnimation>
#include <QDragEnterEvent>
#include <QDragMoveEvent>
#include <QDropEvent>
#include <QMimeData>
#include <QTimer>
#include <QDebug>
#include "DBus/dbusdockedappmanager.h"
#include "appitem.h"

class TransformLabel : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(QTransform transform READ getTransform WRITE setTransform)
public:
    explicit TransformLabel(QWidget *parent=0) : QLabel(parent){}

    QTransform getTransform(){return this->pixTransform;}
    void setTransform(const QTransform &value)
    {
        this->pixTransform = value;
        this->setPixmap(this->pixmap()->transformed(value));
    }

private:
    QTransform pixTransform;
};

class ScreenMask : public QWidget
{
    Q_OBJECT
public:
    explicit ScreenMask(QWidget *parent = 0);

protected:
    void dragEnterEvent(QDragEnterEvent *event);
    void dragLeaveEvent(QDragLeaveEvent *event);
    void dropEvent(QDropEvent *event);

signals:
    void itemDropped(QPoint pos);
    void itemEntered();
    void itemExited();

public slots:
};

#endif // SCREENMASK_H
