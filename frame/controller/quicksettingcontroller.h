#ifndef QUICKSETTINGCONTROLLER_H
#define QUICKSETTINGCONTROLLER_H

#include "abstractpluginscontroller.h"

class QuickSettingItem;

class QuickSettingController : public AbstractPluginsController
{
    Q_OBJECT

public:
    static QuickSettingController *instance();
    const QList<QuickSettingItem *> &settingItems() const { return m_quickSettingItems; }

Q_SIGNALS:
    void pluginInsert(QuickSettingItem *);
    void pluginRemove(QuickSettingItem *);

protected:
    void startLoader();
    QuickSettingController(QObject *parent = Q_NULLPTR);
    ~QuickSettingController() override;

protected:
    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemUpdate(PluginsItemInterface * const itemInter, const QString &) override;
    void itemRemoved(PluginsItemInterface * const itemInter, const QString &) override;
    void requestWindowAutoHide(PluginsItemInterface * const, const QString &, const bool) override {}
    void requestRefreshWindowVisible(PluginsItemInterface * const, const QString &) override {}
    void requestSetAppletVisible(PluginsItemInterface * const, const QString &, const bool) override {}

private:
    QList<QuickSettingItem *> m_quickSettingItems;
};

#endif // CONTAINERPLUGINSCONTROLLER_H
