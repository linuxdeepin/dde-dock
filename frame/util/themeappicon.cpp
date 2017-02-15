#include "themeappicon.h"

#include <QIcon>

ThemeAppIcon::ThemeAppIcon(QObject *parent) : QObject(parent)
{

}

ThemeAppIcon::~ThemeAppIcon()
{

}

QPixmap ThemeAppIcon::getIcon(const QString iconName, const int size)
{
    QIcon icon = QIcon::fromTheme(iconName, QIcon::fromTheme("application-x-desktop"));
    QPixmap pix = icon.pixmap(QSize(size, size));
    if (pix.isNull()) {
        pix = QPixmap(":/icons/resources/application-x-desktop.svg").scaled(size, size);
    }

    return pix;
}

