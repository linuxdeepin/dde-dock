#ifndef SHUTDOWNTRAYLOADER_H
#define SHUTDOWNTRAYLOADER_H

#include "../abstracttrayloader.h"

#include <QObject>

class ShutdownTrayLoader : public AbstractTrayLoader
{
    Q_OBJECT
public:
    explicit ShutdownTrayLoader(QObject *parent = nullptr);

public Q_SLOTS:
    void load() Q_DECL_OVERRIDE;
};

#endif // SHUTDOWNTRAYLOADER_H
