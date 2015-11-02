#ifndef APPPREVIEWSCONTAINER_H
#define APPPREVIEWSCONTAINER_H

#include <QWidget>
#include "dbus/dbusclientmanager.h"

class QHBoxLayout;
class AppPreviewLoaderFrame;

class AppPreviewsContainer : public QWidget
{
    Q_OBJECT
public:
    explicit AppPreviewsContainer(QWidget *parent = 0);
    ~AppPreviewsContainer();

    void addItem(const QString &title,int xid);

protected:
    void leaveEvent(QEvent *);

signals:
    void requestHide();
    void sizeChanged();

public slots:
    void removePreview(int xid);
    void activatePreview(int xid);
    void clearUpPreview();
    QSize getNormalContentSize();

private:
    void setItemCount(int count);

private:
    QMap<int, AppPreviewLoaderFrame *> m_previewMap;
    DBusClientManager *m_clientManager;
    QHBoxLayout *m_mainLayout;
    bool m_isClosing;
};

#endif // APPPREVIEWSCONTAINER_H
