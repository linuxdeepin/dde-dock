#ifndef PLUGINWIDGET_H
#define PLUGINWIDGET_H

#include "constants.h"
#include "dbus/dbuspower.h"

#include <QWidget>

class PluginWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PluginWidget(QWidget *parent = 0);

public slots:
    void displayModeChanged();

protected:
    QSize sizeHint() const;
    void resizeEvent(QResizeEvent *e);
    void paintEvent(QPaintEvent *e);

private:
    const QPixmap loadSvg(const QString &fileName, const QSize &size) const;

private:
    void refershIconPixmap();

private:
    Dock::DisplayMode m_displayMode;
    QPixmap m_iconPixmap;

    DBusPower *m_powerInter;
};

#endif // PLUGINWIDGET_H
