#ifndef ABSTRACTDOCKITEM_H
#define ABSTRACTDOCKITEM_H

#include "Widgets/arrowrectangle.h"
#include "Widgets/highlighteffect.h"
#include <QDebug>

class DBusMenu;
class DBusMenuManager;
class AbstractDockItem : public QFrame
{
    Q_OBJECT
public:
    explicit AbstractDockItem(QWidget *parent = 0);
    virtual ~AbstractDockItem();

    virtual QString getTitle();
    virtual QWidget * getApplet();

    virtual bool moveable();
    virtual bool actived();

    void resize(int width,int height);
    void resize(const QSize &size);

    QPoint getNextPos();
    void setNextPos(const QPoint &value);
    void setNextPos(int x, int y);
    void move(const QPoint &value);
    void move(int x, int y);

    int globalX();
    int globalY();
    QPoint globalPos();

    void showPreview();
    void hidePreview(int interval = 200);
    void cancelHide();
    void resizePreview();

    void showMenu();
    virtual QString getMenuContent();
    virtual void invokeMenuItem(QString itemId, bool checked);

    void setParent(QWidget * parent);

signals:
    void dragStart();
    void dragEntered(QDragEnterEvent * event);
    void dragExited(QDragLeaveEvent * event);
    void drop(QDropEvent * event);
    void mouseEntered();
    void mouseExited();
    void mousePress(int x, int y);
    void mouseRelease(int x, int y);
    void mouseDoubleClick();
    void widthChanged();

protected:

    bool m_moveable = true;
    bool m_isActived = false;
    ArrowRectangle *m_previewAR = new ArrowRectangle();
    HighlightEffect * m_highlight = NULL;

    QPoint m_itemNextPos;

    DBusMenu * m_dbusMenu = NULL;
    DBusMenuManager * m_dbusMenuManager = NULL;

    void initHighlight();
};

#endif // ABSTRACTDOCKITEM_H
