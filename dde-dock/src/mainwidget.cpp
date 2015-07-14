#include "mainwidget.h"

MainWidget::MainWidget(QWidget *parent)
    : QWidget(parent)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockModeData::instance()->getDockHeight());
    mainPanel = new Panel(this);
    mainPanel->resize(this->width(),this->height());
    mainPanel->move(0,0);

    this->setWindowFlags(Qt::WindowStaysOnTopHint | Qt::FramelessWindowHint);
    this->setAttribute(Qt::WA_TranslucentBackground);
    this->move(0,rec.height());

    connect(DockModeData::instance(), SIGNAL(dockModeChanged(Dock::DockMode,Dock::DockMode)),
            this, SLOT(slotDockModeChanged(Dock::DockMode,Dock::DockMode)));

    initHSManager();
    initState();
}

void MainWidget::slotDockModeChanged(Dock::DockMode newMode, Dock::DockMode oldMode)
{
    QRect rec = QApplication::desktop()->screenGeometry();
    this->resize(rec.width(),DockModeData::instance()->getDockHeight());

//    mainPanel->resize(this->width(),this->height());
//    mainPanel->move(0,0);
}

void MainWidget::hasShown()
{
    m_HSManager->SetState(1);
}

void MainWidget::hasHidden()
{
    m_HSManager->SetState(3);
}

void MainWidget::hideStateChanged(int value)
{
    if (value == 0)
    {
        emit startShow();
    }
    else if (value == 1)
    {
        emit startHide();
    }
}

void MainWidget::initHSManager()
{
    m_HSManager = new DBusHideStateManager(this);
    connect(m_HSManager,&DBusHideStateManager::ChangeState,this,&MainWidget::hideStateChanged);
}

void MainWidget::initState()
{
    QRect rec = QApplication::desktop()->screenGeometry();

    //TODO no need "- 100"
    QStateMachine * machine = new QStateMachine(this);
    QState * showState = new QState(machine);
    showState->assignProperty(this,"pos",QPoint(rec.x(),rec.height() - DockModeData::instance()->getDockHeight() - 100));
    QState * hideState = new QState(machine);
    hideState->assignProperty(this,"pos",QPoint(rec.x(),rec.height() - 100));

    machine->setInitialState(showState);

    QPropertyAnimation *sa = new QPropertyAnimation(this, "pos");
    sa->setDuration(200);
    sa->setEasingCurve(QEasingCurve::InSine);
    connect(sa,&QPropertyAnimation::finished,this,&MainWidget::hasShown);

    QPropertyAnimation *ha = new QPropertyAnimation(this, "pos");
    ha->setDuration(200);
    ha->setEasingCurve(QEasingCurve::InSine);
    connect(ha,&QPropertyAnimation::finished,this,&MainWidget::hasHidden);

    QSignalTransition *ts1 = showState->addTransition(this,SIGNAL(startHide()), hideState);
    ts1->addAnimation(ha);
    connect(ts1,&QSignalTransition::triggered,[=](int value = 2){
        m_HSManager->SetState(value);
    });
    QSignalTransition *ts2 = hideState->addTransition(this,SIGNAL(startShow()),showState);
    ts2->addAnimation(sa);
    connect(ts2,&QSignalTransition::triggered,[=](int value = 0){
        m_HSManager->SetState(value);
    });

    machine->start();
}

MainWidget::~MainWidget()
{

}
