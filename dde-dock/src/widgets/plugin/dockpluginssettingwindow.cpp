/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "dockpluginssettingwindow.h"

class PluginSettingLine : public QLabel
{
    Q_OBJECT
public:
    explicit PluginSettingLine(bool checked = false,
                                const QString &id = "",
                                const QString &title = "",
                                const QPixmap &icon = QPixmap(),
                                QWidget *parent = 0);

    void setIcon(const QPixmap &icon);
    void setTitle(const QString &title);
    void setPluginId(const QString &pluginId);
    void setChecked(const bool checked);

    QString pluginId() const;
    bool checked();

signals:
    void checkedChanged(QString id, bool check);

private:
    QLabel *m_iconLabel = NULL;
    QLabel *m_titleLabel = NULL;
    DSwitchButton *m_switchButton = NULL;

    QString m_pluginId = "";

    const int ICON_SIZE = 16;
    const int ICON_SPACING = 6;
    const int CONTENT_MARGIN = 6;
    const int MAX_TEXT_WIDTH = 125;
};

PluginSettingLine::PluginSettingLine(bool checked, const QString &id, const QString &title, const QPixmap &icon, QWidget *parent)
    :QLabel(parent), m_pluginId(id)
{
    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(CONTENT_MARGIN, 0, CONTENT_MARGIN, 0);
    mainLayout->setSpacing(0);

    m_iconLabel = new QLabel;
    m_iconLabel->setFixedSize(ICON_SIZE, ICON_SIZE);
    m_iconLabel->setPixmap(icon.scaled(ICON_SIZE, ICON_SIZE));

    m_titleLabel = new QLabel;
    m_titleLabel->setObjectName("PluginSettingLineTitle");
    m_titleLabel->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
    setTitle(title);

    m_switchButton = new DSwitchButton;
    m_switchButton->setChecked(checked);

    connect(m_switchButton, &DSwitchButton::checkedChanged, [=](bool checked){
        emit checkedChanged(m_pluginId, checked);
    });

    mainLayout->addWidget(m_iconLabel);
    mainLayout->addSpacing(ICON_SPACING);
    mainLayout->addWidget(m_titleLabel);
    mainLayout->addStretch();
    mainLayout->addWidget(m_switchButton);
}

void PluginSettingLine::setTitle(const QString &title)
{
    m_titleLabel->setText(title);

    QFontMetrics fm(m_titleLabel->font());
    m_titleLabel->setText(fm.elidedText(title,Qt::ElideRight, MAX_TEXT_WIDTH));
}

void PluginSettingLine::setIcon(const QPixmap &icon)
{
    m_iconLabel->setPixmap(icon.scaled(ICON_SIZE, ICON_SIZE));
}

QString PluginSettingLine::pluginId() const
{
    return m_pluginId;
}

bool PluginSettingLine::checked()
{
    return m_switchButton->checked();
}

void PluginSettingLine::setPluginId(const QString &pluginId)
{
    m_pluginId = pluginId;
}

void PluginSettingLine::setChecked(const bool checked)
{
    m_switchButton->setChecked(checked);
}

#include "dockpluginssettingwindow.moc"
///////////////////////////////////////////////////////////////////////////////////////////////////


DockPluginsSettingWindow::DockPluginsSettingWindow(QWidget *parent) :
    QFrame(parent)
{
    setObjectName("PluginsSettingFrame");
    setAttribute(Qt::WA_TranslucentBackground);
    setWindowFlags(Qt::FramelessWindowHint | Qt::Dialog);

    QWidget *contentFrame = new QWidget;
    contentFrame->setObjectName("PluginsSettingFrame");
    QVBoxLayout *contentLayout = new QVBoxLayout(this);
    contentLayout->setContentsMargins(0, 0, 0, 0);
    contentLayout->setSpacing(0);
    contentLayout->addWidget(contentFrame);

    m_mainLayout = new QVBoxLayout(contentFrame);
    m_mainLayout->setSpacing(LINE_SPACING);
    m_mainLayout->setContentsMargins(0, 0, 0, 0);

    initCloseTitle();

    setFixedWidth(WIN_WIDTH);

    installEventFilter(this);
}

void DockPluginsSettingWindow::onPluginAdd(bool checked, const QString &id, const QString &title, const QPixmap &icon)
{
    if (m_lineMap.keys().indexOf(id) != -1)
        return;

    PluginSettingLine *line = new PluginSettingLine(checked, id, title, icon);
    connect(line, &PluginSettingLine::checkedChanged, this, &DockPluginsSettingWindow::checkedChanged);

    line->setFixedHeight(LINE_HEIGHT);
    m_mainLayout->addWidget(line, 1, Qt::AlignTop);

    m_lineMap.insert(id, line);

    resizeWithLineCount();
}

void DockPluginsSettingWindow::onPluginRemove(const QString &id)
{
    PluginSettingLine * line = m_lineMap.take(id);
    if (line){
        m_mainLayout->removeWidget(line);
        line->deleteLater();

        resizeWithLineCount();
    }
}

void DockPluginsSettingWindow::onPluginEnabledChanged(const QString &id, bool enabled)
{
    PluginSettingLine * line = m_lineMap.value(id);
    if (line) {
        line->setChecked(enabled);
    }
}

void DockPluginsSettingWindow::onPluginTitleChanged(const QString &id, const QString &title)
{
    PluginSettingLine * line = m_lineMap.value(id);
    if (line) {
        line->setTitle(title);
    }
}

void DockPluginsSettingWindow::mouseMoveEvent(QMouseEvent *event)
{
    if (m_mousePressed)
        move(event->globalPos() - m_pressPosition);
    QFrame::mouseMoveEvent(event);
}

void DockPluginsSettingWindow::mousePressEvent(QMouseEvent *event)
{
    if(event->button() & Qt::LeftButton)
    {
        m_pressPosition = event->globalPos() - frameGeometry().topLeft();
        m_mousePressed = true;
    }
    QFrame::mousePressEvent(event);
}

void DockPluginsSettingWindow::mouseReleaseEvent(QMouseEvent *event)
{
    m_mousePressed = false;
    QFrame::mouseReleaseEvent(event);
}

bool DockPluginsSettingWindow::eventFilter(QObject *obj, QEvent *event)
{
    if (event->type() == QEvent::WindowDeactivate) {
        this->close();
    }

    return QFrame::eventFilter(obj, event);
}

void DockPluginsSettingWindow::resizeWithLineCount()
{
    setFixedHeight((m_lineMap.count() + 1) * (LINE_HEIGHT + LINE_SPACING));
}

void DockPluginsSettingWindow::initCloseTitle()
{
    QLabel *titleLabel = new QLabel(tr("Notification Area Settings"));
    titleLabel->setAlignment(Qt::AlignCenter);
    titleLabel->setObjectName("PluginSettingTitle");
    QPushButton *closeButton = new QPushButton;
    closeButton->setFocusPolicy(Qt::NoFocus);
    closeButton->setFixedSize(ICON_SIZE, ICON_SIZE);
    closeButton->setObjectName("PluginSettingCloseButton");
    connect(closeButton, &QPushButton::clicked, [=]{
        this->hide();
    });

    QHBoxLayout *titleLayout = new QHBoxLayout;
    titleLayout->setAlignment(Qt::AlignVCenter);
    titleLayout->setContentsMargins(0, CONTENT_MARGIN, CONTENT_MARGIN, 0);
    titleLayout->setSpacing(0);

    titleLayout->addWidget(titleLabel, 1);
    titleLayout->addWidget(closeButton, 1);

    m_mainLayout->addLayout(titleLayout, 1);
    DSeparatorHorizontal *sp = new DSeparatorHorizontal;
    m_mainLayout->addWidget(sp);
}
