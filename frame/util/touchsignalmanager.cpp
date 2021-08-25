#include "touchsignalmanager.h"

// 定义手指数
const int SingleFinger = 1;

// 定义长按时间
//const int DragIconTime = 200;
const int DragDockSizeTime = 1000;

TouchSignalManager *TouchSignalManager::m_touchManager = nullptr;

TouchSignalManager::TouchSignalManager(QObject *parent)
    : QObject(parent)
    , m_gestureInter(new Gesture("com.deepin.daemon.Gesture", "/com/deepin/daemon/Gesture", QDBusConnection::systemBus(), this))
    , m_dragIconPressed(false)
{
    // 处理后端触屏信号
    connect(m_gestureInter, &Gesture::TouchSinglePressTimeout, this, &TouchSignalManager::dealShortTouchPress);
    connect(m_gestureInter, &Gesture::TouchUpOrCancel, this, &TouchSignalManager::dealTouchRelease);

    connect(m_gestureInter, &Gesture::TouchPressTimeout, this, &TouchSignalManager::dealTouchPress);
    connect(m_gestureInter, &Gesture::TouchMoving, this, &TouchSignalManager::touchMove);
}

TouchSignalManager *TouchSignalManager::instance()
{
    if (!m_touchManager) {
        m_touchManager = new TouchSignalManager;
    }
    return m_touchManager;
}

bool TouchSignalManager::isDragIconPress() const
{
    return m_dragIconPressed;
}

void TouchSignalManager::dealShortTouchPress(int time, double scaleX, double scaleY)
{
    m_dragIconPressed = true;
    emit shortTouchPress(time, scaleX, scaleY);
}

void TouchSignalManager::dealTouchRelease(double scaleX, double scaleY)
{
    m_dragIconPressed = false;
    emit touchRelease(scaleX, scaleY);
}

void TouchSignalManager::dealMiddleTouchPress(double scaleX, double scaleY)
{
    emit middleTouchPress(scaleX, scaleY);
}

void TouchSignalManager::dealTouchPress(int figerNum, int time, double scaleX, double scaleY)
{
    if (figerNum == SingleFinger && time == DragDockSizeTime) {
        dealMiddleTouchPress(scaleX, scaleY);
    }
}
