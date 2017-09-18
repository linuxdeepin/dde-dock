/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef PREVIEWWIDGET_H
#define PREVIEWWIDGET_H

#include <QWidget>
#include <QDebug>
#include <QDragEnterEvent>
#include <QTimer>
#include <QVBoxLayout>

#include <dimagebutton.h>
#include <dwindowmanagerhelper.h>

DWIDGET_USE_NAMESPACE

class PreviewWidget : public QWidget
{
    Q_OBJECT
public:
    explicit PreviewWidget(const WId wid, QWidget *parent = 0);

    void setTitle(const QString &title);

signals:
    void requestActivateWindow(const WId wid) const;
    void requestPreviewWindow(const WId wid) const;
    void requestCancelPreview() const;
    void requestHidePreview() const;

private slots:
    void refreshImage();
    void closeWindow();
    void showPreview();

    void updatePreviewSize();

private:
    void paintEvent(QPaintEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void dragLeaveEvent(QDragLeaveEvent *e);
    void dropEvent(QDropEvent *e);

private:
    const WId m_wid;
    QImage m_image;
    QString m_title;

    DImageButton *m_closeButton;
    QVBoxLayout *m_centralLayout;

    QTimer *m_droppedDelay;
    QTimer *m_mouseEnterTimer;

    bool m_hovered;

    DWindowManagerHelper *m_wmHelper;
};

#endif // PREVIEWWIDGET_H
