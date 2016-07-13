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
    : DockItem(Plugins, parent),
      m_pluginInter(pluginInter),
      m_centeralWidget(m_pluginInter->itemWidget(itemKey)),
      m_itemKey(itemKey),
      m_draging(false)
{
    Q_ASSERT(m_centeralWidget);

    m_centeralWidget->installEventFilter(this);

    QBoxLayout *hLayout = new QHBoxLayout;
    hLayout->addWidget(m_centeralWidget);
    hLayout->setSpacing(0);
    hLayout->setMargin(0);

    setLayout(hLayout);
    setAttribute(Qt::WA_TranslucentBackground);
}

PluginsItem::~PluginsItem()
{
}

int PluginsItem::itemSortKey() const
{
    return m_pluginInter->itemSortKey(m_itemKey);
}

void PluginsItem::detachPluginWidget()
{
    QWidget *widget = m_pluginInter->itemWidget(m_itemKey);
    if (widget)
        widget->setParent(nullptr);
}

void PluginsItem::mousePressEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton)
    {
        MousePressPoint = e->pos();
    }
    else if (e->button() == Qt::RightButton)
    {

    }
    else
        DockItem::mousePressEvent(e);
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
    if (e->button() != Qt::LeftButton)
        return DockItem::mouseReleaseEvent(e);

    e->accept();

    const QPoint distance = e->pos() - MousePressPoint;
    if (distance.manhattanLength() < PLUGIN_ITEM_DRAG_THRESHOLD)
        mouseClicked();
}

void PluginsItem::paintEvent(QPaintEvent *e)
{
    if (m_draging)
        return;

    DockItem::paintEvent(e);

    // TODO: hover effect
}

bool PluginsItem::eventFilter(QObject *o, QEvent *e)
{
    if (m_draging)
        if (o == m_centeralWidget && e->type() == QEvent::Paint)
            return true;

    return DockItem::eventFilter(o, e);
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

    QDrag *drag = new QDrag(this);
    drag->setPixmap(pixmap);
    drag->setHotSpot(pixmap.rect().center());
    drag->setMimeData(new QMimeData);

    emit dragStarted();
    const Qt::DropAction result = drag->exec(Qt::MoveAction);
    Q_UNUSED(result);

    m_draging = false;
    setVisible(true);
    update();
}

void PluginsItem::mouseClicked()
{
    const QString command = m_pluginInter->itemCommand(m_itemKey);
    if (command.isEmpty())
        return;

    QProcess *proc = new QProcess(this);

    connect(proc, static_cast<void (QProcess::*)(int)>(&QProcess::finished), proc, &QProcess::deleteLater);

    proc->startDetached(command);
}
