#ifndef APPSNAPSHOT_H
#define APPSNAPSHOT_H

#include <QWidget>

class AppSnapshot : public QWidget
{
    Q_OBJECT

public:
    explicit AppSnapshot(QWidget *parent = 0);
};

#endif // APPSNAPSHOT_H
