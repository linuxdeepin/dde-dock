#ifndef TESTPLUGIN_H
#define TESTPLUGIN_H

#include "pluginsiteminterface.h"

const QString Name = "Test";

class TestPlugin : public PluginsItemInterface
{
public:
    TestPlugin();

    virtual const QString pluginName() const override;
    virtual const QString pluginDisplayName() const override;
    virtual void init(PluginProxyInterface *proxyInter) override;
    virtual QWidget *itemWidget(const QString &itemKey) override;
    virtual int itemSortKey(const QString &itemKey) override;
    virtual void setSortKey(const QString &itemKey, const int order) override;
    virtual PluginSizePolicy pluginSizePolicy() const override;
    virtual PluginType type() override;

public:
    void setType(const PluginType type);

private:
    int m_sortKey;
    PluginType m_type;
};

#endif // TESTPLUGIN_H
