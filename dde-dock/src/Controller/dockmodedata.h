#ifndef DOCKMODEDATA_H
#define DOCKMODEDATA_H

#include <QObject>
#include <QStringList>
#include "dockconstants.h"

class DockModeData : public QObject
{
    Q_OBJECT
public:
    static DockModeData * getInstants();

    DockConstants::DockMode getDockMode();
    void setDockMode(DockConstants::DockMode value);

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
    void dockModeChanged(DockConstants::DockMode newMode,DockConstants::DockMode oldMode);

private:
    explicit DockModeData(QObject *parent = 0);

private:
    static DockModeData * dockModeData;

    DockConstants::DockMode m_currentMode = DockConstants::EfficientMode;

};

#endif // DOCKMODEDATA_H
