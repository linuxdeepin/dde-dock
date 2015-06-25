#ifndef APPICON_H
#define APPICON_H

#include <QObject>
#include <QWidget>
#include <QLabel>
#include <QPixmap>

class AppIcon : public QLabel
{
    Q_OBJECT
public:
    explicit AppIcon(QWidget *parent = 0,Qt::WindowFlags f = 0);
    explicit AppIcon(QString iconPath, QWidget *parent = 0,Qt::WindowFlags f = 0);//default search system icon path

    void setIcon(const QString &iconPath);
signals:

public slots:

    QString getSysIcon(const QString &iconName, int size = 42);
};

#endif // APPICON_H
