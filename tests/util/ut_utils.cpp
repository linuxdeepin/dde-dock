#include "utils.h"

#include <gtest/gtest.h>

class Ut_Utils : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_Utils::SetUp()
{
}

void Ut_Utils::TearDown()
{
}

TEST_F(Ut_Utils, comparePluginApi_test)
{
    QString v1("1.0.0");
    QString v2("1.0.1");
    QString v3("1.0.0.0");

    ASSERT_EQ(Utils::comparePluginApi(v1, v1), 0);
    ASSERT_EQ(Utils::comparePluginApi(v1, v2), -1);
    ASSERT_EQ(Utils::comparePluginApi(v2, v1), 1);
    ASSERT_EQ(Utils::comparePluginApi(v1, v3), -1);
    ASSERT_EQ(Utils::comparePluginApi(v3, v1), 1);
}

TEST_F(Ut_Utils, isSettingConfigured_test)
{
//    Utils::isSettingConfigured("com.deepin.dde.dock.mainwindow", "/com/deepin/dde/dock/mainwindow/", "only-show-primary");
    ASSERT_FALSE(Utils::isSettingConfigured("", "", ""));
}

TEST_F(Ut_Utils, screenAt_test)
{
    Utils::screenAt(QPoint(0, 0));
    ASSERT_EQ(Utils::screenAt(QPoint(-1, -1)), nullptr);
}

TEST_F(Ut_Utils, screenAtByScaled_test)
{
    Utils::screenAtByScaled(QPoint(0, 0));
    ASSERT_EQ(Utils::screenAtByScaled(QPoint(-1, -1)), nullptr);
}

TEST_F(Ut_Utils, renderSVG_test)
{
    QPixmap pix(":/res/all_settings_on.png");
    const QSize &size = pix.size();

    ASSERT_TRUE(Utils::renderSVG("", size, 1.0).isNull());
    ASSERT_EQ(Utils::renderSVG(":/res/all_settings_on.png", size, 1.0).size(), size);
    ASSERT_EQ(Utils::renderSVG(":/res/all_settings_on.png", QSize(50, 50), 1.0).size(), QSize(50, 50));
    ASSERT_EQ(Utils::renderSVG(":/res/all_settings_on.png", QSize(50, 50), 0.5).size(), QSize(25, 25));
}
