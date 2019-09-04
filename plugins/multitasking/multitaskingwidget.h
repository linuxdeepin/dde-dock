/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     wangshaojun <wangshaojun_cm@deepin.com>
 *
 * Maintainer: wangshaojun <wangshaojun_cm@deepin.com>
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

#ifndef MultitaskingWidget_H
#define MultitaskingWidget_H

#include <QWidget>

#define PLUGIN_KEY "multitasking"

class MultitaskingWidget : public QWidget
{
    Q_OBJECT

public:
    explicit MultitaskingWidget(QWidget *parent = 0);
    void refreshIcon();
    QSize sizeHint() const override;

signals:
    void requestContextMenu(const QString &itemKey) const;

protected:
    void paintEvent(QPaintEvent *e) override;
    void resizeEvent(QResizeEvent *event) override;
};

#endif // MULTITASKINGWIDGET_H
