#ifndef DOCKPLUGINSSETTINGWINDOW_H
#define DOCKPLUGINSSETTINGWINDOW_H

#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QMouseEvent>
#include <QPushButton>
#include <QFrame>
#include <QLabel>
#include <QDebug>

#include <libdui/dswitchbutton.h>
#include <libdui/dseparatorhorizontal.h>

DUI_USE_NAMESPACE

class PluginSettingLine;

class DockPluginsSettingWindow : public QFrame
{
    Q_OBJECT
public:
    explicit DockPluginsSettingWindow(QWidget *parent = 0);

public slots:
    void onPluginAdd(bool checked = false,
                     const QString &id = "",
                     const QString &title = "",
                     const QPixmap &icon = QPixmap());
    void onPluginRemove(const QString &id);
    void onPluginEnabledChanged(const QString &id, bool enabled);
    void onPluginTitleChanged(const QString &id, const QString &title);

signals:
    void checkedChanged(const QString &id, bool checked);

protected:
    void mouseMoveEvent(QMouseEvent *event);
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);
    bool eventFilter(QObject *obj, QEvent *event);

private:
    void resizeWithLineCount();
    void initCloseTitle();

    bool m_mousePressed;
    QPoint m_pressPosition;
    QVBoxLayout *m_mainLayout;
    QMap<QString, PluginSettingLine *> m_lineMap;

    const int ICON_SIZE = 24;
    const int CONTENT_MARGIN = 6;
    const int LINE_SPACING = 5;
    const int LINE_HEIGHT = 30;
    const int WIN_WIDTH = 230;

};

#endif // DOCKPLUGINSSETTINGWINDOW_H
