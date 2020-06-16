#ifndef HOLDCONTAINER_H
#define HOLDCONTAINER_H

#include "abstractcontainer.h"

class HoldContainer : public AbstractContainer
{
    Q_OBJECT
public:
    explicit HoldContainer(TrayPlugin *trayPlugin, QWidget *parent = nullptr);

public:
    bool acceptWrapper(FashionTrayWidgetWrapper *wrapper) override;
    void addWrapper(FashionTrayWidgetWrapper *wrapper) override;
    void refreshVisible() override;
    void setDockPosition(const Dock::Position pos) override;
//    QSize totalSize() const override;

private:
    QBoxLayout *m_mainBoxLayout;
};

#endif // HOLDCONTAINER_H
