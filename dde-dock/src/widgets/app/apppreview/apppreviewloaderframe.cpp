#include <QLabel>
#include <QPushButton>
#include <QHBoxLayout>
#include <QEvent>
#include <QMouseEvent>
#include <QDebug>

#include "apppreviewloader.h"
#include "apppreviewloaderframe.h"
#include "interfaces/dockconstants.h"

const int BUTTON_SIZE = Dock::APP_PREVIEW_CLOSEBUTTON_SIZE;
const int TITLE_HEIGHT = 25;
const int BORDER_WIDTH = 2;

////////////////////////////////////////
///                             \\// ///
///////////////////////////////////  ///
///                             //\\ ///
///                            ///   ///
///                            ///   ///
///                            ///   ///
///    AppPreviewLoader        ///   ///
///                            ///   ///
///                            ///   ///
///                            ///   ///
///                            ///   ///
//////////////////////////////////   ///
///                                  ///
/// ////////////////////////////////////
PopupFrame::PopupFrame(QWidget *parent) :
    QWidget(parent)
{
    //PopupFrame is a middleware which use for show popup style preview
}

void PopupFrame::mousePressEvent(QMouseEvent *)
{
    emit mousePress();
}

void PopupFrame::leaveEvent(QEvent *)
{
    emit mouseLeave();
}


AppPreviewLoaderFrame::AppPreviewLoaderFrame(const QString &title, int xid, QWidget *parent) :
    QFrame(parent), m_parent(parent), m_inMiniStyle(false),m_canShowTitle(true), m_xid(xid)
{
    this->setFixedSize(Dock::APP_PREVIEW_WIDTH + BUTTON_SIZE / 2, Dock::APP_PREVIEW_HEIGHT + BUTTON_SIZE);

    initPopupWidget();
    initPreviewLoader(xid);
    initCloseButton();
    initTitle(title);
}

AppPreviewLoaderFrame::~AppPreviewLoaderFrame()
{

}

void AppPreviewLoaderFrame::shrink(const QSize &size, bool miniStyle)
{
    this->setFixedSize(size);
    m_inMiniStyle = miniStyle;
    m_canShowTitle = !m_inMiniStyle;

    updatePopWidgetGeometry();
    updateWidgetsGeometry();
}

void AppPreviewLoaderFrame::enterEvent(QEvent *)
{
    QSize ts(this->width() * 1.4, this->height() * 1.4);
    if (m_inMiniStyle && (ts.height() < m_parent->height())) {
        m_popupWidget->setParent(m_parent);
        m_popupWidget->setFixedSize(ts);
        QPoint tp = this->mapToParent(QPoint(0, 0));
        int nx = tp.x() - (m_popupWidget->width() - this->width()) / 2;
        if (nx < 0)
            nx = 0;
        else if (nx + ts.width() > m_parent->width())
            nx = m_parent->width() - ts.width();
        m_popupWidget->move(nx, (m_parent->height() - m_popupWidget->height()) / 2);
        m_popupWidget->show();

        m_canShowTitle = true;
    }

    m_closeButton->show();
    m_previewLoader->setIsHover(true);
    m_previewLoader->requestUpdate();

    updateWidgetsGeometry();
}

void AppPreviewLoaderFrame::initPopupWidget()
{
    m_popupWidget = new PopupFrame(this);
    //TODO add box-shadow

    connect(m_popupWidget, &PopupFrame::mousePress, [=] {
        m_popupWidget->setParent(this); //make sure the popupwidget will be delete with this
        emit requestPreviewActive(m_xid);
    });
    connect(m_popupWidget, &PopupFrame::mouseLeave, [=] {
        m_popupWidget->setFixedSize(this->size());
        m_popupWidget->setParent(this);
        m_layout->addWidget(m_popupWidget);

        m_closeButton->hide();
        m_previewLoader->setIsHover(false);
        m_canShowTitle = !m_inMiniStyle;

        updateWidgetsGeometry();
    });

    m_layout = new QHBoxLayout(this);
    m_layout->setContentsMargins(0, 0, 0, 0);
    m_layout->addWidget(m_popupWidget);

    updatePopWidgetGeometry();
}

void AppPreviewLoaderFrame::initTitle(const QString &t)
{
    m_titleLabel = new QLabel(t, m_popupWidget);
    m_titleLabel->setAlignment(Qt::AlignCenter);
    m_titleLabel->setObjectName("AppPreviewLoaderFrameTitle");

    updateTitleGeometry();
}

void AppPreviewLoaderFrame::initPreviewLoader(int xid)
{
    m_previewLoader = new AppPreviewLoader(xid, m_popupWidget);

    updatePreviewLoaderGeometry();
}

void AppPreviewLoaderFrame::initCloseButton()
{
    m_closeButton = new QPushButton(m_popupWidget);
    m_closeButton->setObjectName("AppPreviewLoaderFrameCloseButton");
    m_closeButton->setFixedSize(BUTTON_SIZE, BUTTON_SIZE);
    m_closeButton->hide();

    connect(m_closeButton, &QPushButton::clicked, [=]{
        m_popupWidget->setParent(this); //make sure the popupwidget will be delete with this
        emit requestPreviewClose(m_xid);
        this->deleteLater();
    });

    updateCloseButtonGeometry();
}

void AppPreviewLoaderFrame::updatePopWidgetGeometry()
{
    m_popupWidget->setFixedSize(this->size());
}

void AppPreviewLoaderFrame::updatePreviewLoaderGeometry()
{
    //left:parent.left
    //leftMargin:0
    //topMargin==rightMargin==bottomMargin==BUTTON_SIZE / 2
    //horizontalCenter:parent.horizontalCenter
    m_previewLoader->setFixedSize(m_popupWidget->width() - BUTTON_SIZE / 2, m_popupWidget->height() - BUTTON_SIZE);
    m_previewLoader->move(0, BUTTON_SIZE / 2);
}

void AppPreviewLoaderFrame::updateCloseButtonGeometry()
{
    //always in the top-right corner
    m_closeButton->move(m_popupWidget->width() - BUTTON_SIZE, 0);
}

void AppPreviewLoaderFrame::updateTitleGeometry()
{
    if (m_canShowTitle) {
        m_titleLabel->setVisible(true);
        QFontMetrics fm(m_titleLabel->font());
        m_titleLabel->setText(fm.elidedText(m_titleLabel->text(), Qt::ElideRight, width() * 4 / 5));

        m_titleLabel->setFixedSize(m_previewLoader->width() - BORDER_WIDTH * 2, TITLE_HEIGHT);
        m_titleLabel->move(BORDER_WIDTH, BUTTON_SIZE / 2 + m_previewLoader->height() - TITLE_HEIGHT - BORDER_WIDTH);
    }
    else
        m_titleLabel->setVisible(false);
}

void AppPreviewLoaderFrame::updateWidgetsGeometry()
{
    updatePreviewLoaderGeometry();
    updateCloseButtonGeometry();
    updateTitleGeometry();
}


