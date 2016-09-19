
#include "constants.h"
#include "trashwidget.h"

#include <QPainter>
#include <QIcon>
#include <QApplication>
#include <QDragEnterEvent>

DWIDGET_USE_NAMESPACE

TrashWidget::TrashWidget(QWidget *parent)
    : QWidget(parent),

      m_popupApplet(new PopupControlWidget(this)),

      m_openAct(tr("Run"), this),
      m_clearAct(tr("Empty"), this)
{
    QIcon::setThemeName("deepin");

    m_popupApplet->setVisible(false);

    connect(m_popupApplet, &PopupControlWidget::emptyChanged, this, &TrashWidget::updateIcon);

    updateIcon();
    setAcceptDrops(true);
}

QWidget *TrashWidget::popupApplet()
{
    return m_popupApplet;
}

QSize TrashWidget::sizeHint() const
{
    return QSize(26, 26);
}

void TrashWidget::dragEnterEvent(QDragEnterEvent *e)
{
    if (e->mimeData()->hasFormat("RequestDock"))
        return e->accept();

    if (e->mimeData()->hasFormat("text/uri-list"))
        return e->accept();
}

void TrashWidget::dropEvent(QDropEvent *e)
{
    if (e->mimeData()->hasFormat("RequestDock"))
        return removeApp(e->mimeData()->data("AppKey"));

    if (e->mimeData()->hasFormat("text/uri-list"))
    {
        const QMimeData *mime = e->mimeData();
        for (auto url : mime->urls())
            moveToTrash(url);
    }
}

void TrashWidget::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_icon.rect().center(), m_icon);
}

void TrashWidget::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    updateIcon();
}

void TrashWidget::mousePressEvent(QMouseEvent *e)
{
    const QPoint dis = e->pos() - rect().center();
    if (e->button() != Qt::RightButton || dis.manhattanLength() > std::min(width(), height()) * 0.8 * 0.5)
        return QWidget::mousePressEvent(e);

//    showMenu();
    emit requestContextMenu();
}

const QPoint TrashWidget::popupMarkPoint()
{
    QPoint p;
    QWidget *w = this;
    do {
        p += w->pos();
        w = qobject_cast<QWidget *>(w->parent());
    } while (w);

    const QRect r = rect();
    switch (qApp->property(PROP_POSITION).value<Dock::Position>())
    {
    case Dock::Top:       p += QPoint(r.width() / 2, r.height());      break;
    case Dock::Bottom:    p += QPoint(r.width() / 2, 0);               break;
    case Dock::Left:      p += QPoint(r.width(), r.height() / 2);      break;
    case Dock::Right:     p += QPoint(0, r.height() / 2);              break;
    }

    return p;
}

void TrashWidget::updateIcon()
{
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    QString iconString = "user-trash";
    if (!m_popupApplet->empty())
        iconString.append("-full");
    if (displayMode == Dock::Efficient)
        iconString.append("-symbolic");

    const int size = displayMode == Dock::Fashion ? std::min(width(), height()) * 0.8 : 16;
    QIcon icon = QIcon::fromTheme(iconString);
    m_icon = icon.pixmap(size, size);

    update();
}

void TrashWidget::showMenu()
{
    DMenu *menu = new DMenu(this);
    menu->setDockMenu(true);

    menu->addAction(&m_openAct);
    if (!m_popupApplet->empty())
        menu->addAction(&m_clearAct);

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    switch (position)
    {
    case Dock::Top:     menu->setDirection(DMenu::Top);         break;
    case Dock::Left:    menu->setDirection(DMenu::Left);        break;
    case Dock::Bottom:  menu->setDirection(DMenu::Bottom);      break;
    case Dock::Right:   menu->setDirection(DMenu::Right);       break;
    default:            Q_UNREACHABLE();
    }

    const QPoint p = popupMarkPoint();

    connect(menu, &DMenu::triggered, this, &TrashWidget::menuTriggered);

    menu->exec(p);

    m_clearAct.setParent(this);
    m_openAct.setParent(this);
    menu->deleteLater();

    emit requestRefershWindowVisible();
}

void TrashWidget::removeApp(const QString &appKey)
{
    const QString cmd("dbus-send --print-reply --dest=com.deepin.dde.Launcher /com/deepin/dde/Launcher com.deepin.dde.Launcher.UninstallApp string:\"" + appKey + "\"");

    QProcess *proc = new QProcess;
    proc->start(cmd);
    proc->waitForFinished();

    proc->deleteLater();
}

void TrashWidget::menuTriggered(DAction *action)
{
    if (action == &m_clearAct)
        m_popupApplet->clearTrashFloder();
    else if (action == &m_openAct)
        m_popupApplet->openTrashFloder();
}

void TrashWidget::moveToTrash(const QUrl &url)
{
    const QFileInfo info = url.toLocalFile();

    QDir trashDir(m_popupApplet->trashDir() + "/files");
    if (!trashDir.exists())
        trashDir.mkpath(".");

    QDir().rename(info.absoluteFilePath(), trashDir.absoluteFilePath(info.fileName()));
}
