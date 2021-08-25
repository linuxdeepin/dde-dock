#include <QTest>

#include <gtest/gtest.h>

#include "testplugin.h"

#define private public
#include "traypluginitem.h"
#undef private

class Ut_TrayPluginItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_TrayPluginItem::SetUp()
{
}

void Ut_TrayPluginItem::TearDown()
{
}

TEST_F(Ut_TrayPluginItem, coverage_test)
{
    TestPlugin plugin;
    TrayPluginItem item(&plugin, "", "");

    ASSERT_EQ(item.itemType(), DockItem::TrayPlugin);

    item.setSuggestIconSize(QSize());
    item.setRightSplitVisible(true);

    ASSERT_EQ(item.trayVisableItemCount(), 0);

    QMouseEvent event(QEvent::MouseButtonPress, QPointF(), Qt::NoButton, Qt::NoButton, Qt::NoModifier);
    qApp->sendEvent(item.centralWidget(), &event);

    QDynamicPropertyChangeEvent event1("TrayVisableItemCount");
    qApp->sendEvent(item.centralWidget(), &event);
}
