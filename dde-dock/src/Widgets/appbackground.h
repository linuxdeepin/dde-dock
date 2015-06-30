#ifndef APPBACKGROUND_H
#define APPBACKGROUND_H

#include <QObject>
#include <QLabel>
#include <QDebug>
#include "dockconstants.h"

class AppBackground : public QLabel
{
    Q_OBJECT
public:
    explicit AppBackground(QWidget *parent = 0);

signals:

public slots:

};

#endif // APPBACKGROUND_H
