#include "panelmenu.h"

PanelMenuItem::PanelMenuItem(QString text, QWidget *parent) : QLabel(text,parent)
{
    this->setAlignment(Qt::AlignCenter);
}

void PanelMenuItem::mousePressEvent(QMouseEvent *event)
{
//    emit itemClicked();
}

void PanelMenuItem::mouseReleaseEvent(QMouseEvent *event)
{
    emit itemClicked();
}

PanelMenu::PanelMenu(QWidget *parent) : QWidget(parent)
{
    this->resize(150,100);
    this->setWindowFlags(Qt::ToolTip);

    QLabel * menuContent = new QLabel(this);
    menuContent->setObjectName("panelMenuContent");
    menuContent->resize(this->width(),this->height());
    menuContent->move(0,0);

    PanelMenuItem *fashionItem = new PanelMenuItem("Fashion Mode",this);
    fashionItem->resize(this->width(),MENU_ITEM_HEIGHT);
    fashionItem->move(0,0);
    connect(fashionItem, SIGNAL(itemClicked()),this, SLOT(changeToFashionMode()));

    PanelMenuItem *efficientItem = new PanelMenuItem("Efficient Mode",this);
    efficientItem->resize(this->width(),MENU_ITEM_HEIGHT);
    efficientItem->move(0,MENU_ITEM_HEIGHT + MENU_ITEM_SPACING);
    connect(efficientItem, SIGNAL(itemClicked()),this, SLOT(changeToEfficientMode()));

    PanelMenuItem *classictItem = new PanelMenuItem("Classic Mode",this);
    classictItem->resize(this->width(),MENU_ITEM_HEIGHT);
    classictItem->move(0,MENU_ITEM_HEIGHT*2 + MENU_ITEM_SPACING*2);
    connect(classictItem, SIGNAL(itemClicked()),this, SLOT(changeToClassicMode()));
}

void PanelMenu::changeToFashionMode()
{
    qWarning() << "Change to fashion mode...";
    dockCons->setDockMode(DockConstants::FashionMode);
    this->hide();
}

void PanelMenu::changeToEfficientMode()
{
    qWarning() << "Change to efficient mode...";
    dockCons->setDockMode(DockConstants::EfficientMode);
    this->hide();
}

void PanelMenu::changeToClassicMode()
{
    qWarning() << "Change to classic mode...";
    dockCons->setDockMode(DockConstants::ClassicMode);
    this->hide();
}
