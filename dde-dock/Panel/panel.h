#ifndef PANEL_H
#define PANEL_H

#include <QWidget>
#include <QPushButton>
#include <QDebug>
#include "Widgets/appitem.h"
#include "Widgets/docklayout.h"

class Panel : public QWidget
{
    Q_OBJECT
public:
    explicit Panel(QWidget *parent = 0);
    ~Panel();

    void resize(const QSize &size);
    void resize(int width,int height);

signals:

public slots:
private:
    DockLayout * leftLayout;
};

#endif // PANEL_H
