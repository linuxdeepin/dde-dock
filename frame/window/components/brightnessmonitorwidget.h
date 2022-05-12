#ifndef BRIGHTNESSMONITORWIDGET_H
#define BRIGHTNESSMONITORWIDGET_H

#include <QMap>
#include <QWidget>

class QLabel;
class CustomSlider;
class BrightnessModel;
class QStandardItemModel;
class QVBoxLayout;
class SliderContainer;
class BrightMonitor;
class SettingDelegate;

namespace Dtk { namespace Widget { class DListView; } }

using namespace Dtk::Widget;

class BrightnessMonitorWidget : public QWidget
{
    Q_OBJECT

public:
    explicit BrightnessMonitorWidget(BrightnessModel *model, QWidget *parent = nullptr);
    ~BrightnessMonitorWidget() override;

private:
    void initUi();
    void initConnection();
    void reloadMonitor();

    void resetHeight();

private Q_SLOTS:
    void onBrightChanged(BrightMonitor *monitor);

private:
    QWidget *m_sliderWidget;
    QVBoxLayout *m_sliderLayout;
    QList<QPair<BrightMonitor *, SliderContainer *>> m_sliderContainers;
    QLabel *m_descriptionLabel;
    DListView *m_deviceList;
    BrightnessModel *m_brightModel;
    QStandardItemModel *m_model;
    SettingDelegate *m_delegate;
};

#endif // BRIGHTNESSMONITORWIDGET_H
