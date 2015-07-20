#ifndef APPPREVIEWS_H
#define APPPREVIEWS_H

#include <QWidget>
#include <QHBoxLayout>
#include <QLabel>
#include <QDebug>
#include "DBus/dbusclientmanager.h"
#include "windowpreview.h"
#include "closebutton.h"
#include "../dockconstants.h"

class AppPreviewFrame : public QWidget
{
    Q_OBJECT

public:
    explicit AppPreviewFrame(QWidget *preview, const QString &title,int xid, QWidget *parent=0);
    void addPreview(QWidget *p);
    void setTitle(const QString &t);

protected:
    void mousePressEvent(QMouseEvent *);
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);

signals:
    void close(int xid);
    void activate(int xid);

private:
    void addCloseButton();
    void showCloseButton();
    void hideCloseButton();

private:
    int xidValue;
    CloseButton *m_cb;
};

class AppPreviews : public QWidget
{
    Q_OBJECT
public:
    explicit AppPreviews(QWidget *parent = 0);

    void addItem(const QString &title,int xid);
    void setTitle(const QString &title);

protected:
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);

signals:
    void mouseEntered();
    void mouseExited();
    void sizeChanged();

public slots:
    void removePreview(int xid);
    void activatePreview(int xid);

private:
    DBusClientManager *m_clientManager = new DBusClientManager(this);
    QHBoxLayout *m_mainLayout = NULL;
    QList<int> m_xidList;
    bool isClosing = false;
};

#endif // APPPREVIEWS_H
