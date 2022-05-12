#include "quicksettingcontroller.h"
#include "quicksettingitem.h"

QuickSettingController::QuickSettingController(QObject *parent)
    : AbstractPluginsController(parent)
{
    // 异步加载本地插件
    QMetaObject::invokeMethod(this, &QuickSettingController::startLoader, Qt::QueuedConnection);
}

QuickSettingController::~QuickSettingController()
{
}

void QuickSettingController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QList<QuickSettingItem *>::iterator findItemIterator = std::find_if(m_quickSettingItems.begin(), m_quickSettingItems.end(),
                 [ = ](QuickSettingItem *item) {
        return item->itemKey() == itemKey;
    });

    if (findItemIterator != m_quickSettingItems.end())
        return;

    QuickSettingItem *quickItem = new QuickSettingItem(itemInter, itemKey);

    m_quickSettingItems << quickItem;

    emit pluginInsert(quickItem);
}

void QuickSettingController::itemUpdate(PluginsItemInterface * const itemInter, const QString &)
{
    auto findItemIterator = std::find_if(m_quickSettingItems.begin(), m_quickSettingItems.end(),
                                         [ = ](QuickSettingItem *item) {
        return item->pluginItem() == itemInter;
    });
    if (findItemIterator != m_quickSettingItems.end()) {
        QuickSettingItem *settingItem = *findItemIterator;
        settingItem->update();
    }
}

void QuickSettingController::itemRemoved(PluginsItemInterface * const itemInter, const QString &)
{
    // 删除本地记录的插件列表
    QList<QuickSettingItem *>::iterator findItemIterator = std::find_if(m_quickSettingItems.begin(), m_quickSettingItems.end(),
                                         [ = ](QuickSettingItem *item) {
            return (item->pluginItem() == itemInter);
    });
    if (findItemIterator != m_quickSettingItems.end()) {
        QuickSettingItem *quickItem = *findItemIterator;
        m_quickSettingItems.removeOne(quickItem);
        Q_EMIT pluginRemove(quickItem);
        quickItem->deleteLater();
    }
}

QuickSettingController *QuickSettingController::instance()
{
    static QuickSettingController instance;
    return &instance;
}

void QuickSettingController::startLoader()
{
    QString pluginsDir("../plugins/quick-trays");
    if (!QDir(pluginsDir).exists())
        pluginsDir = "/usr/lib/dde-dock/plugins/quick-trays";

    AbstractPluginsController::startLoader(new PluginLoader(pluginsDir, this));
}
