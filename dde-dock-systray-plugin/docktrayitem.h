#ifndef DOCKTRAYITEM_H
#define DOCKTRAYITEM_H

#include <QWindow>

#include "abstractdockitem.h"

class DockTrayItem : public AbstractDockItem
{
    Q_OBJECT

    enum Style { Simple, Composite };
public:
    explicit DockTrayItem(QWidget *parent = 0);
    ~DockTrayItem();

    static DockTrayItem* fromWinId(WId winId, QWidget *parent=0);

    void setTitle(const QString &title) Q_DECL_OVERRIDE;
    void setIcon(const QString &iconPath, int size = 42) Q_DECL_OVERRIDE;
    void setMoveable(bool value) Q_DECL_OVERRIDE;
    bool moveable() Q_DECL_OVERRIDE;
    void setActived(bool value) Q_DECL_OVERRIDE;
    bool actived() Q_DECL_OVERRIDE;
    void setIndex(int value) Q_DECL_OVERRIDE;
    int index() Q_DECL_OVERRIDE;

    QWidget * getContents() Q_DECL_OVERRIDE;
};

#endif // DOCKTRAYITEM_H
