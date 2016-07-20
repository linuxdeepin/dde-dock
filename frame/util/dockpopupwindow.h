#ifndef DOCKPOPUPWINDOW_H
#define DOCKPOPUPWINDOW_H

#include "dbus/dbusxmousearea.h"

#include <darrowrectangle.h>

class DockPopupWindow : public Dtk::Widget::DArrowRectangle
{
    Q_OBJECT

public:
    explicit DockPopupWindow(QWidget *parent = 0);
    ~DockPopupWindow();

    bool model() const;

public slots:
    void show(const QPoint &pos, const bool model = false);
    void hide();

signals:
    void accept() const;

protected:
    void mousePressEvent(QMouseEvent *e);

private slots:
    void globalMouseRelease(int button, int x, int y, const QString &id);

private:
    bool m_model;
    QString m_mouseAreaKey;

    QTimer *m_acceptDelayTimer;

    DBusXMouseArea *m_mouseInter;
};

#endif // DOCKPOPUPWINDOW_H
