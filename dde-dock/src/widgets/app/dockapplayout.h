#ifndef DOCKAPPLAYOUT_H
#define DOCKAPPLAYOUT_H

#include "../movablelayout.h"
#include "../../controller/apps/dockappmanager.h"
#include "../../dbus/dbusdockedappmanager.h"

class DropMask;

class DockAppLayout : public MovableLayout
{
    Q_OBJECT
public:
    explicit DockAppLayout(QWidget *parent = 0);

    QSize sizeHint() const;
    void initEntries() const;
signals:
    void needPreviewUpdate();
    void needPreviewHide(bool immediately);
    void needPreviewShow(DockItem *item, QPoint pos);
    void itemHoverableChange(bool v);

protected:
    bool eventFilter(QObject *obj, QEvent *e);

private:
    void initDropMask();
    void initAppManager();

    bool isDraging() const;
    void setIsDraging(bool isDraging);

    void onDrop(QDropEvent *event);
    void onAppItemRemove(const QString &id);
    void onAppItemAdd(DockAppItem *item);
    void onAppAppend(DockAppItem *item);

    QStringList appIds();

private:
    bool m_isDraging;
    DropMask *m_mask;
    DockAppManager *m_appManager;
    DBusDockedAppManager *m_ddam;
};

#endif // DOCKAPPLAYOUT_H
