#include <QGSettings>
#include <QVariant>

#include "qgsettingsinterfacemock.h"

QGSettingsInterfaceMock::QGSettingsInterfaceMock(const QByteArray &schema_id, const QByteArray &path, QObject *parent)
{

}

QGSettingsInterfaceMock::~QGSettingsInterfaceMock()
{

}

QGSettingsInterface::Type QGSettingsInterfaceMock::type()
{
    return Type::FAKE;
}

QGSettings *QGSettingsInterfaceMock::gsettings()
{
    return nullptr;
}

QVariant QGSettingsInterfaceMock::get(const QString &key) const
{
    return QVariant();
}

void QGSettingsInterfaceMock::set(const QString &key, const QVariant &value)
{

}

bool QGSettingsInterfaceMock::trySet(const QString &key, const QVariant &value)
{
    return false;
}

QStringList QGSettingsInterfaceMock::keys() const
{
    return QStringList();
}

QVariantList QGSettingsInterfaceMock::choices(const QString &key) const
{
    return QVariantList();
}

void QGSettingsInterfaceMock::reset(const QString &key)
{

}

bool QGSettingsInterfaceMock::isSchemaInstalled(const QByteArray &schema_id)
{
    return false;
}
