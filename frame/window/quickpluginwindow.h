#ifndef QUICKPLUGINWINDOW_H
#define QUICKPLUGINWINDOW_H

#include "constants.h"

#include <QWidget>

class QuickSettingItem;
class PluginsItemInterface;
class QHBoxLayout;
class QuickSettingContainer;
class QStandardItemModel;
class QStandardItem;
class QMouseEvent;

namespace Dtk { namespace Gui { class DRegionMonitor; }
                namespace Widget { class DListView; class DStandardItem; } }

using namespace Dtk::Widget;

class QuickPluginWindow : public QWidget
{
    Q_OBJECT

Q_SIGNALS:
    void itemCountChanged();

public:
    explicit QuickPluginWindow(QWidget *parent = nullptr);
    ~QuickPluginWindow() override;

    void setPositon(Dock::Position position);
    void addPlugin(QuickSettingItem *item);

    QSize suitableSize();

protected:
    void mouseReleaseEvent(QMouseEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;

private:
    void initUi();
    void initConnection();
    void resetSortRole();
    int fixedItemCount();
    DStandardItem *createStandItem(QuickSettingItem *item);
    void removePlugin(QuickSettingItem *item);
    void startDrag(QuickSettingItem *moveItem);

private:
    DListView *m_listView;
    QStandardItemModel *m_model;
    Dock::Position m_position;
};

#endif // QUICKPLUGINWINDOW_H
