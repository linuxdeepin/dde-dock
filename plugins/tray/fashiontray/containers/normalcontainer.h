#ifndef NORMALCONTAINER_H
#define NORMALCONTAINER_H

#include "abstractcontainer.h"

class NormalContainer : public AbstractContainer
{
    Q_OBJECT
public:
    explicit NormalContainer(TrayPlugin *trayPlugin, QWidget *parent = nullptr);

    // AbstractContainer interface
public:
    bool acceptWrapper(FashionTrayWidgetWrapper *wrapper) Q_DECL_OVERRIDE;
    void addWrapper(FashionTrayWidgetWrapper *wrapper) Q_DECL_OVERRIDE;
    void refreshVisible() Q_DECL_OVERRIDE;
    void setExpand(const bool expand) Q_DECL_OVERRIDE;

protected:
    int whereToInsert(FashionTrayWidgetWrapper *wrapper) Q_DECL_OVERRIDE;

private:
    int whereToInsertByDefault(FashionTrayWidgetWrapper *wrapper) const;
    int whereToInsertAppTrayByDefault(FashionTrayWidgetWrapper *wrapper) const;
    int whereToInsertSystemTrayByDefault(FashionTrayWidgetWrapper *wrapper) const;
};

#endif // NORMALCONTAINER_H
