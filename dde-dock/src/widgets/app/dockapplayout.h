#ifndef DOCKAPPLAYOUT_H
#define DOCKAPPLAYOUT_H

#include "../movablelayout.h"
#include "../../controller/apps/dockappmanager.h"
#include "../../dbus/dbusdockedappmanager.h"

class DockAppLayout : public MovableLayout
{
    Q_OBJECT
public:
    explicit DockAppLayout(QWidget *parent = 0);

    QSize sizeHint() const;
    void initEntries();

signals:
    void needPreviewHide(bool immediately);
    void needPreviewShow(DockItem *item, QPoint pos);
    void needPreviewUpdate();


private:
    void initAppManager();

    void onDrop(QDropEvent *event);
    void onAppItemRemove(const QString &id);
    void onAppItemAdd(DockAppItem *item);
    void onAppAppend(DockAppItem *item);
    QStringList appIds();

private:
    DockAppManager *m_appManager;
    DBusDockedAppManager *m_ddam;
};

#endif // DOCKAPPLAYOUT_H
