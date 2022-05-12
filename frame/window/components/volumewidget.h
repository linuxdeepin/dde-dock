#ifndef VOLUMEWIDGET_H
#define VOLUMEWIDGET_H

#include <DBlurEffectWidget>
#include <QWidget>

class VolumeModel;
class QDBusMessage;
class CustomSlider;
class QLabel;
class AudioSink;

DWIDGET_USE_NAMESPACE

class VolumeWidget : public DBlurEffectWidget
{
    Q_OBJECT

public:
    explicit VolumeWidget(QWidget *parent = nullptr);
    ~VolumeWidget() override;
    VolumeModel *model();

Q_SIGNALS:
    void visibleChanged(bool);
    void rightIconClick();

protected:
    void initUi();
    void initConnection();

    void showEvent(QShowEvent *event) override;
    void hideEvent(QHideEvent *event) override;

private:
    const QString leftIcon();
    const QString rightIcon();

private:
    VolumeModel *m_volumeController;
    CustomSlider *m_volumnCtrl;
    AudioSink *m_defaultSink;
};

#endif // VOLUMEWIDGET_H
