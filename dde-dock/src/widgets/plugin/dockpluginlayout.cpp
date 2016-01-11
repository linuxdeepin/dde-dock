#include "dockpluginlayout.h"
#include "../../panel/panelmenu.h"
#include "../../controller/dockmodedata.h"

DockPluginLayout::DockPluginLayout(QWidget *parent) : MovableLayout(parent)
{
    setAcceptDrops(false);
    initPluginManager();
}

QSize DockPluginLayout::sizeHint() const
{
    QSize size;
    int w = 0;
    int h = 0;
    switch (direction()) {
    case QBoxLayout::LeftToRight:
    case QBoxLayout::RightToLeft:
        size.setHeight(DockModeData::instance()->getAppletsItemHeight());
        for (QWidget * widget : widgets()) {
            w += widget->sizeHint().width();
        }
        size.setWidth(w + getLayoutSpacing() * widgets().count());
        break;
    case QBoxLayout::TopToBottom:
    case QBoxLayout::BottomToTop:
        size.setWidth(DockModeData::instance()->getAppletsItemWidth());
        for (QWidget * widget : widgets()) {
            h += widget->height();
        }
        size.setHeight(h + getLayoutSpacing() * widgets().count());
        break;
    }

    return size;
}

void DockPluginLayout::initAllPlugins()
{
    QTimer::singleShot(500, m_pluginManager, SLOT(initAll()));
}

void DockPluginLayout::initPluginManager()
{
    m_pluginManager = new DockPluginsManager(this);

    connect(m_pluginManager, &DockPluginsManager::itemAppend, [=](DockItem *targetItem){
        this->insertWidget(0, targetItem);
        connect(targetItem, &DockItem::needPreviewShow, this, &DockPluginLayout::needPreviewShow);
        connect(targetItem, &DockItem::needPreviewHide, this, &DockPluginLayout::needPreviewHide);
        connect(targetItem, &DockItem::needPreviewUpdate, this, &DockPluginLayout::needPreviewUpdate);
    });
    connect(m_pluginManager, &DockPluginsManager::itemInsert, [=](DockItem *baseItem, DockItem *targetItem){
        int index = indexOf(baseItem);
        insertWidget(index != -1 ? index : count(), targetItem);
        connect(targetItem, &DockItem::needPreviewShow, this, &DockPluginLayout::needPreviewShow);
        connect(targetItem, &DockItem::needPreviewHide, this, &DockPluginLayout::needPreviewHide);
        connect(targetItem, &DockItem::needPreviewUpdate, this, &DockPluginLayout::needPreviewUpdate);
    });
    connect(m_pluginManager, &DockPluginsManager::itemRemoved, [=](DockItem* item) {
        removeWidget(item);
    });
    connect(PanelMenu::instance(), &PanelMenu::settingPlugin, [=]{
        m_pluginManager->onPluginsSetting(getScreenRect().height - parentWidget()->height());
    });
}

DisplayRect DockPluginLayout::getScreenRect()
{
    DBusDisplay d;
    return d.primaryRect();
}

