#include "appbackground.h"

AppBackground::AppBackground(QWidget *parent) :
    QLabel(parent)
{
    this->setStyleSheet("background:#121922;border-radius: 4px;");
}
