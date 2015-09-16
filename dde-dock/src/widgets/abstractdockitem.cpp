#include "abstractdockitem.h"

ItemTitleLabel::ItemTitleLabel(QWidget *parent) :
    QLabel(parent)
{
    setObjectName("DockAppTitle");
    setAlignment(Qt::AlignCenter);
}

void ItemTitleLabel::setTitle(QString title)
{
    setText(title);

    QFontMetrics fm(font());
    int textWidth = fm.width(title);

    int fitWidth = textWidth + 20;

    resize(fitWidth < 80 ? 80 : fitWidth, 20);
}


AbstractDockItem::AbstractDockItem(QWidget * parent) :
    QFrame(parent)
{

    this->setAttribute(Qt::WA_TranslucentBackground);
    m_previewAR = new PreviewArrowRectangle();

    connect(m_previewAR, &PreviewArrowRectangle::hideFinish, this, &AbstractDockItem::previewHidden);
}

AbstractDockItem::~AbstractDockItem()
{
    delete m_highlight;
    delete m_titleLabel;
}

QString AbstractDockItem::getItemId()
{
    return "";
}

QString AbstractDockItem::getTitle()
{
    return "";
}

QWidget * AbstractDockItem::getApplet()
{
    return NULL;
}

bool AbstractDockItem::moveable()
{
    return m_moveable;
}

bool AbstractDockItem::actived()
{
    return m_isActived;
}

void AbstractDockItem::resize(int width,int height){
    QFrame::resize(width,height);

    if (m_highlight)
        m_highlight->setFixedSize(size());

    emit widthChanged();
}

void AbstractDockItem::resize(const QSize &size){
    QFrame::resize(size);

    if (m_highlight)
        m_highlight->setFixedSize(size);

    emit widthChanged();
}

QPoint AbstractDockItem::getNextPos()
{
    return m_itemNextPos;
}

void AbstractDockItem::setNextPos(const QPoint &value)
{
    m_itemNextPos = value;
}

void AbstractDockItem::setNextPos(int x, int y)
{
    m_itemNextPos.setX(x); m_itemNextPos.setY(y);
}

void AbstractDockItem::move(const QPoint &value)
{
    QWidget::move(value);
    m_highlight->move(pos());

    emit posChanged();
}

void AbstractDockItem::move(int x, int y)
{
    QWidget::move(x,y);
    m_highlight->move(pos());

    emit posChanged();
}

int AbstractDockItem::globalX()
{
    return mapToGlobal(QPoint(0,0)).x();
}

int AbstractDockItem::globalY()
{
    return mapToGlobal(QPoint(0,0)).y();
}

QPoint AbstractDockItem::globalPos()
{
    return mapToGlobal(QPoint(0,0));
}

void AbstractDockItem::showPreview()
{
    if (!m_previewAR->isHidden())
    {
        m_previewAR->resizeWithContent();
        return;
    }

    QWidget *tmpContent = getApplet();

    if (tmpContent == NULL) {
        QString title = getTitle();

        if (!title.isEmpty()) {
            m_titleLabel->setTitle(title);

            tmpContent = m_titleLabel;
        }
    }

    if (tmpContent) {
        m_titleLabel->setParent(NULL);
        m_previewAR->setContent(tmpContent);
        m_previewAR->showPreview(ArrowRectangle::ArrowBottom,
                                 globalX() + width() / 2,
                                 globalY() - 5);
        m_previewPos = QPoint(globalX() + width() / 2,globalY() - 5);
    }
}

void AbstractDockItem::hidePreview(int interval)
{
    m_previewAR->hidePreview(interval);
}

void AbstractDockItem::cancelHide()
{
    m_previewAR->showPreview(ArrowRectangle::ArrowBottom, m_previewPos.x(), m_previewPos.y());
}

void AbstractDockItem::resizePreview()
{
    m_previewAR->resizeWithContent();
    if (!m_previewAR->isHidden())
    {
        m_previewAR->resizeWithContent();
        return;
    }
}

void AbstractDockItem::showMenu()
{
    if (getMenuContent().isEmpty()) return;

    hidePreview(0);

    if (m_dbusMenuManager == NULL) {
        m_dbusMenuManager = new DBusMenuManager(this);
    }

    QDBusPendingReply<QDBusObjectPath> pr = m_dbusMenuManager->RegisterMenu();
    pr.waitForFinished();

    if (pr.isValid()) {
        QDBusObjectPath op = pr.value();

        if (m_dbusMenu != NULL) {
            m_dbusMenu->deleteLater();
        }

        m_dbusMenu = new DBusMenu(op.path(), this);

        connect(m_dbusMenu, &DBusMenu::ItemInvoked, this, &AbstractDockItem::invokeMenuItem);

        QJsonObject targetObj;
        targetObj.insert("x", QJsonValue(globalX() + width() / 2));
        targetObj.insert("y", QJsonValue(globalY() - 5));
        targetObj.insert("isDockMenu", QJsonValue(true));
        targetObj.insert("menuJsonContent", QJsonValue(getMenuContent()));

        m_dbusMenu->ShowMenu(QString(QJsonDocument(targetObj).toJson()));
    }
}

QString AbstractDockItem::getMenuContent()
{
    return "";
}

void AbstractDockItem::invokeMenuItem(QString, bool)
{

}

void AbstractDockItem::setParent(QWidget *parent)
{
    QWidget::setParent(parent);
    initHighlight();
    initTitleLabel();
}

void AbstractDockItem::initHighlight()
{
    //the size and position will update with move() and resize()
    QWidget * lParent = qobject_cast<QWidget *>(parent());
    if (lParent)
    {
        if (!m_highlight)
        {
            m_highlight = new HighlightEffect(this, lParent);
            connect(this, &AbstractDockItem::dragStart, [=](){
                m_highlight->setVisible(false);
            });
            connect(this, &AbstractDockItem::mousePress, [=](){
                m_highlight->showDarker();
                emit frameUpdate();
            });
            connect(this, &AbstractDockItem::mouseRelease, [=](){
                m_highlight->showLighter();
                emit frameUpdate();
            });
            connect(this, &AbstractDockItem::mouseEntered, [=](){
                m_highlight->showLighter();
                emit frameUpdate();
            });
            connect(this, &AbstractDockItem::mouseExited, [=](){
                if (!m_highlight->isVisible())
                    return;
                m_highlight->showNormal();
                emit frameUpdate();
            });
        }
        else
            m_highlight->setParent(lParent);
    }
}

void AbstractDockItem::initTitleLabel()
{
    m_titleLabel = new ItemTitleLabel;
}
