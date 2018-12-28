#ifndef HOLDCONTAINER_H
#define HOLDCONTAINER_H

#include "abstractcontainer.h"

class HoldContainer : public AbstractContainer
{
    Q_OBJECT
public:
    explicit HoldContainer(TrayPlugin *trayPlugin, QWidget *parent = nullptr);

public:
    bool acceptWrapper(FashionTrayWidgetWrapper *wrapper) Q_DECL_OVERRIDE;
    void addWrapper(FashionTrayWidgetWrapper *wrapper) Q_DECL_OVERRIDE;
    void refreshVisible() Q_DECL_OVERRIDE;
    void setDockPosition(const Dock::Position pos) Q_DECL_OVERRIDE;
    void setExpand(const bool expand) Q_DECL_OVERRIDE;
    QSize totalSize() const Q_DECL_OVERRIDE;

protected:
    void resizeEvent(QResizeEvent *event) Q_DECL_OVERRIDE;

private:
    QBoxLayout *m_mainBoxLayout;
    QLabel *m_holdSpliter;
};

#endif // HOLDCONTAINER_H
