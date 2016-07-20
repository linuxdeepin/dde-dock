#include "launcheritem.h"
#include "util/themeappicon.h"
#include "util/imagefactory.h"

#include <QPainter>
#include <QProcess>
#include <QMouseEvent>

LauncherItem::LauncherItem(QWidget *parent)
    : DockItem(DockItem::Launcher, parent),

      m_tips(new QLabel(this))
{
    m_tips->setVisible(false);
    m_tips->setText("Launcher");
    m_tips->setStyleSheet("color:white;"
                          "padding:5px 10px;");
}

void LauncherItem::paintEvent(QPaintEvent *e)
{
    DockItem::paintEvent(e);

    if (!isVisible())
        return;

    QPainter painter(this);

    const QPixmap pixmap = DockDisplayMode == Fashion ? m_largeIcon : m_smallIcon;
    painter.drawPixmap(rect().center() - pixmap.rect().center(), pixmap);

    if (m_hover)
        painter.drawPixmap(rect().center() - pixmap.rect().center(), ImageFactory::lighterEffect(pixmap));
}

void LauncherItem::resizeEvent(QResizeEvent *e)
{
    DockItem::resizeEvent(e);

    const int iconSize = qMin(width(), height());
    if (DockDisplayMode == Efficient)
    {
        m_smallIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.7);
        m_largeIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.9);
    } else {
        m_smallIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.6);
        m_largeIcon = ThemeAppIcon::getIcon("deepin-launcher", iconSize * 0.8);
    }
}

void LauncherItem::mousePressEvent(QMouseEvent *e)
{
    if (e->button() == Qt::RightButton && !perfectIconRect().contains(e->pos()))
        return QWidget::mousePressEvent(e);

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

QWidget *LauncherItem::popupTips()
{
    return m_tips;
}
