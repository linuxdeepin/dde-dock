#ifndef APPICON_H
#define APPICON_H

#include <QObject>
#include <QWidget>
#include <QLabel>
#include <QPixmap>
#include <Controller/dockmodedata.h>
#include <QDebug>

class AppIcon : public QLabel
{
    Q_OBJECT
public:
    explicit AppIcon(QWidget *parent = 0,Qt::WindowFlags f = 0);

    void setIcon(const QString &iconPath);
    QString getIconPath() const;

signals:

public slots:

    QString getSysIcon(const QString &iconName, int size = 48);
private:
    DockModeData *m_modeData = DockModeData::instance();
    QString m_iconPath;
};

#endif // APPICON_H
