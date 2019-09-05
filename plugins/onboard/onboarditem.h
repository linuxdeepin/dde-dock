/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#ifndef ONBOARDITEM_H
#define ONBOARDITEM_H

#include "constants.h"

#include <QWidget>

class OnboardItem : public QWidget
{
    Q_OBJECT

public:
    explicit OnboardItem(QWidget *parent = nullptr);

protected:
    QSize sizeHint() const override;
    void paintEvent(QPaintEvent *e) override;

private:
    const QPixmap loadSvg(const QString &fileName, const QSize &size) const;
    void mousePressEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *event) override;
    void leaveEvent(QEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

private:
    Dock::DisplayMode m_displayMode;
    bool m_hover;
    bool m_pressed;
};

#endif // ONBOARDITEM_H
