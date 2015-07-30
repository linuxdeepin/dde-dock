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
    Q_PROPERTY(int rValue READ getRValue WRITE setRValue)
    Q_PROPERTY(double sValue  READ getSValue WRITE setSValue)
public:
    explicit TransformLabel(QWidget *parent=0) : QLabel(parent){}

    int getRValue(){return m_rValue;}
    double getSValue(){return m_sValue;}
    void setRValue(int value)
    {
        if (!pixmap())
            return;
        QTransform rt;
        rt.translate(width() / 2, height() / 2);
        rt.rotate(value);
        rt.translate(-width() / 2, -height() / 2);
        setPixmap(pixmap()->transformed(rt));
        m_rValue = value;
    }

    void setSValue(double value)
    {
        if (!pixmap())
            return;
        QTransform st(1, 0, 0, 1, width()/2, height()/2);
        st.scale(value, value);
        st.rotate(90);//TODO work around here
        setPixmap(pixmap()->transformed(st));
        m_sValue = value;
    }

private:
    int m_rValue = 0;
    double m_sValue = 0;
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
    void enterEvent(QEvent *);

signals:
    void itemDropped(QPoint pos);
    void itemEntered();
    void itemExited();
    void itemMissing();

private:
    const int ICON_SIZE = 48;
};

#endif // SCREENMASK_H
