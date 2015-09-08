#ifndef ABSTRACTDOCKITEM_H
#define ABSTRACTDOCKITEM_H

#include <QDebug>
#include <QFrame>
#include <QLabel>
#include <QtDBus>
#include <QWidget>
#include <QJsonValue>
#include <QJsonObject>

#include "dbus/dbusmenu.h"
#include "highlighteffect.h"
#include "dbus/dbusmenumanager.h"
#include "previewarrowrectangle.h"
#include "widgets/arrowrectangle.h"
#include "interfaces/dockconstants.h"

class ItemTitleLabel : public QLabel
{
public:
    explicit ItemTitleLabel(QWidget * parent = 0);

    void setTitle(QString title);
};

class AbstractDockItem : public QFrame
{
    Q_OBJECT
public:
    explicit AbstractDockItem(QWidget *parent = 0);
    virtual ~AbstractDockItem();

    virtual QString getTitle();
    virtual QString getItemId();
    virtual QString getMenuContent();
    virtual QWidget * getApplet();
    virtual bool moveable();
    virtual bool actived();
    virtual void invokeMenuItem(QString menuItemId, bool checked);
    virtual void moveWithAnimation(QPoint targetPos, int duration = 100){Q_UNUSED(targetPos) Q_UNUSED(duration)}

    void setNextPos(int x, int y);
    void setNextPos(const QPoint &value);
    void move(int x, int y);
    void move(const QPoint &value);
    void resize(const QSize &size);
    void resize(int width,int height);
    void showMenu();
    void cancelHide();
    void showPreview();
    void resizePreview();
    void hidePreview(int interval = 150);
    void setParent(QWidget * parent);

    int globalX();
    int globalY();
    QPoint globalPos();
    QPoint getNextPos();

signals:
    void dragStart();
    void dragEntered(QDragEnterEvent * event);
    void dragExited(QDragLeaveEvent * event);
    void mouseEntered();
    void mouseExited();
    void mousePress(QMouseEvent *event);
    void mouseRelease(QMouseEvent *event);
    void widthChanged();
    void posChanged();
    void frameUpdate();
    void previewHidden();
    void moveAnimationFinished();

protected:
    bool m_moveable = true;
    bool m_isActived = false;
    PreviewArrowRectangle *m_previewAR = NULL;
    HighlightEffect * m_highlight = NULL;
    ItemTitleLabel *m_titleLabel = NULL;

    QPoint m_itemNextPos;
    QPoint m_previewPos;

    DBusMenu * m_dbusMenu = NULL;
    DBusMenuManager * m_dbusMenuManager = NULL;

    void initHighlight();
    void initTitleLabel();

private:
    const int TITLE_HEIGHT = 20;
};

#endif // ABSTRACTDOCKITEM_H
