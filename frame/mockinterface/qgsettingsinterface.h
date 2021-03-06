#ifndef QGSETTINGSINTERFACE_H
#define QGSETTINGSINTERFACE_H

#include <QVariant>
#include <QStringList>

class QGSettings;
class QGSettingsInterface
{
public:
    enum Type {
        REAL,   // 持有真正的QGSettings指针
        FAKE    // Mock类
    };

    virtual ~QGSettingsInterface() {}

    virtual Type type() = 0;
    virtual QGSettings *gsettings() = 0;
    virtual QVariant get(const QString &key) const = 0;
    virtual void set(const QString &key, const QVariant &value) = 0;
    virtual bool trySet(const QString &key, const QVariant &value) = 0;
    virtual QStringList keys() const = 0;
    virtual QVariantList choices(const QString &key) const = 0;
    virtual void reset(const QString &key) = 0;
    static bool isSchemaInstalled(const QByteArray &schema_id) {Q_UNUSED(schema_id); return false;}

};
#endif // QGSETTINGSINTERFACE_H
