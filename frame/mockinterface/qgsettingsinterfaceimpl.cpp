#include <QGSettings>
#include <QVariant>

#include "qgsettingsinterfaceimpl.h"

QGSettingsInterfaceImpl::QGSettingsInterfaceImpl(const QByteArray &schema_id, const QByteArray &path, QObject *parent)
    : m_gsettings(new QGSettings(schema_id, path, parent))
{

}

QGSettingsInterfaceImpl::~QGSettingsInterfaceImpl()
{

}

QGSettingsInterface::Type QGSettingsInterfaceImpl::type()
{
     return Type::REAL;
}

QGSettings *QGSettingsInterfaceImpl::gsettings()
{
    return m_gsettings;
}

QVariant QGSettingsInterfaceImpl::get(const QString &key) const
{
    return m_gsettings->get(key);
}

void QGSettingsInterfaceImpl::set(const QString &key, const QVariant &value)
{
    return m_gsettings->set(key, value);
}

bool QGSettingsInterfaceImpl::trySet(const QString &key, const QVariant &value)
{
    return m_gsettings->trySet(key, value);
}

QStringList QGSettingsInterfaceImpl::keys() const
{
    return m_gsettings->keys();
}

QVariantList QGSettingsInterfaceImpl::choices(const QString &key) const
{
    return m_gsettings->choices(key);
}

void QGSettingsInterfaceImpl::reset(const QString &key)
{
    return m_gsettings->reset(key);
}

bool QGSettingsInterfaceImpl::isSchemaInstalled(const QByteArray &schema_id)
{
    return QGSettings::isSchemaInstalled(schema_id);
}
