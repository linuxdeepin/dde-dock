#ifndef DOCKPLUGINLAYOUT_H
#define DOCKPLUGINLAYOUT_H

#include "../movablelayout.h"
#include "../../dbus/dbusdisplay.h"
#include "../../controller/plugins/dockpluginsmanager.h"

class DockPluginLayout : public MovableLayout
{
    Q_OBJECT
public:
    explicit DockPluginLayout(QWidget *parent = 0);

    QSize sizeHint() const;
    void initAllPlugins();

signals:
    void needPreviewHide(bool immediately);
    void needPreviewShow(QPoint pos);
    void needPreviewUpdate();

private:
    void initPluginManager();
    DisplayRect getScreenRect();

private:
    DockPluginsManager *m_pluginManager;
};

#endif // DOCKPLUGINLAYOUT_H
