#ifndef DOCKCONSTANTS_H
#define DOCKCONSTANTS_H

#include <QObject>
#include <QStringList>

struct DockItemData {
    QString appTitle;
    QString appIconPath;
    QString appExePath;
    bool appActived;
    QStringList appPreviews;
};

class DockConstants : public QObject
{
    Q_OBJECT
    Q_PROPERTY(int iconSize READ getIconSize WRITE setIconSize)
public:
    static DockConstants * getInstants();

    enum DockMode {
        FashionMode,
        EfficientMode,
        ClassicMode
    };

    int getIconSize();
    void setIconSize(int value);

private:
    explicit DockConstants(QObject *parent = 0);

private:
    static DockConstants * dockConstants;

    DockMode currentMode = DockConstants::FashionMode;
    int iconSize = 42;
};

#endif // DOCKCONSTANTS_H
