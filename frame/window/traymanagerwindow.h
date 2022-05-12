#ifndef TRAYMANAGERWINDOW_H
#define TRAYMANAGERWINDOW_H

#include "constants.h"

#include <QWidget>

#include <com_deepin_daemon_timedate.h>

namespace Dtk { namespace Gui { class DRegionMonitor; };
                namespace Widget { class DBlurEffectWidget; } }

using namespace Dtk::Widget;

using Timedate = com::deepin::daemon::Timedate;

class QuickPluginWindow;
class QBoxLayout;
class TrayGridView;
class TrayModel;
class SystemPluginWindow;
class QLabel;
class QDropEvent;
class DateTimeDisplayer;

class TrayManagerWindow : public QWidget
{
    Q_OBJECT

Q_SIGNALS:
    void sizeChanged();

public:
    explicit TrayManagerWindow(QWidget *parent = nullptr);
    ~TrayManagerWindow() override;
    void setPositon(Dock::Position position);
    QSize suitableSize();

protected:
    void resizeEvent(QResizeEvent *event) override;

private:
    void initUi();
    void initConnection();

    void resetChildWidgetSize();
    void resetMultiDirection();
    void resetSingleDirection();

    void dragEnterEvent(QDragEnterEvent *e) override;
    void dragMoveEvent(QDragMoveEvent *e) override;
    void dropEvent(QDropEvent *e) override;
    void dragLeaveEvent(QDragLeaveEvent *event) override;

    bool showSingleRow();
    int appDatetimeSize();

private:
    DBlurEffectWidget *m_appPluginDatetimeWidget;
    SystemPluginWindow *m_systemPluginWidget;
    QWidget *m_appPluginWidget;
    QuickPluginWindow *m_quickIconWidget;
    DateTimeDisplayer *m_dateTimeWidget;
    QBoxLayout *m_appPluginLayout;
    QBoxLayout *m_appDatetimeLayout;
    QBoxLayout *m_mainLayout;
    TrayGridView *m_trayView;
    TrayModel *m_model;
    Dock::Position m_postion;
};

#endif // PLUGINWINDOW_H
