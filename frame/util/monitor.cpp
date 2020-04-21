/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
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

#include "monitor.h"

const double DoubleZero = 0.000001;

Monitor::Monitor(QObject *parent)
    : QObject(parent)
    , m_scale(-1.0)
    , m_brightness(1.0)
{

}

void Monitor::setX(const int x)
{
    if (m_x == x)
        return;

    m_x = x;

    Q_EMIT xChanged(m_x);
    Q_EMIT geometryChanged();
}

void Monitor::setY(const int y)
{
    if (m_y == y)
        return;

    m_y = y;

    Q_EMIT yChanged(m_y);
    Q_EMIT geometryChanged();
}

void Monitor::setW(const int w)
{
    if (m_w == w)
        return;

    m_w = w;

    Q_EMIT wChanged(m_w);
    Q_EMIT geometryChanged();
}

void Monitor::setH(const int h)
{
    if (m_h == h)
        return;

    m_h = h;

    Q_EMIT hChanged(m_h);
    Q_EMIT geometryChanged();
}

void Monitor::setMmWidth(const uint mmWidth)
{
    m_mmWidth = mmWidth;
}

void Monitor::setMmHeight(const uint mmHeight)
{
    m_mmHeight = mmHeight;
}

void Monitor::setScale(const double scale)
{
    if (fabs(m_scale - scale) < DoubleZero)
        return;

    m_scale = scale;

    Q_EMIT scaleChanged(m_scale);
}

void Monitor::setPrimary(const QString &primaryName)
{
    m_primary = primaryName;
}

void Monitor::setRotate(const quint16 rotate)
{
    if (m_rotate == rotate)
        return;

    m_rotate = rotate;

    Q_EMIT rotateChanged(m_rotate);
}

void Monitor::setBrightness(const double brightness)
{
    if (fabs(m_brightness - brightness) < DoubleZero)
        return;

    m_brightness = brightness;

    Q_EMIT brightnessChanged(m_brightness);
}

void Monitor::setName(const QString &name)
{
    m_name = name;
}

void Monitor::setPath(const QString &path)
{
    m_path = path;
}

void Monitor::setRotateList(const QList<quint16> &rotateList)
{
    m_rotateList = rotateList;
}

void Monitor::setCurrentMode(const Resolution &resolution)
{
    if (m_currentMode == resolution)
        return;

    m_currentMode = resolution;

    Q_EMIT currentModeChanged(m_currentMode);
}

bool compareResolution(const Resolution &first, const Resolution &second)
{
    long firstSum = long(first.width()) * first.height();
    long secondSum = long(second.width()) * second.height();
    if (firstSum > secondSum)
        return true;
    else if (firstSum == secondSum) {
        if (first.rate() - second.rate() > 0.000001)
            return true;
        else
            return false;
    } else
        return false;

}

void Monitor::setModeList(const ResolutionList &modeList)
{
    m_modeList.clear();

    Resolution preResolution;
    // NOTE: ignore resolution less than 1024x768
    for (auto m : modeList) {
        if (m.width() >= 1024 && m.height() >= 768) {
            m_modeList.append(m);
        }
    }
    qSort(m_modeList.begin(), m_modeList.end(), compareResolution);

    Q_EMIT modelListChanged(m_modeList);
}

void Monitor::setMonitorEnable(bool enable)
{
    if (m_enable == enable)
        return;

    m_enable = enable;
    Q_EMIT enableChanged(enable);
}


bool Monitor::isSameResolution(const Resolution &r1, const Resolution &r2)
{
    return r1.width() == r2.width() && r1.height() == r2.height();
}

bool Monitor::isSameRatefresh(const Resolution &r1, const Resolution &r2)
{
    return fabs(r1.rate() - r2.rate()) < 0.000001;
}

bool Monitor::hasResolution(const Resolution &r)
{
    for (auto m : m_modeList) {
        if (isSameResolution(m, r)) {
            return true;
        }
    }

    return false;
}

bool Monitor::hasResolutionAndRate(const Resolution &r)
{
    for (auto m : m_modeList) {
        if (fabs(m.rate() - r.rate()) < 0.000001 &&
                m.width() == r.width() &&
                m.height() == r.height()) {
            return true;
        }
    }

    return false;
}
