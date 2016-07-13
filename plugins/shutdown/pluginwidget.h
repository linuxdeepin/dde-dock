#ifndef PLUGINWIDGET_H
#define PLUGINWIDGET_H

#include <QWidget>

class PluginWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PluginWidget(QWidget *parent = 0);

protected:
    QSize sizeHint() const;
};

#endif // PLUGINWIDGET_H
