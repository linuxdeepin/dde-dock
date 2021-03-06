#ifndef QGSETTINGSINTERFACEIMPL_H
#define QGSETTINGSINTERFACEIMPL_H
#include <QObject>

#include "qgsettingsinterface.h"

class QGSettings;
class QGSettingsInterfaceImpl : public QGSettingsInterface
{
public:
    QGSettingsInterfaceImpl(const QByteArray &schema_id, const QByteArray &path = QByteArray(), QObject *parent = nullptr);
    ~QGSettingsInterfaceImpl() override;

    virtual Type type() override;
    virtual QGSettings *gsettings() override;
    virtual QVariant get(const QString &key) const override;
    virtual void set(const QString &key, const QVariant &value) override;
    virtual bool trySet(const QString &key, const QVariant &value) override;
    virtual QStringList keys() const override;
    virtual QVariantList choices(const QString &key) const override;
    virtual void reset(const QString &key) override;
    static bool isSchemaInstalled(const QByteArray &schema_id);

private:
    QGSettings *m_gsettings;
};

#endif // QGSETTINGSINTERFACEIMPL_H
