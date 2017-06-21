#ifndef PLUGINSITEM_H
#define PLUGINSITEM_H

#include "dockitem.h"
#include "pluginsiteminterface.h"

class PluginsItem : public DockItem
{
    Q_OBJECT

public:
    explicit PluginsItem(PluginsItemInterface* const pluginInter, const QString &itemKey, QWidget *parent = 0);
    ~PluginsItem();

    int itemSortKey() const;
    void setItemSortKey(const int order) const;
    void detachPluginWidget();

    bool allowContainer() const;
    bool isInContainer() const;
    void setInContainer(const bool container);

    using DockItem::showContextMenu;
    using DockItem::hidePopup;

    inline ItemType itemType() const {return Plugins;}

public slots:
    void refershIcon();

private:
    void mousePressEvent(QMouseEvent *e);
    void mouseMoveEvent(QMouseEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    bool eventFilter(QObject *o, QEvent *e);

    void invokedMenuItem(const QString &itemId, const bool checked);
    const QString contextMenu() const;
    QWidget *popupTips();

private:
    void startDrag();
    void mouseClicked();

private:
    PluginsItemInterface * const m_pluginInter;
    QWidget *m_centralWidget;
    const QString m_itemKey;
    bool m_draging;

    static QPoint MousePressPoint;
};

#endif // PLUGINSITEM_H
