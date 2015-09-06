#ifndef COMPOSITETRAYITEM_H
#define COMPOSITETRAYITEM_H

#include <QFrame>
#include <QMap>

#include "interfaces/dockconstants.h"

class TrayIcon;
class CompositeTrayItem : public QFrame
{
    Q_OBJECT
public:
    explicit CompositeTrayItem(QWidget *parent = 0);
    virtual ~CompositeTrayItem();

    void addTrayIcon(QString key, TrayIcon * item);
    void remove(QString key);

    Dock::DockMode mode() const;
    void setMode(const Dock::DockMode &mode);

private:
    Dock::DockMode m_mode;
    QMap<QString, TrayIcon*> m_icons;
    QPixmap m_itemMask;

    void relayout();
};

#endif // COMPOSITETRAYITEM_H
