#ifndef THEMEAPPICON_H
#define THEMEAPPICON_H

#include <QObject>

class ThemeAppIcon : public QObject
{
    Q_OBJECT
public:
    explicit ThemeAppIcon(QObject *parent = 0);
    ~ThemeAppIcon();

    static void gtkInit();

    static QPixmap getIconPixmap(QString iconPath, int width=64, int height=64);
    static QString getThemeIconPath(QString iconName, int size=64);
    static QPixmap getIcon(const QString iconName, const int size);
    static QPixmap loadSvg(const QString &fileName, const int size);

signals:

public slots:
};

#endif // THEMEAPPICON_H
