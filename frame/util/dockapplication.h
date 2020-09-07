/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
 *
 * Author:     liuxing <liuxing@uniontech.com>
 *
 * Maintainer: liuxing <liuxing@uniontech.com>
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

#ifndef DOCKAPPLICATION_H
#define DOCKAPPLICATION_H

#include <DApplication>

DWIDGET_USE_NAMESPACE
#ifdef DCORE_NAMESPACE
DCORE_USE_NAMESPACE
#else
DUTIL_USE_NAMESPACE
#endif

/**
 * @brief The DockApplication class
 * 本类通过重写application的notify函数监控应用的鼠标事件，判断是否为触屏状态
 */
class DockApplication : public DApplication
{
    Q_OBJECT
public:
    explicit DockApplication(int &argc, char **argv);
    virtual bool notify(QObject *obj, QEvent *event) override;

signals:

public slots:
};

#endif // DOCKAPPLICATION_H
