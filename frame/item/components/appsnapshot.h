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

#ifndef APPSNAPSHOT_H
#define APPSNAPSHOT_H

#include <QWidget>
#include <QDebug>
#include <QTimer>
#include <QLabel>

#include <dimagebutton.h>
#include <DWindowManagerHelper>

DWIDGET_USE_NAMESPACE

#define SNAP_WIDTH       200
#define SNAP_HEIGHT      130

class AppSnapshot : public QWidget
{
    Q_OBJECT

public:
    explicit AppSnapshot(const WId wid, QWidget *parent = 0);

    WId wid() const { return m_wid; }
    const QImage snapshot() const { return m_snapshot; }
    const QString title() const { return m_title->text(); }

signals:
    void entered(const WId wid) const;
    void clicked(const WId wid) const;
    void requestCheckWindow() const;

public slots:
    void fetchSnapshot();
    void closeWindow() const;
    void compositeChanged() const;
    void setWindowTitle(const QString &title);

private:
    void dragEnterEvent(QDragEnterEvent *e);
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);

private:
    const WId m_wid;

    QImage m_snapshot;
    QLabel *m_title;
    DImageButton *m_closeBtn;

    DWindowManagerHelper *m_wmHelper;
};

#endif // APPSNAPSHOT_H
