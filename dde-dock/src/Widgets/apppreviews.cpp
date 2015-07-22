#include "apppreviews.h"

AppPreviewFrame::AppPreviewFrame(const QString &title, int xid, QWidget *parent) :
    QWidget(parent),xidValue(xid)
{
    addPreview(xid);
    setTitle(title);
    addCloseButton();
}

void AppPreviewFrame::addPreview(int xid)
{
    WindowPreview * preview = new WindowPreview(xid, this);
    preview->setObjectName("AppPreview");
    preview->resize(Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);

    resize(preview->size());
}

void AppPreviewFrame::setTitle(const QString &t)
{
    QLabel *titleLabel = new QLabel(this);
    titleLabel->setObjectName("AppPreviewFrameTitle");
    QFontMetrics fm(titleLabel->font());
    titleLabel->setText(fm.elidedText(t,Qt::ElideRight,width()));
    titleLabel->setAlignment(Qt::AlignCenter);
    titleLabel->resize(width(),25);
    titleLabel->move(0,height() - titleLabel->height());
}

void AppPreviewFrame::mousePressEvent(QMouseEvent *)
{
    emit activate(xidValue);
}

void AppPreviewFrame::enterEvent(QEvent *)
{
    showCloseButton();
}

void AppPreviewFrame::leaveEvent(QEvent *)
{
    hideCloseButton();
}

void AppPreviewFrame::addCloseButton()
{
    m_cb = new CloseButton(this);
    connect(m_cb,&CloseButton::clicked,[=](){close(this->xidValue);});
    m_cb->resize(28,28);

    m_cb->move(width() - m_cb->width()/* / 2*/,0/*- m_cb->width() / 2*/);
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
    m_mainLayout->setMargin(0);
    setLayout(m_mainLayout);
    resize(Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);
}

void AppPreviews::addItem(const QString &title, int xid)
{
    if (m_xidList.indexOf(xid) != -1)
        return;
    m_mainLayout->setMargin(Dock::APP_PREVIEW_MARGIN);
    m_mainLayout->setSpacing(Dock::APP_PREVIEW_MARGIN);
    m_xidList.append(xid);

    AppPreviewFrame *f = new AppPreviewFrame(title, xid);
    connect(f,&AppPreviewFrame::close,this,&AppPreviews::removePreview);
    connect(f,&AppPreviewFrame::activate,this,&AppPreviews::activatePreview);

    m_mainLayout->addWidget(f);

    int contentWidth = m_mainLayout->count() * (f->width() + Dock::APP_PREVIEW_MARGIN) + Dock::APP_PREVIEW_MARGIN;
    resize(contentWidth,f->height() + Dock::APP_PREVIEW_MARGIN*2);
}

void AppPreviews::setTitle(const QString &title)
{
    m_mainLayout->setMargin(0);
    QLabel *titleLabel = new QLabel(title);
    titleLabel->setObjectName("DockAppTitle");
    titleLabel->setAlignment(Qt::AlignCenter);
    m_mainLayout->addWidget(titleLabel);
    QFontMetrics fm(titleLabel->font());
    int textWidth = fm.width(title);
    resize(textWidth < 80 ? 80 : textWidth,20);
}

void AppPreviews::enterEvent(QEvent *)
{
    emit mouseEntered();
}

void AppPreviews::leaveEvent(QEvent *)
{
    if (isClosing)
    {
        isClosing = false;
        return;
    }
    emit mouseExited();
}

void AppPreviews::removePreview(int xid)
{
    isClosing = true;
    m_mainLayout->removeWidget(qobject_cast<AppPreviewFrame *>(sender()));
    sender()->deleteLater();
    m_clientManager->CloseWindow(xid);
    if (m_mainLayout->count() <= 0)
    {
        emit mouseExited();
        return;
    }

    int contentWidth = m_mainLayout->count() * (Dock::APP_PREVIEW_WIDTH + Dock::APP_PREVIEW_MARGIN) + Dock::APP_PREVIEW_MARGIN;
    resize(contentWidth,Dock::APP_PREVIEW_HEIGHT + Dock::APP_PREVIEW_MARGIN*2);
    emit sizeChanged();
}

void AppPreviews::activatePreview(int xid)
{
    m_clientManager->ActiveWindow(xid);
}







