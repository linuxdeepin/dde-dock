#include "brightnesswidget.h"
#include "customslider.h"
#include "brightnessmodel.h"
#include "brightnessmonitorwidget.h"

#include <QHBoxLayout>
#include <QDebug>

BrightnessWidget::BrightnessWidget(QWidget *parent)
    : DBlurEffectWidget(parent)
    , m_slider(new CustomSlider(CustomSlider::SliderType::Normal, this))
    , m_model(new BrightnessModel(this))
{
    initUi();
    initConenction();
    onUpdateBright();
}

BrightnessWidget::~BrightnessWidget()
{
}

BrightnessModel *BrightnessWidget::model()
{
    return m_model;
}

void BrightnessWidget::showEvent(QShowEvent *event)
{
    DBlurEffectWidget::showEvent(event);
    Q_EMIT visibleChanged(true);
}

void BrightnessWidget::hideEvent(QHideEvent *event)
{
    DBlurEffectWidget::hideEvent(event);
    Q_EMIT visibleChanged(true);
}

void BrightnessWidget::initUi()
{
    QHBoxLayout *layout = new QHBoxLayout(this);
    layout->setContentsMargins(20, 0, 20, 0);
    layout->addWidget(m_slider);

    m_slider->setPageStep(1);
    m_slider->setIconSize(QSize(24, 24));

    m_slider->setLeftIcon(QIcon(":/icons/resources/brightness.svg"));
    m_slider->setRightIcon(QIcon::fromTheme(":/icons/resources/ICON_Device_Laptop.svg"));
    m_slider->setTickPosition(QSlider::TicksBelow);
}

void BrightnessWidget::initConenction()
{
    connect(m_slider, &CustomSlider::iconClicked, this, [ this ](DSlider::SliderIcons icon, bool) {
        if (icon == DSlider::SliderIcons::RightIcon)
            Q_EMIT rightIconClicked();
    });
    connect(m_slider, &CustomSlider::valueChanged, this, [ this ](int value) {

    });

    connect(m_model, &BrightnessModel::brightnessChanged, this, &BrightnessWidget::onUpdateBright);
}

void BrightnessWidget::onUpdateBright()
{

}
