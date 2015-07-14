#ifndef DOCKMODEDATA_H
#define DOCKMODEDATA_H

#include <QObject>
#include <QStringList>
#include <QDebug>
#include "DBus/dbusdocksetting.h"
#include "dockconstants.h"

class DockModeData : public QObject
{
    Q_OBJECT
public:
    static DockModeData * instance();

    Dock::DockMode getDockMode();
    void setDockMode(Dock::DockMode value);
    Dock::HideMode getHideMode();
    void setHideMode(Dock::HideMode value);

    int getDockHeight();
    int getItemHeight();
    int getNormalItemWidth();
    int getActivedItemWidth();
    int getAppItemSpacing();
    int getAppIconSize();
    int getAppletsItemHeight();
    int getAppletsItemWidth();
    int getAppletsItemSpacing();
    int getAppletsIconSize();

signals:
    void dockModeChanged(Dock::DockMode newMode,Dock::DockMode oldMode);
    void hideModeChanged(Dock::HideMode newMode,Dock::HideMode oldMode);

private slots:
    void slotDockModeChanged(int mode);
    void slotHideModeChanged(int mode);

private:
    explicit DockModeData(QObject *parent = 0);

    void initDDS();
private:
    static DockModeData * dockModeData;

    Dock::DockMode m_currentMode = Dock::EfficientMode;
    Dock::HideMode m_hideMode = Dock::KeepShowing;

    DBusDockSetting *m_dds = NULL;
};

#endif // DOCKMODEDATA_H
