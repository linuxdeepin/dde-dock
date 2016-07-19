#ifndef DISKINFO_H
#define DISKINFO_H

#include <QString>
#include <QDataStream>
#include <QDebug>
#include <QtDBus>

class DiskInfo
{
public:
    DiskInfo();
    static void registerMetaType();

    friend QDebug operator<<(QDebug debug, const DiskInfo &info);
    friend QDBusArgument &operator<<(QDBusArgument &args, const DiskInfo &info);
    friend QDataStream &operator<<(QDataStream &args, const DiskInfo &info);
    friend const QDBusArgument &operator>>(const QDBusArgument &args, DiskInfo &info);
    friend const QDataStream &operator>>(QDataStream &args, DiskInfo &info);

public:
    QString m_id;
    QString m_name;
    QString m_type;
    QString m_path;
    QString m_mountPoint;
    QString m_icon;

    bool m_unmountable;
    bool m_ejectable;

    quint64 m_usedSize;
    quint64 m_totalSize;
};

typedef QList<DiskInfo> DiskInfoList;

Q_DECLARE_METATYPE(DiskInfo)
Q_DECLARE_METATYPE(DiskInfoList)

#endif // DISKINFO_H
