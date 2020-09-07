#include "dockapplication.h"
#include "constants.h"

#include <QMouseEvent>

DockApplication::DockApplication(int &argc, char **argv) : DApplication (argc, argv)
{
}

bool DockApplication::notify(QObject *obj, QEvent *event)
{
    QMouseEvent *mouseEvent = dynamic_cast<QMouseEvent *>(event);
    if (mouseEvent) {
        // 鼠标事件可以通过source函数确定是否触屏事件，并将结果写入qApp的动态属性中
        qApp->setProperty(IS_TOUCH_STATE, (mouseEvent->source() == Qt::MouseEventSynthesizedByQt));
    }

    return DApplication::notify(obj, event);
}
