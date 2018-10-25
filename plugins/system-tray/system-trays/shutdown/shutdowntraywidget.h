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

#ifndef SHUTDOWNTRAYWIDGET_H
#define SHUTDOWNTRAYWIDGET_H

#include "../abstractsystemtraywidget.h"
#include "../widgets/tipswidget.h"

#include <QWidget>
#include <QTimer>

class ShutdownTrayWidget : public AbstractSystemTrayWidget
{
    Q_OBJECT

public:
    explicit ShutdownTrayWidget(QWidget *parent = nullptr);

public:
    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

    QWidget *trayTipsWidget() Q_DECL_OVERRIDE;
    const QString trayClickCommand() Q_DECL_OVERRIDE;

    const QString contextMenu() const Q_DECL_OVERRIDE;
    void invokedMenuItem(const QString &menuId, const bool checked) Q_DECL_OVERRIDE;

protected:
    QSize sizeHint() const Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;

private:
    TipsWidget *m_tipsLabel;

    QPixmap m_pixmap;
};

#endif // SHUTDOWNTRAYWIDGET_H
