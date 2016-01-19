#include <QDir>
#include <QDebug>
#include <QSettings>
#include <QApplication>
#include "stylemanager.h"

const QString GLOBAL_STYLE_PATH = "/usr/share/dde-dock/style/";
const QString LOCAL_STYLE_PATH = QDir::homePath() + "/.config/deepin/dde-dock/style/";
const QString CURRENT_STYLE_KEY = "CurrentStyle";
const QString SEARCH_PATH_KEY = "StyleSearchPath";

StyleManager::StyleManager(QObject *parent) : QObject(parent)
{
    initSettings();
}

StyleManager *StyleManager::m_styleManager = NULL;
StyleManager *StyleManager::instance()
{
    if (m_styleManager == NULL)
        m_styleManager = new StyleManager;
    return m_styleManager;
}

QStringList StyleManager::styleNameList()
{
    QStringList nameList;
    nameList << "dark" << "light";  //default style
    nameList << getStyleFromFilesystem();

    return nameList;
}

QString StyleManager::currentStyle()
{
    return m_settings->value(CURRENT_STYLE_KEY).toString();
}

void StyleManager::applyStyle(const QString &styleName)
{
    if (styleName == "dark" || styleName == "light") {
        applyDefaultStyle(styleName);
    }
    else {
        applyThirdPartyStyle(styleName);
    }
}

void StyleManager::initStyleSheet()
{
    applyStyle(currentStyle());
}

QStringList StyleManager::getStyleFromFilesystem()
{
    QStringList list;
    for (QString path : m_settings->value(SEARCH_PATH_KEY).toStringList()) {
        QDir d(path);
        if (d.exists()) {   //read all valuable style
            QFileInfoList nl = d.entryInfoList(QDir::Dirs | QDir::NoDotAndDotDot |QDir::Readable);
            for (QFileInfo p : nl) {
                if (QFile::exists(p.absoluteFilePath() + "/style.qss"))
                    list << p.baseName();
            }
        }
    }

    return list;
}

void StyleManager::initSettings()
{
    m_settings = new QSettings("deepin", "dde-dock");
    if (m_settings->value(CURRENT_STYLE_KEY).toString().isEmpty()) {
        m_settings->setValue(CURRENT_STYLE_KEY, "dark");
        QStringList p;
        p << GLOBAL_STYLE_PATH;
        p << LOCAL_STYLE_PATH;
        m_settings->setValue(SEARCH_PATH_KEY, QVariant(p));
    }
}

void StyleManager::applyDefaultStyle(const QString &name)
{
    QString filePath = QString("://qss/resources/%1/qss/dde-dock.qss").arg(name);
    QFile file(filePath);
    if (file.open(QFile::ReadOnly)) {
        QString styleSheet = QLatin1String(file.readAll());
        qApp->setStyleSheet(styleSheet);
        file.close();
    } else {
        qWarning() << "Dock Open  style file errr!";
    }
}

void StyleManager::applyThirdPartyStyle(const QString &name)
{
    QStringList sp = m_settings->value(SEARCH_PATH_KEY).toStringList();
    for (QString path : sp) {
        QFile file(path + name + "/style.qss");
        if (file.open(QFile::ReadOnly)) {
            QString styleSheet = QLatin1String(file.readAll());
            qApp->setStyleSheet(styleSheet);
            file.close();
            return;
        } else {
            qWarning() << "Dock Open  style file errr!";
        }
    }
}


