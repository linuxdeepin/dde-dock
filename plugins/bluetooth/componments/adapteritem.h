#ifndef ADAPTERITEM_H
#define ADAPTERITEM_H

#include <QScrollArea>
#include <QMap>
#include <QVBoxLayout>
#include <QLabel>

class HorizontalSeparator;
class Adapter;
class SwitchItem;
class DeviceItem;
class Device;
class AdaptersManager;
class MenueItem;
class AdapterItem : public QScrollArea
{
    Q_OBJECT
public:
    explicit AdapterItem(AdaptersManager *a, Adapter *adapter, QWidget *parent = nullptr);
    int pairedDeviceCount();
    int deviceCount();
    void setPowered(bool powered);

signals:
    void deviceStateChanged(int state);
    void powerChanged(bool powered);
    void sizeChange();

private slots:
    void deviceItemPaired(const bool paired);
    void removeDeviceItem(const Device *device);
    void showAndConnect(bool change);
    void addDeviceItem(const Device *constDevice);

private:
    void createDeviceItem(Device *device);
    void updateView();
    void showDevices(bool change);

private:
    QWidget *m_centralWidget;
    HorizontalSeparator *m_line;
    QLabel *m_devGoupName;
    QVBoxLayout *m_deviceLayout;
    MenueItem *m_openControlCenter;

    AdaptersManager *m_adaptersManager;

    Adapter *m_adapter;
    SwitchItem *m_switchItem;
    QMap<QString, DeviceItem*> m_deviceItems;
    QMap<QString, DeviceItem*> m_pairedDeviceItems;
};

#endif // ADAPTERITEM_H
