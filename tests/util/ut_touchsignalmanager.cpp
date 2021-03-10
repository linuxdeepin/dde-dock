#define private public
#include "touchsignalmanager.h"
#undef private

#include <QTest>
#include <QSignalSpy>

#include <gtest/gtest.h>

class Ut_TouchSignalManager : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_TouchSignalManager::SetUp()
{
}

void Ut_TouchSignalManager::TearDown()
{
}

TEST_F(Ut_TouchSignalManager, isDragIconPress_test)
{
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), false);
    TouchSignalManager::instance()->m_dragIconPressed = true;
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), true);
    TouchSignalManager::instance()->m_dragIconPressed = false;
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), false);
}

TEST_F(Ut_TouchSignalManager, dealShortTouchPress_test)
{
    QSignalSpy spy(TouchSignalManager::instance(), SIGNAL(shortTouchPress(int, double, double)));
    TouchSignalManager::instance()->dealShortTouchPress(1, 0, 0);
    QCOMPARE(spy.count(), 1);
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), true);

    const QList<QVariant> &arguments = spy.takeFirst();
    QCOMPARE(arguments.size(), 3);
    QCOMPARE(arguments.at(0), 1);
    ASSERT_TRUE(qAbs(arguments.at(1).toDouble()) < 0.00001);
    ASSERT_TRUE(qAbs(arguments.at(2).toDouble()) < 0.00001);
}

TEST_F(Ut_TouchSignalManager, dealTouchRelease_test)
{
    QSignalSpy spy(TouchSignalManager::instance(), SIGNAL(touchRelease(double, double)));
    TouchSignalManager::instance()->dealTouchRelease(0, 0);
    QCOMPARE(spy.count(), 1);
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), false);

    const QList<QVariant> &arguments = spy.takeFirst();
    QCOMPARE(arguments.size(), 2);
    ASSERT_TRUE(qAbs(arguments.at(0).toDouble()) < 0.00001);
    ASSERT_TRUE(qAbs(arguments.at(1).toDouble()) < 0.00001);
}

TEST_F(Ut_TouchSignalManager, dealTouchPress_test)
{
    QSignalSpy spy(TouchSignalManager::instance(), SIGNAL(middleTouchPress(double, double)));
    TouchSignalManager::instance()->dealTouchPress(1, 1000, 0, 0);
    QCOMPARE(spy.count(), 1);
    const QList<QVariant> &arguments = spy.takeFirst();
    QCOMPARE(arguments.size(), 2);
    ASSERT_TRUE(qAbs(arguments.at(0).toDouble()) < 0.00001);
    ASSERT_TRUE(qAbs(arguments.at(1).toDouble()) < 0.00001);

    TouchSignalManager::instance()->dealTouchPress(1, 2000, 0, 0);
    QCOMPARE(spy.count(), 0);
    TouchSignalManager::instance()->dealTouchPress(1, 500, 0, 0);
    QCOMPARE(spy.count(), 0);
    TouchSignalManager::instance()->dealTouchPress(2, 0000, 0, 0);
    QCOMPARE(spy.count(), 0);
}
