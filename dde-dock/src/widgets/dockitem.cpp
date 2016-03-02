/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QResizeEvent>

#include "dockitem.h"
#include "dbus/dbusdisplay.h"


const int CONTENT_PREVIEW_INTERVAL = 200;
const int DOCK_PREVIEW_MARGIN = 8;

class DockItemTitle : public QLabel
{
public:
    explicit DockItemTitle(QWidget * parent = 0);

    void setTitle(QString title);
};

DockItemTitle::DockItemTitle(QWidget *parent) :
    QLabel(parent)
{
    setObjectName("DockAppTitle");
    setAlignment(Qt::AlignCenter);
}

void DockItemTitle::setTitle(QString title)
{
    setText(title);

    QFontMetrics fm(font());
    int textWidth = fm.width(title);

    int fitWidth = textWidth + 20;

    resize(fitWidth < 80 ? 80 : fitWidth, 20);
}

DockItem::DockItem(QWidget * parent) :
    QFrame(parent),
    m_highlight(nullptr),
    m_dbusMenu(nullptr),
    m_dbusMenuManager(nullptr)
{

    setAttribute(Qt::WA_TranslucentBackground);

    m_titleLabel = new DockItemTitle;
    m_titlePreview = new PreviewWindow(DArrowRectangle::ArrowBottom);
}

DockItem::~DockItem()
{
    delete m_highlight;
    delete m_titleLabel;
}


bool DockItem::hoverable() const
{
    return m_hoverable;
}

void DockItem::setHoverable(bool hoverable)
{
    m_hoverable = hoverable;
}

int DockItem::globalX()
{
    return mapToGlobal(QPoint(0,0)).x();
}

int DockItem::globalY()
{
    return mapToGlobal(QPoint(0,0)).y();
}

QPoint DockItem::globalPos()
{
    return mapToGlobal(QPoint(0,0));
}

void DockItem::showPreview(const QPoint &previewPos)
{
    //make sure parentWidget() is valuable befor init hightlight
    if (!m_highlight && parentWidget()) {
        initHighlight();
    }

    if (!m_titlePreview->isHidden()) {
        m_titlePreview->resizeWithContent();
        return;
    }

    QPoint pos = previewPos.isNull()
            ?  QPoint(globalX() + width() / 2, globalY() - DOCK_PREVIEW_MARGIN)
             : previewPos;

    if (getApplet() == NULL) {
        QString title = getTitle();
        if (!title.isEmpty()) {
            m_titleLabel->setTitle(title);

            m_titlePreview->setArrowX(-1);  //reset position
            m_titlePreview->setContent(m_titleLabel);
            m_titlePreview->showPreview(pos.x(),
                                        pos.y() + DOCK_PREVIEW_MARGIN -
                                        2 - //minute adjustment
                                        m_titlePreview->shadowYOffset() +
                                        m_titlePreview->shadowBlurRadius() +
                                        m_titlePreview->shadowDistance(),
                                        CONTENT_PREVIEW_INTERVAL);
        }
    }
    else {
        m_titleLabel->setParent(NULL);

        emit needPreviewShow(pos);
    }
}

void DockItem::hidePreview(bool immediately)
{
    m_titlePreview->hidePreview(immediately);

    emit needPreviewHide(immediately);
}

void DockItem::setFixedSize(int width, int height)
{
    QFrame::setFixedSize(width, height);

    if (m_highlight) {
        m_highlight->setFixedSize(size());
        m_highlight->setVisible(false);
    }
}

void DockItem::showMenu(const QPoint &menuPos)
{
    if (getMenuContent().isEmpty()) return;

    hidePreview(true);

    if (m_dbusMenuManager == nullptr) {
        m_dbusMenuManager = new DBusMenuManager(this);
    }

    QDBusPendingReply<QDBusObjectPath> pr = m_dbusMenuManager->RegisterMenu();
    pr.waitForFinished();

    if (pr.isValid()) {
        QDBusObjectPath op = pr.value();

        if (m_dbusMenu != nullptr) {
            m_dbusMenu->deleteLater();
        }

        m_dbusMenu = new DBusMenu(op.path(), this);

        connect(m_dbusMenu, &DBusMenu::ItemInvoked, this, &DockItem::invokeMenuItem);
        connect(m_dbusMenu, &DBusMenu::MenuUnregistered, [=] {
            setHoverable(true);
        });

        QPoint pos = menuPos.isNull() ?  QPoint(globalX() + width() / 2, globalY()) : menuPos;
        QJsonObject targetObj;
        targetObj.insert("x", QJsonValue(pos.x()));
        targetObj.insert("y", QJsonValue(pos.y()));
        targetObj.insert("isDockMenu", QJsonValue(true));
        targetObj.insert("menuJsonContent", QJsonValue(getMenuContent()));

        m_dbusMenu->ShowMenu(QString(QJsonDocument(targetObj).toJson()));

        setHoverable(false);
    }
}

QString DockItem::getMenuContent()
{
    return "";
}

void DockItem::invokeMenuItem(QString, bool)
{

}

QPixmap DockItem::grab(const QRect &rectangle)
{
    return QWidget::grab(rectangle);
}

void DockItem::initHighlight()
{
    m_highlight = new HighlightEffect(this, parentWidget());

    connect(this, &DockItem::mousePress, [=](){
        m_highlight->move(pos());
        m_highlight->showDarker();
    });
    connect(this, &DockItem::mouseRelease, [=](){
        m_highlight->move(pos());
        m_highlight->showLighter();
    });
    connect(this, &DockItem::mouseEnter, [=](){
        m_highlight->move(pos());
        m_highlight->showLighter();
        m_highlight->setVisible(true);
    });
    connect(this, &DockItem::mouseLeave, [=](){
        if (!m_highlight->isVisible())
            return;
        m_highlight->setVisible(false);
        m_highlight->move(pos());
    });
}

void DockItem::resizeEvent(QResizeEvent * event)
{
    Q_UNUSED(event)

    if (m_highlight) {
        m_highlight->setFixedSize(size());
        m_highlight->setVisible(false);
    }
}

void DockItem::moveEvent(QMoveEvent *)
{
    if (m_highlight)
        m_highlight->setVisible(false);
}
