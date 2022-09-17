// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef FASHIONTRAYCONTROLWIDGET_H
#define FASHIONTRAYCONTROLWIDGET_H

#include "constants.h"

#include <QLabel>

/**
 * @brief The FashionTrayControlWidget class
 * @note  系统托盘上的展开合并按钮
 */
class FashionTrayControlWidget : public QWidget
{
    Q_OBJECT

public:
    explicit FashionTrayControlWidget(Dock::Position position, QWidget *parent = nullptr);

    void setDockPostion(Dock::Position pos);

    bool expanded() const;
    void setExpanded(const bool &expanded);

Q_SIGNALS:
    void expandChanged(const bool expanded);

protected:
    void paintEvent(QPaintEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

private:
    void refreshArrowPixmap();

private:
    QTimer *m_expandDelayTimer;
    QPixmap m_arrowPix;

    Dock::Position m_dockPosition;
    bool m_expanded;
    bool m_hover;
    bool m_pressed;
};

#endif // FASHIONTRAYCONTROLWIDGET_H
