#ifndef DOCKPLUGINLAYOUT_H
#define DOCKPLUGINLAYOUT_H

#include "../movablelayout.h"

class DockPluginLayout : public MovableLayout
{
    Q_OBJECT
public:
    explicit DockPluginLayout(QWidget *parent = 0);

    QSize sizeHint() const;
};

#endif // DOCKPLUGINLAYOUT_H
