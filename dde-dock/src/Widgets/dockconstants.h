#ifndef DOCKCONSTANTS_H
#define DOCKCONSTANTS_H

#include <QObject>
#include <QStringList>

class DockConstants : public QObject
{
    Q_OBJECT
public:
    static DockConstants * getInstants();

    enum DockMode {
        FashionMode,
        EfficientMode,
        ClassicMode
    };

    DockMode getDockMode();
    void setDockMode(DockMode value);

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
    explicit DockConstants(QObject *parent = 0);

private:
    static DockConstants * dockConstants;

    DockMode m_currentMode = DockConstants::FashionMode;

};

#endif // DOCKCONSTANTS_H
