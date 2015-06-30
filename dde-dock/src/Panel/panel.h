#ifndef PANEL_H
#define PANEL_H

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QDebug>
#include "Widgets/appitem.h"
#include "Widgets/docklayout.h"
#include "Widgets/screenmask.h"

class Panel : public QLabel
{
    Q_OBJECT
public:
    explicit Panel(QWidget *parent = 0);
    ~Panel();

    void resize(const QSize &size);
    void resize(int width,int height);

    void showScreenMask();
    void hideScreenMask();

signals:

public slots:
    void slotDragStarted();
    void slotItemDropped();
    void slotEnteredMask();
    void slotExitedMask();

private:
    DockLayout * leftLayout;
    DockLayout *rightLayout;
    QWidget * parentWidget = NULL;
    ScreenMask * maskWidget = NULL;
};

#endif // PANEL_H
