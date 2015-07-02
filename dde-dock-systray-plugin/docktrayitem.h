#ifndef DOCKTRAYITEM_H
#define DOCKTRAYITEM_H

#include <QWindow>
#include <QWidget>

class DockTrayItem : public QWidget
{
    Q_OBJECT

    enum Style { Simple, Composite };
public:
    ~DockTrayItem();

    static DockTrayItem* fromWinId(WId winId);

private:
    DockTrayItem(QWidget *parent = 0);
};

#endif // DOCKTRAYITEM_H
