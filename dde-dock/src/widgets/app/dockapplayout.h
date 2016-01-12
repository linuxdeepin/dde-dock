#ifndef DOCKAPPLAYOUT_H
#define DOCKAPPLAYOUT_H

#include "../movablelayout.h"
#include "../../controller/apps/dockappmanager.h"

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

    void onAppItemRemove(const QString &id);
    void onAppItemAdd(DockAppItem *item, bool delayShow);

private:
    DockAppManager *m_appManager;
};

#endif // DOCKAPPLAYOUT_H
