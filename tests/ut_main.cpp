#include <gtest/gtest.h>
#include <QApplication>
#include <QDebug>
#include <DLog>

int main(int argc, char **argv)
{
    // gerrit编译时没有显示器，需要指定环境变量
    qputenv("QT_QPA_PLATFORM", "offscreen");

    QApplication app(argc, argv);

    qApp->setProperty("CANSHOW", true);

    ::testing::InitGoogleTest(&argc, argv);

    return RUN_ALL_TESTS();
}
