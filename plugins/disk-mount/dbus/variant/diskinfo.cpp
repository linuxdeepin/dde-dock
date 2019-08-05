/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "diskinfo.h"

DiskInfo::DiskInfo()
{

}

void DiskInfo::registerMetaType()
{
    qRegisterMetaType<DiskInfo>("DiskInfo");
    qDBusRegisterMetaType<DiskInfo>();

    qRegisterMetaType<DiskInfoList>("DiskInfoList");
    qDBusRegisterMetaType<DiskInfoList>();
}

QDebug operator<<(QDebug debug, const DiskInfo &info)
{
    debug << info.m_id << info.m_name << info.m_type << info.m_path << info.m_mountPoint << info.m_icon;
    debug << '\t' << info.m_unmountable << '\t' << info.m_ejectable;
    debug << '\t' << info.m_usedSize << '\t' << info.m_totalSize;
    debug << endl;

    return debug;
}

const QDataStream &operator>>(QDataStream &args, DiskInfo &info)
{
    args >> info.m_id >> info.m_name >> info.m_type >> info.m_path >> info.m_mountPoint >> info.m_icon;
    args >> info.m_unmountable >> info.m_ejectable;
    args >> info.m_usedSize >> info.m_totalSize;

    return args;
}

const QDBusArgument &operator>>(const QDBusArgument &args, DiskInfo &info)
{
    args.beginStructure();
    args >> info.m_id >> info.m_name >> info.m_type >> info.m_path >> info.m_mountPoint >> info.m_icon;
    args >> info.m_unmountable >> info.m_ejectable;
    args >> info.m_usedSize >> info.m_totalSize;
    args.endStructure();

    return args;
}

QDataStream &operator<<(QDataStream &args, const DiskInfo &info)
{
    args << info.m_id << info.m_name << info.m_type << info.m_path << info.m_mountPoint << info.m_icon;
    args << info.m_unmountable << info.m_ejectable;
    args << info.m_usedSize << info.m_totalSize;

    return args;
}

QDBusArgument &operator<<(QDBusArgument &args, const DiskInfo &info)
{
    args.beginStructure();
    args << info.m_id << info.m_name << info.m_type << info.m_path << info.m_mountPoint << info.m_icon;
    args << info.m_unmountable << info.m_ejectable;
    args << info.m_usedSize << info.m_totalSize;
    args.endStructure();

    return args;
}
