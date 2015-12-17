#include "shutdownplugin.h"

#include <QLabel>
#include <QDebug>
#include <QTimer>

const QString PLUGIN_ID = "shutdown";

ShutdownPlugin::ShutdownPlugin(QObject *parent)
    : QObject(parent)
{
    m_mainWidget = new QLabel;
}

QString ShutdownPlugin::getPluginName()
{
    return QString(tr("Shutdown"));
}

void ShutdownPlugin::init(DockPluginProxyInterface *proxy)
{
    m_proxy = proxy;
    m_proxy->itemAddedEvent(PLUGIN_ID);

    changeMode(m_proxy->dockMode(), m_proxy->dockMode());
}

void ShutdownPlugin::changeMode(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    Q_UNUSED(oldMode)

    switch (newMode)
    {
    case Dock::FashionMode:
        m_mainWidget->setFixedSize(48, 48);
        m_mainWidget->setPixmap(QPixmap(":/icons/icons/fashion.svg"));
        break;
    case Dock::ClassicMode:
    case Dock::EfficientMode:
        m_mainWidget->setFixedSize(18, 18);
        m_mainWidget->setPixmap(QPixmap(":/icons/icons/normal.svg"));
        break;
    }

    m_proxy->infoChangedEvent(DockPluginInterface::InfoTypeItemSize, PLUGIN_ID);
}

QStringList ShutdownPlugin::ids()
{
    return QStringList();
}

QString ShutdownPlugin::getName(QString id)
{
    Q_UNUSED(id)

    return getPluginName();
}

QString ShutdownPlugin::getTitle(QString id)
{
    Q_UNUSED(id)

    return getPluginName();
}

QString ShutdownPlugin::getCommand(QString id)
{
    Q_UNUSED(id)

    return QString("dde-shutdown");
}

QWidget *ShutdownPlugin::getItem(QString id)
{
    Q_UNUSED(id)

    return m_mainWidget;
}

QWidget *ShutdownPlugin::getApplet(QString id)
{
    Q_UNUSED(id)

    return nullptr;
}

QString ShutdownPlugin::getMenuContent(QString id)
{
    Q_UNUSED(id)

    return QString();
}

void ShutdownPlugin::invokeMenuItem(QString id, QString itemId, bool checked)
{
    Q_UNUSED(id)
    Q_UNUSED(itemId)
    Q_UNUSED(checked)
}

void ShutdownPlugin::setEnabled(const QString &id, bool enabled)
{
    Q_UNUSED(id)
    Q_UNUSED(enabled)
}

bool ShutdownPlugin::configurable(const QString &id)
{
    Q_UNUSED(id)

    return false;
}

bool ShutdownPlugin::enabled(const QString &id)
{
    Q_UNUSED(id)

    return true;
}

QPixmap ShutdownPlugin::getIcon(QString id)
{
    Q_UNUSED(id);

    return QPixmap();
}
