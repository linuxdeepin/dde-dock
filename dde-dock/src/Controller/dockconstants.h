#ifndef DOCKCONSTANTS_H
#define DOCKCONSTANTS_H

#include <QObject>

class DockConstants : public QObject
{
    Q_OBJECT
public:
    explicit DockConstants(QObject *parent = 0);

    enum DockMode {
        FashionMode,
        EfficientMode,
        ClassicMode
    };

signals:

public slots:
};

#endif // DOCKCONSTANTS_H
