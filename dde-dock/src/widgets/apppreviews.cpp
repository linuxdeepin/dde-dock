#include "apppreviews.h"

AppPreviewFrame::AppPreviewFrame(const QString &title, int xid, QWidget *parent) :
    QFrame(parent),xidValue(xid)
{
    addPreview(xid);
    setTitle(title);
    addCloseButton();
}

AppPreviewFrame::~AppPreviewFrame()
{

}

void AppPreviewFrame::addPreview(int xid)
{
    m_preview = new WindowPreview(xid, this);
    m_preview->resize(Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);

    setFixedSize(m_preview->width() + BUTTON_SIZE / 2, m_preview->height() + BUTTON_SIZE);
    m_preview->move(0, BUTTON_SIZE / 2);
}

void AppPreviewFrame::setTitle(const QString &t)
{
    QLabel *titleLabel = new QLabel(this);
    titleLabel->setObjectName("AppPreviewFrameTitle");
    QFontMetrics fm(titleLabel->font());
    titleLabel->setText(fm.elidedText(t,Qt::ElideRight,width()));
    titleLabel->setAlignment(Qt::AlignCenter);
    titleLabel->resize(width() - BUTTON_SIZE / 2, TITLE_HEIGHT);
    titleLabel->move(0, height() - titleLabel->height() - BUTTON_SIZE / 2);
}

void AppPreviewFrame::mousePressEvent(QMouseEvent *)
{
    emit activate(xidValue);
}

void AppPreviewFrame::enterEvent(QEvent *)
{
    m_preview->setIsHover(true);

    showCloseButton();
}

void AppPreviewFrame::leaveEvent(QEvent *)
{
    m_preview->setIsHover(false);

    hideCloseButton();
}

void AppPreviewFrame::addCloseButton()
{
    m_cb = new QPushButton(this);
    m_cb->setObjectName("PreviewCloseButton");
    m_cb->setFixedSize(BUTTON_SIZE, BUTTON_SIZE);
    m_cb->move(width() - m_cb->width(), 0);
    m_cb->hide();

    connect(m_cb,&QPushButton::clicked,[=]{close(this->xidValue);});
}

void AppPreviewFrame::showCloseButton()
{
    m_cb->show();
}

void AppPreviewFrame::hideCloseButton()
{
    m_cb->hide();
}

AppPreviews::AppPreviews(QWidget *parent) : QWidget(parent)
{
    m_mainLayout = new QHBoxLayout(this);
    m_mainLayout->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
    m_mainLayout->setContentsMargins(20 - PREVIEW_PADDING, 0, 0, 0);
    m_mainLayout->setSpacing(Dock::APP_PREVIEW_MARGIN - BUTTON_SIZE / 2);
    setLayout(m_mainLayout);
    resize(Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);
}

void AppPreviews::addItem(const QString &title, int xid)
{
    if (m_xidList.indexOf(xid) != -1)
        return;
    m_xidList.append(xid);

    AppPreviewFrame *f = new AppPreviewFrame(title, xid);
    connect(f,&AppPreviewFrame::close,this,&AppPreviews::removePreview);
    connect(f,&AppPreviewFrame::activate,this,&AppPreviews::activatePreview);

    m_mainLayout->addWidget(f);

    resize(getContentSize());
}

void AppPreviews::leaveEvent(QEvent *)
{
    if (m_isClosing)
        m_isClosing = false;
}

void AppPreviews::removePreview(int xid)
{
    m_isClosing = true;
    m_mainLayout->removeWidget(qobject_cast<AppPreviewFrame *>(sender()));
    sender()->deleteLater();
    m_clientManager->CloseWindow(xid);
    if (m_mainLayout->count() <= 0)
    {
        emit requestHide();
        return;
    }

    resize(getContentSize());
    emit sizeChanged();
}

void AppPreviews::activatePreview(int xid)
{
    m_clientManager->ActiveWindow(xid);
}

void AppPreviews::clearUpPreview()
{
    QLayoutItem *child;
    while ((child = m_mainLayout->takeAt(0)) != 0){
        child->widget()->deleteLater();
        delete child;
    }

    m_xidList.clear();
}

QSize AppPreviews::getContentSize()
{

    int contentWidth = m_mainLayout->count() * (Dock::APP_PREVIEW_WIDTH + Dock::APP_PREVIEW_MARGIN)
            + Dock::APP_PREVIEW_MARGIN - PREVIEW_PADDING * 2;
    int contentHeight = Dock::APP_PREVIEW_HEIGHT + Dock::APP_PREVIEW_MARGIN*2 - PREVIEW_PADDING * 2;

    return QSize(contentWidth, contentHeight);
}

AppPreviews::~AppPreviews()
{

}





