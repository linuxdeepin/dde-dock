#include "diskmountplugin.h"

DiskMountPlugin::DiskMountPlugin(QObject *parent)
    : QObject(parent),

      m_pluginAdded(false),

      m_tipsLabel(new QLabel),
      m_diskPluginItem(new DiskPluginItem),
      m_diskControlApplet(nullptr)
{
    m_diskPluginItem->setVisible(false);
    m_tipsLabel->setVisible(false);
    m_tipsLabel->setText(tr("Disk"));
    m_tipsLabel->setStyleSheet("color:white;"
                               "padding:5px 10px;");
}

const QString DiskMountPlugin::pluginName() const
{
    return "disk-mount";
}

void DiskMountPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    initCompoments();
    m_diskPluginItem->setDockDisplayMode(displayMode());
}

QWidget *DiskMountPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_diskPluginItem;
}

QWidget *DiskMountPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_tipsLabel;
}

QWidget *DiskMountPlugin::itemPopupApplet(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_diskControlApplet;
}

void DiskMountPlugin::initCompoments()
{
    m_diskControlApplet = new DiskControlWidget;
    m_diskControlApplet->setObjectName("dist-mount");
    m_diskControlApplet->setVisible(false);

    connect(m_diskControlApplet, &DiskControlWidget::diskCountChanged, this, &DiskMountPlugin::diskCountChanged);
}

void DiskMountPlugin::displayModeChanged(const Dock::DisplayMode mode)
{
    m_diskPluginItem->setDockDisplayMode(mode);
}

void DiskMountPlugin::diskCountChanged(const int count)
{
    if (m_pluginAdded == bool(count))
        return;

    m_pluginAdded = bool(count);

    if (m_pluginAdded)
        m_proxyInter->itemAdded(this, QString());
    else
        m_proxyInter->itemRemoved(this, QString());
}
