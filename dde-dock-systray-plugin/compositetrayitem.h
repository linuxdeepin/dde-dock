#ifndef COMPOSITETRAYITEM_H
#define COMPOSITETRAYITEM_H

#include <QFrame>
#include <QMap>

#include "interfaces/dockconstants.h"

class TrayIcon;
class QLabel;
class CompositeTrayItem : public QFrame
{
    Q_OBJECT
public:
    explicit CompositeTrayItem(QWidget *parent = 0);
    virtual ~CompositeTrayItem();

    void addTrayIcon(QString key, TrayIcon * item);
    void remove(QString key);
    void setMode(const Dock::DockMode &mode);
    void clear();

    bool exist(const QString &key);
    QStringList trayIds() const;
    Dock::DockMode mode() const;

    void coverOn();
    void coverOff();

protected:
    void enterEvent(QEvent *) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent *) Q_DECL_OVERRIDE;

private:
    Dock::DockMode m_mode;
    QMap<QString, TrayIcon*> m_icons;
    QPixmap m_itemMask;
    QLabel * m_cover;
    bool m_isCovered;

    void relayout();
};

#endif // COMPOSITETRAYITEM_H
