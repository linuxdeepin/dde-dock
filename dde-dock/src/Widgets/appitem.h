#ifndef APPITEM_H
#define APPITEM_H

#include <QObject>
#include <QWidget>
#include <QPushButton>
#include <QMouseEvent>
#include <QDrag>
#include <QRectF>
#include <QDrag>
#include <QMimeData>
#include <QPixmap>
#include <QImage>
#include <QDebug>
#include "abstractdockitem.h"
#include "Controller/dockmodedata.h"
#include "appicon.h"
#include "appbackground.h"

class AppItem : public AbstractDockItem
{
    Q_OBJECT
    Q_PROPERTY(QPoint pos READ pos WRITE move)
public:
    AppItem(QWidget *parent = 0);
    AppItem(QString title, QWidget *parent = 0);
    AppItem(QString title, QString iconPath, QWidget *parent = 0);
    ~AppItem();

    void setIcon(const QString &iconPath, int size = 42);
    void setActived(bool value);
    void setCurrentOpened(bool value);
    bool currentOpened();

protected:
    void mousePressEvent(QMouseEvent *);
    void mouseReleaseEvent(QMouseEvent *);
    void mouseDoubleClickEvent(QMouseEvent *);
    void mouseMoveEvent(QMouseEvent *);
    void enterEvent(QEvent * event);
    void leaveEvent(QEvent * event);
    void dragEnterEvent(QDragEnterEvent * event);
    void dragLeaveEvent(QDragLeaveEvent * event);
    void dropEvent(QDropEvent * event);

private slots:
    void slotDockModeChanged(DockConstants::DockMode newMode,DockConstants::DockMode oldMode);
    void reanchorIcon();
    void resizeBackground();

private:
    void resizeResources();
    void initBackground();

private:
    DockModeData *dockCons = DockModeData::getInstants();
    AppBackground * appBackground = NULL;

    QLabel * m_appIcon = NULL;
    bool m_isCurrentOpened = false;
    QString m_itemTitle = "";
    QString m_itemIconPath = "";
};

#endif // APPITEM_H
