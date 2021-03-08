#include "themeappicon.h"

#include <QPixmap>

#include <gtest/gtest.h>

class Ut_ThemeAppIcon : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_ThemeAppIcon::SetUp()
{
}

void Ut_ThemeAppIcon::TearDown()
{
}

TEST_F(Ut_ThemeAppIcon, getIcon_test)
{
    ThemeAppIcon appIcon;
    const QPixmap &pix1 = appIcon.getIcon("", 50, 1.0);
    ASSERT_FALSE(pix1.isNull());
    appIcon.getIcon("dde-calendar", 50, 1.0);
    const QPixmap &pix2 = appIcon.getIcon("data:image/test", 50, 1.0);
    ASSERT_FALSE(pix2.isNull());
    const QPixmap &pix3 = appIcon.getIcon(":/res/all_settings_on.png", 50, 1.0);
    ASSERT_FALSE(pix3.isNull());
}
