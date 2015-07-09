#include "apppreviews.h"

AppPreviewFrame::AppPreviewFrame(QWidget *preview, const QString &title, int xid, QWidget *parent) : QWidget(parent),xidValue(xid)
{
    addPreview(preview);
    setTitle(title);
    addCloseButton();
}

void AppPreviewFrame::addPreview(QWidget *p)
{
    this->resize(p->size());
    p->setParent(this);
    p->move(0,0);
}

void AppPreviewFrame::setTitle(const QString &t)
{
    QLabel *titleLabel = new QLabel(this);
    QFontMetrics fm(titleLabel->font());
    titleLabel->setText(fm.elidedText(t,Qt::ElideRight,width() * 4 / 5));
    titleLabel->setStyleSheet("color:white");
    titleLabel->setAlignment(Qt::AlignVCenter | Qt::AlignLeft);
    titleLabel->resize(width() * 4 / 5,20);
    titleLabel->move(width() / 5 / 2,height() - titleLabel->height());
}

void AppPreviewFrame::mousePressEvent(QMouseEvent *)
{
    emit activate(xidValue);
}

void AppPreviewFrame::addCloseButton()
{
    CloseButton *cb = new CloseButton(this);
    connect(cb,&CloseButton::clicked,[=](){close(this->xidValue);});
    cb->resize(28,28);
    cb->move(width() - cb->width(),0);
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
    m_xidList.append(xid);
    WindowPreview * preview = new WindowPreview(xid);
//    QWidget *preview = new QWidget();
    preview->resize(Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);
    AppPreviewFrame *f = new AppPreviewFrame(preview,title,xid);
    connect(f,&AppPreviewFrame::close,this,&AppPreviews::removePreview);
    connect(f,&AppPreviewFrame::activate,this,&AppPreviews::activatePreview);
    m_mainLayout->addWidget(f);

    resize(m_mainLayout->count() * Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);
}

void AppPreviews::setTitle(const QString &title)
{
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

    resize(m_mainLayout->count() * Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);
    emit sizeChanged();
}

void AppPreviews::activatePreview(int xid)
{
    m_clientManager->ActiveWindow(xid);
}







