#ifndef PLUGINWIDGET_H
#define PLUGINWIDGET_H

#include "constants.h"
#include "dbus/dbuspower.h"

#include <QWidget>
#include <QTimer>

class PluginWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PluginWidget(QWidget *parent = 0);

protected:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);

private:
    const QPixmap loadSvg(const QString &fileName, const QSize &size) const;

private:
    void refershIconPixmap();

private:
    bool m_hover;
    Dock::DisplayMode m_displayMode;

    DBusPower *m_powerInter;
};

#endif // PLUGINWIDGET_H
