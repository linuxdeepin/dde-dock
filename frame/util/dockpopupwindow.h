#ifndef DOCKPOPUPWINDOW_H
#define DOCKPOPUPWINDOW_H

#include "dbus/dbusxmousearea.h"
#include "dbus/dbusdisplay.h"

#include <darrowrectangle.h>
#include <DWindowManagerHelper>

DWIDGET_USE_NAMESPACE

class DockPopupWindow : public Dtk::Widget::DArrowRectangle
{
    Q_OBJECT

public:
    explicit DockPopupWindow(QWidget *parent = 0);
    ~DockPopupWindow();

    bool model() const;

    void setContent(QWidget *content);

public slots:
    void show(const QPoint &pos, const bool model = false);
    void show(const int x, const int y);
    void hide();

signals:
    void accept() const;

protected:
    void showEvent(QShowEvent *e);
    void enterEvent(QEvent *e);
    void mousePressEvent(QMouseEvent *e);
    bool eventFilter(QObject *o, QEvent *e);

private slots:
    void globalMouseRelease(int button, int x, int y, const QString &id);
    void registerMouseEvent();
    void unRegisterMouseEvent();
    void compositeChanged();

private:
    bool m_model;
    QPoint m_lastPoint;
    QString m_mouseAreaKey;

    QTimer *m_acceptDelayTimer;

    DBusXMouseArea *m_mouseInter;
    DBusDisplay *m_displayInter;
    DWindowManagerHelper *m_wmHelper;
};

#endif // DOCKPOPUPWINDOW_H
