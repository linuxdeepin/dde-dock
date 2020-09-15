#include "dockapplication.h"
#include "constants.h"

#include <QMouseEvent>
#include <QTouchEvent>

DockApplication::DockApplication(int &argc, char **argv) : DApplication (argc, argv)
{
}

bool DockApplication::notify(QObject *obj, QEvent *event)
{
    QMouseEvent *mouseEvent = dynamic_cast<QMouseEvent *>(event);

    if (mouseEvent) {
        // 鼠标事件可以通过source函数确定是否触屏事件，并将结果写入qApp的动态属性中
        const Qt::MouseEventSource src = mouseEvent->source();
        if (src == Qt::MouseEventSynthesizedByQt || src == Qt::MouseEventSynthesizedByApplication) {
            qApp->setProperty(IS_TOUCH_STATE, true);
        } else {
            qApp->setProperty(IS_TOUCH_STATE, false);
        }
    }

    // 任务栏屏蔽多指触控
    QTouchEvent *touchEvent = dynamic_cast<QTouchEvent *>(event);
    if(touchEvent && (touchEvent->touchPoints().size() > 1)) {
        return true;
    }

    return DApplication::notify(obj, event);
}
