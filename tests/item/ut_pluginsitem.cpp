#include <QTest>

#include <gtest/gtest.h>

#include "testplugin.h"

#include "pluginsitem.h"

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
    item.setDraging(true);
    item.refreshIcon();
    item.onGSettingsChanged("");
    item.startDrag();
    item.mouseClicked();

    QWidget widget;
    item.showPopupWindow(&widget);
    ASSERT_FALSE(item.contextMenu().isEmpty());

    ASSERT_TRUE(item.centralWidget());
}

TEST_F(Ut_PluginsItem, event_test)
{
    TestPlugin plugin;
    PluginsItem item(&plugin, "", "");

    QTest::mousePress(&item, Qt::LeftButton, Qt::NoModifier);
    QTest::mousePress(&item, Qt::RightButton, Qt::NoModifier);
    QTest::mouseMove(&item, QPoint());

    QMouseEvent event1(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
    item.mouseMoveEvent(&event1);

    QMouseEvent event2(QEvent::MouseMove, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    item.mouseMoveEvent(&event2);

    QMouseEvent event3(QEvent::MouseButtonRelease, QPointF(0, 0), Qt::LeftButton, Qt::RightButton, Qt::ControlModifier);
    item.mouseReleaseEvent(&event3);

    QMouseEvent event4(QEvent::MouseButtonRelease, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    item.mouseReleaseEvent(&event4);

    QPointF p;
    QEnterEvent event5(p, p, p);
    item.enterEvent(&event5);

    QEvent event6(QEvent::Leave);
    item.leaveEvent(&event6);

    QShowEvent event7;
    item.showEvent(&event7);

    ASSERT_TRUE(true);
}
