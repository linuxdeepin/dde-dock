#ifndef THEMEICONMANAGER_H
#define THEMEICONMANAGER_H

#include <QIcon>

class QGSettings;
class ThemeIconManager : public QObject
{
    Q_OBJECT
public:
    static ThemeIconManager * instance();

    QIcon getIcon(QString iconName);

    QString getTheme() const;

signals:
    void themeChanged(QString theme);

protected:
    ThemeIconManager(QObject *parent = 0);

private:
    QString m_theme;
    QGSettings * m_gsettings;

    void setTheme(const QString theme);

private slots:
    void settingsChanged(const QString & key);
};

#endif // THEMEICONMANAGER_H
