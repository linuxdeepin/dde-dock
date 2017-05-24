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

    connect(m_closeBtn, &DImageButton::clicked, this, &FloatingPreview::onCloseBtnClicked);
}

void FloatingPreview::trackWindow(AppSnapshot * const snap)
{
    if (!m_tracked.isNull())
        m_tracked->removeEventFilter(this);
    snap->installEventFilter(this);
    m_tracked = snap;

    const QRect r = rect();
    const QRect sr = snap->geometry();
    const QPoint offset = sr.center() - r.center();

    emit requestMove(offset);
}

void FloatingPreview::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    if (m_tracked.isNull())
        return;

    QPainter painter(this);

    const QRect r = rect().marginsRemoved(QMargins(8, 8, 8, 8));
    const QImage snapshot = m_tracked->snapshot();
    const QImage im = snapshot.scaled(r.size(), Qt::KeepAspectRatio, Qt::SmoothTransformation);
    const QRect ir = im.rect();
    const QPoint offset = r.center() - ir.center();
    painter.fillRect(r, Qt::black);
    painter.drawImage(offset.x(), offset.y(), im);

    // bottom black background
    QRect bgr = r;
    bgr.setTop(bgr.bottom() - 25);
    painter.fillRect(bgr, QColor(0, 0, 0, 255 * 0.3));

    // bottom title
    painter.drawText(bgr, Qt::AlignCenter, m_tracked->title());

    // draw border
    const QRect br = r.marginsAdded(QMargins(1, 1, 1, 1));
    painter.setBrush(Qt::transparent);
    painter.setRenderHint(QPainter::Antialiasing);
    painter.drawRoundedRect(br, 3, 3);
}

void FloatingPreview::mouseReleaseEvent(QMouseEvent *e)
{
    QWidget::mouseReleaseEvent(e);

    emit m_tracked->clicked(m_tracked->wid());
}

bool FloatingPreview::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_tracked && event->type() == QEvent::Destroy)
        hide();

    return QWidget::eventFilter(watched, event);
}

void FloatingPreview::onCloseBtnClicked()
{
    Q_ASSERT(!m_tracked.isNull());

    m_tracked->closeWindow();
}
