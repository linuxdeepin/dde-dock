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
    void detachPluginWidget();

private:
    void mousePressEvent(QMouseEvent *e);
    void mouseMoveEvent(QMouseEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    void paintEvent(QPaintEvent *e);
    bool eventFilter(QObject *o, QEvent *e);

    QWidget *popupTips();

private:
    void startDrag();
    void mouseClicked();

private:
    PluginsItemInterface * const m_pluginInter;
    const QString m_itemKey;

    bool m_draging;

    static QPoint MousePressPoint;
};

#endif // PLUGINSITEM_H
