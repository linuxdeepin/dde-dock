#ifndef SOUNDTRAYLOADER_H
#define SOUNDTRAYLOADER_H

#include "../abstracttrayloader.h"

#include <QObject>

class SoundTrayLoader : public AbstractTrayLoader
{
    Q_OBJECT
public:
    explicit SoundTrayLoader(QObject *parent = nullptr);

public Q_SLOTS:
    void load() Q_DECL_OVERRIDE;
};

#endif // SOUNDTRAYLOADER_H
