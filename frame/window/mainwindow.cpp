#include "mainwindow.h"
#include "panel/mainpanel.h"

#include <QDebug>
#include <QResizeEvent>

MainWindow::MainWindow(QWidget *parent)
    : QWidget(parent),
      m_mainPanel(new MainPanel(this)),
      m_positionUpdateTimer(new QTimer(this)),
      m_sizeChangeAni(new QPropertyAnimation(this, "size")),
      m_posChangeAni(new QPropertyAnimation(this, "pos")),
      m_panelShowAni(new QPropertyAnimation(m_mainPanel, "pos")),
      m_panelHideAni(new QPropertyAnimation(m_mainPanel, "pos")),
      m_xcbMisc(XcbMisc::instance())
{
    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowDoesNotAcceptFocus);
    setAttribute(Qt::WA_TranslucentBackground);

    m_settings = new DockSettings(this);
    m_xcbMisc->set_window_type(winId(), XcbMisc::Dock);

    initComponents();
    initConnections();

    m_mainPanel->setFixedSize(m_settings->windowSize());

    updatePanelVisible();

//    setStyleSheet("background-color:red;");
}

MainWindow::~MainWindow()
{
    delete m_xcbMisc;
}

void MainWindow::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);
}

void MainWindow::mousePressEvent(QMouseEvent *e)
{
    e->ignore();

    if (e->button() == Qt::RightButton)
        m_settings->showDockSettingsMenu();
}

void MainWindow::keyPressEvent(QKeyEvent *e)
{
    switch (e->key())
    {
#ifdef QT_DEBUG
    case Qt::Key_Escape:        qApp->quit();       break;
#endif
    default:;
    }
}

void MainWindow::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);

    if (m_settings->hideState() != Show)
        expand();
}

void MainWindow::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);

    updatePanelVisible();
}

void MainWindow::setFixedSize(const QSize &size)
{
    const QPropertyAnimation::State state = m_posChangeAni->state();

    if (state == QPropertyAnimation::Stopped && this->size() == size)
        return;

    if (state == QPropertyAnimation::Running)
        return m_sizeChangeAni->setEndValue(size);

    m_sizeChangeAni->setStartValue(this->size());
    m_sizeChangeAni->setEndValue(size);
    m_sizeChangeAni->start();
}

void MainWindow::move(int x, int y)
{
    const QPropertyAnimation::State state = m_posChangeAni->state();

    if (state == QPropertyAnimation::Stopped && this->pos() == QPoint(x, y))
        return;

    if (state == QPropertyAnimation::Running)
        return m_posChangeAni->setEndValue(QPoint(x, y));

    m_posChangeAni->setStartValue(pos());
    m_posChangeAni->setEndValue(QPoint(x, y));
    m_posChangeAni->start();
}

void MainWindow::initComponents()
{
    m_positionUpdateTimer->setSingleShot(true);
    m_positionUpdateTimer->setInterval(20);
    m_positionUpdateTimer->start();

    m_sizeChangeAni->setDuration(200);
    m_sizeChangeAni->setEasingCurve(QEasingCurve::InOutCubic);

    m_posChangeAni->setDuration(200);
    m_posChangeAni->setEasingCurve(QEasingCurve::InOutCubic);

    m_panelShowAni->setDuration(200);
    m_panelShowAni->setEasingCurve(QEasingCurve::InOutCubic);

    m_panelHideAni->setDuration(200);
    m_panelHideAni->setEasingCurve(QEasingCurve::InOutCubic);
}

void MainWindow::initConnections()
{
    connect(m_settings, &DockSettings::dataChanged, m_positionUpdateTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::windowGeometryChanged, this, &MainWindow::updateGeometry, Qt::DirectConnection);
    connect(m_settings, &DockSettings::windowHideModeChanged, this, &MainWindow::setStrutPartial, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::windowVisibleChanegd, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_settings, &DockSettings::autoHideChanged, this, &MainWindow::updatePanelVisible);

    connect(m_mainPanel, &MainPanel::requestRefershWindowVisible, this, &MainWindow::updatePanelVisible, Qt::QueuedConnection);
    connect(m_mainPanel, &MainPanel::requestWindowAutoHide, m_settings, &DockSettings::setAutoHide);

    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition, Qt::QueuedConnection);

    connect(m_panelHideAni, &QPropertyAnimation::finished, this, &MainWindow::updateGeometry, Qt::QueuedConnection);

    // to fix qt animation bug, sometimes window size not change
    connect(m_sizeChangeAni, &QPropertyAnimation::valueChanged, [this] {
        const QSize size = m_sizeChangeAni->currentValue().toSize();

        QWidget::setFixedSize(size);
        m_mainPanel->setFixedSize(size);
    });
}

void MainWindow::updatePosition()
{
    // all update operation need pass by timer
    Q_ASSERT(sender() == m_positionUpdateTimer);

    clearStrutPartial();
    updateGeometry();
    setStrutPartial();
}

void MainWindow::updateGeometry()
{
    const Position position = m_settings->position();
    QSize size = m_settings->windowSize();

    m_mainPanel->setFixedSize(size);
    m_mainPanel->updateDockPosition(position);
    m_mainPanel->updateDockDisplayMode(m_settings->displayMode());

    if (m_settings->hideState() == Hide)
    {
        m_sizeChangeAni->stop();
        m_posChangeAni->stop();
        switch (position)
        {
        case Top:
        case Bottom:    size.setHeight(1);      break;
        case Left:
        case Right:     size.setWidth(1);       break;
        }
        QWidget::setFixedSize(size);
    }
    else
    {
        setFixedSize(size);
    }

    const QRect primaryRect = m_settings->primaryRect();
    const int offsetX = (primaryRect.width() - size.width()) / 2;
    const int offsetY = (primaryRect.height() - size.height()) / 2;

    switch (position)
    {
    case Top:
        move(primaryRect.topLeft().x() + offsetX, primaryRect.y());                   break;
    case Left:
        move(primaryRect.topLeft().x(), primaryRect.y() + offsetY);                   break;
    case Right:
        move(primaryRect.right() - size.width() + 1, primaryRect.y() + offsetY);      break;
    case Bottom:
        move(primaryRect.x() + offsetX, primaryRect.bottom() - size.height() + 1);    break;
    default:
        Q_ASSERT(false);
    }

    m_mainPanel->update();
}

void MainWindow::clearStrutPartial()
{
    m_xcbMisc->clear_strut_partial(winId());
}

void MainWindow::setStrutPartial()
{
    // first, clear old strut partial
    clearStrutPartial();

    if (m_settings->hideMode() != Dock::KeepShowing)
        return;

    const Position side = m_settings->position();
    const int maxScreenHeight = m_settings->screenHeight();

    XcbMisc::Orientation orientation = XcbMisc::OrientationTop;
    uint strut = 0;
    uint strutStart = 0;
    uint strutEnd = 0;

    const QPoint p = m_posChangeAni->endValue().toPoint();
    const QRect r = QRect(p, m_settings->windowSize());
    switch (side)
    {
    case Position::Top:
        orientation = XcbMisc::OrientationTop;
        strut = r.bottom() + 1;
        strutStart = r.left();
        strutEnd = r.right();
        break;
    case Position::Bottom:
        orientation = XcbMisc::OrientationBottom;
        strut = maxScreenHeight - p.y();
        strutStart = r.left();
        strutEnd = r.right();
        break;
    case Position::Left:
        orientation = XcbMisc::OrientationLeft;
        strut = r.width();
        strutStart = r.top();
        strutEnd = r.bottom();
        break;
    case Position::Right:
        orientation = XcbMisc::OrientationRight;
        strut = r.width();
        strutStart = r.top();
        strutEnd = r.bottom();
        break;
    default:
        Q_ASSERT(false);
    }

    m_xcbMisc->set_strut_partial(winId(), orientation, strut, strutStart, strutEnd);
}

void MainWindow::expand()
{
//    qDebug() << "expand";
    const QPoint finishPos(0, 0);

    if (m_mainPanel->pos() == finishPos && m_mainPanel->size() == this->size())
        return;

    resetPanelEnvironment();

    if (m_panelShowAni->state() == QPropertyAnimation::Running)
        return m_panelShowAni->setEndValue(finishPos);

    const QSize size = m_settings->windowSize();

    QPoint startPos(0, 0);
    switch (m_settings->position())
    {
    case Top:       startPos.setY(-size.height());     break;
    case Bottom:    startPos.setY(size.height());      break;
    case Left:      startPos.setX(-size.width());      break;
    case Right:     startPos.setX(size.width());       break;
    }

    m_panelHideAni->stop();
    m_panelShowAni->setStartValue(startPos);
    m_panelShowAni->setEndValue(finishPos);
    m_panelShowAni->start();
}

void MainWindow::narrow()
{
//    qDebug() << "narrow";
    const QSize size = m_settings->windowSize();

    QPoint finishPos(0, 0);
    switch (m_settings->position())
    {
    case Top:       finishPos.setY(-size.height());     break;
    case Bottom:    finishPos.setY(size.height());      break;
    case Left:      finishPos.setX(-size.width());      break;
    case Right:     finishPos.setX(size.width());       break;
    }

    if (m_panelHideAni->state() == QPropertyAnimation::Running)
        return m_panelHideAni->setEndValue(finishPos);

    m_panelShowAni->stop();
    m_panelHideAni->setStartValue(m_mainPanel->pos());
    m_panelHideAni->setEndValue(finishPos);
    m_panelHideAni->start();
}

void MainWindow::resetPanelEnvironment()
{
    // reset environment
    m_sizeChangeAni->stop();
    m_posChangeAni->stop();
    const QSize size = m_settings->windowSize();
    const QRect primaryRect = m_settings->primaryRect();
    const int offsetX = (primaryRect.width() - size.width()) / 2;
    const int offsetY = (primaryRect.height() - size.height()) / 2;

    QWidget::setFixedSize(size);
    m_mainPanel->setFixedSize(size);
    m_mainPanel->move(0, 0);
    switch (m_settings->position())
    {
    case Top:
        QWidget::move(primaryRect.topLeft().x() + offsetX, 0);               break;
    case Left:
        QWidget::move(primaryRect.topLeft().x(), offsetY);                   break;
    case Right:
        QWidget::move(primaryRect.right() - size.width() + 1, offsetY);      break;
    case Bottom:
        QWidget::move(offsetX, primaryRect.bottom() - size.height() + 1);    break;
    default:
        Q_ASSERT(false);
    }
}

void MainWindow::updatePanelVisible()
{
    if (m_settings->hideMode() == KeepShowing)
        return;

    const Dock::HideState state = m_settings->hideState();

//    qDebug() << state;

    do
    {
        if (state != Hide)
            break;

        if (!m_settings->autoHide())
            break;

        QRect r(pos(), size());
        if (r.contains(QCursor::pos()))
            break;

        return narrow();

    } while (false);

    return expand();
}
