#ifndef APPICON_H
#define APPICON_H

#include <QDebug>
#include <QLabel>
#include <QObject>
#include <QWidget>
#include <QPixmap>
#include <QMouseEvent>

#include "controller/dockmodedata.h"

class AppIcon : public QLabel
{
    Q_OBJECT
public:
    explicit AppIcon(QWidget *parent = 0, Qt::WindowFlags f = 0);

    void setIcon(const QString &iconPath);

signals:
    void mousePress(QMouseEvent *event);
    void mouseRelease(QMouseEvent *event);
    void mouseEnter();
    void mouseLeave();

protected:
    void mousePressEvent(QMouseEvent *ev);
    void mouseReleaseEvent(QMouseEvent *ev);
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);

private:
    DockModeData *m_modeData = DockModeData::instance();
    QString m_iconPath = "";

    QString getThemeIconPath(QString iconName);
    void updateIcon();
};

#endif // APPICON_H
