/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
};

#endif // APPPREVIEWSCONTAINER_H
