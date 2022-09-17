// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#pragma once
#include <memory>

using namespace std;

template <class T>
class Singleton
{
public:
    static inline T *instance() {
        static T*  _instance = new T;
        return _instance;
    }

protected:
    Singleton(void) {}
    ~Singleton(void) {}
    Singleton(const Singleton &) {}
    Singleton &operator= (const Singleton &) {}
};


