#include <QFile>

#include "appicon.h"
#include "Controller/themeiconmanager.h"

AppIcon::AppIcon(QWidget *parent, Qt::WindowFlags f) :
    QLabel(parent, f)
{
    this->setAttribute(Qt::WA_TranslucentBackground);
    this->setAlignment(Qt::AlignCenter);
}

void AppIcon::setIcon(const QString &iconPath)
{
    ThemeIconManager * themeIconManager = ThemeIconManager::instance();

    QIcon icon = themeIconManager->getIcon(iconPath);
    QPixmap pixmap;

    if (icon.isNull()) {
        // iconPath is an absolute path of the system.
        if (QFile::exists(iconPath)) {
            pixmap = QPixmap(iconPath);
        } else if (iconPath.startsWith("data:image/")){
            // iconPath is a string representing an inline image.
            QStringList strs = iconPath.split("base64,");
            if (strs.length() == 2) {
                QByteArray data = QByteArray::fromBase64(strs.at(1).toLatin1());
                pixmap.loadFromData(data);
            }
        }
    } else {
        pixmap = icon.pixmap(48, 48);
    }

    if (!pixmap.isNull()) {
        pixmap = pixmap.scaled(m_modeData->getAppIconSize(),
                               m_modeData->getAppIconSize(),
                               Qt::KeepAspectRatioByExpanding,
                               Qt::SmoothTransformation);

        setPixmap(pixmap);
    }
}
