#ifndef STYLEMANAGER_H
#define STYLEMANAGER_H

#include <QObject>
#include <QSettings>

class StyleManager : public QObject
{
    Q_OBJECT
public:
    static StyleManager *instance();
    QStringList styleNameList();
    QString currentStyle();
    void applyStyle(const QString &styleName);
    void initStyleSheet();

private:
    explicit StyleManager(QObject *parent = 0);
    QStringList getStyleFromFilesystem();
    void initSettings();
    void applyDefaultStyle(const QString &name);
    void applyThirdPartyStyle(const QString &name);

private:
    static StyleManager *m_styleManager;
    QSettings *m_settings;
};

#endif // STYLEMANAGER_H
