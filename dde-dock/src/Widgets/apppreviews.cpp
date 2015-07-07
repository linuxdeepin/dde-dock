#include "apppreviews.h"

AppPreviews::AppPreviews(QWidget *parent) : QWidget(parent)
{
    m_mainLayout = new QHBoxLayout(this);
    setLayout(m_mainLayout);
    resize(Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);
}

void AppPreviews::addItem(const QString &title, int xid)
{
    if (m_xidList.indexOf(xid) != -1)
        return;
    m_xidList.append(xid);
    WindowPreview * preview = new WindowPreview(xid);
    preview->resize(Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);
    m_mainLayout->addWidget(preview);
    resize(m_mainLayout->count() * Dock::APP_PREVIEW_WIDTH,Dock::APP_PREVIEW_HEIGHT);
}

void AppPreviews::setTitle(const QString &title)
{
    QLabel *titleLabel = new QLabel(title);
    titleLabel->setObjectName("DockAppTitle");
    titleLabel->setAlignment(Qt::AlignCenter);
    m_mainLayout->addWidget(titleLabel);
    resize(100,35);
}

