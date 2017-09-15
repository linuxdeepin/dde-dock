#include "themeappicon.h"

#include <QIcon>
#include <QFile>
#include <QDebug>
#include <QApplication>

ThemeAppIcon::ThemeAppIcon(QObject *parent) : QObject(parent)
{

}

ThemeAppIcon::~ThemeAppIcon()
{

}

const QPixmap ThemeAppIcon::getIcon(const QString iconName, const int size)
{
    const auto ratio = qApp->devicePixelRatio();
    const int s = int(size * ratio) & ~1;

    QPixmap pixmap;

    do {

        if (iconName.startsWith("data:image/"))
        {
            const QStringList strs = iconName.split("base64,");
            if (strs.size() == 2)
                pixmap.loadFromData(QByteArray::fromBase64(strs.at(1).toLatin1()));

            if (!pixmap.isNull())
                break;
        }

        if (QFile::exists(iconName))
        {
            pixmap = QPixmap(iconName);
            if (!pixmap.isNull())
                break;
        }

        const QIcon icon = QIcon::fromTheme(iconName, QIcon::fromTheme("application-x-desktop"));
        pixmap = icon.pixmap(QSize(s, s));
        if (!pixmap.isNull())
            break;

        pixmap = QPixmap(":/icons/resources/application-x-desktop.svg");
        if (!pixmap.isNull())
            break;

        Q_UNREACHABLE();

    } while (false);

    pixmap = pixmap.scaled(s, s, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}

