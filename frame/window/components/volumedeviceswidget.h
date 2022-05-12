#ifndef VOLUMEDEVICESWIDGET_H
#define VOLUMEDEVICESWIDGET_H

#include <DStyledItemDelegate>

#include <QWidget>

namespace Dtk { namespace Widget { class DListView; } }

using namespace Dtk::Widget;

class CustomSlider;
class QStandardItemModel;
class QLabel;
class VolumeModel;
class AudioPorts;
class AudioSink;
class SettingDelegate;

class VolumeDevicesWidget : public QWidget
{
    Q_OBJECT

public:
    explicit VolumeDevicesWidget(VolumeModel *model, QWidget *parent = nullptr);
    ~VolumeDevicesWidget() override;

private:
    void initUi();
    void reloadAudioDevices();
    void initConnection();
    QString leftIcon();
    QString rightIcon();
    const QString soundIconFile(AudioPorts *port) const;

    void resizeHeight();

    void resetVolumeInfo();

private:
    CustomSlider *m_volumeSlider;
    QLabel *m_descriptionLabel;
    DListView *m_deviceList;
    VolumeModel *m_volumeModel;
    AudioSink *m_audioSink;
    QStandardItemModel *m_model;
    SettingDelegate *m_delegate;
};

#endif // VOLUMEDEVICESWIDGET_H
