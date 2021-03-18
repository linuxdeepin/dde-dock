#include <QTest>

#include <gtest/gtest.h>

#include "testplugin.h"

#define private public
#include "pluginsitem.h"
#undef private

using namespace ::testing;

class Ut_PluginsItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_PluginsItem::SetUp()
{
}

void Ut_PluginsItem::TearDown()
{
}

TEST_F(Ut_PluginsItem, itemSortKey_test)
{
    TestPlugin plugin;
    PluginsItem item(&plugin, "", "");

    ASSERT_EQ(item.itemSortKey(), 0);

    item.setItemSortKey(99);
    ASSERT_EQ(item.itemSortKey(), 99);

    item.detachPluginWidget();

    ASSERT_EQ(item.pluginName(), Name);
}

TEST_F(Ut_PluginsItem, pluginSizePolicy_test)
{
    TestPlugin plugin;
    PluginsItem item(&plugin, "", "1.2.0");

    ASSERT_EQ(item.pluginSizePolicy(), PluginsItemInterface::System);

    PluginsItem item1(&plugin, "", "1.3.0");
    ASSERT_EQ(item1.pluginSizePolicy(), PluginsItemInterface::Custom);
}

TEST_F(Ut_PluginsItem, itemType_test)
{
    TestPlugin plugin;
    PluginsItem item(&plugin, "", "");

    plugin.setType(PluginsItemInterface::Normal);
    ASSERT_EQ(item.itemType(), PluginsItem::Plugins);

    plugin.setType(PluginsItemInterface::Fixed);
    ASSERT_EQ(item.itemType(), PluginsItem::FixedPlugin);
}

TEST_F(Ut_PluginsItem, cover)
{
    TestPlugin plugin;
    PluginsItem item(&plugin, "", "");

    item.sizeHint();

    ASSERT_TRUE(item.centralWidget());

    item.setDraging(true);
    item.refreshIcon();
    item.onGSettingsChanged("");
    item.startDrag();
    item.mouseClicked();
}

TEST_F(Ut_PluginsItem, event_test)
{
    TestPlugin plugin;
    PluginsItem item(&plugin, "", "");

    QTest::mousePress(&item, Qt::LeftButton, Qt::NoModifier);
    QTest::mousePress(&item, Qt::RightButton, Qt::NoModifier);
    QTest::mouseMove(&item, QPoint());
}
