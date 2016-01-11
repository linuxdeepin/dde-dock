#include <QMouseEvent>
#include <QJsonArray>
#include <QJsonDocument>
#include <QProcess>

#include "dockpluginitemwrapper.h"
#include "../dockmodedata.h"

static const QString MenuItemRun = "id_run";
static const QString MenuItemRemove = "id_remove";

DockPluginItemWrapper::DockPluginItemWrapper(DockPluginInterface *plugin,
                                     QString id, QWidget * parent) :
    DockItem(parent),
    m_plugin(plugin),
    m_id(id)
{
    if (m_plugin) {
        qDebug() << "PluginItemWrapper created " << m_plugin->getPluginName() << m_id;

        QWidget * item = m_plugin->getItem(id);
        m_pluginItemContents = m_plugin->getApplet(id);

        if (item) {
            QHBoxLayout *layout = new QHBoxLayout(this);
            layout->setContentsMargins(0, 0, 0, 0);
            layout->setSpacing(0);
            layout->addWidget(item);
            adjustSize();

            m_display = new DBusDisplay(this);
        }
    }
}


DockPluginItemWrapper::~DockPluginItemWrapper()
{
    qDebug() << "PluginItemWrapper destroyed " << m_plugin->getPluginName() << m_id;
}

QString DockPluginItemWrapper::getTitle()
{
    return m_plugin->getTitle(m_id);
}

QWidget * DockPluginItemWrapper::getApplet()
{
    return m_plugin->getApplet(m_id);
}

QString DockPluginItemWrapper::getItemId()
{
    return m_id;
}

void DockPluginItemWrapper::enterEvent(QEvent *)
{
    if (DockModeData::instance()->getHideState() != Dock::HideStateShown)
        return;

//    emit mouseEntered();

    if (hoverable()) {
        DisplayRect rec = m_display->primaryRect();
        showPreview(QPoint(globalX() + width() / 2, rec.height- DockModeData::instance()->getDockHeight() - DOCK_PREVIEW_MARGIN));
    }
}

void DockPluginItemWrapper::leaveEvent(QEvent *)
{
//    emit mouseExited();

    hidePreview();
}


void DockPluginItemWrapper::mousePressEvent(QMouseEvent * event)
{
    hidePreview(true);

    if (event->button() == Qt::RightButton) {
        DisplayRect rec = m_display->primaryRect();
        this->showMenu(QPoint(rec.x + globalX() + width() / 2,
                              rec.y + rec.height - DockModeData::instance()->getDockHeight()));
    } else if (event->button() == Qt::LeftButton) {
        QString command = m_plugin->getCommand(m_id);
        if (!command.isEmpty()) QProcess::startDetached(command);
    }
}

QString DockPluginItemWrapper::getMenuContent()
{
    QString menuContent = m_plugin->getMenuContent(m_id);

    bool canRun = !m_plugin->getCommand(m_id).isEmpty();
    bool configurable = m_plugin->configurable(m_id);

    if (canRun || configurable) {
        QJsonObject result = QJsonDocument::fromJson(menuContent.toUtf8()).object();
        QJsonArray array = result["items"].toArray();

        QJsonObject itemRun = createMenuItem(MenuItemRun, tr("_Run"), false, false);
        QJsonObject itemRemove = createMenuItem(MenuItemRemove, tr("_Undock"), false, false);

        if (canRun) array.insert(0, itemRun);
        if (configurable) array.append(itemRemove);

        result["items"] = array;

        return QString(QJsonDocument(result).toJson());
    } else {
        return menuContent;
    }
}

void DockPluginItemWrapper::invokeMenuItem(QString itemId, bool checked)
{
    if (itemId == MenuItemRun) {
        QString command = m_plugin->getCommand(m_id);
        QProcess::startDetached(command);
    } else if (itemId == MenuItemRemove){
        m_plugin->setEnabled(m_id, false);
    } else {
        m_plugin->invokeMenuItem(m_id, itemId, checked);
    }
}

QJsonObject DockPluginItemWrapper::createMenuItem(QString itemId, QString itemName, bool checkable, bool checked)
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
