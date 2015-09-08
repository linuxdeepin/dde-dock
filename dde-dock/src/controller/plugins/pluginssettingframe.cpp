#include "pluginssettingframe.h"

PluginsSettingLine::PluginsSettingLine(bool checked, const QString &id, const QString &title, const QPixmap &icon, QWidget *parent)
    :m_pluginId(id), QLabel(parent)
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
        emit disableChanged(m_pluginId, !checked);
    });

    mainLayout->addWidget(m_iconLabel);
    mainLayout->addSpacing(ICON_SPACING);
    mainLayout->addWidget(m_titleLabel);
    mainLayout->addStretch();
    mainLayout->addWidget(m_switchButton);
}

void PluginsSettingLine::setTitle(const QString &title)
{
    m_titleLabel->setText(title);

    QFontMetrics fm(m_titleLabel->font());
    m_titleLabel->setText(fm.elidedText(title,Qt::ElideRight, MAX_TEXT_WIDTH));
}

void PluginsSettingLine::setIcon(const QPixmap &icon)
{
    m_iconLabel->setPixmap(icon.scaled(ICON_SIZE, ICON_SIZE));
}

QString PluginsSettingLine::pluginId() const
{
    return m_pluginId;
}

bool PluginsSettingLine::checked()
{
    return m_switchButton->checked();
}

void PluginsSettingLine::setPluginId(const QString &pluginId)
{
    m_pluginId = pluginId;
}

void PluginsSettingLine::setChecked(const bool checked)
{
    m_switchButton->setChecked(checked);
}


PluginsSettingFrame::PluginsSettingFrame(QWidget *parent) :
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
}

void PluginsSettingFrame::onPluginAdd(bool checked, const QString &id, const QString &title, const QPixmap &icon)
{
    if (m_lineMap.keys().indexOf(id) != -1)
        return;

    PluginsSettingLine *line = new PluginsSettingLine(checked, id, title, icon);
    connect(line, &PluginsSettingLine::disableChanged, this, &PluginsSettingFrame::disableChanged);

    m_mainLayout->addWidget(line, 1, Qt::AlignTop);

    m_lineMap.insert(id, line);

    resizeWithLineCount();
}

void PluginsSettingFrame::onPluginRemove(const QString &id)
{
    PluginsSettingLine * line = m_lineMap.take(id);
    if (line){
        m_mainLayout->removeWidget(line);
        line->deleteLater();

        resizeWithLineCount();
    }
}

void PluginsSettingFrame::clear()
{
    foreach (QString uuid, m_lineMap.keys()) {
        m_lineMap.take(uuid)->deleteLater();
    }
}

void PluginsSettingFrame::mouseMoveEvent(QMouseEvent *event)
{
    move(event->globalPos() - m_pressPosition);
    QFrame::mouseMoveEvent(event);
}

void PluginsSettingFrame::mousePressEvent(QMouseEvent *event)
{
    if(event->button() & Qt::LeftButton)
    {
        m_pressPosition = event->globalPos() - frameGeometry().topLeft();
    }
    QFrame::mousePressEvent(event);
}

void PluginsSettingFrame::mouseReleaseEvent(QMouseEvent *event)
{
    QFrame::mouseReleaseEvent(event);
}

void PluginsSettingFrame::resizeWithLineCount()
{
    setFixedHeight((m_lineMap.count() + 1) * (LINE_HEIGHT + LINE_SPACING));
}

void PluginsSettingFrame::initCloseTitle()
{
    QLabel *titleLabel = new QLabel(tr("Notice Region Setting"));
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
