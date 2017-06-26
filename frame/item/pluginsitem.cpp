#include "constants.h"
#include "pluginsitem.h"
#include "pluginsiteminterface.h"

#include "util/imagefactory.h"

#include <QPainter>
#include <QBoxLayout>
#include <QMouseEvent>
#include <QDrag>
#include <QMimeData>

#define PLUGIN_ITEM_DRAG_THRESHOLD      20

QPoint PluginsItem::MousePressPoint = QPoint();

PluginsItem::PluginsItem(PluginsItemInterface* const pluginInter, const QString &itemKey, QWidget *parent)
    : DockItem(parent),
      m_pluginInter(pluginInter),
      m_centralWidget(m_pluginInter->itemWidget(itemKey)),
      m_itemKey(itemKey),
      m_draging(false)
{
    qDebug() << "load plugins item: " << pluginInter->pluginName() << itemKey << m_centralWidget;

    m_centralWidget->setParent(this);
    m_centralWidget->setVisible(true);
    m_centralWidget->installEventFilter(this);

    QBoxLayout *hLayout = new QHBoxLayout;
    hLayout->addWidget(m_centralWidget);
    hLayout->setSpacing(0);
    hLayout->setMargin(0);

    setLayout(hLayout);
    setAccessibleName(pluginInter->pluginName() + "-" + m_itemKey);
    setAttribute(Qt::WA_TranslucentBackground);
}

PluginsItem::~PluginsItem()
{
}

int PluginsItem::itemSortKey() const
{
    return m_pluginInter->itemSortKey(m_itemKey);
}

void PluginsItem::setItemSortKey(const int order) const
{
    m_pluginInter->setSortKey(m_itemKey, order);
}

void PluginsItem::detachPluginWidget()
{
    QWidget *widget = m_pluginInter->itemWidget(m_itemKey);
    if (widget)
        widget->setParent(nullptr);
}

bool PluginsItem::allowContainer() const
{
    if (DockDisplayMode == Dock::Fashion)
        return false;

    return m_pluginInter->itemAllowContainer(m_itemKey);
}

bool PluginsItem::isInContainer() const
{
    if (DockDisplayMode == Dock::Fashion)
        return false;

    return m_pluginInter->itemIsInContainer(m_itemKey);
}

void PluginsItem::setInContainer(const bool container)
{
    m_pluginInter->setItemIsInContainer(m_itemKey, container);
}

QSize PluginsItem::sizeHint() const
{
    return m_centralWidget->sizeHint();
}

void PluginsItem::refershIcon()
{
    m_pluginInter->refershIcon(m_itemKey);
}

void PluginsItem::mousePressEvent(QMouseEvent *e)
{
    QWidget::mousePressEvent(e);

    if (e->button() == Qt::LeftButton)
        MousePressPoint = e->pos();
}

void PluginsItem::mouseMoveEvent(QMouseEvent *e)
{
    if (e->buttons() != Qt::LeftButton)
        return DockItem::mouseMoveEvent(e);

    e->accept();

    const QPoint distance = e->pos() - MousePressPoint;
    if (distance.manhattanLength() > PLUGIN_ITEM_DRAG_THRESHOLD)
        startDrag();
}

void PluginsItem::mouseReleaseEvent(QMouseEvent *e)
{
    DockItem::mouseReleaseEvent(e);

    if (e->button() != Qt::LeftButton)
        return;

    e->accept();

    const QPoint distance = e->pos() - MousePressPoint;
    if (distance.manhattanLength() < PLUGIN_ITEM_DRAG_THRESHOLD)
        mouseClicked();
}

bool PluginsItem::eventFilter(QObject *o, QEvent *e)
{
    if (m_draging)
        if (o == m_centralWidget && e->type() == QEvent::Paint)
            return true;

    return DockItem::eventFilter(o, e);
}

void PluginsItem::invokedMenuItem(const QString &itemId, const bool checked)
{
    m_pluginInter->invokedMenuItem(m_itemKey, itemId, checked);
}

const QString PluginsItem::contextMenu() const
{
    return m_pluginInter->itemContextMenu(m_itemKey);
}

QWidget *PluginsItem::popupTips()
{
    return m_pluginInter->itemTipsWidget(m_itemKey);
}

void PluginsItem::startDrag()
{
    const QPixmap pixmap = grab();

    m_draging = true;
    update();

    QMimeData *mime = new QMimeData;
    mime->setData(DOCK_PLUGIN_MIME, m_itemKey.toStdString().c_str());

    QDrag *drag = new QDrag(this);
    drag->setPixmap(pixmap);
    drag->setHotSpot(pixmap.rect().center());
    drag->setMimeData(mime);

    emit dragStarted();
    const Qt::DropAction result = drag->exec(Qt::MoveAction);
    Q_UNUSED(result);
    emit itemDropped(drag->target());

    m_draging = false;
    setVisible(true);
    update();
}

void PluginsItem::mouseClicked()
{
    const QString command = m_pluginInter->itemCommand(m_itemKey);
    if (!command.isEmpty())
    {
        QProcess *proc = new QProcess(this);

        connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

        proc->startDetached(command);
        return;
    }

    // request popup applet
    QWidget *w = m_pluginInter->itemPopupApplet(m_itemKey);
    if (w)
        showPopupApplet(w);
}
