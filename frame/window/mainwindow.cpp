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
    if (m_sizeChangeAni->state() == QPropertyAnimation::Running)
        return m_sizeChangeAni->setEndValue(size);

    m_sizeChangeAni->setStartValue(this->size());
    m_sizeChangeAni->setEndValue(size);
    m_sizeChangeAni->start();
}

void MainWindow::move(int x, int y)
{
    if (m_posChangeAni->state() == QPropertyAnimation::Running)
        return m_posChangeAni->setEndValue(QPoint(x, y));

    m_posChangeAni->setStartValue(pos());
    m_posChangeAni->setEndValue(QPoint(x, y));
    m_posChangeAni->start();
}

void MainWindow::initComponents()
{
    m_positionUpdateTimer->setSingleShot(true);
    m_positionUpdateTimer->setInterval(200);
    m_positionUpdateTimer->start();

    m_sizeChangeAni->setDuration(200);
    m_sizeChangeAni->setEasingCurve(QEasingCurve::OutCubic);

    m_posChangeAni->setDuration(200);
    m_posChangeAni->setEasingCurve(QEasingCurve::OutCubic);
}

void MainWindow::initConnections()
{
    connect(m_settings, &DockSettings::dataChanged, m_positionUpdateTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_settings, &DockSettings::windowGeometryChanged, this, &MainWindow::updateGeometry);
    connect(m_settings, &DockSettings::windowHideModeChanged, this, &MainWindow::setStrutPartial);
    connect(m_settings, &DockSettings::windowVisibleChanegd, this, &MainWindow::updatePanelVisible);

    connect(m_panelHideAni, &QPropertyAnimation::finished, this, &MainWindow::updateGeometry);

    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition);
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

    m_mainPanel->setFixedSize(m_settings->windowSize());
    m_mainPanel->updateDockPosition(position);
    m_mainPanel->updateDockDisplayMode(m_settings->displayMode());

    QSize size = m_settings->windowSize();
    if (m_settings->hideState() == Hide)
    {
        m_sizeChangeAni->stop();
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
        move(primaryRect.topLeft().x() + offsetX, 0);               break;
    case Left:
        move(primaryRect.topLeft().x(), offsetY);                   break;
    case Right:
        move(primaryRect.right() - size.width() + 1, offsetY);      break;
    case Bottom:
        move(offsetX, primaryRect.bottom() - size.height() + 1);    break;
    default:
        Q_ASSERT(false);
    }

    update();
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
    if (m_mainPanel->pos() == QPoint(0, 0))
        return;

    m_sizeChangeAni->stop();
    m_posChangeAni->stop();
    const QSize size = m_settings->windowSize();
    const QRect primaryRect = m_settings->primaryRect();
    const int offsetX = (primaryRect.width() - size.width()) / 2;
    const int offsetY = (primaryRect.height() - size.height()) / 2;

    QWidget::setFixedSize(size);
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

    const QPoint finishPos(0, 0);
    if (m_panelShowAni->state() == QPropertyAnimation::Running)
        return m_panelShowAni->setEndValue(finishPos);

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
    const QSize size = m_mainPanel->size();

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

void MainWindow::updatePanelVisible()
{
    const Dock::HideState state = m_settings->hideState();

//    qDebug() << state;

    if (state == Unknown)
        return;

    if (state == Show)
        expand();
    else
        narrow();
}
