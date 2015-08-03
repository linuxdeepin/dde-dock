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

    if (icon.isNull()) {
        qDebug() << iconPath;
    } else {
        QPixmap pixmap = icon.pixmap(48, 48);
        pixmap = pixmap.scaled(m_modeData->getAppIconSize(),
                               m_modeData->getAppIconSize(),
                               Qt::KeepAspectRatioByExpanding,
                               Qt::SmoothTransformation);

        setPixmap(pixmap);
    }
}
