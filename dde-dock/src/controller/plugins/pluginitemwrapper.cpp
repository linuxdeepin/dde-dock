#include <QMouseEvent>
#include <QJsonArray>
#include <QJsonDocument>
#include <QProcess>

#include "pluginitemwrapper.h"

static const QString MenuItemRun = "id_run";
static const QString MenuItemRemove = "id_remove";

PluginItemWrapper::PluginItemWrapper(DockPluginInterface *plugin,
                                     QString id, QWidget * parent) :
    AbstractDockItem(parent),
    m_plugin(plugin),
    m_id(id)
{
    qDebug() << "PluginItemWrapper created " << m_plugin->getPluginName() << m_id;

    if (m_plugin) {
        QWidget * item = m_plugin->getItem(id);
        m_pluginItemContents = m_plugin->getApplet(id);

        if (item) {
            item->setParent(this);
            item->move(0, 0);
            this->adjustSize();

            emit widthChanged();
        }
    }
}


PluginItemWrapper::~PluginItemWrapper()
{
    qDebug() << "PluginItemWrapper destroyed " << m_plugin->getPluginName() << m_id;
}

QString PluginItemWrapper::getTitle()
{
    return m_plugin->getTitle(m_id);
}

QWidget * PluginItemWrapper::getApplet()
{
    return m_plugin->getApplet(m_id);
}

QString PluginItemWrapper::id() const
{
    return m_id;
}

void PluginItemWrapper::enterEvent(QEvent *)
{
    emit mouseEntered();

    if (hoverable())
        showPreview();
}

void PluginItemWrapper::leaveEvent(QEvent *)
{
    emit mouseExited();

    hidePreview();
}


void PluginItemWrapper::mousePressEvent(QMouseEvent * event)
{
    hidePreview(true);

    if (event->button() == Qt::RightButton) {
        this->showMenu();
    } else if (event->button() == Qt::LeftButton) {
        QString command = m_plugin->getCommand(m_id);
        if (!command.isEmpty()) QProcess::startDetached(command);
    }
}

QString PluginItemWrapper::getMenuContent()
{
    QString menuContent = m_plugin->getMenuContent(m_id);

    bool canRun = !m_plugin->getCommand(m_id).isEmpty();
    bool canDisable = m_plugin->canDisable(m_id);

    if (canRun || canDisable) {
        QJsonObject result = QJsonDocument::fromJson(menuContent.toUtf8()).object();
        QJsonArray array = result["items"].toArray();

        QJsonObject itemRun = createMenuItem(MenuItemRun, tr("_Run"), false, false);
        QJsonObject itemRemove = createMenuItem(MenuItemRemove, tr("_Undock"), false, false);

        if (canRun) array.insert(0, itemRun);
        if (canDisable) array.append(itemRemove);

        result["items"] = array;

        return QString(QJsonDocument(result).toJson());
    } else {
        return menuContent;
    }
}

void PluginItemWrapper::invokeMenuItem(QString itemId, bool checked)
{
    if (itemId == MenuItemRun) {
        QString command = m_plugin->getCommand(m_id);
        QProcess::startDetached(command);
    } else if (itemId == MenuItemRemove){
        m_plugin->setDisabled(m_id, true);
    } else {
        m_plugin->invokeMenuItem(m_id, itemId, checked);
    }
}

QJsonObject PluginItemWrapper::createMenuItem(QString itemId, QString itemName, bool checkable, bool checked)
{
    QJsonObject itemObj;

    itemObj.insert("itemId", itemId);
    itemObj.insert("itemText", itemName);
    itemObj.insert("itemIcon", "");
    itemObj.insert("itemIconHover", "");
    itemObj.insert("itemIconInactive", "");
    itemObj.insert("itemExtra", "");
    itemObj.insert("isActive", true);
    itemObj.insert("isCheckable", checkable);
    itemObj.insert("checked", checked);
    itemObj.insert("itemSubMenu", QJsonObject());

    return itemObj;
}
