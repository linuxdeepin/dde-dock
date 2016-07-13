#include "pluginwidget.h"

PluginWidget::PluginWidget(QWidget *parent)
    : QWidget(parent)
{

}

QSize PluginWidget::sizeHint() const
{
    return QSize(24, 24);
}
