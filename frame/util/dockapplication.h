// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
};

#endif // DOCKAPPLICATION_H
