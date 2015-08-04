#ifndef COMPOSITETRAYITEM_H
#define COMPOSITETRAYITEM_H

#include <QFrame>
#include <QMap>

#include <dock/dockconstants.h>

class CompositeTrayItem : public QFrame
{
    Q_OBJECT
public:
    explicit CompositeTrayItem(QWidget *parent = 0);
    virtual ~CompositeTrayItem();

    void addItem(QString key, QWidget * widget);
    void removeItem(QString key);

    Dock::DockMode mode() const;
    void setMode(const Dock::DockMode &mode);

private:
    Dock::DockMode m_mode;
    QMap<QString, QWidget*> m_items;

    void relayout();
};

#endif // COMPOSITETRAYITEM_H
