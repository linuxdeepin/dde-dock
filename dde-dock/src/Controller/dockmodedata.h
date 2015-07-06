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

    Dock::DockMode getDockMode();
    void setDockMode(Dock::DockMode value);

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

private:
    explicit DockModeData(QObject *parent = 0);

private:
    static DockModeData * dockModeData;

    Dock::DockMode m_currentMode = Dock::EfficientMode;

};

#endif // DOCKMODEDATA_H
