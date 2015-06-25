#ifndef DOCKVIEW_H
#define DOCKVIEW_H

#include <QObject>
#include <QWidget>
#include <QAbstractItemView>

class DockView : public QAbstractItemView
{
    Q_OBJECT
public:
    explicit DockView(QWidget *parent = 0);

};

#endif // DOCKVIEW_H
