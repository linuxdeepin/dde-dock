#include <QMouseEvent>

#include "pluginitemwrapper.h"

PluginItemWrapper::PluginItemWrapper(DockPluginInterface *plugin,
                                     QString uuid, QWidget * parent) :
    AbstractDockItem(parent),
    m_plugin(plugin),
    m_uuid(uuid)
{
    qDebug() << "PluginItemWrapper created " << m_plugin->name() << m_uuid;

//    setStyleSheet("PluginItemWrapper { background-color: red } ");

    if (m_plugin) {
        QWidget * item = m_plugin->getItem(uuid);
        m_pluginItemContents = m_plugin->getApplet(uuid);

        if (item) {
            setFixedSize(item->size());
            item->setParent(this);
            item->move(0, 0);

            emit widthChanged();
        }
    }
}


PluginItemWrapper::~PluginItemWrapper()
{
    qDebug() << "PluginItemWrapper destroyed " << m_plugin->name() << m_uuid;
}

QString PluginItemWrapper::getTitle()
{
    return m_plugin->getTitle(m_uuid);
}

QWidget * PluginItemWrapper::getApplet()
{
    return m_plugin->getApplet(m_uuid);
}

QString PluginItemWrapper::uuid() const
{
    return m_uuid;
}

void PluginItemWrapper::enterEvent(QEvent *)
{
    emit mouseEntered();

    showPreview();
}

void PluginItemWrapper::leaveEvent(QEvent *)
{
    emit mouseExited();

    hidePreview();
}


void PluginItemWrapper::mousePressEvent(QMouseEvent * event)
{
    if (event->button() == Qt::RightButton) {
        this->showMenu();
    }
}

QString PluginItemWrapper::getMenuContent()
{
    return m_plugin->getMenuContent(m_uuid);
}

void PluginItemWrapper::invokeMenuItem(QString itemId, bool checked)
{
    m_plugin->invokeMenuItem(m_uuid, itemId, checked);
}
