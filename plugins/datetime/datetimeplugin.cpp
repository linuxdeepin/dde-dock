#include "datetimeplugin.h"

#include <QLabel>

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent),

      m_dateTipsLabel(new QLabel),

      m_refershTimer(new QTimer(this))
{
    m_dateTipsLabel->setObjectName("datetime");
    m_dateTipsLabel->setStyleSheet("color:white;"
                                   "padding:0px 3px;");

    m_refershTimer->setInterval(1000);
    m_refershTimer->start();

    m_centralWidget = new DatetimeWidget;

    connect(m_centralWidget, &DatetimeWidget::requestContextMenu, [this] { m_proxyInter->requestContextMenu(this, QString()); });
    connect(m_centralWidget, &DatetimeWidget::requestUpdateGeometry, [this] { m_proxyInter->itemUpdate(this, QString()); });

    connect(m_refershTimer, &QTimer::timeout, this, &DatetimePlugin::updateCurrentTimeString);
}

DatetimePlugin::~DatetimePlugin()
{
    delete m_centralWidget;
    delete m_dateTipsLabel;
}

const QString DatetimePlugin::pluginName() const
{
    return "datetime";
}

void DatetimePlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    m_proxyInter->itemAdded(this, QString());
}

int DatetimePlugin::itemSortKey(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return -1;
}

QWidget *DatetimePlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_centralWidget;
}

QWidget *DatetimePlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_dateTipsLabel;
}

const QString DatetimePlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return "dde-calendar";
}

const QString DatetimePlugin::itemContextMenu(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    QList<QVariant> items;
    items.reserve(1);

    QMap<QString, QVariant> settings;
    settings["itemId"] = "settings";
    if (m_centralWidget->is24HourFormat())
        settings["itemText"] = tr("12 Hour Time");
    else
        settings["itemText"] = tr("24 Hour Time");
    settings["isActive"] = true;
    items.push_back(settings);

    QMap<QString, QVariant> open;
    open["itemId"] = "open";
    open["itemText"] = tr("Time Settings");
    open["isActive"] = true;
    items.push_back(open);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void DatetimePlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey)
    Q_UNUSED(checked)

    if (menuId == "open")
        QProcess::startDetached("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:datetime\"");
    else
        m_centralWidget->toggleHourFormat();
}

void DatetimePlugin::updateCurrentTimeString()
{
    const QDateTime currentDateTime = QDateTime::currentDateTime();

    if (m_centralWidget->is24HourFormat())
        m_dateTipsLabel->setText(currentDateTime.date().toString(Qt::SystemLocaleLongDate) + currentDateTime.toString(" HH:mm:ss"));
    else
        m_dateTipsLabel->setText(currentDateTime.date().toString(Qt::SystemLocaleLongDate) + currentDateTime.toString(" hh:mm:ss A"));

    const QString currentString = currentDateTime.toString("mm");

    if (currentString == m_currentTimeString)
        return;

    m_currentTimeString = currentString;
    m_centralWidget->update();
}
