#include "launcheritem.h"
#include "util/themeappicon.h"

#include <QPainter>
#include <QProcess>
#include <QMouseEvent>

LauncherItem::LauncherItem(QWidget *parent)
    : DockItem(DockItem::Launcher, parent)
{
}

void LauncherItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_icon.rect().center(), m_icon);
}

void LauncherItem::resizeEvent(QResizeEvent *e)
{
    DockItem::resizeEvent(e);

    m_icon = ThemeAppIcon::getIcon("deepin-launcher", qMin(width(), height()) * 0.7);
}

void LauncherItem::mousePressEvent(QMouseEvent *e)
{
    if (e->button() != Qt::LeftButton)
        return;

    QProcess *proc = new QProcess;

    connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

    QStringList args = QStringList() << "--print-reply"
                                     << "--dest=com.deepin.dde.Launcher"
                                     << "/com/deepin/dde/Launcher"
                                     << "com.deepin.dde.Launcher.Toggle";

    proc->start("dbus-send", args);
}
