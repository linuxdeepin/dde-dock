#include "dockapplication.h"
#include "constants.h"

#include <QObject>
#include <QTouchEvent>
#include <QTest>
#include <QDebug>

#include <gtest/gtest.h>

class Test_DockApplication : public QObject, public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

protected:
    bool eventFilter(QObject* obj, QEvent *event) override
    {
        if (qApp != obj)
            return false;

        if (event->type() == QEvent::TouchUpdate) {
            QTouchEvent *touchEvent = static_cast<QTouchEvent *>(event);
            if (touchEvent)
                m_touchPointNum = touchEvent->touchPoints().size();
        }
        return false;
    }

public:
    int m_touchPointNum; // 触摸点数
};

void Test_DockApplication::SetUp()
{
    qApp->installEventFilter(this);
    m_touchPointNum = 0;
}

void Test_DockApplication::TearDown()
{
    qApp->removeEventFilter(this);
}

// 检测鼠标事件类型
TEST_F(Test_DockApplication, dockapplication_touchstate_test)
{
    // Qt::MouseEventSynthesizedByQt
    QMouseEvent mouseEvent1(QMouseEvent::MouseButtonPress, QPoint(), QPoint(), QPoint(), Qt::LeftButton, Qt::LeftButton, Qt::NoModifier, Qt::MouseEventSynthesizedByQt);
    qApp->sendEvent(qApp, &mouseEvent1);

    ASSERT_TRUE(qApp->property(IS_TOUCH_STATE).toBool());

    // Qt::MouseEventSynthesizedByApplication
    QMouseEvent mouseEvent2(QMouseEvent::MouseButtonPress, QPoint(), QPoint(), QPoint(), Qt::LeftButton, Qt::LeftButton, Qt::NoModifier, Qt::MouseEventSynthesizedByApplication);
    qApp->sendEvent(qApp, &mouseEvent2);

    ASSERT_TRUE(qApp->property(IS_TOUCH_STATE).toBool());

    // 不指定
    QMouseEvent mouseEvent3(QMouseEvent::MouseButtonPress, QPoint(), QPoint(), QPoint(), Qt::LeftButton, Qt::LeftButton, Qt::NoModifier);
    qApp->sendEvent(qApp, &mouseEvent3);

    ASSERT_FALSE(qApp->property(IS_TOUCH_STATE).toBool());
}

// 检测触摸点数，如果是多点，就截获不到,如果是单点，可以正常截获
TEST_F(Test_DockApplication, dockapplication_touchpoints_test)
{
    // 三点触摸
    QList<QTouchEvent::TouchPoint> list;
    list << QTouchEvent::TouchPoint(0) << QTouchEvent::TouchPoint(1) << QTouchEvent::TouchPoint(2);
    QTouchEvent threePointsTouchEvent(QEvent::TouchUpdate, nullptr, Qt::NoModifier, Qt::TouchPointPressed, list);
    QApplication::sendEvent(qApp, &threePointsTouchEvent);
    QTest::qWait(10);

    EXPECT_EQ(m_touchPointNum, 0);

    // 单点触摸
    list.clear();
    list << QTouchEvent::TouchPoint(0);
    QTouchEvent onePointTouchEvent(QEvent::TouchUpdate, nullptr, Qt::NoModifier, Qt::TouchPointPressed, list);
    QApplication::sendEvent(qApp, &onePointTouchEvent);
    QTest::qWait(10);

    EXPECT_EQ(m_touchPointNum, 1);
}



