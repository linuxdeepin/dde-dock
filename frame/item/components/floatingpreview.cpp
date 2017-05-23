#include "floatingpreview.h"
#include "appsnapshot.h"
#include "_previewcontainer.h"

#include <QPainter>
#include <QVBoxLayout>

FloatingPreview::FloatingPreview(QWidget *parent)
    : QWidget(parent),

      m_closeBtn(new DImageButton)
{
    m_closeBtn->setFixedSize(24, 24);
    m_closeBtn->setNormalPic(":/icons/resources/close_round_normal.png");
    m_closeBtn->setHoverPic(":/icons/resources/close_round_hover.png");
    m_closeBtn->setPressPic(":/icons/resources/close_round_press.png");

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addWidget(m_closeBtn);
    centralLayout->setAlignment(m_closeBtn, Qt::AlignRight | Qt::AlignTop);
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setLayout(centralLayout);
    setFixedSize(SNAP_WIDTH, SNAP_HEIGHT);
}

void FloatingPreview::trackWindow(AppSnapshot * const snap)
{
    m_tracked = snap;

    const QRect r = rect();
    const QRect sr = snap->geometry();
    const QPoint offset = sr.center() - r.center();

    emit requestMove(offset);
}

void FloatingPreview::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);

    painter.fillRect(rect(), Qt::red);
}
