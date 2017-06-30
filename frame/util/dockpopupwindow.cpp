#include "dockpopupwindow.h"

#include <QScreen>
#include <QApplication>
#include <QDesktopWidget>

DWIDGET_USE_NAMESPACE

const int MOUSE_BUTTON(1 << 1);

DockPopupWindow::DockPopupWindow(QWidget *parent)
    : DArrowRectangle(ArrowBottom, parent),
      m_model(false),

      m_acceptDelayTimer(new QTimer(this)),

      m_mouseInter(new DBusXMouseArea(this)),
      m_displayInter(new DBusDisplay(this))
{
    m_acceptDelayTimer->setSingleShot(true);
    m_acceptDelayTimer->setInterval(100);

    m_wmHelper = DWindowManagerHelper::instance();

    compositeChanged();

    setBackgroundColor(DBlurEffectWidget::DarkColor);
    setWindowFlags(Qt::X11BypassWindowManagerHint | Qt::WindowStaysOnTopHint);
    setAttribute(Qt::WA_InputMethodEnabled, false);
    setFocusPolicy(Qt::StrongFocus);

    connect(m_acceptDelayTimer, &QTimer::timeout, this, &DockPopupWindow::accept);
    connect(m_wmHelper, &DWindowManagerHelper::hasCompositeChanged, this, &DockPopupWindow::compositeChanged);
}

DockPopupWindow::~DockPopupWindow()
{
}

bool DockPopupWindow::model() const
{
    return m_model;
}

void DockPopupWindow::setContent(QWidget *content)
{
    QWidget *lastWidget = getContent();
    if (lastWidget)
        lastWidget->removeEventFilter(this);
    content->installEventFilter(this);

    setAccessibleName(content->objectName() + "-popup");

    DArrowRectangle::setContent(content);
}

void DockPopupWindow::show(const QPoint &pos, const bool model)
{
    m_model = model;
    m_lastPoint = pos;

    show(pos.x(), pos.y());

    if (!model && !m_mouseAreaKey.isEmpty())
        unRegisterMouseEvent();

    if (model && m_mouseAreaKey.isEmpty())
        registerMouseEvent();
}

void DockPopupWindow::show(const int x, const int y)
{
    m_lastPoint = QPoint(x, y);

    DArrowRectangle::show(x, y);
}

void DockPopupWindow::hide()
{
    if (!m_mouseAreaKey.isEmpty())
        unRegisterMouseEvent();

    DArrowRectangle::hide();
}

void DockPopupWindow::showEvent(QShowEvent *e)
{
    DArrowRectangle::showEvent(e);

    QTimer::singleShot(1, this, [&] {
        raise();
        if (!m_model)
            return;
        activateWindow();
        setFocus(Qt::ActiveWindowFocusReason);
    });
}

void DockPopupWindow::enterEvent(QEvent *e)
{
    DArrowRectangle::enterEvent(e);

    raise();
    if (!m_model)
        return;
    activateWindow();
    setFocus(Qt::ActiveWindowFocusReason);
}

void DockPopupWindow::mousePressEvent(QMouseEvent *e)
{
    DArrowRectangle::mousePressEvent(e);

//    if (e->button() == Qt::LeftButton)
//        m_acceptDelayTimer->start();
}

bool DockPopupWindow::eventFilter(QObject *o, QEvent *e)
{
    if (o != getContent() || e->type() != QEvent::Resize)
        return false;

    // FIXME: ensure position move after global mouse release event
    QTimer::singleShot(100, this, [this] {if (isVisible()) show(m_lastPoint, m_model);});

    return false;
}

void DockPopupWindow::globalMouseRelease(int button, int x, int y, const QString &id)
{
    // button_left
    if (button != 1)
        return;

    if (id != m_mouseAreaKey)
        return;

    Q_ASSERT(m_model);

    const QRect rect = QRect(pos(), size());
    const QPoint pos = QPoint(x, y);

    if (rect.contains(pos))
        return;

    emit accept();

    unRegisterMouseEvent();
}

void DockPopupWindow::registerMouseEvent()
{
    if (!m_mouseAreaKey.isEmpty())
        return;

    // only regist mouse button event
    m_mouseAreaKey = m_mouseInter->RegisterArea(0, 0, m_displayInter->screenWidth(), m_displayInter->screenHeight(), MOUSE_BUTTON);

    connect(m_mouseInter, &DBusXMouseArea::ButtonRelease, this, &DockPopupWindow::globalMouseRelease, Qt::QueuedConnection);
}

void DockPopupWindow::unRegisterMouseEvent()
{
    if (m_mouseAreaKey.isEmpty())
        return;

    disconnect(m_mouseInter, &DBusXMouseArea::ButtonRelease, this, &DockPopupWindow::globalMouseRelease);

    m_mouseInter->UnregisterArea(m_mouseAreaKey);
    m_mouseAreaKey.clear();
}

void DockPopupWindow::compositeChanged()
{
    if (m_wmHelper->hasComposite())
        setBorderColor(QColor(255, 255, 255, 255 * 0.05));
    else
        setBorderColor(QColor("#2C3238"));
}
