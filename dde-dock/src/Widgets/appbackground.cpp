#include "appbackground.h"

AppBackground::AppBackground(QWidget *parent) :
    QLabel(parent)
{
    this->setStyleSheet("QLabel#AppBackground{background: rgba(255,255,255,0.3);border-radius: 4px;}");
}
