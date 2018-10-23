#ifndef POWERTRAYLOADER_H
#define POWERTRAYLOADER_H

#include "../abstracttrayloader.h"
#include "dbus/dbuspower.h"

#include <QObject>

class PowerTrayLoader : public AbstractTrayLoader
{
    Q_OBJECT
public:
    explicit PowerTrayLoader(QObject *parent = nullptr);

public Q_SLOTS:
    void load() Q_DECL_OVERRIDE;

private:
    DBusPower *m_powerInter;
};

#endif // POWERTRAYLOADER_H
