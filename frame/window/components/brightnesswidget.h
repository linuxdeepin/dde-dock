#ifndef BRIGHTNESSWIDGET_H
#define BRIGHTNESSWIDGET_H

#include <DBlurEffectWidget>

DWIDGET_USE_NAMESPACE

class CustomSlider;
class BrightnessModel;
class BrightMonitor;

class BrightnessWidget : public DBlurEffectWidget
{
    Q_OBJECT

public:
    explicit BrightnessWidget(QWidget *parent = nullptr);
    ~BrightnessWidget() override;
    BrightnessModel *model();

Q_SIGNALS:
    void visibleChanged(bool);
    void rightIconClicked();

protected:
    void showEvent(QShowEvent *event) override;
    void hideEvent(QHideEvent *event) override;

private Q_SLOTS:
    void onUpdateBright();

private:
    void initUi();
    void initConenction();

private:
    CustomSlider *m_slider;
    BrightnessModel *m_model;
};

#endif // LIGHTSETTINGWIDGET_H
