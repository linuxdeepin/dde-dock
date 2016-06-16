#include "mainwindow.h"
#include "panel/mainpanel.h"

#include <QDebug>
#include <QResizeEvent>

MainWindow::MainWindow(QWidget *parent)
    : QWidget(parent),
      m_mainPanel(new MainPanel(this)),

      m_settings(new DockSettings(this)),
      m_displayInter(new DBusDisplay(this)),

      m_xcbMisc(XcbMisc::instance()),

      m_positionUpdateTimer(new QTimer(this))
{
    setWindowFlags(Qt::FramelessWindowHint | Qt::WindowDoesNotAcceptFocus);
    setAttribute(Qt::WA_TranslucentBackground);

    m_xcbMisc->set_window_type(winId(), XcbMisc::Dock);

    initComponents();
    initConnections();
}

MainWindow::~MainWindow()
{
    delete m_xcbMisc;
}

void MainWindow::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    m_mainPanel->setFixedSize(e->size());
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

void MainWindow::initComponents()
{
    m_positionUpdateTimer->setSingleShot(true);
    m_positionUpdateTimer->setInterval(200);
    m_positionUpdateTimer->start();
}

void MainWindow::initConnections()
{
    connect(m_displayInter, &DBusDisplay::PrimaryRectChanged, [this] {m_positionUpdateTimer->start();});

    connect(m_positionUpdateTimer, &QTimer::timeout, this, &MainWindow::updatePosition);
}

void MainWindow::updatePosition()
{
    Q_ASSERT(sender() == m_positionUpdateTimer);

    clearStrutPartial();

    const QRect screenRect = m_displayInter->primaryRect();

    setFixedWidth(screenRect.width());
    setFixedHeight(60);

    move(0, screenRect.height() - 60);

    setStrutPartial();
}

void MainWindow::clearStrutPartial()
{
    m_xcbMisc->clear_strut_partial(winId());
}

void MainWindow::setStrutPartial()
{
    // first, clear old strut partial
    clearStrutPartial();

    const DockSide side = m_settings->side();
    const int maxScreenHeight = m_displayInter->screenHeight();

    XcbMisc::Orientation orientation;
    uint strut;
    uint strutStart;
    uint strutEnd;

    const QPoint p = pos();
    const QRect r = rect();
    switch (side)
    {
    case DockSide::Top:
        orientation = XcbMisc::OrientationTop;
        strut = r.bottom();
        strutStart = r.left();
        strutEnd = r.right();
        break;
    case DockSide::Bottom:
        orientation = XcbMisc::OrientationBottom;
        strut = maxScreenHeight - p.y();
        strutStart = r.left();
        strutEnd = r.right();
        break;
    case DockSide::Left:
        orientation = XcbMisc::OrientationLeft;
        strut = r.width();
        strutStart = r.top();
        strutEnd = r.bottom();
        break;
    case DockSide::Right:
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
