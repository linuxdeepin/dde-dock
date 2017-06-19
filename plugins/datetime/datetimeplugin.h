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

    const QString pluginName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    int itemSortKey(const QString &itemKey) override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;

    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;

    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;

private slots:
    void updateCurrentTimeString();

private:
    DatetimeWidget *m_centralWidget;
    QLabel *m_dateTipsLabel;

    QTimer *m_refershTimer;

    QString m_currentTimeString;
};

#endif // DATETIMEPLUGIN_H
