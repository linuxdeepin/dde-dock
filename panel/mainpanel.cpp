#include "mainpanel.h"

#include <QHBoxLayout>

MainPanel::MainPanel(QWidget *parent)
    : QFrame(parent),
      m_itemController(DockItemController::instance(this))
{
    setObjectName("MainPanel");
    setStyleSheet("QWidget #MainPanel {"
                  "border:none;"
                  "background-color:red;"
//                  "border-radius:5px 5px 5px 5px;"
                  "}");

    QHBoxLayout *layout = new QHBoxLayout;
    setLayout(layout);

    const QList<DockItem *> itemList = m_itemController->itemList();
    for (auto item : itemList)
        layout->addWidget(item);
}
