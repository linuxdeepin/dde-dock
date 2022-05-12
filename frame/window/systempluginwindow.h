#ifndef SYSTEMPLUGINWINDOW_H
#define SYSTEMPLUGINWINDOW_H

#include "constants.h"
#include "dockpluginscontroller.h"

#include <DBlurEffectWidget>

class DockPluginsController;
class PluginsItem;
class QBoxLayout;

namespace Dtk { namespace Widget { class DListView; } }

DWIDGET_USE_NAMESPACE

class SystemPluginWindow : public DBlurEffectWidget
{
    Q_OBJECT

Q_SIGNALS:
    void pluginSizeChanged();

public:
    explicit SystemPluginWindow(QWidget *parent = nullptr);
    ~SystemPluginWindow() override;
    void setPositon(Dock::Position position);
    QSize suitableSize();

private:
    void initUi();
    int calcIconSize() const;
    void resizeEvent(QResizeEvent *event) override;

private Q_SLOTS:
    void onPluginItemAdded(PluginsItem *pluginItem);
    void onPluginItemRemoved(PluginsItem *pluginItem);
    void onPluginItemUpdated(PluginsItem *pluginItem);

private:
    DockPluginsController *m_pluginController;
    DListView *m_listView;
    Dock::Position m_position;
    QBoxLayout *m_mainLayout;
};

class FixedPluginController : public DockPluginsController
{
    Q_OBJECT

public:
    FixedPluginController(QObject *parent);

protected:
    const QVariant getValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback = QVariant()) override;
    PluginsItem *createPluginsItem(PluginsItemInterface *const itemInter, const QString &itemKey, const QString &pluginApi) override;
};

#endif // SYSTEMPLUGINWINDOW_H
