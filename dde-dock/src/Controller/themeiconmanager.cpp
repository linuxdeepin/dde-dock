#include <QDebug>
#include <QVariant>
#include <QGlobalStatic>

#include <QGSettings>

#include "themeiconmanager.h"

static const QString KeyIconThemeName = "icon-theme-name";

class ThemeIconManagerPrivate : public ThemeIconManager {

};

Q_GLOBAL_STATIC(ThemeIconManagerPrivate, ThemeIconManagerStatic)

ThemeIconManager * ThemeIconManager::instance()
{
    return ThemeIconManagerStatic;
}

QString ThemeIconManager::getTheme() const
{
    return m_theme;
}

QIcon ThemeIconManager::getIcon(QString iconName)
{
    return QIcon::fromTheme(iconName);
}

// private methods
ThemeIconManager::ThemeIconManager(QObject *parent) :
    QObject(parent)
{
    m_gsettings = new QGSettings("com.deepin.xsettings",
                                 "/com/deepin/xsettings/",
                                 this);
    setTheme(m_gsettings->get(KeyIconThemeName).toString());

    connect(m_gsettings, &QGSettings::changed, this, &ThemeIconManager::settingsChanged);
}

void ThemeIconManager::setTheme(const QString theme)
{
    m_theme = theme;
    QIcon::setThemeName(theme);
}

void ThemeIconManager::settingsChanged(const QString & key)
{
    if (key == KeyIconThemeName) {
        setTheme(m_gsettings->get(KeyIconThemeName).toString());

        emit themeChanged(m_theme);
    }
}
