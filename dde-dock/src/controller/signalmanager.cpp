/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "signalmanager.h"

SignalManager *SignalManager::m_signalManager = NULL;
SignalManager *SignalManager::instance()
{
    if (!m_signalManager)
        m_signalManager = new SignalManager;
    return m_signalManager;
}

SignalManager::SignalManager(QObject *parent) : QObject(parent)
{

}

