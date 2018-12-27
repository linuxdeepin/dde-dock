#include "normalcontainer.h"

NormalContainer::NormalContainer(TrayPlugin *trayPlugin, QWidget *parent)
    : AbstractContainer(trayPlugin, parent)
{

}

bool NormalContainer::acceptWrapper(FashionTrayWidgetWrapper *wrapper)
{
    Q_UNUSED(wrapper);

    return true;
}

void NormalContainer::refreshVisible()
{
    setVisible(expand() && !isEmpty());
}

void NormalContainer::setExpand(const bool expand)
{
    for (auto w : wrapperList()) {
        // reset all tray item attention state
        w->setAttention(false);
    }

    AbstractContainer::setExpand(expand);
}

int NormalContainer::whereToInsert(FashionTrayWidgetWrapper *wrapper)
{
    // 如果已经对图标进行过排序则完全按照从配置文件中获取的顺序来插入图标(即父类的实现)
    if (trayPlugin()->traysSortedInFashionMode()) {
        return AbstractContainer::whereToInsert(wrapper);
    }

    // 如果没有对图标进行过排序则使用下面的默认排序算法:
    // 所有应用图标在系统图标的左侧
    // 新的应用图标在最左侧的应用图标处插入
    // 新的系统图标在最左侧的系统图标处插入
    return whereToInsertByDefault(wrapper);
}

int NormalContainer::whereToInsertByDefault(FashionTrayWidgetWrapper *wrapper) const
{
    int index = 0;
    switch (wrapper->absTrayWidget()->trayTyep()) {
    case AbstractTrayWidget::TrayType::ApplicationTray:
        index = whereToInsertAppTrayByDefault(wrapper);
        break;
    case AbstractTrayWidget::TrayType::SystemTray:
        index = whereToInsertSystemTrayByDefault(wrapper);
        break;
    default:
        Q_UNREACHABLE();
        break;
    }
    return index;
}

int NormalContainer::whereToInsertAppTrayByDefault(FashionTrayWidgetWrapper *wrapper) const
{
    if (wrapperList().isEmpty() || wrapper->absTrayWidget()->trayTyep() != AbstractTrayWidget::TrayType::ApplicationTray) {
        return 0;
    }

    int lastAppTrayIndex = -1;
    for (int i = 0; i < wrapperList().size(); ++i) {
        if (wrapperList().at(i)->absTrayWidget()->trayTyep() == AbstractTrayWidget::TrayType::ApplicationTray) {
            lastAppTrayIndex = i;
            continue;
        }
        break;
    }
    // there is no AppTray
    if (lastAppTrayIndex == -1) {
        return 0;
    }
    // the inserting tray is not a AppTray
    if (wrapper->absTrayWidget()->trayTyep() != AbstractTrayWidget::TrayType::ApplicationTray) {
        return lastAppTrayIndex + 1;
    }

    int insertIndex = trayPlugin()->itemSortKey(wrapper->itemKey());
    // invalid index
    if (insertIndex < -1) {
        return 0;
    }
    for (int i = 0; i < wrapperList().size(); ++i) {
        if (wrapperList().at(i)->absTrayWidget()->trayTyep() != AbstractTrayWidget::TrayType::ApplicationTray) {
            insertIndex = i;
            break;
        }
        if (insertIndex > trayPlugin()->itemSortKey(wrapperList().at(i)->itemKey())) {
            continue;
        }
        insertIndex = i;
        break;
    }
    if (insertIndex > lastAppTrayIndex + 1) {
        insertIndex = lastAppTrayIndex + 1;
    }

    return insertIndex;
}

int NormalContainer::whereToInsertSystemTrayByDefault(FashionTrayWidgetWrapper *wrapper) const
{
    if (wrapperList().isEmpty()) {
        return 0;
    }

    int firstSystemTrayIndex = -1;
    for (int i = 0; i < wrapperList().size(); ++i) {
        if (wrapperList().at(i)->absTrayWidget()->trayTyep() == AbstractTrayWidget::TrayType::SystemTray) {
            firstSystemTrayIndex = i;
            break;
        }
    }
    // there is no SystemTray
    if (firstSystemTrayIndex == -1) {
        return wrapperList().size();
    }
    // the inserting tray is not a SystemTray
    if (wrapper->absTrayWidget()->trayTyep() != AbstractTrayWidget::TrayType::SystemTray) {
        return firstSystemTrayIndex;
    }

    int insertIndex = trayPlugin()->itemSortKey(wrapper->itemKey());
    // invalid index
    if (insertIndex < -1) {
        return firstSystemTrayIndex;
    }
    for (int i = 0; i < wrapperList().size(); ++i) {
        if (wrapperList().at(i)->absTrayWidget()->trayTyep() != AbstractTrayWidget::TrayType::SystemTray) {
            continue;
        }
        if (insertIndex > trayPlugin()->itemSortKey(wrapperList().at(i)->itemKey())) {
            continue;
        }
        insertIndex = i;
        break;
    }
    if (insertIndex < firstSystemTrayIndex) {
        return firstSystemTrayIndex;
    }

    return insertIndex;
}
