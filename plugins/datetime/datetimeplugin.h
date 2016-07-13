#ifndef DATETIMEPLUGIN_H
#define DATETIMEPLUGIN_H

#include "pluginsiteminterface.h"
#include "datetimewidget.h"

#include <QTimer>
#include <QLabel>

class DatetimePlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "datetime.json")

public:
    explicit DatetimePlugin(QObject *parent = 0);
    ~DatetimePlugin();

    const QString pluginName() const;
    ItemType pluginType(const QString &itemKey);
    ItemType tipsType(const QString &itemKey);
    void init(PluginProxyInterface *proxyInter);

    int itemSortKey(const QString &itemKey) const;

    QWidget *itemWidget(const QString &itemKey);
    QWidget *itemTipsWidget(const QString &itemKey);

    const QString itemCommand(const QString &itemKey);

private slots:
    void updateCurrentTimeString();

private:
    DatetimeWidget *m_centeralWidget;
    QLabel *m_dateTipsLabel;

    QTimer *m_refershTimer;

    QString m_currentTimeString;
};

#endif // DATETIMEPLUGIN_H
