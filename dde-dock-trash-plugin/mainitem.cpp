#include "mainitem.h"

MainItem::MainItem(QWidget *parent) : QLabel(parent)
{
    setAcceptDrops(true);
    setFixedSize(Dock::APPLET_FASHION_ITEM_WIDTH, Dock::APPLET_FASHION_ITEM_HEIGHT);

    updateIcon(false);
}

MainItem::~MainItem()
{

}

void MainItem::emptyTrash()
{
    QDBusPendingReply<QString, QDBusObjectPath, QString> tmpReply = m_dfo->NewEmptyTrashJob(false, "", "", "");
    QDBusObjectPath op = tmpReply.argumentAt(1).value<QDBusObjectPath>();
    DBusEmptyTrashJob * detj = new DBusEmptyTrashJob(op.path(), this);
    connect(detj, &DBusEmptyTrashJob::Done, detj, &DBusEmptyTrashJob::deleteLater);
    connect(detj, &DBusEmptyTrashJob::Done, [=](){
        updateIcon(false);
    });

    if (detj->isValid())
        detj->Execute();

}

void MainItem::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
    {
        QProcess * tmpProcess = new QProcess();
        connect(tmpProcess, SIGNAL(finished(int)), tmpProcess, SLOT(deleteLater()));
        tmpProcess->start("nautilus trash://");
    }

    //Makesure it parent can accept the mouse event too
    event->ignore();
}

void MainItem::dragEnterEvent(QDragEnterEvent * event)
{
    if (event->source())
        return;//just accept the object outside this app
    updateIcon(true);

    event->setDropAction(Qt::MoveAction);
    event->accept();
}

void MainItem::dragLeaveEvent(QDragLeaveEvent *)
{
    updateIcon(false);
}

void MainItem::dropEvent(QDropEvent *event)
{
    updateIcon(false);

    if (event->source())
        return;

    QStringList formats = event->mimeData()->formats();
    if (formats.indexOf("_DEEPIN_DND") != -1 && formats.indexOf("text/plain") != -1)//Application
    {
        QString jsonStr = QString(event->mimeData()->data("_DEEPIN_DND")).split("uninstall").last().trimmed();

        //Tim at both ends of the string is added to a character SOH (start of heading)
        jsonStr = jsonStr.mid(1,jsonStr.length() - 2);
        QJsonObject dataObj = QJsonDocument::fromJson(jsonStr.trimmed().toUtf8()).object();
        if (dataObj.isEmpty())
            return;

        qWarning() << "Uninstall application:" << dataObj.value("id").toString();
        m_launcher->RequestUninstall(dataObj.value("id").toString(), true);
    }
    else//File or Dirctory
    {
        QStringList files;
        foreach (QUrl fileUrl, event->mimeData()->urls())
            files << fileUrl.path();

        QDBusPendingReply<QString, QDBusObjectPath, QString> tmpReply = m_dfo->NewTrashJob(files, false, "", "", "");
        QDBusObjectPath op = tmpReply.argumentAt(1).value<QDBusObjectPath>();
        DBusTrashJob * dtj = new DBusTrashJob(op.path(), this);
        connect(dtj, &DBusTrashJob::Done, dtj, &DBusTrashJob::deleteLater);
        connect(dtj, &DBusTrashJob::Done, [=](){
            updateIcon(false);
        });

        if (dtj->isValid())
            dtj->Execute();

        qWarning()<< op.path() << "Move files to trash: "<< files;
    }
}

void MainItem::updateIcon(bool isOpen)
{
    QString iconName = "";
    if (isOpen)
    {
        if (m_dftm->ItemCount() > 0)
            iconName = "user-trash-full-opened";
        else
            iconName = "user-trash-empty-opened";
    }
    else
    {
        if (m_dftm->ItemCount() > 0)
            iconName = "user-trash-full";
        else
            iconName = "user-trash-empty";
    }

    QPixmap pixmap = QIcon::fromTheme(iconName).pixmap(Dock::APPLET_FASHION_ICON_SIZE,Dock::APPLET_FASHION_ICON_SIZE);
    setPixmap(pixmap);
}
