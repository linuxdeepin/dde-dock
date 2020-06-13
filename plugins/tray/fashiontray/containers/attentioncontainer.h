#ifndef ATTENTIONCONTAINER_H
#define ATTENTIONCONTAINER_H

#include "abstractcontainer.h"

class AttentionContainer : public AbstractContainer
{
    Q_OBJECT
public:
    explicit AttentionContainer(TrayPlugin *trayPlugin, QWidget *parent = nullptr);

    FashionTrayWidgetWrapper *takeAttentionWrapper();

    // AbstractContainer interface
public:
    bool acceptWrapper(FashionTrayWidgetWrapper *wrapper) override;
    void refreshVisible() override;
    void addWrapper(FashionTrayWidgetWrapper *wrapper) override;
};

#endif // ATTENTIONCONTAINER_H
