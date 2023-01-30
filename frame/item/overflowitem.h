// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef OVERFLOWITEM_H
#define OVERFLOWITEM_H

#include "dockitem.h"

#include <DWidget>

QT_USE_NAMESPACE
class QScrollArea;
class QPushButton;
class QBoxLayout;

DWIDGET_BEGIN_NAMESPACE
class DIconButton;
DWIDGET_END_NAMESPACE

DWIDGET_USE_NAMESPACE

class OverflowItem : public DockItem
{
    Q_OBJECT

public:
    explicit OverflowItem(QWidget *parent = nullptr);
    inline ItemType itemType() const override { return OverflowIcon; }
    void setPopUpSize(int width, int height);
    void addItem(QWidget *item);
    void hidePopUpWindow();
    void setLayoutPosition(Dock::Position position);

protected:
    void enterEvent(QEvent *e) override;
    void mousePressEvent(QMouseEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;
    void leaveEvent(QEvent *e) override;
    bool eventFilter(QObject *watched, QEvent *e) override;

private:
    void paintEvent(QPaintEvent *e) override;
    void showPopupWindow(QWidget *const content, const bool model = false, const int radius = 6) override;

private:
    QPoint OverflowIconPosition(const QPixmap &pixmap) const;
    void initUI();
    void initLBtnSlot();
    void initRBtnSlot();
    void setbtnsVisible();
    void setbtnsShape();

private:
// status
    int m_width;
    bool m_clicked;
    bool m_showpopup;
// widgets
    QScrollArea *m_scrollarea;
    QWidget *m_centerScroll;
    QBoxLayout *m_popuplayout;
    DockPopupWindow *m_popupwindow;
    QBoxLayout *m_popupbtnslayout;
    DIconButton *m_left;
    DIconButton *m_right;
};

#endif // !OVERFLOWITEM_H
