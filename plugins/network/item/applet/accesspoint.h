#ifndef ACCESSPOINT_H
#define ACCESSPOINT_H

#include <QObject>
#include <QJsonObject>

class AccessPoint : public QObject
{
    Q_OBJECT

public:
    explicit AccessPoint(const QJsonObject &apInfo);
    AccessPoint(const AccessPoint &ap);
    bool operator==(const AccessPoint &ap) const;
    bool operator>(const AccessPoint &ap) const;
    AccessPoint &operator=(const AccessPoint &ap);

    const QString ssid() const;
    int strength() const;
    bool secured() const;

private:
    int m_strength;
    bool m_secured;
    bool m_securedInEap;
    QString m_path;
    QString m_ssid;
};

#endif // ACCESSPOINT_H
