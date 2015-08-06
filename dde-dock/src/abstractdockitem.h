#ifndef ABSTRACTDOCKITEM_H
#define ABSTRACTDOCKITEM_H

#include "Widgets/previewarrowrectangle.h"
#include "Widgets/highlighteffect.h"
#include "dockconstants.h"
#include <QDebug>

class ItemTitleLabel : public QLabel
{
public:
    explicit ItemTitleLabel(QWidget * parent = 0);

    void setTitle(QString title);
};

class DBusMenu;
class DBusMenuManager;
class AbstractDockItem : public QFrame
{
    Q_OBJECT
public:
    explicit AbstractDockItem(QWidget *parent = 0);
    virtual ~AbstractDockItem();

    virtual QString getItemId();
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
    virtual void moveWithAnimation(QPoint targetPos, int duration = 200){}

    int globalX();
    int globalY();
    QPoint globalPos();

    void showPreview();
    void hidePreview(int interval = 150);
    void cancelHide();
    void resizePreview();

    void showMenu();
    virtual QString getMenuContent();
    virtual void invokeMenuItem(QString menuItemId, bool checked);

    void setParent(QWidget * parent);

signals:
    void dragStart();
    void dragEntered(QDragEnterEvent * event);
    void dragExited(QDragLeaveEvent * event);
    void mouseEntered();
    void mouseExited();
    void mousePress(int x, int y);
    void mouseRelease(int x, int y);
    void widthChanged();
    void posChanged();
    void frameUpdate();
    void moveAnimationFinished();
    void previewHidden();

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
